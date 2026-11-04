#!/bin/sh
# 把 third_party/sing-box-sbx 叠到指定的 reF1nd 源码树。
# 用法: scripts/apply-sbx-kernel.sh /path/to/reF1nd-sing-box
set -eu
ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
OVERLAY="$ROOT/third_party/sing-box-sbx"
TARGET=${1:-}

if [ -z "$TARGET" ] || [ ! -d "$TARGET" ]; then
  echo "usage: $0 /path/to/reF1nd-sing-box" >&2
  exit 1
fi
if [ ! -f "$TARGET/go.mod" ]; then
  echo "not a sing-box tree: $TARGET" >&2
  exit 1
fi

find "$OVERLAY" -type f ! -name PIN | while IFS= read -r file; do
  rel=${file#"$OVERLAY"/}
  dest="$TARGET/$rel"
  mkdir -p "$(dirname -- "$dest")"
  cp -a "$file" "$dest"
done
# Reality mldsa65Verify 使用 Go 1.27 的 crypto/mldsa。
if grep -q '^go 1\.' "$TARGET/go.mod"; then
  sed -i 's/^go 1\.[0-9.]*/go 1.27.0/' "$TARGET/go.mod"
fi
echo "applied overlay onto $TARGET"
