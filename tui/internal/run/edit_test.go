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

func TestEnterInsertsNewline(t *testing.T) {
	ed, res := ApplyEdit(Edit{Text: "ab", Cursor: 1}, "enter", true)
	if res != EditNone || ed.Text != "a\nb" || ed.Cursor != 2 {
		t.Fatalf("%q %d %v", ed.Text, ed.Cursor, res)
	}
}

func TestEnterCommitsInput(t *testing.T) {
	ed, res := ApplyEdit(Edit{Text: "ab", Cursor: 2}, "enter", false)
	if res != EditCommit || ed.Text != "ab" {
		t.Fatal(ed.Text, res)
	}
}

func TestCtrlSCommits(t *testing.T) {
	_, res := ApplyEdit(Edit{Text: "ab", Cursor: 2}, "ctrl+s", true)
	if res != EditCommit {
		t.Fatal(res)
	}
}

func TestCtrlCCancels(t *testing.T) {
	ed, res := ApplyEdit(Edit{Text: "ab", Cursor: 1}, "ctrl+c", true)
	if res != EditCancel || ed.Text != "ab" {
		t.Fatal(ed.Text, res)
	}
}

func TestBackspaceAndDelete(t *testing.T) {
	ed, res := ApplyEdit(Edit{Text: "Aa\nBb", Cursor: 4}, "ctrl+h", true)
	if res != EditNone || ed.Text != "Aa\nb" || ed.Cursor != 3 {
		t.Fatalf("%q %d %v", ed.Text, ed.Cursor, res)
	}
	ed, res = ApplyEdit(Edit{Text: "Aa\nBb", Cursor: 4}, "backspace", true)
	if res != EditNone || ed.Text != "Aa\nb" || ed.Cursor != 3 {
		t.Fatalf("%q %d %v", ed.Text, ed.Cursor, res)
	}
	ed, res = ApplyEdit(Edit{Text: "Aa\nBb", Cursor: 4}, "delete", true)
	if res != EditNone || ed.Text != "Aa\nB" || ed.Cursor != 4 {
		t.Fatalf("%q %d %v", ed.Text, ed.Cursor, res)
	}
	ed, res = ApplyEdit(Edit{Text: "ab", Cursor: 0}, "backspace", true)
	if res != EditNone || ed.Text != "ab" || ed.Cursor != 0 {
		t.Fatalf("%q %d", ed.Text, ed.Cursor)
	}
}

func TestAbsentKeyDoesNothing(t *testing.T) {
	ed, res := ApplyEdit(Edit{Text: "ab", Cursor: 1}, "", true)
	if res != EditNone || ed.Text != "ab" || ed.Cursor != 1 {
		t.Fatal(ed, res)
	}
}
