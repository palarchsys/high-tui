// Code généré par go generate. Ne pas modifier.

package ui

import "high-tui/internal/model"

func generatedWindows() []model.Window {
	return []model.Window{
		{
			ID: "main", Label: "high-tui", File: "main.html",
			Style: model.Style{Color: "#e6edf3", Background: "#0d1117", BorderColor: "#e6edf3", BorderWidth: 0, Radius: 0, LineHeight: 1, Gap: 0, MaxWidth: 260, Bold: false, Pad: [4]int{0, 0, 0, 0}, Mar: [4]int{0, 0, 0, 0}},
			Rows: []model.Row{
				{
					ID: "entete", Tag: "header", Lines: 1,
					Style: model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
					Zones: []model.Zone{
						{
							Tag: "header", ID: "titre", Label: "high-tui", Value: "",
							File: "", Name: "", Save: "", Href: "", Run: "", Output: "",
							Fr: 1, Readonly: false,
							Focus:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Enter:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Base:    model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
							FocusOn: model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
						},
					},
				},
				{
					ID: "listes", Tag: "main", Lines: 4,
					Style: model.Style{Color: "#e6edf3", Background: "#0d1117", BorderColor: "#e6edf3", BorderWidth: 0, Radius: 0, LineHeight: 1, Gap: 0, MaxWidth: 260, Bold: false, Pad: [4]int{0, 0, 0, 0}, Mar: [4]int{0, 0, 0, 0}},
					Zones: []model.Zone{
						{
							Tag: "nav", ID: "nav-rubriques", Label: "Rubriques", Value: "",
							File: "", Name: "", Save: "", Href: "", Run: "", Output: "",
							Fr: 1, Readonly: false,
							Focus:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Enter:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Base:    model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
							FocusOn: model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#58a6ff", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: true, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
							Items: []model.Item{
								{ID: "item-reseau", Label: "Réseau", Value: "", Focus: model.Action{Kind: "load", Target: "valeur-hote", File: "samples/passage.env", Name: "HOTE", Href: "", Run: "", Output: "", Value: ""}, Enter: model.Action{Kind: "edit", Target: "valeur-hote", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}},
								{ID: "item-notes", Label: "Notes", Value: "", Focus: model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}, Enter: model.Action{Kind: "window", Target: "", File: "", Name: "", Href: "notes.html", Run: "", Output: "", Value: ""}},
								{ID: "nav-rubriques-retour", Label: "Retour", Value: "", Focus: model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}, Enter: model.Action{Kind: "return", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}},
							},
						},
						{
							Tag: "select", ID: "nav-fiches", Label: "Fiches", Value: "",
							File: "", Name: "", Save: "", Href: "", Run: "", Output: "",
							Fr: 1, Readonly: false,
							Focus:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Enter:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Base:    model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
							FocusOn: model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#58a6ff", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: true, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
							Items: []model.Item{
								{ID: "item-passage", Label: "Passage", Value: "127.0.0.1", Focus: model.Action{Kind: "load", Target: "valeur-hote", File: "samples/passage.env", Name: "HOTE", Href: "", Run: "", Output: "", Value: ""}, Enter: model.Action{Kind: "set", Target: "valeur-hote", File: "samples/passage.env", Name: "HOTE", Href: "", Run: "", Output: "", Value: "127.0.0.1"}},
								{ID: "item-archive", Label: "Archive", Value: "8080", Focus: model.Action{Kind: "load", Target: "valeur-hote", File: "samples/passage.env", Name: "PORT", Href: "", Run: "", Output: "", Value: ""}, Enter: model.Action{Kind: "set", Target: "valeur-hote", File: "samples/passage.env", Name: "PORT", Href: "", Run: "", Output: "", Value: "8080"}},
								{ID: "fiche-jeton", Label: "samples/secrets.conf::HFU_PRO_JETON::label", Value: "samples/secrets.conf::HFU_PRO_JETON::value", Focus: model.Action{Kind: "load", Target: "valeur-hote", File: "samples/passage.env", Name: "HOTE", Href: "", Run: "", Output: "", Value: ""}, Enter: model.Action{Kind: "set", Target: "valeur-hote", File: "samples/passage.env", Name: "HOTE", Href: "", Run: "", Output: "", Value: "samples/secrets.conf::HFU_PRO_JETON::value"}},
							},
						},
						{
							Tag: "select", ID: "nav-champs", Label: "Champs", Value: "",
							File: "", Name: "", Save: "", Href: "", Run: "", Output: "",
							Fr: 1, Readonly: false,
							Focus:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Enter:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Base:    model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
							FocusOn: model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#58a6ff", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: true, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
							Items: []model.Item{
								{ID: "item-nom", Label: "Nom", Value: "", Focus: model.Action{Kind: "load", Target: "valeur-hote", File: "samples/passage.env", Name: "NOM", Href: "", Run: "", Output: "", Value: ""}, Enter: model.Action{Kind: "edit", Target: "valeur-hote", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}},
								{ID: "item-hote", Label: "Hôte", Value: "", Focus: model.Action{Kind: "load", Target: "valeur-hote", File: "samples/passage.env", Name: "HOTE", Href: "", Run: "", Output: "", Value: ""}, Enter: model.Action{Kind: "edit", Target: "valeur-hote", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}},
								{ID: "item-port", Label: "Port", Value: "", Focus: model.Action{Kind: "load", Target: "valeur-hote", File: "samples/passage.env", Name: "PORT", Href: "", Run: "", Output: "", Value: ""}, Enter: model.Action{Kind: "edit", Target: "valeur-hote", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}},
							},
						},
						{
							Tag: "select", ID: "nav-menu", Label: "Menu", Value: "",
							File: "", Name: "", Save: "", Href: "", Run: "", Output: "",
							Fr: 1, Readonly: false,
							Focus:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Enter:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Base:    model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
							FocusOn: model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#58a6ff", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: true, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
							Items: []model.Item{
								{ID: "item-ouvrir", Label: "Ouvrir", Value: "", Focus: model.Action{Kind: "load", Target: "apercu", File: "samples/brouillon.txt", Name: "", Href: "", Run: "", Output: "", Value: ""}, Enter: model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}},
								{ID: "item-modifier", Label: "Modifier", Value: "", Focus: model.Action{Kind: "load", Target: "valeur-hote", File: "samples/passage.env", Name: "HOTE", Href: "", Run: "", Output: "", Value: ""}, Enter: model.Action{Kind: "edit", Target: "valeur-hote", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}},
							},
						},
					},
				},
				{
					ID: "edition", Tag: "main", Lines: 3,
					Style: model.Style{Color: "#e6edf3", Background: "#0d1117", BorderColor: "#e6edf3", BorderWidth: 0, Radius: 0, LineHeight: 1, Gap: 0, MaxWidth: 260, Bold: false, Pad: [4]int{0, 0, 0, 0}, Mar: [4]int{0, 0, 0, 0}},
					Zones: []model.Zone{
						{
							Tag: "input", ID: "valeur-hote", Label: "", Value: "",
							File: "samples/passage.env", Name: "HOTE", Save: "button", Href: "", Run: "", Output: "",
							Fr: 1, Readonly: false,
							Focus:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Enter:   model.Action{Kind: "edit", Target: "valeur-hote", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Base:    model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
							FocusOn: model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#58a6ff", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: true, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
						},
						{
							Tag: "textarea", ID: "apercu", Label: "", Value: "",
							File: "samples/passage.env", Name: "", Save: "", Href: "", Run: "", Output: "",
							Fr: 2, Readonly: false,
							Focus:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Enter:   model.Action{Kind: "edit", Target: "apercu", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Base:    model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
							FocusOn: model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#58a6ff", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: true, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
						},
						{
							Tag: "textarea", ID: "sortie", Label: "", Value: "",
							File: "samples/brouillon.txt", Name: "", Save: "", Href: "", Run: "", Output: "",
							Fr: 1, Readonly: true,
							Focus:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Enter:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Base:    model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
							FocusOn: model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#58a6ff", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: true, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
						},
					},
				},
				{
					ID: "actions", Tag: "main", Lines: 1,
					Style: model.Style{Color: "#e6edf3", Background: "#0d1117", BorderColor: "#e6edf3", BorderWidth: 0, Radius: 0, LineHeight: 1, Gap: 0, MaxWidth: 260, Bold: false, Pad: [4]int{0, 0, 0, 0}, Mar: [4]int{0, 0, 0, 0}},
					Zones: []model.Zone{
						{
							Tag: "button", ID: "enregistrer", Label: "Enregistrer", Value: "",
							File: "", Name: "", Save: "", Href: "", Run: "", Output: "",
							Fr: 1, Readonly: false,
							Focus:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Enter:   model.Action{Kind: "save", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Base:    model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
							FocusOn: model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#58a6ff", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: true, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
						},
						{
							Tag: "button", ID: "recharger", Label: "Recharger", Value: "",
							File: "", Name: "", Save: "", Href: "", Run: "", Output: "",
							Fr: 1, Readonly: false,
							Focus:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Enter:   model.Action{Kind: "refresh", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Base:    model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
							FocusOn: model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#58a6ff", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: true, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
						},
						{
							Tag: "command", ID: "essai", Label: "Essai local", Value: "",
							File: "", Name: "", Save: "", Href: "", Run: "", Output: "",
							Fr: 1, Readonly: false,
							Focus:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Enter:   model.Action{Kind: "run", Target: "", File: "", Name: "", Href: "", Run: "echo ok", Output: "sortie", Value: ""},
							Base:    model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
							FocusOn: model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#58a6ff", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: true, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
						},
					},
				},
				{
					ID: "barre", Tag: "status", Lines: 1,
					Style: model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
					Zones: []model.Zone{
						{
							Tag: "status", ID: "statut", Label: "", Value: "",
							File: "", Name: "", Save: "", Href: "", Run: "", Output: "",
							Fr: 1, Readonly: false,
							Focus:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Enter:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Base:    model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
							FocusOn: model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
						},
					},
				},
			},
		},
		{
			ID: "notes", Label: "Notes", File: "notes.html",
			Style: model.Style{Color: "#e6edf3", Background: "#0d1117", BorderColor: "#e6edf3", BorderWidth: 0, Radius: 0, LineHeight: 1, Gap: 0, MaxWidth: 260, Bold: false, Pad: [4]int{0, 0, 0, 0}, Mar: [4]int{0, 0, 0, 0}},
			Rows: []model.Row{
				{
					ID: "entete-notes", Tag: "header", Lines: 1,
					Style: model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
					Zones: []model.Zone{
						{
							Tag: "header", ID: "titre-notes", Label: "Notes", Value: "",
							File: "", Name: "", Save: "", Href: "", Run: "", Output: "",
							Fr: 1, Readonly: false,
							Focus:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Enter:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Base:    model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
							FocusOn: model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
						},
					},
				},
				{
					ID: "liste-notes", Tag: "main", Lines: 3,
					Style: model.Style{Color: "#e6edf3", Background: "#0d1117", BorderColor: "#e6edf3", BorderWidth: 0, Radius: 0, LineHeight: 1, Gap: 0, MaxWidth: 260, Bold: false, Pad: [4]int{0, 0, 0, 0}, Mar: [4]int{0, 0, 0, 0}},
					Zones: []model.Zone{
						{
							Tag: "select", ID: "nav-notes", Label: "Notes", Value: "",
							File: "", Name: "", Save: "", Href: "", Run: "", Output: "",
							Fr: 1, Readonly: false,
							Focus:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Enter:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Base:    model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
							FocusOn: model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#58a6ff", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: true, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
							Items: []model.Item{
								{ID: "item-brouillon", Label: "Brouillon", Value: "", Focus: model.Action{Kind: "load", Target: "corps", File: "samples/brouillon.txt", Name: "", Href: "", Run: "", Output: "", Value: ""}, Enter: model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}},
								{ID: "nav-notes-retour", Label: "Retour", Value: "", Focus: model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}, Enter: model.Action{Kind: "return", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}},
							},
						},
					},
				},
				{
					ID: "page-notes", Tag: "main", Lines: 4,
					Style: model.Style{Color: "#e6edf3", Background: "#0d1117", BorderColor: "#e6edf3", BorderWidth: 0, Radius: 0, LineHeight: 1, Gap: 0, MaxWidth: 260, Bold: false, Pad: [4]int{0, 0, 0, 0}, Mar: [4]int{0, 0, 0, 0}},
					Zones: []model.Zone{
						{
							Tag: "textarea", ID: "corps", Label: "", Value: "",
							File: "samples/brouillon.txt", Name: "", Save: "", Href: "", Run: "", Output: "",
							Fr: 1, Readonly: true,
							Focus:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Enter:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Base:    model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
							FocusOn: model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#58a6ff", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: true, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
						},
					},
				},
				{
					ID: "barre-notes", Tag: "status", Lines: 1,
					Style: model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
					Zones: []model.Zone{
						{
							Tag: "status", ID: "statut-notes", Label: "", Value: "",
							File: "", Name: "", Save: "", Href: "", Run: "", Output: "",
							Fr: 1, Readonly: false,
							Focus:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Enter:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
							Base:    model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
							FocusOn: model.Style{Color: "#e6edf3", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 0, 1, 0}},
						},
					},
				},
			},
		},
	}
}
