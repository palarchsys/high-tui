// Package run is the live behaviour of the executable: files, actions, editing.
package run

// Edit is the in-memory text being changed.
type Edit struct {
	Text   string
	Cursor int
}

// EditResult says what the key did.
type EditResult int

const (
	EditNone EditResult = iota
	EditCommit
	EditCancel
)

// ApplyEdit handles one key inside an editor.
// multiline is true for a textarea. Enter and Shift+Enter insert a newline then.
// Enter on a single-line input commits. Ctrl+S commits. Ctrl+C and Escape cancel.
func ApplyEdit(ed Edit, key string, multiline bool) (Edit, EditResult) {
	rs := []rune(ed.Text)
	if ed.Cursor < 0 {
		ed.Cursor = 0
	}
	if ed.Cursor > len(rs) {
		ed.Cursor = len(rs)
	}
	switch key {
	case "shift+enter":
		if !multiline {
			return ed, EditNone
		}
		return insertNewline(ed, rs), EditNone
	case "enter":
		if multiline {
			return insertNewline(ed, rs), EditNone
		}
		return ed, EditCommit
	case "ctrl+s":
		return ed, EditCommit
	case "esc", "ctrl+c":
		return ed, EditCancel
	case "left":
		if ed.Cursor > 0 {
			ed.Cursor--
		}
	case "right":
		if ed.Cursor < len(rs) {
			ed.Cursor++
		}
	case "backspace", "ctrl+h":
		if ed.Cursor > 0 {
			rs = append(rs[:ed.Cursor-1], rs[ed.Cursor:]...)
			ed.Cursor--
			ed.Text = string(rs)
		}
	case "delete":
		if ed.Cursor < len(rs) {
			rs = append(rs[:ed.Cursor], rs[ed.Cursor+1:]...)
			ed.Text = string(rs)
		}
	case "up":
		ed.Cursor = moveLine(rs, ed.Cursor, -1)
	case "down":
		ed.Cursor = moveLine(rs, ed.Cursor, 1)
	default:
		if len([]rune(key)) == 1 && key != "space" {
			rs = insert(rs, ed.Cursor, []rune(key)[0])
			ed.Cursor++
			ed.Text = string(rs)
		} else if key == "space" {
			rs = insert(rs, ed.Cursor, ' ')
			ed.Cursor++
			ed.Text = string(rs)
		}
	}
	return ed, EditNone
}

func insertNewline(ed Edit, rs []rune) Edit {
	rs = insert(rs, ed.Cursor, '\n')
	ed.Cursor++
	ed.Text = string(rs)
	return ed
}

func insert(rs []rune, at int, r rune) []rune {
	out := make([]rune, 0, len(rs)+1)
	out = append(out, rs[:at]...)
	out = append(out, r)
	out = append(out, rs[at:]...)
	return out
}

func moveLine(rs []rune, cursor, dir int) int {
	line, col := 0, 0
	for i := 0; i < cursor && i < len(rs); i++ {
		if rs[i] == '\n' {
			line++
			col = 0
			continue
		}
		col++
	}
	target := line + dir
	if target < 0 {
		return cursor
	}
	pos := 0
	curLine := 0
	for pos < len(rs) && curLine < target {
		if rs[pos] == '\n' {
			curLine++
		}
		pos++
	}
	if curLine != target {
		return cursor
	}
	for i := 0; i < col && pos < len(rs) && rs[pos] != '\n'; i++ {
		pos++
	}
	return pos
}
