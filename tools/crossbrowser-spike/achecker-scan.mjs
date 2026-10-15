// M08: IBM accessibility-checker spike over the demo routes (plan μ30–μ33).
// NOT a test suite: diagnostic second a11y opinion; findings go to their OWN
// ledger semantics (M20), never into axe_baseline.json.
//
// Usage:
//   PUPPETEER_EXECUTABLE_PATH=/tmp/pw-browsers/chromium-1243/chrome-linux64/chrome \
//   node achecker-scan.mjs --base http://localhost:8093 [--outdir artifacts]
import { createRequire } from "node:module";
import { readFileSync, mkdirSync, writeFileSync } from "node:fs";
import { join, dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const require = createRequire(import.meta.url);
const aChecker = require("accessibility-checker");

const here = dirname(fileURLToPath(import.meta.url));

function arg(name, fallback) {
  const i = process.argv.indexOf("--" + name);
  return i > -1 && process.argv[i + 1] ? process.argv[i + 1] : fallback;
}

const base = arg("base", "http://localhost:8093");
const outdir = arg("outdir", join(here, "artifacts", "m08-achecker"));
mkdirSync(outdir, { recursive: true });

const routes = JSON.parse(readFileSync(join(here, "routes.json"), "utf8"));
const axeBaseline = JSON.parse(
  readFileSync(resolve(here, "../../visualtest/testdata/axe_baseline.json"), "utf8"),
);

const perRoute = [];
const ruleTotals = {};
for (const route of routes) {
  const url = base + route.path;
  const r = await aChecker.getCompliance(url, "m08");
  const results = (r.report && r.report.results) || [];
  const byRule = {};
  for (const res of results) {
    if (res.ignored) continue;
    const key = res.ruleId;
    byRule[key] = byRule[key] || {
      ruleId: res.ruleId,
      level: res.level,
      message: res.message,
      count: 0,
      sample: null,
    };
    byRule[key].count++;
    if (!byRule[key].sample) byRule[key].sample = { dom: res.path && res.path.dom, snippet: res.snippet && String(res.snippet).slice(0, 160) };
    ruleTotals[key] = ruleTotals[key] || { ruleId: res.ruleId, level: res.level, message: res.message, routes: 0, count: 0 };
    if (byRule[key].count === 1) ruleTotals[key].routes++;
    ruleTotals[key].count++;
  }
  const findings = Object.values(byRule);
  perRoute.push({ path: route.path, findingClasses: findings.length, violations: results.filter((x) => x.level === "violation").length, needsReview: results.filter((x) => x.level === "needsreview").length, recommendations: results.filter((x) => x.level === "recommendation").length, rules: findings });
  console.log(`${route.path}: ${findings.length} rule classes, ${results.length} results`);
}

writeFileSync(join(outdir, "report.json"), JSON.stringify({ perRoute, ruleTotals }, null, 2));

console.log("\n=== rule totals across 16 routes (ruleId | level | routes | count) ===");
for (const t of Object.values(ruleTotals).sort((a, b) => b.count - a.count)) {
  console.log(`${t.ruleId} | ${t.level} | ${t.routes} | ${t.count}`);
}
console.log("\n=== axe baseline entries (for the overlap table) ===");
for (const [route, rules] of Object.entries(axeBaseline)) {
  for (const [rule, info] of Object.entries(rules)) {
    console.log(`${route}: ${rule} (${(info && (info.count ?? JSON.stringify(info).slice(0, 40))) || "?"})`);
  }
}
