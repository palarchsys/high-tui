// Package model holds the window tree produced by the generator.
package model

import "fmt"

// Error is a French diagnostic. The program prints it and exits non-zero.
type Error string

func (e Error) Error() string { return "erreur : " + string(e) }

// Errf formats a French diagnostic without the « erreur : » prefix.
func Errf(format string, args ...any) error {
	return Error(fmt.Sprintf(format, args...))
}

// Action is a focus or enter action. An empty Kind does nothing.
type Action struct {
	Kind   string
	Target string
	File   string
	Name   string
	Href   string
	Run    string
	Output string
	Value  string
}

// Style is a resolved box style. Zero values mean the feature is off.
type Style struct {
	Color, TitleColor, Background, BorderColor string
	BorderWidth                                int
	Radius                                     int
	LineHeight                                 int
	Gap                                        int
	MaxWidth                                   int
	Bold                                       bool
	Pad                                        [4]int // top, right, bottom, left
	Mar                                        [4]int
}

// Item is one row of a nav.
type Item struct {
	ID    string
	Label string
	Value string
	Focus Action
	Enter Action
}

// Zone is the content painted in one cell. It fills that cell.
type Zone struct {
	Tag      string
	ID       string
	Classes  []string
	Label    string
	Value    string
	File     string
	Name     string
	Save     string
	Href     string
	Run      string
	Output   string
	Readonly bool
	Focus    Action
	Enter    Action
	Items    []Item
	Base     Style
	FocusOn  Style
}

// Cell is one slot of a grid row. It holds one zone or one nested grid.
type Cell struct {
	Fr   int
	Grid *Grid
	Zone *Zone
}

// GridRow is one horizontal band. Lines is the text height of its zones.
type GridRow struct {
	ID      string
	Classes []string
	Lines   int
	Cells   []Cell
}

// Grid is a table of cells. It is not a zone.
type Grid struct {
	ID      string
	Classes []string
	Rows    []GridRow
}

// Each visits every zone in document order.
func (g *Grid) Each(fn func(*Zone)) {
	if g == nil {
		return
	}
	for i := range g.Rows {
		for j := range g.Rows[i].Cells {
			cell := &g.Rows[i].Cells[j]
			if cell.Zone != nil {
				fn(cell.Zone)
			}
			cell.Grid.Each(fn)
		}
	}
}

// Zones returns every zone in document order.
func (g *Grid) Zones() []*Zone {
	var out []*Zone
	g.Each(func(zone *Zone) { out = append(out, zone) })
	return out
}

// Window is one screen.
type Window struct {
	ID      string
	Label   string
	File    string
	Classes []string
	Style   Style
	Grid    Grid
}
