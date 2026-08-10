{
  description = "Reusable UI components for Go web apps — templ + Tailwind CSS v4 + HTMX";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
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
        {
          devShells.default = pkgs.mkShellNoCC {
            packages = with pkgs; [
              go_1_26
              gopls
              golangci-lint
              templ
              tailwindcss_4
            ];
            shellHook = ''
              # GOEXPERIMENT=jsonv2: required until Go 1.27 stabilizes it.
              # GOWORK is now ON: go.work lists all 5 modules (root, utils,
              # icons, errorpage, charts/echarts) + visualtest. Use
              # GOWORK=off only for per-module isolation testing.
              export GOEXPERIMENT=jsonv2
            '';
          };

          apps = {
            test = {
              type = "app";
              meta.description = "Run all tests with race detector (all modules via go.work)";
              program = pkgs.writeShellApplication {
                name = "run-tests";
                runtimeInputs = [ pkgs.go_1_26 ];
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
                runtimeInputs = [ pkgs.golangci-lint ];
                text = ''
                  # golangci-lint does not support go.work workspace mode.
                  # Lint each module independently with GOWORK=off.
                  export GOEXPERIMENT=jsonv2
                  echo "==> Linting root module..."
                  GOWORK=off golangci-lint run ./display/... ./feedback/... ./forms/... ./htmx/... ./datastar/... ./integration/... ./layout/... ./navigation/... ./recipes/... ./internal/... ./cmd/...
                  echo "==> Linting utils module..."
                  (cd utils && GOWORK=off golangci-lint run ./...)
                  echo "==> Linting icons module..."
                  (cd icons && GOWORK=off golangci-lint run ./...)
                  echo "==> Linting errorpage module..."
                  (cd errorpage && GOWORK=off golangci-lint run ./...)
                  echo "==> Linting charts/echarts module..."
                  (cd charts/echarts && GOWORK=off golangci-lint run ./...)
                '';
              };
            };

            build = {
              type = "app";
              meta.description = "Regenerate templ + build all packages";
              program = pkgs.writeShellApplication {
                name = "run-build";
                runtimeInputs = [
                  pkgs.go_1_26
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
                  pkgs.go_1_26
                  pkgs.golangci-lint
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
                  for mod in utils icons errorpage charts/echarts; do
                    echo "  -> $mod"
                    (cd "$mod" && GOWORK=off go test -count=1 ./...)
                  done
                  echo "==> Testing visualtest module..."
                  (cd visualtest && GOWORK=off GOEXPERIMENT=jsonv2 go test -count=1 ./...)
                  echo "==> Linting per-module..."
                  echo "  -> root"
                  GOWORK=off golangci-lint run ./display/... ./feedback/... ./forms/... ./htmx/... ./datastar/... ./integration/... ./layout/... ./navigation/... ./recipes/... ./internal/... ./cmd/...
                  for mod in utils icons errorpage charts/echarts; do
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
                runtimeInputs = [ pkgs.go_1_26 ];
                text = ''
                  export GOEXPERIMENT=jsonv2
                  echo "=== Root module ==="
                  go test ./... -count=1 -coverprofile=coverage.out
                  go tool cover -func=coverage.out | tail -1
                  for mod in utils icons errorpage charts/echarts; do
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
                  pkgs.go_1_26
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
                  cd visualtest
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
