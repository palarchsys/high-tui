// Package html parses one window document.
package html

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"high-tui/internal/items"
	"high-tui/internal/model"
)

var actions = map[string]bool{
	"load": true, "edit": true, "window": true, "return": true,
	"save": true, "refresh": true, "run": true, "set": true,
}

var attrs = map[string]bool{
	"id": true, "class": true, "lang": true, "charset": true, "rel": true, "href": true,
	"name": true, "value": true, "type": true, "readonly": true,
	"data-lines": true, "data-fr": true, "data-file": true, "data-name": true,
	"data-save": true, "data-label": true, "data-value": true,
	"data-focus": true, "data-focus-target": true, "data-focus-file": true,
	"data-focus-name": true, "data-focus-href": true, "data-focus-run": true,
	"data-focus-output": true, "data-focus-value": true,
	"data-enter": true, "data-enter-target": true, "data-enter-file": true,
	"data-enter-name": true, "data-enter-href": true, "data-enter-run": true,
	"data-enter-output": true, "data-enter-value": true,
	"data-items": true, "data-next": true,
}

var zoneTags = map[string]bool{
	"header": true, "nav": true, "input": true, "textarea": true,
	"button": true, "footer": true,
}

var voidTags = map[string]bool{"meta": true, "link": true, "input": true}

// Node is a raw element before validation.
type Node struct {
	Tag   string
	Attr  map[string]string
	Text  string
	Kids  []Node
	Where string
}

// Parse reads one HTML document.
func Parse(text, where string) (Node, error) {
	p := &parser{s: text, where: where}
	p.skipSpaceAndComments()
	if !p.consumeDoctype() {
		return Node{}, model.Errf("doctype attendu (%s)", where)
	}
	p.skipSpaceAndComments()
	if p.i >= len(p.s) || p.s[p.i] != '<' {
		return Node{}, model.Errf("balise html attendue (%s)", where)
	}
	root, err := p.element()
	if err != nil {
		return Node{}, err
	}
	p.skipSpaceAndComments()
	if p.i < len(p.s) {
		return Node{}, model.Errf("texte inattendu (%s)", where)
	}
	if root.Tag != "html" {
		return Node{}, model.Errf("balise html attendue (%s)", where)
	}
	return root, nil
}

type parser struct {
	s, where string
	i        int
}

func (p *parser) consumeDoctype() bool {
	if p.i >= len(p.s) || p.s[p.i] != '<' {
		return false
	}
	rest := p.s[p.i:]
	if len(rest) < 9 || !strings.EqualFold(rest[:9], "<!doctype") {
		return false
	}
	j := strings.IndexByte(rest, '>')
	if j < 0 {
		return false
	}
	p.i += j + 1
	return true
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
	if p.i >= len(p.s) || p.s[p.i] != '<' {
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
		if p.i < len(p.s) && p.s[p.i] != '=' {
			if name == "readonly" {
				if _, ok := attr[name]; ok {
					return Node{}, model.Errf("attribut en double « %s » (%s)", name, p.where)
				}
				attr[name] = "readonly"
				continue
			}
			return Node{}, model.Errf("valeur manquante « %s » (%s)", name, p.where)
		}
		if p.i >= len(p.s) || p.s[p.i] != '=' {
			return Node{}, model.Errf("valeur manquante « %s » (%s)", name, p.where)
		}
		p.i++
		p.skipSpace()
		val, err := p.readValue()
		if err != nil {
			return Node{}, err
		}
		if val == "" && name != "readonly" {
			return Node{}, model.Errf("attribut vide « %s » (%s)", name, p.where)
		}
		if _, ok := attr[name]; ok {
			return Node{}, model.Errf("attribut en double « %s » (%s)", name, p.where)
		}
		attr[name] = val
	}
	if voidTags[tag] {
		return Node{Tag: tag, Attr: attr, Where: p.where}, nil
	}
	var b strings.Builder
	var kids []Node
	for {
		if p.i >= len(p.s) {
			return Node{}, model.Errf("balise non fermée « %s » (%s)", tag, p.where)
		}
		if strings.HasPrefix(p.s[p.i:], "<!--") {
			j := strings.Index(p.s[p.i+4:], "-->")
			if j < 0 {
				return Node{}, model.Errf("commentaire non fermé (%s)", p.where)
			}
			p.i += 4 + j + 3
			continue
		}
		if strings.HasPrefix(p.s[p.i:], "</") {
			p.i += 2
			end := p.readName()
			p.skipSpace()
			if p.i >= len(p.s) || p.s[p.i] != '>' || end != tag {
				return Node{}, model.Errf("fermeture « %s » invalide (%s)", tag, p.where)
			}
			p.i++
			return Node{Tag: tag, Attr: attr, Text: strings.TrimSpace(b.String()), Kids: kids, Where: p.where}, nil
		}
		if p.s[p.i] == '<' {
			kid, err := p.element()
			if err != nil {
				return Node{}, err
			}
			kids = append(kids, kid)
			continue
		}
		j := strings.IndexByte(p.s[p.i:], '<')
		if j < 0 {
			return Node{}, model.Errf("balise non fermée « %s » (%s)", tag, p.where)
		}
		b.WriteString(p.s[p.i : p.i+j])
		p.i += j
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

func actionOf(attr map[string]string, prefix, where string, loose bool) (model.Action, error) {
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
		if act.File == "" && !loose {
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
		if loose {
			break
		}
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

func take(attr map[string]string, key string) string {
	v := attr[key]
	delete(attr, key)
	return v
}

func takeOK(attr map[string]string, key string) (string, bool) {
	v, ok := attr[key]
	if ok {
		delete(attr, key)
	}
	return v, ok
}

func leftover(attr map[string]string, where string) error {
	for key := range attr {
		return model.Errf("attribut refusé « %s » (%s)", key, where)
	}
	return nil
}

func keepID(attr map[string]string, ids map[string]bool, where string, required bool) (string, error) {
	id := take(attr, "id")
	if id == "" {
		if required {
			return "", model.Errf("attribut manquant « id » (%s)", where)
		}
		return "", nil
	}
	if ids[id] {
		return "", model.Errf("identifiant en double : %s", id)
	}
	ids[id] = true
	return id, nil
}

// Window converts a parsed document into a model window.
func Window(root Node, file string) (model.Window, error) {
	if root.Tag != "html" {
		return model.Window{}, model.Errf("balise html attendue (%s)", file)
	}
	attr := copyAttr(root.Attr)
	if take(attr, "lang") == "" {
		return model.Window{}, model.Errf("attribut manquant « lang » (%s)", file)
	}
	if err := leftover(attr, "html"); err != nil {
		return model.Window{}, err
	}
	if root.Text != "" || len(root.Kids) != 2 || root.Kids[0].Tag != "head" || root.Kids[1].Tag != "body" {
		return model.Window{}, model.Errf("structure html refusée (%s)", file)
	}
	label, err := headOf(root.Kids[0], file)
	if err != nil {
		return model.Window{}, err
	}
	ids := map[string]bool{}
	grid, err := tableOf(root.Kids[1], ids, file)
	if err != nil {
		return model.Window{}, err
	}
	win := model.Window{
		ID: grid.ID, Label: label, File: file, Classes: grid.Classes, Grid: grid,
	}
	if err := checkTargets(&win); err != nil {
		return model.Window{}, err
	}
	ensureReturn(&win, ids)
	return win, nil
}

func copyAttr(in map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range in {
		out[k] = v
	}
	return out
}

func headOf(head Node, file string) (string, error) {
	attr := copyAttr(head.Attr)
	if err := leftover(attr, "head"); err != nil {
		return "", err
	}
	if head.Text != "" {
		return "", model.Errf("texte inattendu (%s)", file)
	}
	var title, href string
	seen := map[string]bool{}
	for _, kid := range head.Kids {
		if seen[kid.Tag] {
			return "", model.Errf("balise en double « %s » (%s)", kid.Tag, file)
		}
		seen[kid.Tag] = true
		switch kid.Tag {
		case "meta":
			meta := copyAttr(kid.Attr)
			charset := strings.ToLower(take(meta, "charset"))
			if charset != "utf-8" {
				return "", model.Errf("charset refusé (%s)", file)
			}
			if err := leftover(meta, "meta"); err != nil {
				return "", err
			}
			if kid.Text != "" || len(kid.Kids) > 0 {
				return "", model.Errf("enfant refusé (meta)")
			}
		case "title":
			titleAttr := copyAttr(kid.Attr)
			if err := leftover(titleAttr, "title"); err != nil {
				return "", err
			}
			if len(kid.Kids) > 0 || kid.Text == "" {
				return "", model.Errf("titre manquant (%s)", file)
			}
			title = kid.Text
		case "link":
			link := copyAttr(kid.Attr)
			if take(link, "rel") != "stylesheet" {
				return "", model.Errf("feuille de style manquante (%s)", file)
			}
			href = take(link, "href")
			if err := leftover(link, "link"); err != nil {
				return "", err
			}
			want := strings.TrimSuffix(file, ".html") + ".css"
			if href != want {
				return "", model.Errf("feuille de style refusée « %s » (%s)", href, file)
			}
		default:
			return "", model.Errf("balise inconnue « %s » (%s)", kid.Tag, file)
		}
	}
	if !seen["meta"] || !seen["title"] || !seen["link"] {
		return "", model.Errf("en-tête incomplet (%s)", file)
	}
	return title, nil
}

func tableOf(body Node, ids map[string]bool, file string) (model.Grid, error) {
	attr := copyAttr(body.Attr)
	if err := leftover(attr, "body"); err != nil {
		return model.Grid{}, err
	}
	if body.Text != "" || len(body.Kids) != 1 || body.Kids[0].Tag != "table" {
		return model.Grid{}, model.Errf("table attendue (%s)", file)
	}
	return convertGrid(body.Kids[0], ids)
}

func convertGrid(node Node, ids map[string]bool) (model.Grid, error) {
	if node.Tag != "table" {
		return model.Grid{}, model.Errf("balise inconnue « %s » (%s)", node.Tag, node.Attr["id"])
	}
	attr := copyAttr(node.Attr)
	id, err := keepID(attr, ids, "table", true)
	if err != nil {
		return model.Grid{}, err
	}
	classes := classesOf(take(attr, "class"))
	if err := leftover(attr, id); err != nil {
		return model.Grid{}, err
	}
	if node.Text != "" || len(node.Kids) == 0 {
		return model.Grid{}, model.Errf("grille sans rangée (%s)", id)
	}
	grid := model.Grid{ID: id, Classes: classes}
	for _, kid := range node.Kids {
		if kid.Tag != "tr" {
			return model.Grid{}, model.Errf("balise inconnue « %s » (%s)", kid.Tag, id)
		}
		row, err := convertRow(kid, ids)
		if err != nil {
			return model.Grid{}, err
		}
		grid.Rows = append(grid.Rows, row)
	}
	return grid, nil
}

func convertRow(node Node, ids map[string]bool) (model.GridRow, error) {
	attr := copyAttr(node.Attr)
	id, err := keepID(attr, ids, "tr", false)
	if err != nil {
		return model.GridRow{}, err
	}
	where := id
	if where == "" {
		where = "tr"
	}
	classes := classesOf(take(attr, "class"))
	lines := 1
	if raw, ok := takeOK(attr, "data-lines"); ok {
		n, err := positive(raw, where)
		if err != nil {
			return model.GridRow{}, err
		}
		lines = n
	}
	if err := leftover(attr, where); err != nil {
		return model.GridRow{}, err
	}
	if node.Text != "" || len(node.Kids) == 0 {
		return model.GridRow{}, model.Errf("rangée sans colonne (%s)", where)
	}
	row := model.GridRow{ID: id, Classes: classes, Lines: lines}
	for _, kid := range node.Kids {
		if kid.Tag != "td" {
			return model.GridRow{}, model.Errf("balise inconnue « %s » (%s)", kid.Tag, where)
		}
		cell, err := convertCell(kid, ids, where)
		if err != nil {
			return model.GridRow{}, err
		}
		row.Cells = append(row.Cells, cell)
	}
	return row, nil
}

func convertCell(node Node, ids map[string]bool, row string) (model.Cell, error) {
	attr := copyAttr(node.Attr)
	if _, err := keepID(attr, ids, "td", false); err != nil {
		return model.Cell{}, err
	}
	classes := classesOf(take(attr, "class"))
	if classes != nil {
		return model.Cell{}, model.Errf("attribut refusé « class » (%s)", row)
	}
	fr := 1
	if raw, ok := takeOK(attr, "data-fr"); ok {
		n, err := positive(raw, row)
		if err != nil {
			return model.Cell{}, err
		}
		fr = n
	}
	if err := leftover(attr, row); err != nil {
		return model.Cell{}, err
	}
	if node.Text != "" || len(node.Kids) != 1 {
		return model.Cell{}, model.Errf("cellule invalide (%s)", row)
	}
	kid := node.Kids[0]
	cell := model.Cell{Fr: fr}
	if kid.Tag == "table" {
		grid, err := convertGrid(kid, ids)
		if err != nil {
			return model.Cell{}, err
		}
		cell.Grid = &grid
		return cell, nil
	}
	if !zoneTags[kid.Tag] {
		return model.Cell{}, model.Errf("balise inconnue « %s » (%s)", kid.Tag, row)
	}
	zone, err := convertZone(kid, ids)
	if err != nil {
		return model.Cell{}, err
	}
	cell.Zone = &zone
	return cell, nil
}

func convertZone(node Node, ids map[string]bool) (model.Zone, error) {
	attr := copyAttr(node.Attr)
	id, err := keepID(attr, ids, node.Tag, true)
	if err != nil {
		return model.Zone{}, err
	}
	classes := classesOf(take(attr, "class"))
	loose := node.Tag == "nav"
	focus, err := actionOf(attr, "data-focus", id, loose)
	if err != nil {
		return model.Zone{}, err
	}
	enter, err := actionOf(attr, "data-enter", id, loose)
	if err != nil {
		return model.Zone{}, err
	}
	dropAction(attr, "data-focus")
	dropAction(attr, "data-enter")
	readonly, err := readonlyOf(attr, id)
	if err != nil {
		return model.Zone{}, err
	}
	save := take(attr, "data-save")
	if save != "" && save != "enter" && save != "button" {
		return model.Zone{}, model.Errf("save invalide (%s)", id)
	}
	zone := model.Zone{
		Tag: node.Tag, ID: id, Classes: classes, Readonly: readonly, Save: save,
		Focus: focus, Enter: enter,
	}
	switch node.Tag {
	case "header":
		if err := leftover(attr, id); err != nil {
			return model.Zone{}, err
		}
		if node.Text != "" || len(node.Kids) != 1 || node.Kids[0].Tag != "h1" {
			return model.Zone{}, model.Errf("titre manquant (%s)", id)
		}
		h := node.Kids[0]
		if err := plainText(h, "h1"); err != nil {
			return model.Zone{}, err
		}
		zone.Label = h.Text
	case "nav":
		zone.ItemsFrom = take(attr, "data-items")
		zone.NextNav = take(attr, "data-next")
		if err := checkItemsFrom(zone.ItemsFrom, id); err != nil {
			return model.Zone{}, err
		}
		if err := leftover(attr, id); err != nil {
			return model.Zone{}, err
		}
		if node.Text != "" || len(node.Kids) != 2 || node.Kids[0].Tag != "h2" || node.Kids[1].Tag != "ul" {
			return model.Zone{}, model.Errf("liste invalide (%s)", id)
		}
		if err := plainText(node.Kids[0], "h2"); err != nil {
			return model.Zone{}, err
		}
		zone.Label = node.Kids[0].Text
		items, err := listItems(node.Kids[1], ids)
		if err != nil {
			return model.Zone{}, err
		}
		zone.Items = items
	case "input":
		if node.Text != "" || len(node.Kids) > 0 {
			return model.Zone{}, model.Errf("enfant refusé (%s)", id)
		}
		if err := inputType(attr, id); err != nil {
			return model.Zone{}, err
		}
		zone.File = take(attr, "data-file")
		zone.Name = first(take(attr, "data-name"), take(attr, "name"))
		zone.Value = take(attr, "value")
		if err := leftover(attr, id); err != nil {
			return model.Zone{}, err
		}
	case "textarea":
		if len(node.Kids) > 0 {
			return model.Zone{}, model.Errf("enfant refusé (%s)", id)
		}
		zone.File = take(attr, "data-file")
		zone.Name = take(attr, "data-name")
		zone.Value = node.Text
		if err := leftover(attr, id); err != nil {
			return model.Zone{}, err
		}
	case "button":
		if err := buttonType(attr, id); err != nil {
			return model.Zone{}, err
		}
		if err := leftover(attr, id); err != nil {
			return model.Zone{}, err
		}
		if len(node.Kids) > 0 || node.Text == "" {
			return model.Zone{}, model.Errf("label manquant (%s)", id)
		}
		if zone.Enter.Kind == "" {
			return model.Zone{}, model.Errf("action manquante (%s)", id)
		}
		zone.Label = node.Text
	case "footer":
		if err := leftover(attr, id); err != nil {
			return model.Zone{}, err
		}
		if len(node.Kids) > 0 {
			return model.Zone{}, model.Errf("enfant refusé (%s)", id)
		}
		zone.Label = node.Text
	default:
		return model.Zone{}, model.Errf("balise inconnue « %s » (%s)", node.Tag, id)
	}
	return zone, nil
}

func first(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func dropAction(attr map[string]string, prefix string) {
	for _, key := range []string{"", "-target", "-file", "-name", "-href", "-run", "-output", "-value"} {
		delete(attr, prefix+key)
	}
}

func readonlyOf(attr map[string]string, id string) (bool, error) {
	raw, ok := takeOK(attr, "readonly")
	if !ok {
		return false, nil
	}
	if raw != "" && raw != "readonly" && raw != "true" {
		return false, model.Errf("readonly invalide (%s)", id)
	}
	return true, nil
}

func inputType(attr map[string]string, id string) error {
	raw, ok := takeOK(attr, "type")
	if ok && raw != "text" {
		return model.Errf("type refusé (%s)", id)
	}
	return nil
}

func buttonType(attr map[string]string, id string) error {
	if take(attr, "type") != "button" {
		return model.Errf("type refusé (%s)", id)
	}
	return nil
}

func plainText(node Node, tag string) error {
	attr := copyAttr(node.Attr)
	if err := leftover(attr, tag); err != nil {
		return err
	}
	if len(node.Kids) > 0 || node.Text == "" {
		return model.Errf("texte manquant (%s)", tag)
	}
	return nil
}

func listItems(ul Node, ids map[string]bool) ([]model.Item, error) {
	attr := copyAttr(ul.Attr)
	if err := leftover(attr, "ul"); err != nil {
		return nil, err
	}
	if ul.Text != "" {
		return nil, model.Errf("texte inattendu (ul)")
	}
	var items []model.Item
	for _, kid := range ul.Kids {
		if kid.Tag != "li" {
			return nil, model.Errf("balise inconnue « %s » (ul)", kid.Tag)
		}
		item, err := convertItem(kid, ids)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func convertItem(node Node, ids map[string]bool) (model.Item, error) {
	attr := copyAttr(node.Attr)
	id, err := keepID(attr, ids, "li", true)
	if err != nil {
		return model.Item{}, err
	}
	_ = take(attr, "class")
	focus, err := actionOf(attr, "data-focus", id, false)
	if err != nil {
		return model.Item{}, err
	}
	enter, err := actionOf(attr, "data-enter", id, false)
	if err != nil {
		return model.Item{}, err
	}
	dropAction(attr, "data-focus")
	dropAction(attr, "data-enter")
	label := take(attr, "data-label")
	value := take(attr, "data-value")
	if err := leftover(attr, id); err != nil {
		return model.Item{}, err
	}
	if len(node.Kids) > 1 {
		return model.Item{}, model.Errf("enfant refusé (%s)", id)
	}
	if len(node.Kids) == 1 {
		if node.Text != "" || node.Kids[0].Tag != "a" {
			return model.Item{}, model.Errf("enfant refusé (%s)", id)
		}
		anchor := node.Kids[0]
		aAttr := copyAttr(anchor.Attr)
		if take(aAttr, "href") == "" {
			return model.Item{}, model.Errf("href manquant (%s)", id)
		}
		if err := leftover(aAttr, id); err != nil {
			return model.Item{}, err
		}
		if len(anchor.Kids) > 0 || anchor.Text == "" {
			return model.Item{}, model.Errf("label manquant (%s)", id)
		}
		if label == "" {
			label = anchor.Text
		}
	} else if label == "" {
		label = node.Text
	}
	if label == "" {
		return model.Item{}, model.Errf("label manquant (%s)", id)
	}
	if enter.Kind == "set" && enter.Value == "" {
		enter.Value = value
		if enter.Value == "" {
			return model.Item{}, model.Errf("valeur manquante (%s)", id)
		}
	}
	return model.Item{ID: id, Label: label, Value: value, Focus: focus, Enter: enter}, nil
}

func zonesOf(win *model.Window) []*model.Zone { return win.Grid.Zones() }

func checkTargets(win *model.Window) error {
	byID := map[string]model.Zone{}
	for _, zone := range zonesOf(win) {
		byID[zone.ID] = *zone
	}
	one := func(act model.Action) error {
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
		if zone.NextNav != "" {
			target, ok := byID[zone.NextNav]
			if !ok || target.Tag != "nav" {
				return model.Errf("cible absente : %s", zone.NextNav)
			}
		}
		if err := one(zone.Focus); err != nil {
			return err
		}
		if err := one(zone.Enter); err != nil {
			return err
		}
		for _, item := range zone.Items {
			if err := one(item.Focus); err != nil {
				return err
			}
			if err := one(item.Enter); err != nil {
				return err
			}
		}
	}
	return listCycles(win)
}

func checkItemsFrom(source, id string) error {
	if source == "" {
		return nil
	}
	return items.CheckSource(source, id)
}

func listCycles(win *model.Window) error {
	byID := map[string]model.Zone{}
	for _, zone := range zonesOf(win) {
		byID[zone.ID] = *zone
	}
	for id, zone := range byID {
		if zone.NextNav == "" {
			continue
		}
		seen := map[string]bool{}
		cur := id
		for cur != "" {
			if seen[cur] {
				return model.Errf("cycle de listes (%s)", id)
			}
			seen[cur] = true
			next, ok := byID[cur]
			if !ok {
				break
			}
			cur = next.NextNav
		}
	}
	return nil
}

func ensureReturn(win *model.Window, ids map[string]bool) {
	found := false
	win.Grid.Each(func(zone *model.Zone) {
		if zone.Focus.Kind == "return" || zone.Enter.Kind == "return" {
			found = true
		}
		for _, item := range zone.Items {
			if item.Focus.Kind == "return" || item.Enter.Kind == "return" {
				found = true
			}
		}
	})
	if found {
		return
	}
	added := false
	win.Grid.Each(func(zone *model.Zone) {
		if added || zone.Tag != "nav" {
			return
		}
		id := zone.ID + "-retour"
		for ids[id] {
			id += "-r"
		}
		ids[id] = true
		zone.Items = append(zone.Items, model.Item{
			ID: id, Label: "Retour", Enter: model.Action{Kind: "return"},
		})
		added = true
	})
}
