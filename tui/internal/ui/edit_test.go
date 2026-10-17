package ui

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"high-tui/internal/model"
	"high-tui/internal/run"
)

func TestShiftEnterAddsLine(t *testing.T) {
	a := editApp(t)
	a.keyMode = keysEnhanced
	a.onKey(keyMsg("enter"))
	if a.edit == nil {
		t.Fatal("édition non ouverte")
	}
	a.edit.state.Cursor = 1
	a.onKey(keyMsg("shift+enter"))
	if a.edit == nil || a.edit.state.Text != "a\nb" {
		t.Fatalf("texte inattendu")
	}
	a.onKey(keyMsg("enter"))
	if a.edit != nil {
		t.Fatal("entrée n'a pas fermé")
	}
	if a.val[key("main", "corps")] != "a\nb" {
		t.Fatal(a.val[key("main", "corps")])
	}
}

func TestShiftEnterWithoutEnhancements(t *testing.T) {
	a := editApp(t)
	a.keyMode = keysPlain
	a.onKey(keyMsg("enter"))
	a.edit.state.Cursor = 1
	a.onKey(keyMsg("shift+enter"))
	if a.edit.state.Text != "ab" {
		t.Fatal(a.edit.state.Text)
	}
	if a.status != noShiftMsg {
		t.Fatal(a.status)
	}
	a.onKey(keyMsg("enter"))
	if a.edit != nil {
		t.Fatal("entrée doit fermer")
	}
}

func TestAbsentFocusDoesNothing(t *testing.T) {
	a := editApp(t)
	before := a.val[key("main", "corps")]
	if err := a.applyFocus(); err != nil {
		t.Fatal(err)
	}
	if a.val[key("main", "corps")] != before {
		t.Fatal(a.val[key("main", "corps")])
	}
	if a.status != "Prêt" {
		t.Fatal(a.status)
	}
}

func TestScreenWords(t *testing.T) {
	root := filepath.Join("..", "..")
	frame, err := Screen(root, 80)
	if err != nil {
		t.Fatal(err)
	}
	for _, word := range []string{"Rubriques", "Fiches", "Champs", "Menu", "127.0.0.1", "Jeton pro"} {
		if !strings.Contains(frame, word) {
			t.Fatalf("manque %s", word)
		}
	}
	for _, line := range strings.Split(frame, "\n") {
		if lipgloss.Width(line) != 80 {
			t.Fatalf("largeur %d pour %q", lipgloss.Width(line), line)
		}
	}
}

func TestSetWritesTempFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "vars.env")
	if err := os.WriteFile(path, []byte("HOTE=ancien\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	a := editApp(t)
	a.root = dir
	err := a.set(model.Action{Kind: "set", File: "vars.env", Name: "HOTE", Value: "nouveau", Target: "corps"})
	if err != nil {
		t.Fatal(err)
	}
	got, err := run.ReadVar(path, "HOTE")
	if err != nil || got != "nouveau" {
		t.Fatal(got, err)
	}
	if a.status != "Enregistré" {
		t.Fatal(a.status)
	}
}

func editApp(t *testing.T) *app {
	t.Helper()
	style := model.Style{
		Color: "#e6edf3", Background: "#161b22", BorderColor: "#30363d",
		BorderWidth: 1, LineHeight: 1, Pad: [4]int{0, 1, 0, 1},
	}
	zone := model.Zone{
		Tag: "textarea", ID: "corps", Fr: 1, Readonly: false,
		Base: style, FocusOn: style,
		Enter: model.Action{Kind: "edit", Target: "corps"},
	}
	win := model.Window{
		ID: "main", File: "main.html",
		Style: model.Style{
			Color: "#e6edf3", Background: "#0d1117", BorderColor: "#e6edf3",
			LineHeight: 1, MaxWidth: 80,
		},
		Rows: []model.Row{{
			ID: "page", Tag: "main", Lines: 4, Zones: []model.Zone{zone},
		}},
	}
	a := &app{
		root: t.TempDir(), wins: []model.Window{win}, origin: []model.Window{win}, width: 40,
		sel: map[string]int{}, scroll: map[string]int{}, val: map[string]string{},
		file: map[string]string{}, name: map[string]string{}, dirty: map[string]bool{},
		resolved: map[string]string{}, focusAt: map[string]int{}, status: "Prêt",
	}
	a.val[key("main", "corps")] = "ab"
	if err := a.layoutAll(); err != nil {
		t.Fatal(err)
	}
	return a
}

func keyMsg(name string) tea.KeyPressMsg {
	switch name {
	case "enter":
		return tea.KeyPressMsg{Code: tea.KeyEnter}
	case "shift+enter":
		return tea.KeyPressMsg{Code: tea.KeyEnter, Mod: tea.ModShift}
	case "esc":
		return tea.KeyPressMsg{Code: tea.KeyEsc}
	default:
		return tea.KeyPressMsg{Text: name, Code: rune([]rune(name)[0])}
	}
}
