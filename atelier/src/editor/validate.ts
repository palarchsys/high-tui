import { zonesOf } from "./model";
import type { Action, Source, Win, Zone } from "./types";

export function validate(win: Win, sources: Source[]): string[] {
  const errors: string[] = [];
  const zones = zonesOf(win);
  const byId = new Map<string, Zone>();
  const seen = new Set<string>();
  const take = (id: string, where: string) => {
    if (!id) errors.push(`erreur : attribut manquant « id » (${where})`);
    else if (seen.has(id)) errors.push(`erreur : identifiant en double : ${id}`);
    else seen.add(id);
  };
  take(win.grid.id, "table");
  const walk = (gridId: string, rows: Win["grid"]["rows"]) => {
    for (const row of rows) {
      if (row.id) take(row.id, "tr");
      if (row.lines < 1) errors.push(`erreur : nombre invalide (${row.id || gridId})`);
      if (!row.cells.length) errors.push(`erreur : rangée sans colonne (${row.id || gridId})`);
      for (const cell of row.cells) {
        if (cell.fr < 1) errors.push(`erreur : nombre invalide (${row.id || gridId})`);
        if (!cell.zone && !cell.grid) errors.push(`erreur : cellule invalide (${row.id || gridId})`);
        if (cell.grid) {
          take(cell.grid.id, "table");
          walk(cell.grid.id, cell.grid.rows);
        }
      }
    }
  };
  walk(win.grid.id, win.grid.rows);
  for (const zone of zones) {
    take(zone.id, zone.tag);
    for (const item of zone.items) take(item.id, "li");
    byId.set(zone.id, zone);
  }
  const fileOk = (path: string, where: string) => {
    if (!path) return;
    if (path.startsWith("cmd:")) {
      const script = path.slice(4);
      if (!script || /[\t\n]/.test(script)) errors.push(`erreur : commande illégale (${where})`);
      return;
    }
    if (!sources.some((source) => source.path === path)) errors.push(`erreur : fichier introuvable : ${path}`);
  };
  const one = (action: Action, loose: boolean, where: string) => {
    if (!action.kind) return;
    if (action.kind === "load" || action.kind === "edit") {
      const target = byId.get(action.target);
      if (!target) errors.push(`erreur : cible absente : ${action.target || where}`);
      else if (target.tag !== "input" && target.tag !== "textarea") errors.push(`erreur : cible incompatible : ${action.target}`);
      else if (action.kind === "edit" && target.tag === "textarea" && target.readonly) {
        errors.push(`erreur : cible non éditable : ${action.target}`);
      }
      if (action.kind === "load") {
        if (!action.file && !loose) errors.push(`erreur : fichier manquant (${where})`);
        fileOk(action.file, where);
      }
    }
    if (action.kind === "window" && !action.href) errors.push(`erreur : href manquant (${where})`);
    if (action.kind === "run") {
      if (!action.run) errors.push(`erreur : commande manquante (${where})`);
      if (action.output) {
        const target = byId.get(action.output);
        if (!target || target.tag !== "textarea") errors.push(`erreur : cible absente : ${action.output}`);
      }
    }
    if (action.kind === "set") {
      if (!loose) {
        if (!action.file) errors.push(`erreur : fichier manquant (${where})`);
        if (!action.name) errors.push(`erreur : variable manquante (${where})`);
        if (!action.value) errors.push(`erreur : valeur manquante (${where})`);
      }
      fileOk(action.file, where);
      if (action.target) {
        const target = byId.get(action.target);
        if (!target || (target.tag !== "input" && target.tag !== "textarea")) {
          errors.push(`erreur : cible incompatible : ${action.target}`);
        }
      }
    }
  };
  for (const zone of zones) {
    if (zone.tag === "button" && !zone.enter.kind) errors.push(`erreur : action manquante (${zone.id})`);
    if (zone.tag === "header" && !zone.label) errors.push(`erreur : titre manquant (${zone.id})`);
    if (zone.tag === "nav") {
      if (!zone.label) errors.push(`erreur : liste invalide (${zone.id})`);
      if (!zone.items.length && !zone.itemsFrom) errors.push(`erreur : liste invalide (${zone.id})`);
      fileOk(zone.itemsFrom, zone.id);
      if (zone.nextNav) {
        const target = byId.get(zone.nextNav);
        if (!target || target.tag !== "nav") errors.push(`erreur : cible absente : ${zone.nextNav}`);
      }
    }
    fileOk(zone.file, zone.id);
    one(zone.focus, zone.tag === "nav", zone.id);
    one(zone.enter, zone.tag === "nav", zone.id);
    for (const item of zone.items) {
      if (!item.label) errors.push(`erreur : label manquant (${item.id || zone.id})`);
      one(item.focus, false, item.id || zone.id);
      one(item.enter, false, item.id || zone.id);
    }
  }
  for (const zone of zones) {
    if (!zone.nextNav) continue;
    const seenNav = new Set<string>();
    let cur = zone.id;
    while (cur) {
      if (seenNav.has(cur)) {
        errors.push(`erreur : cycle de listes (${zone.id})`);
        break;
      }
      seenNav.add(cur);
      cur = byId.get(cur)?.nextNav ?? "";
    }
  }
  return errors;
}
