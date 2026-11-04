//go:build with_utls

package tls

import (
	"context"
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"net"
	"testing"

	"github.com/sagernet/sing-box/option"
	"github.com/sagernet/sing/common/logger"

	utls "github.com/metacubex/utls"
)

func TestRealityChromeClientHelloKeepsX25519MLKEM768(t *testing.T) {
	privateKey, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	config, err := NewRealityClient(context.Background(), logger.NOP(), "www.example.com", option.OutboundTLSOptions{
		Enabled:    true,
		ServerName: "www.example.com",
		UTLS:       &option.OutboundUTLSOptions{Enabled: true, Fingerprint: "chrome"},
		Reality: &option.OutboundRealityOptions{
			Enabled:   true,
			PublicKey: base64.RawURLEncoding.EncodeToString(privateKey.PublicKey().Bytes()),
			ShortID:   "0123abcd",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	client, ok := config.(*RealityClientConfig)
	if !ok {
		t.Fatalf("expected *RealityClientConfig, got %T", config)
	}
	conn, peer := net.Pipe()
	defer conn.Close()
	defer peer.Close()
	uConn, _, err := client.prepareClientHello(conn)
	if err != nil {
		t.Fatal(err)
	}
	hybridIndex := -1
	classicalIndex := -1
	for i, share := range uConn.HandshakeState.Hello.KeyShares {
		switch share.Group {
		case utls.X25519MLKEM768:
			if hybridIndex == -1 {
				hybridIndex = i
			}
		case utls.X25519:
			if classicalIndex == -1 {
				classicalIndex = i
			}
		}
	}
	if hybridIndex < 0 {
		t.Fatal("chrome Reality ClientHello missing X25519MLKEM768 key share")
	}
	if classicalIndex < 0 {
		t.Fatal("chrome Reality ClientHello missing X25519 key share")
	}
	if hybridIndex > classicalIndex {
		t.Fatalf("X25519MLKEM768 must precede X25519, got %d then %d", hybridIndex, classicalIndex)
	}
	hasHybridGroup := false
	for _, curve := range uConn.HandshakeState.Hello.SupportedCurves {
		if curve == utls.X25519MLKEM768 {
			hasHybridGroup = true
			break
		}
	}
	if !hasHybridGroup {
		t.Fatal("chrome Reality ClientHello missing X25519MLKEM768 in supported_groups")
	}
}

func TestRealityClientHelloReportsXrayCompatibleVersion(t *testing.T) {
	privateKey, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	config, err := NewRealityClient(context.Background(), logger.NOP(), "www.example.com", option.OutboundTLSOptions{
		Enabled:    true,
		ServerName: "www.example.com",
		UTLS:       &option.OutboundUTLSOptions{Enabled: true, Fingerprint: "chrome"},
		Reality: &option.OutboundRealityOptions{
			Enabled:   true,
			PublicKey: base64.RawURLEncoding.EncodeToString(privateKey.PublicKey().Bytes()),
			ShortID:   "0123abcd",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	client, ok := config.(*RealityClientConfig)
	if !ok {
		t.Fatalf("expected *RealityClientConfig, got %T", config)
	}
	conn, peer := net.Pipe()
	defer conn.Close()
	defer peer.Close()
	uConn, _, err := client.prepareClientHello(conn)
	if err != nil {
		t.Fatal(err)
	}
	// SessionId is AEAD-sealed before return; assert the plaintext version helper.
	_ = uConn
	got := make([]byte, 16)
	applyRealityReportedVersion(got)
	if got[0] != 26 || got[1] != 10 || got[2] != 31 {
		t.Fatalf("Reality SessionId version = %d.%d.%d, want 26.10.31 to pass Xray default minClientVer", got[0], got[1], got[2])
	}
}

func TestRealityClientKeepsMLDSAVerifyKey(t *testing.T) {
	privateKey, err := ecdh.X25519().GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	config, err := NewRealityClient(context.Background(), logger.NOP(), "www.example.com", option.OutboundTLSOptions{
		Enabled:    true,
		ServerName: "www.example.com",
		UTLS:       &option.OutboundUTLSOptions{Enabled: true, Fingerprint: "chrome"},
		Reality: &option.OutboundRealityOptions{
			Enabled:        true,
			PublicKey:      base64.RawURLEncoding.EncodeToString(privateKey.PublicKey().Bytes()),
			ShortID:        "0123abcd",
			Mldsa65Verify:  "dW3KmTdkxdnW2nSvu_qY88c7ibp2-wPik9mKBf0kvA",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	client, ok := config.(*RealityClientConfig)
	if !ok {
		t.Fatalf("expected *RealityClientConfig, got %T", config)
	}
	if len(client.mldsa65Verify) == 0 {
		t.Fatal("expected decoded mldsa65Verify")
	}
}
