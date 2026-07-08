#!/usr/bin/env node
import { spawn } from "node:child_process";
import { createReadStream } from "node:fs";
import fs from "node:fs/promises";
import http from "node:http";
import path from "node:path";
import process from "node:process";

const ROOT = process.cwd();
const DEFAULT_HOST = "127.0.0.1";
const DEFAULT_PORT = 4321;

const MIME_TYPES = new Map([
  [".html", "text/html; charset=utf-8"],
  [".css", "text/css; charset=utf-8"],
  [".js", "text/javascript; charset=utf-8"],
  [".mjs", "text/javascript; charset=utf-8"],
  [".json", "application/json; charset=utf-8"],
  [".svg", "image/svg+xml"],
  [".png", "image/png"],
  [".jpg", "image/jpeg"],
  [".jpeg", "image/jpeg"],
  [".gif", "image/gif"],
  [".webp", "image/webp"],
  [".ico", "image/x-icon"],
  [".txt", "text/plain; charset=utf-8"],
  [".woff", "font/woff"],
  [".woff2", "font/woff2"]
]);

function usage() {
  return `Design HTML server

Usage:
  node .opencode/skills/sc-design/scripts/serve-html.mjs [html-file-or-dir] [options]

Options:
  --open             Open the preview URL in the default browser.
  --host <host>      Host to bind. Default: ${DEFAULT_HOST}
  --port <port>      Port to bind. Default: ${DEFAULT_PORT}; auto-increments if busy.
  -h, --help         Show this help.

Examples:
  node .opencode/skills/sc-design/scripts/serve-html.mjs --open
  node .opencode/skills/sc-design/scripts/serve-html.mjs okinawa-ui-directions.html --open
  node .opencode/skills/sc-design/scripts/serve-html.mjs .space/missions/M-123/design/example.html --port 4330
`;
}

function fail(message, code = 1) {
  const error = new Error(message);
  error.exitCode = code;
  throw error;
}

function parseArgs(argv) {
  const options = {
    host: DEFAULT_HOST,
    port: DEFAULT_PORT,
    open: false,
    portExplicit: false,
    target: null
  };

  for (let index = 0; index < argv.length; index += 1) {
    const arg = argv[index];
    if (arg === "--open") {
      options.open = true;
      continue;
    }
    if (arg === "--host") {
      options.host = argv[index + 1];
      index += 1;
      continue;
    }
    if (arg === "--port") {
      options.port = Number(argv[index + 1]);
      options.portExplicit = true;
      index += 1;
      continue;
    }
    if (arg === "-h" || arg === "--help" || arg === "help") {
      options.help = true;
      continue;
    }
    if (arg.startsWith("-")) {
      fail(`Unknown option "${arg}".\n\n${usage()}`);
    }
    if (options.target) {
      fail(`Only one html file or directory can be served.\n\n${usage()}`);
    }
    options.target = arg;
  }

  if (!options.host) {
    fail("Missing value for --host.");
  }
  if (!Number.isInteger(options.port) || options.port < 0 || options.port > 65535) {
    fail("Invalid --port value. Use an integer from 0 to 65535.");
  }

  return options;
}

async function exists(filePath) {
  try {
    await fs.access(filePath);
    return true;
  } catch {
    return false;
  }
}

async function currentMissionId() {
  const currentFile = path.join(ROOT, ".space", "current");
  if (!(await exists(currentFile))) {
    return null;
  }
  const value = (await fs.readFile(currentFile, "utf8")).trim();
  return value || null;
}

async function currentDesignDir() {
  const id = await currentMissionId();
  if (!id) {
    fail("No current mission. Pass an HTML file or run /sc-start first.");
  }
  return path.join(ROOT, ".space", "missions", id, "design");
}

async function resolveTarget(target) {
  if (!target) {
    return currentDesignDir();
  }

  const direct = path.resolve(ROOT, target);
  if (await exists(direct)) {
    return direct;
  }

  if (!path.isAbsolute(target)) {
    const designDir = await currentDesignDir();
    const designRelative = path.join(designDir, target);
    if (await exists(designRelative)) {
      return designRelative;
    }
  }

  fail(`Cannot find HTML file or directory: ${target}`);
}

function toUrlPath(filePath, rootDir) {
  const relative = path.relative(rootDir, filePath);
  if (!relative) {
    return "/";
  }
  return `/${relative.split(path.sep).map(encodeURIComponent).join("/")}`;
}

async function chooseOpenPath(targetPath, stat, rootDir) {
  if (stat.isFile()) {
    return toUrlPath(targetPath, rootDir);
  }

  const indexPath = path.join(targetPath, "index.html");
  if (await exists(indexPath)) {
    return toUrlPath(indexPath, rootDir);
  }

  return toUrlPath(targetPath, rootDir);
}

function escapeHtml(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;");
}

function formatBytes(bytes) {
  if (bytes < 1024) {
    return `${bytes} B`;
  }
  if (bytes < 1024 * 1024) {
    return `${(bytes / 1024).toFixed(1)} KB`;
  }
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

async function directoryListing(dirPath, urlPath) {
  const entries = await fs.readdir(dirPath, { withFileTypes: true });
  const rows = await Promise.all(entries
    .filter((entry) => !entry.name.startsWith("."))
    .sort((left, right) => {
      if (left.isDirectory() !== right.isDirectory()) {
        return left.isDirectory() ? -1 : 1;
      }
      return left.name.localeCompare(right.name);
    })
    .map(async (entry) => {
      const filePath = path.join(dirPath, entry.name);
      const stat = await fs.stat(filePath);
      const hrefBase = urlPath.endsWith("/") ? urlPath : `${urlPath}/`;
      const href = `${hrefBase}${encodeURIComponent(entry.name)}${entry.isDirectory() ? "/" : ""}`;
      const type = entry.isDirectory() ? "folder" : path.extname(entry.name).replace(".", "") || "file";
      return `<a class="row" href="${href}">
        <span class="name">${escapeHtml(entry.name)}${entry.isDirectory() ? "/" : ""}</span>
        <span>${escapeHtml(type)}</span>
        <span>${entry.isDirectory() ? "-" : formatBytes(stat.size)}</span>
        <span>${escapeHtml(stat.mtime.toLocaleString())}</span>
      </a>`;
    }));

  return `<!doctype html>
<html lang="th">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Design Preview</title>
  <style>
    :root {
      color-scheme: light dark;
      --bg: #f6f4ee;
      --text: #172018;
      --muted: #647064;
      --line: #d8d3c7;
      --accent: #1d7f6e;
      --panel: #fffdfa;
    }
    @media (prefers-color-scheme: dark) {
      :root {
        --bg: #111613;
        --text: #edf2ea;
        --muted: #a6b0a5;
        --line: #2c342f;
        --accent: #55c3a8;
        --panel: #171d19;
      }
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      min-height: 100vh;
      background: var(--bg);
      color: var(--text);
      font: 15px/1.5 ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
    }
    main {
      width: min(1040px, calc(100vw - 32px));
      margin: 48px auto;
    }
    h1 {
      margin: 0 0 8px;
      font-size: 40px;
      line-height: 1.05;
      letter-spacing: 0;
    }
    p {
      margin: 0 0 24px;
      color: var(--muted);
    }
    .list {
      overflow: hidden;
      border: 1px solid var(--line);
      border-radius: 8px;
      background: var(--panel);
    }
    .row {
      display: grid;
      grid-template-columns: minmax(160px, 1fr) 110px 110px 190px;
      gap: 16px;
      align-items: center;
      padding: 12px 16px;
      color: inherit;
      text-decoration: none;
      border-top: 1px solid var(--line);
    }
    .row:first-child { border-top: 0; }
    .row:hover, .row:focus-visible {
      background: color-mix(in srgb, var(--accent) 10%, transparent);
      outline: none;
    }
    .name {
      color: var(--accent);
      font-weight: 650;
      overflow-wrap: anywhere;
    }
    @media (max-width: 720px) {
      main { width: min(100% - 24px, 1040px); margin: 28px auto; }
      h1 { font-size: 28px; }
      .row { grid-template-columns: 1fr; gap: 4px; }
    }
  </style>
</head>
<body>
  <main>
    <h1>ไฟล์ Design (Preview)</h1>
    <p>เลือกไฟล์ HTML เพื่อเปิดดูตัวอย่างผ่าน local server.</p>
    <section class="list" aria-label="Design files">
      ${rows.join("\n")}
    </section>
  </main>
</body>
</html>`;
}

async function sendFile(response, filePath, stat) {
  const contentType = MIME_TYPES.get(path.extname(filePath).toLowerCase()) || "application/octet-stream";
  response.writeHead(200, {
    "content-type": contentType,
    "content-length": stat.size,
    "cache-control": "no-store"
  });
  createReadStream(filePath).pipe(response);
}

async function handleRequest(rootDir, request, response) {
  try {
    const requestUrl = new URL(request.url || "/", "http://localhost");
    const decodedPath = decodeURIComponent(requestUrl.pathname);
    const relativePath = decodedPath.replace(/^\/+/, "");
    const requestedPath = path.resolve(rootDir, relativePath);
    const rootBoundary = rootDir.endsWith(path.sep) ? rootDir : `${rootDir}${path.sep}`;

    if (requestedPath !== rootDir && !requestedPath.startsWith(rootBoundary)) {
      response.writeHead(403, { "content-type": "text/plain; charset=utf-8" });
      response.end("Forbidden");
      return;
    }

    const stat = await fs.stat(requestedPath).catch(() => null);
    if (!stat) {
      response.writeHead(404, { "content-type": "text/plain; charset=utf-8" });
      response.end("Not found");
      return;
    }

    if (stat.isDirectory()) {
      const indexPath = path.join(requestedPath, "index.html");
      const indexStat = await fs.stat(indexPath).catch(() => null);
      if (indexStat?.isFile()) {
        await sendFile(response, indexPath, indexStat);
        return;
      }

      const html = await directoryListing(requestedPath, requestUrl.pathname);
      response.writeHead(200, {
        "content-type": "text/html; charset=utf-8",
        "cache-control": "no-store"
      });
      response.end(html);
      return;
    }

    await sendFile(response, requestedPath, stat);
  } catch (error) {
    response.writeHead(500, { "content-type": "text/plain; charset=utf-8" });
    response.end(`Server error: ${error.message}`);
  }
}

function listen(server, host, port) {
  return new Promise((resolve, reject) => {
    server.once("error", reject);
    server.listen(port, host, () => {
      server.removeListener("error", reject);
      resolve(server.address());
    });
  });
}

async function startServer(rootDir, host, port, portExplicit) {
  const attempts = portExplicit || port === 0 ? 1 : 20;
  for (let attempt = 0; attempt < attempts; attempt += 1) {
    const candidatePort = port === 0 ? 0 : port + attempt;
    const server = http.createServer((request, response) => {
      void handleRequest(rootDir, request, response);
    });

    try {
      const address = await listen(server, host, candidatePort);
      return { server, address };
    } catch (error) {
      server.close();
      if (error.code !== "EADDRINUSE" || portExplicit) {
        throw error;
      }
    }
  }

  fail(`No available port found from ${port} to ${port + attempts - 1}.`);
}

function openBrowser(url) {
  const platform = process.platform;
  const command = platform === "darwin" ? "open" : platform === "win32" ? "cmd" : "xdg-open";
  const args = platform === "win32" ? ["/c", "start", "", url] : [url];
  const child = spawn(command, args, {
    detached: true,
    stdio: "ignore"
  });
  child.unref();
}

async function main() {
  const options = parseArgs(process.argv.slice(2));
  if (options.help) {
    console.log(usage());
    return;
  }

  const targetPath = await resolveTarget(options.target);
  const stat = await fs.stat(targetPath);
  if (!stat.isDirectory() && !stat.isFile()) {
    fail(`Target must be an HTML file or directory: ${targetPath}`);
  }
  if (stat.isFile() && path.extname(targetPath).toLowerCase() !== ".html") {
    fail(`Target file must be .html: ${targetPath}`);
  }

  const rootDir = stat.isDirectory() ? targetPath : path.dirname(targetPath);
  const openPath = await chooseOpenPath(targetPath, stat, rootDir);
  const { server, address } = await startServer(rootDir, options.host, options.port, options.portExplicit);
  const actualPort = typeof address === "object" && address ? address.port : options.port;
  const url = `http://${options.host}:${actualPort}${openPath}`;

  console.log("Design HTML server");
  console.log(`Root: ${path.relative(ROOT, rootDir) || "."}`);
  console.log(`URL: ${url}`);
  console.log("Press Ctrl+C to stop.");

  if (options.open) {
    openBrowser(url);
  }

  process.on("SIGINT", () => {
    server.close(() => {
      console.log("\nStopped design HTML server.");
      process.exit(0);
    });
  });
}

main().catch((error) => {
  console.error(`Design server error: ${error.message}`);
  process.exit(error.exitCode ?? 1);
});
