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
	zone := func(id string) model.Zone {
		return model.Zone{Tag: "nav", ID: id, Fr: 1, Base: box, FocusOn: box}
	}
	win := model.Window{
		ID: "main",
		Style: model.Style{
			Color: "#e6edf3", Background: "#0d1117", BorderColor: "#e6edf3",
			LineHeight: 1, MaxWidth: 260,
		},
		Rows: []model.Row{
			{ID: "entete", Tag: "header", Lines: 1, Zones: []model.Zone{zone("titre")}},
			{ID: "listes", Tag: "main", Lines: 4, Zones: []model.Zone{
				zone("a"), zone("b"), zone("c"), zone("d"),
			}},
		},
	}
	placed, err := Place(win, 260)
	if err != nil {
		t.Fatal(err)
	}
	if placed.Bands[0].Height != 3 || placed.Bands[0].Boxes[0].OuterH != 3 {
		t.Fatal(placed.Bands[0].Height, placed.Bands[0].Boxes[0].OuterH)
	}
	if placed.Bands[1].Height != 6 || placed.Bands[1].Y != 5 {
		t.Fatal(placed.Bands[1].Height, placed.Bands[1].Y)
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
