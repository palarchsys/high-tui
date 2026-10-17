package html

import "testing"

func TestAbsentActionDoesNothing(t *testing.T) {
	doc, err := Parse(`<window id="w" label="W">
	  <main id="r" lines="1">
	    <button id="b" label="Enregistrer" enter="save"></button>
	  </main>
	</window>`, "t.html")
	if err != nil {
		t.Fatal(err)
	}
	win, err := Window(doc, "t.html")
	if err != nil {
		t.Fatal(err)
	}
	zone := win.Rows[0].Zones[0]
	if zone.Focus.Kind != "" {
		t.Fatalf("focus déclenché : %q", zone.Focus.Kind)
	}
	if zone.Fr != 1 || zone.Readonly || zone.File != "" || zone.Value != "" {
		t.Fatalf("défauts %#v", zone)
	}
	if zone.Enter.Kind != "save" {
		t.Fatal(zone.Enter)
	}
}

func TestEmptyAttribute(t *testing.T) {
	_, err := Parse(`<window id="w"><main id="r"><button id="b" label="" enter="save"></button></main></window>`, "t.html")
	if err == nil || !stringsContains(err.Error(), "attribut vide") {
		t.Fatal(err)
	}
}

func stringsContains(s, part string) bool {
	return len(s) >= len(part) && (func() bool {
		for i := 0; i+len(part) <= len(s); i++ {
			if s[i:i+len(part)] == part {
				return true
			}
		}
		return false
	})()
}
