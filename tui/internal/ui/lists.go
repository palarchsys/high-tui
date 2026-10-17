package ui

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"high-tui/internal/items"
	"high-tui/internal/layout"
	"high-tui/internal/model"
	"high-tui/internal/run"
)

type listSpin struct{ gen int }

type listDone struct {
	gen       int
	nav       string
	items     []model.Item
	err       error
	cancelled bool
}

func fillAction(act model.Action, item model.Item) model.Action {
	if act.File == "" {
		act.File = item.File
	}
	if act.Name == "" {
		act.Name = item.Name
	}
	if act.Value == "" {
		act.Value = item.Value
	}
	return act
}

func (a *app) loadDeclaredItems() error {
	var failed error
	for wi := range a.wins {
		win := &a.wins[wi]
		win.Grid.Each(func(zone *model.Zone) {
			if failed != nil || zone.Tag != "nav" || zone.ItemsFrom == "" {
				return
			}
			prev := append([]model.Item(nil), zone.Items...)
			loaded, err := a.readSource(zone.ItemsFrom, zone.ID)
			if err != nil {
				zone.Items = keepReturn(prev, nil)
				failed = err
				return
			}
			zone.Items = keepReturn(prev, loaded)
			if err := a.finishItems(win.ID, zone); err != nil {
				zone.Items = keepReturn(prev, nil)
				failed = err
			}
		})
	}
	return failed
}

func (a *app) readSource(source, where string) ([]model.Item, error) {
	if err := items.CheckSource(source, where); err != nil {
		return nil, err
	}
	if strings.HasPrefix(source, "cmd:") {
		return runScript(strings.TrimPrefix(source, "cmd:"), where)
	}
	return a.readItemFile(source, where)
}

func (a *app) readItemFile(file, where string) ([]model.Item, error) {
	raw, err := os.ReadFile(run.Path(a.root, file))
	if err != nil {
		return nil, model.Errf("fichier introuvable : %s", file)
	}
	return items.Parse(string(raw), where)
}

func runScript(script, where string) ([]model.Item, error) {
	out, err := exec.Command("sh", "-c", script).Output()
	if err != nil {
		return nil, model.Errf("commande échouée (%s)", where)
	}
	return items.Parse(string(out), where)
}

func keepReturn(prev, loaded []model.Item) []model.Item {
	for _, item := range loaded {
		if item.Focus.Kind == "return" || item.Enter.Kind == "return" {
			return loaded
		}
	}
	for _, item := range prev {
		if item.Focus.Kind == "return" || item.Enter.Kind == "return" {
			return append(loaded, item)
		}
	}
	return loaded
}

func (a *app) finishItems(winID string, zone *model.Zone) error {
	n := 0
	for i := range zone.Items {
		item := &zone.Items[i]
		if item.ID == "" {
			n++
			item.ID = fmt.Sprintf("%s-%d", zone.ID, n)
		}
		if err := a.resolveItem(item); err != nil {
			return err
		}
	}
	a.sel[key(winID, zone.ID)] = 0
	return nil
}

func (a *app) resolveItem(item *model.Item) error {
	var err error
	item.Label, err = run.Resolve(item.Label, a.root, item.ID)
	if err != nil {
		return err
	}
	item.Value, err = run.Resolve(item.Value, a.root, item.ID)
	if err != nil {
		return err
	}
	item.Focus.Value, err = run.Resolve(item.Focus.Value, a.root, item.ID)
	if err != nil {
		return err
	}
	item.Enter.Value, err = run.Resolve(item.Enter.Value, a.root, item.ID)
	return err
}

func (a *app) findZone(win int, id string) *model.Zone {
	var found *model.Zone
	a.wins[win].Grid.Each(func(zone *model.Zone) {
		if zone.ID == id {
			found = zone
		}
	})
	return found
}

func (a *app) syncZone(win int, id string) {
	zone := a.findZone(win, id)
	if zone == nil || win >= len(a.placed) {
		return
	}
	a.placed[win].EachBox(func(box *layout.Box) {
		if box.Zone.ID == id {
			box.Zone = *zone
		}
	})
}

func (a *app) clearChain(id string, seen map[string]bool) error {
	if id == "" {
		return nil
	}
	if seen[id] {
		return model.Errf("cycle de listes (%s)", id)
	}
	seen[id] = true
	zone := a.findZone(a.cur, id)
	if zone == nil || zone.Tag != "nav" {
		return model.Errf("cible absente : %s", id)
	}
	next := zone.NextNav
	zone.Items = nil
	a.sel[key(a.winID(), id)] = 0
	a.syncZone(a.cur, id)
	return a.clearChain(next, seen)
}

func (a *app) installItems(navID string, loaded []model.Item) error {
	zone := a.findZone(a.cur, navID)
	if zone == nil || zone.Tag != "nav" {
		return model.Errf("cible absente : %s", navID)
	}
	zone.Items = cloneItems(loaded)
	if err := a.finishItems(a.winID(), zone); err != nil {
		zone.Items = nil
		a.syncZone(a.cur, navID)
		return err
	}
	a.syncZone(a.cur, navID)
	return nil
}

func cloneItems(in []model.Item) []model.Item {
	if in == nil {
		return nil
	}
	out := make([]model.Item, len(in))
	for i := range in {
		out[i] = in[i]
		out[i].NextItems = cloneItems(in[i].NextItems)
	}
	return out
}

func (a *app) follow(zone model.Zone) (tea.Cmd, error) {
	if zone.NextNav == "" {
		return nil, nil
	}
	if err := a.clearChain(zone.NextNav, map[string]bool{}); err != nil {
		return nil, err
	}
	sel := a.sel[key(a.winID(), zone.ID)]
	if sel < 0 || sel >= len(zone.Items) {
		return nil, nil
	}
	item := zone.Items[sel]
	if len(item.NextItems) > 0 {
		if err := a.installItems(zone.NextNav, item.NextItems); err != nil {
			a.status = err.Error()
			return nil, err
		}
		return nil, nil
	}
	if item.Next == "" {
		return nil, nil
	}
	if strings.HasPrefix(item.Next, "cmd:") {
		return a.startList(zone.NextNav, strings.TrimPrefix(item.Next, "cmd:")), nil
	}
	loaded, err := a.readItemFile(item.Next, zone.NextNav)
	if err != nil {
		a.status = err.Error()
		return nil, nil
	}
	if err := a.installItems(zone.NextNav, loaded); err != nil {
		a.status = err.Error()
		return nil, err
	}
	return nil, nil
}

func (a *app) startList(navID, script string) tea.Cmd {
	if script == "" || strings.ContainsAny(script, "\t\n") {
		a.status = model.Errf("commande illégale (%s)", navID).Error()
		return nil
	}
	a.listGen++
	gen := a.listGen
	ctx, cancel := context.WithCancel(context.Background())
	a.listCancel = cancel
	a.listLabel = script
	a.status = spinFrames[0] + " En cours : " + script
	return tea.Batch(listSpinCmd(gen), execList(ctx, gen, navID, script))
}

func (a *app) stopList() {
	a.listGen++
	if a.listCancel != nil {
		a.listCancel()
		a.listCancel = nil
	}
}

func listSpinCmd(gen int) tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(time.Time) tea.Msg {
		return listSpin{gen: gen}
	})
}

func execList(ctx context.Context, gen int, nav, script string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.CommandContext(ctx, "sh", "-c", script)
		out, err := cmd.Output()
		if ctx.Err() != nil {
			return listDone{gen: gen, nav: nav, cancelled: true}
		}
		if err != nil {
			return listDone{gen: gen, nav: nav, err: model.Errf("commande échouée (%s)", nav)}
		}
		parsed, err := items.Parse(string(out), nav)
		if err != nil {
			return listDone{gen: gen, nav: nav, err: err}
		}
		return listDone{gen: gen, nav: nav, items: parsed}
	}
}

func (a *app) onListDone(msg listDone) {
	if msg.gen != a.listGen {
		return
	}
	a.listCancel = nil
	if msg.cancelled {
		a.status = "Annulé"
		return
	}
	if msg.err != nil {
		a.status = msg.err.Error()
		return
	}
	if err := a.installItems(msg.nav, msg.items); err != nil {
		a.status = err.Error()
		return
	}
	a.status = "Prêt"
}
