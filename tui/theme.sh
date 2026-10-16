#!/usr/bin/env bash
# Dark palette, widths, rounded borders.
# The only gum style, gum table, and gum input wrappers.
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
# Padding is inside --width. Gum wraps past this inner width.
PAD_X=1
INNER_COL=$((W_COL - PAD_X * 2))
INNER_HEADER=$((W_HEADER - PAD_X * 2))
BORDER_ROUNDED=rounded

# One width rule: cut, then pad with spaces. Newlines cannot survive.
clip_pad() {
  local text=$1
  local width=$2
  text=${text//$'\t'/ }
  text=${text//$'\n'/ }
  text=${text//$'\r'/ }
  if (( ${#text} > width )); then
    text=${text:0:width}
  fi
  printf -v text '%s%*s' "$text" $((width - ${#text})) ''
  printf '%s' "$text"
}

clip_lines() {
  local text=$1
  local width=$2
  local slots=$3
  local -a lines=()
  local line n i out=""
  while IFS= read -r line || [[ -n ${line:-} ]]; do
    lines+=("$(clip_pad "$line" "$width")")
  done <<< "$text"
  n=${#lines[@]}
  if (( n > slots )); then
    lines=("${lines[@]:0:slots}")
    n=$slots
  fi
  while (( n < slots )); do
    lines+=("$(clip_pad "" "$width")")
    n=$((n + 1))
  done
  for i in "${!lines[@]}"; do
    if (( i > 0 )); then
      out+=$'\n'
    fi
    out+="${lines[$i]}"
  done
  printf '%s' "$out"
}

style_panel() {
  local focused=$1
  local text=$2
  local border_fg=$C_BORDER
  local -a extra=()
  text=$(clip_lines "$text" "$INNER_COL" $((ROW_SLOTS + 1)))
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
  local text
  text=$(clip_pad "$1" "$INNER_HEADER")
  gum style \
    --foreground "$C_FG" \
    --background "$C_BG" \
    --border "$BORDER_ROUNDED" \
    --border-foreground "$C_ACCENT" \
    --padding "0 1" \
    --width "$W_HEADER" \
    --bold \
    "$text"
}

style_status() {
  local text
  text=$(clip_pad "$1" "$INNER_HEADER")
  gum style \
    --foreground "$C_FG" \
    --background "$C_PANEL" \
    --border "$BORDER_ROUNDED" \
    --border-foreground "$C_BORDER" \
    --padding "0 1" \
    --width "$W_HEADER" \
    "$text"
}

table_print() {
  local champ val
  local -a rows=()
  local n=0
  while IFS=$'\t' read -r champ val || [[ -n ${champ:-}${val:-} ]]; do
    if (( n >= ROW_SLOTS )); then
      break
    fi
    rows+=("$(clip_pad "${champ:-}" "$W_CHAMP")"$'\t'"$(clip_pad "${val:-}" "$W_VALEUR")")
    n=$((n + 1))
  done
  while (( n < ROW_SLOTS )); do
    rows+=("$(clip_pad "" "$W_CHAMP")"$'\t'"$(clip_pad "" "$W_VALEUR")")
    n=$((n + 1))
  done
  printf '%s\n' "${rows[@]}" | gum table --print --separator $'\t' \
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
