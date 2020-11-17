#!/bin/bash

set -euo pipefail

GOARCH=arm go build -v
scp .env spotify-screen pi@raspberrypi:
ssh -t pi@raspberrypi ./spotify-screen
