import { emptyAction, emptyItem } from "./model";
import type { Action, Item, Source } from "./types";

const REF = /^(.+)::([A-Za-z_][A-Za-z0-9_]*)::([A-Za-z0-9_]+)$/;
const VAR = /^([A-Za-z_][A-Za-z0-9_]*)=(.*)$/;

export type Pair = { name: string; label: string; value: string };

export function parseEnv(body: string): Pair[] {
  const labels = new Map<string, string>();
  const pairs: Pair[] = [];
  const seen = new Set<string>();
  for (const raw of body.split("\n")) {
    const line = raw.replace(/\r$/, "");
    const label = line.match(/^# label ([A-Za-z_][A-Za-z0-9_]*) (.*)$/);
    if (label) {
      if (!labels.has(label[1])) labels.set(label[1], label[2]);
      continue;
    }
    const variable = line.match(VAR);
    if (variable && !seen.has(variable[1])) {
      seen.add(variable[1]);
      pairs.push({ name: variable[1], label: labels.get(variable[1]) || variable[1], value: variable[2] });
    }
  }
  return pairs;
}

export function writeEnv(body: string, pairs: Pair[]): string {
  const preamble: string[] = [];
  for (const raw of body.split("\n")) {
    const line = raw.replace(/\r$/, "");
    if (!line.trim()) continue;
    if (/^# label [A-Za-z_]/.test(line)) continue;
    if (VAR.test(line)) continue;
    preamble.push(line);
  }
  const lines = [...preamble];
  for (const pair of pairs) {
    if (pair.label && pair.label !== pair.name) lines.push(`# label ${pair.name} ${pair.label}`);
    lines.push(`${pair.name}=${pair.value}`);
  }
  return `${lines.join("\n")}\n`;
}

export function readVar(body: string, name: string): { label: string; value: string; found: boolean } {
  const pair = parseEnv(body).find((item) => item.name === name);
  if (!pair) return { label: "", value: "", found: false };
  return { label: pair.label, value: pair.value, found: true };
}

export function resolveRef(text: string, sources: Source[], where: string): { text: string; error?: string } {
  if ((text.match(/::/g) || []).length < 2) return { text };
  const match = REF.exec(text.trim());
  if (!match) return { text };
  const [, path, name, kind] = match;
  if (kind !== "label" && kind !== "value") {
    return { text: "", error: `erreur : le dernier mot n'est ni label ni value (${where})` };
  }
  const source = sources.find((item) => item.path === path);
  if (!source) return { text: "", error: `erreur : fichier introuvable : ${path}` };
  const pair = readVar(source.body, name);
  if (kind === "value") {
    if (!pair.found) return { text: "", error: `erreur : variable manquante (${where})` };
    return { text: pair.value };
  }
  if (!pair.label || pair.label === name) {
    const explicit = source.body.split("\n").some((line) => line.startsWith(`# label ${name} `));
    if (!explicit) return { text: "", error: `erreur : label manquant (${where})` };
  }
  return { text: pair.label };
}

function actionFrom(raw: Record<string, unknown> | undefined, where: string): { action: Action; error?: string } {
  if (!raw) return { action: emptyAction() };
  const action: Action = {
    kind: String(raw.kind ?? "") as Action["kind"],
    target: String(raw.target ?? ""),
    file: String(raw.file ?? ""),
    name: String(raw.name ?? ""),
    href: String(raw.href ?? ""),
    run: String(raw.run ?? ""),
    output: String(raw.output ?? ""),
    value: String(raw.value ?? ""),
  };
  const known = ["", "load", "edit", "set", "window", "return", "save", "refresh", "run"];
  if (!known.includes(action.kind)) return { action, error: `erreur : action inconnue « ${action.kind} » (${where})` };
  if (action.kind === "load" || action.kind === "edit") {
    if (!action.target) return { action, error: `erreur : cible manquante (${where})` };
  }
  if (action.kind === "window" && !action.href) return { action, error: `erreur : href manquant (${where})` };
  if (action.kind === "run" && !action.run) return { action, error: `erreur : commande manquante (${where})` };
  if (action.run && /[\t\n]/.test(action.run)) return { action, error: `erreur : commande illégale (${where})` };
  return { action };
}

const ITEM_FIELDS = new Set(["label", "value", "file", "name", "next", "focus", "enter"]);

function itemFrom(raw: Record<string, unknown>, where: string): { item: Item; error?: string } {
  for (const key of Object.keys(raw)) {
    if (!ITEM_FIELDS.has(key)) return { item: emptyItem(), error: `erreur : champ inconnu (${where})` };
  }
  const label = String(raw.label ?? "").trim();
  if (!label) return { item: emptyItem(), error: `erreur : label manquant (${where})` };
  const focus = actionFrom(raw.focus as Record<string, unknown> | undefined, where);
  if (focus.error) return { item: emptyItem(), error: focus.error };
  const enter = actionFrom(raw.enter as Record<string, unknown> | undefined, where);
  if (enter.error) return { item: emptyItem(), error: enter.error };
  const item = emptyItem({
    label,
    value: String(raw.value ?? ""),
    file: String(raw.file ?? ""),
    name: String(raw.name ?? ""),
    focus: focus.action,
    enter: enter.action,
  });
  const next = raw.next;
  if (next == null) return { item };
  if (typeof next === "string") {
    if (!next) return { item, error: `erreur : source manquante (${where})` };
    if (next.startsWith("cmd:")) {
      const script = next.slice(4);
      if (!script || /[\t\n]/.test(script)) return { item, error: `erreur : commande illégale (${where})` };
      item.nextKind = "cmd";
    } else item.nextKind = "file";
    item.next = next;
    return { item };
  }
  if (Array.isArray(next)) {
    const nested = itemsFromArray(next, where);
    if (nested.error) return { item: emptyItem(), error: nested.error };
    item.nextKind = "inline";
    item.nextItems = nested.items;
    return { item };
  }
  return { item, error: `erreur : JSON invalide (${where})` };
}

function itemsFromArray(rows: unknown[], where: string): { items: Item[]; error?: string } {
  const items: Item[] = [];
  for (const row of rows) {
    if (!row || typeof row !== "object" || Array.isArray(row)) return { items: [], error: `erreur : JSON invalide (${where})` };
    const parsed = itemFrom(row as Record<string, unknown>, where);
    if (parsed.error) return { items: [], error: parsed.error };
    items.push(parsed.item);
  }
  return { items };
}

export function parseItems(body: string, where: string): { items: Item[]; error?: string } {
  const text = body.trim();
  if (!text) return { items: [], error: `erreur : JSON invalide (${where})` };
  try {
    if (text.startsWith("[")) {
      const rows = JSON.parse(text) as unknown;
      if (!Array.isArray(rows)) return { items: [], error: `erreur : JSON invalide (${where})` };
      return itemsFromArray(rows, where);
    }
  } catch {
    return { items: [], error: `erreur : JSON invalide (${where})` };
  }
  const items: Item[] = [];
  for (const line of text.split("\n")) {
    const trimmed = line.trim();
    if (!trimmed) continue;
    try {
      const row = JSON.parse(trimmed) as unknown;
      if (!row || typeof row !== "object" || Array.isArray(row)) return { items: [], error: `erreur : JSON invalide (${where})` };
      const parsed = itemFrom(row as Record<string, unknown>, where);
      if (parsed.error) return { items: [], error: parsed.error };
      items.push(parsed.item);
    } catch {
      return { items: [], error: `erreur : JSON invalide (${where})` };
    }
  }
  if (items.length === 0) return { items: [], error: `erreur : JSON invalide (${where})` };
  return { items };
}

export function splitLines(body: string, separator: string): { label: string; value: string }[] {
  const rows: { label: string; value: string }[] = [];
  for (const raw of body.split("\n")) {
    const line = raw.replace(/\r$/, "").trim();
    if (!line) continue;
    if (!separator) {
      rows.push({ label: line, value: "" });
      continue;
    }
    const at = line.indexOf(separator);
    if (at < 0) rows.push({ label: line, value: "" });
    else rows.push({ label: line.slice(0, at).trim(), value: line.slice(at + separator.length).trim() });
  }
  return rows;
}

export function itemsFromSource(source: Source, where: string): { items: Item[]; error?: string } {
  if (source.read === "json") return parseItems(source.body, where);
  if (source.read === "env") {
    return {
      items: parseEnv(source.body).map((pair) =>
        emptyItem({ label: pair.label, value: pair.value, name: pair.name, file: source.path }),
      ),
    };
  }
  if (source.read === "lines") {
    return {
      items: splitLines(source.body, source.separator).map((row) => emptyItem({ label: row.label, value: row.value })),
    };
  }
  const items = source.body
    .split("\n")
    .map((line) => line.replace(/\r$/, "").trim())
    .filter(Boolean)
    .map((label) => emptyItem({ label }));
  return { items };
}

export function itemToRaw(item: Item): Record<string, unknown> {
  const raw: Record<string, unknown> = { label: item.label };
  if (item.value) raw.value = item.value;
  if (item.file) raw.file = item.file;
  if (item.name) raw.name = item.name;
  if (item.nextKind === "file" || item.nextKind === "cmd") raw.next = item.next;
  if (item.nextKind === "inline") raw.next = item.nextItems.map(itemToRaw);
  const focus = actionToRaw(item.focus);
  const enter = actionToRaw(item.enter);
  if (focus) raw.focus = focus;
  if (enter) raw.enter = enter;
  return raw;
}

function actionToRaw(action: Action): Record<string, string> | undefined {
  if (!action.kind) return undefined;
  const raw: Record<string, string> = { kind: action.kind };
  for (const key of ["target", "file", "name", "href", "run", "output", "value"] as const) {
    if (action[key]) raw[key] = action[key];
  }
  return raw;
}

export function writeItems(items: Item[]): string {
  return `${JSON.stringify(items.map(itemToRaw), null, 2)}\n`;
}

export function projectPath(path: string): string {
  const cut = path.replace(/\.[^.]+$/, "");
  return `${cut}.items.json`;
}

export function sortLabels(items: Item[], enabled: boolean): Item[] {
  if (!enabled) return items;
  return [...items].sort((a, b) => a.label.localeCompare(b.label, "fr"));
}
