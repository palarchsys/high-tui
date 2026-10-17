// Package build reads the HTML and CSS of every reached window.
package build

import (
	"os"
	"path/filepath"
	"strings"

	"high-tui/internal/css"
	"high-tui/internal/html"
	"high-tui/internal/items"
	"high-tui/internal/model"
)

// Load reads windows/main.html and every window it opens.
func Load(root string) ([]model.Window, error) {
	var wins []model.Window
	seen := map[string]bool{}
	jsonSeen := map[string]bool{}
	if err := loadFile(root, filepath.Join(root, "windows", "main.html"), &wins, seen, jsonSeen); err != nil {
		return nil, err
	}
	if len(wins) == 0 {
		return nil, model.Errf("aucune fenêtre")
	}
	return wins, nil
}

func loadFile(root, path string, wins *[]model.Window, seen, jsonSeen map[string]bool) error {
	path = filepath.Clean(path)
	if seen[path] {
		return nil
	}
	seen[path] = true
	raw, err := os.ReadFile(path)
	if err != nil {
		return model.Errf("fichier introuvable : %s", filepath.Base(path))
	}
	doc, err := html.Parse(string(raw), filepath.Base(path))
	if err != nil {
		return err
	}
	win, err := html.Window(doc, filepath.Base(path))
	if err != nil {
		return err
	}
	for _, have := range *wins {
		if have.ID == win.ID {
			return model.Errf("identifiant de fenêtre en double : %s", win.ID)
		}
	}
	cssPath := strings.TrimSuffix(path, filepath.Ext(path)) + ".css"
	cssRaw, err := os.ReadFile(cssPath)
	if err != nil {
		return model.Errf("fichier introuvable : %s", filepath.Base(cssPath))
	}
	rules, err := css.Parse(string(cssRaw), filepath.Base(cssPath))
	if err != nil {
		return err
	}
	inherit := css.Blank()
	win.Style, err = css.Resolve(rules, "body", nil, "", false, inherit)
	if err != nil {
		return err
	}
	parent := model.Style{Color: win.Style.Color, Background: win.Style.Background}
	var hrefs []string
	for _, zone := range win.Grid.Zones() {
		zone.Base, err = css.Resolve(rules, zone.Tag, zone.Classes, zone.ID, false, parent)
		if err != nil {
			return err
		}
		zone.FocusOn, err = css.Resolve(rules, zone.Tag, zone.Classes, zone.ID, true, parent)
		if err != nil {
			return err
		}
		if zone.File != "" {
			if err := exist(root, zone.File); err != nil {
				return err
			}
		}
		if zone.ItemsFrom != "" && !strings.HasPrefix(zone.ItemsFrom, "cmd:") {
			if err := exist(root, zone.ItemsFrom); err != nil {
				return err
			}
			more, err := hrefsFromSource(root, zone.ItemsFrom, jsonSeen)
			if err != nil {
				return err
			}
			hrefs = append(hrefs, more...)
		}
		hrefs = append(hrefs, hrefsOf(zone.Focus)...)
		hrefs = append(hrefs, hrefsOf(zone.Enter)...)
		for _, item := range zone.Items {
			hrefs = append(hrefs, hrefsOf(item.Focus)...)
			hrefs = append(hrefs, hrefsOf(item.Enter)...)
			if item.Focus.File != "" {
				if err := exist(root, item.Focus.File); err != nil {
					return err
				}
			}
			if item.Enter.File != "" {
				if err := exist(root, item.Enter.File); err != nil {
					return err
				}
			}
		}
	}
	*wins = append(*wins, win)
	dir := filepath.Dir(path)
	for _, href := range hrefs {
		next := href
		if !filepath.IsAbs(next) {
			next = filepath.Join(dir, href)
		}
		if err := loadFile(root, next, wins, seen, jsonSeen); err != nil {
			return err
		}
	}
	return nil
}

func hrefsOf(act model.Action) []string {
	if act.Kind == "window" && act.Href != "" {
		return []string{act.Href}
	}
	return nil
}

// hrefsFromSource reads a JSON list and every file it names. A cmd: source is left for runtime.
// Window hrefs stay relative to the HTML file that named the list, not to the JSON file.
func hrefsFromSource(root, source string, seen map[string]bool) ([]string, error) {
	if source == "" || strings.HasPrefix(source, "cmd:") {
		return nil, nil
	}
	path := source
	if !filepath.IsAbs(path) {
		path = filepath.Join(root, source)
	}
	path = filepath.Clean(path)
	if seen[path] {
		return nil, nil
	}
	seen[path] = true
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, model.Errf("fichier introuvable : %s", source)
	}
	list, err := items.Parse(string(raw), source)
	if err != nil {
		return nil, err
	}
	return hrefsFromItems(root, list, seen)
}

func hrefsFromItems(root string, list []model.Item, seen map[string]bool) ([]string, error) {
	var hrefs []string
	for _, item := range list {
		hrefs = append(hrefs, hrefsOf(item.Focus)...)
		hrefs = append(hrefs, hrefsOf(item.Enter)...)
		if len(item.NextItems) > 0 {
			more, err := hrefsFromItems(root, item.NextItems, seen)
			if err != nil {
				return nil, err
			}
			hrefs = append(hrefs, more...)
		}
		more, err := hrefsFromSource(root, item.Next, seen)
		if err != nil {
			return nil, err
		}
		hrefs = append(hrefs, more...)
	}
	return hrefs, nil
}

func exist(root, file string) error {
	path := file
	if !filepath.IsAbs(path) {
		path = filepath.Join(root, file)
	}
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		return model.Errf("fichier introuvable : %s", file)
	}
	return nil
}
