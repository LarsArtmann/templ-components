// M07: base-select Firefox fallback assert (plan μ27–μ28).
// NOT a test suite: diagnostic assert over the /stylable-select demo route.
// Verifies that in engines WITHOUT `appearance: base-select` support the
// Stylable selects fall back to native rendering without layout breakage.
//
// Usage:
//   PLAYWRIGHT_BROWSERS_PATH=<nix browsers dir> node stylable-assert.mjs \
//     --base http://localhost:8093 [--outdir artifacts]
//
// Selectors come from examples/demo/stylable_select_demo.templ:
//   select[name=stylable-region]  Stylable, flat options, pre-selected
//   select[name=stylable-channel] Stylable, optgroups
//   select[name=stylable-env]     Stylable, validation error
//   select[name=native-region]    native comparison select
import { chromium, firefox } from "playwright-core";
import { mkdirSync, writeFileSync } from "node:fs";
import { join, dirname } from "node:path";
import { fileURLToPath } from "node:url";

const here = dirname(fileURLToPath(import.meta.url));

function arg(name, fallback) {
  const i = process.argv.indexOf("--" + name);
  return i > -1 && process.argv[i + 1] ? process.argv[i + 1] : fallback;
}

const base = arg("base", "http://localhost:8093");
const outdir = arg("outdir", join(here, "artifacts", "m07-fallback"));
mkdirSync(outdir, { recursive: true });

const SELECTS = ["stylable-region", "stylable-channel", "stylable-env", "native-region"];

const probe = (names) => `
  (() => {
    const out = {};
    for (const name of ${JSON.stringify(SELECTS)}) {
      const el = document.querySelector('select[name="' + name + '"]');
      if (!el) { out[name] = null; continue; }
      const cs = getComputedStyle(el);
      const rect = el.getBoundingClientRect();
      out[name] = {
        appearance: cs.appearance,
        width: Math.round(rect.width),
        height: Math.round(rect.height),
        value: el.value,
        optionCount: el.options.length,
      };
    }
    out._page = {
      hScroll: document.documentElement.scrollWidth > window.innerWidth,
      scrollWidth: document.documentElement.scrollWidth,
      innerWidth: window.innerWidth,
    };
    return out;
  })()
`;

const report = [];
for (const [engineName, engine] of [["chromium", chromium], ["firefox", firefox]]) {
  const browser = await engine.launch();
  const context = await browser.newContext({ viewport: { width: 1280, height: 800 } });
  await context.addInitScript((t) => {
    try { localStorage.setItem("theme", t); } catch (e) {}
  }, "light");
  const page = await context.newPage();
  await page.goto(base + "/stylable-select", { waitUntil: "load" });
  await page.waitForTimeout(800);

  const result = await page.evaluate(probe(SELECTS));
  const entry = { engine: engineName, selects: result, layoutOk: !result._page.hScroll };

  // μ27: per-select screenshots of the flat Stylable select + native counterpart.
  for (const name of ["stylable-region", "native-region"]) {
    const el = page.locator(`select[name="${name}"]`);
    await el.screenshot({ path: join(outdir, `${engineName}-${name}.png`) });
  }

  // Functional fallback check: selectOption still works in both dialects.
  await page.selectOption('select[name="stylable-region"]', "us-east");
  entry.functionalPick = (await page.inputValue('select[name="stylable-region"]')) === "us-east";

  report.push(entry);
  await browser.close();
}

writeFileSync(join(outdir, "report.json"), JSON.stringify(report, null, 2));
for (const r of report) {
  console.log(`--- ${r.engine} (layoutOk=${r.layoutOk}, functionalPick=${r.functionalPick})`);
  for (const [name, v] of Object.entries(r.selects)) {
    if (name === "_page") continue;
    if (!v) { console.log(`  ${name}: MISSING`); continue; }
    console.log(`  ${name}: appearance=${v.appearance} ${v.width}x${v.height} value=${v.value} options=${v.optionCount}`);
  }
}
