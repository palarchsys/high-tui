import { emptyAction, emptyItem, focusStyle, windowStyle, zoneStyle } from "./model";
import type { Action, Cell, Grid, GridRow, Item, Style, Win, Zone, ZoneTag } from "./types";

function esc(text: string): string {
  return text
    .replace(/&/g, "\u0026amp;")
    .replace(/</g, "\u0026lt;")
    .replace(/>/g, "\u0026gt;")
    .replace(/"/g, "\u0026quot;");
}

function attr(name: string, value: string): string {
  if (!value) return "";
  return ` ${name}="${esc(value)}"`;
}

function actionAttrs(prefix: string, action: Action): string {
  if (!action.kind) return "";
  return (
    ` ${prefix}="${action.kind}"` +
    attr(`${prefix}-target`, action.target) +
    attr(`${prefix}-file`, action.file) +
    attr(`${prefix}-name`, action.name) +
    attr(`${prefix}-href`, action.href) +
    attr(`${prefix}-run`, action.run) +
    attr(`${prefix}-output`, action.output) +
    attr(`${prefix}-value`, action.value)
  );
}

function emitItem(item: Item, pad: string): string {
  const inner = item.enter.href ? `<a href="${esc(item.enter.href)}">${esc(item.label)}</a>` : esc(item.label);
  return `${pad}<li id="${esc(item.id)}"${attr("data-label", item.label && item.enter.href ? "" : "")}${attr("data-value", item.value)}${actionAttrs("data-focus", item.focus)}${actionAttrs("data-enter", item.enter)}>${inner}</li>`;
}

function emitZone(zone: Zone, pad: string): string {
  const id = ` id="${esc(zone.id)}"`;
  const classes = zone.classes.length ? ` class="${esc(zone.classes.join(" "))}"` : "";
  const common =
    id +
    classes +
    (zone.readonly ? " readonly" : "") +
    attr("data-file", zone.file) +
    attr("name", zone.tag === "input" ? zone.name : "") +
    attr("data-name", zone.tag === "textarea" ? zone.name : "") +
    attr("data-save", zone.save) +
    attr("data-items", zone.itemsFrom) +
    attr("data-next", zone.nextNav) +
    actionAttrs("data-focus", zone.focus) +
    actionAttrs("data-enter", zone.enter);
  if (zone.tag === "header") return `${pad}<header${common}><h1>${esc(zone.label)}</h1></header>`;
  if (zone.tag === "footer") return `${pad}<footer${common}>${esc(zone.label)}</footer>`;
  if (zone.tag === "button") return `${pad}<button${common} type="button">${esc(zone.label)}</button>`;
  if (zone.tag === "input") return `${pad}<input${common}${attr("value", zone.value)}>`;
  if (zone.tag === "textarea") return `${pad}<textarea${common}>${esc(zone.value)}</textarea>`;
  const items = zone.items.map((item) => emitItem(item, `${pad}    `)).join("\n");
  return `${pad}<nav${common}>\n${pad}  <h2>${esc(zone.label)}</h2>\n${pad}  <ul>\n${items}\n${pad}  </ul>\n${pad}</nav>`;
}

function emitGrid(grid: Grid, pad: string): string {
  const rows = grid.rows
    .map((row) => {
      const cells = row.cells
        .map((cell) => {
          const fr = cell.fr !== 1 ? ` data-fr="${cell.fr}"` : "";
          const inner = cell.grid ? emitGrid(cell.grid, `${pad}      `) : cell.zone ? emitZone(cell.zone, `${pad}      `) : "";
          return `${pad}    <td${fr}>\n${inner}\n${pad}    </td>`;
        })
        .join("\n");
      const id = row.id ? ` id="${esc(row.id)}"` : "";
      const lines = row.lines !== 1 ? ` data-lines="${row.lines}"` : "";
      return `${pad}  <tr${id}${lines}>\n${cells}\n${pad}  </tr>`;
    })
    .join("\n");
  return `${pad}<table id="${esc(grid.id)}">\n${rows}\n${pad}</table>`;
}

function decl(style: Style, withMax: boolean): string[] {
  const lines = [
    `color: ${style.color}`,
    `background: ${style.background}`,
    `border: ${style.borderWidth}ch ${style.borderWidth ? "solid" : "none"} ${style.borderColor}`,
    `border-radius: ${style.radius}ch`,
    `padding: ${style.pad.map((n) => `${n}ch`).join(" ")}`,
    `margin: ${style.mar.map((n) => `${n}ch`).join(" ")}`,
    `line-height: ${style.lineHeight}`,
    `column-gap: ${style.gap}ch`,
    `font-weight: ${style.bold ? "bold" : "normal"}`,
    `--title-color: ${style.titleColor}`,
  ];
  if (withMax && style.maxWidth) lines.splice(6, 0, `max-width: ${style.maxWidth}ch`);
  return lines;
}

function rule(selector: string, lines: string[]): string {
  return `${selector} {\n${lines.map((line) => `  ${line};`).join("\n")}\n}`;
}

export function emitCss(win: Win): string {
  const zones: Zone[] = [];
  const walk = (grid: Grid) => {
    for (const row of grid.rows)
      for (const cell of row.cells) {
        if (cell.zone) zones.push(cell.zone);
        if (cell.grid) walk(cell.grid);
      }
  };
  walk(win.grid);
  const parts = [
    rule("body", decl(win.style, true)),
    rule("table", [
      "width: 100%",
      "border-collapse: collapse",
      "table-layout: fixed",
      `background: ${win.style.background}`,
    ]),
    "td { vertical-align: top; }",
    "h1, h2 { margin: 0; color: var(--title-color); }",
    "ul, li { margin: 0; padding: 0; list-style: none; }",
  ];
  for (const zone of zones) {
    parts.push(rule(`#${zone.id}`, decl(zone.base, false)));
    parts.push(rule(`#${zone.id}:focus`, decl(zone.focusOn, false)));
  }
  return `${parts.join("\n\n")}\n`;
}

export function emitHtml(win: Win): string {
  const css = win.file.replace(/\.html$/, ".css");
  return `<!DOCTYPE html>
<html lang="fr">
<head>
  <meta charset="utf-8">
  <title>${esc(win.label)}</title>
  <link rel="stylesheet" href="${esc(css)}">
</head>
<body>
  ${emitGrid(win.grid, "  ").trim()}
</body>
</html>
`;
}

type Node = { tag: string; attr: Record<string, string>; text: string; kids: Node[] };

const KNOWN = new Set([
  "id", "class", "lang", "charset", "rel", "href", "name", "value", "type", "readonly",
  "data-lines", "data-fr", "data-file", "data-name", "data-save", "data-label", "data-value",
  "data-focus", "data-focus-target", "data-focus-file", "data-focus-name", "data-focus-href",
  "data-focus-run", "data-focus-output", "data-focus-value",
  "data-enter", "data-enter-target", "data-enter-file", "data-enter-name", "data-enter-href",
  "data-enter-run", "data-enter-output", "data-enter-value",
  "data-items", "data-next",
]);

class Parser {
  i = 0;
  constructor(public s: string, private where: string) {}

  fail(message: string): Error {
    return new Error(`erreur : ${message} (${this.where})`);
  }

  skip() {
    for (;;) {
      while (this.i < this.s.length && /\s/.test(this.s[this.i])) this.i += 1;
      if (!this.s.startsWith("<!--", this.i)) return;
      const end = this.s.indexOf("-->", this.i + 4);
      if (end < 0) throw this.fail("commentaire non fermé");
      this.i = end + 3;
    }
  }

  name(): string {
    const start = this.i;
    while (this.i < this.s.length && /[a-z0-9-]/.test(this.s[this.i])) this.i += 1;
    return this.s.slice(start, this.i);
  }

  value(): string {
    const q = this.s[this.i];
    if (q !== '"' && q !== "'") throw this.fail("valeur non citée");
    this.i += 1;
    const start = this.i;
    while (this.i < this.s.length && this.s[this.i] !== q) {
      if (this.s[this.i] === "\n" || this.s[this.i] === "\t") throw this.fail("valeur illégale");
      this.i += 1;
    }
    if (this.i >= this.s.length) throw this.fail("guillemet non fermé");
    const val = this.s.slice(start, this.i);
    this.i += 1;
    return val;
  }

  element(): Node {
    if (this.s[this.i] !== "<") throw this.fail("balise attendue");
    this.i += 1;
    const tag = this.name();
    if (!tag) throw this.fail("balise attendue");
    const attr: Record<string, string> = {};
    for (;;) {
      this.skip();
      if (this.s.startsWith("/>", this.i)) {
        this.i += 2;
        return { tag, attr, text: "", kids: [] };
      }
      if (this.s[this.i] === ">") {
        this.i += 1;
        break;
      }
      const key = this.name();
      if (!key || !KNOWN.has(key)) throw this.fail(`attribut inconnu « ${key || "?"} »`);
      this.skip();
      if (this.s[this.i] !== "=") {
        if (key !== "readonly") throw this.fail(`valeur manquante « ${key} »`);
        if (key in attr) throw this.fail(`attribut en double « ${key} »`);
        attr[key] = "readonly";
        continue;
      }
      this.i += 1;
      this.skip();
      const val = this.value();
      if (!val && key !== "readonly") throw this.fail(`attribut vide « ${key} »`);
      if (key in attr) throw this.fail(`attribut en double « ${key} »`);
      attr[key] = val;
    }
    if (tag === "meta" || tag === "link" || tag === "input") return { tag, attr, text: "", kids: [] };
    let text = "";
    const kids: Node[] = [];
    for (;;) {
      if (this.i >= this.s.length) throw this.fail(`balise non fermée « ${tag} »`);
      if (this.s.startsWith("<!--", this.i)) {
        const end = this.s.indexOf("-->", this.i + 4);
        if (end < 0) throw this.fail("commentaire non fermé");
        this.i = end + 3;
        continue;
      }
      if (this.s.startsWith("</", this.i)) {
        this.i += 2;
        const end = this.name();
        this.skip();
        if (this.s[this.i] !== ">" || end !== tag) throw this.fail(`fermeture « ${tag} » invalide`);
        this.i += 1;
        return { tag, attr, text: text.trim(), kids };
      }
      if (this.s[this.i] === "<") {
        kids.push(this.element());
        continue;
      }
      const next = this.s.indexOf("<", this.i);
      if (next < 0) throw this.fail(`balise non fermée « ${tag} »`);
      text += this.s.slice(this.i, next);
      this.i = next;
    }
  }
}

function take(attr: Record<string, string>, key: string): string {
  const value = attr[key] ?? "";
  delete attr[key];
  return value;
}

function leftover(attr: Record<string, string>, where: string) {
  const key = Object.keys(attr)[0];
  if (key) throw new Error(`erreur : attribut refusé « ${key} » (${where})`);
}

function positive(raw: string, where: string): number {
  if (!/^[1-9][0-9]*$/.test(raw)) throw new Error(`erreur : nombre invalide « ${raw} » (${where})`);
  return Number(raw);
}

function actionOf(attr: Record<string, string>, prefix: string, where: string, loose: boolean): Action {
  const action: Action = {
    kind: take(attr, prefix) as Action["kind"],
    target: take(attr, `${prefix}-target`),
    file: take(attr, `${prefix}-file`),
    name: take(attr, `${prefix}-name`),
    href: take(attr, `${prefix}-href`),
    run: take(attr, `${prefix}-run`),
    output: take(attr, `${prefix}-output`),
    value: take(attr, `${prefix}-value`),
  };
  if (!action.kind) return action;
  const known = ["load", "edit", "set", "window", "return", "save", "refresh", "run"];
  if (!known.includes(action.kind)) throw new Error(`erreur : action inconnue « ${action.kind} » (${where})`);
  if ((action.kind === "load" || action.kind === "edit") && !action.target) {
    throw new Error(`erreur : cible manquante (${where})`);
  }
  if (action.kind === "load" && !action.file && !loose) throw new Error(`erreur : fichier manquant (${where})`);
  if (action.kind === "window" && !action.href) throw new Error(`erreur : href manquant (${where})`);
  if (action.kind === "run" && !action.run) throw new Error(`erreur : commande manquante (${where})`);
  if (action.kind === "set" && !loose) {
    if (!action.file) throw new Error(`erreur : fichier manquant (${where})`);
    if (!action.name) throw new Error(`erreur : variable manquante (${where})`);
    if (!action.value) throw new Error(`erreur : valeur manquante (${where})`);
  }
  return action;
}

function zoneFrom(node: Node, ids: Set<string>): Zone {
  const attr = { ...node.attr };
  const id = take(attr, "id");
  if (!id) throw new Error(`erreur : attribut manquant « id » (${node.tag})`);
  if (ids.has(id)) throw new Error(`erreur : identifiant en double : ${id}`);
  ids.add(id);
  const classes = take(attr, "class").split(/\s+/).filter(Boolean);
  const loose = node.tag === "nav";
  const focus = actionOf(attr, "data-focus", id, loose);
  const enter = actionOf(attr, "data-enter", id, loose);
  const readonly = take(attr, "readonly");
  const save = take(attr, "data-save");
  if (save && save !== "enter" && save !== "button") throw new Error(`erreur : save invalide (${id})`);
  const base = zoneStyle();
  const zone: Zone = {
    tag: node.tag as ZoneTag,
    id,
    classes,
    label: "",
    value: "",
    file: "",
    name: "",
    save: save as Zone["save"],
    readonly: readonly === "readonly" || readonly === "true",
    focus,
    enter,
    itemsFrom: "",
    nextNav: "",
    items: [],
    base,
    focusOn: focusStyle(base),
  };
  if (node.tag === "header") {
    leftover(attr, id);
    const h1 = node.kids[0];
    if (node.kids.length !== 1 || h1?.tag !== "h1" || !h1.text) throw new Error(`erreur : titre manquant (${id})`);
    zone.label = h1.text;
  } else if (node.tag === "nav") {
    zone.itemsFrom = take(attr, "data-items");
    zone.nextNav = take(attr, "data-next");
    leftover(attr, id);
    const [h2, ul] = node.kids;
    if (node.kids.length !== 2 || h2?.tag !== "h2" || ul?.tag !== "ul" || !h2.text) {
      throw new Error(`erreur : liste invalide (${id})`);
    }
    zone.label = h2.text;
    zone.items = (ul.kids ?? []).map((li) => itemFrom(li, ids));
  } else if (node.tag === "input") {
    const type = take(attr, "type");
    if (type && type !== "text") throw new Error(`erreur : type refusé (${id})`);
    zone.file = take(attr, "data-file");
    zone.name = take(attr, "data-name") || take(attr, "name");
    zone.value = take(attr, "value");
    leftover(attr, id);
  } else if (node.tag === "textarea") {
    if (node.kids.length) throw new Error(`erreur : enfant refusé (${id})`);
    zone.file = take(attr, "data-file");
    zone.name = take(attr, "data-name");
    zone.value = node.text;
    leftover(attr, id);
  } else if (node.tag === "button") {
    if (take(attr, "type") !== "button") throw new Error(`erreur : type refusé (${id})`);
    leftover(attr, id);
    if (!node.text) throw new Error(`erreur : label manquant (${id})`);
    if (!zone.enter.kind) throw new Error(`erreur : action manquante (${id})`);
    zone.label = node.text;
  } else if (node.tag === "footer") {
    leftover(attr, id);
    zone.label = node.text;
  } else throw new Error(`erreur : balise inconnue « ${node.tag} » (${id})`);
  return zone;
}

function itemFrom(node: Node, ids: Set<string>): Item {
  if (node.tag !== "li") throw new Error("erreur : balise inconnue « " + node.tag + " » (ul)");
  const attr = { ...node.attr };
  const id = take(attr, "id");
  if (!id) throw new Error("erreur : attribut manquant « id » (li)");
  if (ids.has(id)) throw new Error(`erreur : identifiant en double : ${id}`);
  ids.add(id);
  take(attr, "class");
  const focus = actionOf(attr, "data-focus", id, false);
  const enter = actionOf(attr, "data-enter", id, false);
  let label = take(attr, "data-label");
  const value = take(attr, "data-value");
  leftover(attr, id);
  let href = "";
  if (node.kids.length === 1 && node.kids[0].tag === "a") {
    href = node.kids[0].attr.href ?? "";
    if (!href) throw new Error(`erreur : href manquant (${id})`);
    if (!label) label = node.kids[0].text;
    if (enter.href === "") enter.href = href;
  } else if (!label) label = node.text;
  if (!label) throw new Error(`erreur : label manquant (${id})`);
  return emptyItem({ id, label, value, focus, enter });
}

function gridFrom(node: Node, ids: Set<string>): Grid {
  if (node.tag !== "table") throw new Error(`erreur : balise inconnue « ${node.tag} »`);
  const attr = { ...node.attr };
  const id = take(attr, "id");
  if (!id) throw new Error("erreur : attribut manquant « id » (table)");
  if (ids.has(id)) throw new Error(`erreur : identifiant en double : ${id}`);
  ids.add(id);
  take(attr, "class");
  leftover(attr, id);
  if (!node.kids.length) throw new Error(`erreur : grille sans rangée (${id})`);
  const rows: GridRow[] = node.kids.map((tr) => {
    if (tr.tag !== "tr") throw new Error(`erreur : balise inconnue « ${tr.tag} » (${id})`);
    const rowAttr = { ...tr.attr };
    const rowId = take(rowAttr, "id");
    if (rowId) {
      if (ids.has(rowId)) throw new Error(`erreur : identifiant en double : ${rowId}`);
      ids.add(rowId);
    }
    take(rowAttr, "class");
    const linesRaw = take(rowAttr, "data-lines");
    leftover(rowAttr, rowId || "tr");
    const cells: Cell[] = tr.kids.map((td) => {
      if (td.tag !== "td") throw new Error(`erreur : balise inconnue « ${td.tag} »`);
      const cellAttr = { ...td.attr };
      take(cellAttr, "id");
      const frRaw = take(cellAttr, "data-fr");
      leftover(cellAttr, rowId || "tr");
      const kid = td.kids[0];
      if (td.kids.length !== 1 || !kid) throw new Error(`erreur : cellule invalide (${rowId || "tr"})`);
      if (kid.tag === "table") return { fr: frRaw ? positive(frRaw, rowId || "tr") : 1, grid: gridFrom(kid, ids) };
      return { fr: frRaw ? positive(frRaw, rowId || "tr") : 1, zone: zoneFrom(kid, ids) };
    });
    return { id: rowId, lines: linesRaw ? positive(linesRaw, rowId || "tr") : 1, cells };
  });
  return { id, rows };
}

export function parseHtml(text: string, file: string): Win {
  const parser = new Parser(text, file);
  parser.skip();
  if (!parser.s.slice(parser.i, parser.i + 9).toLowerCase().startsWith("<!doctype")) {
    throw new Error(`erreur : doctype attendu (${file})`);
  }
  const end = parser.s.indexOf(">", parser.i);
  parser.i = end + 1;
  parser.skip();
  const root = parser.element();
  if (root.tag !== "html") throw new Error(`erreur : balise html attendue (${file})`);
  const [head, body] = root.kids;
  if (!head || !body || head.tag !== "head" || body.tag !== "body") {
    throw new Error(`erreur : structure html refusée (${file})`);
  }
  const title = head.kids.find((kid) => kid.tag === "title");
  if (!title?.text) throw new Error(`erreur : titre manquant (${file})`);
  const table = body.kids[0];
  if (body.kids.length !== 1 || table?.tag !== "table") throw new Error(`erreur : table attendue (${file})`);
  const grid = gridFrom(table, new Set());
  const win: Win = { id: grid.id, label: title.text, file, style: windowStyle(), grid };
  ensureReturn(win);
  return win;
}

export function ensureReturn(win: Win) {
  const zones: Zone[] = [];
  const walk = (grid: Grid) => {
    for (const row of grid.rows)
      for (const cell of row.cells) {
        if (cell.zone) zones.push(cell.zone);
        if (cell.grid) walk(cell.grid);
      }
  };
  walk(win.grid);
  const has = zones.some(
    (zone) =>
      zone.focus.kind === "return" ||
      zone.enter.kind === "return" ||
      zone.items.some((item) => item.focus.kind === "return" || item.enter.kind === "return"),
  );
  if (has) return;
  const nav = zones.find((zone) => zone.tag === "nav");
  if (!nav) return;
  nav.items.push(emptyItem({ id: `${nav.id}-retour`, label: "Retour", enter: { ...emptyAction(), kind: "return" } }));
}

function ch(raw: string): number | null {
  const match = raw.trim().match(/^(\d+)(?:ch)?$/);
  return match ? Number(match[1]) : null;
}

function quad(raw: string): [number, number, number, number] | null {
  const parts = raw.trim().split(/\s+/).map(ch);
  if (parts.some((part) => part == null)) return null;
  const n = parts as number[];
  if (n.length === 1) return [n[0], n[0], n[0], n[0]];
  if (n.length === 2) return [n[0], n[1], n[0], n[1]];
  if (n.length === 3) return [n[0], n[1], n[2], n[1]];
  if (n.length === 4) return [n[0]!, n[1]!, n[2]!, n[3]!];
  return null;
}

function paint(style: Style, decls: Record<string, string>) {
  if (decls.color && !decls.color.startsWith("var(")) style.color = decls.color;
  const bg = decls.background || decls["background-color"];
  if (bg) style.background = bg;
  if (decls["--title-color"]) style.titleColor = decls["--title-color"];
  if (decls["font-weight"]) style.bold = decls["font-weight"] === "bold";
  if (decls["line-height"]) style.lineHeight = Number(decls["line-height"]) || 1;
  if (decls["column-gap"]) style.gap = ch(decls["column-gap"]) ?? style.gap;
  if (decls["max-width"]) style.maxWidth = ch(decls["max-width"]) ?? style.maxWidth;
  if (decls["border-radius"]) style.radius = ch(decls["border-radius"]) ?? style.radius;
  if (decls.padding) style.pad = quad(decls.padding) ?? style.pad;
  if (decls.margin) style.mar = quad(decls.margin) ?? style.mar;
  if (decls.border) {
    const match = decls.border.match(/^(\d+)(?:ch)?\s+(solid|none)\s+(\S+)$/);
    if (match) {
      style.borderWidth = Number(match[1]);
      style.borderColor = match[3];
    }
  }
  if (decls["border-width"]) style.borderWidth = ch(decls["border-width"]) ?? style.borderWidth;
  if (decls["border-color"]) style.borderColor = decls["border-color"];
}

export function applyCss(win: Win, css: string): Win {
  const chunks = css
    .replace(/\/\*[\s\S]*?\*\//g, "")
    .split("}")
    .map((part) => part.trim())
    .filter(Boolean);
  const rules: { sel: string; decls: Record<string, string> }[] = [];
  for (const chunk of chunks) {
    const cut = chunk.indexOf("{");
    if (cut < 0) continue;
    const decls: Record<string, string> = {};
    for (const piece of chunk.slice(cut + 1).split(";")) {
      const colon = piece.indexOf(":");
      if (colon < 0) continue;
      decls[piece.slice(0, colon).trim()] = piece.slice(colon + 1).trim();
    }
    for (const sel of chunk.slice(0, cut).split(",")) rules.push({ sel: sel.trim(), decls });
  }
  const next: Win = structuredClone(win);
  const zones: Zone[] = [];
  const walk = (grid: Grid) => {
    for (const row of grid.rows)
      for (const cell of row.cells) {
        if (cell.zone) zones.push(cell.zone);
        if (cell.grid) walk(cell.grid);
      }
  };
  walk(next.grid);
  for (const rule of rules) {
    const focus = rule.sel.endsWith(":focus");
    const sel = focus ? rule.sel.slice(0, -6) : rule.sel;
    if (sel === "body") paint(next.style, rule.decls);
    if (sel.startsWith("#")) {
      const zone = zones.find((item) => item.id === sel.slice(1));
      if (zone) paint(focus ? zone.focusOn : zone.base, rule.decls);
    }
  }
  return next;
}
