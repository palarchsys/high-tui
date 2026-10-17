// Package html parses one window document.
package html

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"high-tui/internal/model"
)

var zoneTags = map[string]bool{
	"header": true, "nav": true, "select": true, "input": true,
	"textarea": true, "button": true, "command": true, "status": true,
}

var rowTags = map[string]bool{"header": true, "main": true, "status": true}

var actions = map[string]bool{
	"load": true, "edit": true, "window": true, "return": true,
	"save": true, "refresh": true, "run": true, "set": true,
}

var attrs = map[string]bool{
	"id": true, "class": true, "label": true, "value": true, "fr": true, "lines": true,
	"file": true, "name": true, "save": true, "readonly": true, "href": true, "run": true,
	"output": true, "target": true, "focus": true, "enter": true,
	"focus-target": true, "focus-file": true, "focus-name": true, "focus-href": true,
	"focus-run": true, "focus-output": true, "focus-value": true,
	"enter-target": true, "enter-file": true, "enter-name": true, "enter-href": true,
	"enter-run": true, "enter-output": true, "enter-value": true,
}

// Node is a raw element before validation.
type Node struct {
	Tag   string
	Attr  map[string]string
	Kids  []Node
	Where string
}

// Parse reads a single window document.
func Parse(text, where string) (Node, error) {
	p := &parser{s: text, where: where}
	var roots []Node
	for {
		p.skipSpaceAndComments()
		if p.i >= len(p.s) {
			break
		}
		if !strings.HasPrefix(p.s[p.i:], "<") {
			return Node{}, model.Errf("texte inattendu (%s)", where)
		}
		node, err := p.element()
		if err != nil {
			return Node{}, err
		}
		roots = append(roots, node)
	}
	if len(roots) != 1 || roots[0].Tag != "window" {
		return Node{}, model.Errf("balise window attendue (%s)", where)
	}
	return roots[0], nil
}

type parser struct {
	s, where string
	i        int
}

func (p *parser) skipSpaceAndComments() {
	for p.i < len(p.s) {
		if strings.HasPrefix(p.s[p.i:], "<!--") {
			j := strings.Index(p.s[p.i+4:], "-->")
			if j < 0 {
				p.i = len(p.s)
				return
			}
			p.i += 4 + j + 3
			continue
		}
		r, w := utf8.DecodeRuneInString(p.s[p.i:])
		if !unicode.IsSpace(r) {
			return
		}
		p.i += w
	}
}

func (p *parser) element() (Node, error) {
	if p.s[p.i] != '<' {
		return Node{}, model.Errf("balise attendue (%s)", p.where)
	}
	p.i++
	if p.i < len(p.s) && p.s[p.i] == '/' {
		return Node{}, model.Errf("balise fermante inattendue (%s)", p.where)
	}
	tag := p.readName()
	if tag == "" {
		return Node{}, model.Errf("balise attendue (%s)", p.where)
	}
	attr := map[string]string{}
	for {
		p.skipSpace()
		if p.i >= len(p.s) {
			return Node{}, model.Errf("balise non fermée « %s » (%s)", tag, p.where)
		}
		if strings.HasPrefix(p.s[p.i:], "/>") {
			p.i += 2
			return Node{Tag: tag, Attr: attr, Where: p.where}, nil
		}
		if p.s[p.i] == '>' {
			p.i++
			break
		}
		name := p.readName()
		if name == "" || !attrs[name] {
			shown := name
			if shown == "" {
				shown = "?"
			}
			return Node{}, model.Errf("attribut inconnu « %s » (%s)", shown, p.where)
		}
		p.skipSpace()
		if p.i >= len(p.s) || p.s[p.i] != '=' {
			return Node{}, model.Errf("valeur manquante « %s » (%s)", name, p.where)
		}
		p.i++
		p.skipSpace()
		val, err := p.readValue()
		if err != nil {
			return Node{}, err
		}
		if val == "" {
			return Node{}, model.Errf("attribut vide « %s » (%s)", name, p.where)
		}
		if _, ok := attr[name]; ok {
			return Node{}, model.Errf("attribut en double « %s » (%s)", name, p.where)
		}
		attr[name] = val
	}
	var kids []Node
	for {
		p.skipSpaceAndComments()
		if p.i >= len(p.s) {
			return Node{}, model.Errf("balise non fermée « %s » (%s)", tag, p.where)
		}
		if strings.HasPrefix(p.s[p.i:], "</") {
			p.i += 2
			end := p.readName()
			p.skipSpace()
			if p.i >= len(p.s) || p.s[p.i] != '>' || end != tag {
				return Node{}, model.Errf("fermeture « %s » invalide (%s)", tag, p.where)
			}
			p.i++
			return Node{Tag: tag, Attr: attr, Kids: kids, Where: p.where}, nil
		}
		if p.s[p.i] != '<' {
			return Node{}, model.Errf("texte inattendu (%s)", p.where)
		}
		kid, err := p.element()
		if err != nil {
			return Node{}, err
		}
		kids = append(kids, kid)
	}
}

func (p *parser) skipSpace() {
	for p.i < len(p.s) {
		r, w := utf8.DecodeRuneInString(p.s[p.i:])
		if !unicode.IsSpace(r) {
			return
		}
		p.i += w
	}
}

func (p *parser) readName() string {
	start := p.i
	for p.i < len(p.s) {
		c := p.s[p.i]
		if (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-' {
			p.i++
			continue
		}
		break
	}
	return p.s[start:p.i]
}

func (p *parser) readValue() (string, error) {
	if p.i >= len(p.s) || (p.s[p.i] != '"' && p.s[p.i] != '\'') {
		return "", model.Errf("valeur non citée (%s)", p.where)
	}
	q := p.s[p.i]
	p.i++
	start := p.i
	for p.i < len(p.s) && p.s[p.i] != q {
		if p.s[p.i] == '\n' || p.s[p.i] == '\t' {
			return "", model.Errf("valeur illégale (%s)", p.where)
		}
		p.i++
	}
	if p.i >= len(p.s) {
		return "", model.Errf("guillemet non fermé (%s)", p.where)
	}
	val := p.s[start:p.i]
	p.i++
	return val, nil
}

func positive(raw, where string) (int, error) {
	if raw == "" {
		return 0, model.Errf("nombre invalide (%s)", where)
	}
	n := 0
	for _, c := range raw {
		if c < '0' || c > '9' {
			return 0, model.Errf("nombre invalide « %s » (%s)", raw, where)
		}
		n = n*10 + int(c-'0')
	}
	if n < 1 {
		return 0, model.Errf("nombre invalide « %s » (%s)", raw, where)
	}
	return n, nil
}

func actionOf(attr map[string]string, prefix, where string) (model.Action, error) {
	kind := attr[prefix]
	act := model.Action{
		Kind:   kind,
		Target: attr[prefix+"-target"],
		File:   attr[prefix+"-file"],
		Name:   attr[prefix+"-name"],
		Href:   attr[prefix+"-href"],
		Run:    attr[prefix+"-run"],
		Output: attr[prefix+"-output"],
		Value:  attr[prefix+"-value"],
	}
	if kind == "" {
		return act, nil
	}
	if !actions[kind] {
		return act, model.Errf("action inconnue « %s » (%s)", kind, where)
	}
	if strings.ContainsAny(act.Run, "\t\n") {
		return act, model.Errf("commande illégale (%s)", where)
	}
	switch kind {
	case "load":
		if act.Target == "" {
			return act, model.Errf("cible manquante (%s)", where)
		}
		if act.File == "" {
			return act, model.Errf("fichier manquant (%s)", where)
		}
	case "edit":
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
	case "set":
		if act.File == "" {
			return act, model.Errf("fichier manquant (%s)", where)
		}
		if act.Name == "" {
			return act, model.Errf("variable manquante (%s)", where)
		}
		if act.Value == "" {
			return act, model.Errf("valeur manquante (%s)", where)
		}
	}
	return act, nil
}

func classesOf(raw string) []string {
	if raw == "" {
		return nil
	}
	return strings.Fields(raw)
}

// Window converts a parsed document into a model window. File checks and
// target checks are done by the builder, which knows the directory.
func Window(root Node, file string) (model.Window, error) {
	id := root.Attr["id"]
	if id == "" {
		return model.Window{}, model.Errf("attribut manquant « id » (%s)", file)
	}
	win := model.Window{
		ID: id, Label: root.Attr["label"], File: file, Classes: classesOf(root.Attr["class"]),
	}
	ids := map[string]bool{id: true}
	if len(root.Kids) == 0 {
		return model.Window{}, model.Errf("fenêtre sans rangée (%s)", file)
	}
	for _, rowNode := range root.Kids {
		if !rowTags[rowNode.Tag] {
			return model.Window{}, model.Errf("balise inconnue « %s » (%s)", rowNode.Tag, id)
		}
		row, err := convertRow(rowNode, ids)
		if err != nil {
			return model.Window{}, err
		}
		win.Rows = append(win.Rows, row)
	}
	if err := checkTargets(&win); err != nil {
		return model.Window{}, err
	}
	ensureReturn(&win, ids)
	return win, nil
}

func convertRow(node Node, ids map[string]bool) (model.Row, error) {
	id := node.Attr["id"]
	if id == "" {
		return model.Row{}, model.Errf("attribut manquant « id » (%s)", node.Tag)
	}
	if ids[id] {
		return model.Row{}, model.Errf("identifiant en double : %s", id)
	}
	ids[id] = true
	lines := 1
	if raw, ok := node.Attr["lines"]; ok {
		n, err := positive(raw, id)
		if err != nil {
			return model.Row{}, err
		}
		lines = n
	}
	row := model.Row{ID: id, Tag: node.Tag, Classes: classesOf(node.Attr["class"]), Lines: lines}
	if len(node.Kids) == 0 {
		return model.Row{}, model.Errf("rangée sans colonne (%s)", id)
	}
	for _, kid := range node.Kids {
		zone, err := convertZone(kid, ids)
		if err != nil {
			return model.Row{}, err
		}
		row.Zones = append(row.Zones, zone)
	}
	return row, nil
}

func convertZone(node Node, ids map[string]bool) (model.Zone, error) {
	if !zoneTags[node.Tag] {
		return model.Zone{}, model.Errf("balise inconnue « %s » (%s)", node.Tag, node.Attr["id"])
	}
	id := node.Attr["id"]
	if id == "" {
		return model.Zone{}, model.Errf("attribut manquant « id » (%s)", node.Tag)
	}
	if ids[id] {
		return model.Zone{}, model.Errf("identifiant en double : %s", id)
	}
	ids[id] = true
	fr := 1
	if raw, ok := node.Attr["fr"]; ok {
		n, err := positive(raw, id)
		if err != nil {
			return model.Zone{}, err
		}
		fr = n
	}
	readonly := false
	if raw, ok := node.Attr["readonly"]; ok {
		if raw != "true" && raw != "false" {
			return model.Zone{}, model.Errf("readonly invalide (%s)", id)
		}
		readonly = raw == "true"
	}
	if raw, ok := node.Attr["save"]; ok && raw != "enter" && raw != "button" {
		return model.Zone{}, model.Errf("save invalide (%s)", id)
	}
	focus, err := actionOf(node.Attr, "focus", id)
	if err != nil {
		return model.Zone{}, err
	}
	enter, err := actionOf(node.Attr, "enter", id)
	if err != nil {
		return model.Zone{}, err
	}
	zone := model.Zone{
		Tag: node.Tag, ID: id, Classes: classesOf(node.Attr["class"]),
		Label: node.Attr["label"], Value: node.Attr["value"], File: node.Attr["file"],
		Name: node.Attr["name"], Save: node.Attr["save"], Href: node.Attr["href"],
		Run: node.Attr["run"], Output: node.Attr["output"], Fr: fr, Readonly: readonly,
		Focus: focus, Enter: enter,
	}
	if node.Tag == "nav" || node.Tag == "select" {
		if zone.Label == "" {
			return model.Zone{}, model.Errf("label manquant (%s)", id)
		}
		for _, kid := range node.Kids {
			item, err := convertItem(kid, ids, zone.Value)
			if err != nil {
				return model.Zone{}, err
			}
			zone.Items = append(zone.Items, item)
		}
	} else if len(node.Kids) > 0 {
		return model.Zone{}, model.Errf("enfant refusé (%s)", id)
	}
	if node.Tag == "header" && zone.Label == "" {
		return model.Zone{}, model.Errf("label manquant (%s)", id)
	}
	if node.Tag == "button" || node.Tag == "command" {
		if zone.Label == "" {
			return model.Zone{}, model.Errf("label manquant (%s)", id)
		}
	}
	if node.Tag == "button" && zone.Enter.Kind == "" {
		return model.Zone{}, model.Errf("action manquante (%s)", id)
	}
	if node.Tag == "command" && zone.Enter.Kind != "run" {
		return model.Zone{}, model.Errf("commande manquante (%s)", id)
	}
	return zone, nil
}

func convertItem(node Node, ids map[string]bool, fallback string) (model.Item, error) {
	if node.Tag != "item" {
		return model.Item{}, model.Errf("balise inconnue « %s » (%s)", node.Tag, node.Attr["id"])
	}
	if len(node.Kids) > 0 {
		return model.Item{}, model.Errf("enfant refusé (%s)", node.Attr["id"])
	}
	id := node.Attr["id"]
	if id == "" {
		return model.Item{}, model.Errf("attribut manquant « id » (item)")
	}
	if ids[id] {
		return model.Item{}, model.Errf("identifiant en double : %s", id)
	}
	ids[id] = true
	if node.Attr["label"] == "" {
		return model.Item{}, model.Errf("label manquant (%s)", id)
	}
	focus, err := actionOf(node.Attr, "focus", id)
	if err != nil {
		return model.Item{}, err
	}
	enter, err := actionOf(node.Attr, "enter", id)
	if err != nil {
		return model.Item{}, err
	}
	value := node.Attr["value"]
	if enter.Kind == "set" && enter.Value == "" {
		enter.Value = value
		if enter.Value == "" {
			enter.Value = fallback
		}
		if enter.Value == "" {
			return model.Item{}, model.Errf("valeur manquante (%s)", id)
		}
	}
	return model.Item{ID: id, Label: node.Attr["label"], Value: value, Focus: focus, Enter: enter}, nil
}

func zonesOf(win *model.Window) []model.Zone {
	var out []model.Zone
	for _, row := range win.Rows {
		out = append(out, row.Zones...)
	}
	return out
}

func checkTargets(win *model.Window) error {
	byID := map[string]model.Zone{}
	for _, zone := range zonesOf(win) {
		byID[zone.ID] = zone
	}
	one := func(act model.Action, where string) error {
		switch act.Kind {
		case "load", "edit":
			zone, ok := byID[act.Target]
			if !ok {
				return model.Errf("cible absente : %s", act.Target)
			}
			if zone.Tag != "input" && zone.Tag != "textarea" {
				return model.Errf("cible incompatible : %s", act.Target)
			}
			if act.Kind == "edit" && zone.Tag == "textarea" && zone.Readonly {
				return model.Errf("cible non éditable : %s", act.Target)
			}
		case "set":
			if act.Target == "" {
				return nil
			}
			zone, ok := byID[act.Target]
			if !ok || (zone.Tag != "input" && zone.Tag != "textarea") {
				return model.Errf("cible incompatible : %s", act.Target)
			}
		case "run":
			if act.Output == "" {
				return nil
			}
			zone, ok := byID[act.Output]
			if !ok || zone.Tag != "textarea" {
				return model.Errf("cible absente : %s", act.Output)
			}
		}
		return nil
	}
	for _, zone := range zonesOf(win) {
		if err := one(zone.Focus, zone.ID); err != nil {
			return err
		}
		if err := one(zone.Enter, zone.ID); err != nil {
			return err
		}
		for _, item := range zone.Items {
			if err := one(item.Focus, item.ID); err != nil {
				return err
			}
			if err := one(item.Enter, item.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

func ensureReturn(win *model.Window, ids map[string]bool) {
	for _, zone := range zonesOf(win) {
		if zone.Focus.Kind == "return" || zone.Enter.Kind == "return" {
			return
		}
		for _, item := range zone.Items {
			if item.Focus.Kind == "return" || item.Enter.Kind == "return" {
				return
			}
		}
	}
	for ri := range win.Rows {
		for zi := range win.Rows[ri].Zones {
			zone := &win.Rows[ri].Zones[zi]
			if zone.Tag != "nav" && zone.Tag != "select" {
				continue
			}
			id := zone.ID + "-retour"
			for ids[id] {
				id += "-r"
			}
			ids[id] = true
			zone.Items = append(zone.Items, model.Item{
				ID: id, Label: "Retour", Enter: model.Action{Kind: "return"},
			})
			return
		}
	}
}
