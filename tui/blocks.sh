#!/usr/bin/env bash
# Cached Gum blocks. A key reprints them.
# Shellcheck: this file is sourced.

declare -gA BLOCK=()

zone_payload() {
  local id=$1 tag
  tag=${Z_TAG[${ZIX[$CUR|$id]}]}
  case $tag in
    nav) printf '%s' "${SEL[$CUR|$id]:-0}" ;;
    input|textarea) printf '%s' "${VAL[$CUR|$id]//$'\n'/$'\x1e'}" ;;
    status) printf '%s' "$status_msg" ;;
    *) printf '%s%s' "${PROP[$CUR|$id|text]:-}" "${PROP[$CUR|$id|title]:-}" ;;
  esac
}

logical_text() {
  local id=$1 tag lines sel slots start i shown count pos
  tag=${Z_TAG[${ZIX[$CUR|$id]}]}
  lines=${Z_LINES[${ZIX[$CUR|$id]}]}
  case $tag in
    header|button) printf '%s' "${PROP[$CUR|$id|text]:-}" ;;
    command) printf '%s' "${PROP[$CUR|$id|title]:-}" ;;
    status) printf '%s' "$status_msg" ;;
    input) printf '%s\n%s' "${PROP[$CUR|$id|name]}" "${VAL[$CUR|$id]:-}" ;;
    textarea) printf '%s' "${VAL[$CUR|$id]:-}" ;;
    nav)
      sel=${SEL[$CUR|$id]:-0}
      count=$(nav_count "$id")
      slots=$((lines - 1))
      if (( slots < 1 )); then slots=1; fi
      start=0
      if (( sel >= slots )); then start=$((sel - slots + 1)); fi
      printf '%s' "${PROP[$CUR|$id|title]}"
      shown=0
      for ((i = start; i < count && shown < slots; i++)); do
        pos=$(nav_item_at "$id" "$i")
        if (( i == sel )); then
          printf '\n› %s' "${N_TITLE[$pos]}"
        else
          printf '\n  %s' "${N_TITLE[$pos]}"
        fi
        shown=$((shown + 1))
      done
      ;;
  esac
}

compose() {
  local id=$1
  local idx=${ZIX[$CUR|$id]}
  local lines=${Z_LINES[$idx]} text_w=${Z_TW[$idx]} content_h=${Z_CH[$idx]}
  local content_w=$((Z_OW[$idx] - 2 * Z_BW[$idx]))
  local pad_t=${Z_PT[$idx]} pad_b=${Z_PB[$idx]} pad_l=${Z_PL[$idx]} pad_r=${Z_PR[$idx]}
  local lh=${Z_LH[$idx]}
  local -a logical=() out=()
  local raw line i n blank left right body joined="" sep=""
  raw=$(logical_text "$id")
  if [[ -n $raw ]]; then
    while IFS= read -r line || [[ -n ${line:-} ]]; do
      logical+=("$(clip_pad "$line" "$text_w")")
    done <<< "$raw"
  fi
  while (( ${#logical[@]} < lines )); do
    logical+=("$(clip_pad "" "$text_w")")
  done
  if (( ${#logical[@]} > lines )); then
    logical=("${logical[@]:0:lines}")
  fi
  printf -v blank '%*s' "$content_w" ""
  for ((i = 0; i < pad_t; i++)); do out+=("$blank"); done
  for i in "${!logical[@]}"; do
    printf -v left '%*s' "$pad_l" ""
    printf -v right '%*s' "$pad_r" ""
    body=${logical[$i]}
    out+=("${left}${body}${right}")
    if (( i + 1 < ${#logical[@]} )); then
      for ((n = 1; n < lh; n++)); do out+=("$blank"); done
    fi
  done
  for ((i = 0; i < pad_b; i++)); do out+=("$blank"); done
  while (( ${#out[@]} < content_h )); do out+=("$blank"); done
  if (( ${#out[@]} > content_h )); then
    out=("${out[@]:0:content_h}")
  fi
  for line in "${out[@]}"; do
    joined+="${sep}${line}"
    sep=$'\n'
  done
  printf '%s' "$joined"
}

cache_put() {
  local id=$1 focus=$2
  local idx=${ZIX[$CUR|$id]}
  local fg bg bc bold payload key text block cw
  payload=$(zone_payload "$id")
  key="$CUR|$id|$focus|$payload"
  if [[ -n ${BLOCK[$key]+x} ]]; then
    return 0
  fi
  if (( focus == 1 )); then
    fg=${Z_FFG[$idx]}; bg=${Z_FBG[$idx]}; bc=${Z_FBC[$idx]}; bold=${Z_FBOLD[$idx]}
  else
    fg=${Z_FG[$idx]}; bg=${Z_BG[$idx]}; bc=${Z_BC[$idx]}; bold=${Z_BOLD[$idx]}
  fi
  text=$(compose "$id")
  cw=$((Z_OW[$idx] - 2 * Z_BW[$idx]))
  block=$(style_box "$fg" "$bg" "$bc" "${Z_BW[$idx]}" "$bold" "$cw" "$text")
  block=${block%$'\n'}
  BLOCK[$key]=$block
}

refresh_zone() {
  local id=$1 tag
  tag=${Z_TAG[${ZIX[$CUR|$id]}]}
  cache_put "$id" 0
  if [[ $tag == nav || $tag == input || $tag == button || $tag == command || ( $tag == textarea && ${PROP[$CUR|$id|readonly]} != 1 ) ]]; then
    cache_put "$id" 1
  fi
}

zone_view() {
  local id=$1 focus=0 tag payload key
  tag=${Z_TAG[${ZIX[$CUR|$id]}]}
  if [[ $(focused_id) == "$id" ]]; then
    focus=1
  fi
  if [[ $tag == header || $tag == status || ( $tag == textarea && ${PROP[$CUR|$id|readonly]} == 1 ) ]]; then
    focus=0
  fi
  payload=$(zone_payload "$id")
  key="$CUR|$id|$focus|$payload"
  if [[ -z ${BLOCK[$key]+x} ]]; then
    cache_put "$id" "$focus"
  fi
  printf '%s' "${BLOCK[$key]}"
}

blank_block() {
  local n=$1 w=$2 i s out="" sep=""
  printf -v s '%*s' "$w" ""
  for ((i = 0; i < n; i++)); do
    out+="${sep}${s}"
    sep=$'\n'
  done
  printf '%s' "$out"
}

warm_window() {
  local win=$1 save_cur=$CUR save_msg=$status_msg
  local i id tag j k item target name file saved
  local -a msgs titles=()
  CUR=$win
  for i in "${!Z_WIN[@]}"; do
    [[ ${Z_WIN[$i]} == "$win" ]] || continue
    id=${Z_ID[$i]}
    tag=${Z_TAG[$i]}
    case $tag in
      nav)
        local count
        count=$(nav_count "$id")
        for ((j = 0; j < count; j++)); do
          SEL["$win|$id"]=$j
          cache_put "$id" 0
          cache_put "$id" 1
        done
        SEL["$win|$id"]=0
        ;;
      input)
        cache_put "$id" 0
        cache_put "$id" 1
        saved=${VAL[$win|$id]}
        for k in "${!N_NAV[@]}"; do
          [[ ${N_WIN[$k]} == "$win" && ${N_ACTION[$k]} == load && ${N_TARGET[$k]} == "$id" && -n ${N_NAME[$k]} ]] || continue
          VAL[$win|$id]=$(read_var "${N_FILE[$k]}" "${N_NAME[$k]}")
          cache_put "$id" 0
          cache_put "$id" 1
        done
        VAL[$win|$id]=$saved
        ;;
      textarea)
        cache_put "$id" 0
        saved=${VAL[$win|$id]}
        for k in "${!N_NAV[@]}"; do
          [[ ${N_WIN[$k]} == "$win" && ${N_ACTION[$k]} == load && ${N_TARGET[$k]} == "$id" && -z ${N_NAME[$k]} ]] || continue
          VAL[$win|$id]=$(<"${N_FILE[$k]}")
          cache_put "$id" 0
        done
        VAL[$win|$id]=$saved
        ;;
      button)
        cache_put "$id" 0
        cache_put "$id" 1
        ;;
      command)
        cache_put "$id" 0
        cache_put "$id" 1
        titles+=("${PROP[$win|$id|title]}")
        ;;
      header)
        cache_put "$id" 0
        ;;
    esac
  done
  msgs=("Prêt" "Déjà à la première fenêtre" "Enregistré" "Rechargé" "Annulé" "Modifié" "Terminé 0" "Terminé 1")
  for title in "${titles[@]+"${titles[@]}"}"; do
    local s
    for s in "${SPINNER[@]}"; do
      msgs+=("$s En cours : ${title}")
    done
  done
  for i in "${!Z_WIN[@]}"; do
    [[ ${Z_WIN[$i]} == "$win" && ${Z_TAG[$i]} == status ]] || continue
    id=${Z_ID[$i]}
    for status_msg in "${msgs[@]}"; do
      cache_put "$id" 0
    done
  done
  status_msg=$save_msg
  CUR=$save_cur
}

warm_narrow() {
  local save_cur=$CUR id="" saved key payload i
  CUR=${WIN_ID[0]}
  for i in "${!Z_WIN[@]}"; do
    if [[ ${Z_WIN[$i]} == "$CUR" && ${Z_TAG[$i]} == header ]]; then
      id=${Z_ID[$i]}
      break
    fi
  done
  saved=${PROP[$CUR|$id|text]:-}
  PROP[$CUR|$id|text]="Terminal étroit"
  cache_put "$id" 0
  payload=$(zone_payload "$id")
  BLOCK[narrow]=${BLOCK[$CUR|$id|0|$payload]}
  PROP[$CUR|$id|text]=$saved
  CUR=$save_cur
}

warm_blocks() {
  local win
  for win in "${WIN_ID[@]}"; do
    warm_window "$win"
  done
  warm_narrow
  CUR=${WIN_ID[0]}
  FOCUS_I=0
  status_msg="Prêt"
}

frame_text() {
  local -a blocks=() rows=() parts=()
  local i id row y h gap prev=0 part blanks
  declare -A seen=()
  if [[ -n ${COLUMNS:-} && ${COLUMNS} -lt ${WIN_MAX[$CUR]} ]]; then
    blocks+=("${BLOCK[narrow]}")
  fi
  for i in "${!Z_WIN[@]}"; do
    [[ ${Z_WIN[$i]} == "$CUR" ]] || continue
    row=${Z_ROW[$i]}
    if [[ -z ${seen[$row]+x} ]]; then
      seen[$row]=1
      rows+=("$row")
    fi
  done
  for row in "${rows[@]}"; do
    y=0; h=0
    parts=()
    for i in "${!Z_WIN[@]}"; do
      [[ ${Z_WIN[$i]} == "$CUR" && ${Z_ROW[$i]} == "$row" ]] || continue
      id=${Z_ID[$i]}
      y=${Z_Y[$i]}
      h=${Z_OH[$i]}
      parts+=("$(zone_view "$id")")
      gap=${Z_GAP[$i]}
      if (( gap > 0 )); then
        parts+=("$(blank_block "$h" "$gap")")
      fi
    done
    blanks=$((y - prev))
    if (( blanks > 0 )); then
      blocks+=("$(blank_block "$blanks" "${WIN_MAX[$CUR]}")")
    fi
    blocks+=("$(paste_horizontal "${parts[@]}")")
    prev=$((y + h))
  done
  blanks=$(( ${WIN_H[$CUR]} - prev ))
  if (( blanks > 0 )); then
    blocks+=("$(blank_block "$blanks" "${WIN_MAX[$CUR]}")")
  fi
  stack_vertical "${blocks[@]}"
}
