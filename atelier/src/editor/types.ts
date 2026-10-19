export type ActionKind =
  | ""
  | "load"
  | "edit"
  | "set"
  | "window"
  | "return"
  | "save"
  | "refresh"
  | "run";

export type Action = {
  kind: ActionKind;
  target: string;
  file: string;
  name: string;
  href: string;
  run: string;
  output: string;
  value: string;
};

export type NextKind = "none" | "file" | "cmd" | "inline";

export type Item = {
  id: string;
  label: string;
  value: string;
  file: string;
  name: string;
  nextKind: NextKind;
  next: string;
  nextItems: Item[];
  focus: Action;
  enter: Action;
};

export type Style = {
  color: string;
  titleColor: string;
  background: string;
  borderColor: string;
  borderWidth: number;
  radius: number;
  lineHeight: number;
  gap: number;
  maxWidth: number;
  bold: boolean;
  pad: [number, number, number, number];
  mar: [number, number, number, number];
};

export type ZoneTag = "header" | "nav" | "input" | "textarea" | "button" | "footer";

export type Zone = {
  tag: ZoneTag;
  id: string;
  classes: string[];
  label: string;
  value: string;
  file: string;
  name: string;
  save: "" | "enter" | "button";
  readonly: boolean;
  focus: Action;
  enter: Action;
  itemsFrom: string;
  nextNav: string;
  items: Item[];
  base: Style;
  focusOn: Style;
};

export type Cell = { fr: number; grid?: Grid; zone?: Zone };
export type GridRow = { id: string; lines: number; cells: Cell[] };
export type Grid = { id: string; rows: GridRow[] };

export type Win = {
  id: string;
  label: string;
  file: string;
  style: Style;
  grid: Grid;
};

export type SourceRead = "text" | "env" | "json" | "lines";

export type Source = {
  id: string;
  path: string;
  body: string;
  read: SourceRead;
  separator: string;
  projectedFrom?: string;
};

export type Sel =
  | { kind: "zone"; id: string }
  | { kind: "row"; id: string }
  | { kind: "grid"; id: string };

export const ACTION_KINDS: ActionKind[] = [
  "",
  "load",
  "edit",
  "set",
  "window",
  "return",
  "save",
  "refresh",
  "run",
];

export const ZONE_TAGS: ZoneTag[] = ["header", "nav", "input", "textarea", "button", "footer"];
