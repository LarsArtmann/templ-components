{
  description = "Reusable UI components for Go web apps — templ + Tailwind CSS v4 + HTMX";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
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
        { config, pkgs, ... }:
        {
          devShells.default = pkgs.mkShellNoCC {
            packages = with pkgs; [
              go_1_26
              gopls
              golangci-lint
              templ
              tailwindcss_4
            ];
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
                  golangci-lint run ./display/... ./errorpage/... ./feedback/... ./forms/... ./htmx/... ./icons/... ./layout/... ./navigation/... ./utils/... ./internal/...
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
                  echo "==> Linting..."
                  golangci-lint run ./display/... ./errorpage/... ./feedback/... ./forms/... ./htmx/... ./icons/... ./layout/... ./navigation/... ./utils/... ./internal/...
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
          };

          # treefmt: format .nix (nixfmt) and .go (gofmt + goimports).
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
              gofmt.enable = true;
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
