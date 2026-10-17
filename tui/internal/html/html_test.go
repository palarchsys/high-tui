package html

import (
	"strings"
	"testing"
)

func doc(body string) string {
	return `<!DOCTYPE html>
<html lang="fr">
<head>
<meta charset="utf-8">
<title>W</title>
<link rel="stylesheet" href="t.css">
</head>
<body>
` + body + `
</body>
</html>`
}

func TestAbsentActionDoesNothing(t *testing.T) {
	raw := doc(`<table id="w">
<tr id="r" data-lines="1">
<td><button id="b" type="button" data-enter="save">Enregistrer</button></td>
</tr>
</table>`)
	tree, err := Parse(raw, "t.html")
	if err != nil {
		t.Fatal(err)
	}
	win, err := Window(tree, "t.html")
	if err != nil {
		t.Fatal(err)
	}
	cell := win.Grid.Rows[0].Cells[0]
	zone := cell.Zone
	if zone.Focus.Kind != "" {
		t.Fatalf("focus déclenché : %q", zone.Focus.Kind)
	}
	if cell.Fr != 1 || zone.Readonly || zone.File != "" || zone.Value != "" {
		t.Fatalf("défauts %#v", zone)
	}
	if zone.Enter.Kind != "save" {
		t.Fatal(zone.Enter)
	}
}

func TestEmptyAttribute(t *testing.T) {
	raw := doc(`<table id="w"><tr id="r"><td><button id="b" type="button" data-enter="">Enregistrer</button></td></tr></table>`)
	_, err := Parse(raw, "t.html")
	if err == nil || !strings.Contains(err.Error(), "attribut vide") {
		t.Fatal(err)
	}
}

func TestReadonlyFalseRefused(t *testing.T) {
	raw := doc(`<table id="w"><tr><td><textarea id="c" readonly="false"></textarea></td></tr></table>`)
	tree, err := Parse(raw, "t.html")
	if err != nil {
		t.Fatal(err)
	}
	_, err = Window(tree, "t.html")
	if err == nil || !strings.Contains(err.Error(), "readonly invalide") {
		t.Fatal(err)
	}
}

func TestNestedTable(t *testing.T) {
	raw := doc(`<table id="w">
<tr>
<td>
<table id="dedans">
<tr data-lines="2"><td><textarea id="corps">placeholder</textarea></td></tr>
</table>
</td>
</tr>
</table>`)
	tree, err := Parse(raw, "t.html")
	if err != nil {
		t.Fatal(err)
	}
	win, err := Window(tree, "t.html")
	if err != nil {
		t.Fatal(err)
	}
	cell := win.Grid.Rows[0].Cells[0]
	if cell.Zone != nil || cell.Grid == nil || cell.Grid.ID != "dedans" {
		t.Fatal("grille imbriquée absente")
	}
	zone := cell.Grid.Rows[0].Cells[0].Zone
	if zone == nil || zone.ID != "corps" || zone.Value != "placeholder" || zone.File != "" {
		t.Fatalf("%#v", zone)
	}
}

func TestRowTagRefused(t *testing.T) {
	raw := doc(`<table id="w"><row id="r"></row></table>`)
	tree, err := Parse(raw, "t.html")
	if err != nil {
		t.Fatal(err)
	}
	_, err = Window(tree, "t.html")
	if err == nil || !strings.Contains(err.Error(), "balise inconnue") {
		t.Fatal(err)
	}
}
