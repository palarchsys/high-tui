import { create } from "zustand";
import { persist } from "zustand/middleware";
import { itemsFromSource, projectPath, writeItems } from "./data";
import { exampleSources, exampleWindows } from "./fixtures";
import { applyCss, emitCss, emitHtml, ensureReturn, parseHtml } from "./html";
import {
  findZone,
  freshId,
  gridOfZone,
  makeZone,
  mapZones,
  replaceZone,
  updateGrid,
  zonesOf,
} from "./model";
import type { Action, Cell, GridRow, Item, Sel, Source, Style, Win, Zone, ZoneTag } from "./types";
import { validate } from "./validate";

type Snap = { windows: Win[]; sources: Source[]; activeId: string };

type Drag =
  | { kind: "zone"; id: string }
  | { kind: "tag"; tag: ZoneTag }
  | null;

type State = {
  windows: Win[];
  sources: Source[];
  activeId: string;
  gridId: string;
  sel: Sel | null;
  tab: "display" | "data" | "preview";
  focusId: string | null;
  navSel: Record<string, number>;
  columns: number;
  past: Snap[];
  future: Snap[];
  drag: Drag;
  hydrated: boolean;
  reset: () => void;
  undo: () => void;
  redo: () => void;
  setTab: (tab: State["tab"]) => void;
  setColumns: (columns: number) => void;
  selectWindow: (id: string) => void;
  select: (sel: Sel | null) => void;
  setDrag: (drag: Drag) => void;
  focusZone: (id: string) => void;
  pickItem: (zoneId: string, index: number) => void;
  patchWindow: (partial: Partial<Pick<Win, "label" | "file">> & { style?: Style }) => void;
  patchZone: (id: string, partial: Partial<Zone>) => void;
  patchStyle: (id: string, key: keyof Style, value: Style[keyof Style], focus: boolean) => void;
  setRowLines: (rowId: string, lines: number) => void;
  setCellFr: (zoneId: string, fr: number) => void;
  addRow: (tag: ZoneTag) => void;
  addColumn: (rowId: string, tag: ZoneTag) => void;
  swapZones: (a: string, b: string) => void;
  removeZone: (id: string) => void;
  nestZone: (id: string) => void;
  addWindow: () => void;
  removeWindow: (id: string) => void;
  importText: (html: string, css: string) => string | null;
  setSource: (id: string, partial: Partial<Source>) => void;
  addSource: () => void;
  removeSource: (id: string) => void;
  bind: (zoneId: string, mode: "raw" | "var" | "list" | "none", sourceId: string, varName?: string) => void;
  setReference: (zoneId: string, sourceId: string, varName: string, slot: "label" | "value") => void;
  updateItems: (zoneId: string, items: Item[]) => void;
};

function snap(state: Pick<State, "windows" | "sources" | "activeId">): Snap {
  return {
    windows: structuredClone(state.windows),
    sources: structuredClone(state.sources),
    activeId: state.activeId,
  };
}

function active(state: Pick<State, "windows" | "activeId">): Win {
  return state.windows.find((win) => win.id === state.activeId) ?? state.windows[0];
}

function commit(state: State, next: Partial<State> & Pick<State, "windows" | "sources" | "activeId">): Partial<State> {
  return {
    ...next,
    past: [...state.past, snap(state)].slice(-40),
    future: [],
  };
}

function mapActive(state: State, fn: (win: Win) => Win): Pick<State, "windows" | "sources" | "activeId"> {
  return {
    windows: state.windows.map((win) => (win.id === state.activeId ? fn(win) : win)),
    sources: state.sources,
    activeId: state.activeId,
  };
}

function firstFocus(win: Win): string | null {
  return zonesOf(win).find((zone) => zone.tag !== "header" && zone.tag !== "footer")?.id ?? null;
}

export const useEditor = create<State>()(
  persist(
    (set, get) => ({
      windows: exampleWindows(),
      sources: exampleSources(),
      activeId: "main",
      gridId: "corps",
      sel: { kind: "zone", id: "nav-rubriques" },
      tab: "display",
      focusId: "nav-rubriques",
      navSel: {},
      columns: 88,
      past: [],
      future: [],
      drag: null,
      hydrated: false,
      reset: () =>
        set({
          windows: exampleWindows(),
          sources: exampleSources(),
          activeId: "main",
          gridId: "corps",
          sel: { kind: "zone", id: "nav-rubriques" },
          focusId: "nav-rubriques",
          navSel: {},
          past: [],
          future: [],
        }),
      undo: () => {
        const { past, future } = get();
        const prev = past[past.length - 1];
        if (!prev) return;
        set({
          ...prev,
          past: past.slice(0, -1),
          future: [snap(get()), ...future].slice(0, 40),
        });
      },
      redo: () => {
        const { future, past } = get();
        const next = future[0];
        if (!next) return;
        set({
          ...next,
          future: future.slice(1),
          past: [...past, snap(get())].slice(-40),
        });
      },
      setTab: (tab) => set({ tab }),
      setColumns: (columns) => set({ columns: Math.max(40, Math.min(160, columns)) }),
      selectWindow: (id) => {
        const win = get().windows.find((item) => item.id === id);
        if (!win) return;
        set({ activeId: id, gridId: win.grid.id, sel: null, focusId: firstFocus(win) });
      },
      select: (sel) => {
        if (sel?.kind === "zone") set({ sel, gridId: gridOfZone(active(get()), sel.id) });
        else if (sel?.kind === "grid") set({ sel, gridId: sel.id });
        else set({ sel });
      },
      setDrag: (drag) => set({ drag }),
      focusZone: (id) => set({ focusId: id, sel: { kind: "zone", id }, gridId: gridOfZone(active(get()), id) }),
      pickItem: (zoneId, index) => {
        const win = active(get());
        const zone = findZone(win, zoneId);
        const navSel = { ...get().navSel, [zoneId]: index };
        if (zone?.nextNav) {
          let cur: string | undefined = zone.nextNav;
          const seen = new Set<string>();
          while (cur && !seen.has(cur)) {
            seen.add(cur);
            navSel[cur] = 0;
            cur = findZone(win, cur)?.nextNav;
          }
        }
        set({ navSel, focusId: zoneId, sel: { kind: "zone", id: zoneId } });
      },
      patchWindow: (partial) => {
        const state = get();
        set(
          commit(state, {
            ...mapActive(state, (win) => ({
              ...win,
              ...partial,
              style: partial.style ?? win.style,
              id: partial.file ? partial.file.replace(/\.html$/, "") : win.id,
            })),
            activeId: partial.file ? partial.file.replace(/\.html$/, "") : state.activeId,
          }),
        );
      },
      patchZone: (id, partial) => {
        const state = get();
        set(
          commit(
            state,
            mapActive(state, (win) =>
              mapZones(win, (zone) => {
                if (zone.id !== id) {
                  let next = zone;
                  if (partial.id && zone.nextNav === id) next = { ...next, nextNav: partial.id };
                  if (partial.id) {
                    const retarget = (action: Action): Action => ({
                      ...action,
                      target: action.target === id ? partial.id! : action.target,
                      output: action.output === id ? partial.id! : action.output,
                    });
                    next = { ...next, focus: retarget(next.focus), enter: retarget(next.enter) };
                  }
                  return next;
                }
                return { ...zone, ...partial };
              }),
            ),
          ),
        );
      },
      patchStyle: (id, key, value, focus) => {
        const state = get();
        set(
          commit(
            state,
            mapActive(state, (win) =>
              mapZones(win, (zone) => {
                if (zone.id !== id) return zone;
                const style = { ...(focus ? zone.focusOn : zone.base), [key]: value };
                return focus ? { ...zone, focusOn: style } : { ...zone, base: style };
              }),
            ),
          ),
        );
      },
      setRowLines: (rowId, lines) => {
        const state = get();
        const apply = (rows: GridRow[]): GridRow[] =>
          rows.map((row) => ({
            ...row,
            lines: row.id === rowId ? Math.max(1, lines) : row.lines,
            cells: row.cells.map((cell) => (cell.grid ? { ...cell, grid: { ...cell.grid, rows: apply(cell.grid.rows) } } : cell)),
          }));
        set(commit(state, mapActive(state, (win) => ({ ...win, grid: { ...win.grid, rows: apply(win.grid.rows) } }))));
      },
      setCellFr: (zoneId, fr) => {
        const state = get();
        const apply = (rows: GridRow[]): GridRow[] =>
          rows.map((row) => ({
            ...row,
            cells: row.cells.map((cell) => {
              if (cell.zone?.id === zoneId) return { ...cell, fr: Math.max(1, fr) };
              if (cell.grid) return { ...cell, grid: { ...cell.grid, rows: apply(cell.grid.rows) } };
              return cell;
            }),
          }));
        set(commit(state, mapActive(state, (win) => ({ ...win, grid: { ...win.grid, rows: apply(win.grid.rows) } }))));
      },
      addRow: (tag) => {
        const state = get();
        set(
          commit(
            state,
            mapActive(state, (win) => {
              const id = freshId(win, tag);
              const row: GridRow = { id: freshId(win, "rang"), lines: tag === "textarea" ? 3 : 1, cells: [{ fr: 1, zone: makeZone(tag, id) }] };
              return updateGrid(win, state.gridId || win.grid.id, (grid) => ({ ...grid, rows: [...grid.rows, row] }));
            }),
          ),
        );
      },
      addColumn: (rowId, tag) => {
        const state = get();
        set(
          commit(
            state,
            mapActive(state, (win) => {
              const id = freshId(win, tag);
              const zone = makeZone(tag, id);
              const apply = (rows: GridRow[]): GridRow[] =>
                rows.map((row) => ({
                  ...row,
                  cells:
                    row.id === rowId
                      ? [...row.cells, { fr: 1, zone }]
                      : row.cells.map((cell) => (cell.grid ? { ...cell, grid: { ...cell.grid, rows: apply(cell.grid.rows) } } : cell)),
                }));
              return { ...win, grid: { ...win.grid, rows: apply(win.grid.rows) } };
            }),
          ),
        );
      },
      swapZones: (a, b) => {
        if (a === b) return;
        const state = get();
        set(
          commit(
            state,
            mapActive(state, (win) => {
              const left = findZone(win, a);
              const right = findZone(win, b);
              if (!left || !right) return win;
              return mapZones(win, (zone) => {
                if (zone.id === a) return right;
                if (zone.id === b) return left;
                return zone;
              });
            }),
          ),
        );
      },
      removeZone: (id) => {
        const state = get();
        set(
          commit(
            state,
            mapActive(state, (win) => {
              const drop = (rows: GridRow[]): GridRow[] =>
                rows
                  .map((row) => ({
                    ...row,
                    cells: row.cells
                      .filter((cell) => cell.zone?.id !== id)
                      .map((cell) => (cell.grid ? { ...cell, grid: { ...cell.grid, rows: drop(cell.grid.rows) } } : cell)),
                  }))
                  .filter((row) => row.cells.length > 0);
              const rows = drop(win.grid.rows);
              if (!rows.length) return win;
              return { ...win, grid: { ...win.grid, rows } };
            }),
          ),
        );
      },
      nestZone: (id) => {
        const state = get();
        set(
          commit(
            state,
            mapActive(state, (win) => {
              const zone = findZone(win, id);
              if (!zone) return win;
              const gridId = freshId(win, "grille");
              const apply = (rows: GridRow[]): GridRow[] =>
                rows.map((row) => ({
                  ...row,
                  cells: row.cells.map((cell) => {
                    if (cell.zone?.id === id) {
                      const nested: Cell = {
                        fr: cell.fr,
                        grid: { id: gridId, rows: [{ id: freshId(win, "rang"), lines: row.lines, cells: [{ fr: 1, zone }] }] },
                      };
                      return nested;
                    }
                    if (cell.grid) return { ...cell, grid: { ...cell.grid, rows: apply(cell.grid.rows) } };
                    return cell;
                  }),
                }));
              return { ...win, grid: { ...win.grid, rows: apply(win.grid.rows) } };
            }),
          ),
        );
      },
      addWindow: () => {
        const state = get();
        const draft = exampleWindows()[0];
        const stampId = `page-${state.windows.length + 1}`;
        const win: Win = {
          id: stampId,
          label: "Nouvelle fenêtre",
          file: `${stampId}.html`,
          style: draft.style,
          grid: {
            id: stampId,
            rows: [
              { id: `${stampId}-entete`, lines: 1, cells: [{ fr: 1, zone: makeZone("header", `${stampId}-titre`, { label: "Nouvelle fenêtre" }) }] },
              {
                id: `${stampId}-liste`,
                lines: 3,
                cells: [{ fr: 1, zone: makeZone("nav", `${stampId}-nav`, { label: "Liste", items: [emptyReturn(`${stampId}-nav`)] }) }],
              },
              { id: `${stampId}-barre`, lines: 1, cells: [{ fr: 1, zone: makeZone("footer", `${stampId}-statut`) }] },
            ],
          },
        };
        ensureReturn(win);
        set(commit(state, { windows: [...state.windows, win], sources: state.sources, activeId: win.id }));
        set({ gridId: win.id, sel: { kind: "zone", id: `${stampId}-nav` }, focusId: `${stampId}-nav` });
      },
      removeWindow: (id) => {
        const state = get();
        if (state.windows.length < 2) return;
        const windows = state.windows.filter((win) => win.id !== id);
        const activeId = state.activeId === id ? windows[0].id : state.activeId;
        set(commit(state, { windows, sources: state.sources, activeId }));
      },
      importText: (html, css) => {
        try {
          const parsed = applyCss(parseHtml(html, "import.html"), css);
          const state = get();
          const file = parsed.file === "import.html" ? `${parsed.id}.html` : parsed.file;
          const win = { ...parsed, file };
          const windows = state.windows.some((item) => item.id === win.id)
            ? state.windows.map((item) => (item.id === win.id ? win : item))
            : [...state.windows, win];
          set(commit(state, { windows, sources: state.sources, activeId: win.id }));
          set({ gridId: win.grid.id, focusId: firstFocus(win), sel: null });
          return null;
        } catch (error) {
          return error instanceof Error ? error.message : "erreur : import refusé";
        }
      },
      setSource: (id, partial) => {
        const state = get();
        const sources = state.sources.map((source) => (source.id === id ? { ...source, ...partial } : source));
        const current = sources.find((source) => source.id === id);
        let next = sources;
        if (current) {
          next = sources.map((source) => {
            if (source.projectedFrom !== current.path) return source;
            const projected = itemsFromSource(current, current.path);
            if (projected.error) return source;
            return { ...source, body: writeItems(projected.items) };
          });
        }
        set(commit(state, { windows: state.windows, sources: next, activeId: state.activeId }));
      },
      addSource: () => {
        const state = get();
        const id = `source-${state.sources.length + 1}`;
        const source: Source = { id, path: `samples/${id}.json`, body: '[{"label":"Ligne"}]\n', read: "json", separator: "|" };
        set(commit(state, { windows: state.windows, sources: [...state.sources, source], activeId: state.activeId }));
      },
      removeSource: (id) => {
        const state = get();
        set(
          commit(state, {
            windows: state.windows,
            sources: state.sources.filter((source) => source.id !== id),
            activeId: state.activeId,
          }),
        );
      },
      bind: (zoneId, mode, sourceId, varName) => {
        const state = get();
        const source = state.sources.find((item) => item.id === sourceId);
        if (!source && mode !== "none") return;
        let sources = state.sources;
        let path = source?.path ?? "";
        if (source && mode === "list" && source.read !== "json") {
          path = projectPath(source.path);
          const projected = itemsFromSource(source, source.path);
          if (projected.error) return;
          const body = writeItems(projected.items);
          const existing = sources.find((item) => item.path === path);
          sources = existing
            ? sources.map((item) => (item.path === path ? { ...item, body, read: "json", projectedFrom: source.path } : item))
            : [...sources, { id: `proj-${source.id}`, path, body, read: "json", separator: "", projectedFrom: source.path }];
        }
        set(
          commit(state, {
            sources,
            activeId: state.activeId,
            windows: state.windows.map((win) =>
              win.id !== state.activeId
                ? win
                : mapZones(win, (zone) => {
                    if (zone.id !== zoneId) return zone;
                    if (mode === "none") return { ...zone, file: "", name: "", itemsFrom: zone.tag === "nav" ? "" : zone.itemsFrom };
                    if (mode === "list" && zone.tag === "nav") return { ...zone, itemsFrom: path };
                    if (mode === "raw" && (zone.tag === "input" || zone.tag === "textarea")) {
                      return { ...zone, file: path, name: "" };
                    }
                    if (mode === "var" && (zone.tag === "input" || zone.tag === "textarea")) {
                      return { ...zone, file: path, name: varName || zone.name };
                    }
                    return zone;
                  }),
            ),
          }),
        );
      },
      setReference: (zoneId, sourceId, varName, slot) => {
        const state = get();
        const source = state.sources.find((item) => item.id === sourceId);
        if (!source || !varName) return;
        const token = `${source.path}::${varName}::${slot === "label" ? "label" : "value"}`;
        set(
          commit(
            state,
            mapActive(state, (win) =>
              mapZones(win, (zone) => {
                if (zone.id !== zoneId) return zone;
                if (slot === "label") return { ...zone, label: token };
                return { ...zone, value: token };
              }),
            ),
          ),
        );
      },
      updateItems: (zoneId, items) => {
        const state = get();
        set(commit(state, mapActive(state, (win) => replaceZone(win, zoneId, { ...findZone(win, zoneId)!, items }))));
      },
    }),
    {
      name: "atelier-high-tui",
      version: 1,
      skipHydration: true,
      partialize: (state) => ({
        windows: state.windows,
        sources: state.sources,
        activeId: state.activeId,
        gridId: state.gridId,
        columns: state.columns,
      }),
    },
  ),
);

function emptyReturn(navId: string): Item {
  return { ...makeZone("nav", navId).items[0], id: `${navId}-ligne`, label: "Ligne" };
}

export function currentWin(state: Pick<State, "windows" | "activeId">): Win {
  return active(state);
}

export function issuesOf(state: Pick<State, "windows" | "sources" | "activeId">): string[] {
  const win = active(state);
  return win ? validate(win, state.sources) : [];
}

export function exported(win: Win): { html: string; css: string } {
  return { html: emitHtml(win), css: emitCss(win) };
}
