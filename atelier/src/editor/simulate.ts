import { itemsFromSource, parseItems, resolveRef } from "./data";
import type { LiveZone } from "./layout";
import { zonesOf } from "./model";
import type { Action, Item, Source, Win, Zone } from "./types";

export type Sim = {
  live: Map<string, LiveZone>;
  status: string;
  errors: string[];
};

function cloneItem(item: Item): Item {
  return {
    ...item,
    focus: { ...item.focus },
    enter: { ...item.enter },
    nextItems: item.nextItems.map(cloneItem),
  };
}

function hasReturn(items: Item[]): boolean {
  return items.some((item) => item.focus.kind === "return" || item.enter.kind === "return");
}

function keepReturn(prev: Item[], loaded: Item[]): Item[] {
  if (hasReturn(loaded)) return loaded;
  const back = prev.find((item) => item.focus.kind === "return" || item.enter.kind === "return");
  return back ? [...loaded, cloneItem(back)] : loaded;
}

function fill(action: Action, item: Item): Action {
  return {
    ...action,
    file: action.file || item.file,
    name: action.name || item.name,
    value: action.value || item.value,
  };
}

function sourceOf(sources: Source[], path: string): Source | undefined {
  return sources.find((source) => source.path === path);
}

export function simulate(win: Win, sources: Source[], navSel: Record<string, number>, focusId: string | null): Sim {
  const errors: string[] = [];
  const note = (error?: string) => {
    if (error && !errors.includes(error)) errors.push(error);
  };
  const zones = zonesOf(win);
  const byId = new Map(zones.map((zone) => [zone.id, zone]));
  const items = new Map<string, Item[]>();
  for (const zone of zones) items.set(zone.id, zone.items.map(cloneItem));

  for (const zone of zones) {
    if (zone.tag !== "nav" || !zone.itemsFrom) continue;
    const source = sourceOf(sources, zone.itemsFrom);
    if (!source) {
      note(`erreur : fichier introuvable : ${zone.itemsFrom}`);
      items.set(zone.id, keepReturn(zone.items, []));
      continue;
    }
    const parsed = zone.itemsFrom.startsWith("cmd:") ? parseItems(source.body, zone.id) : parseItems(source.body, zone.id);
    if (parsed.error) {
      note(parsed.error);
      items.set(zone.id, keepReturn(zone.items, []));
      continue;
    }
    items.set(zone.id, stamp(zone.id, keepReturn(zone.items, parsed.items), sources, note));
  }

  for (let pass = 0; pass < zones.length; pass += 1) {
    for (const zone of zones) {
      if (zone.tag !== "nav" || !zone.nextNav) continue;
      const cycle = clearChain(zone.nextNav, byId, items, new Set());
      note(cycle || undefined);
      const rows = items.get(zone.id) ?? [];
      if (!rows.length) continue;
      const sel = Math.min(Math.max(0, navSel[zone.id] ?? 0), rows.length - 1);
      const item = rows[sel];
      if (item.nextKind === "inline") {
        items.set(zone.nextNav, stamp(zone.nextNav, item.nextItems.map(cloneItem), sources, note));
        continue;
      }
      if (item.nextKind === "none" || !item.next) continue;
      if (item.next.startsWith("cmd:")) {
        const fixture = sourceOf(sources, item.next);
        if (!fixture) continue;
        const parsed = parseItems(fixture.body, zone.nextNav);
        note(parsed.error);
        if (!parsed.error) items.set(zone.nextNav, stamp(zone.nextNav, parsed.items, sources, note));
        continue;
      }
      const file = sourceOf(sources, item.next);
      if (!file) {
        note(`erreur : fichier introuvable : ${item.next}`);
        continue;
      }
      const parsed = parseItems(file.body, zone.nextNav);
      note(parsed.error);
      if (!parsed.error) items.set(zone.nextNav, stamp(zone.nextNav, parsed.items, sources, note));
    }
  }

  const overlay = new Map<string, { value: string; name: string }>();
  const focus = focusId ? byId.get(focusId) : undefined;
  if (focus) {
    let action = focus.focus;
    if (focus.tag === "nav") {
      const rows = items.get(focus.id) ?? [];
      const sel = Math.min(Math.max(0, navSel[focus.id] ?? 0), Math.max(0, rows.length - 1));
      const item = rows[sel];
      if (item) {
        const chosen = item.focus.kind ? item.focus : focus.focus;
        action = fill(chosen, item);
      }
    }
    if (action.kind === "load" && action.target) {
      const file = sourceOf(sources, action.file);
      if (!action.file || !file) note(action.file ? `erreur : fichier introuvable : ${action.file}` : `erreur : fichier manquant (${action.target})`);
      else if (action.name) {
        const pair = itemsFromSource({ ...file, read: "env" }, action.target);
        const found = pair.items.find((item) => item.name === action.name);
        overlay.set(action.target, { value: found?.value ?? "", name: action.name });
        if (!found) note(`erreur : variable manquante (${action.target})`);
      } else overlay.set(action.target, { value: file.body.replace(/\n$/, ""), name: "" });
    }
  }

  const live = new Map<string, LiveZone>();
  for (const zone of zones) {
    const label = resolved(zone.label, sources, zone.id, note);
    let value = zone.value;
    let name = zone.name;
    const over = overlay.get(zone.id);
    if (over) {
      value = over.value;
      name = over.name || name;
    } else if (zone.file && (zone.tag === "input" || zone.tag === "textarea")) {
      const file = sourceOf(sources, zone.file);
      if (!file) note(`erreur : fichier introuvable : ${zone.file}`);
      else if (zone.tag === "input" && zone.name) {
        const found = itemsFromSource({ ...file, read: "env" }, zone.id).items.find((item) => item.name === zone.name);
        value = found?.value ?? "";
        if (!found) note(`erreur : variable manquante (${zone.id})`);
      } else if (zone.tag === "textarea" && !zone.name) value = file.body.replace(/\n$/, "");
    }
    value = resolved(value, sources, zone.id, note);
    const rows = (items.get(zone.id) ?? []).map((item) => ({
      id: item.id,
      label: resolved(item.label, sources, item.id || zone.id, note),
    }));
    live.set(zone.id, { label, value, name, items: rows });
  }
  const status = errors[0] ?? "Prêt";
  return { live, status, errors };
}

function stamp(navId: string, rows: Item[], sources: Source[], note: (error?: string) => void): Item[] {
  let n = 0;
  return rows.map((item) => {
    const next = cloneItem(item);
    if (!next.id) {
      n += 1;
      next.id = `${navId}-${n}`;
    }
    next.label = resolved(next.label, sources, next.id, note);
    next.value = resolved(next.value, sources, next.id, note);
    next.focus.value = resolved(next.focus.value, sources, next.id, note);
    next.enter.value = resolved(next.enter.value, sources, next.id, note);
    return next;
  });
}

function resolved(text: string, sources: Source[], where: string, note: (error?: string) => void): string {
  const out = resolveRef(text, sources, where);
  note(out.error);
  return out.error ? text : out.text;
}

function clearChain(id: string, zones: Map<string, Zone>, items: Map<string, Item[]>, seen: Set<string>): string | null {
  if (!id) return null;
  if (seen.has(id)) return `erreur : cycle de listes (${id})`;
  seen.add(id);
  const zone = zones.get(id);
  if (!zone || zone.tag !== "nav") return `erreur : cible absente : ${id}`;
  items.set(id, []);
  return clearChain(zone.nextNav, zones, items, seen);
}
