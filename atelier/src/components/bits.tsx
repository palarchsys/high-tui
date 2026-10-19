import type { ReactNode } from "react";
import type { ZoneTag } from "@/editor/types";

export const TAGS: { tag: ZoneTag; label: string }[] = [
  { tag: "header", label: "Titre" },
  { tag: "nav", label: "Liste" },
  { tag: "input", label: "Champ" },
  { tag: "textarea", label: "Texte" },
  { tag: "button", label: "Bouton" },
  { tag: "footer", label: "Statut" },
];

export function tagLabel(tag: ZoneTag): string {
  return TAGS.find((item) => item.tag === tag)?.label ?? tag;
}

export function Field({ label, children }: { label: string; children: ReactNode }) {
  return (
    <label className="flex flex-col gap-1 text-xs text-muted">
      <span>{label}</span>
      {children}
    </label>
  );
}

export const control =
  "min-h-11 w-full rounded-sm border border-line bg-bg px-2 text-sm text-fg outline-none focus:border-accent";

export function Button({
  children,
  onClick,
  tone = "ghost",
  type = "button",
  disabled,
}: {
  children: ReactNode;
  onClick?: () => void;
  tone?: "ghost" | "accent" | "quiet";
  type?: "button" | "submit";
  disabled?: boolean;
}) {
  const look =
    tone === "accent"
      ? "bg-accent text-bg"
      : tone === "quiet"
        ? "text-muted"
        : "border border-line bg-surface text-fg";
  return (
    <button
      type={type}
      disabled={disabled}
      onClick={onClick}
      className={`inline-flex min-h-11 items-center justify-center gap-2 rounded-sm px-3 text-sm disabled:opacity-40 ${look}`}
    >
      {children}
    </button>
  );
}
