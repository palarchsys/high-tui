#!/usr/bin/env bash
# Horizontal line paste, vertical stack, in-place paint. No gum.
# Shellcheck: this file is sourced.

split_lines() {
  local data=${1-}
  local -n _out=$2
  local line
  _out=()
  data=${data%$'\n'}
  [[ -z $data ]] && return 0
  while IFS= read -r line || [[ -n ${line:-} ]]; do
    _out+=("$line")
  done <<< "$data"
}

paste_horizontal() {
  local -a blocks=("$@")
  local -a chunk=() flat=() counts=() bases=()
  local b line acc=0 maxh=0
  local row i out first=1
  for b in "${blocks[@]}"; do
    split_lines "$b" chunk
    counts+=("${#chunk[@]}")
    if (( ${#chunk[@]} > maxh )); then
      maxh=${#chunk[@]}
    fi
    flat+=("${chunk[@]}")
  done
  acc=0
  for b in "${counts[@]}"; do
    bases+=("$acc")
    acc=$((acc + b))
  done
  for ((row = 0; row < maxh; row++)); do
    out=""
    for i in "${!blocks[@]}"; do
      if (( row < counts[i] )); then
        out+="${flat[$((bases[i] + row))]}"
      fi
    done
    if (( first == 0 )); then
      printf '\n'
    fi
    first=0
    printf '%s' "$out"
  done
}

stack_vertical() {
  local -a blocks=("$@")
  local -a chunk=()
  local b line first=1
  for b in "${blocks[@]}"; do
    split_lines "$b" chunk
    for line in "${chunk[@]}"; do
      if (( first == 0 )); then
        printf '\n'
      fi
      first=0
      printf '%s' "$line"
    done
  done
}

paint_inplace() {
  local screen=$1
  if [[ -t 1 && ${DRAW_ONCE:-} != 1 ]]; then
    printf '\e[?2026h\e[H%s\e[0J\e[?2026l\n' "$screen"
  else
    printf '%s\n' "$screen"
  fi
}

restore_screen() {
  printf '\e[?25h\e[?1049l'
}
