{ pkgs ? import <nixpkgs> {} }:
  pkgs.mkShell {
    # nativeBuildInputs is usually what you want -- tools you need to run
    nativeBuildInputs = with pkgs.buildPackages; [
      go
      libGL
      pkg-config
      xorg.libX11.dev
      xorg.libX11
      xorg.libXcursor
      xorg.libXi
      xorg.libXinerama
      xorg.libXrandr 
      xorg.libXxf86vm
      gcc
    ];
  hardeningDisable = [ "fortify" ];
}
