import { layoutWindow } from "@/editor/layout";
import { zonesOf } from "@/editor/model";
import { simulate } from "@/editor/simulate";
import { currentWin, useEditor } from "@/editor/store";

export function PreviewPane() {
  const win = useEditor((state) => currentWin(state));
  const sources = useEditor((state) => state.sources);
  const focusId = useEditor((state) => state.focusId);
  const navSel = useEditor((state) => state.navSel);
  const columns = useEditor((state) => state.columns);
  const setColumns = useEditor((state) => state.setColumns);
  const focusZone = useEditor((state) => state.focusZone);
  const pickItem = useEditor((state) => state.pickItem);
  const focusable = zonesOf(win).filter((zone) => zone.tag !== "header" && zone.tag !== "footer");
  const focus = focusId && focusable.some((zone) => zone.id === focusId) ? focusId : focusable[0]?.id ?? null;
  const sim = simulate(win, sources, navSel, focus);
  const laid = layoutWindow(win, sim.live, columns, focus, navSel, sim.status);
  const index = Math.max(0, focusable.findIndex((zone) => zone.id === focus));

  return (
    <section className="flex h-[46vh] min-h-48 min-w-0 max-w-full shrink-0 flex-col overflow-hidden border-t border-line lg:h-[38vh]">
      <div className="flex flex-wrap items-center gap-2 border-b border-line px-3 py-1">
        <h2 className="text-sm text-accent">Aperçu</h2>
        <span className="text-sm text-fg">{sim.status}</span>
        <label className="ml-auto flex items-center gap-2 text-xs text-muted">
          Colonnes
          <input
            className="w-24 accent-accent"
            type="range"
            min={40}
            max={140}
            value={columns}
            onChange={(event) => setColumns(Number(event.target.value))}
          />
          <span className="tabular-nums">{columns}</span>
        </label>
        <button type="button" className="min-h-11 px-2 text-sm" onClick={() => focusable[index - 1] && focusZone(focusable[index - 1].id)}>
          Zone précédente
        </button>
        <button type="button" className="min-h-11 px-2 text-sm" onClick={() => focusable[index + 1] && focusZone(focusable[(index + 1) % focusable.length].id)}>
          Zone suivante
        </button>
      </div>
      <div className="min-h-0 flex-1 overflow-auto bg-bg p-3">
        <div
          className="relative font-mono text-xs leading-none"
          style={{
            width: `${laid.width}ch`,
            height: `${laid.height}em`,
            background: win.style.background,
            color: win.style.color,
          }}
        >
          {laid.boxes.map((box) => (
            <div
              key={box.id}
              role="button"
              tabIndex={0}
              onClick={() => focusZone(box.id)}
              onKeyDown={(event) => {
                if (event.key === "Enter") focusZone(box.id);
              }}
              className="absolute overflow-hidden text-left"
              style={{
                left: `${box.x}ch`,
                top: `${box.y}em`,
                width: `${box.w}ch`,
                height: `${box.h}em`,
                boxSizing: "border-box",
                background: box.style.background,
                color: box.style.color,
                borderStyle: "solid",
                borderWidth: `${box.style.borderWidth}ch`,
                borderColor: box.style.borderColor,
                borderRadius: `${box.style.radius}ch`,
                padding: `${box.style.pad[0]}em ${box.style.pad[1]}ch ${box.style.pad[2]}em ${box.style.pad[3]}ch`,
                fontWeight: box.style.bold ? 600 : 400,
              }}
            >
              {box.lines.slice(0, Math.max(box.textH, 1)).map((line, lineIndex) => {
                const item = lineIndex > 0 && box.tag === "nav";
                const selected = item && lineIndex - 1 === box.selected;
                const title = lineIndex === 0 && box.title;
                return (
                  <div
                    key={`${box.id}-${lineIndex}`}
                    className="whitespace-pre"
                    style={{
                      height: `${Math.max(1, box.style.lineHeight)}em`,
                      color: title ? box.style.titleColor : undefined,
                      background: selected ? box.style.titleColor : undefined,
                    }}
                    onClick={(event) => {
                      if (!item) return;
                      event.stopPropagation();
                      pickItem(box.id, lineIndex - 1);
                    }}
                  >
                    {line || " "}
                  </div>
                );
              })}
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
