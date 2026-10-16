#!/usr/bin/env bash
# Dark palette, widths, rounded borders, Gum helpers.
# Shellcheck: this file is sourced.

C_BG="#0d1117"
C_PANEL="#161b22"
C_FG="#e6edf3"
C_BORDER="#30363d"
C_ACCENT="#58a6ff"

W_COL=22
W_CHAMP=12
W_VALEUR=28
W_INPUT=40
W_HEADER=$((W_COL * 4 + 6))
ROW_SLOTS=3
BORDER_ROUNDED=rounded

style_panel() {
  local focused=$1
  local text=$2
  local border_fg=$C_BORDER
  local -a extra=()
  if [[ $focused == 1 ]]; then
    border_fg=$C_ACCENT
    extra+=(--bold)
  fi
  gum style \
    --foreground "$C_FG" \
    --background "$C_PANEL" \
    --border "$BORDER_ROUNDED" \
    --border-foreground "$border_fg" \
    --padding "0 1" \
    --width "$W_COL" \
    "${extra[@]}" \
    "$text"
}

style_banner() {
  gum style \
    --foreground "$C_FG" \
    --background "$C_BG" \
    --border "$BORDER_ROUNDED" \
    --border-foreground "$C_ACCENT" \
    --padding "0 1" \
    --width "$W_HEADER" \
    --bold \
    "$1"
}

style_status() {
  gum style \
    --foreground "$C_FG" \
    --background "$C_PANEL" \
    --border "$BORDER_ROUNDED" \
    --border-foreground "$C_BORDER" \
    --padding "0 1" \
    --width "$W_HEADER" \
    "$1"
}

join_horizontal() {
  gum join --horizontal "$@"
}

join_vertical() {
  gum join --vertical "$@"
}

table_print() {
  gum table --print --separator $'\t' \
    --columns "Champ,Valeur" \
    --widths "${W_CHAMP},${W_VALEUR}" \
    --border "$BORDER_ROUNDED" \
    --border.foreground "$C_BORDER" \
    --header.foreground "$C_ACCENT" \
    --header.background "$C_BG" \
    --cell.foreground "$C_FG" \
    --cell.background "$C_PANEL"
}

input_value() {
  local header=$1
  local value=$2
  gum input \
    --header "$header" \
    --value "$value" \
    --width "$W_INPUT" \
    --prompt "› " \
    --header.foreground "$C_FG" \
    --prompt.foreground "$C_ACCENT" \
    --cursor.foreground "$C_ACCENT"
}
