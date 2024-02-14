{
  description = "Description for the project";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-23.11";
    denv = {
      url = "github:iliayar/env.nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = inputs@{ flake-parts, denv, ... }:
    flake-parts.lib.mkFlake { inputs = denv.inputs; } {
      imports = [ denv.flakeModules.default ];
      systems =
        [ "x86_64-linux" "aarch64-linux" "aarch64-darwin" "x86_64-darwin" ];
      perSystem = { config, self', inputs', pkgs, system, ... }: {
        denvs.default = { 
            langs.go.enable = true;

            langs.ocaml.enable = true;
            denv.packages = with pkgs; [
              ocamlPackages.batteries
              ocamlPackages.core
              ocamlPackages.ppx_deriving
              ocamlPackages.ppx_sexp_conv
            
              nodePackages.esy
            ];
        };
      };
      flake = { };
    };
}
