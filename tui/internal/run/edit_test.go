package run

import "testing"

func TestShiftEnterInsertsNewline(t *testing.T) {
	ed, res := ApplyEdit(Edit{Text: "ab", Cursor: 1}, "shift+enter", true)
	if res != EditNone || ed.Text != "a\nb" || ed.Cursor != 2 {
		t.Fatalf("%q %d %v", ed.Text, ed.Cursor, res)
	}
}

func TestShiftEnterIgnoredOnInput(t *testing.T) {
	ed, res := ApplyEdit(Edit{Text: "ab", Cursor: 1}, "shift+enter", false)
	if res != EditNone || ed.Text != "ab" {
		t.Fatalf("%q %v", ed.Text, res)
	}
}

func TestEnterCommits(t *testing.T) {
	_, res := ApplyEdit(Edit{Text: "ab", Cursor: 2}, "enter", true)
	if res != EditCommit {
		t.Fatal(res)
	}
}

func TestAbsentKeyDoesNothing(t *testing.T) {
	ed, res := ApplyEdit(Edit{Text: "ab", Cursor: 1}, "", true)
	if res != EditNone || ed.Text != "ab" || ed.Cursor != 1 {
		t.Fatal(ed, res)
	}
}
