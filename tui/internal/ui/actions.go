package ui

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"high-tui/internal/layout"
	"high-tui/internal/model"
	"high-tui/internal/run"
	"high-tui/internal/text"
)

const (
	keysUnknown = iota
	keysPlain
	keysEnhanced
)

const noShiftMsg = "Ce terminal ne distingue pas Maj+Entrée"

var spinFrames = []string{"⠋", "⠙", "⠹", "⠸"}

func (a *app) applyFocus() error {
	box, ok := a.focused()
	if !ok {
		return nil
	}
	mark := a.winID() + "\x1f" + box.Zone.ID
	if box.Zone.Tag == "nav" || box.Zone.Tag == "select" {
		mark += fmt.Sprintf("\x1f%d", a.sel[key(a.winID(), box.Zone.ID)])
	}
	if mark == a.focusMark {
		return nil
	}
	a.focusMark = mark
	act, _ := a.action("focus")
	cmd, err := a.perform(act)
	if err != nil {
		return err
	}
	a.pending = cmd
	return nil
}

func (a *app) action(which string) (model.Action, string) {
	box, ok := a.focused()
	if !ok {
		return model.Action{}, ""
	}
	zone := box.Zone
	if zone.Tag == "nav" || zone.Tag == "select" {
		sel := a.sel[key(a.winID(), zone.ID)]
		if sel < 0 || sel >= len(zone.Items) {
			return model.Action{}, zone.Label
		}
		item := zone.Items[sel]
		if which == "focus" {
			return item.Focus, item.Label
		}
		return item.Enter, item.Label
	}
	if which == "focus" {
		return zone.Focus, zone.Label
	}
	return zone.Enter, zone.Label
}

func (a *app) perform(act model.Action) (tea.Cmd, error) {
	switch act.Kind {
	case "", "return":
		if act.Kind == "return" {
			return nil, a.goBack()
		}
		return nil, nil
	case "load":
		return nil, a.load(act)
	case "edit":
		return nil, a.beginEdit(act.Target)
	case "set":
		return nil, a.set(act)
	case "window":
		err := a.open(act.Href)
		return a.takePending(), err
	case "save":
		return nil, a.save()
	case "refresh":
		err := a.refresh()
		return a.takePending(), err
	case "run":
		_, label := a.action("enter")
		return a.startRun(act, label), nil
	default:
		return nil, model.Errf("action inconnue « %s »", act.Kind)
	}
}

func (a *app) takePending() tea.Cmd {
	cmd := a.pending
	a.pending = nil
	return cmd
}

func (a *app) zoneByID(id string) (model.Zone, layout.Box, bool) {
	for _, band := range a.current().Bands {
		for _, box := range band.Boxes {
			if box.Zone.ID == id {
				return box.Zone, box, true
			}
		}
	}
	return model.Zone{}, layout.Box{}, false
}

func (a *app) readZone(winID string, zone model.Zone) error {
	k := key(winID, zone.ID)
	file := a.file[k]
	if file == "" {
		file = zone.File
	}
	if file == "" || (zone.Tag != "input" && zone.Tag != "textarea") {
		return nil
	}
	path := run.Path(a.root, file)
	if zone.Tag == "input" {
		name := a.name[k]
		if name == "" {
			name = zone.Name
		}
		if name == "" {
			return nil
		}
		val, err := run.ReadVar(path, name)
		if err != nil {
			return err
		}
		a.val[k] = val
		return nil
	}
	body, err := readBody(path, file)
	if err != nil {
		return err
	}
	a.val[k] = body
	return nil
}

func readBody(path, shown string) (string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return "", model.Errf("fichier introuvable : %s", shown)
	}
	return strings.TrimRight(string(raw), "\n"), nil
}

func (a *app) load(act model.Action) error {
	zone, _, ok := a.zoneByID(act.Target)
	if !ok {
		return model.Errf("cible absente : %s", act.Target)
	}
	if zone.Tag != "input" && zone.Tag != "textarea" {
		return model.Errf("cible incompatible : %s", act.Target)
	}
	if act.File == "" {
		return model.Errf("fichier manquant (%s)", act.Target)
	}
	k := key(a.winID(), zone.ID)
	a.file[k] = act.File
	path := run.Path(a.root, act.File)
	if act.Name != "" {
		a.name[k] = act.Name
		val, err := run.ReadVar(path, act.Name)
		if err != nil {
			return err
		}
		a.val[k] = val
	} else {
		a.name[k] = ""
		body, err := readBody(path, act.File)
		if err != nil {
			return err
		}
		a.val[k] = body
	}
	a.dirty[k] = false
	a.scroll[k] = 0
	return nil
}

func (a *app) set(act model.Action) error {
	if act.File == "" {
		return model.Errf("fichier manquant")
	}
	if act.Name == "" {
		return model.Errf("variable manquante")
	}
	if err := run.WriteVar(run.Path(a.root, act.File), act.Name, act.Value); err != nil {
		return err
	}
	if act.Target != "" {
		k := key(a.winID(), act.Target)
		a.val[k] = act.Value
		a.file[k] = act.File
		a.name[k] = act.Name
		a.dirty[k] = false
		a.scroll[k] = 0
	}
	a.status = "Enregistré"
	return nil
}

func (a *app) save() error {
	win := a.wins[a.cur]
	for _, row := range win.Rows {
		for _, zone := range row.Zones {
			if zone.Readonly || (zone.Tag != "input" && zone.Tag != "textarea") {
				continue
			}
			k := key(win.ID, zone.ID)
			if !a.dirty[k] {
				continue
			}
			if err := a.writeZone(zone, k); err != nil {
				return err
			}
			a.dirty[k] = false
		}
	}
	a.status = "Enregistré"
	return nil
}

func (a *app) writeZone(zone model.Zone, k string) error {
	file := a.file[k]
	if file == "" {
		file = zone.File
	}
	if file == "" {
		return model.Errf("fichier manquant (%s)", zone.ID)
	}
	path := run.Path(a.root, file)
	if zone.Tag == "input" {
		name := a.name[k]
		if name == "" {
			name = zone.Name
		}
		if name == "" {
			return model.Errf("variable manquante (%s)", zone.ID)
		}
		return run.WriteVar(path, name, a.val[k])
	}
	body := a.val[k]
	if !strings.HasSuffix(body, "\n") {
		body += "\n"
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		return model.Errf("fichier introuvable : %s", file)
	}
	return nil
}

func (a *app) refresh() error {
	a.edit = nil
	if err := a.resolveAll(); err != nil {
		return err
	}
	for wi := range a.wins {
		win := a.wins[wi]
		for _, row := range win.Rows {
			for _, zone := range row.Zones {
				k := key(win.ID, zone.ID)
				a.file[k] = zone.File
				a.name[k] = zone.Name
				a.dirty[k] = false
				a.scroll[k] = 0
				if err := a.readZone(win.ID, zone); err != nil {
					return err
				}
			}
		}
	}
	a.focusMark = ""
	if err := a.applyFocus(); err != nil {
		return err
	}
	a.status = "Rechargé"
	return nil
}

func (a *app) beginEdit(id string) error {
	zone, _, ok := a.zoneByID(id)
	if !ok {
		return model.Errf("cible absente : %s", id)
	}
	if zone.Tag != "input" && zone.Tag != "textarea" {
		return model.Errf("cible incompatible : %s", id)
	}
	if zone.Readonly {
		return model.Errf("cible non éditable : %s", id)
	}
	k := key(a.winID(), zone.ID)
	text := a.val[k]
	a.edit = &editor{
		zone: id, multiline: zone.Tag == "textarea",
		original: text, state: run.Edit{Text: text, Cursor: len([]rune(text))},
	}
	a.reveal()
	return nil
}

func (a *app) commitEdit() error {
	ed := a.edit
	if ed == nil {
		return nil
	}
	zone, _, ok := a.zoneByID(ed.zone)
	if !ok {
		a.edit = nil
		return model.Errf("cible absente : %s", ed.zone)
	}
	k := key(a.winID(), ed.zone)
	prev := a.val[k]
	a.val[k] = ed.state.Text
	file := a.file[k]
	if file == "" {
		file = zone.File
	}
	if file != "" {
		if err := a.writeZone(zone, k); err != nil {
			a.val[k] = prev
			a.status = err.Error()
			return nil
		}
		a.dirty[k] = false
		a.status = "Enregistré"
		a.edit = nil
		return nil
	}
	if ed.state.Text != ed.original {
		a.dirty[k] = true
		a.status = "Modifié"
	}
	a.edit = nil
	return nil
}

func (a *app) cancelEdit() {
	if a.edit == nil {
		return
	}
	a.edit = nil
	a.status = "Annulé"
}

func (a *app) open(href string) error {
	idx := a.windowIndex(href)
	if idx < 0 {
		return model.Errf("fenêtre introuvable : %s", href)
	}
	a.focusAt[a.winID()] = a.focus
	a.stack = append(a.stack, a.cur)
	a.cur = idx
	a.focus = a.focusAt[a.winID()]
	a.edit = nil
	a.focusMark = ""
	a.status = "Prêt"
	return a.applyFocus()
}

func (a *app) goBack() error {
	if len(a.stack) == 0 {
		a.status = "Déjà à la première fenêtre"
		return nil
	}
	a.focusAt[a.winID()] = a.focus
	a.cur = a.stack[len(a.stack)-1]
	a.stack = a.stack[:len(a.stack)-1]
	a.focus = a.focusAt[a.winID()]
	a.edit = nil
	a.focusMark = ""
	a.status = "Prêt"
	return a.applyFocus()
}

func (a *app) windowIndex(href string) int {
	base := filepath.Base(href)
	id := strings.TrimSuffix(base, filepath.Ext(base))
	for i, win := range a.wins {
		if win.File == href || win.File == base || win.ID == id {
			return i
		}
	}
	return -1
}

func (a *app) moveFocus(delta int) error {
	boxes := a.boxes()
	if len(boxes) == 0 {
		return nil
	}
	a.focus += delta
	if a.focus < 0 {
		a.focus = len(boxes) - 1
	}
	if a.focus >= len(boxes) {
		a.focus = 0
	}
	return a.applyFocus()
}

func (a *app) moveVertical(delta int) error {
	box, ok := a.focused()
	if !ok {
		return nil
	}
	zone := box.Zone
	k := key(a.winID(), zone.ID)
	if zone.Tag == "nav" || zone.Tag == "select" {
		if len(zone.Items) == 0 {
			return nil
		}
		sel := a.sel[k] + delta
		if sel < 0 || sel >= len(zone.Items) {
			return nil
		}
		a.sel[k] = sel
		return a.applyFocus()
	}
	body := logical(viewState{
		placed: a.current(), winID: a.winID(), sel: a.sel, scroll: a.scroll,
		val: a.val, name: a.name, resolved: a.resolved, status: a.status,
	}, box)
	_, maxScroll := text.View(body, box.TextW, box.Lines, a.scroll[k])
	next := a.scroll[k] + delta
	if next < 0 || next > maxScroll {
		return a.moveFocus(delta)
	}
	a.scroll[k] = next
	return nil
}

func (a *app) editKey(keyName string) tea.Cmd {
	if a.edit == nil {
		return nil
	}
	next, result := run.ApplyEdit(a.edit.state, keyName, a.edit.multiline)
	a.edit.state = next
	a.reveal()
	switch result {
	case run.EditCommit:
		_ = a.commitEdit()
	case run.EditCancel:
		a.cancelEdit()
	}
	return nil
}

func (a *app) paste(content string) {
	if a.edit == nil {
		return
	}
	content = strings.ReplaceAll(content, "\r", "")
	if !a.edit.multiline {
		content = strings.ReplaceAll(content, "\n", "")
	}
	ed := a.edit.state
	rs := []rune(ed.Text)
	if ed.Cursor < 0 {
		ed.Cursor = 0
	}
	if ed.Cursor > len(rs) {
		ed.Cursor = len(rs)
	}
	ins := []rune(content)
	out := make([]rune, 0, len(rs)+len(ins))
	out = append(out, rs[:ed.Cursor]...)
	out = append(out, ins...)
	out = append(out, rs[ed.Cursor:]...)
	ed.Text = string(out)
	ed.Cursor += len(ins)
	a.edit.state = ed
	a.reveal()
}

func (a *app) reveal() {
	if a.edit == nil {
		return
	}
	zone, box, ok := a.zoneByID(a.edit.zone)
	if !ok {
		return
	}
	body, cursor := a.editBody(zone)
	width := box.TextW
	rows := text.Wrap(body, width)
	if len(rows) > box.Lines && box.TextW >= 2 {
		width = box.TextW - 1
	}
	line, _ := text.Locate(body, width, cursor)
	k := key(a.winID(), zone.ID)
	scroll := a.scroll[k]
	if line < scroll {
		scroll = line
	}
	if line >= scroll+box.Lines {
		scroll = line - box.Lines + 1
	}
	if scroll < 0 {
		scroll = 0
	}
	a.scroll[k] = scroll
}

func (a *app) editBody(zone model.Zone) (string, int) {
	text := a.edit.state.Text
	cursor := a.edit.state.Cursor
	if zone.Tag == "input" {
		name := a.name[key(a.winID(), zone.ID)]
		if name == "" {
			name = zone.Name
		}
		text = name + "\n" + text
		cursor += len([]rune(name)) + 1
	}
	return text, cursor
}

func (a *app) cursor() (int, int, bool) {
	if a.edit == nil {
		return 0, 0, false
	}
	zone, box, ok := a.zoneByID(a.edit.zone)
	if !ok {
		return 0, 0, false
	}
	body, pos := a.editBody(zone)
	width := box.TextW
	rows := text.Wrap(body, width)
	if len(rows) > box.Lines && box.TextW >= 2 {
		width = box.TextW - 1
	}
	line, col := text.Locate(body, width, pos)
	if col > width {
		col = width
	}
	scroll := a.scroll[key(a.winID(), zone.ID)]
	vis := line - scroll
	if vis < 0 {
		vis = 0
	}
	if vis >= box.Lines {
		vis = box.Lines - 1
	}
	style := zone.Base
	lh := style.LineHeight
	if lh < 1 {
		lh = 1
	}
	x := box.X + style.BorderWidth + style.Pad[3] + col
	y := box.Y + style.BorderWidth + style.Pad[0] + vis*lh
	return x, y, true
}

type cmdDone struct {
	gen       int
	zone      string
	out       string
	code      int
	cancelled bool
}

type spinMsg struct{ gen int }

func (a *app) startRun(act model.Action, label string) tea.Cmd {
	if label == "" {
		label = act.Run
	}
	a.runGen++
	gen := a.runGen
	ctx, cancel := context.WithCancel(context.Background())
	a.cancel = cancel
	a.running = true
	a.spin = 0
	a.cmdLabel = label
	a.status = spinFrames[0] + " En cours : " + label
	return tea.Batch(spinCmd(gen), execCmd(ctx, gen, act.Run, act.Output))
}

func (a *app) stopRun() {
	a.runGen++
	a.running = false
	if a.cancel != nil {
		a.cancel()
		a.cancel = nil
	}
}

func spinCmd(gen int) tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(time.Time) tea.Msg {
		return spinMsg{gen}
	})
}

func execCmd(ctx context.Context, gen int, script, zone string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.CommandContext(ctx, "sh", "-c", script)
		out, err := cmd.CombinedOutput()
		code := 0
		cancelled := ctx.Err() != nil
		if err != nil && !cancelled {
			if ee, ok := err.(*exec.ExitError); ok {
				code = ee.ExitCode()
			} else {
				code = 1
			}
		}
		return cmdDone{gen: gen, zone: zone, out: string(out), code: code, cancelled: cancelled}
	}
}
