#!/bin/sh
# install-machine.sh - one-shot User-layer spacecraft install for a fresh machine.
#
# Clones (or updates) spacecraft into a durable path, runs `make install-global`,
# then installs companion CLIs (caveman, rtk, codegraph) with Cursor wire and
# prints a tokless-like Tools status box.
#
# Env:
#   SPACECRAFT_INSTALL_DIR  Durable clone dir (default: $HOME/.local/share/spacecraft)
#   SPACECRAFT_CLONE_SRC    Local filesystem clone source (no network). When unset,
#                           clones from SPACECRAFT_REPO / SPACECRAFT_REF.
#   SPACECRAFT_REPO         Remote URL when CLONE_SRC is unset (default: github)
#   SPACECRAFT_REF          Remote ref when CLONE_SRC is unset (default: main)
#   SPACECRAFT_CAVEMAN_INSTALL     Override caveman install.sh (skip curl|sh)
#   SPACECRAFT_RTK_INSTALL         Override rtk install.sh
#   SPACECRAFT_CODEGRAPH_INSTALL   Override codegraph install.sh
#   SPACECRAFT_COMPANION_MARKERS   Optional dir for test fixture markers
set -e
# Prefer pipefail so curl|sh does not mask curl failures (dash may lack it).
(set -o pipefail) 2>/dev/null && set -o pipefail

# Keep HOME absolute before make -C / companion `cd` (relative HOME nests under clone).
case "$HOME" in
  /*) ;;
  *) HOME="$(pwd)/$HOME"; export HOME ;;
esac

# Absolutize PATH entries so fixture fake-bin survives companion `cd`.
_abs_path=
IFS=:
for _p in ${PATH:-/usr/bin:/bin}; do
  [ -n "$_p" ] || continue
  case "$_p" in
    /*) _a=$_p ;;
    *) _a="$(pwd)/$_p" ;;
  esac
  if [ -n "$_abs_path" ]; then
    _abs_path="$_abs_path:$_a"
  else
    _abs_path="$_a"
  fi
done
unset IFS _p _a
PATH="$_abs_path"
export PATH
unset _abs_path

INSTALL_DIR="${SPACECRAFT_INSTALL_DIR:-$HOME/.local/share/spacecraft}"
REPO_URL="${SPACECRAFT_REPO:-https://github.com/xiivthx/spacecraft.git}"
REPO_REF="${SPACECRAFT_REF:-main}"

# Keep paths absolute: companion install `cd`s away from the caller cwd.
case "$INSTALL_DIR" in
  /*) ;;
  *) INSTALL_DIR="$(pwd)/$INSTALL_DIR" ;;
esac
if [ -n "${SPACECRAFT_CLONE_SRC:-}" ]; then
  case "$SPACECRAFT_CLONE_SRC" in
    /*) ;;
    *) SPACECRAFT_CLONE_SRC="$(pwd)/$SPACECRAFT_CLONE_SRC" ;;
  esac
fi
if [ -n "${SPACECRAFT_CAVEMAN_INSTALL:-}" ]; then
  case "$SPACECRAFT_CAVEMAN_INSTALL" in
    /*) ;;
    *) SPACECRAFT_CAVEMAN_INSTALL="$(pwd)/$SPACECRAFT_CAVEMAN_INSTALL" ;;
  esac
fi
if [ -n "${SPACECRAFT_RTK_INSTALL:-}" ]; then
  case "$SPACECRAFT_RTK_INSTALL" in
    /*) ;;
    *) SPACECRAFT_RTK_INSTALL="$(pwd)/$SPACECRAFT_RTK_INSTALL" ;;
  esac
fi
if [ -n "${SPACECRAFT_CODEGRAPH_INSTALL:-}" ]; then
  case "$SPACECRAFT_CODEGRAPH_INSTALL" in
    /*) ;;
    *) SPACECRAFT_CODEGRAPH_INSTALL="$(pwd)/$SPACECRAFT_CODEGRAPH_INSTALL" ;;
  esac
fi
if [ -n "${SPACECRAFT_COMPANION_MARKERS:-}" ]; then
  case "$SPACECRAFT_COMPANION_MARKERS" in
    /*) ;;
    *) SPACECRAFT_COMPANION_MARKERS="$(pwd)/$SPACECRAFT_COMPANION_MARKERS" ;;
  esac
  export SPACECRAFT_COMPANION_MARKERS
fi

CAVEMAN_INSTALL_URL="https://raw.githubusercontent.com/JuliusBrussee/caveman/main/install.sh"
RTK_INSTALL_URL="https://raw.githubusercontent.com/rtk-ai/rtk/refs/heads/master/install.sh"
CODEGRAPH_INSTALL_URL="https://raw.githubusercontent.com/colbymchenry/codegraph/main/install.sh"

# Non-interactive by default (never prompt). --agents / --yes are explicit flags.
agents=cursor
non_interactive=yes
flagged=0

while [ $# -gt 0 ]; do
  case "$1" in
    --agents)
      if [ $# -lt 2 ]; then
        echo "error: --agents requires a value" >&2
        exit 1
      fi
      agents=$2
      flagged=1
      shift 2
      ;;
    --yes)
      non_interactive=yes
      flagged=1
      shift
      ;;
    *)
      echo "error: unknown argument: $1" >&2
      exit 1
      ;;
  esac
done

case "$agents" in
  cursor) ;;
  *)
    echo "error: unsupported --agents value: $agents (supported: cursor)" >&2
    exit 1
    ;;
esac

if [ "$flagged" -eq 1 ]; then
  echo "agents: $agents"
  echo "non-interactive: $non_interactive"
fi

mkdir -p "$(dirname -- "$INSTALL_DIR")"

if [ ! -d "$INSTALL_DIR/.git" ]; then
  if [ -e "$INSTALL_DIR" ]; then
    echo "error: $INSTALL_DIR exists but is not a git clone" >&2
    exit 1
  fi
  if [ -n "${SPACECRAFT_CLONE_SRC:-}" ]; then
    echo "Cloning local source: $SPACECRAFT_CLONE_SRC -> $INSTALL_DIR"
    git clone "$SPACECRAFT_CLONE_SRC" "$INSTALL_DIR"
  else
    echo "Cloning $REPO_URL ($REPO_REF) -> $INSTALL_DIR"
    git clone --depth 1 --branch "$REPO_REF" "$REPO_URL" "$INSTALL_DIR" \
      || git clone --depth 1 "$REPO_URL" "$INSTALL_DIR"
  fi
else
  echo "Updating existing clone: $INSTALL_DIR"
  git -C "$INSTALL_DIR" pull --ff-only
fi

echo "Running User-layer install from $INSTALL_DIR"
make -C "$INSTALL_DIR" install-global

# Prefer freshly installed User-layer binaries.
export PATH="$HOME/.local/bin:${PATH:-/usr/bin:/bin}"

# --- companions (soft-fail: spacecraft layer already done) ---
# Official installers (caveman skills, etc.) may write into cwd. Leave the
# caller's working tree alone by running companions under a durable cache dir.
companion_cwd="$HOME/.cache/spacecraft-machine-install"
mkdir -p "$companion_cwd"
cd "$companion_cwd"

_tool_version() {
  cmd=$1
  if ! command -v "$cmd" >/dev/null 2>&1; then
    return 1
  fi
  ver=$("$cmd" --version 2>/dev/null | sed -n '1p') || ver=""
  if [ -z "$ver" ]; then
    ver=$("$cmd" version 2>/dev/null | sed -n '1p') || ver=""
  fi
  if [ -z "$ver" ]; then
    return 1
  fi
  printf '%s' "$ver"
}

_spacecraft_version() {
  if [ -d "$INSTALL_DIR/.git" ]; then
    ver=$(git -C "$INSTALL_DIR" describe --tags --always 2>/dev/null) || ver=""
    if [ -n "$ver" ]; then
      printf '%s' "$ver"
      return 0
    fi
  fi
  printf '%s' "dev"
}

# Run override script or official curl|sh.
# $1=label $2=override path $3=url $4=shell (sh|bash)
# Remaining args are forwarded only on the official path (`shell -s -- ...`).
# Override path always runs `sh "$override"` with no extra flags (fixtures).
_run_installer() {
  label=$1
  override=$2
  url=$3
  shell=${4:-sh}
  shift 4
  if [ -n "$override" ]; then
    echo "Installing $label (override: $override)"
    sh "$override"
  elif [ "$#" -gt 0 ]; then
    echo "Installing $label (official installer)"
    curl -fsSL "$url" | "$shell" -s -- "$@"
  else
    echo "Installing $label (official installer)"
    curl -fsSL "$url" | "$shell"
  fi
}

caveman_state=fail
caveman_pct=0
caveman_ver="—"
rtk_state=fail
rtk_pct=0
rtk_ver="—"
codegraph_state=fail
codegraph_pct=0
codegraph_ver="—"

# Official caveman is Cursor-scoped; override fixtures get no extra flags.
if _run_installer "caveman" "${SPACECRAFT_CAVEMAN_INSTALL:-}" "$CAVEMAN_INSTALL_URL" bash --only cursor; then
  # Caveman is primarily a skill installer; CLI may be absent after success.
  if ver=$(_tool_version caveman); then
    caveman_ver=$ver
  else
    caveman_ver=skill
  fi
  caveman_state=ok
  caveman_pct=100
else
  echo "warning: caveman install failed (continuing)" >&2
fi

if _run_installer "rtk" "${SPACECRAFT_RTK_INSTALL:-}" "$RTK_INSTALL_URL" sh; then
  if rtk init -g --agent cursor; then
    if ver=$(_tool_version rtk); then
      rtk_state=ok
      rtk_pct=100
      rtk_ver=$ver
    fi
  else
    echo "warning: rtk Cursor wire failed (continuing)" >&2
  fi
else
  echo "warning: rtk install failed (continuing)" >&2
fi

if _run_installer "codegraph" "${SPACECRAFT_CODEGRAPH_INSTALL:-}" "$CODEGRAPH_INSTALL_URL" sh; then
  if codegraph install --target=cursor --yes; then
    if ver=$(_tool_version codegraph); then
      codegraph_state=ok
      codegraph_pct=100
      codegraph_ver=$ver
    fi
  else
    echo "warning: codegraph Cursor wire failed (continuing)" >&2
  fi
else
  echo "warning: codegraph install failed (continuing)" >&2
fi

# shellcheck disable=SC1091
. "$INSTALL_DIR/scripts/tools-status.sh"

sc_ver=$(_spacecraft_version)
row_sc=$(tools_status_row "Spacecraft" ok 100 "$sc_ver")
row_cv=$(tools_status_row "Caveman" "$caveman_state" "$caveman_pct" "$caveman_ver")
row_rtk=$(tools_status_row "RTK" "$rtk_state" "$rtk_pct" "$rtk_ver")
row_cg=$(tools_status_row "CodeGraph" "$codegraph_state" "$codegraph_pct" "$codegraph_ver")
tools_status_box "$row_sc" "$row_cv" "$row_rtk" "$row_cg"

echo ""
echo "Next steps:"
echo "  - Ensure ~/.local/bin is on your PATH (CLIs land there)."
echo "  - Paste ~/.cursor/spacecraft/USER-RULES.txt into Cursor Settings > Rules > User Rules."
echo "  - Restart Cursor to pick up skills, hooks, and companion wiring."
echo "  - Per project: run codegraph init in the repo you want indexed."
