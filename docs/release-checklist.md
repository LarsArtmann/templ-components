# Release Checklist — templ-components

The operational, lesson-forged checklist for cutting a release. `scripts/release.sh`
automates the mechanics; this page is the human audit trail: what the script does,
which incident each hardening step exists because of, and what must be verified by
hand before pushing.

**Sources:** the v1.9.0 (`unknown revision` verify trap), v1.10.0 (daemon race +
stale CSS + leaked self-replace), and v1.11.0 (9-day red master from a skipped
re-tidy) release cuts. The script encodes all of them — read this page to
understand _why_, so a refactor never quietly removes a guard.

---

## Before you cut

- [ ] **Master is green.** Check CI + Website workflows on the tip commit. A cut
      from a red master ships the redness in a tag (immutable once pushed).
- [ ] **Working tree is clean and on `master`.** The script enforces this, but
      check first so you don't discover it after staging notes.
- [ ] **`[Unreleased]` in CHANGELOG.md is warm.** The script refuses to cut with
      an empty section. Entries land with their feature/fix commits, never at
      release time.
- [ ] **The daemon is not mid-commit.** The auto-commit daemon has raced a cut
      (v1.10.0): it committed the in-flight version bumps and pushed tags before
      the script finished. `git log --oneline -3` — if the daemon just committed,
      wait or pause it before starting.
- [ ] **Version number sanity.** New > current (`sort -V`), correct semver
      increment for the changes (breaking changes → major per the v2 policy).

## What release.sh does, step by step

| Step | What                                                                                                                   | Incident it exists because of                                                                                                                                                                                                      |
| ---- | ---------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| 1–3  | Validate clean tree, master, version ordering                                                                          | routine                                                                                                                                                                                                                            |
| 4    | Collect notes from `--notes-file` or `[Unreleased]`                                                                    | hostile stdin prompts                                                                                                                                                                                                              |
| —    | Install ONE `EXIT` trap that restores version files, go.mods, CHANGELOG, FEATURES on failure                           | v1.9.0: a second `trap` silently disabled the rollback trap, leaving a dirty tree                                                                                                                                                  |
| 5    | Bump `utils.Version`                                                                                                   | —                                                                                                                                                                                                                                  |
| 5b   | Bump sibling `require` pins in every published go.mod **and `visualtest/go.mod`**                                      | v1.11.0: visualtest pins were bumped by hand after push; the gap left master red for 9 days                                                                                                                                        |
| 6/6b | Move `[Unreleased]` under the new heading; bump FEATURES.md version+date                                               | drift guards `TestVersionMatches(Changelog\|Features)`                                                                                                                                                                             |
| 7    | Full verify: `templ generate` + CSS recompile + build/test/lint, **with local `replace` directives still present**     | v1.9.0: stripping replaces before verify made every build fail with `unknown revision <sub>/v<version>` — go1.26 workspace mode still resolves `require` entries at unpushed versions; proven with a 2×2 replaces/GOPRIVATE matrix |
| 7b   | Strip replaces from the go.mods being tagged (consumer-clean tags)                                                     | v1.10.0: the root self-replace leaked into the tag; the strip pattern matches `[ /]` to catch both                                                                                                                                 |
| —    | Parse-check every stripped go.mod (`go mod edit -json`)                                                                | —                                                                                                                                                                                                                                  |
| 8    | Commit `release: <version> — <summary>` with notes in the body                                                         | one-commit release convention                                                                                                                                                                                                      |
| 8b   | Assert the COMMITTED TREE: version files agree, no replaces in any tagged go.mod                                       | v1.10.0: the daemon split the cut across commits, so the script's commit didn't carry the bumps. Tags are immutable — verify the tree, not the diff                                                                                |
| 9    | Annotated SSH-signed tags: root `v<x>` + one `<sub>/v<x>` per published module, derived from root go.mod's replace set | v1.8.3: a root-only release broke every consumer; the hardcoded module list had drifted from reality                                                                                                                               |
| 10   | Re-add replaces (follow-up commit) for local dev                                                                       | remove-at-release strategy                                                                                                                                                                                                         |
| 10b  | `GOWORK=off go mod tidy` in all modules **including visualtest**; stage the diff                                       | v1.11.0: skipping the tidy left go.sums at the previous version — 9 days red                                                                                                                                                       |
| 10c  | Second tidy must be a no-op (`git diff --exit-code`)                                                                   | same invariant CI enforces in "Verify no untracked changes"                                                                                                                                                                        |

Why pre-push tidying is only PARTIALLY safe: `visualtest/go.mod` carries local
replace directives for its siblings, so ITS tidy needs no proxy. The five
published sub-modules (icons, errorpage, charts/echarts, datastar, htmx) have no
replaces — their 10b tidy resolves siblings against the module proxy, where the
new tags do not exist yet. Pre-push, their `go.sum` files therefore CANNOT be
moved to the new version: the sweep no-ops and the stale checksums ship. That is
why the 24-hour watch repeats the tidy after propagation (v1.11.0 red for 9 days;
v1.12.0 repeated it). See `docs/modularization/README.md` ("visualtest
sibling-pin policy").

### If the cut aborts after step 8 (the release commit exists)

Never re-run the script from the top and never `git reset`. The release commit is
created only after the full verify suite passes, so it is almost always correct:
verify its tree by hand (`git show <commit>:utils/version.go`, the CHANGELOG
heading, `grep -c '^replace' <sub>/go.mod` per tagged module), create the tags on
it (`git tag -a ...` per module), then perform steps 10/10b manually. v1.12.0's
8b assertions aborted exactly this way (SIGPIPE — see incident index) and the
commit was verified byte-correct.

## After the script, before pushing

- [ ] `git show v<version>` — tag message, then `git show <release-commit>`:
      version files moved together, replaces absent from every tagged go.mod,
      CHANGELOG heading + fresh empty `[Unreleased]`, release notes in the body.
- [ ] `git log --oneline -3` — no daemon commit wedged between the release
      commit and the tags.
- [ ] **CSS byte-stability:** `nix run .#css` produces no diff. v1.10.0 shipped
      demo CSS one release stale because BuildFlow recompiled only after the
      release commit.
- [ ] **website typescript pin intact** (`website/package.json`): `astro check`
      crashes on TS 7 — needs the 6.x pin until `astro check` supports the
      native compiler. The daemon has flipped it back before.
- [ ] `scripts/check-release-tags.sh <version>` — every published module has its
      tag on the release commit.
- [ ] All four fast guards pass: `check-lint-config.sh`, `check-templ-sync.sh`,
      `check-version-sync.sh`, `check-module-sync.sh`.

## Push (manual — house rule: never push without the human)

```bash
git push origin master --follow-tags
```

The proxy picks up the tags asynchronously; a `go get` immediately after push can
still 404 until ingestion completes.

## After pushing — the 24-hour watch

- [ ] **CI + Website green on master AND on the tag.**
- [ ] **Post-propagation tidy sweep:** once `go list -m
      github.com/larsartmann/templ-components@<version>` resolves, re-run the
      `GOWORK=off go mod tidy` sweep across all 8 modules and commit the go.sum
      refresh. The pre-push 10b sweep cannot update the five replace-less
      sub-modules (see the note above the step table); CI's tidy can, and fails
      "Verify no untracked changes" + Lint until the refresh lands. v1.11.0 sat
      red 9 days for exactly this; v1.12.0 repeated it.
- [ ] **GitHub Release page** created from the release notes (`gh release create
      v<version> --notes-file ... --latest`); v1.11.0 and v1.12.0 shipped
      tag-only for days before the pages were backfilled.
- [ ] **Dependabot alerts resolved** — a security-update workflow can run
      "successfully" while changing nothing; verify open alerts actually close
      after the fix lands on master (`gh api repos/.../dependabot/alerts`).
- [ ] `go list -m github.com/larsartmann/templ-components@<version>` (and one
      sub-module) resolves from the proxy.
- [ ] **pkg.go.dev** shows the new version for root + sub-modules.
- [ ] If anything is red: fix on master, and only cut `+0.0.1` — **never retag a
      version** (the proxy caches tags permanently; `retract` is the escape hatch
      of last resort).

## Incident index (why this page is long)

| Version | Incident                                                                        | Guard added                                                                          |
| ------- | ------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------ |
| v1.8.3  | Root-only release; sub-module consumers broke                                   | tag set derived from go.mod replaces; `check-release-tags.sh` + CI gate              |
| v1.9.0  | Verify after replace-strip → `unknown revision` everywhere                      | verify-before-strip ordering; single shared EXIT trap                                |
| v1.10.0 | Daemon race: split commits, pushed tags mid-cut, stale CSS, leaked self-replace | 8b tree assertions; CSS recompile inside step 7; `[ /]` strip pattern                |
| v1.11.0 | Re-add commit without re-tidy → go.sum drift → 9-day red master                 | 10b tidy (visualtest included) + 10c idempotence gate; visualtest sibling-pin policy |
| v1.12.0 | 8b assertion SIGPIPE (`git show \| grep -q` exit 141 past the 64KB pipe buffer, CHANGELOG at 155KB) aborted the cut post-commit; recovery skipped no module, but the pre-push tidy still left sub-module go.sums stale → CI red ×3 (Build & Test drift, Lint missing go.sum, Visual Regression mid-animation capture) | 8b assertions use `grep -c`; late-abort recovery procedure documented above this table; post-propagation tidy added to the 24-hour watch; visualtest harness settles all finite animations before capture |
