package ui

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"high-tui/internal/layout"
	"high-tui/internal/model"
)

func TestFocusFillsNextAndClearsDownstream(t *testing.T) {
	a := cascadeApp(t)
	if err := a.applyFocus(); err != nil {
		t.Fatal(err)
	}
	fiches := a.findZone(0, "nav-fiches")
	champs := a.findZone(0, "nav-champs")
	menu := a.findZone(0, "nav-menu")
	if fiches == nil || len(fiches.Items) != 1 || fiches.Items[0].Label != "Jeton pro" {
		t.Fatalf("fiches %#v", fiches)
	}
	if len(champs.Items) != 0 {
		t.Fatalf("champs %#v", champs.Items)
	}
	if len(menu.Items) != 1 || menu.Items[0].Label != "Ouvrir" {
		t.Fatal(menu.Items)
	}
	box := placedNav(t, a, "nav-fiches")
	if len(box.Zone.Items) != 1 || box.Zone.Items[0].Label != "Jeton pro" {
		t.Fatal(box.Zone.Items)
	}

	a.sel[key("main", "nav-rubriques")] = 1
	a.focusMark = ""
	if err := a.applyFocus(); err != nil {
		t.Fatal(err)
	}
	if len(a.findZone(0, "nav-fiches").Items) != 0 {
		t.Fatal(a.findZone(0, "nav-fiches").Items)
	}
}

func TestFocusReadsFileAndCompletesAction(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "suite.json"), []byte(`[{"label":"Suite","file":"note.env","name":"HOTE","value":"1"}]`), 0o644); err != nil {
		t.Fatal(err)
	}
	a := cascadeApp(t)
	a.root = dir
	rub := a.findZone(0, "nav-rubriques")
	rub.Items[0].NextItems = nil
	rub.Items[0].Next = "suite.json"
	a.syncZone(0, "nav-rubriques")
	a.focusMark = ""
	if err := a.applyFocus(); err != nil {
		t.Fatal(err)
	}
	fiches := a.findZone(0, "nav-fiches")
	if len(fiches.Items) != 1 || fiches.Items[0].Label != "Suite" {
		t.Fatal(fiches.Items)
	}
	a.focus = 1
	a.focusMark = ""
	if err := a.applyFocus(); err != nil {
		t.Fatal(err)
	}
	act, label := a.action("enter")
	if label != "Suite" || act.Kind != "set" || act.File != "note.env" || act.Name != "HOTE" || act.Value != "1" {
		t.Fatalf("%s %#v", label, act)
	}
}

func TestCommandListAndStaleResult(t *testing.T) {
	a := cascadeApp(t)
	msg := execList(context.Background(), 3, "nav-fiches", `printf '%s' '[{"label":"Echo"}]'`)().(listDone)
	if msg.err != nil || msg.cancelled {
		t.Fatal(msg.err, msg.cancelled)
	}
	a.listGen = 4
	a.onListDone(msg)
	if a.findZone(0, "nav-fiches").Items[0].Label != "Ancien" {
		t.Fatal("un résultat ancien a remplacé la liste")
	}
	a.listGen = msg.gen
	a.onListDone(msg)
	got := a.findZone(0, "nav-fiches")
	if len(got.Items) != 1 || got.Items[0].Label != "Echo" {
		t.Fatal(got.Items)
	}
	if placedNav(t, a, "nav-fiches").Zone.Items[0].Label != "Echo" {
		t.Fatal("liste non affichée")
	}
}

func TestBadFileClearsTarget(t *testing.T) {
	a := cascadeApp(t)
	rub := a.findZone(0, "nav-rubriques")
	rub.Items[0].NextItems = nil
	rub.Items[0].Next = "absent.json"
	a.syncZone(0, "nav-rubriques")
	a.focusMark = ""
	if err := a.applyFocus(); err != nil {
		t.Fatal(err)
	}
	if len(a.findZone(0, "nav-fiches").Items) != 0 {
		t.Fatal(a.findZone(0, "nav-fiches").Items)
	}
	if !strings.Contains(a.status, "fichier introuvable") {
		t.Fatal(a.status)
	}
}

func cascadeApp(t *testing.T) *app {
	t.Helper()
	style := model.Style{Color: "#e6edf3", Background: "#0d1117", LineHeight: 1}
	rubriques := model.Zone{
		Tag: "nav", ID: "nav-rubriques", Label: "Rubriques", NextNav: "nav-fiches",
		Base: style, FocusOn: style,
		Items: []model.Item{
			{ID: "reseau", Label: "Réseau", NextItems: []model.Item{{Label: "Jeton pro"}}},
			{ID: "notes", Label: "Notes"},
		},
	}
	fiches := model.Zone{
		Tag: "nav", ID: "nav-fiches", Label: "Fiches", NextNav: "nav-champs",
		Base: style, FocusOn: style,
		Enter: model.Action{Kind: "set", Target: "valeur"},
		Items: []model.Item{{ID: "ancien", Label: "Ancien"}},
	}
	champs := model.Zone{
		Tag: "nav", ID: "nav-champs", Label: "Champs", Base: style, FocusOn: style,
		Items: []model.Item{{ID: "nom", Label: "Nom"}},
	}
	menu := model.Zone{
		Tag: "nav", ID: "nav-menu", Label: "Menu", Base: style, FocusOn: style,
		Items: []model.Item{{ID: "ouvrir", Label: "Ouvrir"}},
	}
	valeur := model.Zone{Tag: "input", ID: "valeur", Name: "HOTE", Base: style, FocusOn: style}
	win := model.Window{
		ID: "main", File: "main.html",
		Style: model.Style{Color: "#e6edf3", Background: "#0d1117", LineHeight: 1, MaxWidth: 80},
		Grid: model.Grid{ID: "main", Rows: []model.GridRow{{
			ID: "listes", Lines: 4,
			Cells: []model.Cell{
				{Fr: 1, Zone: &rubriques},
				{Fr: 1, Zone: &fiches},
				{Fr: 1, Zone: &champs},
				{Fr: 1, Zone: &menu},
				{Fr: 1, Zone: &valeur},
			},
		}}},
	}
	a := &app{
		root: t.TempDir(), wins: []model.Window{win}, origin: []model.Window{win}, width: 100,
		sel: map[string]int{}, scroll: map[string]int{}, val: map[string]string{},
		file: map[string]string{}, name: map[string]string{}, dirty: map[string]bool{},
		resolved: map[string]string{}, focusAt: map[string]int{}, status: "Prêt",
	}
	if err := a.layoutAll(); err != nil {
		t.Fatal(err)
	}
	return a
}

func placedNav(t *testing.T, a *app, id string) layout.Box {
	t.Helper()
	var found layout.Box
	ok := false
	a.current().EachBox(func(box *layout.Box) {
		if box.Zone.ID == id {
			found = *box
			ok = true
		}
	})
	if !ok {
		t.Fatalf("zone absente %s", id)
	}
	return found
}
