// Package items reads a nav list from JSON.
// A source is one array, or one object per line.
package items

import (
	"bytes"
	"encoding/json"
	"strings"

	"high-tui/internal/model"
)

var kinds = map[string]bool{
	"load": true, "edit": true, "window": true, "return": true,
	"save": true, "refresh": true, "run": true, "set": true,
}

type rawAction struct {
	Kind   string `json:"kind"`
	Target string `json:"target"`
	File   string `json:"file"`
	Name   string `json:"name"`
	Href   string `json:"href"`
	Run    string `json:"run"`
	Output string `json:"output"`
	Value  string `json:"value"`
}

type rawItem struct {
	Label string          `json:"label"`
	Value string          `json:"value"`
	File  string          `json:"file"`
	Name  string          `json:"name"`
	Next  json.RawMessage `json:"next"`
	Focus *rawAction      `json:"focus"`
	Enter *rawAction      `json:"enter"`
}

// Parse reads a JSON array or one JSON object per line.
func Parse(text, where string) ([]model.Item, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, model.Errf("JSON invalide (%s)", where)
	}
	if text[0] == '[' {
		return parseArray([]byte(text), where)
	}
	var items []model.Item
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		item, err := parseOne([]byte(line), where)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		return nil, model.Errf("JSON invalide (%s)", where)
	}
	return items, nil
}

func parseArray(raw []byte, where string) ([]model.Item, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	var rows []rawItem
	if err := dec.Decode(&rows); err != nil {
		return nil, jsonErr(err, where)
	}
	if dec.More() {
		return nil, model.Errf("JSON invalide (%s)", where)
	}
	return convertAll(rows, where)
}

func parseOne(raw []byte, where string) (model.Item, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	var row rawItem
	if err := dec.Decode(&row); err != nil {
		return model.Item{}, jsonErr(err, where)
	}
	if dec.More() {
		return model.Item{}, model.Errf("JSON invalide (%s)", where)
	}
	return convert(row, where)
}

func convertAll(rows []rawItem, where string) ([]model.Item, error) {
	items := make([]model.Item, 0, len(rows))
	for _, row := range rows {
		item, err := convert(row, where)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func convert(row rawItem, where string) (model.Item, error) {
	if strings.TrimSpace(row.Label) == "" {
		return model.Item{}, model.Errf("label manquant (%s)", where)
	}
	next, nested, err := parseNext(row.Next, where)
	if err != nil {
		return model.Item{}, err
	}
	focus, err := convertAction(row.Focus, where)
	if err != nil {
		return model.Item{}, err
	}
	enter, err := convertAction(row.Enter, where)
	if err != nil {
		return model.Item{}, err
	}
	return model.Item{
		Label: row.Label, Value: row.Value, File: row.File, Name: row.Name,
		Next: next, NextItems: nested, Focus: focus, Enter: enter,
	}, nil
}

func parseNext(raw json.RawMessage, where string) (string, []model.Item, error) {
	if len(bytes.TrimSpace(raw)) == 0 || string(raw) == "null" {
		return "", nil, nil
	}
	if raw[0] == '"' {
		var source string
		if err := json.Unmarshal(raw, &source); err != nil {
			return "", nil, model.Errf("JSON invalide (%s)", where)
		}
		if err := CheckSource(source, where); err != nil {
			return "", nil, err
		}
		return source, nil, nil
	}
	if raw[0] == '[' {
		items, err := parseArray(raw, where)
		return "", items, err
	}
	return "", nil, model.Errf("JSON invalide (%s)", where)
}

// CheckSource accepts a file path or a cmd: command.
func CheckSource(source, where string) error {
	if source == "" {
		return model.Errf("source manquante (%s)", where)
	}
	if strings.HasPrefix(source, "cmd:") {
		script := strings.TrimPrefix(source, "cmd:")
		if script == "" || strings.ContainsAny(script, "\t\n") {
			return model.Errf("commande illégale (%s)", where)
		}
	}
	return nil
}

func convertAction(raw *rawAction, where string) (model.Action, error) {
	if raw == nil {
		return model.Action{}, nil
	}
	act := model.Action{
		Kind: raw.Kind, Target: raw.Target, File: raw.File, Name: raw.Name,
		Href: raw.Href, Run: raw.Run, Output: raw.Output, Value: raw.Value,
	}
	if act.Kind == "" {
		return act, nil
	}
	if !kinds[act.Kind] {
		return act, model.Errf("action inconnue « %s » (%s)", act.Kind, where)
	}
	if strings.ContainsAny(act.Run, "\t\n") {
		return act, model.Errf("commande illégale (%s)", where)
	}
	switch act.Kind {
	case "load", "edit":
		if act.Target == "" {
			return act, model.Errf("cible manquante (%s)", where)
		}
	case "window":
		if act.Href == "" {
			return act, model.Errf("href manquant (%s)", where)
		}
	case "run":
		if act.Run == "" {
			return act, model.Errf("commande manquante (%s)", where)
		}
	}
	return act, nil
}

func jsonErr(err error, where string) error {
	if strings.Contains(err.Error(), "unknown field") {
		return model.Errf("champ inconnu (%s)", where)
	}
	return model.Errf("JSON invalide (%s)", where)
}
