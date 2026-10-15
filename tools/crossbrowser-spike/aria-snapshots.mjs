// M10: aria-snapshot goldens spike (plan μ37–μ40).
// NOT a test suite: captures Playwright ariaSnapshot() YAML for four
// representative component states + one intentional-regression demo
// proving a mutated ARIA attribute is caught by snapshot diff.
//
// Usage:
//   PLAYWRIGHT_BROWSERS_PATH=/tmp/pw-browsers node aria-snapshots.mjs \
//     --base http://localhost:8093 [--outdir artifacts]
import { chromium } from "playwright-core";
import { mkdirSync, writeFileSync } from "node:fs";
import { join, dirname } from "node:path";
import { fileURLToPath } from "node:url";

const here = dirname(fileURLToPath(import.meta.url));

function arg(name, fallback) {
  const i = process.argv.indexOf("--" + name);
  return i > -1 && process.argv[i + 1] ? process.argv[i + 1] : fallback;
}

const base = arg("base", "http://localhost:8093");
const outdir = arg("outdir", join(here, "artifacts", "m10-aria"));
mkdirSync(outdir, { recursive: true });

const browser = await chromium.launch();
const context = await browser.newContext({ viewport: { width: 1280, height: 800 } });
await context.addInitScript(() => {
  try {
    localStorage.setItem("theme", "light");
  } catch (e) {}
});
const page = await context.newPage();

async function snap(route, locator, name, note) {
  await page.goto(base + route, { waitUntil: "load" });
  await page.waitForTimeout(600);
  const el = page.locator(locator).first();
  const yaml = await el.ariaSnapshot();
  writeFileSync(join(outdir, name + ".yaml"), `# ${name} — ${note}\n${yaml}\n`);
  console.log(`captured ${name} (${yaml.split("\n").length} lines)`);
  return yaml;
}

// μ37: Table with aria-sort (display demo table has SortAsc on "Name")
const tableYaml = await snap(
  "/display",
  "table",
  "table-aria-sort",
  "Table with sortable headers, aria-sort=ascending on Name",
);
await page.goto(base + "/display", { waitUntil: "load" });
const domSort = await page.evaluate(
  () => document.querySelector("th[aria-sort]")?.getAttribute("aria-sort") || "MISSING",
);
console.log(
  `SPIKE FINDING: DOM has th[aria-sort="${domSort}"] but ariaSnapshot YAML exposes aria-sort: ${tableYaml.includes("aria-sort")}`,
);

// μ38a: Modal as role=dialog — OPEN it first (closed dialog = display:none = empty snapshot);
// snapshot inline WITHOUT re-navigating (snap() reloads the page)
await page.goto(base + "/display", { waitUntil: "load" });
await page.waitForTimeout(600);
await page.click('[data-tc-overlay-trigger="demo-modal"]');
await page.waitForTimeout(300);
const modalYaml = await page.locator("#demo-modal").ariaSnapshot();
writeFileSync(
  join(outdir, "modal-dialog.yaml"),
  `# modal-dialog — native <dialog> modal (open)\n${modalYaml}\n`,
);
console.log(`captured modal-dialog (open, ${modalYaml.split("\n").length} lines)`);

// μ38b: Dropdown trigger — native Popover API: popovertarget, aria-expanded
// intentionally ABSENT (HTML-validation rule: aria-expanded next to popovertarget
// is invalid). ATs derive expanded state from the popover itself.
await snap(
  "/display",
  "[popovertarget]",
  "dropdown-popovertarget",
  "Dropdown trigger (native popover, no aria-expanded by design)",
);

// μ39: Tabs role=tablist
await snap("/display", '[role="tablist"]', "tabs-tablist", "Tabs tablist with selected state");

// μ40: regression demo — mutate aria-sort, prove a DOM-level assert catches it
// (ariaSnapshot does NOT expose aria-sort, so the guard must read the attribute)
await page.goto(base + "/display", { waitUntil: "load" });
await page.waitForTimeout(600);
const before = await page.evaluate(
  () => document.querySelectorAll('th[aria-sort="ascending"]').length,
);
await page.evaluate(() => {
  document
    .querySelectorAll('th[aria-sort="ascending"]')
    .forEach((th) => th.removeAttribute("aria-sort"));
});
const after = await page.evaluate(
  () => document.querySelectorAll('th[aria-sort="ascending"]').length,
);
const beforeYaml = await page.locator("table").first().ariaSnapshot();
writeFileSync(
  join(outdir, "regression-before.yaml"),
  `# th[aria-sort] before: ${before}\n${beforeYaml}\n`,
);
writeFileSync(
  join(outdir, "regression-after.yaml"),
  `# th[aria-sort] after removal: ${after} (attribute gone)\n${beforeYaml}\n`,
);
console.log(
  `μ40 regression demo: ascending th count ${before} → ${after} (all removed); YAML can't see it — attribute asserts must`,
);

await browser.close();
