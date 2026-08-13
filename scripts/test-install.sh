#!/bin/sh
# test-install.sh - install/bootstrap smoke test in a throwaway temp dir.
#
# Seeds a pre-existing unrelated hook so the hooks merge must preserve user
# hooks alongside project session-start (safety hooks stay User-layer only),
# then runs install-global against a fake HOME and checks the core skills land.
# Never writes ~/.cursorrules.
#
# Usage: test-install.sh <repo-root> <spacecraft-binary>
set -e

ROOT="${1:?usage: test-install.sh <repo-root> <spacecraft-binary>}"
BIN="${2:?usage: test-install.sh <repo-root> <spacecraft-binary>}"

mkdir -p "$ROOT/.tmp"
tmp=$(mktemp -d "$ROOT/.tmp/install-smoke.XXXXXX")
fake_home="$tmp/home"
# Always succeed so EXIT trap cannot flip a green run to status 1.
cleanup() { rm -rf "$tmp"; }
trap cleanup EXIT INT TERM

echo "Install smoke in $tmp"
mkdir -p "$tmp/.cursor" "$fake_home"
printf '%s\n' '{"version":1,"hooks":{"beforeShellExecution":[{"command":".cursor/hooks/user-unrelated.sh","matcher":"echo"}]}}' \
  > "$tmp/.cursor/hooks.json"

HOME="$fake_home" sh "$ROOT/scripts/install-cursor.sh" "$tmp" "$ROOT"
sh "$ROOT/scripts/smoke.sh" "$tmp" "$BIN"

# First .space create via install-cursor: git init (if needed) + .space/ in .gitignore.
if [ ! -d "$tmp/.git" ] && ! git -C "$tmp" rev-parse --git-dir >/dev/null 2>&1; then
  echo "FAIL: install into empty project did not initialize git (missing .git)"
  exit 1
fi
if [ ! -f "$tmp/.gitignore" ] || ! grep -Eq '^[[:space:]]*\.space/?[[:space:]]*$' "$tmp/.gitignore"; then
  echo "FAIL: install into empty project .gitignore missing .space/ entry"
  exit 1
fi
echo "  ok   install-project ensures git + .space/ in .gitignore"

hooks="$tmp/.cursor/hooks.json"
if ! grep -q 'user-unrelated.sh' "$hooks"; then
  echo "FAIL: install clobbered pre-existing hooks (user-unrelated.sh missing)"
  exit 1
fi
if ! grep -q 'session-start.sh' "$hooks"; then
  echo "FAIL: install-project missing session-start hook in $hooks"
  exit 1
fi
if grep -q 'check-ship-commands.sh' "$hooks"; then
  echo "FAIL: install-project merged User-layer safety hook check-ship-commands.sh into project $hooks"
  exit 1
fi
if grep -q 'check-main-write.sh' "$hooks"; then
  echo "FAIL: install-project merged User-layer safety hook check-main-write.sh into project $hooks"
  exit 1
fi
echo "  ok   hooks merge preserves user + session-start; omits safety hooks"

if [ -f "$fake_home/.cursorrules" ]; then
  echo "FAIL: install wrote legacy ~/.cursorrules"
  exit 1
fi
echo "  ok   no legacy ~/.cursorrules"

# T4/T5: install-project places domain rules 300-620 + the session-start
# hook into a project target distinct from the source repo. Runs before any
# User-layer install-global / --full, so project domain does not depend on it.
for rule in 300-security 400-performance 500-database 600-firmware \
  610-firmware-peripherals 620-firmware-testing; do
  test -f "$tmp/.cursor/rules/$rule.mdc" \
    || { echo "FAIL: install-project missing domain rule $rule.mdc in $tmp/.cursor/rules"; exit 1; }
done
grep -q 'session-start.sh' "$hooks" \
  || { echo "FAIL: install-project did not wire session-start hook into $hooks"; exit 1; }
test -f "$tmp/.cursor/hooks/session-start.sh" \
  || { echo "FAIL: install-project did not copy session-start.sh into $tmp/.cursor/hooks/"; exit 1; }
echo "  ok   install-project places domain rules 300-620 + session-start hook"

# Project-layer must not copy User-layer safety hook scripts (global owns them).
for safety_hook in check-main-write.sh check-ship-commands.sh; do
  if [ -f "$tmp/.cursor/hooks/$safety_hook" ]; then
    echo "FAIL: install-project copied User-layer safety hook $safety_hook into $tmp/.cursor/hooks/"
    exit 1
  fi
done
test -f "$tmp/.cursor/hooks/session-start.sh" \
  || { echo "FAIL: install-project missing session-start.sh under $tmp/.cursor/hooks/ after safety omit"; exit 1; }
echo "  ok   install-project omits safety hook scripts; keeps session-start.sh"

# Project-layer must not install agents (User-layer / global owns sc-*.md agents).
if [ -f "$tmp/.cursor/agents/sc-coder.md" ]; then
  echo "FAIL: install-project copied agent sc-coder.md into $tmp/.cursor/agents (User-layer only)"
  exit 1
fi
if [ -d "$tmp/.cursor/agents" ]; then
  leftover_agents=$(find "$tmp/.cursor/agents" -maxdepth 1 -type f -name 'sc-*.md' 2>/dev/null | head -n 1)
  if [ -n "$leftover_agents" ]; then
    echo "FAIL: install-project left spacecraft agent under $tmp/.cursor/agents: $leftover_agents"
    exit 1
  fi
fi
echo "  ok   install-project omits agents (no sc-*.md under project .cursor/agents)"

# T5: project-layer install also lands domain encyclopedia skills (lean User
# layer omits these; project bootstrap must not require SPACECRAFT_SKILL_PROFILE=full).
# Domain packs stay project-local; lean-core lifecycle skills stay User-layer only.
project_skills="$tmp/.cursor/skills"
project_domain_skills="sc-solid sc-security sc-performance sc-web-backend sc-web-frontend sc-database sc-firmware sc-browser-probe"
for skill in $project_domain_skills; do
  test -f "$project_skills/$skill/SKILL.md" \
    || { echo "FAIL: install-project missing domain skill $skill under $project_skills (project layer must include domain packs without User --full)"; exit 1; }
done
echo "  ok   install-project places domain encyclopedia skills (no User --full required)"

# T4: install-project omits User-layer basenames via explicit list (aligned with
# gen-user-rules SOURCES: 000/026/027/050/100/200) — not alwaysApply: true only.
for rule in 000-spacecraft 026-intent-coach 027-th-en-hil 050-style 100-conventions 200-workflow; do
  if [ -f "$tmp/.cursor/rules/$rule.mdc" ]; then
    echo "FAIL: install-project copied User-layer rule $rule.mdc into project target $tmp (must exclude by basename list)"
    exit 1
  fi
done
# install-cursor.sh must name that exclude (USER_LAYER or the basename list).
if ! grep -Eq 'USER_LAYER|000-spacecraft' "$ROOT/scripts/install-cursor.sh"; then
  echo "FAIL: install-cursor.sh missing explicit User-layer basename exclude (USER_LAYER or 000-spacecraft)"
  exit 1
fi
echo "  ok   install-project excludes User-layer basenames from project target"

# Lean-core skills (must stay in sync with scripts/global-sync.sh LEAN_SKILLS) stay
# User-layer only (~/.cursor/skills via lean install-global). Project install must
# omit them — and install-cursor.sh must name the exclude like USER_LAYER for rules.
project_lean_skills="sc-discuss sc-run sc-ship sc-quick sc-mission sc-planning sc-tdd sc-verification sc-judge sc-clarify sc-git sc-search sc-storm sc-writer"
for skill in $project_lean_skills; do
  if [ -e "$project_skills/$skill" ]; then
    echo "FAIL: install-project copied lean-core skill $skill into $project_skills (User-layer only)"
    exit 1
  fi
done
# Spot-check nested lean references that used to be required under project.
if [ -e "$project_skills/sc-discuss/references/lens-pass.md" ]; then
  echo "FAIL: install-project left sc-discuss/references/lens-pass.md under project (lean-core)"
  exit 1
fi
if [ -e "$project_skills/sc-judge/references/judge-break" ]; then
  echo "FAIL: install-project left sc-judge judge-break fixtures under project (lean-core)"
  exit 1
fi
if ! grep -Eq 'LEAN_SKILLS' "$ROOT/scripts/install-cursor.sh"; then
  echo "FAIL: install-cursor.sh missing explicit lean-skill exclude (LEAN_SKILLS)"
  exit 1
fi
echo "  ok   install-project omits lean-core skills (User-layer only)"

# Re-running project install prunes lean-core skills previously copied into the
# project skills dir (mirrors lean install-global pruning domain packs).
# Also prunes leftover spacecraft agents + User-layer safety hook scripts.
mkdir -p "$project_skills/sc-run" "$tmp/.cursor/agents" "$tmp/.cursor/hooks"
printf '%s\n' '# seeded-for-project-lean-prune' > "$project_skills/sc-run/SKILL.md"
printf '%s\n' '# seeded-for-project-agent-prune' > "$tmp/.cursor/agents/sc-coder.md"
printf '%s\n' '#!/bin/sh' > "$tmp/.cursor/hooks/check-ship-commands.sh"
printf '%s\n' '#!/bin/sh' > "$tmp/.cursor/hooks/check-main-write.sh"
HOME="$fake_home" sh "$ROOT/scripts/install-cursor.sh" "$tmp" "$ROOT"
if [ -e "$project_skills/sc-run" ]; then
  echo "FAIL: re-run install-project left lean-core skill sc-run under $project_skills (expected prune)"
  exit 1
fi
# Domain packs must still be present after prune re-run.
test -f "$project_skills/sc-browser-probe/SKILL.md" \
  || { echo "FAIL: after lean prune re-run, missing domain skill sc-browser-probe under $project_skills"; exit 1; }
echo "  ok   re-run install-project prunes lean-core skills; keeps domain packs"

if [ -e "$tmp/.cursor/agents/sc-coder.md" ]; then
  echo "FAIL: re-run install-project left seeded agent sc-coder.md under $tmp/.cursor/agents (expected prune)"
  exit 1
fi
echo "  ok   re-run install-project prunes leftover spacecraft agents"

for safety_hook in check-main-write.sh check-ship-commands.sh; do
  if [ -f "$tmp/.cursor/hooks/$safety_hook" ]; then
    echo "FAIL: re-run install-project left seeded safety hook $safety_hook under $tmp/.cursor/hooks (expected prune)"
    exit 1
  fi
done
test -f "$tmp/.cursor/hooks/session-start.sh" \
  || { echo "FAIL: after safety-hook prune re-run, missing session-start.sh under $tmp/.cursor/hooks/"; exit 1; }
echo "  ok   re-run install-project prunes safety hook scripts; keeps session-start.sh"

mkdir -p "$fake_home/.cursor"
printf '%s\n' '{"version":1,"hooks":{"beforeShellExecution":[{"command":"~/.cursor/hooks/user-unrelated-global.sh","matcher":"echo"}]}}' \
  > "$fake_home/.cursor/hooks.json"

HOME="$fake_home" make -C "$ROOT" install-global \
  GLOBAL="$fake_home/.cursor" LOCAL_BIN="$fake_home/.local/bin" BIN="$BIN"
for skill in sc-run sc-ship sc-quick sc-storm; do
  test -f "$fake_home/.cursor/skills/$skill/SKILL.md" \
    || { echo "FAIL: install-global missing $skill skill"; exit 1; }
done
test -f "$fake_home/.cursor/skills/sc-discuss/references/lens-pass.md" \
  || { echo "FAIL: install-global missing sc-discuss/references/lens-pass.md"; exit 1; }
test -f "$fake_home/.cursor/skills/sc-judge/references/judge-break/empty-evidence/expect.json" \
  || { echo "FAIL: install-global missing sc-judge judge-break fixtures"; exit 1; }
echo "  ok   install-global installs sc-run, sc-ship, sc-quick, sc-storm, lens-pass, and judge-break"

# Lean User-layer skill allowlist (default install-global / global-sync):
# lifecycle + process only; domain encyclopedias stay out of ~/.cursor/skills.
global_skills="$fake_home/.cursor/skills"
lean_skills="sc-discuss sc-run sc-ship sc-quick sc-mission sc-planning sc-tdd sc-verification sc-judge sc-clarify sc-git sc-search sc-storm sc-writer"
domain_skills="sc-solid sc-security sc-performance sc-web-backend sc-web-frontend sc-database sc-firmware"
for skill in $lean_skills; do
  test -f "$global_skills/$skill/SKILL.md" \
    || { echo "FAIL: lean install-global missing lean-core skill $skill under $global_skills"; exit 1; }
done
for skill in $domain_skills; do
  if [ -e "$global_skills/$skill" ]; then
    echo "FAIL: lean install-global installed domain skill $skill under $global_skills (expected omit)"
    exit 1
  fi
done
echo "  ok   install-global lean allowlist: lean-core present, domain encyclopedias omitted"

# T3: re-running lean install-global prunes spacecraft-owned domain skills that
# sit outside the lean allowlist; unrelated paths under GLOBAL stay put.
mkdir -p "$global_skills/sc-solid"
printf '%s\n' '# seeded-for-lean-prune' > "$global_skills/sc-solid/SKILL.md"
printf '%s\n' 'user-keep' > "$fake_home/.cursor/user-unrelated-global.txt"
HOME="$fake_home" make -C "$ROOT" install-global \
  GLOBAL="$fake_home/.cursor" LOCAL_BIN="$fake_home/.local/bin" BIN="$BIN"
if [ -e "$global_skills/sc-solid" ]; then
  echo "FAIL: lean install-global left domain skill sc-solid under $global_skills (expected prune)"
  exit 1
fi
test -f "$fake_home/.cursor/user-unrelated-global.txt" \
  || { echo "FAIL: lean install-global removed unrelated $fake_home/.cursor/user-unrelated-global.txt"; exit 1; }
echo "  ok   lean install-global prunes domain skills; preserves unrelated GLOBAL paths"

# T3 docs: warn that lean reconcile removes spacecraft-owned domain packs.
if ! grep -Eqi 'prune|destructive' "$ROOT/docs/installation.md"; then
  echo "FAIL: docs/installation.md missing prune/destructive warning for lean reconcile"
  exit 1
fi
echo "  ok   installation.md warns lean reconcile is destructive for domain packs"

# T4: SPACECRAFT_SKILL_PROFILE=full (documented --full equivalent) installs
# domain encyclopedias; lean default inventory diverges (omits after reconcile).
HOME="$fake_home" SPACECRAFT_SKILL_PROFILE=full make -C "$ROOT" install-global \
  GLOBAL="$fake_home/.cursor" LOCAL_BIN="$fake_home/.local/bin" BIN="$BIN"
test -f "$global_skills/sc-solid/SKILL.md" \
  || { echo "FAIL: install-global with SPACECRAFT_SKILL_PROFILE=full missing domain skill sc-solid under $global_skills"; exit 1; }
HOME="$fake_home" make -C "$ROOT" install-global \
  GLOBAL="$fake_home/.cursor" LOCAL_BIN="$fake_home/.local/bin" BIN="$BIN"
if [ -e "$global_skills/sc-solid" ]; then
  echo "FAIL: lean install-global after full left domain skill sc-solid under $global_skills (lean/full inventories must diverge)"
  exit 1
fi
echo "  ok   install-global full keeps domain encyclopedias; lean omits"

user_rules="$fake_home/.cursor/spacecraft/USER-RULES.txt"
test -f "$user_rules" \
  || { echo "FAIL: install-global did not write $user_rules"; exit 1; }
for marker in 'Spacecraft' 'Intent coach' 'Agent chat language' 'Coding Standards' 'Project Structure' 'Lane Detection' 'Graph vs Loop' 'Context budget'; do
  grep -q "$marker" "$user_rules" \
    || { echo "FAIL: USER-RULES.txt missing marker: $marker"; exit 1; }
done
echo "  ok   install-global generates USER-RULES.txt with six-source markers (+ Graph vs Loop, Context budget)"

if [ -f "$fake_home/.cursorrules" ]; then
  echo "FAIL: install-global wrote legacy ~/.cursorrules"
  exit 1
fi
echo "  ok   install-global still writes no legacy ~/.cursorrules"

# T3 acceptance 1: global hooks.json gets the safety hooks merged in, and the
# pre-existing unrelated global hook (seeded above) survives the merge.
global_hooks="$fake_home/.cursor/hooks.json"
if ! grep -q 'user-unrelated-global.sh' "$global_hooks"; then
  echo "FAIL: install-global clobbered pre-existing global hook (user-unrelated-global.sh missing)"
  exit 1
fi
if ! grep -q 'check-main-write.sh' "$global_hooks"; then
  echo "FAIL: install-global did not merge check-main-write.sh into $global_hooks"
  exit 1
fi
if ! grep -q 'check-ship-commands.sh' "$global_hooks"; then
  echo "FAIL: install-global did not merge check-ship-commands.sh into $global_hooks"
  exit 1
fi
echo "  ok   install-global hooks merge preserves unrelated hook + adds safety hooks"

# T3 acceptance 2: installed hook commands use absolute or ~ paths (never the
# project-layer's repo-relative .cursor/hooks/...), and the scripts themselves
# are copied into ~/.cursor/hooks/ so they work outside this repo.
if ! grep -qE '"command": "(~|/)[^"]*check-main-write\.sh"' "$global_hooks"; then
  echo "FAIL: check-main-write.sh command in $global_hooks is not an absolute or ~ path"
  exit 1
fi
if ! grep -qE '"command": "(~|/)[^"]*check-ship-commands\.sh"' "$global_hooks"; then
  echo "FAIL: check-ship-commands.sh command in $global_hooks is not an absolute or ~ path"
  exit 1
fi
for hook_script in check-main-write.sh check-ship-commands.sh; do
  test -f "$fake_home/.cursor/hooks/$hook_script" \
    || { echo "FAIL: install-global did not copy $hook_script into $fake_home/.cursor/hooks/"; exit 1; }
done
echo "  ok   install-global hook paths are absolute/~ and scripts land in ~/.cursor/hooks/"

# T1 acceptance 1: tools_status_row prints ✔, tool name, green progress bar,
# percent, and version for a given success tool state (ANSI-stripped ok).
tools_status="$ROOT/scripts/tools-status.sh"
if [ ! -f "$tools_status" ]; then
  echo "FAIL: $tools_status does not exist yet (tools_status_row helper)"
  exit 1
fi
# shellcheck disable=SC1090
. "$tools_status"
if ! command -v tools_status_row >/dev/null 2>&1; then
  echo "FAIL: tools_status_row is not defined after sourcing $tools_status"
  exit 1
fi
row=$(tools_status_row "Spacecraft" ok 100 "v0.42.0") || {
  echo "FAIL: tools_status_row exited non-zero for success state"
  exit 1
}
esc=$(printf '\033')
# Require green ANSI on the raw row (progress bar / mark chrome).
printf '%s' "$row" | grep -Eq "${esc}\\[(32|92)m" \
  || { echo "FAIL: tools_status_row output missing green ANSI for progress/mark"; printf 'got: %s\n' "$row"; exit 1; }
# Strip CSI sequences for content asserts.
plain=$(printf '%s' "$row" | sed "s/${esc}\\[[0-9;]*[a-zA-Z]//g")
printf '%s' "$plain" | grep -q '✔' \
  || { echo "FAIL: tools_status_row missing success mark ✔"; printf 'got: %s\n' "$plain"; exit 1; }
printf '%s' "$plain" | grep -q 'Spacecraft' \
  || { echo "FAIL: tools_status_row missing tool name Spacecraft"; printf 'got: %s\n' "$plain"; exit 1; }
printf '%s' "$plain" | grep -Eq '\[[^]]*█[^]]*\]|\[[=#]+\]' \
  || { echo "FAIL: tools_status_row missing progress bar"; printf 'got: %s\n' "$plain"; exit 1; }
printf '%s' "$plain" | grep -Eq '100[[:space:]]*%' \
  || { echo "FAIL: tools_status_row missing percent 100%"; printf 'got: %s\n' "$plain"; exit 1; }
printf '%s' "$plain" | grep -q 'v0.42.0' \
  || { echo "FAIL: tools_status_row missing version v0.42.0"; printf 'got: %s\n' "$plain"; exit 1; }
echo "  ok   tools_status_row prints ✔, name, green bar, percent, version"

# T1 acceptance 2: tools_status_box wraps rows under a box-drawing header
# labeled Tools (ANSI-stripped ok).
if ! command -v tools_status_box >/dev/null 2>&1; then
  echo "FAIL: tools_status_box is not defined after sourcing $tools_status"
  exit 1
fi
box=$(tools_status_box "$row") || {
  echo "FAIL: tools_status_box exited non-zero"
  exit 1
}
box_plain=$(printf '%s' "$box" | sed "s/${esc}\\[[0-9;]*[a-zA-Z]//g")
# Box-drawing: at least one of ─ │ ┌ ┐ └ ┘ ┬ ┴ ├ ┤ ┼ ╔ ╗ ╚ ╝ ═ ║
printf '%s' "$box_plain" | grep -Eq '[─│┌┐└┘┬┴├┤┼╔╗╚╝═║]' \
  || { echo "FAIL: tools_status_box output missing box-drawing characters"; printf 'got: %s\n' "$box_plain"; exit 1; }
printf '%s' "$box_plain" | grep -q 'Tools' \
  || { echo "FAIL: tools_status_box output missing header label Tools"; printf 'got: %s\n' "$box_plain"; exit 1; }
# Wrapped row content still present inside the box.
printf '%s' "$box_plain" | grep -q 'Spacecraft' \
  || { echo "FAIL: tools_status_box did not wrap the tools_status_row content"; printf 'got: %s\n' "$box_plain"; exit 1; }
echo "  ok   tools_status_box wraps rows under box-drawing header Tools"

# T2 acceptance 1: install-machine.sh clones spacecraft into a durable path
# and completes an install-global-equivalent User-layer install under fake HOME.
#
# Network-free smoke env contract (for sc-coder / CI):
#   SPACECRAFT_INSTALL_DIR  Durable clone directory (default: $HOME/.local/share/spacecraft).
#   SPACECRAFT_CLONE_SRC    Local filesystem path used as clone source when set
#                           (e.g. this repo root). Prefer `git clone "$SPACECRAFT_CLONE_SRC"
#                           "$SPACECRAFT_INSTALL_DIR"` or an equivalent local copy so CI
#                           never hits the network. When unset, production may clone from
#                           the remote (see SPACECRAFT_REPO / SPACECRAFT_REF in bootstrap.sh).
# Companion fixture overrides (required on every install-machine invocation so CI
# never hits curl|sh; T2 asserts stay about clone/User-layer, not companion markers):
#   SPACECRAFT_CAVEMAN_INSTALL / SPACECRAFT_RTK_INSTALL / SPACECRAFT_CODEGRAPH_INSTALL
#   SPACECRAFT_COMPANION_MARKERS  Shared markers dir under $tmp
#   PATH                          Prepend scripts/testdata/install-machine/fake-bin
# Companions / Tools status rows are out of scope for this acceptance (see T3).
machine_installer="$ROOT/scripts/install-machine.sh"
if [ ! -f "$machine_installer" ]; then
  echo "FAIL: $machine_installer does not exist yet (machine User-layer installer)"
  exit 1
fi

fixtures_dir="$ROOT/scripts/testdata/install-machine"
companion_markers="$tmp/companion-markers"
mkdir -p "$companion_markers"
companion_path="$fixtures_dir/fake-bin:${PATH:-/usr/bin:/bin}"

machine_home="$tmp/machine-home"
mkdir -p "$machine_home"
machine_install_dir="$machine_home/.local/share/spacecraft"

HOME="$machine_home" \
  PATH="$companion_path" \
  SPACECRAFT_INSTALL_DIR="$machine_install_dir" \
  SPACECRAFT_CLONE_SRC="$ROOT" \
  SPACECRAFT_COMPANION_MARKERS="$companion_markers" \
  SPACECRAFT_CAVEMAN_INSTALL="$fixtures_dir/caveman-install.sh" \
  SPACECRAFT_RTK_INSTALL="$fixtures_dir/rtk-install.sh" \
  SPACECRAFT_CODEGRAPH_INSTALL="$fixtures_dir/codegraph-install.sh" \
  sh "$machine_installer" \
  || { echo "FAIL: install-machine.sh exited non-zero under fake HOME"; exit 1; }

# Durable clone has spacecraft content.
test -d "$machine_install_dir" \
  || { echo "FAIL: install-machine did not create durable clone dir $machine_install_dir"; exit 1; }
if [ ! -f "$machine_install_dir/Makefile" ] && [ ! -d "$machine_install_dir/.cursor/agents" ]; then
  echo "FAIL: durable clone $machine_install_dir missing spacecraft content (Makefile or .cursor/agents)"
  exit 1
fi

# User-layer outcomes equivalent to install-global (under fake HOME).
test -e "$machine_home/.local/bin/spacecraft" \
  || { echo "FAIL: install-machine did not link CLI at $machine_home/.local/bin/spacecraft"; exit 1; }
for skill in sc-run sc-ship sc-quick; do
  test -f "$machine_home/.cursor/skills/$skill/SKILL.md" \
    || { echo "FAIL: install-machine User layer missing $skill skill under fake HOME"; exit 1; }
done
test -d "$machine_home/.cursor/agents" \
  || { echo "FAIL: install-machine User layer missing agents under fake HOME"; exit 1; }
machine_user_rules="$machine_home/.cursor/spacecraft/USER-RULES.txt"
test -f "$machine_user_rules" \
  || { echo "FAIL: install-machine did not write $machine_user_rules"; exit 1; }
echo "  ok   install-machine clones durable path + User-layer install under fake HOME"

# T2 acceptance 2: re-run against an existing clone updates in place (git pull)
# instead of re-cloning, and still completes the User-layer install.
# An untracked marker in the working tree must survive (re-clone would wipe it).
machine_marker="$machine_install_dir/.install-machine-rerun-marker"
machine_marker_token="rerun-marker-$(date +%s)-$$"
printf '%s\n' "$machine_marker_token" > "$machine_marker"
test -f "$machine_marker" \
  || { echo "FAIL: could not plant untracked marker in durable clone"; exit 1; }

machine_rerun_log="$tmp/install-machine-rerun.log"
HOME="$machine_home" \
  PATH="$companion_path" \
  SPACECRAFT_INSTALL_DIR="$machine_install_dir" \
  SPACECRAFT_CLONE_SRC="$ROOT" \
  SPACECRAFT_COMPANION_MARKERS="$companion_markers" \
  SPACECRAFT_CAVEMAN_INSTALL="$fixtures_dir/caveman-install.sh" \
  SPACECRAFT_RTK_INSTALL="$fixtures_dir/rtk-install.sh" \
  SPACECRAFT_CODEGRAPH_INSTALL="$fixtures_dir/codegraph-install.sh" \
  sh "$machine_installer" > "$machine_rerun_log" 2>&1 \
  || {
    echo "FAIL: install-machine.sh re-run exited non-zero under fake HOME"
    cat "$machine_rerun_log"
    exit 1
  }

grep -q 'Updating existing clone' "$machine_rerun_log" \
  || {
    echo "FAIL: re-run stdout missing 'Updating existing clone' (expected in-place pull, not re-clone)"
    cat "$machine_rerun_log"
    exit 1
  }
grep -qiE 'Cloning (local source|https?://|git@)' "$machine_rerun_log" \
  && {
    echo "FAIL: re-run stdout mentions Cloning — expected update-in-place, not re-clone"
    cat "$machine_rerun_log"
    exit 1
  }

test -f "$machine_marker" \
  || { echo "FAIL: re-run wiped untracked marker $machine_marker (looks like a fresh clone)"; exit 1; }
grep -qxF "$machine_marker_token" "$machine_marker" \
  || { echo "FAIL: untracked marker content changed after re-run"; exit 1; }

test -e "$machine_home/.local/bin/spacecraft" \
  || { echo "FAIL: after re-run, CLI missing at $machine_home/.local/bin/spacecraft"; exit 1; }
for skill in sc-run sc-ship sc-quick; do
  test -f "$machine_home/.cursor/skills/$skill/SKILL.md" \
    || { echo "FAIL: after re-run, User layer missing $skill skill under fake HOME"; exit 1; }
done
test -d "$machine_home/.cursor/agents" \
  || { echo "FAIL: after re-run, User layer missing agents under fake HOME"; exit 1; }
test -f "$machine_user_rules" \
  || { echo "FAIL: after re-run, missing $machine_user_rules"; exit 1; }
echo "  ok   install-machine re-run updates in place + preserves User-layer"

# T3 acceptance 1: successful path runs each companion install.sh then Cursor-wire
# (rtk init -g --agent cursor; codegraph install --target=cursor --yes) and marks
# Tools rows done with a version. Network-free via fixture overrides.
#
# Env / fixture contract for sc-coder (honor in install-machine.sh):
#   SPACECRAFT_CAVEMAN_INSTALL     Path to caveman install.sh substitute (skip curl|sh)
#   SPACECRAFT_RTK_INSTALL         Path to rtk install.sh substitute
#   SPACECRAFT_CODEGRAPH_INSTALL   Path to codegraph install.sh substitute
#   SPACECRAFT_COMPANION_MARKERS   Dir where fixtures write invocation markers
#   PATH                           Prepend scripts/testdata/install-machine/fake-bin so
#                                  wire steps hit fake `rtk` / `codegraph` / `caveman`
#                                  (--version -> fixture versions below).
#
# Production (env unset): official curl|sh installers, then real wire commands.
# Caveman has no separate Cursor-wire beyond install.sh.
#
# Marker files under SPACECRAFT_COMPANION_MARKERS:
#   caveman.install | rtk.install | codegraph.install
#   rtk.wire | codegraph.wire
#
# Fixture versions (Tools rows must show these on success):
#   Caveman v0.0.1-fixture | RTK v0.0.2-fixture | CodeGraph v0.0.3-fixture
# Spacecraft row must also show ✔ + some non-empty version.
# Reuse shared fixtures_dir / companion_markers / companion_path from T2.
rm -f "$companion_markers"/*
companion_home="$tmp/companion-home"
companion_install_dir="$companion_home/.local/share/spacecraft"
mkdir -p "$companion_home"
companion_log="$tmp/install-machine-companions.log"

HOME="$companion_home" \
  PATH="$companion_path" \
  SPACECRAFT_INSTALL_DIR="$companion_install_dir" \
  SPACECRAFT_CLONE_SRC="$ROOT" \
  SPACECRAFT_COMPANION_MARKERS="$companion_markers" \
  SPACECRAFT_CAVEMAN_INSTALL="$fixtures_dir/caveman-install.sh" \
  SPACECRAFT_RTK_INSTALL="$fixtures_dir/rtk-install.sh" \
  SPACECRAFT_CODEGRAPH_INSTALL="$fixtures_dir/codegraph-install.sh" \
  sh "$machine_installer" > "$companion_log" 2>&1 \
  || {
    echo "FAIL: install-machine.sh (companions path) exited non-zero under fake HOME"
    cat "$companion_log"
    exit 1
  }

for marker in caveman.install rtk.install codegraph.install rtk.wire codegraph.wire; do
  test -f "$companion_markers/$marker" \
    || {
      echo "FAIL: companion orchestration missing marker $companion_markers/$marker"
      echo "--- install-machine stdout/stderr ---"
      cat "$companion_log"
      exit 1
    }
done

companion_plain=$(sed "s/$(printf '\033')\\[[0-9;]*[a-zA-Z]//g" "$companion_log")
printf '%s' "$companion_plain" | grep -q 'Tools' \
  || {
    echo "FAIL: install-machine stdout missing Tools status header"
    cat "$companion_log"
    exit 1
  }

# Each companion row: ✔ + name + fixture version (ANSI-stripped).
for pair in "Caveman|v0.0.1-fixture" "RTK|v0.0.2-fixture" "CodeGraph|v0.0.3-fixture"; do
  name=${pair%%|*}
  ver=${pair##*|}
  printf '%s\n' "$companion_plain" | grep -E "✔[^[:cntrl:]]*$name" | grep -q "$ver" \
    || {
      echo "FAIL: Tools status missing success row for $name with version $ver"
      echo "--- stripped stdout ---"
      printf '%s\n' "$companion_plain"
      exit 1
    }
done

# Spacecraft row done with some version token.
printf '%s\n' "$companion_plain" | grep -E "✔[^[:cntrl:]]*Spacecraft" | grep -Eq 'v?[0-9]|fixture|dev|main' \
  || {
    echo "FAIL: Tools status missing Spacecraft success row with a version"
    echo "--- stripped stdout ---"
    printf '%s\n' "$companion_plain"
    exit 1
  }

echo "  ok   install-machine companions install+wire + Tools ✔/versions (fixtures)"

# T3 acceptance 2: per-companion SPACECRAFT_*_INSTALL overrides prove orchestration
# without live curl (fixtures on PATH; curl spy must not be invoked).
printf '%s\n' "$companion_plain" | grep -q 'override:' \
  || {
    echo "FAIL: companion success stdout missing 'override:' (expected fixture installer path)"
    cat "$companion_log"
    exit 1
  }
printf '%s\n' "$companion_plain" | grep -qi 'official installer' \
  && {
    echo "FAIL: companion success stdout mentions 'official installer' despite overrides"
    cat "$companion_log"
    exit 1
  }
test ! -f "$companion_markers/curl.called" \
  || {
    echo "FAIL: curl was invoked despite SPACECRAFT_*_INSTALL overrides"
    echo "--- curl spy args ---"
    cat "$companion_markers/curl.called"
    echo "--- install-machine stdout/stderr ---"
    cat "$companion_log"
    exit 1
  }
echo "  ok   companion installer overrides skip curl (override path, no official installer)"

# T3 acceptance 3: one companion failing (fixture exits non-zero) soft-fails —
# overall exit 0, that Tools row shows ✘, other companions + Spacecraft complete,
# User-layer still present, other companions' markers still written.
rm -f "$companion_markers"/*
soft_fail_home="$tmp/soft-fail-home"
soft_fail_install_dir="$soft_fail_home/.local/share/spacecraft"
mkdir -p "$soft_fail_home"
soft_fail_log="$tmp/install-machine-soft-fail.log"
soft_fail_rc=0
HOME="$soft_fail_home" \
  PATH="$companion_path" \
  SPACECRAFT_INSTALL_DIR="$soft_fail_install_dir" \
  SPACECRAFT_CLONE_SRC="$ROOT" \
  SPACECRAFT_COMPANION_MARKERS="$companion_markers" \
  SPACECRAFT_CAVEMAN_INSTALL="$fixtures_dir/caveman-install.sh" \
  SPACECRAFT_RTK_INSTALL="$fixtures_dir/rtk-install-fail.sh" \
  SPACECRAFT_CODEGRAPH_INSTALL="$fixtures_dir/codegraph-install.sh" \
  sh "$machine_installer" > "$soft_fail_log" 2>&1 \
  || soft_fail_rc=$?

test "$soft_fail_rc" -eq 0 \
  || {
    echo "FAIL: soft-fail path: install-machine.sh exited $soft_fail_rc (expected 0)"
    cat "$soft_fail_log"
    exit 1
  }

# Failing companion must not write success install/wire markers.
test ! -f "$companion_markers/rtk.install" \
  || {
    echo "FAIL: soft-fail path wrote rtk.install despite failing fixture"
    cat "$soft_fail_log"
    exit 1
  }
test ! -f "$companion_markers/rtk.wire" \
  || {
    echo "FAIL: soft-fail path ran rtk wire despite install failure"
    cat "$soft_fail_log"
    exit 1
  }

# Other companions still completed (install + wire markers).
for marker in caveman.install codegraph.install codegraph.wire; do
  test -f "$companion_markers/$marker" \
    || {
      echo "FAIL: soft-fail path missing other-companion marker $companion_markers/$marker"
      echo "--- install-machine stdout/stderr ---"
      cat "$soft_fail_log"
      exit 1
    }
done

# User-layer still completed under soft-fail HOME.
test -e "$soft_fail_home/.local/bin/spacecraft" \
  || {
    echo "FAIL: soft-fail path missing User-layer CLI at $soft_fail_home/.local/bin/spacecraft"
    cat "$soft_fail_log"
    exit 1
  }
test -f "$soft_fail_home/.cursor/spacecraft/USER-RULES.txt" \
  || {
    echo "FAIL: soft-fail path missing USER-RULES.txt"
    cat "$soft_fail_log"
    exit 1
  }

soft_fail_plain=$(sed "s/$(printf '\033')\\[[0-9;]*[a-zA-Z]//g" "$soft_fail_log")

# RTK Tools row: fail marker ✘ (not ✔).
printf '%s\n' "$soft_fail_plain" | grep -E "✘[^[:cntrl:]]*RTK" >/dev/null \
  || {
    echo "FAIL: soft-fail Tools status missing ✘ RTK row"
    echo "--- stripped stdout ---"
    printf '%s\n' "$soft_fail_plain"
    exit 1
  }
printf '%s\n' "$soft_fail_plain" | grep -E "✔[^[:cntrl:]]*RTK" >/dev/null \
  && {
    echo "FAIL: soft-fail Tools status still shows ✔ RTK (expected fail marker only)"
    echo "--- stripped stdout ---"
    printf '%s\n' "$soft_fail_plain"
    exit 1
  }

# Caveman + CodeGraph + Spacecraft still ✔.
for name in Caveman CodeGraph Spacecraft; do
  printf '%s\n' "$soft_fail_plain" | grep -E "✔[^[:cntrl:]]*$name" >/dev/null \
    || {
      echo "FAIL: soft-fail Tools status missing ✔ $name row"
      echo "--- stripped stdout ---"
      printf '%s\n' "$soft_fail_plain"
      exit 1
    }
done

echo "  ok   companion soft-fail: exit 0, ✘ RTK, others+Spacecraft ✔, User-layer intact"

# T4 acceptance 1 (install-machine flags): default run never prompts; --agents cursor
# and --yes are accepted as explicit non-interactive flags with the same Cursor-wire
# effect. Production must parse the flags (echo confirmation) — silent ignore is not enough.
#
# Prompt patterns that must never appear (stdout/stderr, ANSI-stripped):
#   select | Choose | [Y/n] | read
# Parsed-flag confirmation (required once near start of a flagged run):
#   agents: cursor
#   non-interactive: yes
_assert_no_install_prompts() {
  log_file=$1
  label=$2
  plain=$(sed "s/$(printf '\033')\\[[0-9;]*[a-zA-Z]//g" "$log_file")
  if printf '%s\n' "$plain" | grep -Eqi 'select|Choose|\[Y/n\]|[[:space:]]read[[:space:]]'; then
    echo "FAIL: $label showed an interactive prompt pattern"
    echo "--- stripped stdout/stderr ---"
    printf '%s\n' "$plain"
    exit 1
  fi
}

_assert_no_install_prompts "$companion_log" "default install-machine run (no --agents)"

rm -f "$companion_markers"/*
flags_home="$tmp/flags-home"
flags_install_dir="$flags_home/.local/share/spacecraft"
mkdir -p "$flags_home"
flags_log="$tmp/install-machine-flags.log"
flags_rc=0
HOME="$flags_home" \
  PATH="$companion_path" \
  SPACECRAFT_INSTALL_DIR="$flags_install_dir" \
  SPACECRAFT_CLONE_SRC="$ROOT" \
  SPACECRAFT_COMPANION_MARKERS="$companion_markers" \
  SPACECRAFT_CAVEMAN_INSTALL="$fixtures_dir/caveman-install.sh" \
  SPACECRAFT_RTK_INSTALL="$fixtures_dir/rtk-install.sh" \
  SPACECRAFT_CODEGRAPH_INSTALL="$fixtures_dir/codegraph-install.sh" \
  sh "$machine_installer" --agents cursor --yes > "$flags_log" 2>&1 \
  || flags_rc=$?

test "$flags_rc" -eq 0 \
  || {
    echo "FAIL: install-machine.sh --agents cursor --yes exited $flags_rc (expected 0)"
    cat "$flags_log"
    exit 1
  }

_assert_no_install_prompts "$flags_log" "install-machine --agents cursor --yes"

flags_plain=$(sed "s/$(printf '\033')\\[[0-9;]*[a-zA-Z]//g" "$flags_log")
printf '%s\n' "$flags_plain" | grep -Eq '^agents:[[:space:]]*cursor$' \
  || {
    echo "FAIL: flagged run missing parsed confirmation line 'agents: cursor'"
    echo "--- stripped stdout/stderr ---"
    printf '%s\n' "$flags_plain"
    exit 1
  }
printf '%s\n' "$flags_plain" | grep -Eq '^non-interactive:[[:space:]]*yes$' \
  || {
    echo "FAIL: flagged run missing parsed confirmation line 'non-interactive: yes'"
    echo "--- stripped stdout/stderr ---"
    printf '%s\n' "$flags_plain"
    exit 1
  }

# Same Cursor-wire effect as default (fixtures write rtk.wire / codegraph.wire).
for marker in rtk.wire codegraph.wire; do
  test -f "$companion_markers/$marker" \
    || {
      echo "FAIL: --agents cursor --yes missing Cursor-wire marker $companion_markers/$marker"
      echo "--- install-machine stdout/stderr ---"
      cat "$flags_log"
      exit 1
    }
done

echo "  ok   install-machine --agents cursor --yes: parsed, non-interactive, Cursor-wired"

# T4 acceptance 2: install-machine never invokes Project-layer bootstrap.sh and
# never attempts to install the Cursor app itself.
#
# Behavioral seam: captured fixture-run logs under fake-bin PATH (not source greps).
# A fake bootstrap.sh on PATH is weak (bootstrap is not looked up on PATH).
# Assert execution-shaped lines only — bare "bootstrap.sh" in clone docs must not
# false-positive (require sh/bash/./ runner or "Running ... bootstrap").
# Cursor app: log patterns + brew/curl PATH spies (fake-bin/brew, fake-bin/curl).
# Note: machine installer logs from fixture runs must show no bootstrap execution.
_assert_no_bootstrap_exec_or_cursor_app() {
  log_file=$1
  label=$2
  plain=$(sed "s/$(printf '\033')\\[[0-9;]*[a-zA-Z]//g" "$log_file")

  # Execution-shaped bootstrap invocations (not mere filename mentions).
  if printf '%s\n' "$plain" | grep -Eqi \
    '(^|[[:space:];|&])(\./|sh[[:space:]]+|bash[[:space:]]+)[^[:space:]]*bootstrap\.sh|Running[[:space:]].*bootstrap'; then
    echo "FAIL: $label log shows bootstrap.sh execution"
    echo "--- stripped stdout/stderr ---"
    printf '%s\n' "$plain"
    exit 1
  fi

  # Cursor app install attempts (app binary, download URL, brew cask, install*app).
  if printf '%s\n' "$plain" | grep -Eqi \
    'Cursor\.app|cursor\.com/downloads|brew[[:space:]]+install[[:space:]]+--cask[[:space:]]+cursor|install[^[:cntrl:]]*cursor[^[:cntrl:]]*app'; then
    echo "FAIL: $label log shows Cursor app install attempt"
    echo "--- stripped stdout/stderr ---"
    printf '%s\n' "$plain"
    exit 1
  fi

  # fake-bin/brew prints this if invoked — survives even after marker dirs are wiped.
  if printf '%s\n' "$plain" | grep -qi 'fake brew:'; then
    echo "FAIL: $label invoked brew PATH spy (Cursor/app package managers must not run)"
    echo "--- stripped stdout/stderr ---"
    printf '%s\n' "$plain"
    exit 1
  fi
}

_assert_no_bootstrap_exec_or_cursor_app "$companion_log" "companion success install-machine run"
_assert_no_bootstrap_exec_or_cursor_app "$flags_log" "install-machine --agents cursor --yes"

# PATH spy markers from the last fixture run (flags): no brew, no Cursor curl.
if [ -f "$companion_markers/brew.called" ]; then
  echo "FAIL: brew PATH spy was invoked during install-machine (unexpected)"
  echo "--- brew.called ---"
  cat "$companion_markers/brew.called"
  exit 1
fi
if [ -f "$companion_markers/curl.called" ] && \
  grep -Eqi 'cursor\.com/downloads|Cursor\.app' "$companion_markers/curl.called"; then
  echo "FAIL: curl PATH spy recorded Cursor app download"
  echo "--- curl.called ---"
  cat "$companion_markers/curl.called"
  exit 1
fi

echo "  ok   install-machine: no bootstrap.sh exec, no Cursor app install (log + brew spy)"

# T4 acceptance 3: install-machine final stdout includes its own post-install notes
# (spec Output): PATH (~/.local/bin), USER-RULES.txt paste, restart Cursor, and
# per-project codegraph init. Assert on companion success log (ANSI-stripped).
# Nested install-global may already print paste/Restart; PATH tilde form and
# codegraph init must still appear from install-machine's final messaging.
#
# Exact substrings (coder must emit; independent of expanded $HOME paths):
#   1. ~/.local/bin   plus PATH on the same line (PATH note)
#   2. USER-RULES.txt with paste (paste instruction)
#   3. Restart Cursor
#   4. codegraph init
_assert_post_install_notes() {
  log_file=$1
  label=$2
  plain=$(sed "s/$(printf '\033')\\[[0-9;]*[a-zA-Z]//g" "$log_file")

  printf '%s\n' "$plain" | grep -E '~/.local/bin' | grep -Eqi 'PATH' \
    || {
      echo "FAIL: $label missing PATH (~/.local/bin) note"
      echo "--- stripped stdout/stderr ---"
      printf '%s\n' "$plain"
      exit 1
    }

  printf '%s\n' "$plain" | grep -F 'USER-RULES.txt' | grep -Eqi 'paste' \
    || {
      echo "FAIL: $label missing USER-RULES.txt paste instruction"
      echo "--- stripped stdout/stderr ---"
      printf '%s\n' "$plain"
      exit 1
    }

  printf '%s\n' "$plain" | grep -Fq 'Restart Cursor' \
    || {
      echo "FAIL: $label missing Restart Cursor reminder"
      echo "--- stripped stdout/stderr ---"
      printf '%s\n' "$plain"
      exit 1
    }

  printf '%s\n' "$plain" | grep -Fq 'codegraph init' \
    || {
      echo "FAIL: $label missing per-project codegraph init note"
      echo "--- stripped stdout/stderr ---"
      printf '%s\n' "$plain"
      exit 1
    }
}

_assert_post_install_notes "$companion_log" "companion success install-machine run"
echo "  ok   install-machine post-install notes (PATH, USER-RULES paste, Restart Cursor, codegraph init)"

# T5 acceptance 1: make install-machine builds the CLI and invokes
# scripts/install-machine.sh. Use make -n only (no second full install); A2
# covers end-to-end smoke. HOME override works the same way as install-global
# under fake HOME in this script.
makefile_dry="$tmp/install-machine-make-n.log"
if ! make -n -C "$ROOT" install-machine >"$makefile_dry" 2>&1; then
  echo "FAIL: make -n install-machine failed (target missing or recipe error)"
  echo "--- make -n stdout/stderr ---"
  cat "$makefile_dry"
  exit 1
fi
grep -Fq 'scripts/install-machine.sh' "$makefile_dry" \
  || {
    echo "FAIL: make -n install-machine dry-run missing scripts/install-machine.sh"
    echo "--- make -n stdout/stderr ---"
    cat "$makefile_dry"
    exit 1
  }
# Must wire the Node CLI (recipe shows install-cli and/or spacecraft.mjs).
if ! grep -Eqi 'install-cli|spacecraft\.mjs|node' "$makefile_dry"; then
  echo "FAIL: make -n install-machine dry-run does not show Node CLI wire (install-cli / spacecraft.mjs)"
  echo "--- make -n stdout/stderr ---"
  cat "$makefile_dry"
  exit 1
fi
if grep -Fq 'go build' "$makefile_dry"; then
  echo "FAIL: make -n install-machine dry-run still invokes go build"
  echo "--- make -n stdout/stderr ---"
  cat "$makefile_dry"
  exit 1
fi
make -C "$ROOT" help 2>/dev/null | grep -Fq 'install-machine' \
  || {
    echo "FAIL: make help does not mention install-machine"
    exit 1
  }
echo "  ok   make install-machine target builds CLI + invokes scripts/install-machine.sh"

# T5 acceptance 2: end-to-end consolidation gate. Re-check critical markers from
# T2–T4 in one place so a skipped prior section cannot leave make test-install
# green. Uses the companion success log + companion_home / companion_install_dir
# (later soft-fail/flags runs wipe companion_markers).
test -f "$companion_log" \
  || {
    echo "FAIL: T5 e2e gate missing companion success log $companion_log (prior section skipped?)"
    exit 1
  }

# Durable clone (filesystem + log evidence).
test -d "$companion_install_dir" \
  || {
    echo "FAIL: T5 e2e gate missing durable clone dir $companion_install_dir"
    exit 1
  }
if [ ! -f "$companion_install_dir/Makefile" ] && [ ! -d "$companion_install_dir/.cursor/agents" ]; then
  echo "FAIL: T5 e2e gate durable clone missing spacecraft content (Makefile or .cursor/agents)"
  exit 1
fi
e2e_plain=$(sed "s/$(printf '\033')\\[[0-9;]*[a-zA-Z]//g" "$companion_log")
printf '%s\n' "$e2e_plain" | grep -Eqi 'Cloning (local source|https?://|git@)|Updating existing clone' \
  || {
    echo "FAIL: T5 e2e gate companion_log missing durable clone evidence (Cloning/Updating)"
    echo "--- stripped stdout/stderr ---"
    printf '%s\n' "$e2e_plain"
    exit 1
  }

# User-layer CLI / agents / USER-RULES.txt under companion success HOME.
test -e "$companion_home/.local/bin/spacecraft" \
  || {
    echo "FAIL: T5 e2e gate missing User-layer CLI at $companion_home/.local/bin/spacecraft"
    exit 1
  }
test -d "$companion_home/.cursor/agents" \
  || {
    echo "FAIL: T5 e2e gate missing User-layer agents under $companion_home/.cursor/agents"
    exit 1
  }
test -f "$companion_home/.cursor/spacecraft/USER-RULES.txt" \
  || {
    echo "FAIL: T5 e2e gate missing $companion_home/.cursor/spacecraft/USER-RULES.txt"
    exit 1
  }

# Stubbed companion orchestration + Tools status markers in captured stdout.
printf '%s\n' "$e2e_plain" | grep -q 'override:' \
  || {
    echo "FAIL: T5 e2e gate companion_log missing stubbed companion 'override:' evidence"
    echo "--- stripped stdout/stderr ---"
    printf '%s\n' "$e2e_plain"
    exit 1
  }
printf '%s\n' "$e2e_plain" | grep -q 'Tools' \
  || {
    echo "FAIL: T5 e2e gate companion_log missing Tools status header"
    echo "--- stripped stdout/stderr ---"
    printf '%s\n' "$e2e_plain"
    exit 1
  }
for pair in "Caveman|v0.0.1-fixture" "RTK|v0.0.2-fixture" "CodeGraph|v0.0.3-fixture"; do
  name=${pair%%|*}
  ver=${pair##*|}
  printf '%s\n' "$e2e_plain" | grep -E "✔[^[:cntrl:]]*$name" | grep -q "$ver" \
    || {
      echo "FAIL: T5 e2e gate Tools status missing ✔ $name with version $ver"
      echo "--- stripped stdout/stderr ---"
      printf '%s\n' "$e2e_plain"
      exit 1
    }
done
printf '%s\n' "$e2e_plain" | grep -E "✔[^[:cntrl:]]*Spacecraft" | grep -Eq 'v?[0-9]|fixture|dev|main' \
  || {
    echo "FAIL: T5 e2e gate Tools status missing ✔ Spacecraft with a version"
    echo "--- stripped stdout/stderr ---"
    printf '%s\n' "$e2e_plain"
    exit 1
  }

echo "  ok   T5 e2e gate: durable clone + User-layer + stubbed companions + Tools marks"

# Permanent smoke: install-machine must not leave caveman/agentkit debris at repo root.
if [ -e "$ROOT/.agents" ]; then
  echo "FAIL: repo root has .agents/ (install-machine must not leave companion debris)"
  exit 1
fi
if [ -f "$ROOT/skills-lock.json" ]; then
  echo "FAIL: repo root has skills-lock.json (install-machine must not leave companion debris)"
  exit 1
fi
echo "  ok   repo root has no .agents/ and no skills-lock.json"

# Permanent smoke: installation docs mention install-machine, Tools status, and companions.
for doc in "$ROOT/docs/installation.md" "$ROOT/README.md"; do
  test -f "$doc" || { echo "FAIL: missing doc $doc"; exit 1; }
  for marker in install-machine Tools caveman rtk codegraph; do
    grep -qi "$marker" "$doc" \
      || { echo "FAIL: $doc missing marker: $marker"; exit 1; }
  done
done
echo "  ok   docs/installation.md + README.md mention install-machine, Tools, companions"

# Critical fix: official caveman path must pass --only cursor (Cursor-only install).
# Use -- so BSD grep does not treat the pattern as a long option.
grep -Fq -- '--only cursor' "$machine_installer" \
  || {
    echo "FAIL: $machine_installer missing '--only cursor' (caveman official path)"
    exit 1
  }
echo "  ok   install-machine.sh mentions --only cursor for caveman"

echo "Install smoke OK"
