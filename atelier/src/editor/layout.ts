import type { Grid, GridRow, Style, Win, ZoneTag } from "./types";

export type LiveItem = { id: string; label: string };

export type LiveZone = {
  label: string;
  value: string;
  name: string;
  items: LiveItem[];
};

export type Box = {
  id: string;
  x: number;
  y: number;
  w: number;
  h: number;
  textW: number;
  textH: number;
  style: Style;
  tag: ZoneTag;
  title: boolean;
  lines: string[];
  itemIds: string[];
  selected: number;
};

function edge(row: GridRow, side: 0 | 2): number {
  let margin = 0;
  for (const cell of row.cells) {
    if (cell.zone) margin = Math.max(margin, cell.zone.base.mar[side]);
    if (cell.grid?.rows.length) {
      const sub = side === 0 ? cell.grid.rows[0] : cell.grid.rows[cell.grid.rows.length - 1];
      margin = Math.max(margin, edge(sub, side));
    }
  }
  return margin;
}

function flank(row: GridRow, side: 1 | 3): number {
  let margin = 0;
  for (const cell of row.cells) if (cell.zone) margin = Math.max(margin, cell.zone.base.mar[side]);
  return margin;
}

function gap(row: GridRow, index: number): number {
  const left = row.cells[index]?.zone?.base;
  const right = row.cells[index + 1]?.zone?.base;
  return Math.max(left?.gap ?? 0, right?.gap ?? 0, left?.mar[1] ?? 0, right?.mar[3] ?? 0);
}

function columns(row: GridRow, width: number): number[] {
  const count = row.cells.length;
  if (!count) return [];
  let gaps = 0;
  for (let i = 0; i < count - 1; i += 1) gaps += gap(row, i);
  let useful = width - flank(row, 3) - flank(row, 1) - gaps;
  if (useful < count) useful = count;
  const sum = row.cells.reduce((total, cell) => total + Math.max(1, cell.fr), 0);
  const widths: number[] = [];
  let used = 0;
  for (let i = 0; i < count; i += 1) {
    if (i === count - 1) widths.push(Math.max(1, useful - used));
    else {
      const size = Math.max(1, Math.floor((useful * Math.max(1, row.cells[i].fr)) / sum));
      widths.push(size);
      used += size;
    }
  }
  return widths;
}

function outerHeight(style: Style, lines: number): number {
  const leading = Math.max(1, style.lineHeight);
  const text = lines + Math.max(0, lines - 1) * (leading - 1);
  return style.borderWidth * 2 + style.pad[0] + style.pad[2] + text;
}

function wrap(text: string, width: number): string[] {
  if (width < 1) return [];
  const out: string[] = [];
  for (const para of text.split("\n")) {
    const chars = [...para];
    if (chars.length === 0) {
      out.push("");
      continue;
    }
    for (let i = 0; i < chars.length; i += width) out.push(chars.slice(i, i + width).join(""));
  }
  return out.length ? out : [""];
}

export function layoutWindow(
  win: Win,
  live: Map<string, LiveZone>,
  columnsCount: number,
  focusId: string | null,
  navSel: Record<string, number>,
  status: string,
): { width: number; height: number; boxes: Box[] } {
  const width = Math.max(24, Math.min(win.style.maxWidth || 260, columnsCount));
  const frame = win.style;
  const x = frame.mar[3] + frame.borderWidth + frame.pad[3];
  const y = frame.mar[0] + frame.borderWidth + frame.pad[0];
  const inner = Math.max(8, width - x - frame.mar[1] - frame.borderWidth - frame.pad[1]);
  const laid = layoutGrid(win.grid, x, y, inner, 0, live, focusId, navSel, status);
  const height = y + laid.height + frame.pad[2] + frame.borderWidth + frame.mar[2];
  return { width, height, boxes: laid.boxes };
}

function layoutGrid(
  grid: Grid,
  x: number,
  y: number,
  width: number,
  extra: number,
  live: Map<string, LiveZone>,
  focusId: string | null,
  navSel: Record<string, number>,
  status: string,
): { boxes: Box[]; height: number } {
  const boxes: Box[] = [];
  const heights = grid.rows.map((row) => rowHeight(row, width, live, focusId, navSel, status));
  if (extra > 0 && heights.length) heights[heights.length - 1] += extra;
  let cursor = y;
  let previous = 0;
  for (let i = 0; i < grid.rows.length; i += 1) {
    const top = edge(grid.rows[i], 0);
    cursor += i === 0 ? top : Math.max(top, previous);
    boxes.push(...placeRow(grid.rows[i], x, cursor, width, heights[i], live, focusId, navSel, status));
    cursor += heights[i];
    previous = edge(grid.rows[i], 2);
  }
  return { boxes, height: cursor - y + previous };
}

function rowHeight(
  row: GridRow,
  width: number,
  live: Map<string, LiveZone>,
  focusId: string | null,
  navSel: Record<string, number>,
  status: string,
): number {
  const cols = columns(row, width);
  let height = 1;
  row.cells.forEach((cell, index) => {
    if (cell.zone) height = Math.max(height, outerHeight(cell.zone.base, row.lines));
    if (cell.grid) {
      const nested = layoutGrid(cell.grid, 0, 0, cols[index] ?? width, 0, live, focusId, navSel, status);
      height = Math.max(height, nested.height);
    }
  });
  return height;
}

function placeRow(
  row: GridRow,
  x: number,
  y: number,
  width: number,
  height: number,
  live: Map<string, LiveZone>,
  focusId: string | null,
  navSel: Record<string, number>,
  status: string,
): Box[] {
  const boxes: Box[] = [];
  const cols = columns(row, width);
  let cursor = x + flank(row, 3);
  row.cells.forEach((cell, index) => {
    const col = cols[index] ?? 1;
    if (cell.zone) {
      const focused = cell.zone.id === focusId;
      const style = focused ? cell.zone.focusOn : cell.zone.base;
      const textW = Math.max(1, col - style.borderWidth * 2 - style.pad[1] - style.pad[3]);
      const painted = paintZone(cell.zone.tag, cell.zone.id, live.get(cell.zone.id), textW, navSel[cell.zone.id] ?? 0, status);
      boxes.push({
        id: cell.zone.id,
        x: cursor,
        y,
        w: col,
        h: height,
        textW,
        textH: row.lines,
        style,
        tag: cell.zone.tag,
        title: painted.title,
        lines: painted.lines,
        itemIds: painted.itemIds,
        selected: painted.selected,
      });
    }
    if (cell.grid) {
      const measured = layoutGrid(cell.grid, cursor, y, col, 0, live, focusId, navSel, status);
      const extra = Math.max(0, height - measured.height);
      const nested =
        extra > 0 ? layoutGrid(cell.grid, cursor, y, col, extra, live, focusId, navSel, status) : measured;
      boxes.push(...nested.boxes);
    }
    cursor += col + (index < row.cells.length - 1 ? gap(row, index) : 0);
  });
  return boxes;
}

function paintZone(
  tag: ZoneTag,
  id: string,
  live: LiveZone | undefined,
  textW: number,
  selected: number,
  status: string,
): { lines: string[]; itemIds: string[]; title: boolean; selected: number } {
  const label = live?.label ?? "";
  const value = live?.value ?? "";
  if (tag === "footer") return { lines: wrap(status, textW), itemIds: [], title: false, selected: -1 };
  if (tag === "header" || tag === "button") return { lines: wrap(label, textW), itemIds: [], title: tag === "header", selected: -1 };
  if (tag === "input") {
    const name = live?.name || "";
    const body = wrap(value, textW);
    return { lines: [ [...name].slice(0, textW).join(""), ...body ], itemIds: [], title: true, selected: -1 };
  }
  if (tag === "textarea") return { lines: wrap(value, textW), itemIds: [], title: false, selected: -1 };
  const items = live?.items ?? [];
  const lines = [[...label].slice(0, textW).join("")];
  const itemIds = [""];
  const sel = items.length ? Math.min(Math.max(0, selected), items.length - 1) : -1;
  items.forEach((item, index) => {
    lines.push([...item.label].slice(0, textW).join(""));
    itemIds.push(item.id);
    if (index === sel) {
      /* selection is recorded, the view highlights this row */
    }
  });
  return { lines, itemIds, title: true, selected: sel };
}

export const textOf = paintZone;
