{
  description = "Reusable UI components for Go web apps — templ + Tailwind CSS v4 + HTMX";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    # Separate rolling nixpkgs pin for the GO TOOLCHAIN. The security-fix
    # patches (GO-2026-5972 encoding/asn1, GO-2026-6089 net/http,
    # GO-2026-6090 crypto/tls — all reachable via the demo's http.Server)
    # landed in go 1.26.6, but the locked `nixpkgs` input still packages
    # 1.26.5. Updating `nixpkgs` wholesale would ALSO drift pkgs.templ
    # (breaking the v0.3.1020 generate pin — see AGENTS.md "templ Version
    # Pin") and pkgs.golangci-lint (breaking lint reproducibility), so the
    # toolchain rides its own input — same isolation pattern as
    # nixpkgs-chromium. Fold this input back into `nixpkgs` at the next
    # deliberate full-flake update.
    nixpkgs-go.url = "github:NixOS/nixpkgs/nixos-unstable";
    # Separate nixpkgs pin for Chromium. Visual regression tests are
    # pixel-sensitive: a Chromium major bump can shift font AA, sub-pixel
    # layout, or rendering timings enough to flip goldens. Pinning Chromium
    # to its own nixpkgs revision (updated independently and deliberately)
    # insulates the visual goldens from routine `nix flake update` bumps.
    # To update: change the rev below to a nixpkgs commit with the desired
    # Chromium, run `nix flake lock`, then `nix run .#visual -- -update`.
    nixpkgs-chromium.url = "github:NixOS/nixpkgs/148bab9c1c3c53136ecb44a6ea356a0ed5b39b06";
    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };
    systems.url = "github:nix-systems/default";
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs@{
      self,
      flake-parts,
      ...
    }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import inputs.systems;

      # treefmt-nix provides the `treefmt` config module and a `formatter` app
      # automatically (replacing the former bare `formatter = pkgs.nixfmt;`).
      imports = [ inputs.treefmt-nix.flakeModule ];

      perSystem =
        {
          config,
          pkgs,
          inputs',
          ...
        }:
        let
          # Go toolchain from the dedicated nixpkgs-go input (1.26.6+ for the
          # security fixes above); everything else stays on the locked
          # nixpkgs so the templ pin holds. golangci-lint also rides this
          # input: the locked nixpkgs ships 2.12.2, which predates the
          # exhaustruct_v5 migration used by .golangci.yml (CI pins 2.13.2;
          # 2.13.1 accepts the same config).
          goToolchain = inputs'.nixpkgs-go.legacyPackages.go_1_26;
          golangciLint = inputs'.nixpkgs-go.legacyPackages.golangci-lint;
        in
        {
          devShells.default = pkgs.mkShellNoCC {
            packages = with pkgs; [
              goToolchain
              gopls
              golangciLint
              templ
              tailwindcss_4
            ];
            shellHook = ''
              # GOEXPERIMENT=jsonv2: required until Go 1.27 stabilizes it.
              # GOWORK is now ON: go.work lists all 7 modules (root, utils,
              # icons, errorpage, charts/echarts, datastar, htmx) + visualtest.
              # Use GOWORK=off only for per-module isolation testing.
              export GOEXPERIMENT=jsonv2
            '';
          };

          apps = {
            test = {
              type = "app";
              meta.description = "Run all tests with race detector (all modules via go.work)";
              program = pkgs.writeShellApplication {
                name = "run-tests";
                runtimeInputs = [ goToolchain ];
                text = ''
                  export GOEXPERIMENT=jsonv2
                  go test ./... -count=1 -race
                '';
              };
            };

            lint = {
              type = "app";
              meta.description = "Run golangci-lint across all modules";
              program = pkgs.writeShellApplication {
                name = "run-lint";
                runtimeInputs = [
                  golangciLint
                  # Same version CI pins (.github/workflows/ci.yaml Actionlint
                  # step) so local `nix run .#lint` and CI agree on workflow
                  # findings.
                  pkgs.actionlint
                ];
                text = ''
                  # golangci-lint does not support go.work workspace mode.
                  # Lint each module independently with GOWORK=off.
                  export GOEXPERIMENT=jsonv2
                  echo "==> Linting root module..."
                  GOWORK=off golangci-lint run ./display/... ./feedback/... ./forms/... ./integration/... ./layout/... ./navigation/... ./recipes/... ./internal/... ./cmd/...
                  echo "==> Linting utils module..."
                  (cd utils && GOWORK=off golangci-lint run ./...)
                  echo "==> Linting icons module..."
                  (cd icons && GOWORK=off golangci-lint run ./...)
                  echo "==> Linting errorpage module..."
                  (cd errorpage && GOWORK=off golangci-lint run ./...)
                  echo "==> Linting charts/echarts module..."
                  (cd charts/echarts && GOWORK=off golangci-lint run ./...)
                  echo "==> Linting datastar module..."
                  (cd datastar && GOWORK=off golangci-lint run ./...)
                  echo "==> Linting htmx module..."
                  (cd htmx && GOWORK=off golangci-lint run ./...)
                  echo "==> Actionlint (GitHub Actions workflows)..."
                  actionlint
                '';
              };
            };

            build = {
              type = "app";
              meta.description = "Regenerate templ + build all packages";
              program = pkgs.writeShellApplication {
                name = "run-build";
                runtimeInputs = [
                  goToolchain
                  pkgs.templ
                ];
                text = ''
                  export GOEXPERIMENT=jsonv2
                  find . -name '*_templ.go' -print0 | xargs -0 rm
                  templ generate ./...
                  go build ./...
                  echo "Build successful."
                '';
              };
            };

            verify = {
              type = "app";
              meta.description = "Full verification: generate + build + test + lint (all modules)";
              program = pkgs.writeShellApplication {
                name = "run-verify";
                runtimeInputs = [
                  goToolchain
                  golangciLint
                  pkgs.templ
                ];
                text = ''
                  export GOEXPERIMENT=jsonv2
                  echo "==> Regenerating templ..."
                  find . -name '*_templ.go' -print0 | xargs -0 rm
                  templ generate ./...
                  echo "==> Building (workspace)..."
                  go build ./...
                  echo "==> Testing (workspace)..."
                  go test ./... -count=1
                  echo "==> Per-module GOWORK=off isolation tests..."
                  for mod in utils icons errorpage charts/echarts datastar htmx; do
                    echo "  -> $mod"
                    (cd "$mod" && GOWORK=off go test -count=1 ./...)
                  done
                  echo "==> Testing visualtest module..."
                  (cd visualtest && GOWORK=off GOEXPERIMENT=jsonv2 go test -count=1 ./...)
                  echo "==> Linting per-module..."
                  echo "  -> root"
                  GOWORK=off golangci-lint run ./display/... ./feedback/... ./forms/... ./integration/... ./layout/... ./navigation/... ./recipes/... ./internal/... ./cmd/...
                  for mod in utils icons errorpage charts/echarts datastar htmx; do
                    echo "  -> $mod"
                    (cd "$mod" && GOWORK=off golangci-lint run ./...)
                  done
                  echo "==> All checks passed."
                '';
              };
            };

            coverage = {
              type = "app";
              meta.description = "Run tests with coverage report (all modules)";
              program = pkgs.writeShellApplication {
                name = "run-coverage";
                runtimeInputs = [ goToolchain ];
                text = ''
                  export GOEXPERIMENT=jsonv2
                  echo "=== Root module ==="
                  go test ./... -count=1 -coverprofile=coverage.out
                  go tool cover -func=coverage.out | tail -1
                  for mod in utils icons errorpage charts/echarts datastar htmx; do
                    echo "=== $mod ==="
                    (cd "$mod" && go test ./... -count=1 -coverprofile=coverage.out && go tool cover -func=coverage.out | tail -1)
                  done
                '';
              };
            };

            css = {
              type = "app";
              meta.description = "Recompile the demo CSS (examples/demo/static/app.css) via tailwindcss --minify";
              program = pkgs.writeShellApplication {
                name = "run-css";
                runtimeInputs = [ pkgs.tailwindcss_4 ];
                text = ''
                  tailwindcss \
                    --input examples/demo/demo.css \
                    --output examples/demo/static/app.css \
                    --minify
                  echo "CSS compiled: examples/demo/static/app.css"
                '';
              };
            };

            visual = {
              type = "app";
              meta.description = "Run pixel-level visual regression tests (headless Chromium via chromedp)";
              program = pkgs.writeShellApplication {
                name = "run-visual";
                runtimeInputs = [
                  goToolchain
                  # fc-match for the font diagnostics below.
                  pkgs.fontconfig
                  # Use Chromium from the pinned nixpkgs-chromium input (not the
                  # main nixpkgs) so visual goldens don't shift on routine
                  # `nix flake update`. Update deliberately: see nixpkgs-chromium
                  # input comment in flake.nix.
                  inputs'.nixpkgs-chromium.legacyPackages.chromium
                ];
                text = ''
                  export GOEXPERIMENT=jsonv2
                  # visualtest is its own module with a local replace directive;
                  # the parent go.work would shadow it, so disable workspace mode.
                  export GOWORK=off
                  export CHROMEDP_CHROME_PATH="${inputs'.nixpkgs-chromium.legacyPackages.chromium}/bin/chromium"
                  # Font determinism: the demo CSS declares Inter, JetBrains
                  # Mono, and Space Grotesk (headings). Neither dev machines
                  # nor CI runners reliably have them installed, and host
                  # fontconfig resolves the fallbacks differently in every
                  # environment, which shifted text rendering and flipped
                  # goldens between local and CI.
                  #
                  # Fix: a FULLY PURE fonts.conf. Note: pkgs.makeFontsConf is
                  # unsuitable here — it includes /etc/fonts/conf.d and impure
                  # FHS/profile font dirs by default, reintroducing exactly
                  # the host drift this pin exists to eliminate. This conf
                  # lists only nix store font directories, so rendering is
                  # identical on every machine. Space Grotesk is not in nixpkgs
                  # (headings fall back to Inter, next in the CSS stack) and
                  # the jetbrains-mono derivation is currently broken upstream
                  # (gftools dependency fails; mono text falls back to DejaVu
                  # Sans Mono). DejaVu covers glyphs Inter lacks.
                  #
                  # After changing the font list, regenerate ALL goldens:
                  #   nix run .#visual -- -update
                  export FONTCONFIG_FILE="${pkgs.writeText "tc-visual-fonts.conf" ''
                    <?xml version="1.0"?>
                    <!DOCTYPE fontconfig SYSTEM "urn:fontconfig:fonts.dtd">
                    <fontconfig>
                      <dir>${pkgs.inter}</dir>
                      <dir>${pkgs.dejavu_fonts}</dir>
                      <cachedir>/tmp/tc-visualtest-fontconfig-cache</cachedir>
                    </fontconfig>
                  ''}"
                  cd visualtest
                  # Font guard: under the pinned FONTCONFIG_FILE every CSS
                  # generic must resolve to Inter (Space Grotesk and JetBrains
                  # Mono are absent, so headings and code fall back to it).
                  # If a generic resolves to anything else, the pure pin is
                  # broken (host fonts leaked in) and goldens WILL shift, so
                  # fail fast with a clear message instead of pixel diffs.
                  for family in sans-serif serif monospace; do
                    result=$(fc-match "$family")
                    echo "font pin: $family -> $result"
                    case "$result" in
                      Inter*) ;;
                      *)
                        echo "ERROR: $family resolves away from Inter ($result)." >&2
                        echo "The pure fontconfig pin is broken; fix fonts.conf before trusting goldens." >&2
                        exit 1
                        ;;
                    esac
                  done
                  # Forward extra args (e.g. -update, -run TestButtons) to go test.
                  go test ./... -count=1 "$@"
                '';
              };
            };
          };

          # treefmt: format .nix (nixfmt) and .go (gofumpt + goimports).
          # gofumpt (not gofmt) aligns with .golangci.yml's gofumpt linter,
          # preventing a latent conflict where treefmt and golangci-lint
          # disagree on formatting.
          # Generated *_templ.go files are excluded — they are templ output and
          # must not be hand-reformatted (would cause perpetual churn vs the
          # generator). Format enforcement runs via `nix flake check` (see checks
          # below) and `nix fmt`.
          treefmt = {
            settings.excludes = [
              "**/*_templ.go"
              "website/**"
              "examples/demo/static/**"
            ];
            programs = {
              nixfmt.enable = true;
              gofumpt.enable = true;
              goimports.enable = true;
              # gotools (goimports) shells out to `go` at format time: that
              # go must match the go.mod toolchain (1.26.7) or it tries to
              # DOWNLOAD the newer toolchain, which fails inside the
              # flake-check sandbox (no network). Same input as
              # goToolchain.
              goimports.package = inputs'.nixpkgs-go.legacyPackages.gotools;
            };
          };

          # `nix flake check` runs these. format = treefmt verification (catches
          # unformatted nix/go files before they land). build = hermetic
          # templ-generate + go build (catches compile errors without needing a
          # developer shell).
          checks = {
            format = config.treefmt.build.check self;
          };
        };
    };
}
