import { useMemo, useState } from "react";
import { Button, control, Field } from "@/components/bits";
import { itemsFromSource, parseEnv, parseItems, writeEnv, writeItems } from "@/editor/data";
import { emptyItem, findZone, zonesOf } from "@/editor/model";
import { currentWin, useEditor } from "@/editor/store";
import type { Item, Source, Zone } from "@/editor/types";

export function DataPane() {
  const win = useEditor((state) => currentWin(state));
  const sources = useEditor((state) => state.sources);
  const sel = useEditor((state) => state.sel);
  const [sourceId, setSourceId] = useState(sources[0]?.id ?? "");
  const source = sources.find((item) => item.id === sourceId) ?? sources[0];
  const zone = sel?.kind === "zone" ? findZone(win, sel.id) : undefined;
  const wires = useMemo(() => connections(win), [win]);

  return (
    <section className="flex min-h-0 min-w-0 flex-1 flex-col overflow-auto border-line lg:border-l">
      <div className="border-b border-line p-3">
        <h2 className="text-sm text-accent">Fils</h2>
        {wires.length === 0 ? <p className="mt-2 text-sm text-muted">Aucun branchement.</p> : null}
        <ul className="mt-2 flex flex-col gap-1">
          {wires.map((wire) => (
            <li key={wire} className="text-sm text-fg">
              {wire}
            </li>
          ))}
        </ul>
      </div>
      <div className="grid min-h-0 flex-1 lg:grid-cols-[16rem_1fr]">
        <div className="border-b border-line p-3 lg:border-r lg:border-b-0">
          <div className="mb-2 flex items-center justify-between gap-2">
            <h2 className="text-sm text-accent">Sources</h2>
            <Button onClick={() => useEditor.getState().addSource()}>Ajouter</Button>
          </div>
          <ul className="flex flex-col">
            {sources.map((item) => (
              <li key={item.id}>
                <button
                  type="button"
                  className={`min-h-11 w-full truncate rounded-sm px-2 text-left text-sm ${item.id === source?.id ? "bg-surface text-accent" : "text-fg"}`}
                  onClick={() => setSourceId(item.id)}
                >
                  {item.path}
                </button>
              </li>
            ))}
          </ul>
        </div>
        <div className="flex flex-col gap-4 p-3">
          {source ? <SourceEditor source={source} zone={zone} /> : <p className="text-sm text-muted">Aucune source.</p>}
          {zone?.tag === "nav" ? <NextData zone={zone} navs={zonesOf(win).filter((item) => item.tag === "nav")} /> : null}
          {zone && zone.tag !== "nav" ? <BindZone zone={zone} sources={sources} /> : null}
        </div>
      </div>
    </section>
  );
}

function connections(win: ReturnType<typeof currentWin>): string[] {
  const lines: string[] = [];
  for (const zone of zonesOf(win)) {
    if (zone.itemsFrom) lines.push(`${zone.itemsFrom}  →  ${zone.id} · lignes`);
    if (zone.file && zone.name) lines.push(`${zone.file} · ${zone.name}  →  ${zone.id} · variable`);
    else if (zone.file) lines.push(`${zone.file}  →  ${zone.id} · brut`);
    if (zone.nextNav) lines.push(`${zone.id}  →  ${zone.nextNav} · next-data`);
    if (zone.label.includes("::")) lines.push(`${zone.label}  →  ${zone.id} · référence`);
    if (zone.value.includes("::")) lines.push(`${zone.value}  →  ${zone.id} · référence`);
  }
  return lines;
}

function SourceEditor({ source, zone }: { source: Source; zone: Zone | undefined }) {
  const setSource = useEditor((state) => state.setSource);
  const bind = useEditor((state) => state.bind);
  const setReference = useEditor((state) => state.setReference);
  const removeSource = useEditor((state) => state.removeSource);
  const pairs = source.read === "env" ? parseEnv(source.body) : [];
  const parsed = source.read === "json" ? parseItems(source.body, source.path) : itemsFromSource(source, source.path);
  const [name, setName] = useState(pairs[0]?.name ?? "");

  return (
    <div className="flex flex-col gap-3">
      <div className="grid gap-2 sm:grid-cols-2">
        <Field label="Chemin">
          <input className={control} value={source.path} onChange={(event) => setSource(source.id, { path: event.target.value })} />
        </Field>
        <Field label="Lecture">
          <select
            className={control}
            value={source.read}
            onChange={(event) => setSource(source.id, { read: event.target.value as Source["read"] })}
          >
            <option value="text">brut, sans split</option>
            <option value="env">split env label/value</option>
            <option value="json">liste json</option>
            <option value="lines">lignes séparées</option>
          </select>
        </Field>
      </div>
      {source.read === "lines" ? (
        <Field label="Séparateur">
          <input className={control} value={source.separator} onChange={(event) => setSource(source.id, { separator: event.target.value })} />
        </Field>
      ) : null}
      {source.projectedFrom ? <p className="text-xs text-muted">Projection de {source.projectedFrom}. Le générateur lit ce JSON.</p> : null}
      {source.read === "env" ? (
        <EnvRows
          body={source.body}
          onChange={(body) => setSource(source.id, { body })}
        />
      ) : source.read === "json" && !parsed.error ? (
        <JsonRows items={parsed.items} onChange={(items) => setSource(source.id, { body: writeItems(items) })} />
      ) : (
        <Field label="Contenu">
          <textarea className={`${control} min-h-40 py-2 font-mono`} value={source.body} onChange={(event) => setSource(source.id, { body: event.target.value })} />
        </Field>
      )}
      {parsed.error ? <p className="text-sm text-accent">{parsed.error}</p> : <p className="text-xs text-muted">{parsed.items.length} enregistrement(s)</p>}
      {zone ? (
        <div className="flex flex-wrap gap-2">
          {zone.tag === "nav" ? <Button onClick={() => bind(zone.id, "list", source.id)}>Boucle vers {zone.id}</Button> : null}
          {zone.tag === "input" || zone.tag === "textarea" ? (
            <>
              <Button onClick={() => bind(zone.id, "raw", source.id)}>Brut vers {zone.id}</Button>
              {source.read === "env" ? (
                <>
                  <select className={`${control} w-auto`} value={name} onChange={(event) => setName(event.target.value)}>
                    {pairs.map((pair) => (
                      <option key={pair.name} value={pair.name}>
                        {pair.name}
                      </option>
                    ))}
                  </select>
                  <Button onClick={() => bind(zone.id, "var", source.id, name)}>Variable vers {zone.id}</Button>
                  <Button onClick={() => setReference(zone.id, source.id, name, "label")}>Libellé ::label</Button>
                  <Button onClick={() => setReference(zone.id, source.id, name, "value")}>Valeur ::value</Button>
                </>
              ) : null}
              <Button onClick={() => bind(zone.id, "none", source.id)}>Détacher</Button>
            </>
          ) : null}
        </div>
      ) : (
        <p className="text-sm text-muted">Choisis une zone dans l'affichage pour la brancher.</p>
      )}
      <Button tone="quiet" onClick={() => removeSource(source.id)}>Retirer la source</Button>
    </div>
  );
}

function EnvRows({ body, onChange }: { body: string; onChange: (body: string) => void }) {
  const pairs = parseEnv(body);
  const write = (next: typeof pairs) => onChange(writeEnv(body, next));
  return (
    <div className="flex flex-col gap-2">
      {pairs.map((pair, index) => (
        <div key={pair.name} className="grid gap-2 sm:grid-cols-3">
          <input className={control} aria-label="Nom" value={pair.name} onChange={(event) => write(pairs.map((item, i) => (i === index ? { ...item, name: event.target.value } : item)))} />
          <input className={control} aria-label="Label" value={pair.label} onChange={(event) => write(pairs.map((item, i) => (i === index ? { ...item, label: event.target.value } : item)))} />
          <input className={control} aria-label="Valeur" value={pair.value} onChange={(event) => write(pairs.map((item, i) => (i === index ? { ...item, value: event.target.value } : item)))} />
        </div>
      ))}
      <Button
        onClick={() => write([...pairs, { name: `NOM_${pairs.length + 1}`, label: "Label", value: "" }])}
      >
        Ajouter une variable
      </Button>
    </div>
  );
}

function JsonRows({ items, onChange, depth = 0 }: { items: Item[]; onChange: (items: Item[]) => void; depth?: number }) {
  const set = (index: number, item: Item) => onChange(items.map((current, i) => (i === index ? item : current)));
  return (
    <div className="flex flex-col gap-2">
      {items.map((item, index) => (
        <ItemCard key={`${item.id}-${index}`} item={item} depth={depth} onChange={(next) => set(index, next)} onRemove={() => onChange(items.filter((_, i) => i !== index))} />
      ))}
      <Button onClick={() => onChange([...items, emptyItem({ label: "Ligne" })])}>Ajouter une ligne</Button>
    </div>
  );
}

function ItemCard({ item, onChange, onRemove, depth }: { item: Item; onChange: (item: Item) => void; onRemove: () => void; depth: number }) {
  return (
    <div className="flex flex-col gap-2 border border-line p-2">
      <div className="grid gap-2 sm:grid-cols-2">
        <Field label="Label">
          <input className={control} value={item.label} onChange={(event) => onChange({ ...item, label: event.target.value })} />
        </Field>
        <Field label="Valeur">
          <input className={control} value={item.value} onChange={(event) => onChange({ ...item, value: event.target.value })} />
        </Field>
        <Field label="Fichier">
          <input className={control} value={item.file} onChange={(event) => onChange({ ...item, file: event.target.value })} />
        </Field>
        <Field label="Nom">
          <input className={control} value={item.name} onChange={(event) => onChange({ ...item, name: event.target.value })} />
        </Field>
      </div>
      <Field label="next-data">
        <select
          className={control}
          value={item.nextKind}
          onChange={(event) => onChange({ ...item, nextKind: event.target.value as Item["nextKind"] })}
        >
          <option value="none">aucun</option>
          <option value="file">fichier</option>
          <option value="cmd">commande</option>
          <option value="inline">liste inline</option>
        </select>
      </Field>
      {item.nextKind === "file" || item.nextKind === "cmd" ? (
        <input className={control} value={item.next} onChange={(event) => onChange({ ...item, next: event.target.value })} />
      ) : null}
      {item.nextKind === "inline" && depth < 2 ? (
        <JsonRows items={item.nextItems} depth={depth + 1} onChange={(nextItems) => onChange({ ...item, nextItems })} />
      ) : null}
      <Button tone="quiet" onClick={onRemove}>Retirer la ligne</Button>
    </div>
  );
}

function NextData({ zone, navs }: { zone: Zone; navs: Zone[] }) {
  const patchZone = useEditor((state) => state.patchZone);
  return (
    <div className="border border-line p-3">
      <h3 className="text-sm text-accent">next-data</h3>
      <p className="mt-1 text-xs text-pretty text-muted">
        data-next nomme la nav qui reçoit le next de la ligne choisie. Le contenu se règle dans la source JSON, pas dans les lignes d'aperçu.
      </p>
      <Field label="Nav cible">
        <select className={control} value={zone.nextNav} onChange={(event) => patchZone(zone.id, { nextNav: event.target.value })}>
          <option value="">aucune</option>
          {navs
            .filter((nav) => nav.id !== zone.id)
            .map((nav) => (
              <option key={nav.id} value={nav.id}>
                {nav.id}
              </option>
            ))}
        </select>
      </Field>
    </div>
  );
}

function BindZone({ zone, sources }: { zone: Zone; sources: Source[] }) {
  return (
    <p className="text-sm text-muted">
      {zone.id}
      {zone.file ? ` lit ${zone.file}${zone.name ? ` · ${zone.name}` : ""}` : " n'a pas encore de fichier."} {sources.length ? "" : ""}
    </p>
  );
}
