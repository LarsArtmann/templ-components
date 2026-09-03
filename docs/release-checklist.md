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

Why pre-push tidying is safe: `visualtest/go.mod` carries local replace
directives for its siblings (bumped together with its pins in 5b), so nothing in
10b/10c needs the new tags to exist on the module proxy yet. See
`docs/modularization/README.md` ("visualtest sibling-pin policy").

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
- [ ] **`go mod tidy` round-trip:** after the re-add commit lands, CI's
      per-module tidy must leave the tree clean. v1.11.0's 9-day red was exactly
      this check failing silently.
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
