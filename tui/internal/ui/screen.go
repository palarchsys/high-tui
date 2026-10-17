package ui

import (
	tea "charm.land/bubbletea/v2"
	"high-tui/internal/layout"
	"high-tui/internal/model"
	"high-tui/internal/run"
)

// Screen renders one frame from the generated windows. It does not parse HTML.
func Screen(root string, width int) (string, error) {
	app, err := newApp(root, width)
	if err != nil {
		return "", err
	}
	return app.frame(), nil
}

type app struct {
	root      string
	wins      []model.Window
	origin    []model.Window
	placed    []layout.Placed
	cur       int
	stack     []int
	focus     int
	sel       map[string]int
	scroll    map[string]int
	val       map[string]string
	file      map[string]string
	name      map[string]string
	dirty     map[string]bool
	rawLabel  map[string]string
	rawValue  map[string]string
	resolved  map[string]string
	status    string
	edit      *editor
	width     int
	termW     int
	keyMode   int
	focusAt   map[string]int
	spin      int
	running   bool
	runGen    int
	cancel    func()
	cmdLabel  string
	focusMark string
	pending   tea.Cmd
}

type editor struct {
	zone      string
	multiline bool
	original  string
	state     run.Edit
}

func newApp(root string, width int) (*app, error) {
	a := &app{
		root: root, wins: generatedWindows(), origin: generatedWindows(), width: width,
		sel: map[string]int{}, scroll: map[string]int{}, val: map[string]string{},
		file: map[string]string{}, name: map[string]string{}, dirty: map[string]bool{},
		rawLabel: map[string]string{}, rawValue: map[string]string{}, resolved: map[string]string{},
		focusAt: map[string]int{}, status: "Prêt",
	}
	for wi := range a.wins {
		win := &a.wins[wi]
		a.rawLabel[win.ID] = win.Label
		for ri := range win.Rows {
			for zi := range win.Rows[ri].Zones {
				zone := &win.Rows[ri].Zones[zi]
				k := key(win.ID, zone.ID)
				a.rawLabel[k] = zone.Label
				a.rawValue[k] = zone.Value
				a.file[k] = zone.File
				a.name[k] = zone.Name
				if err := a.readZone(win.ID, *zone); err != nil {
					return nil, err
				}
				for ii := range zone.Items {
					item := &zone.Items[ii]
					ik := key(win.ID, item.ID)
					a.rawLabel[ik] = item.Label
					a.rawValue[ik] = item.Value
				}
			}
		}
	}
	if err := a.resolveAll(); err != nil {
		return nil, err
	}
	if err := a.layoutAll(); err != nil {
		return nil, err
	}
	if err := a.applyFocus(); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *app) resolveAll() error {
	resolve := func(raw, where string) (string, error) {
		if raw == "" {
			return "", nil
		}
		got, err := run.Resolve(raw, a.root, where)
		if err != nil {
			return "", err
		}
		a.resolved[raw] = got
		return got, nil
	}
	var err error
	for wi := range a.wins {
		win := &a.wins[wi]
		src := a.origin[wi]
		win.Label, err = resolve(src.Label, win.ID)
		if err != nil {
			return err
		}
		for ri := range win.Rows {
			for zi := range win.Rows[ri].Zones {
				zone := &win.Rows[ri].Zones[zi]
				from := src.Rows[ri].Zones[zi]
				zone.Label, err = resolve(from.Label, zone.ID)
				if err != nil {
					return err
				}
				zone.Value, err = resolve(from.Value, zone.ID)
				if err != nil {
					return err
				}
				zone.Focus = from.Focus
				zone.Enter = from.Enter
				zone.Focus.Value, err = resolve(from.Focus.Value, zone.ID)
				if err != nil {
					return err
				}
				zone.Enter.Value, err = resolve(from.Enter.Value, zone.ID)
				if err != nil {
					return err
				}
				for ii := range zone.Items {
					item := &zone.Items[ii]
					rawItem := from.Items[ii]
					item.Label, err = resolve(rawItem.Label, item.ID)
					if err != nil {
						return err
					}
					item.Value, err = resolve(rawItem.Value, item.ID)
					if err != nil {
						return err
					}
					item.Focus = rawItem.Focus
					item.Enter = rawItem.Enter
					item.Focus.Value, err = resolve(rawItem.Focus.Value, item.ID)
					if err != nil {
						return err
					}
					item.Enter.Value, err = resolve(rawItem.Enter.Value, item.ID)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (a *app) layoutAll() error {
	a.placed = nil
	for _, win := range a.wins {
		width := a.width
		if width == 0 && a.termW > 0 {
			max := win.Style.MaxWidth
			if max < 1 {
				max = 260
			}
			if a.termW < max {
				max = a.termW
			}
			n, err := layout.Fitted(max, win.ID)
			if err != nil {
				return err
			}
			width = n
		}
		placed, err := layout.Place(win, width)
		if err != nil {
			return err
		}
		a.placed = append(a.placed, placed)
	}
	return nil
}

func (a *app) relayout(termW int) error {
	prevW := a.termW
	prev := append([]layout.Placed(nil), a.placed...)
	a.termW = termW
	if err := a.layoutAll(); err != nil {
		a.termW = prevW
		a.placed = prev
		return err
	}
	return nil
}

func (a *app) current() layout.Placed { return a.placed[a.cur] }

func (a *app) boxes() []layout.Box {
	var out []layout.Box
	for _, band := range a.current().Bands {
		for _, box := range band.Boxes {
			switch box.Zone.Tag {
			case "nav", "select", "input", "textarea", "button", "command":
				out = append(out, box)
			}
		}
	}
	return out
}

func (a *app) focused() (layout.Box, bool) {
	boxes := a.boxes()
	if len(boxes) == 0 {
		return layout.Box{}, false
	}
	if a.focus < 0 {
		a.focus = 0
	}
	if a.focus >= len(boxes) {
		a.focus = len(boxes) - 1
	}
	return boxes[a.focus], true
}

func (a *app) frame() string {
	box, ok := a.focused()
	focusID := ""
	if ok {
		focusID = box.Zone.ID
	}
	val := a.val
	editID := ""
	if a.edit != nil {
		editID = a.edit.zone
		val = make(map[string]string, len(a.val)+1)
		for k, v := range a.val {
			val[k] = v
		}
		val[key(a.winID(), a.edit.zone)] = a.edit.state.Text
	}
	return paint(viewState{
		placed: a.current(), winID: a.wins[a.cur].ID, focusID: focusID, editID: editID,
		status: a.status, sel: a.sel, scroll: a.scroll, val: val, name: a.name,
		resolved: a.resolved,
	})
}

func (a *app) winID() string { return a.wins[a.cur].ID }
