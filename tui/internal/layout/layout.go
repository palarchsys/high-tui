// Package layout places a window inside a terminal width.
package layout

import (
	"os"
	"strconv"

	"high-tui/internal/model"
)

// Box is one placed zone. X and Y are the top-left cell of the border.
type Box struct {
	Zone   model.Zone
	X, Y   int
	OuterW int
	OuterH int
	TextW  int
	Lines  int
	Gap    int
}

// Sub is a grid placed inside a cell. Its bands start at the top of that cell.
type Sub struct {
	Bands         []Band
	Width, Height int
}

// Slot is one cell of a placed row.
type Slot struct {
	Box           *Box
	Sub           *Sub
	Width, Height int
	Gap           int
}

// Band is one placed row. InsetL and InsetR are the side margins, in cells.
type Band struct {
	Y              int
	Height         int
	InsetL, InsetR int
	Cells          []Slot
}

// Placed is a window ready to draw.
type Placed struct {
	Window model.Window
	Width  int
	Height int
	Inner  int
	Bands  []Band
}

// EachBox visits placed zones in document order.
func (p Placed) EachBox(fn func(*Box)) {
	var walk func([]Band)
	walk = func(bands []Band) {
		for i := range bands {
			for j := range bands[i].Cells {
				cell := &bands[i].Cells[j]
				if cell.Box != nil {
					fn(cell.Box)
				}
				if cell.Sub != nil {
					walk(cell.Sub.Bands)
				}
			}
		}
	}
	walk(p.Bands)
}

// Widths splits inner across fr weights. The last column takes the remainder.
func Widths(inner int, frs []int, gap int, where string) ([]int, error) {
	n := len(frs)
	if n < 1 {
		return nil, model.Errf("rangée sans colonne (%s)", where)
	}
	if n == 1 {
		if inner < 1 {
			return nil, model.Errf("largeur insuffisante (%s)", where)
		}
		return []int{inner}, nil
	}
	usable := inner - (n-1)*gap
	total := 0
	for _, fr := range frs {
		total += fr
	}
	if total < 1 || usable < n {
		return nil, model.Errf("largeur insuffisante (%s)", where)
	}
	widths := make([]int, 0, n)
	used := 0
	for _, fr := range frs[:n-1] {
		width := (usable * fr) / total
		if width < 1 {
			return nil, model.Errf("colonne trop étroite (%s)", where)
		}
		widths = append(widths, width)
		used += width
	}
	last := usable - used
	if last < 1 {
		return nil, model.Errf("dernière colonne trop étroite (%s)", where)
	}
	return append(widths, last), nil
}

// Fitted is min(css max-width, COLUMNS) when COLUMNS is a positive integer.
func Fitted(cssMax int, where string) (int, error) {
	raw := os.Getenv("COLUMNS")
	if raw == "" {
		return cssMax, nil
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 {
		return 0, model.Errf("COLUMNS invalide « %s » (%s)", raw, where)
	}
	if n < cssMax {
		return n, nil
	}
	return cssMax, nil
}

func textW(outer int, style model.Style) int {
	return outer - style.BorderWidth*2 - style.Pad[1] - style.Pad[3]
}

func outerH(style model.Style, lines int) int {
	between := (lines - 1) * (style.LineHeight - 1)
	return style.BorderWidth*2 + style.Pad[0] + lines + between + style.Pad[2]
}

func rowName(row model.GridRow) string {
	if row.ID != "" {
		return row.ID
	}
	return "sans nom"
}

func directInsets(chrome model.Style, row model.GridRow, root bool) (int, int) {
	left, right := 0, 0
	if root {
		left = chrome.Mar[3] + chrome.BorderWidth + chrome.Pad[3]
		right = chrome.Mar[1] + chrome.BorderWidth + chrome.Pad[1]
	}
	maxL, maxR := 0, 0
	for _, cell := range row.Cells {
		if cell.Zone == nil {
			continue
		}
		if cell.Zone.Base.Mar[3] > maxL {
			maxL = cell.Zone.Base.Mar[3]
		}
		if cell.Zone.Base.Mar[1] > maxR {
			maxR = cell.Zone.Base.Mar[1]
		}
	}
	return left + maxL, right + maxR
}

func columnGaps(row model.GridRow) []int {
	base := 0
	for _, cell := range row.Cells {
		if cell.Zone != nil && cell.Zone.Base.Gap > base {
			base = cell.Zone.Base.Gap
		}
	}
	if len(row.Cells) < 2 {
		return nil
	}
	gaps := make([]int, len(row.Cells)-1)
	for i := range gaps {
		g := base
		if row.Cells[i].Zone != nil && row.Cells[i].Zone.Base.Mar[1] > g {
			g = row.Cells[i].Zone.Base.Mar[1]
		}
		if row.Cells[i+1].Zone != nil && row.Cells[i+1].Zone.Base.Mar[3] > g {
			g = row.Cells[i+1].Zone.Base.Mar[3]
		}
		gaps[i] = g
	}
	return gaps
}

// offsetBoxes moves every zone inside bands by dx, dy. Nested grids keep the same shift.
func offsetBoxes(bands []Band, dx, dy int) {
	if dx == 0 && dy == 0 {
		return
	}
	for i := range bands {
		for j := range bands[i].Cells {
			cell := &bands[i].Cells[j]
			if cell.Box != nil {
				cell.Box.X += dx
				cell.Box.Y += dy
			}
			if cell.Sub != nil {
				offsetBoxes(cell.Sub.Bands, dx, dy)
			}
		}
	}
}

func growSub(sub *Sub, delta int) {
	if sub == nil || delta <= 0 {
		return
	}
	if len(sub.Bands) == 0 {
		sub.Height += delta
		return
	}
	last := &sub.Bands[len(sub.Bands)-1]
	last.Height += delta
	for i := range last.Cells {
		last.Cells[i].Height += delta
		if last.Cells[i].Box != nil {
			last.Cells[i].Box.OuterH += delta
		}
		if last.Cells[i].Sub != nil {
			growSub(last.Cells[i].Sub, delta)
		}
	}
	sub.Height += delta
}

// Place lays win out at the given width. Width 0 uses Fitted.
func Place(win model.Window, width int) (Placed, error) {
	if width == 0 {
		n, err := Fitted(win.Style.MaxWidth, win.ID)
		if err != nil {
			return Placed{}, err
		}
		width = n
	}
	chrome := win.Style.Pad[1] + win.Style.Pad[3] + win.Style.BorderWidth*2 + win.Style.Mar[1] + win.Style.Mar[3]
	bands, height, err := placeGrid(win.Grid, width, win.Style, true)
	if err != nil {
		return Placed{}, err
	}
	return Placed{Window: win, Width: width, Height: height, Inner: width - chrome, Bands: bands}, nil
}

func placeGrid(grid model.Grid, width int, chrome model.Style, root bool) ([]Band, int, error) {
	if len(grid.Rows) == 0 {
		where := grid.ID
		if where == "" {
			where = "grille"
		}
		return nil, 0, model.Errf("grille sans rangée (%s)", where)
	}
	y := 0
	var prevMB *int
	var bands []Band
	for _, row := range grid.Rows {
		band, next, mb, err := placeRow(row, width, chrome, root, y, prevMB)
		if err != nil {
			return nil, 0, err
		}
		bands = append(bands, band)
		y = next
		kept := mb
		prevMB = &kept
	}
	if prevMB != nil {
		y += *prevMB
	}
	return bands, y, nil
}

func placeRow(row model.GridRow, width int, chrome model.Style, root bool, y int, prevMB *int) (Band, int, int, error) {
	where := rowName(row)
	if len(row.Cells) == 0 {
		return Band{}, 0, 0, model.Errf("rangée sans colonne (%s)", where)
	}
	insetL, insetR := directInsets(chrome, row, root)
	content := width - insetL - insetR
	if content < 1 {
		return Band{}, 0, 0, model.Errf("largeur insuffisante (%s)", where)
	}
	frs := make([]int, len(row.Cells))
	for i, cell := range row.Cells {
		fr := cell.Fr
		if fr < 1 {
			fr = 1
		}
		frs[i] = fr
		if (cell.Zone == nil) == (cell.Grid == nil) {
			return Band{}, 0, 0, model.Errf("cellule invalide (%s)", where)
		}
	}
	gaps := columnGaps(row)
	sumGap := 0
	same := true
	for i, g := range gaps {
		sumGap += g
		if i > 0 && g != gaps[0] {
			same = false
		}
	}
	var widths []int
	var err error
	if same {
		gap := 0
		if len(gaps) > 0 {
			gap = gaps[0]
		}
		widths, err = Widths(content, frs, gap, where)
	} else {
		widths, err = Widths(content-sumGap, frs, 0, where)
	}
	if err != nil {
		return Band{}, 0, 0, err
	}
	type measured struct {
		box             *Box
		sub             *Sub
		top, bottom, in int
	}
	ms := make([]measured, len(row.Cells))
	rowH := 0
	lines := row.Lines
	if lines < 1 {
		lines = 1
	}
	for i, cell := range row.Cells {
		if cell.Zone != nil {
			tw := textW(widths[i], cell.Zone.Base)
			if tw < 1 {
				return Band{}, 0, 0, model.Errf("zone trop étroite : %s", cell.Zone.ID)
			}
			oh := outerH(cell.Zone.Base, lines)
			ms[i].box = &Box{
				Zone: *cell.Zone, OuterW: widths[i], OuterH: oh, TextW: tw, Lines: lines,
			}
			ms[i].in = oh
			ms[i].top = cell.Zone.Base.Mar[0]
			ms[i].bottom = cell.Zone.Base.Mar[2]
		} else {
			nb, nh, err := placeGrid(*cell.Grid, widths[i], model.Style{}, false)
			if err != nil {
				return Band{}, 0, 0, err
			}
			top, bottom := 0, 0
			if len(nb) > 0 {
				top = nb[0].Y
				last := nb[len(nb)-1]
				bottom = nh - (last.Y + last.Height)
			}
			for bi := range nb {
				nb[bi].Y -= top
			}
			offsetBoxes(nb, 0, -top)
			contentH := nh - top - bottom
			ms[i].sub = &Sub{Bands: nb, Width: widths[i], Height: contentH}
			ms[i].in = contentH
			ms[i].top = top
			ms[i].bottom = bottom
		}
		if ms[i].in > rowH {
			rowH = ms[i].in
		}
	}
	for i := range ms {
		box := ms[i].box
		if box == nil {
			continue
		}
		oh := box.OuterH
		for oh < rowH {
			box.Lines++
			next := outerH(box.Zone.Base, box.Lines)
			if next <= oh {
				oh = rowH
				break
			}
			oh = next
		}
		box.OuterH = oh
		if oh > rowH {
			rowH = oh
		}
	}
	for i := range ms {
		if ms[i].sub == nil {
			continue
		}
		if ms[i].sub.Height < rowH {
			growSub(ms[i].sub, rowH-ms[i].sub.Height)
		}
		if ms[i].sub.Height > rowH {
			rowH = ms[i].sub.Height
		}
	}
	for i := range ms {
		if ms[i].box != nil && ms[i].box.OuterH < rowH {
			ms[i].box.OuterH = rowH
		}
		if ms[i].sub != nil && ms[i].sub.Height < rowH {
			growSub(ms[i].sub, rowH-ms[i].sub.Height)
		}
	}
	mt, mb := 0, 0
	for i, item := range ms {
		if i == 0 || item.top > mt {
			mt = item.top
		}
		if i == 0 || item.bottom > mb {
			mb = item.bottom
		}
	}
	if prevMB == nil {
		y += mt
	} else if *prevMB > mt {
		y += *prevMB
	} else {
		y += mt
	}
	band := Band{Y: y, Height: rowH, InsetL: insetL, InsetR: insetR}
	x := insetL
	for i := range ms {
		g := 0
		if i < len(gaps) {
			g = gaps[i]
		}
		if ms[i].box != nil {
			ms[i].box.X = x
			ms[i].box.Y = y
			ms[i].box.Gap = g
			ms[i].box.OuterW = widths[i]
			if ms[i].box.OuterH < rowH {
				ms[i].box.OuterH = rowH
			}
		}
		if ms[i].sub != nil {
			offsetBoxes(ms[i].sub.Bands, x, y)
		}
		band.Cells = append(band.Cells, Slot{
			Box: ms[i].box, Sub: ms[i].sub, Width: widths[i], Height: rowH, Gap: g,
		})
		x += widths[i] + g
	}
	y += rowH
	used := insetL + insetR
	for i, w := range widths {
		used += w
		if i < len(gaps) {
			used += gaps[i]
		}
	}
	if used != width {
		return Band{}, 0, 0, model.Errf("la rangée %s ne fait pas 100 %% (%d au lieu de %d)", where, used, width)
	}
	return band, y, mb, nil
}
