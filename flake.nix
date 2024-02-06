{
  outputs = { self, nixpkgs }: {
    devShells.x86_64-linux.default = let
      pkgs = nixpkgs.legacyPackages.x86_64-linux;
    in pkgs.mkShellNoCC {
      buildInputs = with pkgs; [
        (python3.withPackages (p: with p; [
          requests
        ]))
        go
        gnumake
        wget
        nodePackages.prettier
        typescript
      ];
    };
  };
}
