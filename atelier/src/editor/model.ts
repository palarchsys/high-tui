import type { Action, Cell, Grid, GridRow, Item, Style, Win, Zone, ZoneTag } from "./types";

export function emptyAction(): Action {
  return { kind: "", target: "", file: "", name: "", href: "", run: "", output: "", value: "" };
}

export function windowStyle(): Style {
  return {
    color: "#e6edf3",
    titleColor: "#ffa657",
    background: "#0d1117",
    borderColor: "#e6edf3",
    borderWidth: 0,
    radius: 0,
    lineHeight: 1,
    gap: 0,
    maxWidth: 260,
    bold: false,
    pad: [0, 0, 0, 0],
    mar: [0, 0, 0, 0],
  };
}

export function zoneStyle(): Style {
  return {
    color: "#e6edf3",
    titleColor: "#ffa657",
    background: "#161b22",
    borderColor: "#30363d",
    borderWidth: 1,
    radius: 1,
    lineHeight: 1,
    gap: 1,
    maxWidth: 0,
    bold: false,
    pad: [0, 1, 0, 1],
    mar: [1, 1, 1, 1],
  };
}

export function focusStyle(base: Style): Style {
  return { ...base, borderColor: "#58a6ff", bold: true, maxWidth: 0 };
}

export function emptyItem(partial: Partial<Item> = {}): Item {
  return {
    id: "",
    label: "",
    value: "",
    file: "",
    name: "",
    nextKind: "none",
    next: "",
    nextItems: [],
    focus: emptyAction(),
    enter: emptyAction(),
    ...partial,
  };
}

export function makeZone(tag: ZoneTag, id: string, partial: Partial<Zone> = {}): Zone {
  const base = zoneStyle();
  const zone: Zone = {
    tag,
    id,
    classes: [],
    label: "",
    value: "",
    file: "",
    name: "",
    save: "",
    readonly: false,
    focus: emptyAction(),
    enter: emptyAction(),
    itemsFrom: "",
    nextNav: "",
    items: [],
    base,
    focusOn: focusStyle(base),
    ...partial,
  };
  if (tag === "header" && !zone.label) zone.label = "Titre";
  if (tag === "footer" && !zone.label) zone.label = "Prêt";
  if (tag === "button") {
    if (!zone.label) zone.label = "Action";
    if (!zone.enter.kind) zone.enter = { ...emptyAction(), kind: "save" };
  }
  if (tag === "nav" && zone.items.length === 0 && !partial.items) {
    zone.items = [emptyItem({ id: `${id}-1`, label: "Ligne" })];
  }
  return zone;
}

export function eachCell(grid: Grid, fn: (cell: Cell, row: GridRow, grid: Grid) => void) {
  for (const row of grid.rows) {
    for (const cell of row.cells) {
      fn(cell, row, grid);
      if (cell.grid) eachCell(cell.grid, fn);
    }
  }
}

export function eachZone(grid: Grid, fn: (zone: Zone, row: GridRow, cell: Cell) => void) {
  eachCell(grid, (cell, row) => {
    if (cell.zone) fn(cell.zone, row, cell);
  });
}

export function zonesOf(win: Win): Zone[] {
  const out: Zone[] = [];
  eachZone(win.grid, (zone) => out.push(zone));
  return out;
}

export function findZone(win: Win, id: string): Zone | undefined {
  return zonesOf(win).find((zone) => zone.id === id);
}

export function collectIds(win: Win): string[] {
  const ids: string[] = [];
  const walk = (grid: Grid) => {
    ids.push(grid.id);
    for (const row of grid.rows) {
      if (row.id) ids.push(row.id);
      for (const cell of row.cells) {
        if (cell.zone) {
          ids.push(cell.zone.id);
          for (const item of cell.zone.items) if (item.id) ids.push(item.id);
        }
        if (cell.grid) walk(cell.grid);
      }
    }
  };
  walk(win.grid);
  return ids;
}

export function freshId(win: Win, prefix: string): string {
  const ids = new Set(collectIds(win));
  let n = 1;
  while (ids.has(`${prefix}-${n}`)) n += 1;
  return `${prefix}-${n}`;
}

export function mapWin(win: Win, fn: (grid: Grid) => Grid): Win {
  return { ...win, grid: fn(win.grid) };
}

export function updateGrid(win: Win, gridId: string, fn: (grid: Grid) => Grid): Win {
  const walk = (grid: Grid): Grid => {
    const next = grid.id === gridId ? fn(grid) : grid;
    return {
      ...next,
      rows: next.rows.map((row) => ({
        ...row,
        cells: row.cells.map((cell) => (cell.grid ? { ...cell, grid: walk(cell.grid) } : cell)),
      })),
    };
  };
  return { ...win, grid: walk(win.grid) };
}

export function mapZones(win: Win, fn: (zone: Zone) => Zone): Win {
  const walk = (grid: Grid): Grid => ({
    ...grid,
    rows: grid.rows.map((row) => ({
      ...row,
      cells: row.cells.map((cell) => {
        if (cell.zone) return { ...cell, zone: fn(cell.zone) };
        if (cell.grid) return { ...cell, grid: walk(cell.grid) };
        return cell;
      }),
    })),
  });
  return { ...win, grid: walk(win.grid) };
}

export function replaceZone(win: Win, id: string, zone: Zone): Win {
  return mapZones(win, (current) => (current.id === id ? zone : current));
}

export function findGrid(win: Win, id: string): Grid | undefined {
  let found: Grid | undefined;
  const walk = (grid: Grid) => {
    if (grid.id === id) found = grid;
    for (const row of grid.rows) for (const cell of row.cells) if (cell.grid) walk(cell.grid);
  };
  walk(win.grid);
  return found;
}

export function gridOfZone(win: Win, zoneId: string): string {
  let found = win.grid.id;
  const walk = (grid: Grid) => {
    for (const row of grid.rows) {
      for (const cell of row.cells) {
        if (cell.zone?.id === zoneId) found = grid.id;
        if (cell.grid) walk(cell.grid);
      }
    }
  };
  walk(win.grid);
  return found;
}

export function rowOfZone(win: Win, zoneId: string): GridRow | undefined {
  let found: GridRow | undefined;
  eachZone(win.grid, (zone, row) => {
    if (zone.id === zoneId) found = row;
  });
  return found;
}
