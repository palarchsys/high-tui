#!/usr/bin/env bash
# Gum wrappers. Colors and widths are arguments, never literals.
# Shellcheck: this file is sourced.

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

style_box() {
  local fg=$1 bg=$2 bc=$3 bw=$4 bold=$5 width=$6 text=$7
  local -a args
  args=(--foreground "$fg" --background "$bg" --width "$width")
  if (( bw == 1 )); then
    args+=(--border rounded --border-foreground "$bc")
  fi
  if (( bold == 1 )); then
    args+=(--bold)
  fi
  gum style "${args[@]}" "$text"
}

input_value() {
  local header=$1 value=$2 width=$3
  gum input \
    --header "$header" \
    --value "$value" \
    --width "$width" \
    --prompt "› " \
    --header.foreground "$FG_UI" \
    --prompt.foreground "$ACCENT_UI" \
    --cursor.foreground "$ACCENT_UI"
}

write_value() {
  local header=$1 value=$2 width=$3 height=$4
  gum write \
    --header "$header" \
    --value "$value" \
    --width "$width" \
    --height "$height" \
    --no-show-help \
    --header.foreground "$FG_UI" \
    --prompt.foreground "$ACCENT_UI" \
    --cursor.foreground "$ACCENT_UI"
}
