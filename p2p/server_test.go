package p2p

import (
	"testing"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

func TestServerListen(t *testing.T) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Error("failed to generate privkey:", err)
	}

	server := &Server{
		Config: Config{
			MaxPeers:   10,
			ListenAddr: "127.0.0.1:99999",
			PrivateKey: privKey,
		},
	}

	err = server.Start()
	if err != nil {
		t.Error("failed to start server:", err)
	}
}
