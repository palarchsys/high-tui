package layout

import (
	"testing"

	"high-tui/internal/model"
)

func TestPlaceHeaderAndNav(t *testing.T) {
	box := model.Style{
		Color: "#e6edf3", Background: "#161b22", BorderColor: "#30363d",
		BorderWidth: 1, LineHeight: 1, Gap: 1,
		Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0},
	}
	zone := func(id string) *model.Zone {
		return &model.Zone{Tag: "nav", ID: id, Base: box, FocusOn: box}
	}
	win := model.Window{
		ID: "main",
		Style: model.Style{
			Color: "#e6edf3", Background: "#0d1117", BorderColor: "#e6edf3",
			LineHeight: 1, MaxWidth: 260,
		},
		Grid: model.Grid{Rows: []model.GridRow{
			{ID: "entete", Lines: 1, Cells: []model.Cell{{Fr: 1, Zone: zone("titre")}}},
			{ID: "listes", Lines: 4, Cells: []model.Cell{
				{Fr: 1, Zone: zone("a")}, {Fr: 1, Zone: zone("b")},
				{Fr: 1, Zone: zone("c")}, {Fr: 1, Zone: zone("d")},
			}},
		}},
	}
	placed, err := Place(win, 260)
	if err != nil {
		t.Fatal(err)
	}
	if placed.Bands[0].Height != 3 || placed.Bands[0].Cells[0].Box.OuterH != 3 {
		t.Fatal(placed.Bands[0].Height, placed.Bands[0].Cells[0].Box.OuterH)
	}
	if placed.Bands[1].Height != 6 || placed.Bands[1].Y != 5 {
		t.Fatal(placed.Bands[1].Height, placed.Bands[1].Y)
	}
}

func TestPlaceSideMargins(t *testing.T) {
	box := model.Style{
		BorderWidth: 1, LineHeight: 1, Gap: 1,
		Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 1, 1, 1},
	}
	zone := func(id string) *model.Zone {
		return &model.Zone{Tag: "nav", ID: id, Base: box, FocusOn: box}
	}
	win := model.Window{
		ID:    "main",
		Style: model.Style{LineHeight: 1, MaxWidth: 80},
		Grid: model.Grid{Rows: []model.GridRow{{
			ID: "listes", Lines: 4,
			Cells: []model.Cell{
				{Fr: 1, Zone: zone("a")}, {Fr: 1, Zone: zone("b")},
				{Fr: 1, Zone: zone("c")}, {Fr: 1, Zone: zone("d")},
			},
		}}},
	}
	placed, err := Place(win, 80)
	if err != nil {
		t.Fatal(err)
	}
	band := placed.Bands[0]
	if band.InsetL != 1 || band.InsetR != 1 {
		t.Fatal(band.InsetL, band.InsetR)
	}
	if band.Y != 1 || band.Height != 6 || band.Cells[0].Box.X != 1 {
		t.Fatal(band.Y, band.Height, band.Cells[0].Box.X)
	}
	if band.Cells[0].Gap != 1 || band.Cells[1].Gap != 1 || band.Cells[2].Gap != 1 || band.Cells[3].Gap != 0 {
		t.Fatal(band.Cells[0].Gap, band.Cells[1].Gap, band.Cells[2].Gap, band.Cells[3].Gap)
	}
	used := band.InsetL + band.InsetR
	for _, item := range band.Cells {
		used += item.Width + item.Gap
	}
	if used != 80 {
		t.Fatal(used)
	}
}

func TestNestedGridFillsCell(t *testing.T) {
	box := model.Style{
		BorderWidth: 1, LineHeight: 1, Gap: 1,
		Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 1, 1, 1},
	}
	zone := func(id string) *model.Zone {
		return &model.Zone{Tag: "nav", ID: id, Base: box, FocusOn: box}
	}
	inner := model.Grid{ID: "corps", Rows: []model.GridRow{
		{ID: "haut", Lines: 1, Cells: []model.Cell{{Fr: 1, Zone: zone("haut")}}},
		{ID: "bas", Lines: 1, Cells: []model.Cell{{Fr: 1, Zone: zone("bas")}}},
	}}
	win := model.Window{
		ID:    "main",
		Style: model.Style{LineHeight: 1, MaxWidth: 40},
		Grid: model.Grid{ID: "main", Rows: []model.GridRow{{
			ID: "bande", Lines: 1, Cells: []model.Cell{
				{Fr: 1, Zone: zone("cote")},
				{Fr: 1, Grid: &inner},
			},
		}}},
	}
	placed, err := Place(win, 40)
	if err != nil {
		t.Fatal(err)
	}
	band := placed.Bands[0]
	if len(band.Cells) != 2 || band.Cells[1].Sub == nil || band.Cells[1].Box != nil {
		t.Fatal("grille imbriquée absente")
	}
	if band.Cells[0].Box.OuterH != band.Cells[1].Sub.Height {
		t.Fatal(band.Cells[0].Box.OuterH, band.Cells[1].Sub.Height)
	}
	if band.Cells[1].Sub.Width != band.Cells[1].Width {
		t.Fatal(band.Cells[1].Sub.Width, band.Cells[1].Width)
	}
	if band.Cells[1].Sub.Bands[0].Y != 0 {
		t.Fatal(band.Cells[1].Sub.Bands[0].Y)
	}
	used := band.InsetL + band.InsetR
	for _, item := range band.Cells {
		used += item.Width + item.Gap
	}
	if used != 40 {
		t.Fatal(used)
	}
}

func TestWidths(t *testing.T) {
	one, err := Widths(10, []int{1}, 0, "essai")
	if err != nil || len(one) != 1 || one[0] != 10 {
		t.Fatal(one, err)
	}
	two, err := Widths(10, []int{1, 1, 2}, 0, "essai")
	if err != nil || two[0] != 2 || two[1] != 2 || two[2] != 6 {
		t.Fatal(two, err)
	}
	three, err := Widths(10, []int{1, 1, 2}, 1, "essai")
	if err != nil || three[0] != 2 || three[1] != 2 || three[2] != 4 {
		t.Fatal(three, err)
	}
}
