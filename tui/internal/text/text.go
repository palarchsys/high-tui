// Package text wraps, clips, and draws a scrollbar inside a fixed box.
package text

import "strings"

// Row is one visual line and the source line it came from.
type Row struct {
	Text string
	From int
}

func runes(s string) []rune { return []rune(s) }

// ClipPad left-aligns s in width cells.
func ClipPad(s string, width int) string {
	r := runes(s)
	if width < 0 {
		width = 0
	}
	if len(r) > width {
		r = r[:width]
	}
	if len(r) < width {
		r = append(r, []rune(strings.Repeat(" ", width-len(r)))...)
	}
	return string(r)
}

func wrapOne(text string, width int) []string {
	text = strings.ReplaceAll(text, "\t", " ")
	text = strings.ReplaceAll(text, "\r", "")
	if width < 1 {
		width = 1
	}
	var out []string
	for len(runes(text)) > width {
		rs := runes(text)
		if width < len(rs) && rs[width] == ' ' {
			out = append(out, string(rs[:width]))
			text = string(rs[width+1:])
			continue
		}
		chunk := rs[:width]
		cut := -1
		for i := len(chunk) - 1; i >= 0; i-- {
			if chunk[i] == ' ' {
				cut = i
				break
			}
		}
		if cut > 0 {
			out = append(out, string(chunk[:cut]))
			text = string(rs[cut+1:])
			continue
		}
		out = append(out, string(chunk))
		text = string(rs[width:])
	}
	out = append(out, text)
	return out
}

// Wrap splits logical lines. Each visual row remembers its source line.
func Wrap(text string, width int) []Row {
	if text == "" {
		return []Row{{Text: "", From: 0}}
	}
	parts := strings.Split(text, "\n")
	var rows []Row
	for src, line := range parts {
		for _, piece := range wrapOne(line, width) {
			rows = append(rows, Row{Text: piece, From: src})
		}
	}
	return rows
}

// Bar is the scrollbar glyph for one visual row.
func Bar(row, visible, total, scroll int) string {
	if total <= visible {
		return ""
	}
	thumb := visible * visible / total
	if thumb < 1 {
		thumb = 1
	}
	if thumb > visible {
		thumb = visible
	}
	maxScroll := total - visible
	top := 0
	if maxScroll > 0 {
		top = scroll * (visible - thumb) / maxScroll
	}
	if row >= top && row < top+thumb {
		return "█"
	}
	return "│"
}

// Page returns exactly lines rows of width textW, with a scrollbar when needed.
// A row past the text has From -1. The visible text keeps its source line.
func Page(text string, textW, lines, scroll int) ([]Row, int) {
	if lines < 1 {
		lines = 1
	}
	width := textW
	rows := Wrap(text, width)
	bar := false
	if len(rows) > lines && textW >= 2 {
		width = textW - 1
		rows = Wrap(text, width)
		bar = true
	}
	maxScroll := 0
	if len(rows) > lines {
		maxScroll = len(rows) - lines
	}
	if scroll < 0 {
		scroll = 0
	}
	if scroll > maxScroll {
		scroll = maxScroll
	}
	out := make([]Row, lines)
	for i := 0; i < lines; i++ {
		piece := ""
		from := -1
		if scroll+i < len(rows) {
			piece = rows[scroll+i].Text
			from = rows[scroll+i].From
		}
		piece = ClipPad(piece, width)
		if bar {
			piece += Bar(i, lines, len(rows), scroll)
		}
		out[i] = Row{Text: piece, From: from}
	}
	return out, maxScroll
}

// View returns exactly lines strings of width textW, with a scrollbar when needed.
func View(text string, textW, lines, scroll int) ([]string, int) {
	rows, maxScroll := Page(text, textW, lines, scroll)
	out := make([]string, len(rows))
	for i, row := range rows {
		out[i] = row.Text
	}
	return out, maxScroll
}

// Locate maps a rune cursor to a visual line and column, using the same breaks as Wrap.
func Locate(text string, width, cursor int) (int, int) {
	text = strings.ReplaceAll(text, "\t", " ")
	text = strings.ReplaceAll(text, "\r", "")
	if width < 1 {
		width = 1
	}
	rs := []rune(text)
	if cursor < 0 {
		cursor = 0
	}
	if cursor > len(rs) {
		cursor = len(rs)
	}
	line := 0
	i := 0
	for i < len(rs) {
		end := i
		for end < len(rs) && rs[end] != '\n' {
			end++
		}
		pos := i
		for pos < end {
			shown, adv := consume(rs[pos:end], width)
			next := pos + adv
			if cursor < next {
				col := cursor - pos
				if col > shown {
					col = shown
				}
				return line, col
			}
			if next >= end {
				break
			}
			line++
			pos = next
		}
		if cursor <= end {
			return line, cursor - pos
		}
		if end < len(rs) && rs[end] == '\n' {
			if cursor == end {
				return line, end - pos
			}
			i = end + 1
			line++
			continue
		}
		break
	}
	return line, 0
}

func consume(rs []rune, width int) (shown, advance int) {
	if len(rs) <= width {
		return len(rs), len(rs)
	}
	if rs[width] == ' ' {
		return width, width + 1
	}
	cut := -1
	for i := width - 1; i >= 0; i-- {
		if rs[i] == ' ' {
			cut = i
			break
		}
	}
	if cut > 0 {
		return cut, cut + 1
	}
	return width, width
}

// ListScroll keeps the selected source line inside the viewport.
// Title is source 0. Item i is source i+1.
func ListScroll(rows []Row, lines, sel int) int {
	if lines < 1 || len(rows) <= lines {
		return 0
	}
	source := sel + 1
	start, end := -1, -1
	for i, row := range rows {
		if row.From == source {
			if start < 0 {
				start = i
			}
			end = i
		}
	}
	if start < 0 {
		return 0
	}
	scroll := 0
	if end-start+1 >= lines {
		scroll = start
	} else if end >= lines {
		scroll = end - lines + 1
	}
	maxScroll := len(rows) - lines
	if scroll > maxScroll {
		scroll = maxScroll
	}
	if scroll < 0 {
		scroll = 0
	}
	return scroll
}
