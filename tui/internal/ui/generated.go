// Code généré par go generate. Ne pas modifier.

package ui

import "high-tui/internal/model"

func generatedWindows() []model.Window {
	return []model.Window{
		{
			ID: "main", Label: "high-tui", File: "main.html",
			Style: model.Style{Color: "#e6edf3", TitleColor: "", Background: "#0d1117", BorderColor: "#e6edf3", BorderWidth: 0, Radius: 0, LineHeight: 1, Gap: 0, MaxWidth: 260, Bold: false, Pad: [4]int{0, 0, 0, 0}, Mar: [4]int{0, 0, 0, 0}},
			Grid: model.Grid{
				ID: "main",
				Rows: []model.GridRow{
					{
						ID: "entete", Lines: 1,
						Cells: []model.Cell{
							{
								Fr: 1,
								Zone: &model.Zone{
									Tag: "header", ID: "titre", Label: "high-tui", Value: "",
									File: "", Name: "", Save: "", Href: "", Run: "", Output: "",
									Readonly: false, ItemsFrom: "", NextNav: "",
									Focus:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
									Enter:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
									Base:    model.Style{Color: "#e6edf3", TitleColor: "#ffa657", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 1, 1, 1}},
									FocusOn: model.Style{Color: "#e6edf3", TitleColor: "#ffa657", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 1, 1, 1}},
								},
							},
						},
					},
					{
						ID: "", Lines: 1,
						Cells: []model.Cell{
							{
								Fr: 1,
								Grid: &model.Grid{
									ID: "corps",
									Rows: []model.GridRow{
										{
											ID: "listes", Lines: 4,
											Cells: []model.Cell{
												{
													Fr: 1,
													Zone: &model.Zone{
														Tag: "nav", ID: "nav-rubriques", Label: "Rubriques", Value: "",
														File: "", Name: "", Save: "", Href: "", Run: "", Output: "",
														Readonly: false, ItemsFrom: "samples/rubriques.json", NextNav: "nav-fiches",
														Focus:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
														Enter:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
														Base:    model.Style{Color: "#e6edf3", TitleColor: "#ffa657", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 1, 1, 1}},
														FocusOn: model.Style{Color: "#e6edf3", TitleColor: "#ffa657", Background: "#161b22", BorderColor: "#58a6ff", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: true, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 1, 1, 1}},
														Items: []model.Item{
															{ID: "item-reseau", Label: "Réseau", Value: "", Focus: model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}, Enter: model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}},
															{ID: "item-notes", Label: "Notes", Value: "", Focus: model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}, Enter: model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}},
															{ID: "item-dossier", Label: "Dossier", Value: "", Focus: model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}, Enter: model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}},
															{ID: "nav-rubriques-retour", Label: "Retour", Value: "", Focus: model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}, Enter: model.Action{Kind: "return", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}},
														},
													},
												},
												{
													Fr: 1,
													Zone: &model.Zone{
														Tag: "nav", ID: "nav-fiches", Label: "Fiches", Value: "",
														File: "", Name: "", Save: "", Href: "", Run: "", Output: "",
														Readonly: false, ItemsFrom: "", NextNav: "nav-champs",
														Focus:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
														Enter:   model.Action{Kind: "set", Target: "valeur-hote", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
														Base:    model.Style{Color: "#e6edf3", TitleColor: "#ffa657", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 1, 1, 1}},
														FocusOn: model.Style{Color: "#e6edf3", TitleColor: "#ffa657", Background: "#161b22", BorderColor: "#58a6ff", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: true, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 1, 1, 1}},
														Items: []model.Item{
															{ID: "item-passage", Label: "Passage", Value: "", Focus: model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}, Enter: model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}},
															{ID: "item-archive", Label: "Archive", Value: "", Focus: model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}, Enter: model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}},
															{ID: "fiche-jeton", Label: "Jeton pro", Value: "", Focus: model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}, Enter: model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}},
														},
													},
												},
												{
													Fr: 1,
													Zone: &model.Zone{
														Tag: "nav", ID: "nav-champs", Label: "Champs", Value: "",
														File: "", Name: "", Save: "", Href: "", Run: "", Output: "",
														Readonly: false, ItemsFrom: "", NextNav: "",
														Focus:   model.Action{Kind: "load", Target: "valeur-hote", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
														Enter:   model.Action{Kind: "edit", Target: "valeur-hote", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
														Base:    model.Style{Color: "#e6edf3", TitleColor: "#ffa657", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 1, 1, 1}},
														FocusOn: model.Style{Color: "#e6edf3", TitleColor: "#ffa657", Background: "#161b22", BorderColor: "#58a6ff", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: true, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 1, 1, 1}},
														Items: []model.Item{
															{ID: "item-nom", Label: "Nom", Value: "", Focus: model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}, Enter: model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}},
															{ID: "item-hote", Label: "Hôte", Value: "", Focus: model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}, Enter: model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}},
															{ID: "item-port", Label: "Port", Value: "", Focus: model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}, Enter: model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}},
														},
													},
												},
												{
													Fr: 1,
													Zone: &model.Zone{
														Tag: "nav", ID: "nav-menu", Label: "Menu", Value: "",
														File: "", Name: "", Save: "", Href: "", Run: "", Output: "",
														Readonly: false, ItemsFrom: "", NextNav: "",
														Focus:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
														Enter:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
														Base:    model.Style{Color: "#e6edf3", TitleColor: "#ffa657", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 1, 1, 1}},
														FocusOn: model.Style{Color: "#e6edf3", TitleColor: "#ffa657", Background: "#161b22", BorderColor: "#58a6ff", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: true, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 1, 1, 1}},
														Items: []model.Item{
															{ID: "item-ouvrir", Label: "Ouvrir", Value: "", Focus: model.Action{Kind: "load", Target: "apercu", File: "samples/brouillon.txt", Name: "", Href: "", Run: "", Output: "", Value: ""}, Enter: model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}},
															{ID: "item-modifier", Label: "Modifier", Value: "", Focus: model.Action{Kind: "load", Target: "valeur-hote", File: "samples/passage.env", Name: "HOTE", Href: "", Run: "", Output: "", Value: ""}, Enter: model.Action{Kind: "edit", Target: "valeur-hote", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""}},
														},
													},
												},
											},
										},
										{
											ID: "edition", Lines: 3,
											Cells: []model.Cell{
												{
													Fr: 1,
													Zone: &model.Zone{
														Tag: "input", ID: "valeur-hote", Label: "", Value: "127.0.0.1",
														File: "samples/passage.env", Name: "HOTE", Save: "button", Href: "", Run: "", Output: "",
														Readonly: false, ItemsFrom: "", NextNav: "",
														Focus:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
														Enter:   model.Action{Kind: "edit", Target: "valeur-hote", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
														Base:    model.Style{Color: "#e6edf3", TitleColor: "#ffa657", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 1, 1, 1}},
														FocusOn: model.Style{Color: "#e6edf3", TitleColor: "#ffa657", Background: "#161b22", BorderColor: "#58a6ff", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: true, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 1, 1, 1}},
													},
												},
												{
													Fr: 2,
													Zone: &model.Zone{
														Tag: "textarea", ID: "apercu", Label: "", Value: "HOTE=127.0.0.1",
														File: "samples/passage.env", Name: "", Save: "", Href: "", Run: "", Output: "",
														Readonly: false, ItemsFrom: "", NextNav: "",
														Focus:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
														Enter:   model.Action{Kind: "edit", Target: "apercu", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
														Base:    model.Style{Color: "#e6edf3", TitleColor: "#ffa657", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 1, 1, 1}},
														FocusOn: model.Style{Color: "#e6edf3", TitleColor: "#ffa657", Background: "#161b22", BorderColor: "#58a6ff", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: true, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 1, 1, 1}},
													},
												},
												{
													Fr: 1,
													Zone: &model.Zone{
														Tag: "textarea", ID: "sortie", Label: "", Value: "Brouillon",
														File: "samples/brouillon.txt", Name: "", Save: "", Href: "", Run: "", Output: "",
														Readonly: true, ItemsFrom: "", NextNav: "",
														Focus:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
														Enter:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
														Base:    model.Style{Color: "#e6edf3", TitleColor: "#ffa657", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 1, 1, 1}},
														FocusOn: model.Style{Color: "#e6edf3", TitleColor: "#ffa657", Background: "#161b22", BorderColor: "#58a6ff", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: true, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 1, 1, 1}},
													},
												},
											},
										},
										{
											ID: "actions", Lines: 1,
											Cells: []model.Cell{
												{
													Fr: 1,
													Zone: &model.Zone{
														Tag: "button", ID: "enregistrer", Label: "Enregistrer", Value: "",
														File: "", Name: "", Save: "", Href: "", Run: "", Output: "",
														Readonly: false, ItemsFrom: "", NextNav: "",
														Focus:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
														Enter:   model.Action{Kind: "save", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
														Base:    model.Style{Color: "#e6edf3", TitleColor: "", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 1, 1, 1}},
														FocusOn: model.Style{Color: "#e6edf3", TitleColor: "", Background: "#161b22", BorderColor: "#58a6ff", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: true, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 1, 1, 1}},
													},
												},
												{
													Fr: 1,
													Zone: &model.Zone{
														Tag: "button", ID: "recharger", Label: "Recharger", Value: "",
														File: "", Name: "", Save: "", Href: "", Run: "", Output: "",
														Readonly: false, ItemsFrom: "", NextNav: "",
														Focus:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
														Enter:   model.Action{Kind: "refresh", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
														Base:    model.Style{Color: "#e6edf3", TitleColor: "", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 1, 1, 1}},
														FocusOn: model.Style{Color: "#e6edf3", TitleColor: "", Background: "#161b22", BorderColor: "#58a6ff", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: true, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 1, 1, 1}},
													},
												},
												{
													Fr: 1,
													Zone: &model.Zone{
														Tag: "button", ID: "essai", Label: "Essai local", Value: "",
														File: "", Name: "", Save: "", Href: "", Run: "", Output: "",
														Readonly: false, ItemsFrom: "", NextNav: "",
														Focus:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
														Enter:   model.Action{Kind: "run", Target: "", File: "", Name: "", Href: "", Run: "echo ok", Output: "sortie", Value: ""},
														Base:    model.Style{Color: "#e6edf3", TitleColor: "", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 1, 1, 1}},
														FocusOn: model.Style{Color: "#e6edf3", TitleColor: "", Background: "#161b22", BorderColor: "#58a6ff", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: true, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 1, 1, 1}},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					{
						ID: "barre", Lines: 1,
						Cells: []model.Cell{
							{
								Fr: 1,
								Zone: &model.Zone{
									Tag: "footer", ID: "statut", Label: "Prêt", Value: "",
									File: "", Name: "", Save: "", Href: "", Run: "", Output: "",
									Readonly: false, ItemsFrom: "", NextNav: "",
									Focus:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
									Enter:   model.Action{Kind: "", Target: "", File: "", Name: "", Href: "", Run: "", Output: "", Value: ""},
									Base:    model.Style{Color: "#e6edf3", TitleColor: "", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 1, 1, 1}},
									FocusOn: model.Style{Color: "#e6edf3", TitleColor: "", Background: "#161b22", BorderColor: "#30363d", BorderWidth: 1, Radius: 1, LineHeight: 1, Gap: 1, MaxWidth: 260, Bold: false, Pad: [4]int{0, 1, 0, 1}, Mar: [4]int{1, 1, 1, 1}},
								},
							},
						},
					},
				},
			},
		},
	}
}
