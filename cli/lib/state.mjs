import { readFileSync, writeFileSync } from 'node:fs';
import path from 'node:path';
import { missionDir, missionExists, resolveActive } from './resolve.mjs';

const VALID_STATES = new Set([
  'active',
  'planned',
  'in_progress',
  'ready',
  'blocked',
  'shipped',
]);

const VALID_TRANSITIONS = {
  active: ['planned', 'blocked'],
  planned: ['in_progress', 'blocked'],
  in_progress: ['ready', 'blocked'],
  ready: ['shipped', 'blocked'],
  blocked: ['active', 'in_progress'],
  shipped: [],
};

const VALID_CLARIFY = new Set(['open', 'clear', 'deferred']);

function writeJSON(filePath, value) {
  writeFileSync(filePath, `${JSON.stringify(value, null, 2)}\n`);
}

function readMissionJSON(spaceDir, id) {
  const data = readFileSync(path.join(missionDir(spaceDir, id), 'mission.json'), 'utf8');
  return JSON.parse(data);
}

export function stateCmd(args, spaceDir, mid) {
  let missionID = '';
  let newState = '';

  switch (args.length) {
    case 1: {
      newState = args[0];
      if (!VALID_STATES.has(newState)) {
        console.error(`spacecraft state: invalid state ${JSON.stringify(newState)}`);
        return 1;
      }
      missionID = resolveActive(spaceDir, mid);
      if (!missionID) {
        console.error(
          "spacecraft state: no active mission - provide mission-id or select one with 'spacecraft use'",
        );
        return 1;
      }
      break;
    }
    case 2:
      missionID = args[0];
      newState = args[1];
      break;
    default:
      console.error('Usage: spacecraft set-state [mission-id] <new-state>');
      console.error('Valid states: active → planned → in_progress → ready → shipped');
      return 1;
  }

  if (!VALID_STATES.has(newState)) {
    console.error(`spacecraft state: invalid state ${JSON.stringify(newState)}`);
    return 1;
  }

  if (!missionExists(spaceDir, missionID)) {
    console.error(`spacecraft state: mission not found: ${missionID}`);
    return 1;
  }

  let m;
  try {
    m = readMissionJSON(spaceDir, missionID);
  } catch (err) {
    console.error(`spacecraft state: invalid mission.json: ${err.message}`);
    return 1;
  }

  const oldState = typeof m.state === 'string' ? m.state : '';
  if (oldState === newState) {
    console.log(`${missionID} already ${newState} - no change`);
    return 0;
  }

  if (oldState !== '') {
    const allowed = VALID_TRANSITIONS[oldState] ?? [];
    if (!allowed.includes(newState)) {
      console.error(`spacecraft state: invalid transition ${oldState} → ${newState}`);
      console.error(`Allowed: ${JSON.stringify(allowed)}`);
      return 1;
    }
  }

  m.state = newState;
  try {
    writeJSON(path.join(missionDir(spaceDir, missionID), 'mission.json'), m);
  } catch (err) {
    console.error('spacecraft state:', err.message);
    return 1;
  }

  console.log(`${missionID}: ${oldState} → ${newState}`);
  return 0;
}

export function clarifyStatusCmd(args, spaceDir, mid) {
  if (args.length === 0) {
    console.error('Usage: spacecraft clarify-status <open|clear|deferred>');
    return 1;
  }
  const status = args[0];
  if (!VALID_CLARIFY.has(status)) {
    console.error(
      `spacecraft clarify-status: invalid status ${JSON.stringify(status)} (open|clear|deferred)`,
    );
    return 1;
  }
  const id = resolveActive(spaceDir, mid);
  if (!id) {
    console.error(
      "spacecraft clarify-status: no active mission - pass a mission via branch or 'use'",
    );
    return 1;
  }
  try {
    writeFileSync(path.join(missionDir(spaceDir, id), 'clarify-status'), `${status}\n`);
  } catch (err) {
    console.error('spacecraft clarify-status:', err.message);
    return 1;
  }
  console.log(`Mission ${id} clarification: ${status}`);
  return 0;
}
