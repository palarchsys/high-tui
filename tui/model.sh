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
N_FACT=() N_FTGT=() N_FFILE=() N_FNAME=() N_FHREF=() N_FRUN=() N_FOUT=()
N_EACT=() N_ETGT=() N_EFILE=() N_ENAME=() N_EHREF=() N_ERUN=() N_EOUT=() N_EVALUE=()
F_WIN=() F_ID=()
declare -gA ZIX=() PROP=() WIN_MAX=() WIN_TITLE=() WIN_FILE=()
declare -gA WIN_BG=() WIN_FG=() WIN_H=() WIN_BY_FILE=()
declare -gA SEL=() VAL=() DIRTY=() SCROLL=()
CUR=
FOCUS_I=0
STACK=()
status_msg="Prêt"
FG_UI=
ACCENT_UI=
MAXW=0
FOCUS_MARK=
RETURNING=0
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
        N_WIN+=("$w"); N_NAV+=("$id"); N_INDEX+=("$row"); N_ID+=("$idx"); N_TITLE+=("$tag")
        N_FACT+=("$x"); N_FTGT+=("$y"); N_FFILE+=("$ow"); N_FNAME+=("$oh"); N_FHREF+=("$tw")
        N_FRUN+=("$lines"); N_FOUT+=("$ch")
        N_EACT+=("$fg"); N_ETGT+=("$bg"); N_EFILE+=("$bc"); N_ENAME+=("$bw"); N_EHREF+=("$bold")
        N_ERUN+=("$ffg"); N_EOUT+=("$fbg"); N_EVALUE+=("$fbc")
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
    if [[ ${Z_WIN[$i]} == "$CUR" && ( ${Z_TAG[$i]} == nav || ${Z_TAG[$i]} == select ) ]]; then
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
    if [[ $tag == nav || $tag == select ]]; then
      SEL["$w|$id"]=0
    fi
    DIRTY["$w|$id"]=0
    SCROLL["$w|$id"]=0
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
  apply_focus
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
  if [[ $tag == nav || $tag == select ]]; then
    sel=${SEL[$CUR|$id]:-0}
    count=$(nav_count "$id")
    sel=$((sel + dir))
    if (( sel < 0 )); then sel=0; fi
    if (( count > 0 && sel >= count )); then sel=$((count - 1)); fi
    SEL[$CUR|$id]=$sel
    apply_focus
  elif [[ $tag == textarea || $tag == input ]]; then
    local max s
    max=$(scroll_max "$id")
    s=${SCROLL[$CUR|$id]:-0}
    if (( max > 0 && dir < 0 && s > 0 )); then
      SCROLL[$CUR|$id]=$((s - 1))
    elif (( max > 0 && dir > 0 && s < max )); then
      SCROLL[$CUR|$id]=$((s + 1))
    else
      move_focus "$dir"
    fi
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
  SCROLL[$CUR|$target]=0
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
  apply_focus
}

nav_return() {
  if (( RETURNING == 1 )); then
    return 0
  fi
  RETURNING=1
  if (( ${#STACK[@]} == 0 )); then
    status_msg="Déjà à la première fenêtre"
    RETURNING=0
    return 0
  fi
  CUR=${STACK[-1]}
  STACK=("${STACK[@]:0:${#STACK[@]}-1}")
  MAXW=${WIN_MAX[$CUR]}
  FOCUS_I=0
  status_msg="Prêt"
  apply_focus
  RETURNING=0
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
    SCROLL[$CUR|$id]=0
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
  if new=$(write_value "${PROP[$CUR|$id|label]:-${id}}" "${VAL[$CUR|$id]}" "${Z_TW[$idx]}" "${Z_LINES[$idx]}"); then
    VAL[$CUR|$id]=$new
    if [[ ${PROP[$CUR|$id|save]} == enter ]]; then
      printf '%s' "$new" > "${PROP[$CUR|$id|file]}"
      DIRTY[$CUR|$id]=0
      status_msg="Enregistré"
    else
      DIRTY[$CUR|$id]=1
      status_msg="Modifié"
    fi
    SCROLL[$CUR|$id]=0
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
      SCROLL[$CUR|$id]=0
      refresh_zone "$id"
    elif [[ $tag == textarea ]]; then
      VAL[$CUR|$id]=$(<"${PROP[$CUR|$id|file]}")
      DIRTY[$CUR|$id]=0
      SCROLL[$CUR|$id]=0
      refresh_zone "$id"
    fi
  done
  status_msg="Rechargé"
}

edit_target() {
  local id=$1 tag
  [[ -z $id ]] && return 0
  if [[ -z ${ZIX[$CUR|$id]+x} ]]; then
    status_msg="Cible absente"
    return 0
  fi
  tag=${Z_TAG[${ZIX[$CUR|$id]}]}
  case $tag in
    input) edit_input "$id" ;;
    textarea) edit_textarea "$id" ;;
  esac
}

apply_set() {
  local target=$1 file=$2 name=$3 value=$4 tag
  write_var "$file" "$name" "$value"
  if [[ -n $target && -n ${ZIX[$CUR|$target]+x} ]]; then
    tag=${Z_TAG[${ZIX[$CUR|$target]}]}
    PROP[$CUR|$target|file]=$file
    if [[ $tag == input ]]; then
      PROP[$CUR|$target|name]=$name
      VAL[$CUR|$target]=$value
    elif [[ $tag == textarea ]]; then
      VAL[$CUR|$target]=$(<"$file")
    fi
    DIRTY[$CUR|$target]=0
    SCROLL[$CUR|$target]=0
    refresh_zone "$target"
  fi
  status_msg="Enregistré"
}

dispatch_action() {
  local zone=$1 action=$2 target=$3 file=$4 name=$5 href=$6 run=$7 output=$8 value=${9:-}
  case $action in
    "") return 0 ;;
    load) nav_load_zone "$target" "$file" "$name" ;;
    edit)
      if [[ -z $target ]]; then target=$zone; fi
      edit_target "$target"
      ;;
    window) nav_open_window "$href" ;;
    return) nav_return ;;
    save) save_active ;;
    refresh) refresh_active ;;
    set) apply_set "$target" "$file" "$name" "$value" ;;
    run)
      if [[ -n $run ]]; then PROP[$CUR|$zone|run]=$run; fi
      if [[ -n $output ]]; then PROP[$CUR|$zone|output]=$output; fi
      run_command "$zone"
      ;;
  esac
}

apply_focus() {
  local id tag sel pos
  id=$(focused_id)
  [[ -z $id ]] && return 0
  tag=${Z_TAG[${ZIX[$CUR|$id]}]}
  if [[ $tag == nav || $tag == select ]]; then
    sel=${SEL[$CUR|$id]:-0}
    if [[ $FOCUS_MARK == "$CUR|$id|$sel" ]]; then
      return 0
    fi
    FOCUS_MARK="$CUR|$id|$sel"
    pos=$(nav_item_at "$id" "$sel")
    if (( pos < 0 )); then
      return 0
    fi
    dispatch_action "$id" \
      "${N_FACT[$pos]}" "${N_FTGT[$pos]}" "${N_FFILE[$pos]}" "${N_FNAME[$pos]}" \
      "${N_FHREF[$pos]}" "${N_FRUN[$pos]}" "${N_FOUT[$pos]}" ""
  else
    if [[ $FOCUS_MARK == "$CUR|$id" ]]; then
      return 0
    fi
    FOCUS_MARK="$CUR|$id"
    dispatch_action "$id" \
      "${PROP[$CUR|$id|focus.action]:-}" \
      "${PROP[$CUR|$id|focus.target]:-}" \
      "${PROP[$CUR|$id|focus.file]:-}" \
      "${PROP[$CUR|$id|focus.name]:-}" \
      "${PROP[$CUR|$id|focus.href]:-}" \
      "${PROP[$CUR|$id|focus.run]:-}" \
      "${PROP[$CUR|$id|focus.output]:-}" \
      "${PROP[$CUR|$id|focus.value]:-}"
  fi
}

activate() {
  local id tag pos
  id=$(focused_id)
  [[ -z $id ]] && return 0
  tag=${Z_TAG[${ZIX[$CUR|$id]}]}
  if [[ $tag == nav || $tag == select ]]; then
    pos=$(nav_item_at "$id" "${SEL[$CUR|$id]:-0}")
    if (( pos < 0 )); then
      return 0
    fi
    dispatch_action "$id" \
      "${N_EACT[$pos]}" "${N_ETGT[$pos]}" "${N_EFILE[$pos]}" "${N_ENAME[$pos]}" \
      "${N_EHREF[$pos]}" "${N_ERUN[$pos]}" "${N_EOUT[$pos]}" "${N_EVALUE[$pos]}"
  else
    dispatch_action "$id" \
      "${PROP[$CUR|$id|enter.action]:-}" \
      "${PROP[$CUR|$id|enter.target]:-}" \
      "${PROP[$CUR|$id|enter.file]:-}" \
      "${PROP[$CUR|$id|enter.name]:-}" \
      "${PROP[$CUR|$id|enter.href]:-}" \
      "${PROP[$CUR|$id|enter.run]:-}" \
      "${PROP[$CUR|$id|enter.output]:-}" \
      "${PROP[$CUR|$id|enter.value]:-}"
  fi
}

run_command() {
  local id=$1 title run out tmp pid frame=0 key rest code
  title=${PROP[$CUR|$id|label]:-}
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
