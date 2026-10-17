package items

import (
	"strings"
	"testing"
)

func TestParseInlineAndFile(t *testing.T) {
	raw := `[
	  {"label":"Réseau","next":[{"label":"Passage","value":"127.0.0.1","file":"samples/passage.env","name":"HOTE"}]},
	  {"label":"Notes","next":"samples/fiches-notes.json"}
	]`
	items, err := Parse(raw, "nav")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 || items[0].Label != "Réseau" || len(items[0].NextItems) != 1 {
		t.Fatalf("%#v", items)
	}
	if items[0].NextItems[0].File != "samples/passage.env" || items[0].NextItems[0].Value != "127.0.0.1" {
		t.Fatal(items[0].NextItems[0])
	}
	if items[1].Next != "samples/fiches-notes.json" {
		t.Fatal(items[1].Next)
	}
}

func TestParseJSONL(t *testing.T) {
	raw := "{\"label\":\"Un\"}\n{\"label\":\"Deux\",\"enter\":{\"kind\":\"edit\",\"target\":\"corps\"}}\n"
	items, err := Parse(raw, "nav")
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 2 || items[1].Enter.Kind != "edit" || items[1].Enter.Target != "corps" {
		t.Fatal(items)
	}
}

func TestRejects(t *testing.T) {
	cases := []struct {
		raw  string
		want string
	}{
		{`[{"label":""}]`, "label manquant"},
		{`[{"label":"A","id":"x"}]`, "champ inconnu"},
		{`[{"label":"A","next":"cmd:"}]`, "commande illégale"},
		{`not json`, "JSON invalide"},
		{`[{"label":"A","next":1}]`, "JSON invalide"},
	}
	for _, tc := range cases {
		_, err := Parse(tc.raw, "nav")
		if err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("%s → %v", tc.raw, err)
		}
	}
}
