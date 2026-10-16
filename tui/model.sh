#!/usr/bin/env bash
# Layout, focus, files, navigation, commands.
# Shellcheck: this file is sourced.

WIN_ID=()
Z_WIN=() Z_ID=() Z_ROW=() Z_TAG=()
Z_X=() Z_Y=() Z_OW=() Z_OH=() Z_TW=() Z_LINES=() Z_CH=()
Z_FG=() Z_BG=() Z_BC=() Z_BW=() Z_BOLD=()
Z_FFG=() Z_FBG=() Z_FBC=() Z_FBOLD=()
Z_PT=() Z_PB=() Z_PL=() Z_PR=() Z_LH=() Z_GAP=()
N_WIN=() N_NAV=() N_INDEX=() N_ID=() N_TITLE=()
N_ACTION=() N_TARGET=() N_FILE=() N_NAME=() N_HREF=()
F_WIN=() F_ID=()
declare -gA ZIX=() PROP=() WIN_MAX=() WIN_TITLE=() WIN_FILE=()
declare -gA WIN_BG=() WIN_FG=() WIN_H=() WIN_BY_FILE=()
declare -gA SEL=() VAL=() DIRTY=()
CUR=
FOCUS_I=0
STACK=()
status_msg="Prêt"
FG_UI=
ACCENT_UI=
MAXW=0
SPINNER=(⠋ ⠙ ⠹ ⠸)

load_layout() {
  local kind w id row idx tag x y ow oh tw lines ch
  local fg bg bc bw bold ffg fbg fbc fbold pt pb pl pr lh gap
  local file title bgc fgc total inner
  local nav index iid action target name href key val order
  [[ -f layout.tsv ]] || return 1
  while IFS=$'\x1f' read -r kind w id row idx tag x y ow oh tw lines ch \
      fg bg bc bw bold ffg fbg fbc fbold pt pb pl pr lh gap; do
    [[ -z ${kind:-} ]] && continue
    case $kind in
      W)
        WIN_ID+=("$w")
        WIN_MAX[$w]=$id
        WIN_TITLE[$w]=$row
        WIN_FILE[$w]=$idx
        WIN_BG[$w]=$tag
        WIN_FG[$w]=$x
        WIN_H[$w]=$y
        WIN_BY_FILE[$idx]=$w
        ;;
      Z)
        ZIX["$w|$id"]=${#Z_WIN[@]}
        Z_WIN+=("$w"); Z_ID+=("$id"); Z_ROW+=("$row"); Z_TAG+=("$tag")
        Z_X+=("$x"); Z_Y+=("$y"); Z_OW+=("$ow"); Z_OH+=("$oh")
        Z_TW+=("$tw"); Z_LINES+=("$lines"); Z_CH+=("$ch")
        Z_FG+=("$fg"); Z_BG+=("$bg"); Z_BC+=("$bc"); Z_BW+=("$bw"); Z_BOLD+=("$bold")
        Z_FFG+=("$ffg"); Z_FBG+=("$fbg"); Z_FBC+=("$fbc"); Z_FBOLD+=("$fbold")
        Z_PT+=("$pt"); Z_PB+=("$pb"); Z_PL+=("$pl"); Z_PR+=("$pr")
        Z_LH+=("$lh"); Z_GAP+=("$gap")
        ;;
      P)
        PROP["$w|$id|$row"]=$idx
        ;;
      N)
        N_WIN+=("$w"); N_NAV+=("$id"); N_INDEX+=("$row"); N_ID+=("$idx")
        N_TITLE+=("$tag"); N_ACTION+=("$x"); N_TARGET+=("$y"); N_FILE+=("$ow")
        N_NAME+=("$oh"); N_HREF+=("$tw")
        ;;
      F)
        F_WIN+=("$w"); F_ID+=("$row")
        ;;
    esac
  done < layout.tsv
  CUR=${WIN_ID[0]}
  MAXW=${WIN_MAX[$CUR]}
  FG_UI=${WIN_FG[$CUR]}
  local i
  for i in "${!Z_WIN[@]}"; do
    if [[ ${Z_WIN[$i]} == "$CUR" && ${Z_TAG[$i]} == nav ]]; then
      ACCENT_UI=${Z_FBC[$i]}
      break
    fi
  done
}

init_values() {
  local i id w tag
  for i in "${!Z_WIN[@]}"; do
    w=${Z_WIN[$i]}
    id=${Z_ID[$i]}
    tag=${Z_TAG[$i]}
    if [[ $tag == input ]]; then
      VAL["$w|$id"]=$(read_var "${PROP[$w|$id|file]}" "${PROP[$w|$id|name]}")
    elif [[ $tag == textarea ]]; then
      VAL["$w|$id"]=$(<"${PROP[$w|$id|file]}")
    fi
    if [[ $tag == nav ]]; then
      SEL["$w|$id"]=0
    fi
    DIRTY["$w|$id"]=0
  done
}

read_var() {
  local file=$1 name=$2 line key
  while IFS= read -r line || [[ -n ${line:-} ]]; do
    [[ $line =~ ^([A-Za-z_][A-Za-z0-9_]*)=(.*)$ ]] || continue
    key=${BASH_REMATCH[1]}
    if [[ $key == "$name" ]]; then
      printf '%s' "${BASH_REMATCH[2]}"
      return 0
    fi
  done < "$file"
}

write_var() {
  local file=$1 name=$2 value=$3
  local tmp line key found=0
  tmp=$(mktemp)
  while IFS= read -r line || [[ -n ${line:-} ]]; do
    if [[ $line =~ ^([A-Za-z_][A-Za-z0-9_]*)=(.*)$ ]]; then
      key=${BASH_REMATCH[1]}
      if [[ $key == "$name" ]]; then
        printf '%s=%s\n' "$name" "$value"
        found=1
        continue
      fi
    fi
    printf '%s\n' "$line"
  done < "$file" > "$tmp"
  if (( found == 0 )); then
    printf '%s=%s\n' "$name" "$value" >> "$tmp"
  fi
  mv "$tmp" "$file"
}

focused_id() {
  local -a ids=()
  local i
  for i in "${!F_WIN[@]}"; do
    if [[ ${F_WIN[$i]} == "$CUR" ]]; then
      ids+=("${F_ID[$i]}")
    fi
  done
  if (( ${#ids[@]} == 0 )); then
    printf ''
    return 0
  fi
  if (( FOCUS_I < 0 )); then FOCUS_I=0; fi
  if (( FOCUS_I >= ${#ids[@]} )); then FOCUS_I=$(( ${#ids[@]} - 1 )); fi
  printf '%s' "${ids[$FOCUS_I]}"
}

move_focus() {
  local dir=$1 n=0 i
  for i in "${!F_WIN[@]}"; do
    if [[ ${F_WIN[$i]} == "$CUR" ]]; then
      n=$((n + 1))
    fi
  done
  FOCUS_I=$((FOCUS_I + dir))
  if (( FOCUS_I < 0 )); then FOCUS_I=0; fi
  if (( n > 0 && FOCUS_I >= n )); then FOCUS_I=$((n - 1)); fi
}

nav_count() {
  local id=$1 i n=0
  for i in "${!N_NAV[@]}"; do
    if [[ ${N_WIN[$i]} == "$CUR" && ${N_NAV[$i]} == "$id" ]]; then
      n=$((n + 1))
    fi
  done
  printf '%s' "$n"
}

move_vert() {
  local dir=$1 id tag sel count
  id=$(focused_id)
  [[ -z $id ]] && return 0
  tag=${Z_TAG[${ZIX[$CUR|$id]}]}
  if [[ $tag == nav ]]; then
    sel=${SEL[$CUR|$id]:-0}
    count=$(nav_count "$id")
    sel=$((sel + dir))
    if (( sel < 0 )); then sel=0; fi
    if (( count > 0 && sel >= count )); then sel=$((count - 1)); fi
    SEL[$CUR|$id]=$sel
  else
    move_focus "$dir"
  fi
}

nav_item_at() {
  local id=$1 sel=$2 i seen=0
  for i in "${!N_NAV[@]}"; do
    if [[ ${N_WIN[$i]} == "$CUR" && ${N_NAV[$i]} == "$id" ]]; then
      if (( seen == sel )); then
        printf '%s' "$i"
        return 0
      fi
      seen=$((seen + 1))
    fi
  done
  printf '%s' -1
}

nav_load_zone() {
  local target=$1 file=$2 name=$3
  PROP[$CUR|$target|file]=$file
  if [[ -n $name ]]; then
    PROP[$CUR|$target|name]=$name
    VAL[$CUR|$target]=$(read_var "$file" "$name")
  else
    VAL[$CUR|$target]=$(<"$file")
  fi
  DIRTY[$CUR|$target]=0
  refresh_zone "$target"
}

nav_open_window() {
  local href=$1 id
  id=${WIN_BY_FILE[$href]}
  STACK+=("$CUR")
  CUR=$id
  MAXW=${WIN_MAX[$CUR]}
  FOCUS_I=0
  status_msg="Prêt"
}

nav_return() {
  if (( ${#STACK[@]} == 0 )); then
    status_msg="Déjà à la première fenêtre"
    return 0
  fi
  CUR=${STACK[-1]}
  STACK=("${STACK[@]:0:${#STACK[@]}-1}")
  MAXW=${WIN_MAX[$CUR]}
  FOCUS_I=0
  status_msg="Prêt"
}

activate_nav() {
  local id=$1 sel pos
  sel=${SEL[$CUR|$id]:-0}
  pos=$(nav_item_at "$id" "$sel")
  if (( pos < 0 )); then
    return 0
  fi
  case ${N_ACTION[$pos]} in
    load) nav_load_zone "${N_TARGET[$pos]}" "${N_FILE[$pos]}" "${N_NAME[$pos]}" ;;
    window) nav_open_window "${N_HREF[$pos]}" ;;
    return) nav_return ;;
  esac
}

edit_input() {
  local id=$1 idx new
  idx=${ZIX[$CUR|$id]}
  if [[ -t 1 ]]; then
    printf '\e[?25h'
  fi
  if new=$(input_value "${PROP[$CUR|$id|name]}" "${VAL[$CUR|$id]}" "${Z_TW[$idx]}"); then
    new=${new//$'\t'/ }
    new=${new//$'\n'/ }
    new=${new//$'\r'/ }
    VAL[$CUR|$id]=$new
    if [[ ${PROP[$CUR|$id|save]} == enter ]]; then
      write_var "${PROP[$CUR|$id|file]}" "${PROP[$CUR|$id|name]}" "$new"
      DIRTY[$CUR|$id]=0
      status_msg="Enregistré"
    else
      DIRTY[$CUR|$id]=1
      status_msg="Modifié"
    fi
    refresh_zone "$id"
  else
    status_msg="Annulé"
  fi
  if [[ -t 1 ]]; then
    printf '\e[?25l'
  fi
}

edit_textarea() {
  local id=$1 idx new
  if [[ ${PROP[$CUR|$id|readonly]} == 1 ]]; then
    return 0
  fi
  idx=${ZIX[$CUR|$id]}
  if [[ -t 1 ]]; then
    printf '\e[?25h'
  fi
  if new=$(write_value "${PROP[$CUR|$id|title]:-${id}}" "${VAL[$CUR|$id]}" "${Z_TW[$idx]}" "${Z_LINES[$idx]}"); then
    VAL[$CUR|$id]=$new
    if [[ ${PROP[$CUR|$id|save]} == enter ]]; then
      printf '%s' "$new" > "${PROP[$CUR|$id|file]}"
      DIRTY[$CUR|$id]=0
      status_msg="Enregistré"
    else
      DIRTY[$CUR|$id]=1
      status_msg="Modifié"
    fi
    refresh_zone "$id"
  else
    status_msg="Annulé"
  fi
  if [[ -t 1 ]]; then
    printf '\e[?25l'
  fi
}

save_active() {
  local i id tag
  for i in "${!Z_WIN[@]}"; do
    [[ ${Z_WIN[$i]} == "$CUR" ]] || continue
    id=${Z_ID[$i]}
    [[ ${DIRTY[$CUR|$id]:-0} == 1 ]] || continue
    tag=${Z_TAG[$i]}
    if [[ $tag == input ]]; then
      write_var "${PROP[$CUR|$id|file]}" "${PROP[$CUR|$id|name]}" "${VAL[$CUR|$id]}"
    elif [[ $tag == textarea && ${PROP[$CUR|$id|readonly]} != 1 ]]; then
      printf '%s' "${VAL[$CUR|$id]}" > "${PROP[$CUR|$id|file]}"
    fi
    DIRTY[$CUR|$id]=0
  done
  status_msg="Enregistré"
}

refresh_active() {
  local i id tag
  for i in "${!Z_WIN[@]}"; do
    [[ ${Z_WIN[$i]} == "$CUR" ]] || continue
    id=${Z_ID[$i]}
    tag=${Z_TAG[$i]}
    if [[ $tag == input ]]; then
      VAL[$CUR|$id]=$(read_var "${PROP[$CUR|$id|file]}" "${PROP[$CUR|$id|name]}")
      DIRTY[$CUR|$id]=0
      refresh_zone "$id"
    elif [[ $tag == textarea ]]; then
      VAL[$CUR|$id]=$(<"${PROP[$CUR|$id|file]}")
      DIRTY[$CUR|$id]=0
      refresh_zone "$id"
    fi
  done
  status_msg="Rechargé"
}

activate() {
  local id tag
  id=$(focused_id)
  [[ -z $id ]] && return 0
  tag=${Z_TAG[${ZIX[$CUR|$id]}]}
  case $tag in
    nav) activate_nav "$id" ;;
    input) edit_input "$id" ;;
    textarea) edit_textarea "$id" ;;
    button)
      case ${PROP[$CUR|$id|action]} in
        save) save_active ;;
        refresh) refresh_active ;;
      esac
      ;;
    command) run_command "$id" ;;
  esac
}

run_command() {
  local id=$1 title run out tmp pid frame=0 key rest code
  title=${PROP[$CUR|$id|title]}
  run=${PROP[$CUR|$id|run]}
  out=${PROP[$CUR|$id|output]:-}
  tmp=$(mktemp)
  bash -c "$run" >"$tmp" 2>&1 &
  pid=$!
  while kill -0 "$pid" 2>/dev/null; do
    status_msg="${SPINNER[$frame]} En cours : ${title}"
    paint_inplace "$(frame_text)"
    frame=$(( (frame + 1) % 4 ))
    key=
    rc=0
    IFS= read -rsn1 -t 0.1 key || rc=$?
    if (( rc == 0 )) && [[ $key == $'\e' ]]; then
      kill "$pid" 2>/dev/null || true
      wait "$pid" 2>/dev/null || true
      status_msg="Annulé"
      rm -f "$tmp"
      return 0
    fi
    if (( rc != 0 && rc <= 128 )); then
      sleep 0.1
    fi
  done
  set +e
  wait "$pid"
  code=$?
  set -e
  if [[ -n $out ]]; then
    VAL[$CUR|$out]=$(<"$tmp")
    refresh_zone "$out"
  fi
  status_msg="Terminé ${code}"
  rm -f "$tmp"
}
