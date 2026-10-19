import { ensureReturn } from "./html";
import { emptyAction, emptyItem, makeZone, windowStyle } from "./model";
import type { Action, Item, Source, Win, Zone } from "./types";

function act(partial: Partial<Action>): Action {
  return { ...emptyAction(), ...partial };
}

function rowItem(partial: Partial<Item>): Item {
  return emptyItem(partial);
}

function zone(partial: Parameters<typeof makeZone>[2] & { tag: Zone["tag"]; id: string }): Zone {
  return makeZone(partial.tag, partial.id, partial);
}

const RUBRIQUES = `[
  {
    "label": "Réseau",
    "next": [
      {
        "label": "Passage",
        "value": "127.0.0.1",
        "file": "samples/passage.env",
        "name": "HOTE",
        "focus": {"kind": "load", "target": "valeur-hote"},
        "next": [
          {"label": "Nom", "file": "samples/passage.env", "name": "NOM"},
          {"label": "Hôte", "file": "samples/passage.env", "name": "HOTE"},
          {"label": "Port", "file": "samples/passage.env", "name": "PORT"}
        ]
      },
      {
        "label": "Archive",
        "value": "8080",
        "file": "samples/passage.env",
        "name": "PORT",
        "focus": {"kind": "load", "target": "valeur-hote"}
      },
      {
        "label": "Jeton pro",
        "file": "samples/passage.env",
        "name": "HOTE",
        "value": "samples/secrets.conf::HFU_PRO_JETON::value",
        "focus": {"kind": "load", "target": "valeur-hote"}
      }
    ]
  },
  {"label": "Notes", "next": "samples/fiches-notes.json"},
  {"label": "Dossier", "next": "cmd:printf '%s' '[{\\"label\\":\\"Echo\\"}]'"}
]
`;

const FICHES = `[
  {
    "label": "Brouillon",
    "enter": {"kind": "window", "href": "notes.html"}
  }
]
`;

const DOSSIER_CMD = `cmd:printf '%s' '[{"label":"Echo"}]'`;

export function exampleSources(): Source[] {
  return [
    { id: "rubriques", path: "samples/rubriques.json", body: RUBRIQUES, read: "json", separator: "" },
    { id: "fiches", path: "samples/fiches-notes.json", body: FICHES, read: "json", separator: "" },
    {
      id: "passage",
      path: "samples/passage.env",
      body: "# Passage\nNOM=Passage\nHOTE=127.0.0.1\nPORT=8080\n",
      read: "env",
      separator: "",
    },
    {
      id: "secrets",
      path: "samples/secrets.conf",
      body: "# label HFU_PRO_JETON Jeton pro\nHFU_PRO_JETON=jeton-demo\n",
      read: "env",
      separator: "",
    },
    {
      id: "brouillon",
      path: "samples/brouillon.txt",
      body: "Brouillon\nCorps de la note.\nÉtat : ouvert\n",
      read: "text",
      separator: "",
    },
    { id: "dossier-cmd", path: DOSSIER_CMD, body: '[{"label":"Echo"}]\n', read: "json", separator: "" },
  ];
}

export function exampleWindows(): Win[] {
  const main: Win = {
    id: "main",
    label: "high-tui",
    file: "main.html",
    style: windowStyle(),
    grid: {
      id: "main",
      rows: [
        {
          id: "entete",
          lines: 1,
          cells: [{ fr: 1, zone: zone({ tag: "header", id: "titre", label: "high-tui" }) }],
        },
        {
          id: "cadre",
          lines: 1,
          cells: [
            {
              fr: 1,
              grid: {
                id: "corps",
                rows: [
                  {
                    id: "listes",
                    lines: 4,
                    cells: [
                      {
                        fr: 1,
                        zone: zone({
                          tag: "nav",
                          id: "nav-rubriques",
                          label: "Rubriques",
                          itemsFrom: "samples/rubriques.json",
                          nextNav: "nav-fiches",
                          items: [
                            rowItem({ id: "item-reseau", label: "Réseau" }),
                            rowItem({ id: "item-notes", label: "Notes" }),
                            rowItem({ id: "item-dossier", label: "Dossier" }),
                          ],
                        }),
                      },
                      {
                        fr: 1,
                        zone: zone({
                          tag: "nav",
                          id: "nav-fiches",
                          label: "Fiches",
                          nextNav: "nav-champs",
                          enter: act({ kind: "set", target: "valeur-hote" }),
                          items: [
                            rowItem({ id: "item-passage", label: "Passage" }),
                            rowItem({ id: "item-archive", label: "Archive" }),
                            rowItem({ id: "fiche-jeton", label: "Jeton pro" }),
                          ],
                        }),
                      },
                      {
                        fr: 1,
                        zone: zone({
                          tag: "nav",
                          id: "nav-champs",
                          label: "Champs",
                          focus: act({ kind: "load", target: "valeur-hote" }),
                          enter: act({ kind: "edit", target: "valeur-hote" }),
                          items: [
                            rowItem({ id: "item-nom", label: "Nom" }),
                            rowItem({ id: "item-hote", label: "Hôte" }),
                            rowItem({ id: "item-port", label: "Port" }),
                          ],
                        }),
                      },
                      {
                        fr: 1,
                        zone: zone({
                          tag: "nav",
                          id: "nav-menu",
                          label: "Menu",
                          items: [
                            rowItem({
                              id: "item-ouvrir",
                              label: "Ouvrir",
                              focus: act({ kind: "load", target: "apercu", file: "samples/brouillon.txt" }),
                            }),
                            rowItem({
                              id: "item-modifier",
                              label: "Modifier",
                              focus: act({
                                kind: "load",
                                target: "valeur-hote",
                                file: "samples/passage.env",
                                name: "HOTE",
                              }),
                              enter: act({ kind: "edit", target: "valeur-hote" }),
                            }),
                          ],
                        }),
                      },
                    ],
                  },
                  {
                    id: "edition",
                    lines: 3,
                    cells: [
                      {
                        fr: 1,
                        zone: zone({
                          tag: "input",
                          id: "valeur-hote",
                          name: "HOTE",
                          value: "127.0.0.1",
                          file: "samples/passage.env",
                          save: "button",
                          enter: act({ kind: "edit", target: "valeur-hote" }),
                        }),
                      },
                      {
                        fr: 2,
                        zone: zone({
                          tag: "textarea",
                          id: "apercu",
                          file: "samples/passage.env",
                          value: "HOTE=127.0.0.1",
                          enter: act({ kind: "edit", target: "apercu" }),
                        }),
                      },
                      {
                        fr: 1,
                        zone: zone({
                          tag: "textarea",
                          id: "sortie",
                          readonly: true,
                          file: "samples/brouillon.txt",
                          value: "Brouillon",
                        }),
                      },
                    ],
                  },
                  {
                    id: "actions",
                    lines: 1,
                    cells: [
                      {
                        fr: 1,
                        zone: zone({
                          tag: "button",
                          id: "enregistrer",
                          label: "Enregistrer",
                          enter: act({ kind: "save" }),
                        }),
                      },
                      {
                        fr: 1,
                        zone: zone({
                          tag: "button",
                          id: "recharger",
                          label: "Recharger",
                          enter: act({ kind: "refresh" }),
                        }),
                      },
                      {
                        fr: 1,
                        zone: zone({
                          tag: "button",
                          id: "essai",
                          label: "Essai local",
                          enter: act({ kind: "run", run: "echo ok", output: "sortie" }),
                        }),
                      },
                    ],
                  },
                ],
              },
            },
          ],
        },
        {
          id: "barre",
          lines: 1,
          cells: [{ fr: 1, zone: zone({ tag: "footer", id: "statut", label: "Prêt" }) }],
        },
      ],
    },
  };
  ensureReturn(main);

  const notes: Win = {
    id: "notes",
    label: "Notes",
    file: "notes.html",
    style: windowStyle(),
    grid: {
      id: "notes",
      rows: [
        {
          id: "entete-notes",
          lines: 1,
          cells: [{ fr: 1, zone: zone({ tag: "header", id: "titre-notes", label: "Notes" }) }],
        },
        {
          id: "liste-notes",
          lines: 3,
          cells: [
            {
              fr: 1,
              zone: zone({
                tag: "nav",
                id: "nav-notes",
                label: "Notes",
                items: [
                  rowItem({
                    id: "item-brouillon",
                    label: "Brouillon",
                    focus: act({ kind: "load", target: "corps", file: "samples/brouillon.txt" }),
                  }),
                  rowItem({
                    id: "nav-notes-retour",
                    label: "Retour",
                    enter: act({ kind: "return", href: "main.html" }),
                  }),
                ],
              }),
            },
          ],
        },
        {
          id: "page-notes",
          lines: 4,
          cells: [
            {
              fr: 1,
              zone: zone({
                tag: "textarea",
                id: "corps",
                readonly: true,
                file: "samples/brouillon.txt",
                value: "Brouillon",
              }),
            },
          ],
        },
        {
          id: "barre-notes",
          lines: 1,
          cells: [{ fr: 1, zone: zone({ tag: "footer", id: "statut-notes", label: "Prêt" }) }],
        },
      ],
    },
  };
  return [main, notes];
}

export const DOSSIER_COMMAND = DOSSIER_CMD;
