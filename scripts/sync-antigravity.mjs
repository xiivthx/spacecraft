import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);
const ROOT = path.resolve(__dirname, '..');
const PLUGIN_DIR = path.join(ROOT, 'plugins', 'spacecraft');

console.log('Syncing Spacecraft for Antigravity...');

// 1. Ensure plugin directory structure
const dirs = [
  PLUGIN_DIR,
  path.join(PLUGIN_DIR, 'rules'),
  path.join(PLUGIN_DIR, 'skills'),
  path.join(PLUGIN_DIR, 'agents'),
  path.join(PLUGIN_DIR, 'hooks'),
];
for (const dir of dirs) {
  fs.mkdirSync(dir, { recursive: true });
}

// 2. Generate plugin.json
const pluginJson = {
  name: 'spacecraft',
  version: '1.0.0',
  description: 'Spacecraft mission-control harness for Antigravity. Features always-on rules, specialized subagents, UX/UI/Frontend design gates, Chrome DevTools MCP live probes, safety hooks, and strict Source of Trust.',
  author: {
    name: 'Spacecraft Team',
  },
  repository: 'https://github.com/xiivthx/spacecraft',
  license: 'MIT',
};
fs.writeFileSync(path.join(PLUGIN_DIR, 'plugin.json'), JSON.stringify(pluginJson, null, 2) + '\n');
console.log('  ok: generated plugins/spacecraft/plugin.json');

// 3. Generate hooks.json for Antigravity
const hooksJson = {
  "spacecraft-safety-gate": {
    "PreToolUse": [
      {
        "matcher": "run_command",
        "hooks": [
          {
            "type": "command",
            "command": "node ./hooks/safety-check.mjs",
            "timeout": 15
          }
        ]
      }
    ]
  }
};
fs.writeFileSync(path.join(PLUGIN_DIR, 'hooks.json'), JSON.stringify(hooksJson, null, 2) + '\n');
console.log('  ok: generated plugins/spacecraft/hooks.json');

// 4. Generate hooks/safety-check.mjs (keep short; mirror plugins/spacecraft/hooks/safety-check.mjs)
const safetyCheckMjs = `#!/usr/bin/env node
/**
 * Spacecraft Antigravity Safety Hook (PreToolUse on run_command).
 * Blocks: secrets in command paths, force-push, catastrophic rm, mutating git on main,
 * and ship ops without SPACECRAFT_SHIP=1. git push is denied in-agent (no Cursor-style ask).
 */
import { execSync } from 'node:child_process';

function readStdin() {
  return new Promise((resolve) => {
    let data = '';
    process.stdin.setEncoding('utf8');
    process.stdin.on('data', (chunk) => { data += chunk; });
    process.stdin.on('end', () => resolve(data));
  });
}

function respond(decision, reason) {
  const out = { decision };
  if (reason) out.reason = reason;
  console.log(JSON.stringify(out));
  process.exit(0);
}

function isSecretPath(p) {
  const base = p.split(/[\\\\/]/).pop() || '';
  if (/^\\.env\\.example$/i.test(base) || /\\.env\\.sample$/i.test(base)) return false;
  if (/^\\.env$/i.test(base) || /^\\.env\\./i.test(base)) return true;
  if (/(^|[/\\\\])(credentials|id_rsa|id_ed25519)(\\.|$)/i.test(p)) return true;
  if (/\\.(pem|p12|pfx|key)$/i.test(base) && !/\\.pub$/i.test(base)) return true;
  return false;
}

async function main() {
  const raw = await readStdin();
  if (!raw.trim()) return respond('allow');
  let payload;
  try {
    payload = JSON.parse(raw);
  } catch {
    return respond('allow');
  }

  const toolCall = payload.toolCall || {};
  const name = toolCall.name || '';
  const args = toolCall.args || {};
  const cmd = String(args.CommandLine || args.command || '');

  const filePath = String(args.path || args.file_path || args.filePath || '');
  if (filePath && isSecretPath(filePath)) {
    return respond('deny', 'Spacecraft Safety Gate: secret file read blocked (.env / credentials / keys).');
  }

  if (name !== 'run_command') {
    return respond('allow');
  }
  if (!cmd.trim()) return respond('allow');

  if (/hooks_test|safety-check|check-main-write|check-ship-commands|block-secrets|block-destructive/.test(cmd)) {
    return respond('allow');
  }

  if (/\\bgit\\s+push\\b[^\\n;|&]*(-f|--force|--force-with-lease)\\b/.test(cmd) ||
      /\\bgit\\s+push\\s+(-f|--force|--force-with-lease)\\b/.test(cmd)) {
    return respond('deny', 'Spacecraft Safety Gate: force push blocked.');
  }
  if (/\\brm\\s+[^\\n;|&]*-[a-zA-Z]*r/.test(cmd) &&
      (/\\brm\\s+(?:-[a-zA-Z0-9]+\\s+)*\\/(?:\\s|$)/.test(cmd) ||
       /\\brm\\s+(?:-[a-zA-Z0-9]+\\s+)*(?:~|\\$HOME)(?:\\s|$)/.test(cmd))) {
    return respond('deny', 'Spacecraft Safety Gate: catastrophic rm blocked.');
  }
  if (/[\\s'"\`](\\.env(?:\\.[\\w.-]+)?|id_rsa|id_ed25519|credentials(?:\\.json)?)[\\s'"\`]/.test(\` \${cmd} \`) &&
      !/\\.env\\.example|\\.env\\.sample/.test(cmd)) {
    return respond('deny', 'Spacecraft Safety Gate: command references a secret path.');
  }

  const isGit = /(?:^|[;&|]|&&|\\|\\|)\\s*(?:export\\s+)?(?:[A-Za-z_][A-Za-z0-9_]*=\\S+\\s+)*(?:env\\s+)?(?:[^\\s;|&]+\\/)?git\\b/.test(cmd);
  if (!isGit) return respond('allow');

  const isMutate = /(?:^|[;&|]|&&|\\|\\|)\\s*(?:export\\s+)?(?:[A-Za-z_][A-Za-z0-9_]*=\\S+\\s+)*(?:env\\s+)?(?:[^\\s;|&]+\\/)?git\\s+(commit|merge|rebase|cherry-pick|reset|push|am|pull|tag|branch\\s+(-D|--delete\\s+-f))\\b/.test(cmd);
  if (!isMutate) return respond('allow');

  const hasShipFlag = /(?:^|[\\s])SPACECRAFT_SHIP=1(?:\\s|$)/.test(cmd) || process.env.SPACECRAFT_SHIP === '1';
  const isPush = /(?:^|[;&|]|&&|\\|\\|)\\s*(?:export\\s+)?(?:[A-Za-z_][A-Za-z0-9_]*=\\S+\\s+)*(?:env\\s+)?(?:[^\\s;|&]+\\/)?git\\s+push\\b/.test(cmd);
  const isShipCmd = /(?:^|[;&|]|&&|\\|\\|)\\s*(?:export\\s+)?(?:[A-Za-z_][A-Za-z0-9_]*=\\S+\\s+)*(?:env\\s+)?(?:[^\\s;|&]+\\/)?git\\s+(merge|push|tag)\\b/.test(cmd);

  if (isPush) {
    return respond('deny', 'Spacecraft Safety Gate: git push blocked in agent. Human must push after AUTH (Cursor ask / terminal).');
  }

  if (hasShipFlag && isShipCmd) {
    return respond('allow');
  }

  let branch = '';
  try {
    const cwd = args.Cwd || (payload.workspacePaths && payload.workspacePaths[0]) || process.cwd();
    branch = execSync('git branch --show-current', { cwd, encoding: 'utf8', stdio: ['ignore', 'pipe', 'ignore'] }).trim();
  } catch {
    // ignore
  }

  if (branch === 'main' || branch === 'master' || branch === '') {
    return respond('deny', 'Spacecraft Safety Gate: Mutating git on main/master blocked. Use feat/<id>/<title>. Ship merge/tag requires SPACECRAFT_SHIP=1.');
  }

  if (isShipCmd && !hasShipFlag) {
    return respond('deny', 'Spacecraft Safety Gate: merge/tag require /sc-ship or /sc-quick with SPACECRAFT_SHIP=1.');
  }

  return respond('allow');
}

main().catch(() => respond('allow'));
`;
fs.writeFileSync(path.join(PLUGIN_DIR, 'hooks', 'safety-check.mjs'), safetyCheckMjs, { mode: 0o755 });
console.log('  ok: generated plugins/spacecraft/hooks/safety-check.mjs');

// 5. Generate short AGENTS.md / GEMINI.md (hard contract only — skills hold depth)
const agentsRuleMd = `# Spacecraft hard contract (Antigravity)

Rules are context; hooks are enforcement. Depth lives in skills - keep this file short.

## Hard gates (hooks)

- Secrets / force-push / catastrophic rm / main mutate / ship without \`SPACECRAFT_SHIP=1\`: \`hooks/safety-check.mjs\`
- \`git push\` is **denied in-agent** (no auto-push). Human pushes after AUTH.

## Soft contract

- **AUTH:** Quoted user authorization before outward push/deploy/publish/send.
- **INTENT:** Class + intended behavior before behavior-changing edits.
- **Commander:** No product code/tests - Task-delegate (\`sc-coder\` / \`sc-tester\` / \`sc-firmware\` / \`sc-rtl\`; prose → \`sc-writer\`).
- **Lanes:** \`/sc-discuss\` → \`/sc-run\` → human check → \`/sc-ship\`. Small edits: \`/sc-quick\`.
- **SoT:** explicit user > approved draft + spec > DESIGN.md > process rules > evidence > code.
- **Language:** English technical substance; Thai for HIL / status / handoff.

See plugin skills for UX, TDD, judge, and domain encyclopedias.
`;
fs.writeFileSync(path.join(PLUGIN_DIR, 'rules', 'AGENTS.md'), agentsRuleMd);
console.log('  ok: generated plugins/spacecraft/rules/AGENTS.md');

fs.writeFileSync(path.join(ROOT, 'GEMINI.md'), agentsRuleMd);
console.log('  ok: generated workspace GEMINI.md');

// 6. Sync Subagents from .cursor/agents/ → plugins/spacecraft/agents/
// SoT copy only (no fallback stubs). Intentional Antigravity patch:
// sc-browser-probe description mentions Chrome DevTools MCP for that host.
const agentsSourceDir = path.join(ROOT, '.cursor', 'agents');
const agentsTargetDir = path.join(PLUGIN_DIR, 'agents');

const agentFiles = fs.readdirSync(agentsSourceDir).filter((f) => f.endsWith('.md'));
for (const file of agentFiles) {
  let content = fs.readFileSync(path.join(agentsSourceDir, file), 'utf8');

  if (file === 'sc-browser-probe.md') {
    content = content.replace(
      'Live browser probe + AFK fix-loop',
      'Live browser probe (Chrome DevTools MCP / browser automation) + AFK fix-loop'
    );
  }

  fs.writeFileSync(path.join(agentsTargetDir, file), content);
  console.log(`  ok: synced agent ${file}`);
}

for (const entry of fs.readdirSync(agentsTargetDir)) {
  if (!agentFiles.includes(entry)) {
    const orphan = path.join(agentsTargetDir, entry);
    fs.rmSync(orphan, { recursive: true, force: true });
    console.log(`  ok: pruned agent ${entry}`);
  }
}

// 7. Sync all skills from .cursor/skills/ → plugins/spacecraft/skills/
// Full tree per skill (SKILL.md + references/). Wipe-then-copy so deleted refs prune.
// Do not mirror .cursor/deprecated (graveyard stays local / manually curated).
const skillsSourceDir = path.join(ROOT, '.cursor', 'skills');
const skillsTargetDir = path.join(PLUGIN_DIR, 'skills');

const skillDirs = fs
  .readdirSync(skillsSourceDir)
  .filter((f) => fs.statSync(path.join(skillsSourceDir, f)).isDirectory());

for (const skill of skillDirs) {
  const src = path.join(skillsSourceDir, skill);
  const dest = path.join(skillsTargetDir, skill);
  fs.rmSync(dest, { recursive: true, force: true });
  fs.cpSync(src, dest, { recursive: true });
  console.log(`  ok: synced skill ${skill}`);
}

for (const entry of fs.readdirSync(skillsTargetDir)) {
  if (!skillDirs.includes(entry)) {
    const orphan = path.join(skillsTargetDir, entry);
    fs.rmSync(orphan, { recursive: true, force: true });
    console.log(`  ok: pruned skill ${entry}`);
  }
}

console.log('\nAntigravity sync complete! All assets are up to date.');
