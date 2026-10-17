package build

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestJSONHrefLoadsWindows(t *testing.T) {
	root := t.TempDir()
	win := filepath.Join(root, "windows")
	if err := os.MkdirAll(win, 0o755); err != nil {
		t.Fatal(err)
	}
	writePage(t, win, "main", `<table id="main"><tr><td>
<nav id="liste" data-items="liste.json"><h2>L</h2><ul><li id="a">A</li></ul></nav>
</td></tr></table>`)
	writePage(t, win, "notes", `<table id="notes"><tr><td><header id="titre-notes"><h1>Notes</h1></header></td></tr></table>`)
	writePage(t, win, "autre", `<table id="autre"><tr><td><header id="titre-autre"><h1>Autre</h1></header></td></tr></table>`)
	writePage(t, win, "encore", `<table id="encore"><tr><td><header id="titre-encore"><h1>Encore</h1></header></td></tr></table>`)
	liste := `[
{"label":"Notes","next":"suite.json"},
{"label":"Direct","enter":{"kind":"window","href":"autre.html"}},
{"label":"Cmd","next":"cmd:false"}
]`
	suite := `[{"label":"Brouillon","enter":{"kind":"window","href":"notes.html"},"next":[{"label":"Encore","focus":{"kind":"window","href":"encore.html"}}]}]`
	if err := os.WriteFile(filepath.Join(root, "liste.json"), []byte(liste), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "suite.json"), []byte(suite), 0o644); err != nil {
		t.Fatal(err)
	}
	wins, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]bool{}
	for _, item := range wins {
		got[item.ID] = true
	}
	for _, id := range []string{"main", "notes", "autre", "encore"} {
		if !got[id] {
			t.Fatalf("fenêtre %s absente dans %v", id, got)
		}
	}
}

func TestJSONNextMissing(t *testing.T) {
	root := t.TempDir()
	win := filepath.Join(root, "windows")
	if err := os.MkdirAll(win, 0o755); err != nil {
		t.Fatal(err)
	}
	writePage(t, win, "main", `<table id="main"><tr><td>
<nav id="liste" data-items="liste.json"><h2>L</h2><ul><li id="a">A</li></ul></nav>
</td></tr></table>`)
	liste := `[{"label":"Notes","next":"absent.json"}]`
	if err := os.WriteFile(filepath.Join(root, "liste.json"), []byte(liste), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(root)
	if err == nil || !strings.Contains(err.Error(), "fichier introuvable") {
		t.Fatal(err)
	}
}

func TestJSONNextInvalid(t *testing.T) {
	root := t.TempDir()
	win := filepath.Join(root, "windows")
	if err := os.MkdirAll(win, 0o755); err != nil {
		t.Fatal(err)
	}
	writePage(t, win, "main", `<table id="main"><tr><td>
<nav id="liste" data-items="liste.json"><h2>L</h2><ul><li id="a">A</li></ul></nav>
</td></tr></table>`)
	if err := os.WriteFile(filepath.Join(root, "liste.json"), []byte(`[{"label":"Notes","next":"suite.json"}]`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "suite.json"), []byte(`{`), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := Load(root)
	if err == nil || !strings.Contains(err.Error(), "JSON invalide") {
		t.Fatal(err)
	}
}

func writePage(t *testing.T, dir, name, body string) {
	t.Helper()
	html := "<!DOCTYPE html>\n<html lang=\"fr\">\n<head>\n<meta charset=\"utf-8\">\n<title>" + name + "</title>\n<link rel=\"stylesheet\" href=\"" + name + ".css\">\n</head>\n<body>\n" + body + "\n</body>\n</html>\n"
	if err := os.WriteFile(filepath.Join(dir, name+".html"), []byte(html), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, name+".css"), []byte("body { max-width: 40ch; color: #e6edf3; background: #0d1117; }\n"), 0o644); err != nil {
		t.Fatal(err)
	}
}
