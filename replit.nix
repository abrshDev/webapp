{ pkgs }: {
  deps = [
    pkgs.go
    pkgs.chromium
    pkgs.git
    pkgs.wget
    pkgs.curl
    pkgs.nodejs
    pkgs.python3
    pkgs.unzip
    pkgs.gcc
    pkgs.make
    pkgs.ngrok
  ];
}
