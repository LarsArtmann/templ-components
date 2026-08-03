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
    nixpkgs-chromium.url = "github:NixOS/nixpkgs/bfb0bf3c2c9aa3c8d8dc4b97a6f0e5e6f0121c10";
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
        { config, pkgs, inputs, ... }:
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
              # GOWORK=off: the root go.work references an absolute path to
              # go-error-family (Lars' machine only). Setting GOWORK=off in the
              # devShell ensures `go build ./...`, `go generate ./...`, and
              # BuildFlow all operate on the main module standalone. The
              # visualtest module builds via its own go.mod with a local replace.
              export GOEXPERIMENT=jsonv2
              export GOWORK=off
            '';
          };

          apps = {
            test = {
              type = "app";
              meta.description = "Run all tests with race detector";
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
              meta.description = "Run golangci-lint across all packages";
              program = pkgs.writeShellApplication {
                name = "run-lint";
                runtimeInputs = [ pkgs.golangci-lint ];
                text = ''
                  golangci-lint run ./datastar/... ./display/... ./errorpage/... ./feedback/... ./forms/... ./htmx/... ./icons/... ./integration/... ./layout/... ./navigation/... ./recipes/... ./utils/... ./internal/...
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
              meta.description = "Full verification: generate + build + test + lint";
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
                  echo "==> Building..."
                  go build ./...
                  echo "==> Testing..."
                  go test ./... -count=1
                  echo "==> Testing visualtest module (separate go.mod)..."
                  # visualtest is a separate Go module; ./... from repo root skips it.
                  # Tests skip cleanly when Chromium is absent, so this never red-lines CI.
                  cd visualtest && GOWORK=off GOEXPERIMENT=jsonv2 go test -count=1 ./... && cd ..
                  echo "==> Linting..."
                  golangci-lint run ./datastar/... ./display/... ./errorpage/... ./feedback/... ./forms/... ./htmx/... ./icons/... ./integration/... ./layout/... ./navigation/... ./recipes/... ./utils/... ./internal/...
                  echo "==> All checks passed."
                '';
              };
            };

            coverage = {
              type = "app";
              meta.description = "Run tests with coverage report";
              program = pkgs.writeShellApplication {
                name = "run-coverage";
                runtimeInputs = [ pkgs.go_1_26 ];
                text = ''
                  export GOEXPERIMENT=jsonv2
                  go test ./... -count=1 -coverprofile=coverage.out
                  go tool cover -func=coverage.out | tail -1
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
                  pkgs.chromium
                ];
                text = ''
                  export GOEXPERIMENT=jsonv2
                  # visualtest is its own module with a local replace directive;
                  # the parent go.work would shadow it, so disable workspace mode.
                  export GOWORK=off
                  export CHROMEDP_CHROME_PATH="${pkgs.chromium}/bin/chromium"
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
