#!/usr/bin/env node
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
  const base = p.split(/[\\/]/).pop() || '';
  if (/^\.env\.example$/i.test(base) || /\.env\.sample$/i.test(base)) return false;
  if (/^\.env$/i.test(base) || /^\.env\./i.test(base)) return true;
  if (/(^|[/\\])(credentials|id_rsa|id_ed25519)(\.|$)/i.test(p)) return true;
  if (/\.(pem|p12|pfx|key)$/i.test(base) && !/\.pub$/i.test(base)) return true;
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

  if (/\bgit\s+push\b[^\n;|&]*(-f|--force|--force-with-lease)\b/.test(cmd) ||
      /\bgit\s+push\s+(-f|--force|--force-with-lease)\b/.test(cmd)) {
    return respond('deny', 'Spacecraft Safety Gate: force push blocked.');
  }
  if (/\brm\s+[^\n;|&]*-[a-zA-Z]*r/.test(cmd) &&
      (/\brm\s+(?:-[a-zA-Z0-9]+\s+)*\/(?:\s|$)/.test(cmd) ||
       /\brm\s+(?:-[a-zA-Z0-9]+\s+)*(?:~|\$HOME)(?:\s|$)/.test(cmd))) {
    return respond('deny', 'Spacecraft Safety Gate: catastrophic rm blocked.');
  }
  if (/[\s'"`](\.env(?:\.[\w.-]+)?|id_rsa|id_ed25519|credentials(?:\.json)?)[\s'"`]/.test(` ${cmd} `) &&
      !/\.env\.example|\.env\.sample/.test(cmd)) {
    return respond('deny', 'Spacecraft Safety Gate: command references a secret path.');
  }

  const isGit = /(?:^|[;&|]|&&|\|\|)\s*(?:export\s+)?(?:[A-Za-z_][A-Za-z0-9_]*=\S+\s+)*(?:env\s+)?(?:[^\s;|&]+\/)?git\b/.test(cmd);
  if (!isGit) return respond('allow');

  const isMutate = /(?:^|[;&|]|&&|\|\|)\s*(?:export\s+)?(?:[A-Za-z_][A-Za-z0-9_]*=\S+\s+)*(?:env\s+)?(?:[^\s;|&]+\/)?git\s+(commit|merge|rebase|cherry-pick|reset|push|am|pull|tag|branch\s+(-D|--delete\s+-f))\b/.test(cmd);
  if (!isMutate) return respond('allow');

  const hasShipFlag = /(?:^|[\s])SPACECRAFT_SHIP=1(?:\s|$)/.test(cmd) || process.env.SPACECRAFT_SHIP === '1';
  const isPush = /(?:^|[;&|]|&&|\|\|)\s*(?:export\s+)?(?:[A-Za-z_][A-Za-z0-9_]*=\S+\s+)*(?:env\s+)?(?:[^\s;|&]+\/)?git\s+push\b/.test(cmd);
  const isShipCmd = /(?:^|[;&|]|&&|\|\|)\s*(?:export\s+)?(?:[A-Za-z_][A-Za-z0-9_]*=\S+\s+)*(?:env\s+)?(?:[^\s;|&]+\/)?git\s+(merge|push|tag)\b/.test(cmd);

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
