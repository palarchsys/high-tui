import { Button, control, Field, tagLabel } from "@/components/bits";
import { findZone, rowOfZone } from "@/editor/model";
import { currentWin, useEditor } from "@/editor/store";
import type { Action, ActionKind, Zone } from "@/editor/types";
import { ACTION_KINDS } from "@/editor/types";

export function Inspector() {
  const win = useEditor((state) => currentWin(state));
  const sel = useEditor((state) => state.sel);
  const patchWindow = useEditor((state) => state.patchWindow);
  const patchZone = useEditor((state) => state.patchZone);
  const patchStyle = useEditor((state) => state.patchStyle);
  const removeZone = useEditor((state) => state.removeZone);
  const nestZone = useEditor((state) => state.nestZone);
  const setRowLines = useEditor((state) => state.setRowLines);
  const setCellFr = useEditor((state) => state.setCellFr);
  const zone = sel?.kind === "zone" ? findZone(win, sel.id) : undefined;
  const row = zone ? rowOfZone(win, zone.id) : undefined;

  return (
    <aside className="flex w-full shrink-0 flex-col gap-4 overflow-auto border-line p-3 lg:w-80 lg:border-l">
      <Field label="Fenêtre">
        <input className={control} value={win.label} onChange={(event) => patchWindow({ label: event.target.value })} />
      </Field>
      <Field label="Fichier">
        <input className={control} value={win.file} onChange={(event) => patchWindow({ file: event.target.value })} />
      </Field>
      <Field label="Largeur max (ch)">
        <input
          className={control}
          type="number"
          min={24}
          value={win.style.maxWidth}
          onChange={(event) => patchWindow({ style: { ...win.style, maxWidth: Number(event.target.value) || 24 } })}
        />
      </Field>
      <Colors
        style={win.style}
        onChange={(key, value) => patchWindow({ style: { ...win.style, [key]: value } })}
      />
      {zone && row ? (
        <>
          <h2 className="text-sm text-accent">{tagLabel(zone.tag)}</h2>
          <Field label="Identifiant">
            <input className={control} value={zone.id} onChange={(event) => patchZone(zone.id, { id: event.target.value.trim() })} />
          </Field>
          <Field label={zone.tag === "input" || zone.tag === "textarea" ? "Placeholder" : "Libellé"}>
            <input
              className={control}
              value={zone.tag === "input" || zone.tag === "textarea" ? zone.value : zone.label}
              onChange={(event) =>
                patchZone(zone.id, zone.tag === "input" || zone.tag === "textarea" ? { value: event.target.value } : { label: event.target.value })
              }
            />
          </Field>
          <div className="grid grid-cols-2 gap-2">
            <Field label="Lignes">
              <input className={control} type="number" min={1} value={row.lines} onChange={(event) => row.id && setRowLines(row.id, Number(event.target.value) || 1)} />
            </Field>
            <Field label="fr">
              <input className={control} type="number" min={1} value={frOf(win.grid, zone.id) || 1} onChange={(event) => setCellFr(zone.id, Number(event.target.value) || 1)} />
            </Field>
          </div>
          {(zone.tag === "input" || zone.tag === "textarea") && (
            <label className="flex min-h-11 items-center gap-2 text-sm">
              <input type="checkbox" checked={zone.readonly} onChange={(event) => patchZone(zone.id, { readonly: event.target.checked })} />
              Lecture seule
            </label>
          )}
          {(zone.tag === "input" || zone.tag === "textarea") && (
            <Field label="Enregistrer">
              <select className={control} value={zone.save} onChange={(event) => patchZone(zone.id, { save: event.target.value as Zone["save"] })}>
                <option value="">non</option>
                <option value="enter">entrée</option>
                <option value="button">bouton</option>
              </select>
            </Field>
          )}
          <ActionEditor title="Focus" action={zone.focus} loose={zone.tag === "nav"} onChange={(focus) => patchZone(zone.id, { focus })} />
          <ActionEditor title="Entrée" action={zone.enter} loose={zone.tag === "nav"} onChange={(enter) => patchZone(zone.id, { enter })} />
          <ZoneStyle zone={zone} onChange={(key, value, focus) => patchStyle(zone.id, key, value, focus)} />
          <div className="flex flex-wrap gap-2">
            <Button onClick={() => nestZone(zone.id)}>Imbriquer</Button>
            <Button onClick={() => removeZone(zone.id)}>Retirer</Button>
          </div>
        </>
      ) : (
        <p className="text-sm text-pretty text-muted">Choisis une zone pour régler son texte, ses actions et son cadre.</p>
      )}
    </aside>
  );
}

function frOf(grid: { rows: { cells: { fr: number; zone?: { id: string }; grid?: { rows: { cells: { fr: number; zone?: { id: string } }[] }[] } }[] }[] }, id: string): number {
  for (const row of grid.rows) {
    for (const cell of row.cells) {
      if (cell.zone?.id === id) return cell.fr;
      if (cell.grid) {
        const nested = frOf(cell.grid, id);
        if (nested > 0) return nested;
      }
    }
  }
  return 0;
}

function Colors({ style, onChange }: { style: { color: string; background: string; titleColor: string }; onChange: (key: "color" | "background" | "titleColor", value: string) => void }) {
  return (
    <div className="grid grid-cols-3 gap-2">
      <Color label="Texte" value={style.color} onChange={(value) => onChange("color", value)} />
      <Color label="Fond" value={style.background} onChange={(value) => onChange("background", value)} />
      <Color label="Titre" value={style.titleColor} onChange={(value) => onChange("titleColor", value)} />
    </div>
  );
}

function Color({ label, value, onChange }: { label: string; value: string; onChange: (value: string) => void }) {
  return (
    <Field label={label}>
      <input className={`${control} px-1`} type="color" value={toHex(value)} onChange={(event) => onChange(event.target.value)} />
    </Field>
  );
}

function toHex(value: string): string {
  return /^#[0-9a-fA-F]{6}$/.test(value) ? value : "#e6edf3";
}

function ActionEditor({ title, action, onChange, loose }: { title: string; action: Action; onChange: (action: Action) => void; loose: boolean }) {
  const set = (partial: Partial<Action>) => onChange({ ...action, ...partial });
  return (
    <fieldset className="flex flex-col gap-2 border border-line p-2">
      <legend className="px-1 text-xs text-muted">{title}</legend>
      <select className={control} value={action.kind} onChange={(event) => set({ kind: event.target.value as ActionKind })}>
        {ACTION_KINDS.map((kind) => (
          <option key={kind || "none"} value={kind}>
            {kind || (loose ? "rien" : "rien")}
          </option>
        ))}
      </select>
      {action.kind && action.kind !== "return" && action.kind !== "save" && action.kind !== "refresh" ? (
        <div className="grid grid-cols-2 gap-2">
          <Field label="Cible">
            <input className={control} value={action.target} onChange={(event) => set({ target: event.target.value })} />
          </Field>
          <Field label="Fichier">
            <input className={control} value={action.file} onChange={(event) => set({ file: event.target.value })} />
          </Field>
          <Field label="Nom">
            <input className={control} value={action.name} onChange={(event) => set({ name: event.target.value })} />
          </Field>
          <Field label="Valeur">
            <input className={control} value={action.value} onChange={(event) => set({ value: event.target.value })} />
          </Field>
          <Field label="href">
            <input className={control} value={action.href} onChange={(event) => set({ href: event.target.value })} />
          </Field>
          <Field label="Commande">
            <input className={control} value={action.run} onChange={(event) => set({ run: event.target.value })} />
          </Field>
          <Field label="Sortie">
            <input className={control} value={action.output} onChange={(event) => set({ output: event.target.value })} />
          </Field>
        </div>
      ) : null}
    </fieldset>
  );
}

function ZoneStyle({ zone, onChange }: { zone: Zone; onChange: (key: keyof Zone["base"], value: Zone["base"][keyof Zone["base"]], focus: boolean) => void }) {
  const style = zone.base;
  return (
    <div className="flex flex-col gap-2">
      <h3 className="text-xs text-muted">Cadre</h3>
      <Colors style={style} onChange={(key, value) => onChange(key, value, false)} />
      <div className="grid grid-cols-2 gap-2">
        <Field label="Bord">
          <input className={control} type="number" min={0} max={1} value={style.borderWidth} onChange={(event) => onChange("borderWidth", Number(event.target.value) ? 1 : 0, false)} />
        </Field>
        <Field label="Rayon">
          <input className={control} type="number" min={0} max={1} value={style.radius} onChange={(event) => onChange("radius", Number(event.target.value) ? 1 : 0, false)} />
        </Field>
        <Field label="Interligne">
          <input className={control} type="number" min={1} value={style.lineHeight} onChange={(event) => onChange("lineHeight", Number(event.target.value) || 1, false)} />
        </Field>
        <Field label="Écart">
          <input className={control} type="number" min={0} value={style.gap} onChange={(event) => onChange("gap", Number(event.target.value) || 0, false)} />
        </Field>
      </div>
      <Quad label="Marge" value={style.mar} onChange={(mar) => onChange("mar", mar, false)} />
      <Quad label="Padding" value={style.pad} onChange={(pad) => onChange("pad", pad, false)} />
      <Field label="Focus, bordure">
        <input className={`${control} px-1`} type="color" value={toHex(zone.focusOn.borderColor)} onChange={(event) => onChange("borderColor", event.target.value, true)} />
      </Field>
    </div>
  );
}

function Quad({
  label,
  value,
  onChange,
}: {
  label: string;
  value: [number, number, number, number];
  onChange: (value: [number, number, number, number]) => void;
}) {
  const names = ["H", "D", "B", "G"];
  return (
    <fieldset>
      <legend className="text-xs text-muted">{label}</legend>
      <div className="grid grid-cols-4 gap-1">
        {names.map((name, index) => (
          <input
            key={name}
            aria-label={`${label} ${name}`}
            className={control}
            type="number"
            min={0}
            value={value[index]}
            onChange={(event) => {
              const next = [...value] as [number, number, number, number];
              next[index] = Number(event.target.value) || 0;
              onChange(next);
            }}
          />
        ))}
      </div>
    </fieldset>
  );
}
