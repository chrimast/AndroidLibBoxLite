#!/bin/sh
set -eu
export PATH="/opt/data/home/go-sdk/go/bin:$PATH"
export GOTOOLCHAIN=local
export GOPATH=/opt/data/home/gopath
export GOCACHE=/tmp/gocache-netbox
export GOMODCACHE=/opt/data/home/gopath/pkg/mod
export JAVA_HOME=/opt/data/home/jdk/jdk-25.0.5+8
export ANDROID_HOME=/opt/data/home/android-sdk
export ANDROID_NDK_HOME=/opt/data/home/android-sdk/ndk/28.2.13676358
mkdir -p /tmp/gocache-netbox
STAGE=/tmp/boxlite-work/androidlibbox-staging/sing-box
BOUND="$STAGE/build/arm64/shared"
rm -rf "$BOUND"
mkdir -p "$BOUND"
cp -a /opt/data/workspace/AndroidLibBoxLite/.agents/gobind/. "$BOUND/"
cp /opt/data/workspace/AndroidLibBoxLite/native/cli/export.go.txt "$BOUND/cli_export.go"
cp "$STAGE/go.mod" "$BOUND/go.mod"
cp "$STAGE/go.sum" "$BOUND/go.sum"
python3 - "$BOUND/go.mod" "$STAGE" <<'PY'
from pathlib import Path
import sys
path, stage = Path(sys.argv[1]), sys.argv[2]
text = path.read_text()
text = text.replace('module github.com/sagernet/sing-box', 'module gobind', 1)
if 'replace github.com/sagernet/sing-box =>' not in text:
    text += f'\nreplace github.com/sagernet/sing-box => {stage}\n'
path.write_text(text)
PY
TAGS='with_quic,with_dhcp,with_wireguard,with_utls,with_tailscale,with_openvpn,with_openconnect,with_naive_outbound,with_ebpf,with_xhttp,badlinkname,tfogo_checklinkname0,ts_omit_logtail,ts_omit_ssh,ts_omit_drive,ts_omit_taildrop,ts_omit_webclient,ts_omit_doctor,ts_omit_capture,ts_omit_kube,ts_omit_aws,ts_omit_synology,ts_omit_bird'
LDFLAGS='-s -w -buildid= -X github.com/sagernet/sing-box/constant.Version=1.15.0-alpha.13-reF1nd -X runtime.godebugDefault=multipathtcp=0,tlssha1=1 -checklinkname=0'
CLANG="$ANDROID_NDK_HOME/toolchains/llvm/prebuilt/linux-x86_64/bin/clang"
OUT=/opt/data/workspace/AndroidLibBoxLite/.agents/dist/libbox.so
export GOOS=android GOARCH=arm64 GOARM=7 CGO_ENABLED=1
export CC="$CLANG --target=aarch64-linux-android26"
export CGO_LDFLAGS='-Wl,-z,max-page-size=16384'
cd "$BOUND"
go mod download github.com/sagernet/gomobile
go build -mod=mod -trimpath -buildvcs=false -tags "$TAGS" -ldflags "$LDFLAGS" -buildmode=c-shared -o "$OUT" .
ls -lh "$OUT"
