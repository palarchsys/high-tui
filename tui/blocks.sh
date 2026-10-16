#!/usr/bin/env bash
# One associative cache. One column renderer. Warm every variant once.
# Shellcheck: this file is sourced.

declare -gA BLOCK=()
declare -gA ZONE_REF=()
declare -gA ZONE_W=()
declare -gA ZONE_H=()
STEADY_STATUS=(Prêt Menu Ouvrir Retour Échap)
STATUS_TEMPLATE_TEXT=""
STATUS_TEMPLATE_BLOCK=""

column_key() {
  local focused=$1
  local text=$2
  text=${text//$'\n'/$'\x1e'}
  printf 'column:%s:%s' "$focused" "$text"
}

detail_key() {
  printf 'detail:%s:%s' "$1" "$2"
}

status_key() {
  printf 'status:%s' "$1"
}

visible_cols() {
  local s=$1
  local out="" ch i n=${#s} in_csi=0
  for ((i = 0; i < n; i++)); do
    ch=${s:i:1}
    if [[ $ch == $'\e' ]]; then
      in_csi=1
      continue
    fi
    if (( in_csi )); then
      if [[ $ch =~ [[:alpha:]] ]]; then
        in_csi=0
      fi
      continue
    fi
    out+=$ch
  done
  printf '%s' "${#out}"
}

last_visible_index() {
  local s=$1
  local i n=${#s} ch in_csi=0 last=-1
  for ((i = 0; i < n; i++)); do
    ch=${s:i:1}
    if [[ $ch == $'\e' ]]; then
      in_csi=1
      continue
    fi
    if (( in_csi )); then
      if [[ $ch =~ [[:alpha:]] ]]; then
        in_csi=0
      fi
      continue
    fi
    last=$i
  done
  printf '%s' "$last"
}

drop_visible() {
  local prefix=$1
  local drop=$2
  local i n=${#prefix} ch in_csi=0
  local -a chars=() visible_idx=()
  for ((i = 0; i < n; i++)); do
    chars+=("${prefix:i:1}")
  done
  for ((i = 0; i < n; i++)); do
    ch=${chars[$i]}
    if [[ $ch == $'\e' ]]; then
      in_csi=1
      continue
    fi
    if (( in_csi )); then
      if [[ $ch =~ [[:alpha:]] ]]; then
        in_csi=0
      fi
      continue
    fi
    visible_idx+=("$i")
  done
  local remain=$drop vi idx
  for ((vi = ${#visible_idx[@]} - 1; vi >= 0 && remain > 0; vi--)); do
    idx=${visible_idx[$vi]}
    if [[ ${chars[$idx]} == " " ]]; then
      chars[idx]=$'\x01'
      remain=$((remain - 1))
    fi
  done
  for ((vi = ${#visible_idx[@]} - 1; vi >= 0 && remain > 0; vi--)); do
    idx=${visible_idx[$vi]}
    if [[ ${chars[idx]} != $'\x01' ]]; then
      chars[idx]=$'\x01'
      remain=$((remain - 1))
    fi
  done
  local out=""
  for ((i = 0; i < n; i++)); do
    if [[ ${chars[$i]} != $'\x01' ]]; then
      out+=${chars[$i]}
    fi
  done
  printf '%s' "$out"
}

# Spaces go in before the right border, and extra rows before the bottom border.
fit_width() {
  local line=$1
  local want=$2
  local have idx pad spaces prefix suffix
  have=$(visible_cols "$line")
  if (( have == want || want <= 0 )); then
    printf '%s' "$line"
    return 0
  fi
  idx=$(last_visible_index "$line")
  if (( idx < 0 )); then
    clip_pad "$line" "$want"
    return 0
  fi
  if (( have < want )); then
    pad=$((want - have))
    printf -v spaces '%*s' "$pad" ""
    printf '%s' "${line:0:idx}${spaces}${line:idx}"
    return 0
  fi
  prefix=${line:0:idx}
  suffix=${line:idx}
  prefix=$(drop_visible "$prefix" $((have - want)))
  printf '%s' "${prefix}${suffix}"
}

zone_set() {
  local role=$1
  local block=$2
  local -a lines=()
  ZONE_REF[$role]=$block
  split_lines "$block" lines
  ZONE_H[$role]=${#lines[@]}
  if (( ${#lines[@]} > 0 )); then
    ZONE_W[$role]=$(visible_cols "${lines[0]}")
  else
    ZONE_W[$role]=0
  fi
}

init_zones() {
  zone_set column "$(style_panel 0 " ")"
  zone_set banner "$(style_banner " ")"
  zone_set status "$(style_status " ")"
  zone_set detail "$(table_print </dev/null)"
}

lock_zone() {
  local role=$1
  local block=$2
  local -a lines=() ref=() kept=()
  local want_h=${ZONE_H[$role]}
  local want_w=${ZONE_W[$role]}
  local i w h same=1 line pad_line out="" sep=""
  split_lines "$block" lines
  h=${#lines[@]}
  if (( h != want_h )); then
    same=0
  else
    for line in "${lines[@]}"; do
      w=$(visible_cols "$line")
      if (( w != want_w )); then
        same=0
        break
      fi
    done
  fi
  if (( same == 1 )); then
    printf '%s' "$block"
    return 0
  fi
  if (( h > want_h )); then
    for ((i = 0; i < want_h - 1; i++)); do
      kept+=("${lines[$i]}")
    done
    kept+=("${lines[$((h - 1))]}")
    lines=("${kept[@]}")
  elif (( h < want_h )); then
    split_lines "${ZONE_REF[$role]}" ref
    pad_line=${ref[$(( ${#ref[@]} - 2 ))]}
    for ((i = 0; i < h - 1; i++)); do
      kept+=("${lines[$i]}")
    done
    while (( ${#kept[@]} + 1 < want_h )); do
      kept+=("$pad_line")
    done
    if (( h > 0 )); then
      kept+=("${lines[$((h - 1))]}")
    fi
    lines=("${kept[@]}")
  fi
  for i in "${!lines[@]}"; do
    lines[i]=$(fit_width "${lines[$i]}" "$want_w")
  done
  for line in "${lines[@]}"; do
    out+="${sep}${line}"
    sep=$'\n'
  done
  printf '%s' "$out"
}

render_column() {
  local title=$1
  local selected=$2
  local focused=$3
  shift 3
  local text key
  text=$(column_text "$title" "$selected" "$@")
  key=$(column_key "$focused" "$text")
  if [[ -z ${BLOCK[$key]+x} ]]; then
    BLOCK[$key]=$(lock_zone column "$(style_panel "$focused" "$text")")
  fi
}

column_view() {
  local title=$1
  local selected=$2
  local focused=$3
  shift 3
  local text key
  text=$(column_text "$title" "$selected" "$@")
  key=$(column_key "$focused" "$text")
  printf '%s' "${BLOCK[$key]}"
}

refresh_detail() {
  local rub=$1
  local fic=$2
  local key
  key=$(detail_key "$rub" "$fic")
  BLOCK[$key]=$(lock_zone detail "$(detail_tsv "$rub" "$fic" | table_print)")
}

derive_status() {
  local new_text=$1
  local -a lines=()
  local mid before after target
  local delta best_start best_len run_start run_len i c n new_len spaces
  if [[ $new_text == "$STATUS_TEMPLATE_TEXT" ]]; then
    printf '%s' "$STATUS_TEMPLATE_BLOCK"
    return 0
  fi
  split_lines "$STATUS_TEMPLATE_BLOCK" lines
  mid=${lines[1]}
  delta=$((${#new_text} - ${#STATUS_TEMPLATE_TEXT}))
  before=${mid%%"$STATUS_TEMPLATE_TEXT"*}
  after=${mid#*"$STATUS_TEMPLATE_TEXT"}
  target=$after
  best_start=0
  best_len=0
  run_start=0
  run_len=0
  n=${#target}
  for ((i = 0; i < n; i++)); do
    c=${target:i:1}
    if [[ $c == " " ]]; then
      if (( run_len == 0 )); then
        run_start=$i
      fi
      run_len=$((run_len + 1))
      if (( run_len > best_len )); then
        best_len=$run_len
        best_start=$run_start
      fi
    else
      run_len=0
    fi
  done
  new_len=$((best_len - delta))
  if (( best_len == 0 || new_len < 0 )); then
    style_status "$new_text"
    return 0
  fi
  printf -v spaces '%*s' "$new_len" ""
  after="${target:0:best_start}${spaces}${target:best_start+best_len}"
  lines[1]="${before}${new_text}${after}"
  for i in "${!lines[@]}"; do
    if (( i > 0 )); then
      printf '\n'
    fi
    printf '%s' "${lines[$i]}"
  done
}

cache_status() {
  local text=$1
  local key
  key=$(status_key "$text")
  if [[ -n ${BLOCK[$key]+x} ]]; then
    return 0
  fi
  BLOCK[$key]=$(lock_zone status "$(derive_status "$text")")
}

warm_list() {
  local title=$1
  shift
  local -a items=("$@")
  local i on
  for i in "${!items[@]}"; do
    for on in 0 1; do
      render_column "$title" "$i" "$on" "${items[@]}"
    done
  done
}

warm_columns() {
  local save_ri=$ri save_fi=$fi save_ci=$ci
  local r f
  local -a rubs=() fics=() champs=()
  rebuild
  rubs=("${RUBS[@]}")
  warm_list Rubriques "${rubs[@]}"
  for r in "${!rubs[@]}"; do
    ri=$r
    fi=0
    ci=0
    rebuild
    fics=("${FICS[@]}")
    warm_list Fiches "${fics[@]}"
    for f in "${!fics[@]}"; do
      ri=$r
      fi=$f
      ci=0
      rebuild
      champs=("${CHAMPS[@]}")
      warm_list Champs "${champs[@]}"
    done
  done
  warm_list Menu "${MENU[@]}"
  ri=$save_ri
  fi=$save_fi
  ci=$save_ci
  rebuild
}

warm_details() {
  local -A seen=()
  local i rub fic mark
  for i in "${!RUB[@]}"; do
    rub=${RUB[$i]}
    fic=${FIC[$i]}
    mark=$'\t'"${rub}"$'\t'"${fic}"
    if [[ -z ${seen[$mark]+x} ]]; then
      seen[$mark]=1
      refresh_detail "$rub" "$fic"
    fi
  done
}

warm_status() {
  local save_ri=$ri save_fi=$fi save_ci=$ci save_mi=$mi
  local save_focus=$focus save_msg=$status_msg
  local -a messages=("${STEADY_STATUS[@]}" Annulé)
  local -A seen=()
  local -a rubs=() fics=() champs=()
  local name text r f c m msg fic_n champ_n rub fic
  # Edit outcomes stay in the cache, so a later key does not call gum.
  for name in "${CHAMP[@]}"; do
    if [[ -z ${seen[$name]+x} ]]; then
      seen[$name]=1
      messages+=("Modifié : ${name}")
    fi
  done

  rebuild
  text=$(status_text "${RUBS[$ri]}" "${FICS[$fi]}" "${CHAMPS[$ci]}")
  STATUS_TEMPLATE_TEXT=$text
  STATUS_TEMPLATE_BLOCK=$(lock_zone status "$(style_status "$text")")
  BLOCK[$(status_key "$text")]=$STATUS_TEMPLATE_BLOCK

  rubs=("${RUBS[@]}")
  for r in "${!rubs[@]}"; do
    ri=$r
    fi=0
    ci=0
    rebuild
    fic_n=${#FICS[@]}
    for ((f = 0; f < fic_n; f++)); do
      ri=$r
      fi=$f
      ci=0
      rebuild
      fics=("${FICS[@]}")
      champs=("${CHAMPS[@]}")
      champ_n=${#champs[@]}
      rub=${RUBS[$ri]}
      fic=${fics[$f]}
      for ((c = 0; c < champ_n; c++)); do
        for focus in 0 1 2; do
          for msg in "${messages[@]}"; do
            status_msg=$msg
            text=$(status_text "$rub" "$fic" "${champs[$c]}")
            cache_status "$text"
          done
        done
      done
    done
  done

  for ((m = 0; m < ${#MENU[@]}; m++)); do
    mi=$m
    focus=3
    for msg in "${messages[@]}"; do
      status_msg=$msg
      text=$(status_text - - -)
      cache_status "$text"
    done
  done

  ri=$save_ri
  fi=$save_fi
  ci=$save_ci
  mi=$save_mi
  focus=$save_focus
  status_msg=$save_msg
  rebuild
}

warm_blocks() {
  rebuild
  init_zones
  BLOCK[header]=$(lock_zone banner "$(style_banner "high-tui")")
  BLOCK[narrow]=$(lock_zone banner "$(style_banner "Terminal étroit")")
  warm_columns
  warm_details
  warm_status
}

frame_text() {
  local rub fic champ text
  local on0=0 on1=0 on2=0 on3=0
  local c0 c1 c2 c3 cols header detail status dkey skey
  rebuild
  rub=${RUBS[$ri]}
  fic=${FICS[$fi]}
  champ=${CHAMPS[$ci]}
  case $focus in
    0) on0=1 ;;
    1) on1=1 ;;
    2) on2=1 ;;
    3) on3=1 ;;
  esac
  c0=$(column_view Rubriques "$ri" "$on0" "${RUBS[@]}")
  c1=$(column_view Fiches "$fi" "$on1" "${FICS[@]}")
  c2=$(column_view Champs "$ci" "$on2" "${CHAMPS[@]}")
  c3=$(column_view Menu "$mi" "$on3" "${MENU[@]}")
  cols=$(paste_horizontal "$c0" "$c1" "$c2" "$c3")
  header=${BLOCK[header]}
  dkey=$(detail_key "$rub" "$fic")
  detail=${BLOCK[$dkey]}
  text=$(status_text "$rub" "$fic" "$champ")
  skey=$(status_key "$text")
  status=${BLOCK[$skey]}
  if [[ -n ${COLUMNS:-} && ${COLUMNS} -lt 80 ]]; then
    stack_vertical "${BLOCK[narrow]}" "$header" "$cols" "$detail" "$status"
  else
    stack_vertical "$header" "$cols" "$detail" "$status"
  fi
}
