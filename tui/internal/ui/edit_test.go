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
	if a.status != noShiftMsg {
		t.Fatal(a.status)
	}
	a.edit.state.Cursor = 1
	a.onKey(keyMsg("shift+enter"))
	if a.edit == nil || a.edit.state.Text != "a\nb" {
		t.Fatal(a.edit)
	}
	if a.status == noShiftMsg {
		t.Fatal(a.status)
	}
	a.onKey(keyMsg("enter"))
	if a.edit != nil {
		t.Fatal("entrée doit fermer")
	}
	if a.val[key("main", "corps")] != "a\nb" {
		t.Fatal(a.val[key("main", "corps")])
	}
}

func TestShiftEnterAfterOtherKeys(t *testing.T) {
	a := editApp(t)
	a.onKey(keyMsg("left"))
	if a.keyMode != keysPlain {
		t.Fatal(a.keyMode)
	}
	a.onKey(keyMsg("enter"))
	a.edit.state.Cursor = 1
	a.onKey(keyMsg("shift+enter"))
	if a.edit == nil || a.edit.state.Text != "a\nb" {
		t.Fatal(a.edit)
	}
}

func TestEnterClosesWithoutNewline(t *testing.T) {
	a := editApp(t)
	a.keyMode = keysEnhanced
	a.onKey(keyMsg("enter"))
	a.edit.state.Cursor = 1
	a.onKey(keyMsg("enter"))
	if a.edit != nil {
		t.Fatal("entrée n'a pas fermé")
	}
	if a.val[key("main", "corps")] != "ab" {
		t.Fatal(a.val[key("main", "corps")])
	}
}

func TestQuitIsCtrlQ(t *testing.T) {
	a := editApp(t)
	if cmd := a.onKey(keyMsg("q")); cmd != nil {
		t.Fatal("q quitte")
	}
	a.onKey(keyMsg("enter"))
	a.onKey(keyMsg("q"))
	if a.edit == nil || a.edit.state.Text != "abq" {
		t.Fatal(a.edit)
	}
	cmd := a.onKey(tea.KeyPressMsg{Code: 'q', Text: "q", Mod: tea.ModCtrl})
	if cmd == nil {
		t.Fatal("ctrl+q ignoré")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("pas une sortie")
	}
	b := editApp(t)
	b.running = true
	if cmd := b.onKey(keyMsg("q")); cmd != nil {
		t.Fatal("q quitte pendant une commande")
	}
	if !b.running {
		t.Fatal("q a annulé la commande")
	}
	cmd = b.onKey(tea.KeyPressMsg{Code: 'q', Mod: tea.ModCtrl})
	if cmd == nil {
		t.Fatal("ctrl+q ignoré pendant une commande")
	}
	if _, ok := cmd().(tea.QuitMsg); !ok {
		t.Fatal("pas une sortie")
	}
	if b.running {
		t.Fatal("la commande court encore")
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
	title := stylePrefix(t, "#ffa657", true)
	plainTitle := stylePrefix(t, "#ffa657", false)
	body := stylePrefix(t, "#e6edf3", false)
	for _, line := range strings.Split(frame, "\n") {
		if lipgloss.Width(line) != 80 {
			t.Fatalf("largeur %d pour %q", lipgloss.Width(line), line)
		}
		if strings.Contains(line, "Rubriques") && !strings.Contains(line, title) {
			t.Fatal("titre de liste sans couleur")
		}
		if strings.Contains(line, "high-tui") && !strings.Contains(line, plainTitle) {
			t.Fatal("titre de fenêtre sans couleur")
		}
		if strings.Contains(line, "Jeton pro") && strings.Contains(line, title) {
			t.Fatal("ligne d'item colorée comme un titre")
		}
		if strings.Contains(line, "127.0.0.1") && strings.Contains(line, title) {
			t.Fatal("valeur colorée comme un titre")
		}
		if strings.Contains(line, "Rubriques") && !strings.HasPrefix(visible(line), " ") {
			t.Fatal("marge gauche absente")
		}
		if strings.Contains(line, "Jeton pro") && !strings.Contains(line, body) {
			t.Fatal("couleur du texte absente")
		}
	}
}

func stylePrefix(t *testing.T, fg string, bold bool) string {
	t.Helper()
	s := lipgloss.NewStyle().
		Foreground(lipgloss.Color(fg)).
		Background(lipgloss.Color("#161b22")).
		Inline(true)
	if bold {
		s = s.Bold(true)
	}
	out := s.Render("X")
	i := strings.Index(out, "X")
	if i < 1 {
		t.Fatal(out)
	}
	return out[:i]
}

func visible(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); {
		if s[i] == 0x1b {
			j := i + 1
			if j < len(s) && s[j] == '[' {
				j++
				for j < len(s) && ((s[j] >= '0' && s[j] <= '9') || s[j] == ';') {
					j++
				}
				if j < len(s) {
					j++
				}
				i = j
				continue
			}
		}
		b.WriteByte(s[i])
		i++
	}
	return b.String()
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
