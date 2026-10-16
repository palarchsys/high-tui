#!/usr/bin/env bash
# TSV load, cursor, list text, detail rows, status text, in-memory edit.
# Shellcheck: this file is sourced.

DATA=data.tsv
RUB=()
FIC=()
CHAMP=()
VAL=()
RUBS=()
FICS=()
CHAMPS=()
CIDX=()
MENU=(Ouvrir Modifier Retour)
FOCUS_NAME=(Rubriques Fiches Champs Menu)
ri=0
fi=0
ci=0
mi=0
focus=0
origin=0
status_msg="Prêt"

load_data() {
  RUB=()
  FIC=()
  CHAMP=()
  VAL=()
  local rub fiche champ valeur
  while IFS=$'\t' read -r rub fiche champ valeur || [[ -n ${rub:-} ]]; do
    [[ -z ${rub:-} ]] && continue
    [[ $rub == rubrique ]] && continue
    RUB+=("$rub")
    FIC+=("$fiche")
    CHAMP+=("$champ")
    VAL+=("$valeur")
  done < "$DATA"
}

rebuild() {
  local -A seen_rub=()
  local -A seen_fic=()
  local i name rub fic
  RUBS=()
  for i in "${!RUB[@]}"; do
    name=${RUB[$i]}
    if [[ -z ${seen_rub[$name]+x} ]]; then
      seen_rub[$name]=1
      RUBS+=("$name")
    fi
  done
  if (( ri >= ${#RUBS[@]} )); then
    ri=$((${#RUBS[@]} - 1))
  fi
  if (( ri < 0 )); then
    ri=0
  fi

  FICS=()
  rub=${RUBS[$ri]}
  for i in "${!RUB[@]}"; do
    if [[ ${RUB[$i]} == "$rub" ]]; then
      name=${FIC[$i]}
      if [[ -z ${seen_fic[$name]+x} ]]; then
        seen_fic[$name]=1
        FICS+=("$name")
      fi
    fi
  done
  if (( fi >= ${#FICS[@]} )); then
    fi=$((${#FICS[@]} - 1))
  fi
  if (( fi < 0 )); then
    fi=0
  fi

  CHAMPS=()
  CIDX=()
  fic=${FICS[$fi]}
  for i in "${!RUB[@]}"; do
    if [[ ${RUB[$i]} == "$rub" && ${FIC[$i]} == "$fic" ]]; then
      CHAMPS+=("${CHAMP[$i]}")
      CIDX+=("$i")
    fi
  done
  if (( ci >= ${#CHAMPS[@]} )); then
    ci=$((${#CHAMPS[@]} - 1))
  fi
  if (( ci < 0 )); then
    ci=0
  fi
  if (( mi < 0 )); then
    mi=0
  fi
  if (( mi >= ${#MENU[@]} )); then
    mi=$((${#MENU[@]} - 1))
  fi
}

column_text() {
  local title=$1
  local selected=$2
  shift 2
  local -a items=("$@")
  local text=$title
  local i n
  for i in "${!items[@]}"; do
    if [[ $i -eq $selected ]]; then
      text+=$'\n'"› ${items[$i]}"
    else
      text+=$'\n'"  ${items[$i]}"
    fi
  done
  n=${#items[@]}
  while [[ $n -lt $ROW_SLOTS ]]; do
    text+=$'\n'" "
    n=$((n + 1))
  done
  printf '%s' "$text"
}

detail_tsv() {
  local rub=$1
  local fic=$2
  local i
  for i in "${!RUB[@]}"; do
    if [[ ${RUB[$i]} == "$rub" && ${FIC[$i]} == "$fic" ]]; then
      printf '%s\t%s\n' "${CHAMP[$i]}" "${VAL[$i]}"
    fi
  done
}

status_text() {
  local rub=$1
  local fic=$2
  local champ=$3
  local where=${FOCUS_NAME[$focus]}
  local extra note
  if [[ $focus -eq 3 ]]; then
    extra=${MENU[$mi]}
  else
    extra="${rub} / ${fic} / ${champ}"
  fi
  note=$status_msg
  if (( ${#note} > 40 )); then
    note=${note:0:39}
  fi
  printf '%s · %s · %s · ↑↓ ←→ entrée q' "$where" "$extra" "$note"
}

move_up() {
  case $focus in
    0)
      if (( ri > 0 )); then
        ri=$((ri - 1))
        fi=0
        ci=0
      fi
      ;;
    1)
      if (( fi > 0 )); then
        fi=$((fi - 1))
        ci=0
      fi
      ;;
    2)
      if (( ci > 0 )); then
        ci=$((ci - 1))
      fi
      ;;
    3)
      if (( mi > 0 )); then
        mi=$((mi - 1))
      fi
      ;;
  esac
}

move_down() {
  rebuild
  case $focus in
    0)
      if (( ri + 1 < ${#RUBS[@]} )); then
        ri=$((ri + 1))
        fi=0
        ci=0
      fi
      ;;
    1)
      if (( fi + 1 < ${#FICS[@]} )); then
        fi=$((fi + 1))
        ci=0
      fi
      ;;
    2)
      if (( ci + 1 < ${#CHAMPS[@]} )); then
        ci=$((ci + 1))
      fi
      ;;
    3)
      if (( mi + 1 < ${#MENU[@]} )); then
        mi=$((mi + 1))
      fi
      ;;
  esac
}

move_right() {
  if (( focus < 3 )); then
    focus=$((focus + 1))
  fi
}

move_left() {
  if (( focus > 0 )); then
    focus=$((focus - 1))
  fi
}

ouvrir() {
  if (( origin < 2 )); then
    focus=$((origin + 1))
  else
    focus=3
  fi
  status_msg="Ouvrir"
}

edit_champ() {
  local idx name old value rub fic
  rebuild
  idx=${CIDX[$ci]}
  name=${CHAMP[$idx]}
  old=${VAL[$idx]}
  rub=${RUB[$idx]}
  fic=${FIC[$idx]}
  if [[ -t 1 ]]; then
    printf '\e[?25h'
  fi
  if value=$(input_value "$name" "$old"); then
    value=${value//$'\t'/ }
    value=${value//$'\r'/ }
    value=${value//$'\n'/ }
    VAL[$idx]=$value
    status_msg="Modifié : ${name}"
    refresh_detail "$rub" "$fic"
  else
    status_msg="Annulé"
  fi
  if [[ -t 1 ]]; then
    printf '\e[?25l'
  fi
}

on_enter() {
  if (( focus != 3 )); then
    origin=$focus
    focus=3
    status_msg="Menu"
    return
  fi
  rebuild
  case ${MENU[$mi]} in
    Ouvrir) ouvrir ;;
    Modifier) edit_champ ;;
    Retour)
      focus=2
      status_msg="Retour"
      ;;
  esac
}

on_escape() {
  if (( focus == 3 )); then
    focus=2
    status_msg="Échap"
  fi
}
