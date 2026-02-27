{
  description = "Steampipe plugin for Redmine";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      supportedSystems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
    in
    {
      packages = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          default = self.packages.${system}.steampipe-plugin-redmine;

          steampipe-plugin-redmine = pkgs.buildGoModule.override { go = pkgs.go_1_26; } {
            pname = "steampipe-plugin-redmine";
            version = "0.1.0";

            src = pkgs.lib.cleanSource ./.;

            vendorHash = "sha256-d/ZCR1hl+MUFCSMyESZbFmeGtQIb+Vp1hnqEc3cTm6g=";

            ldflags = [
              "-s"
              "-w"
            ];

            doCheck = true;

            installPhase = ''
              runHook preInstall

              mkdir -p $out
              cp $GOPATH/bin/steampipe-plugin-redmine $out/steampipe-plugin-redmine.plugin
              cp -R config $out/.

              runHook postInstall
            '';

            meta = {
              description = "Redmine Plugin for Steampipe";
              license = pkgs.lib.licenses.asl20;
            };
          };
        }
      );

      devShells = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          default = pkgs.mkShell {
            buildInputs = [
              pkgs.go_1_26
              pkgs.gopls
              pkgs.gotools
              pkgs.go-tools
              pkgs.delve
              pkgs.just
              pkgs.steampipe
              pkgs.patchelf
            ];

            shellHook = ''
              echo "steampipe-plugin-redmine dev shell"
              echo "Go $(go version | cut -d' ' -f3)"
              echo "Run 'just' to see available commands"

              # NixOS: patch steampipe's bundled postgres binaries if needed
              _sp_db_dir="''${STEAMPIPE_INSTALL_DIR:-$HOME/.steampipe}/db"
              if [ -d "$_sp_db_dir" ]; then
                _interp="$(cat ${pkgs.stdenv.cc}/nix-support/dynamic-linker)"
                for bin in "$_sp_db_dir"/*/postgres/bin/*; do
                  if file "$bin" 2>/dev/null | grep -q "ELF.*dynamically linked" && \
                     readelf -l "$bin" 2>/dev/null | grep -q "interpreter: /lib64"; then
                    patchelf --set-interpreter "$_interp" "$bin" 2>/dev/null || true
                  fi
                done
              fi
              unset _sp_db_dir _interp
            '';
          };
        }
      );
    };
}
