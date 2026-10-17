package css

import "testing"

func TestUnknownProperty(t *testing.T) {
	_, err := Parse("window { foo: 1; }", "theme.css")
	if err == nil || !contains(err.Error(), "propriété inconnue") {
		t.Fatal(err)
	}
}

func TestCascadeFocus(t *testing.T) {
	rules, err := Parse(`
window { color: #e6edf3; background: #0d1117; }
nav { color: #e6edf3; background: #161b22; border: 1ch solid #30363d; }
nav:focus { border-color: #58a6ff; font-weight: bold; }
`, "main.css")
	if err != nil {
		t.Fatal(err)
	}
	inherit, err := Resolve(rules, "window", nil, "main", false, Blank())
	if err != nil {
		t.Fatal(err)
	}
	base, err := Resolve(rules, "nav", nil, "nav-rubriques", false, inherit)
	if err != nil {
		t.Fatal(err)
	}
	if base.BorderColor != "#30363d" || base.Bold {
		t.Fatalf("base %#v", base)
	}
	focus, err := Resolve(rules, "nav", nil, "nav-rubriques", true, inherit)
	if err != nil {
		t.Fatal(err)
	}
	if focus.BorderColor != "#58a6ff" || !focus.Bold {
		t.Fatalf("focus %#v", focus)
	}
}

func TestTitleColorStaysOnFocus(t *testing.T) {
	rules, err := Parse(`
window { color: #e6edf3; background: #0d1117; }
nav, header, input { title-color: #ffa657; }
nav:focus { border-color: #58a6ff; font-weight: bold; }
`, "main.css")
	if err != nil {
		t.Fatal(err)
	}
	inherit, err := Resolve(rules, "window", nil, "main", false, Blank())
	if err != nil {
		t.Fatal(err)
	}
	base, err := Resolve(rules, "nav", nil, "nav-rubriques", false, inherit)
	if err != nil {
		t.Fatal(err)
	}
	if base.TitleColor != "#ffa657" {
		t.Fatal(base.TitleColor)
	}
	focus, err := Resolve(rules, "nav", nil, "nav-rubriques", true, inherit)
	if err != nil {
		t.Fatal(err)
	}
	if focus.TitleColor != "#ffa657" || focus.BorderColor != "#58a6ff" {
		t.Fatalf("focus %#v", focus)
	}
	area, err := Resolve(rules, "textarea", nil, "corps", false, inherit)
	if err != nil {
		t.Fatal(err)
	}
	if area.TitleColor != "" {
		t.Fatal(area.TitleColor)
	}
}

func contains(s, part string) bool {
	return len(s) >= len(part) && (s == part || len(part) == 0 || indexOf(s, part) >= 0)
}

func indexOf(s, part string) int {
	for i := 0; i+len(part) <= len(s); i++ {
		if s[i:i+len(part)] == part {
			return i
		}
	}
	return -1
}
