package ui

import (
	"fmt"
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
	a.onKey(keyMsg("ctrl+s"))
	if a.edit != nil {
		t.Fatal("ctrl+s n'a pas fermé")
	}
	if a.val[key("main", "corps")] != "a\nb" {
		t.Fatal(a.val[key("main", "corps")])
	}
}

func TestEnterAddsLineWhenTerminalIsPlain(t *testing.T) {
	a := editApp(t)
	a.keyMode = keysPlain
	a.onKey(keyMsg("enter"))
	if a.edit == nil {
		t.Fatal("édition non ouverte")
	}
	if a.status == noShiftMsg {
		t.Fatal(a.status)
	}
	a.edit.state.Cursor = 1
	a.onKey(keyMsg("enter"))
	if a.edit == nil || a.edit.state.Text != "a\nb" {
		t.Fatal(a.edit)
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

func TestCtrlSAndCtrlC(t *testing.T) {
	dir := t.TempDir()
	a := editApp(t)
	a.root = dir
	a.file[key("main", "corps")] = "note.txt"
	a.onKey(keyMsg("enter"))
	a.edit.state.Cursor = 1
	a.onKey(keyMsg("enter"))
	a.onKey(keyMsg("ctrl+s"))
	if a.edit != nil || a.status != "Enregistré" {
		t.Fatal(a.edit, a.status)
	}
	got, err := os.ReadFile(filepath.Join(dir, "note.txt"))
	if err != nil || string(got) != "a\nb\n" {
		t.Fatalf("%q %v", got, err)
	}
	b := editApp(t)
	b.onKey(keyMsg("enter"))
	b.onKey(keyMsg("x"))
	if cmd := b.onKey(keyMsg("ctrl+c")); cmd != nil {
		t.Fatal("ctrl+c quitte l'édition")
	}
	if b.edit != nil || b.val[key("main", "corps")] != "ab" || b.status != "Annulé" {
		t.Fatal(b.edit, b.val[key("main", "corps")], b.status)
	}
	c := editApp(t)
	if cmd := c.onKey(keyMsg("ctrl+c")); cmd != nil {
		t.Fatal("ctrl+c quitte hors édition")
	}
	c.running = true
	if cmd := c.onKey(keyMsg("ctrl+c")); cmd != nil || !c.running {
		t.Fatal("ctrl+c quitte pendant une commande")
	}
}

func TestEditBorderUsesTitleColor(t *testing.T) {
	a := editApp(t)
	zone := &a.placed[0].Bands[0].Cells[0].Box.Zone
	zone.Base.TitleColor = "#ffa657"
	zone.Base.BorderColor = "#30363d"
	zone.FocusOn.TitleColor = "#ffa657"
	zone.FocusOn.BorderColor = "#58a6ff"
	before := a.frame()
	if strings.Contains(before, "255;166;87") {
		t.Fatal("couleur de titre avant l'édition")
	}
	a.onKey(keyMsg("enter"))
	frame := a.frame()
	if a.edit == nil {
		t.Fatal("édition non ouverte")
	}
	if !strings.Contains(frame, "255;166;87") {
		t.Fatal("contour d'édition sans la couleur des titres")
	}
	if strings.Contains(frame, "88;166;255") {
		t.Fatal("contour de focus conservé")
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

func TestGeneratedNotesWindow(t *testing.T) {
	root := filepath.Join("..", "..")
	app, err := newApp(root, 80)
	if err != nil {
		t.Fatal(err)
	}
	if err := app.open("notes.html"); err != nil {
		t.Fatal(err)
	}
	if app.winID() != "notes" {
		t.Fatal(app.winID())
	}
}

func TestCursorSitsOnNestedText(t *testing.T) {
	style := model.Style{
		Color: "#e6edf3", Background: "#161b22", BorderColor: "#30363d",
		BorderWidth: 1, LineHeight: 1, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 1, 1, 1},
	}
	header := model.Zone{Tag: "header", ID: "titre", Label: "Titre", Base: style, FocusOn: style}
	area := model.Zone{Tag: "textarea", ID: "corps", Base: style, FocusOn: style}
	inner := model.Grid{ID: "corps", Rows: []model.GridRow{{
		ID: "edition", Lines: 2, Cells: []model.Cell{{Fr: 1, Zone: &area}},
	}}}
	win := model.Window{
		ID: "main", File: "main.html",
		Style: model.Style{
			Color: "#e6edf3", Background: "#0d1117", LineHeight: 1, MaxWidth: 40,
		},
		Grid: model.Grid{Rows: []model.GridRow{
			{ID: "entete", Lines: 1, Cells: []model.Cell{{Fr: 1, Zone: &header}}},
			{ID: "page", Lines: 2, Cells: []model.Cell{{Fr: 1, Grid: &inner}}},
		}},
	}
	a := &app{
		root: t.TempDir(), wins: []model.Window{win}, origin: []model.Window{win}, width: 40,
		sel: map[string]int{}, scroll: map[string]int{}, val: map[string]string{},
		file: map[string]string{}, name: map[string]string{}, dirty: map[string]bool{},
		resolved: map[string]string{}, focusAt: map[string]int{}, status: "Prêt",
	}
	a.val[key("main", "corps")] = "HELLO"
	if err := a.layoutAll(); err != nil {
		t.Fatal(err)
	}
	if err := a.beginEdit("corps"); err != nil {
		t.Fatal(err)
	}
	a.edit.state.Cursor = 1
	a.reveal()
	assertCursorRune(t, a, 'E')
}

func TestCursorCellIsTheEditedCharacter(t *testing.T) {
	root := filepath.Join("..", "..")
	a, err := newApp(root, 80)
	if err != nil {
		t.Fatal(err)
	}
	a.val[key("main", "apercu")] = "Aa\nBb\nCc"
	focusZone(t, a, "apercu")
	if err := a.beginEdit("apercu"); err != nil {
		t.Fatal(err)
	}
	a.edit.state.Cursor = 3
	a.reveal()
	assertCursorRune(t, a, 'B')
	a.edit.state.Cursor = 4
	a.reveal()
	assertCursorRune(t, a, 'b')

	a.edit.state.Cursor = 2
	a.reveal()
	assertCursorAfter(t, a, 'a')

	a.edit = nil
	a.val[key("main", "valeur-hote")] = "QZQ"
	focusZone(t, a, "valeur-hote")
	if err := a.beginEdit("valeur-hote"); err != nil {
		t.Fatal(err)
	}
	a.edit.state.Cursor = 1
	a.reveal()
	assertCursorRune(t, a, 'Z')

	_, box, ok := a.zoneByID("apercu")
	if !ok {
		t.Fatal("apercu")
	}
	text := strings.Repeat("Q", box.TextW) + "Z"
	a.edit = nil
	a.val[key("main", "apercu")] = text
	focusZone(t, a, "apercu")
	if err := a.beginEdit("apercu"); err != nil {
		t.Fatal(err)
	}
	a.edit.state.Cursor = box.TextW
	a.reveal()
	assertCursorRune(t, a, 'Z')

	var long strings.Builder
	for i := 0; i < 12; i++ {
		fmt.Fprintf(&long, "L%d\n", i)
	}
	long.WriteString("ZZ")
	a.edit = nil
	a.val[key("main", "apercu")] = long.String()
	focusZone(t, a, "apercu")
	if err := a.beginEdit("apercu"); err != nil {
		t.Fatal(err)
	}
	a.edit.state.Cursor = strings.Index(long.String(), "Z")
	a.reveal()
	assertCursorRune(t, a, 'Z')
}

func TestBackspaceErasesThePreviousRune(t *testing.T) {
	root := filepath.Join("..", "..")
	a, err := newApp(root, 80)
	if err != nil {
		t.Fatal(err)
	}
	const body = "Aa\nBb\nCc"
	open := func() {
		t.Helper()
		a.edit = nil
		a.val[key("main", "apercu")] = body
		focusZone(t, a, "apercu")
		if err := a.beginEdit("apercu"); err != nil {
			t.Fatal(err)
		}
		a.edit.state.Cursor = 4
		a.reveal()
		assertCursorRune(t, a, 'b')
	}
	press := func(msg tea.KeyPressMsg) {
		t.Helper()
		open()
		a.onKey(msg)
		if a.edit == nil || a.edit.state.Text != "Aa\nb\nCc" || a.edit.state.Cursor != 3 {
			t.Fatalf("texte %q curseur %d", a.edit.state.Text, a.edit.state.Cursor)
		}
		assertCursorRune(t, a, 'b')
		assertCursorInside(t, a)
	}
	press(tea.KeyPressMsg{Code: tea.KeyBackspace})
	press(tea.KeyPressMsg{Code: 'h', Mod: tea.ModCtrl})
	press(tea.KeyPressMsg{Code: tea.KeyBackspace, Text: "\x7f"})

	open()
	a.edit.state.Cursor = len([]rune(body))
	a.reveal()
	assertCursorAfter(t, a, 'c')
	a.onKey(tea.KeyPressMsg{Code: tea.KeyBackspace})
	if a.edit.state.Text != "Aa\nBb\nC" || a.edit.state.Cursor != 7 {
		t.Fatalf("fin %q %d", a.edit.state.Text, a.edit.state.Cursor)
	}
	assertCursorAfter(t, a, 'C')
	assertCursorInside(t, a)

	open()
	a.edit.state.Cursor = 3
	a.reveal()
	assertCursorRune(t, a, 'B')
	a.onKey(tea.KeyPressMsg{Code: tea.KeyBackspace})
	if a.edit.state.Text != "AaBb\nCc" || a.edit.state.Cursor != 2 {
		t.Fatalf("jointure %q %d", a.edit.state.Text, a.edit.state.Cursor)
	}
	assertCursorRune(t, a, 'B')
	assertCursorInside(t, a)

	open()
	a.onKey(tea.KeyPressMsg{Code: tea.KeyDelete})
	if a.edit.state.Text != "Aa\nB\nCc" || a.edit.state.Cursor != 4 {
		t.Fatalf("suppr %q %d", a.edit.state.Text, a.edit.state.Cursor)
	}
	assertCursorAfter(t, a, 'B')
	assertCursorInside(t, a)

	_, box, ok := a.zoneByID("apercu")
	if !ok {
		t.Fatal("apercu")
	}
	wrapped := strings.Repeat("Q", box.TextW) + "Z"
	a.edit = nil
	a.val[key("main", "apercu")] = wrapped
	focusZone(t, a, "apercu")
	if err := a.beginEdit("apercu"); err != nil {
		t.Fatal(err)
	}
	a.edit.state.Cursor = box.TextW
	a.reveal()
	assertCursorRune(t, a, 'Z')
	a.onKey(tea.KeyPressMsg{Code: 'h', Mod: tea.ModCtrl})
	if a.edit.state.Text != strings.Repeat("Q", box.TextW-1)+"Z" {
		t.Fatal(a.edit.state.Text)
	}
	assertCursorRune(t, a, 'Z')
	assertCursorInside(t, a)

	var long strings.Builder
	for i := 0; i < 12; i++ {
		fmt.Fprintf(&long, "L%d\n", i)
	}
	long.WriteString("ZZ")
	a.edit = nil
	a.val[key("main", "apercu")] = long.String()
	focusZone(t, a, "apercu")
	if err := a.beginEdit("apercu"); err != nil {
		t.Fatal(err)
	}
	at := strings.Index(long.String(), "Z")
	a.edit.state.Cursor = at + 1
	a.reveal()
	assertCursorRune(t, a, 'Z')
	a.onKey(tea.KeyPressMsg{Code: tea.KeyBackspace})
	if !strings.HasSuffix(a.edit.state.Text, "Z") || a.edit.state.Cursor != at {
		t.Fatalf("défilement %q %d", a.edit.state.Text, a.edit.state.Cursor)
	}
	assertCursorRune(t, a, 'Z')
	assertCursorInside(t, a)

	idle := editApp(t)
	if cmd := idle.onKey(tea.KeyPressMsg{Code: 'h', Mod: tea.ModCtrl}); cmd != nil || idle.edit != nil {
		t.Fatal("ctrl+h hors édition")
	}
}

func focusZone(t *testing.T, a *app, id string) {
	t.Helper()
	for i, box := range a.boxes() {
		if box.Zone.ID == id {
			a.focus = i
			return
		}
	}
	t.Fatalf("zone absente : %s", id)
}

func assertCursorRune(t *testing.T, a *app, want rune) {
	t.Helper()
	frame := a.frame()
	x, y, ok := a.cursor()
	if !ok {
		t.Fatal("curseur absent")
	}
	got := runeAt(frame, x, y)
	if got != want {
		t.Fatalf("cellule %q en %d,%d, attendu %q\n%s", got, x, y, want, cursorMark(frame, x, y))
	}
}

func assertCursorAfter(t *testing.T, a *app, prev rune) {
	t.Helper()
	frame := a.frame()
	x, y, ok := a.cursor()
	if !ok {
		t.Fatal("curseur absent")
	}
	got := runeAt(frame, x, y)
	before := runeAt(frame, x-1, y)
	if before != prev || isBorder(got) || isBorder(before) {
		t.Fatalf("curseur %q après %q en %d,%d\n%s", got, before, x, y, cursorMark(frame, x, y))
	}
}

func runeAt(frame string, x, y int) rune {
	lines := strings.Split(frame, "\n")
	if y < 0 || y >= len(lines) || x < 0 {
		return 0
	}
	col := 0
	for _, r := range visible(lines[y]) {
		w := lipgloss.Width(string(r))
		if w < 1 {
			w = 1
		}
		if x >= col && x < col+w {
			return r
		}
		col += w
	}
	return 0
}

func assertCursorInside(t *testing.T, a *app) {
	t.Helper()
	if a.edit == nil {
		t.Fatal("édition absente")
	}
	_, box, ok := a.zoneByID(a.edit.zone)
	if !ok {
		t.Fatal("zone absente")
	}
	x, y, ok := a.cursor()
	if !ok {
		t.Fatal("curseur absent")
	}
	style := box.Zone.Base
	left := box.X + style.BorderWidth + style.Pad[3]
	top := box.Y + style.BorderWidth + style.Pad[0]
	width := box.TextW
	if y < top || y >= top+box.Lines || x < left || x >= left+width {
		t.Fatalf("curseur %d,%d hors du texte %d,%d %dx%d", x, y, left, top, width, box.Lines)
	}
	got := runeAt(a.frame(), x, y)
	if isBorder(got) || got == '█' {
		t.Fatalf("curseur sur %q en %d,%d", got, x, y)
	}
}

func isBorder(r rune) bool {
	switch r {
	case '─', '│', '┌', '┐', '└', '┘', '├', '┤', '┬', '┴', '┼', '╭', '╮', '╯', '╰':
		return true
	default:
		return false
	}
}

func cursorMark(frame string, x, y int) string {
	lines := strings.Split(frame, "\n")
	var b strings.Builder
	for i := y - 1; i <= y+1; i++ {
		if i < 0 || i >= len(lines) {
			continue
		}
		plain := visible(lines[i])
		fmt.Fprintf(&b, "%2d %s\n", i, plain)
		if i == y {
			fmt.Fprintf(&b, "   %s^\n", strings.Repeat(" ", x))
		}
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
		Tag: "textarea", ID: "corps", Readonly: false,
		Base: style, FocusOn: style,
		Enter: model.Action{Kind: "edit", Target: "corps"},
	}
	win := model.Window{
		ID: "main", File: "main.html",
		Style: model.Style{
			Color: "#e6edf3", Background: "#0d1117", BorderColor: "#e6edf3",
			LineHeight: 1, MaxWidth: 80,
		},
		Grid: model.Grid{ID: "main", Rows: []model.GridRow{{
			ID: "page", Lines: 4, Cells: []model.Cell{{Fr: 1, Zone: &zone}},
		}}},
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
	case "ctrl+s":
		return tea.KeyPressMsg{Code: 's', Mod: tea.ModCtrl}
	case "ctrl+c":
		return tea.KeyPressMsg{Code: 'c', Mod: tea.ModCtrl}
	default:
		return tea.KeyPressMsg{Text: name, Code: rune([]rune(name)[0])}
	}
}
