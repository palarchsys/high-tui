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

// Item is one row of a nav or a select.
type Item struct {
	ID    string
	Label string
	Value string
	Focus Action
	Enter Action
}

// Zone is one cell of a row.
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
	Fr       int
	Readonly bool
	Focus    Action
	Enter    Action
	Items    []Item
	Base     Style
	FocusOn  Style
}

// Row is a horizontal band. Lines is the text height inside every cell.
type Row struct {
	ID      string
	Tag     string
	Classes []string
	Lines   int
	Style   Style
	Zones   []Zone
}

// Window is one screen.
type Window struct {
	ID      string
	Label   string
	File    string
	Classes []string
	Style   Style
	Rows    []Row
}
