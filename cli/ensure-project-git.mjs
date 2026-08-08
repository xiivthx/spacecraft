#!/usr/bin/env node
/**
 * Thin CLI entry so shell installers share project-git ensure helpers.
 *
 * Usage: ensure-project-git.mjs --ready|--ignore <project-root>
 */
import path from 'node:path';
import { ensureProjectReady, ensureSpaceIgnored } from './lib/project-git.mjs';

const mode = process.argv[2];
const projectRoot = process.argv[3];

if (
  (mode !== '--ready' && mode !== '--ignore') ||
  !projectRoot ||
  process.argv.length !== 4
) {
  console.error('usage: ensure-project-git.mjs --ready|--ignore <project-root>');
  process.exit(1);
}

const root = path.resolve(projectRoot);

try {
  if (mode === '--ready') {
    ensureProjectReady(root);
  } else {
    ensureSpaceIgnored(root);
  }
} catch (err) {
  console.error(`ensure-project-git: ${err.message}`);
  process.exit(1);
}
