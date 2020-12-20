#!/bin/bash -x

set -euo pipefail
shopt -s expand_aliases
alias 'rs'='rsync --rsync-path="sudo /usr/bin/rsync" --info=progress2 --checksum'

docker build . -o build
rs ./build/spotify-screen pi@raspberrypi:/usr/local/bin/
rs ./spotify-screen.service pi@raspberrypi:/etc/systemd/system/
ssh pi@raspberrypi -- 'bash -xc "sudo systemctl daemon-reload && sudo systemctl restart spotify-screen && systemctl status spotify-screen"'
