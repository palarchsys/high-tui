#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "$0")"
python3 engine.py
source theme.sh
source model.sh
source blocks.sh
source screen.sh

load_layout
init_values
warm_blocks

if [[ ${DRAW_ONCE:-} == 1 ]]; then
  paint_inplace "$(frame_text)"
  exit 0
fi

if [[ -t 1 ]]; then
  printf '\e[?1049h\e[?25l'
  trap restore_screen EXIT
fi

while true; do
  paint_inplace "$(frame_text)"
  read -rsn1 key
  if [[ $key == $'\e' ]]; then
    read -rsn2 -t 0.05 rest || rest=
    case $rest in
      '[A') move_vert -1 ;;
      '[B') move_vert 1 ;;
      '[C') move_focus 1 ;;
      '[D') move_focus -1 ;;
      '') nav_return ;;
    esac
  elif [[ $key == q ]]; then
    exit 0
  elif [[ -z $key || $key == $'\r' ]]; then
    activate
  fi
done
