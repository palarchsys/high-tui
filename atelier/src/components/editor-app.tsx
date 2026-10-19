import { useEffect, useMemo, useState, type ReactNode } from "react";
import { Database, Download, Eye, LayoutGrid, Redo2, RotateCcw, Undo2 } from "lucide-react";
import { Button } from "@/components/bits";
import { DataPane } from "@/components/data-pane";
import { DisplayPane } from "@/components/display-pane";
import { PreviewPane } from "@/components/preview-pane";
import { currentWin, exported, issuesOf, useEditor } from "@/editor/store";

export function EditorApp() {
  const [open, setOpen] = useState(false);
  const tab = useEditor((state) => state.tab);
  const setTab = useEditor((state) => state.setTab);
  const windows = useEditor((state) => state.windows);
  const sources = useEditor((state) => state.sources);
  const activeId = useEditor((state) => state.activeId);
  const undo = useEditor((state) => state.undo);
  const redo = useEditor((state) => state.redo);
  const reset = useEditor((state) => state.reset);
  const selectWindow = useEditor((state) => state.selectWindow);
  const addWindow = useEditor((state) => state.addWindow);
  const setDrag = useEditor((state) => state.setDrag);
  const issues = useMemo(() => issuesOf({ windows, sources, activeId }), [windows, sources, activeId]);

  useEffect(() => {
    void Promise.resolve(useEditor.persist.rehydrate());
  }, []);

  useEffect(() => {
    const onKey = (event: KeyboardEvent) => {
      if (!(event.metaKey || event.ctrlKey) || event.key.toLowerCase() !== "z") return;
      event.preventDefault();
      if (event.shiftKey) useEditor.getState().redo();
      else useEditor.getState().undo();
    };
    const onUp = () => setDrag(null);
    window.addEventListener("keydown", onKey);
    window.addEventListener("pointerup", onUp);
    return () => {
      window.removeEventListener("keydown", onKey);
      window.removeEventListener("pointerup", onUp);
    };
  }, [setDrag]);

  return (
    <main className="flex h-dvh flex-col bg-bg text-fg">
      <header className="flex flex-wrap items-center gap-2 border-b border-line px-3 py-2">
        <div className="mr-2 flex items-center gap-2">
          <span className="size-3 bg-accent" aria-hidden="true" />
          <p className="text-sm">
            Atelier <span className="text-muted">high-tui</span>
          </p>
        </div>
        <div className="flex gap-1 overflow-x-auto">
          {windows.map((win) => (
            <button
              key={win.id}
              type="button"
              className={`min-h-11 rounded-sm px-3 text-sm ${win.id === activeId ? "bg-surface text-accent" : "text-fg"}`}
              onClick={() => selectWindow(win.id)}
            >
              {win.label}
            </button>
          ))}
          <button type="button" className="min-h-11 px-3 text-sm text-muted" onClick={addWindow}>
            Fenêtre
          </button>
        </div>
        <div className="ml-auto flex flex-wrap gap-1">
          <IconButton label="Annuler" onClick={undo}>
            <Undo2 className="size-4" />
          </IconButton>
          <IconButton label="Rétablir" onClick={redo}>
            <Redo2 className="size-4" />
          </IconButton>
          <IconButton label="Exemples" onClick={reset}>
            <RotateCcw className="size-4" />
          </IconButton>
          <Button tone={issues.length ? "ghost" : "accent"} onClick={() => setOpen(true)}>
            {issues.length ? `${issues.length} erreurs` : "Exporter"}
          </Button>
        </div>
      </header>
      <nav className="flex border-b border-line lg:hidden">
        <Tab on={tab === "display"} onClick={() => setTab("display")} icon={<LayoutGrid className="size-4" />} label="Affichage" />
        <Tab on={tab === "data"} onClick={() => setTab("data")} icon={<Database className="size-4" />} label="Données" />
        <Tab on={tab === "preview"} onClick={() => setTab("preview")} icon={<Eye className="size-4" />} label="Aperçu" />
      </nav>
      <div className="flex min-h-0 min-w-0 flex-1 flex-col overflow-hidden">
        <div className="flex min-h-0 flex-1 flex-col lg:flex-row">
          <div className={tab === "display" ? "flex min-h-0 min-w-0 flex-1" : "hidden min-h-0 min-w-0 flex-1 lg:flex"}>
            <DisplayPane />
          </div>
          <div className={tab === "data" ? "flex min-h-0 min-w-0 flex-1" : "hidden min-h-0 min-w-0 flex-1 lg:flex"}>
            <DataPane />
          </div>
        </div>
        <div className={tab === "preview" ? "block" : "hidden lg:block"}>
          <PreviewPane />
        </div>
      </div>
      {open ? <ExportDialog issues={issues} onClose={() => setOpen(false)} /> : null}
    </main>
  );
}

function Tab({ on, onClick, label, icon }: { on: boolean; onClick: () => void; label: string; icon: ReactNode }) {
  return (
    <button type="button" className={`flex min-h-11 flex-1 items-center justify-center gap-2 text-sm ${on ? "text-accent" : "text-muted"}`} onClick={onClick}>
      {icon}
      {label}
    </button>
  );
}

function IconButton({ label, onClick, children }: { label: string; onClick: () => void; children: ReactNode }) {
  return (
    <button type="button" aria-label={label} title={label} className="inline-flex size-11 items-center justify-center rounded-sm border border-line" onClick={onClick}>
      {children}
    </button>
  );
}

function ExportDialog({ issues, onClose }: { issues: string[]; onClose: () => void }) {
  const win = useEditor((state) => currentWin(state));
  const importText = useEditor((state) => state.importText);
  const files = exported(win);
  const [html, setHtml] = useState("");
  const [css, setCss] = useState("");
  const [message, setMessage] = useState<string | null>(null);
  return (
    <div className="fixed inset-0 z-20 flex items-end justify-center bg-bg/80 p-3 sm:items-center">
      <div className="flex max-h-[90vh] w-full max-w-3xl flex-col gap-3 overflow-auto border border-line bg-surface p-4">
        <div className="flex items-center justify-between gap-2">
          <h2 className="text-sm text-accent">Export</h2>
          <Button onClick={onClose}>Fermer</Button>
        </div>
        <p className="text-sm text-pretty text-muted">
          {win.file} et {win.file.replace(/\.html$/, ".css")}. Le générateur Go n'est pas lancé ici.
        </p>
        {issues.length ? (
          <ul className="flex flex-col gap-1 text-sm text-accent">
            {issues.map((issue) => (
              <li key={issue}>{issue}</li>
            ))}
          </ul>
        ) : (
          <p className="text-sm text-fg">Aucune erreur. Les fichiers suivent le dialecte actuel.</p>
        )}
        <div className="flex flex-wrap gap-2">
          <Button
            tone="accent"
            disabled={issues.length > 0}
            onClick={() => {
              download(win.file, files.html);
              download(win.file.replace(/\.html$/, ".css"), files.css);
            }}
          >
            <Download className="size-4" />
            Télécharger
          </Button>
        </div>
        <pre className="max-h-48 overflow-auto border border-line bg-bg p-3 font-mono text-xs text-fg">{files.html}</pre>
        <h3 className="text-sm">Importer</h3>
        <textarea className="min-h-24 border border-line bg-bg p-2 font-mono text-xs" placeholder="HTML" value={html} onChange={(event) => setHtml(event.target.value)} />
        <textarea className="min-h-24 border border-line bg-bg p-2 font-mono text-xs" placeholder="CSS" value={css} onChange={(event) => setCss(event.target.value)} />
        <Button
          onClick={() => {
            const error = importText(html, css);
            setMessage(error);
            if (!error) onClose();
          }}
        >
          Remplacer ou ajouter
        </Button>
        {message ? <p className="text-sm text-accent">{message}</p> : null}
      </div>
    </div>
  );
}

function download(name: string, text: string) {
  const url = URL.createObjectURL(new Blob([text], { type: "text/plain;charset=utf-8" }));
  const link = document.createElement("a");
  link.href = url;
  link.download = name;
  link.click();
  URL.revokeObjectURL(url);
}
