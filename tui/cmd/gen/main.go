// Command gen reads windows/*.html and the matching CSS once and writes
// internal/ui/generated.go. The executable does not parse HTML again.
package main

import (
	"fmt"
	"go/format"
	"os"
	"strconv"
	"strings"

	"high-tui/internal/build"
	"high-tui/internal/model"
)

func main() {
	root, err := os.Getwd()
	if err != nil {
		fail(err)
	}
	wins, err := build.Load(root)
	if err != nil {
		fail(err)
	}
	var b strings.Builder
	b.WriteString("// Code généré par go generate. Ne pas modifier.\n\n")
	b.WriteString("package ui\n\n")
	b.WriteString("import \"high-tui/internal/model\"\n\n")
	b.WriteString("func generatedWindows() []model.Window {\n\treturn []model.Window{\n")
	for _, win := range wins {
		writeWindow(&b, win)
	}
	b.WriteString("\t}\n}\n")
	src, err := format.Source([]byte(b.String()))
	if err != nil {
		fail(err)
	}
	if err := os.WriteFile("internal/ui/generated.go", src, 0o644); err != nil {
		fail(err)
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}

func q(s string) string { return strconv.Quote(s) }

func writeWindow(b *strings.Builder, win model.Window) {
	b.WriteString("\t\t{\n")
	fmt.Fprintf(b, "\t\t\tID: %s, Label: %s, File: %s,\n", q(win.ID), q(win.Label), q(win.File))
	writeStrings(b, "Classes", win.Classes, 3)
	b.WriteString("\t\t\tStyle: ")
	writeStyle(b, win.Style)
	b.WriteString(",\n\t\t\tGrid: ")
	writeGrid(b, win.Grid, 3)
	b.WriteString(",\n\t\t},\n")
}

func writeGrid(b *strings.Builder, grid model.Grid, indent int) {
	pad := strings.Repeat("\t", indent)
	b.WriteString("model.Grid{\n")
	fmt.Fprintf(b, pad+"\tID: %s,\n", q(grid.ID))
	writeStrings(b, "Classes", grid.Classes, indent+1)
	b.WriteString(pad + "\tRows: []model.GridRow{\n")
	for _, row := range grid.Rows {
		b.WriteString(pad + "\t\t{\n")
		fmt.Fprintf(b, pad+"\t\t\tID: %s, Lines: %d,\n", q(row.ID), row.Lines)
		writeStrings(b, "Classes", row.Classes, indent+3)
		b.WriteString(pad + "\t\t\tCells: []model.Cell{\n")
		for _, cell := range row.Cells {
			b.WriteString(pad + "\t\t\t\t{\n")
			fmt.Fprintf(b, pad+"\t\t\t\t\tFr: %d,\n", cell.Fr)
			if cell.Zone != nil {
				b.WriteString(pad + "\t\t\t\t\tZone: &")
				writeZone(b, *cell.Zone)
			}
			if cell.Grid != nil {
				b.WriteString(pad + "\t\t\t\t\tGrid: &")
				writeGrid(b, *cell.Grid, indent+5)
			}
			b.WriteString(",\n" + pad + "\t\t\t\t},\n")
		}
		b.WriteString(pad + "\t\t\t},\n" + pad + "\t\t},\n")
	}
	b.WriteString(pad + "\t},\n" + pad + "}")
}

func writeZone(b *strings.Builder, zone model.Zone) {
	b.WriteString("model.Zone{\n")
	fmt.Fprintf(b, "\t\t\t\t\t\t\tTag: %s, ID: %s, Label: %s, Value: %s,\n", q(zone.Tag), q(zone.ID), q(zone.Label), q(zone.Value))
	fmt.Fprintf(b, "\t\t\t\t\t\t\tFile: %s, Name: %s, Save: %s, Href: %s, Run: %s, Output: %s,\n", q(zone.File), q(zone.Name), q(zone.Save), q(zone.Href), q(zone.Run), q(zone.Output))
	fmt.Fprintf(b, "\t\t\t\t\t\t\tReadonly: %t, ItemsFrom: %s, NextNav: %s,\n", zone.Readonly, q(zone.ItemsFrom), q(zone.NextNav))
	writeStrings(b, "Classes", zone.Classes, 7)
	b.WriteString("\t\t\t\t\t\t\tFocus: ")
	writeAction(b, zone.Focus)
	b.WriteString(",\n\t\t\t\t\t\t\tEnter: ")
	writeAction(b, zone.Enter)
	b.WriteString(",\n\t\t\t\t\t\t\tBase: ")
	writeStyle(b, zone.Base)
	b.WriteString(",\n\t\t\t\t\t\t\tFocusOn: ")
	writeStyle(b, zone.FocusOn)
	b.WriteString(",\n")
	if len(zone.Items) > 0 {
		b.WriteString("\t\t\t\t\t\t\tItems: []model.Item{\n")
		for _, item := range zone.Items {
			b.WriteString("\t\t\t\t\t\t\t\t{")
			fmt.Fprintf(b, "ID: %s, Label: %s, Value: %s, Focus: ", q(item.ID), q(item.Label), q(item.Value))
			writeAction(b, item.Focus)
			b.WriteString(", Enter: ")
			writeAction(b, item.Enter)
			b.WriteString("},\n")
		}
		b.WriteString("\t\t\t\t\t\t\t},\n")
	}
	b.WriteString("\t\t\t\t\t\t}")
}

func writeAction(b *strings.Builder, act model.Action) {
	fmt.Fprintf(b, "model.Action{Kind: %s, Target: %s, File: %s, Name: %s, Href: %s, Run: %s, Output: %s, Value: %s}",
		q(act.Kind), q(act.Target), q(act.File), q(act.Name), q(act.Href), q(act.Run), q(act.Output), q(act.Value))
}

func writeStyle(b *strings.Builder, st model.Style) {
	fmt.Fprintf(b, "model.Style{Color: %s, TitleColor: %s, Background: %s, BorderColor: %s, BorderWidth: %d, Radius: %d, LineHeight: %d, Gap: %d, MaxWidth: %d, Bold: %t, Pad: [4]int{%d, %d, %d, %d}, Mar: [4]int{%d, %d, %d, %d}}",
		q(st.Color), q(st.TitleColor), q(st.Background), q(st.BorderColor), st.BorderWidth, st.Radius, st.LineHeight, st.Gap, st.MaxWidth, st.Bold,
		st.Pad[0], st.Pad[1], st.Pad[2], st.Pad[3], st.Mar[0], st.Mar[1], st.Mar[2], st.Mar[3])
}

func writeStrings(b *strings.Builder, name string, values []string, tabs int) {
	pad := strings.Repeat("\t", tabs)
	if len(values) == 0 {
		return
	}
	b.WriteString(pad + name + ": []string{")
	for i, v := range values {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(q(v))
	}
	b.WriteString("},\n")
}
