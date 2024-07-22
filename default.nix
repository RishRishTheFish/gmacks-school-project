{ pkgs ? import <nixpkgs> {} }:

pkgs.mkShell {
  buildInputs = [
    pkgs.go
    pkgs.libGL
    pkgs.libGLU
    pkgs.xorg.libXxf86vm
    pkgs.gcc
    pkgs.pkg-config
    pkgs.xorg.libX11
    pkgs.xorg.libXi
    pkgs.xorg.libXcursor
    pkgs.xorg.libXrandr
    pkgs.xorg.libXinerama
    pkgs.wayland
    pkgs.wayland-protocols
    pkgs.mesa
    pkgs.gtk3
  ];
  shellHook = ''
    export CGO_ENABLED=1
  '';
}
