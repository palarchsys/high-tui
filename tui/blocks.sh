#!/usr/bin/env bash
# Cached Gum blocks. A key reprints them.
# Shellcheck: this file is sourced.

declare -gA BLOCK=()

zone_payload() {
  local id=$1 tag
  tag=${Z_TAG[${ZIX[$CUR|$id]}]}
  case $tag in
    nav|select) printf '%s' "${SEL[$CUR|$id]:-0}" ;;
    input) printf '%s\x1f%s' "${PROP[$CUR|$id|name]:-}" "${VAL[$CUR|$id]//$'\n'/$'\x1e'}" ;;
    textarea) printf '%s' "${VAL[$CUR|$id]//$'\n'/$'\x1e'}" ;;
    status) printf '%s' "$status_msg" ;;
    *) printf '%s' "${PROP[$CUR|$id|label]:-}" ;;
  esac
}

block_key() {
  local id=$1 focus=$2
  printf '%s\x1d%s\x1d%s' "$CUR|$id|$focus" "$(zone_payload "$id")" "${SCROLL[$CUR|$id]:-0}"
}

logical_text() {
  local id=$1 tag sel i count pos
  tag=${Z_TAG[${ZIX[$CUR|$id]}]}
  case $tag in
    header|button|command) printf '%s' "${PROP[$CUR|$id|label]:-}" ;;
    status) printf '%s' "$status_msg" ;;
    input) printf '%s\n%s' "${PROP[$CUR|$id|name]}" "${VAL[$CUR|$id]:-}" ;;
    textarea) printf '%s' "${VAL[$CUR|$id]:-}" ;;
    nav|select)
      sel=${SEL[$CUR|$id]:-0}
      count=$(nav_count "$id")
      printf '%s\n' "${PROP[$CUR|$id|label]:-}"
      for ((i = 0; i < count; i++)); do
        pos=$(nav_item_at "$id" "$i")
        if (( i == sel )); then
          printf '› %s\n' "${N_TITLE[$pos]}"
        else
          printf '  %s\n' "${N_TITLE[$pos]}"
        fi
      done
      ;;
  esac
}

wrap_one() {
  local text=$1 width=$2 chunk cut
  text=${text//$'\t'/ }
  text=${text//$'\r'/ }
  if (( width < 1 )); then width=1; fi
  while (( ${#text} > width )); do
    if [[ ${text:width:1} == ' ' ]]; then
      printf '%s\n' "${text:0:width}"
      text=${text:width+1}
      continue
    fi
    chunk=${text:0:width}
    if [[ $chunk == *' '* ]]; then
      cut=${chunk% *}
      if (( ${#cut} > 0 )); then
        printf '%s\n' "$cut"
        text=${text:${#cut}+1}
        continue
      fi
    fi
    printf '%s\n' "$chunk"
    text=${text:width}
  done
  printf '%s\n' "$text"
}

fill_wrapped() {
  local id=$1 width=$2
  local -n _rows=$3
  local -n _from=$4
  local raw line piece src=0
  _rows=()
  _from=()
  raw=$(logical_text "$id")
  if [[ -z $raw ]]; then
    _rows+=("")
    _from+=(0)
    return 0
  fi
  while IFS= read -r line || [[ -n ${line:-} ]]; do
    while IFS= read -r piece || [[ -n ${piece:-} ]]; do
      _rows+=("$piece")
      _from+=("$src")
    done < <(wrap_one "$line" "$width")
    src=$((src + 1))
  done <<< "$raw"
}

wrapped_rows() {
  local id=$1
  local -n _rows=$2
  local idx=${ZIX[$CUR|$id]}
  local text_w=${Z_TW[$idx]} lines=${Z_LINES[$idx]}
  fill_wrapped "$id" "$text_w" "$2" "$3"
  if (( ${#_rows[@]} > lines && text_w >= 2 )); then
    fill_wrapped "$id" $((text_w - 1)) "$2" "$3"
  fi
}

wrapped_total() {
  local id=$1
  local -a rows=() from=()
  wrapped_rows "$id" rows from
  printf '%s' "${#rows[@]}"
}

scroll_max() {
  local id=$1 total lines max
  total=$(wrapped_total "$id")
  lines=${Z_LINES[${ZIX[$CUR|$id]}]}
  max=$((total - lines))
  if (( max < 0 )); then max=0; fi
  printf '%s' "$max"
}

view_scroll() {
  local id=$1 total=$2 lines=$3
  local -n _from=$4
  local tag sel target=0 end=0 span=0 s=0 i
  tag=${Z_TAG[${ZIX[$CUR|$id]}]}
  if [[ $tag == nav || $tag == select ]]; then
    sel=${SEL[$CUR|$id]:-0}
    for i in "${!_from[@]}"; do
      if (( _from[i] == sel + 1 )); then
        target=$i
        break
      fi
    done
    end=$target
    while (( end + 1 < total && _from[end + 1] == sel + 1 )); do
      end=$((end + 1))
    done
    span=$((end - target + 1))
    if (( span >= lines )); then
      s=$target
    elif (( end >= lines )); then
      s=$((end - lines + 1))
    fi
  else
    s=${SCROLL[$CUR|$id]:-0}
  fi
  if (( s > total - lines )); then s=$((total - lines)); fi
  if (( s < 0 )); then s=0; fi
  if [[ $tag != nav && $tag != select ]]; then
    SCROLL[$CUR|$id]=$s
  fi
  printf '%s' "$s"
}

bar_glyph() {
  local row=$1 visible=$2 total=$3 scroll=$4
  local thumb_h thumb_top max_scroll
  thumb_h=$((visible * visible / total))
  if (( thumb_h < 1 )); then thumb_h=1; fi
  if (( thumb_h > visible )); then thumb_h=$visible; fi
  max_scroll=$((total - visible))
  if (( max_scroll <= 0 )); then
    thumb_top=0
  else
    thumb_top=$((scroll * (visible - thumb_h) / max_scroll))
  fi
  if (( row >= thumb_top && row < thumb_top + thumb_h )); then
    printf '█'
  else
    printf '│'
  fi
}

compose() {
  local id=$1
  local idx=${ZIX[$CUR|$id]}
  local lines=${Z_LINES[$idx]} text_w=${Z_TW[$idx]} content_h=${Z_CH[$idx]}
  local content_w=$((Z_OW[$idx] - 2 * Z_BW[$idx]))
  local pad_t=${Z_PT[$idx]} pad_b=${Z_PB[$idx]} pad_l=${Z_PL[$idx]} pad_r=${Z_PR[$idx]}
  local lh=${Z_LH[$idx]}
  local -a rows=() from=() logical=() out=()
  local width=$text_w bar=0 total scroll=0 i src piece glyph
  local line n blank left right body joined="" sep=""
  fill_wrapped "$id" "$width" rows from
  total=${#rows[@]}
  if (( total > lines && text_w >= 2 )); then
    width=$((text_w - 1))
    fill_wrapped "$id" "$width" rows from
    total=${#rows[@]}
    bar=1
  fi
  scroll=$(view_scroll "$id" "$total" "$lines" from)
  for ((i = 0; i < lines; i++)); do
    src=$((scroll + i))
    piece=""
    if (( src < total )); then
      piece=${rows[$src]}
    fi
    piece=$(clip_pad "$piece" "$width")
    if (( bar == 1 )); then
      glyph=$(bar_glyph "$i" "$lines" "$total" "$scroll")
      piece="${piece}${glyph}"
    fi
    logical+=("$piece")
  done
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
  key=$(block_key "$id" "$focus")
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

warm_view() {
  local id=$1 with_focus=$2 max s saved
  saved=${SCROLL[$CUR|$id]:-0}
  max=$(scroll_max "$id")
  for ((s = 0; s <= max; s++)); do
    SCROLL[$CUR|$id]=$s
    cache_put "$id" 0
    if (( with_focus == 1 )); then
      cache_put "$id" 1
    fi
  done
  if (( saved > max )); then saved=$max; fi
  SCROLL[$CUR|$id]=$saved
}

refresh_zone() {
  local id=$1 tag focusable=0
  tag=${Z_TAG[${ZIX[$CUR|$id]}]}
  if [[ $tag == nav || $tag == select || $tag == input || $tag == button || $tag == command || $tag == textarea ]]; then
    focusable=1
  fi
  if [[ $tag == input || $tag == textarea ]]; then
    warm_view "$id" "$focusable"
  else
    cache_put "$id" 0
    if (( focusable == 1 )); then
      cache_put "$id" 1
    fi
  fi
}

zone_view() {
  local id=$1 focus=0 tag key
  tag=${Z_TAG[${ZIX[$CUR|$id]}]}
  if [[ $(focused_id) == "$id" ]]; then
    focus=1
  fi
  if [[ $tag == header || $tag == status ]]; then
    focus=0
  fi
  key=$(block_key "$id" "$focus")
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
      nav|select)
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
        local saved_name=${PROP[$win|$id|name]}
        warm_view "$id" 1
        saved=${VAL[$win|$id]}
        for k in "${!N_NAV[@]}"; do
          [[ ${N_WIN[$k]} == "$win" ]] || continue
          if [[ ${N_FACT[$k]} == load && ${N_FTGT[$k]} == "$id" && -n ${N_FNAME[$k]} ]]; then
            PROP[$win|$id|name]=${N_FNAME[$k]}
            VAL[$win|$id]=$(read_var "${N_FFILE[$k]}" "${N_FNAME[$k]}")
            warm_view "$id" 1
          fi
          if [[ ${N_EACT[$k]} == load && ${N_ETGT[$k]} == "$id" && -n ${N_ENAME[$k]} ]]; then
            PROP[$win|$id|name]=${N_ENAME[$k]}
            VAL[$win|$id]=$(read_var "${N_EFILE[$k]}" "${N_ENAME[$k]}")
            warm_view "$id" 1
          fi
          if [[ ${N_EACT[$k]} == set && ${N_ETGT[$k]} == "$id" && -n ${N_EVALUE[$k]} ]]; then
            PROP[$win|$id|name]=${N_ENAME[$k]}
            VAL[$win|$id]=${N_EVALUE[$k]}
            warm_view "$id" 1
          fi
        done
        VAL[$win|$id]=$saved
        PROP[$win|$id|name]=$saved_name
        ;;
      textarea)
        local editable=0
        if [[ ${PROP[$win|$id|readonly]} != 1 ]]; then editable=1; fi
        warm_view "$id" "$editable"
        saved=${VAL[$win|$id]}
        for k in "${!N_NAV[@]}"; do
          [[ ${N_WIN[$k]} == "$win" ]] || continue
          if [[ ${N_FACT[$k]} == load && ${N_FTGT[$k]} == "$id" && -z ${N_FNAME[$k]} && -n ${N_FFILE[$k]} ]]; then
            VAL[$win|$id]=$(<"${N_FFILE[$k]}")
            warm_view "$id" "$editable"
          fi
          if [[ ${N_EACT[$k]} == load && ${N_ETGT[$k]} == "$id" && -z ${N_ENAME[$k]} && -n ${N_EFILE[$k]} ]]; then
            VAL[$win|$id]=$(<"${N_EFILE[$k]}")
            warm_view "$id" "$editable"
          fi
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
        titles+=("${PROP[$win|$id|label]:-}")
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

warm_blocks() {
  local win
  for win in "${WIN_ID[@]}"; do
    warm_window "$win"
  done
  CUR=${WIN_ID[0]}
  FOCUS_I=0
  status_msg="Prêt"
  FOCUS_MARK=
  apply_focus
}

frame_text() {
  local -a blocks=() rows=() parts=()
  local i id row y h gap prev=0 part blanks
  declare -A seen=()
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
