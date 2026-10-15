// Cross-engine diagnostic screenshots of the demo routes (plan task M06).
// NOT a test suite: produces artifacts + an error report; never gates CI.
//
// Usage:
//   PLAYWRIGHT_BROWSERS_PATH=<nix browsers dir> node screenshot.mjs \
//     --engine chromium|firefox|webkit --base http://localhost:8093 \
//     --theme light|dark [--outdir artifacts]
//
// Engine pinning: PLAYWRIGHT_BROWSERS_PATH points at nixpkgs
// playwright-driver.browsers; the playwright-core version in package.json
// must equal that driver's version (revision numbers are hardcoded per
// release). Theme pinning rides context.addInitScript so the choice is set
// BEFORE any page script runs (no reload flash; headless Chromium defaults
// to prefers-color-scheme: dark without this).
import { chromium, firefox, webkit } from 'playwright-core';
import { readFileSync, mkdirSync, writeFileSync } from 'node:fs';
import { join, dirname } from 'node:path';
import { fileURLToPath } from 'node:url';

const here = dirname(fileURLToPath(import.meta.url));

function arg(name, fallback) {
  const i = process.argv.indexOf('--' + name);
  return i > -1 && process.argv[i + 1] ? process.argv[i + 1] : fallback;
}

const engineName = arg('engine', 'chromium');
const base = arg('base', 'http://localhost:8093');
const theme = arg('theme', 'light');
const outdir = arg('outdir', join(here, 'artifacts'));

const routes = JSON.parse(readFileSync(join(here, 'routes.json'), 'utf8'));
const engines = { chromium, firefox, webkit };
const engine = engines[engineName];
if (!engine) {
  console.error(`unknown engine ${engineName} (want chromium|firefox|webkit)`);
  process.exit(2);
}

const outRoot = join(outdir, `${engineName}-${theme}`);
mkdirSync(outRoot, { recursive: true });

const report = [];
const browser = await engine.launch();
try {
  for (const route of routes) {
    const slug = route.path === '/' ? 'home' : route.path.replaceAll('/', '_').replace(/^_/, '');
    const entry = { engine: engineName, theme, path: route.path, console: [], pageerrors: [] };
    const context = await browser.newContext({
      viewport: { width: 1280, height: 800 },
    });
    // Set the pinned theme before any page script executes.
    await context.addInitScript(
      (t) => { try { localStorage.setItem('theme', t); } catch (e) { /* private mode */ } },
      theme,
    );
    const page = await context.newPage();
    page.on('console', (msg) => {
      if (msg.type() === 'error' || msg.type() === 'warning') {
        entry.console.push(`${msg.type()}: ${msg.text().slice(0, 300)}`);
      }
    });
    page.on('pageerror', (err) => entry.pageerrors.push(String(err).slice(0, 300)));
    try {
      // 'load' + settle: SSE pages (LiveRegion) never reach networkidle.
      await page.goto(base + route.path, { waitUntil: 'load', timeout: 30000 });
      await page.waitForTimeout(1200);
      await page.screenshot({ path: join(outRoot, `${slug}.png`), fullPage: true });
      entry.ok = true;
    } catch (err) {
      entry.ok = false;
      entry.error = String(err).slice(0, 300);
    }
    await context.close();
    report.push(entry);
    const bad = entry.error || entry.pageerrors.length || entry.console.length;
    console.log(`${bad ? '!' : ' '} ${engineName} ${theme} ${route.path}${entry.error ? ' ERROR ' + entry.error : ''}`);
  }
} finally {
  await browser.close();
}

writeFileSync(join(outRoot, 'report.json'), JSON.stringify(report, null, 2));
const failed = report.filter((r) => !r.ok || r.pageerrors.length);
console.log(`done: ${report.length} routes, ${failed.length} with errors/pageerrors`);
