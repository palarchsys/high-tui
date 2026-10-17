package ui

import (
	"strings"

	"charm.land/lipgloss/v2"
	"high-tui/internal/layout"
	"high-tui/internal/text"
)

func paint(view viewState) string {
	if view.err != "" {
		return view.err
	}
	win := view.placed
	bg := lipgloss.NewStyle().
		Foreground(lipgloss.Color(win.Window.Style.Color)).
		Background(lipgloss.Color(win.Window.Style.Background))
	blank := fill(bg, win.Width, 1)
	var lines []string
	y := 0
	for _, band := range win.Bands {
		for y < band.Y {
			lines = append(lines, blank)
			y++
		}
		var cells []string
		if band.InsetL > 0 {
			cells = append(cells, fill(bg, band.InsetL, band.Height))
		}
		for _, box := range band.Boxes {
			cells = append(cells, paintBox(view, box))
			if box.Gap > 0 {
				cells = append(cells, fill(bg, box.Gap, band.Height))
			}
		}
		if band.InsetR > 0 {
			cells = append(cells, fill(bg, band.InsetR, band.Height))
		}
		row := lipgloss.JoinHorizontal(lipgloss.Top, cells...)
		for _, line := range strings.Split(row, "\n") {
			lines = append(lines, line)
		}
		y += band.Height
	}
	for y < win.Height {
		lines = append(lines, blank)
		y++
	}
	return strings.Join(lines, "\n")
}

func fill(style lipgloss.Style, width, height int) string {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	row := strings.Repeat(" ", width)
	rows := make([]string, height)
	for i := range rows {
		rows[i] = row
	}
	return style.Width(width).Render(strings.Join(rows, "\n"))
}

func paintBox(view viewState, box layout.Box) string {
	zone := box.Zone
	style := zone.Base
	focused := view.focusID == zone.ID
	if focused {
		style = zone.FocusOn
	}
	body := logical(view, box)
	scroll := 0
	if zone.Tag == "nav" || zone.Tag == "select" {
		width := box.TextW
		rows := text.Wrap(body, width)
		if len(rows) > box.Lines && box.TextW >= 2 {
			rows = text.Wrap(body, box.TextW-1)
		}
		scroll = text.ListScroll(rows, box.Lines, view.sel[key(view.winID, zone.ID)])
	} else {
		scroll = view.scroll[key(view.winID, zone.ID)]
	}
	page, _ := text.Page(body, box.TextW, box.Lines, scroll)
	if style.LineHeight > 1 {
		var spaced []text.Row
		for i, row := range page {
			spaced = append(spaced, row)
			if i < len(page)-1 {
				for n := 1; n < style.LineHeight; n++ {
					spaced = append(spaced, text.Row{From: -1})
				}
			}
		}
		page = spaced
	}
	title := style.TitleColor
	if title == "" {
		title = style.Color
	}
	lines := make([]string, len(page))
	for i, row := range page {
		fg := style.Color
		if titleZone(zone.Tag) && row.From == 0 {
			fg = title
		}
		line := lipgloss.NewStyle().
			Foreground(lipgloss.Color(fg)).
			Background(lipgloss.Color(style.Background)).
			Inline(true)
		if style.Bold {
			line = line.Bold(true)
		}
		lines[i] = line.Render(row.Text)
	}
	s := lipgloss.NewStyle().
		Background(lipgloss.Color(style.Background)).
		Width(box.OuterW)
	if style.Bold {
		s = s.Bold(true)
	}
	if style.BorderWidth == 1 {
		border := lipgloss.NormalBorder()
		if style.Radius >= 1 {
			border = lipgloss.RoundedBorder()
		}
		s = s.Border(border).BorderForeground(lipgloss.Color(style.BorderColor)).
			BorderBackground(lipgloss.Color(style.Background))
	}
	if style.Pad != ([4]int{}) {
		s = s.Padding(style.Pad[0], style.Pad[1], style.Pad[2], style.Pad[3])
	}
	return s.Render(strings.Join(lines, "\n"))
}

func titleZone(tag string) bool {
	switch tag {
	case "header", "nav", "select", "input":
		return true
	default:
		return false
	}
}

func logical(view viewState, box layout.Box) string {
	zone := box.Zone
	switch zone.Tag {
	case "header", "button", "command":
		return shown(view, zone.Label)
	case "status":
		return view.status
	case "input":
		name := view.name[key(view.winID, zone.ID)]
		if name == "" {
			name = zone.Name
		}
		return name + "\n" + view.val[key(view.winID, zone.ID)]
	case "textarea":
		return view.val[key(view.winID, zone.ID)]
	case "nav", "select":
		var b strings.Builder
		b.WriteString(shown(view, zone.Label))
		b.WriteByte('\n')
		sel := view.sel[key(view.winID, zone.ID)]
		for i, item := range zone.Items {
			mark := "  "
			if i == sel {
				mark = "› "
			}
			b.WriteString(mark)
			b.WriteString(shown(view, item.Label))
			if i < len(zone.Items)-1 {
				b.WriteByte('\n')
			}
		}
		return b.String()
	default:
		return shown(view, zone.Label)
	}
}

func shown(view viewState, raw string) string {
	if v, ok := view.resolved[raw]; ok {
		return v
	}
	return raw
}

type viewState struct {
	placed   layout.Placed
	winID    string
	focusID  string
	status   string
	sel      map[string]int
	scroll   map[string]int
	val      map[string]string
	name     map[string]string
	resolved map[string]string
	err      string
}

func key(win, id string) string { return win + "\x1f" + id }
