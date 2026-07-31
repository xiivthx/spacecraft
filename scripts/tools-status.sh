#!/bin/sh
# tools-status.sh - tokless-like Tools status line helpers.
# Sourceable: . scripts/tools-status.sh
#
# tools_status_row <name> <ok|fail> <percent> <version>
# tools_status_box <row> [row...]

# Visible length after stripping CSI sequences (for box width).
_tools_status_visible_len() {
  esc=$(printf '\033')
  plain=$(printf '%s' "$1" | sed "s/${esc}\\[[0-9;]*[a-zA-Z]//g")
  printf '%s' "$plain" | wc -m | tr -d ' '
}

tools_status_row() {
  name=$1
  state=$2
  percent=$3
  version=$4

  green=$(printf '\033[32m')
  reset=$(printf '\033[0m')

  case "$state" in
    ok) mark='✔' ;;
    fail) mark='✘' ;;
    *) mark='?' ;;
  esac

  # 20-cell bar filled by percent (0-100).
  width=20
  filled=$((percent * width / 100))
  if [ "$filled" -gt "$width" ]; then
    filled=$width
  fi
  empty=$((width - filled))
  bar='['
  i=0
  while [ "$i" -lt "$filled" ]; do
    bar="${bar}█"
    i=$((i + 1))
  done
  i=0
  while [ "$i" -lt "$empty" ]; do
    bar="${bar} "
    i=$((i + 1))
  done
  bar="${bar}]"

  printf '%s %s %s%s%s %s%% %s\n' "$mark" "$name" "$green" "$bar" "$reset" "$percent" "$version"
}

# Wrap one or more status rows under a box-drawing frame labeled Tools.
tools_status_box() {
  if [ "$#" -lt 1 ]; then
    echo "tools_status_box: need at least one row" >&2
    return 1
  fi

  # Inner width: max visible row length, at least enough for "─ Tools ".
  inner=8
  for row in "$@"; do
    # Drop trailing newline from tools_status_row capture if present.
    row=$(printf '%s' "$row" | tr -d '\n')
    len=$(_tools_status_visible_len "$row")
    if [ "$len" -gt "$inner" ]; then
      inner=$len
    fi
  done

  # Top: ┌─ Tools ─…─┐  ("─ Tools " is 8 cells; pad remaining with ─).
  pad=$((inner - 8))
  top='┌─ Tools '
  i=0
  while [ "$i" -lt "$pad" ]; do
    top="${top}─"
    i=$((i + 1))
  done
  top="${top}┐"
  printf '%s\n' "$top"

  for row in "$@"; do
    row=$(printf '%s' "$row" | tr -d '\n')
    len=$(_tools_status_visible_len "$row")
    spaces=$((inner - len))
    pad_str=
    i=0
    while [ "$i" -lt "$spaces" ]; do
      pad_str="${pad_str} "
      i=$((i + 1))
    done
    printf '│%s%s│\n' "$row" "$pad_str"
  done

  # Bottom: └─…─┘
  bottom='└'
  i=0
  while [ "$i" -lt "$inner" ]; do
    bottom="${bottom}─"
    i=$((i + 1))
  done
  bottom="${bottom}┘"
  printf '%s\n' "$bottom"
}
