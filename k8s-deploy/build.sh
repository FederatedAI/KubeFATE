#!/usr/bin/env bash

ldflags+=$(sh ./generate-ldflags.sh)

echo ldflags:${ldflags}
if [ "$1" = "win" ]; then
  go build -o kubefate.exe -i -v -gcflags='-N -l' -ldflags="${ldflags}"
else
  go build -o kubefate -i -v -gcflags='-N -l' -ldflags="${ldflags}"
fi