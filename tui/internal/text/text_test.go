package text

import "testing"

func TestLocate(t *testing.T) {
	line, col := Locate("ab", 10, 1)
	if line != 0 || col != 1 {
		t.Fatal(line, col)
	}
	line, col = Locate("a\nb", 10, 2)
	if line != 1 || col != 0 {
		t.Fatal(line, col)
	}
	line, col = Locate("alpha beta gamma", 10, 11)
	if line != 1 || col != 0 {
		t.Fatal(line, col)
	}
}

func TestWrapKeepsExactWord(t *testing.T) {
	rows := Wrap("alpha beta gamma", 10)
	if len(rows) != 2 || rows[0].Text != "alpha beta" || rows[1].Text != "gamma" {
		t.Fatalf("%#v", rows)
	}
}

func TestPageKeepsSourceLine(t *testing.T) {
	rows, max := Page("titre\nitem", 10, 2, 0)
	if max != 0 || len(rows) != 2 || rows[0].From != 0 || rows[1].From != 1 {
		t.Fatalf("%#v %d", rows, max)
	}
	if rows[0].Text != "titre     " || rows[1].Text != "item      " {
		t.Fatalf("%#v", rows)
	}
}

func TestShiftLineIsTwoRows(t *testing.T) {
	rows := Wrap("une\ndeux", 10)
	if len(rows) != 2 || rows[0].Text != "une" || rows[1].Text != "deux" {
		t.Fatalf("%#v", rows)
	}
}
