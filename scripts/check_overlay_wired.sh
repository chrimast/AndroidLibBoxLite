#!/bin/sh
# Fail if NetBOX overlay is not wired into the AAR build.
set -eu
ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$ROOT"
test -f third_party/sing-box-sbx/transport/v2rayxhttp/register.go
test -f third_party/sing-box-sbx/common/tls/reality_client.go
test -f scripts/apply-sbx-kernel.sh
grep -q with_xhttp build_tags.txt
grep -q apply-sbx-kernel.sh scripts/prepare_sources.py
grep -q apply_sbx_overlay scripts/prepare_sources.py
grep -q MlkemEcdhe scripts/prepare_sources.py
echo "overlay wiring ok"
