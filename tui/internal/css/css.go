// Package css parses the window stylesheet subset.
package css

import (
	"sort"
	"strings"
	"unicode"

	"high-tui/internal/model"
)

var known = map[string]bool{
	"color": true, "background": true, "background-color": true,
	"border": true, "border-width": true, "border-style": true, "border-color": true,
	"border-radius": true, "padding": true, "margin": true, "max-width": true,
	"column-gap": true, "line-height": true, "font-weight": true,
}

var tags = map[string]bool{
	"window": true, "header": true, "main": true, "nav": true, "select": true,
	"input": true, "textarea": true, "button": true, "command": true, "status": true,
}

// Rule is one selector plus its declarations, in source order.
type Rule struct {
	Tag, Class, ID string
	Focus          bool
	Spec           [3]int
	Decls          [][2]string
	Order          int
}

// Parse reads a window CSS file. where is the file name used in errors.
func Parse(text, where string) ([]Rule, error) {
	raw, err := stripComments(text)
	if err != nil {
		return nil, err
	}
	var rules []Rule
	i := 0
	order := 0
	for {
		for i < len(raw) && unicode.IsSpace(rune(raw[i])) {
			i++
		}
		if i >= len(raw) {
			break
		}
		j := strings.IndexByte(raw[i:], '{')
		if j < 0 {
			return nil, model.Errf("accolade ouvrante manquante dans le CSS")
		}
		j += i
		selText := strings.TrimSpace(raw[i:j])
		k := strings.IndexByte(raw[j:], '}')
		if k < 0 {
			return nil, model.Errf("accolade fermante manquante dans le CSS")
		}
		k += j
		body := raw[j+1 : k]
		i = k + 1
		if strings.HasPrefix(selText, "@") {
			name := selText
			if f := strings.Fields(selText); len(f) > 0 {
				name = f[0]
			}
			return nil, model.Errf("règle @ refusée (%s)", name)
		}
		decls, err := parseDecls(body, where)
		if err != nil {
			return nil, err
		}
		for _, part := range strings.Split(selText, ",") {
			sel, err := parseSelector(part, where)
			if err != nil {
				return nil, err
			}
			sel.Decls = decls
			sel.Order = order
			order++
			rules = append(rules, sel)
		}
	}
	return rules, nil
}

func stripComments(text string) (string, error) {
	var b strings.Builder
	for i := 0; i < len(text); {
		if strings.HasPrefix(text[i:], "/*") {
			j := strings.Index(text[i+2:], "*/")
			if j < 0 {
				return "", model.Errf("commentaire CSS non fermé")
			}
			i += 2 + j + 2
			continue
		}
		b.WriteByte(text[i])
		i++
	}
	return b.String(), nil
}

func parseSelector(part, where string) (Rule, error) {
	text := strings.TrimSpace(part)
	if text == "" {
		return Rule{}, model.Errf("sélecteur vide (%s)", where)
	}
	var sel Rule
	i := 0
	if i < len(text) && (unicode.IsLetter(rune(text[i])) || text[i] == '_') {
		j := i + 1
		for j < len(text) && (unicode.IsLetter(rune(text[j])) || unicode.IsDigit(rune(text[j])) || text[j] == '-' || text[j] == '_') {
			j++
		}
		sel.Tag = text[i:j]
		i = j
	}
	if i < len(text) && text[i] == '.' {
		j := i + 1
		for j < len(text) && (unicode.IsLetter(rune(text[j])) || unicode.IsDigit(rune(text[j])) || text[j] == '-' || text[j] == '_') {
			j++
		}
		if j == i+1 {
			return Rule{}, model.Errf("classe vide (%s)", where)
		}
		sel.Class = text[i+1 : j]
		i = j
	}
	if i < len(text) && text[i] == '#' {
		j := i + 1
		for j < len(text) && (unicode.IsLetter(rune(text[j])) || unicode.IsDigit(rune(text[j])) || text[j] == '-' || text[j] == '_') {
			j++
		}
		if j == i+1 {
			return Rule{}, model.Errf("identifiant CSS vide (%s)", where)
		}
		sel.ID = text[i+1 : j]
		i = j
	}
	if strings.HasPrefix(text[i:], ":focus") {
		sel.Focus = true
		i += len(":focus")
	}
	if i != len(text) || (sel.Tag == "" && sel.Class == "" && sel.ID == "" && !sel.Focus) {
		return Rule{}, model.Errf("sélecteur refusé « %s » (%s)", text, where)
	}
	if sel.Tag != "" && !tags[sel.Tag] {
		return Rule{}, model.Errf("balise inconnue « %s » (%s)", sel.Tag, where)
	}
	idBit, classBit, tagBit := 0, 0, 0
	if sel.ID != "" {
		idBit = 1
	}
	if sel.Class != "" {
		classBit++
	}
	if sel.Focus {
		classBit++
	}
	if sel.Tag != "" {
		tagBit = 1
	}
	sel.Spec = [3]int{idBit, classBit, tagBit}
	return sel, nil
}

func parseDecls(body, where string) ([][2]string, error) {
	var props [][2]string
	for _, piece := range strings.Split(body, ";") {
		piece = strings.TrimSpace(piece)
		if piece == "" {
			continue
		}
		name, val, ok := strings.Cut(piece, ":")
		if !ok {
			return nil, model.Errf("déclaration invalide « %s » (%s)", piece, where)
		}
		name = strings.ToLower(strings.TrimSpace(name))
		val = strings.Join(strings.Fields(val), " ")
		if !known[name] {
			return nil, model.Errf("propriété inconnue « %s » (%s)", name, where)
		}
		props = append(props, [2]string{name, val})
	}
	return props, nil
}

func lengthOf(token, where string) (int, error) {
	raw := strings.TrimSpace(token)
	raw = strings.TrimSuffix(raw, "ch")
	if raw == "" || strings.IndexFunc(raw, func(r rune) bool { return r < '0' || r > '9' }) >= 0 {
		return 0, model.Errf("longueur invalide « %s » (%s)", strings.TrimSpace(token), where)
	}
	n := 0
	for _, c := range raw {
		n = n*10 + int(c-'0')
	}
	return n, nil
}

func fourEdges(parts []string, where string) ([4]int, error) {
	var nums []int
	for _, p := range parts {
		n, err := lengthOf(p, where)
		if err != nil {
			return [4]int{}, err
		}
		nums = append(nums, n)
	}
	switch len(nums) {
	case 1:
		return [4]int{nums[0], nums[0], nums[0], nums[0]}, nil
	case 2:
		return [4]int{nums[0], nums[1], nums[0], nums[1]}, nil
	case 3:
		return [4]int{nums[0], nums[1], nums[2], nums[1]}, nil
	case 4:
		return [4]int{nums[0], nums[1], nums[2], nums[3]}, nil
	default:
		return [4]int{}, model.Errf("trop de valeurs (%s)", where)
	}
}

// Blank is the non-use style. MaxWidth applies to a window only.
func Blank() model.Style {
	return model.Style{LineHeight: 1, MaxWidth: 260}
}

func applyDecl(style *model.Style, name, val, where string) error {
	switch name {
	case "color":
		style.Color = val
	case "background", "background-color":
		style.Background = val
	case "border-color":
		style.BorderColor = val
	case "border-width":
		n, err := lengthOf(val, where)
		if err != nil {
			return err
		}
		style.BorderWidth = n
	case "border-style":
		if val != "solid" && val != "none" {
			return model.Errf("border-style refusé « %s » (%s)", val, where)
		}
		if val == "none" {
			style.BorderWidth = 0
		}
	case "border-radius":
		n, err := lengthOf(val, where)
		if err != nil {
			return err
		}
		style.Radius = n
	case "border":
		widthSet := false
		width, color := 0, ""
		for _, bit := range strings.Fields(val) {
			switch {
			case bit == "solid":
				continue
			case bit == "none":
				width, widthSet = 0, true
			case strings.HasPrefix(bit, "#"):
				color = bit
			default:
				n, err := lengthOf(bit, where)
				if err != nil {
					return err
				}
				width, widthSet = n, true
			}
		}
		if !widthSet || color == "" {
			return model.Errf("border incomplet « %s » (%s)", val, where)
		}
		style.BorderWidth = width
		style.BorderColor = color
	case "padding":
		edges, err := fourEdges(strings.Fields(val), where)
		if err != nil {
			return err
		}
		style.Pad = edges
	case "margin":
		edges, err := fourEdges(strings.Fields(val), where)
		if err != nil {
			return err
		}
		style.Mar = edges
	case "max-width":
		n, err := lengthOf(val, where)
		if err != nil {
			return err
		}
		style.MaxWidth = n
	case "column-gap":
		n, err := lengthOf(val, where)
		if err != nil {
			return err
		}
		style.Gap = n
	case "line-height":
		n, err := lengthOf(val, where)
		if err != nil || n < 1 || strings.Contains(val, "ch") {
			return model.Errf("line-height refusé « %s » (%s)", val, where)
		}
		style.LineHeight = n
	case "font-weight":
		if val != "bold" && val != "normal" {
			return model.Errf("font-weight refusé « %s » (%s)", val, where)
		}
		style.Bold = val == "bold"
	}
	if style.BorderWidth != 0 && style.BorderWidth != 1 {
		return model.Errf("border-width doit valoir 0 ou 1 (%s)", where)
	}
	return nil
}

func matches(rule Rule, tag string, classes []string, id string, focus bool) bool {
	if rule.Focus && !focus {
		return false
	}
	if rule.Tag != "" && rule.Tag != tag {
		return false
	}
	if rule.ID != "" && rule.ID != id {
		return false
	}
	if rule.Class != "" {
		found := false
		for _, c := range classes {
			if c == rule.Class {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

// Resolve cascades tag, class, id, then :focus. inherit supplies color and background.
func Resolve(rules []Rule, tag string, classes []string, id string, focus bool, inherit model.Style) (model.Style, error) {
	style := Blank()
	style.Color = inherit.Color
	style.Background = inherit.Background
	style.BorderColor = inherit.Color
	type hit struct {
		spec  [3]int
		order int
		decls [][2]string
	}
	var chosen []hit
	for _, rule := range rules {
		if matches(rule, tag, classes, id, focus) {
			chosen = append(chosen, hit{rule.Spec, rule.Order, rule.Decls})
		}
	}
	sort.SliceStable(chosen, func(i, j int) bool {
		a, b := chosen[i], chosen[j]
		if a.spec != b.spec {
			if a.spec[0] != b.spec[0] {
				return a.spec[0] < b.spec[0]
			}
			if a.spec[1] != b.spec[1] {
				return a.spec[1] < b.spec[1]
			}
			return a.spec[2] < b.spec[2]
		}
		return a.order < b.order
	})
	where := id
	if where == "" {
		where = tag
	}
	for _, item := range chosen {
		for _, decl := range item.decls {
			if err := applyDecl(&style, decl[0], decl[1], where); err != nil {
				return model.Style{}, err
			}
		}
	}
	if style.Color == "" || style.Background == "" {
		return model.Style{}, model.Errf("couleur manquante pour %s", where)
	}
	if style.BorderColor == "" {
		style.BorderColor = style.Color
	}
	if style.LineHeight < 1 {
		style.LineHeight = 1
	}
	return style, nil
}
