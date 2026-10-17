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

// Band is one placed row.
type Band struct {
	Row    model.Row
	Y      int
	Height int
	Boxes  []Box
}

// Placed is a window ready to draw.
type Placed struct {
	Window model.Window
	Width  int
	Height int
	Inner  int
	Bands  []Band
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

// Place lays win out at the given width. Width 0 uses Fitted.
func Place(win model.Window, width int) (Placed, error) {
	if width == 0 {
		n, err := Fitted(win.Style.MaxWidth, win.ID)
		if err != nil {
			return Placed{}, err
		}
		width = n
	}
	pad := win.Style.Pad
	inner := width - pad[1] - pad[3] - win.Style.BorderWidth*2
	placed := Placed{Window: win, Width: width, Inner: inner}
	y := 0
	var prevMB *int
	for _, row := range win.Rows {
		gap := row.Style.Gap
		frs := make([]int, len(row.Zones))
		for i, zone := range row.Zones {
			if zone.Base.Gap > gap {
				gap = zone.Base.Gap
			}
			frs[i] = zone.Fr
		}
		widths, err := Widths(inner, frs, gap, row.ID)
		if err != nil {
			return Placed{}, err
		}
		x := pad[3] + win.Style.BorderWidth
		band := Band{Row: row}
		maxH := 0
		for i, zone := range row.Zones {
			tw := textW(widths[i], zone.Base)
			if tw < 1 {
				return Placed{}, model.Errf("zone trop étroite : %s", zone.ID)
			}
			oh := outerH(zone.Base, row.Lines)
			if oh > maxH {
				maxH = oh
			}
			g := 0
			if i < len(row.Zones)-1 {
				g = gap
			}
			band.Boxes = append(band.Boxes, Box{
				Zone: zone, X: x, OuterW: widths[i], OuterH: oh, TextW: tw, Lines: row.Lines, Gap: g,
			})
			x += widths[i] + g
		}
		mt, mb := 0, 0
		for i, zone := range row.Zones {
			if i == 0 || zone.Base.Mar[0] > mt {
				mt = zone.Base.Mar[0]
			}
			if i == 0 || zone.Base.Mar[2] > mb {
				mb = zone.Base.Mar[2]
			}
		}
		if prevMB == nil {
			y += mt
		} else if *prevMB > mt {
			y += *prevMB
		} else {
			y += mt
		}
		band.Y = y
		band.Height = maxH
		for i := range band.Boxes {
			band.Boxes[i].Y = y
			band.Boxes[i].OuterH = maxH
		}
		y += maxH
		prev := mb
		prevMB = &prev
		used := 0
		for i, w := range widths {
			used += w
			if i > 0 {
				used += gap
			}
		}
		if used != inner {
			return Placed{}, model.Errf("la rangée %s ne fait pas 100 %% (%d au lieu de %d)", row.ID, used, inner)
		}
		placed.Bands = append(placed.Bands, band)
	}
	if prevMB != nil {
		y += *prevMB
	}
	placed.Height = y
	return placed, nil
}
