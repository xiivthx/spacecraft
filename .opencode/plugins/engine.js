import path from 'path';
import fs from 'fs';
import { fileURLToPath } from 'url';

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const projectRoot = path.resolve(__dirname, '../..');

const engineDir = path.join(projectRoot, '.engine');
const skillsDir = path.join(engineDir, 'skills');
const commandsDir = path.join(engineDir, 'commands');

let _bootstrapCache = undefined;

const readFileIfExists = (filePath) => {
  try {
    return fs.readFileSync(filePath, 'utf8');
  } catch {
    return '';
  }
};

const parseFrontmatter = (content) => {
  const match = content.match(/^---\n([\s\S]*?)\n---\n([\s\S]*)$/);
  if (!match) return { frontmatter: {}, body: content };

  const frontmatter = {};
  for (const line of match[1].split('\n')) {
    const colonIdx = line.indexOf(':');
    if (colonIdx < 0) continue;
    const key = line.slice(0, colonIdx).trim();
    let value = line.slice(colonIdx + 1).trim().replace(/^["']|["']$/g, '');

    if (value === 'true' || value === 'false') {
      value = value === 'true';
    }

    frontmatter[key] = value;
  }

  return { frontmatter, body: match[2] };
};

const getBootstrapContent = () => {
  if (_bootstrapCache !== undefined) return _bootstrapCache;

  const persona = readFileIfExists(path.join(engineDir, 'PERSONA.md'));
  const agents = readFileIfExists(path.join(engineDir, 'AGENTS.md'));
  const design = readFileIfExists(path.join(engineDir, 'DESIGN.md'));

  const parts = [];
  if (persona) parts.push(persona);
  if (agents) parts.push(agents);
  if (design) parts.push(design);

  if (parts.length === 0) {
    _bootstrapCache = null;
    return null;
  }

  _bootstrapCache = `<EXTREMELY_IMPORTANT>
${parts.join('\n\n---\n\n')}
</EXTREMELY_IMPORTANT>`;

  return _bootstrapCache;
};

export const EnginePlugin = async ({ client, directory }) => {
  return {
    config: async (config) => {
      config.skills = config.skills || {};
      config.skills.paths = config.skills.paths || [];
      if (!config.skills.paths.includes(skillsDir)) {
        config.skills.paths.push(skillsDir);
      }

      config.command = config.command || {};
      try {
        const files = fs.readdirSync(commandsDir);
        for (const file of files) {
          if (!file.endsWith('.md')) continue;
          const name = file.replace(/\.md$/, '');
          const content = readFileIfExists(path.join(commandsDir, file));
          if (!content) continue;

          const { frontmatter, body } = parseFrontmatter(content);

          const commandDef = {
            template: body.trim() || content.trim(),
            description: frontmatter.description || '',
          };

          if (frontmatter.agent) commandDef.agent = frontmatter.agent;
          if (frontmatter.subtask !== undefined) commandDef.subtask = frontmatter.subtask;
          if (frontmatter.model) commandDef.model = frontmatter.model;

          config.command[name] = commandDef;
        }
      } catch {}
    },

    'experimental.chat.messages.transform': async (_input, output) => {
      const bootstrap = getBootstrapContent();
      if (!bootstrap || !output.messages.length) return;

      const firstUser = output.messages.find(m => m.info?.role === 'user');
      if (!firstUser || !firstUser.parts?.length) return;

      if (firstUser.parts.some(p => p.type === 'text' && p.text.includes('EXTREMELY_IMPORTANT'))) return;

      const ref = firstUser.parts[0];
      firstUser.parts.unshift({ ...ref, type: 'text', text: bootstrap });
    }
  };
};
