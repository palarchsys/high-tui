import { useState } from "react";
import { TAGS, tagLabel } from "@/components/bits";
import { Inspector } from "@/components/inspector";
import type { Cell, Grid, GridRow, ZoneTag } from "@/editor/types";
import { currentWin, useEditor } from "@/editor/store";

export function DisplayPane() {
  const win = useEditor((state) => currentWin(state));
  const drag = useEditor((state) => state.drag);
  const setDrag = useEditor((state) => state.setDrag);
  const addRow = useEditor((state) => state.addRow);
  const select = useEditor((state) => state.select);
  const armed = drag?.kind === "tag" ? drag.tag : null;

  return (
    <section className="flex min-h-0 min-w-0 flex-1 flex-col lg:flex-row">
      <div className="flex min-h-0 min-w-0 flex-1 flex-col">
        <div className="flex gap-2 overflow-x-auto border-b border-line p-2">
          {TAGS.map((item) => (
            <button
              key={item.tag}
              type="button"
              className={`min-h-11 shrink-0 rounded-sm border px-3 text-sm ${armed === item.tag ? "border-accent bg-accent text-bg" : "border-line bg-surface text-fg"}`}
              onPointerDown={(event) => {
                event.preventDefault();
                setDrag({ kind: "tag", tag: item.tag });
              }}
            >
              {item.label}
            </button>
          ))}
        </div>
        <div className="min-h-0 flex-1 overflow-auto p-3">
          <GridView grid={win.grid} />
          <button
            type="button"
            data-new-row="true"
            className="mt-3 flex min-h-11 w-full items-center justify-center rounded-sm border border-dashed border-line text-sm text-muted"
            onClick={() => armed && addRow(armed)}
          >
            {armed ? `Déposer ${tagLabel(armed)} dans une rangée` : "Arme une zone, puis dépose-la ici"}
          </button>
        </div>
      </div>
      <Inspector />
    </section>
  );
}

function GridView({ grid }: { grid: Grid }) {
  const gridId = useEditor((state) => state.gridId);
  const select = useEditor((state) => state.select);
  return (
    <div className={`flex flex-col gap-2 rounded-sm border p-2 ${gridId === grid.id ? "border-accent" : "border-line"}`}>
      <button type="button" className="self-start text-xs text-muted" onClick={() => select({ kind: "grid", id: grid.id })}>
        {grid.id}
      </button>
      {grid.rows.map((row, index) => (
        <RowView key={row.id || index} row={row} />
      ))}
    </div>
  );
}

function RowView({ row }: { row: GridRow }) {
  const addColumn = useEditor((state) => state.addColumn);
  const setRowLines = useEditor((state) => state.setRowLines);
  const drag = useEditor((state) => state.drag);
  const [lines, setLines] = useState(row.lines);

  return (
    <div className="flex flex-col gap-1">
      <div className="flex items-center gap-2 text-xs text-muted">
        <span>{row.id || "rangée"}</span>
        <span className="tabular-nums">{row.lines} lignes</span>
        <button
          type="button"
          className="min-h-11 px-2"
          onPointerDown={(event) => {
            event.preventDefault();
            const startY = event.clientY;
            const start = row.lines;
            let next = start;
            const move = (ev: PointerEvent) => {
              next = Math.max(1, start + Math.round((ev.clientY - startY) / 16));
              setLines(next);
            };
            const up = () => {
              setRowLines(row.id, next);
              window.removeEventListener("pointermove", move);
              window.removeEventListener("pointerup", up);
            };
            window.addEventListener("pointermove", move);
            window.addEventListener("pointerup", up);
          }}
        >
          hauteur
        </button>
      </div>
      <div className="flex items-stretch gap-1">
        {row.cells.map((cell, index) => (
          <CellView key={cell.zone?.id || cell.grid?.id || index} cell={cell} />
        ))}
        <button
          type="button"
          data-row={row.id}
          className="flex min-h-11 w-11 shrink-0 items-center justify-center rounded-sm border border-dashed border-line text-muted"
          onClick={() => drag?.kind === "tag" && row.id && addColumn(row.id, drag.tag)}
        >
          +
        </button>
      </div>
    </div>
  );
}

function CellView({ cell }: { cell: Cell }) {
  const setCellFr = useEditor((state) => state.setCellFr);
  if (cell.grid) {
    return (
      <div className="min-w-0 flex-1" style={{ flexGrow: cell.fr }}>
        <GridView grid={cell.grid} />
      </div>
    );
  }
  if (!cell.zone) return null;
  const zone = cell.zone;
  return (
    <div className="flex min-w-0 flex-1 items-stretch" style={{ flexGrow: cell.fr }}>
      <ZoneCard id={zone.id} tag={zone.tag} label={zone.label || zone.value || zone.id} hint={hint(zone)} fr={cell.fr} />
      <button
        type="button"
        aria-label="Largeur"
        className="w-3 cursor-col-resize text-muted"
        onPointerDown={(event) => {
          event.preventDefault();
          const startX = event.clientX;
          const start = cell.fr;
          let next = start;
          const move = (ev: PointerEvent) => {
            next = Math.max(1, start + Math.round((ev.clientX - startX) / 28));
          };
          const up = () => {
            setCellFr(zone.id, next);
            window.removeEventListener("pointermove", move);
            window.removeEventListener("pointerup", up);
          };
          window.addEventListener("pointermove", move);
          window.addEventListener("pointerup", up);
        }}
      />
    </div>
  );
}

function hint(zone: { itemsFrom: string; nextNav: string; file: string; name: string }): string {
  if (zone.itemsFrom) return zone.itemsFrom;
  if (zone.nextNav) return `next-data → ${zone.nextNav}`;
  if (zone.file && zone.name) return `${zone.name} · ${zone.file}`;
  if (zone.file) return zone.file;
  return "";
}

function ZoneCard({ id, tag, label, hint: note, fr }: { id: string; tag: ZoneTag; label: string; hint: string; fr: number }) {
  const sel = useEditor((state) => state.sel);
  const select = useEditor((state) => state.select);
  const setDrag = useEditor((state) => state.setDrag);
  const drag = useEditor((state) => state.drag);
  const swapZones = useEditor((state) => state.swapZones);
  const selected = sel?.kind === "zone" && sel.id === id;
  const hot = drag?.kind === "zone" && drag.id !== id;
  return (
    <article
      data-drop={id}
      className={`min-w-0 flex-1 rounded-sm border p-2 ${selected ? "border-accent" : "border-line"} ${hot ? "bg-bg" : "bg-surface"}`}
      onPointerDown={(event) => {
        select({ kind: "zone", id });
        if ((event.target as HTMLElement).closest("[data-grip]")) {
          event.preventDefault();
          setDrag({ kind: "zone", id });
        }
      }}
      onPointerUp={() => {
        if (drag?.kind === "zone" && drag.id !== id) swapZones(drag.id, id);
      }}
    >
      <header className="flex items-center gap-2 text-xs text-muted">
        <button type="button" data-grip="true" className="min-h-11 cursor-grab px-1" aria-label="Déplacer">
          ::
        </button>
        <span>{tagLabel(tag)}</span>
        <span className="ml-auto tabular-nums">{fr}fr</span>
      </header>
      <p className="truncate text-sm text-fg">{label}</p>
      <p className="truncate text-xs text-accent">{id}</p>
      {note ? <p className="truncate text-xs text-muted">{note}</p> : null}
    </article>
  );
}
