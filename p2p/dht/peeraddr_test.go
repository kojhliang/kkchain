package dht

import (
	"testing"

	"encoding/hex"
	"fmt"
	"github.com/invin/kkchain/crypto/ed25519"
)

func TestFormatAddr(t *testing.T) {
	t.Parallel()

	//key, _ := p2p.GenerateKey(crypto.Secp256k1)
	//pbKey, _ := key.GetPublic().Bytes()

	kp := ed25519.RandomKeyPair()
	pbKey := kp.PublicKey

	peerA := CreateID("/ip4/127.0.0.1/tcp/8860", pbKey)

	address := FormatPeerAddr(peerA)

	expected := fmt.Sprintf("%s@%s", hex.EncodeToString(pbKey), "/ip4/127.0.0.1/tcp/8860")
	if address != expected {
		t.Errorf("FormatPeerAddr() = %+v, expected %+v\n", address, expected)
	}
	t.Logf("address :%s\n", address)
}

func TestParsePeerAddr(t *testing.T) {
	t.Parallel()

	address := peeraddr("c0bc7c08b52dc1df44dcc450068171f7039ea89c6ce9c678908a4f76c5f8d2f4@/ip4/127.0.0.1/tcp/8860")

	id, err := ParsePeerAddr(address)
	if err != nil {
		t.Errorf("ParseAddress() = <nil>, expected %+v (%s)\n", err, address)
	}

	t.Logf("%s\n", id.String())
}

func TestToNetAddr(t *testing.T) {
	t.Parallel()

	address := "/ip4/127.0.0.1/tcp/8860"
	addr, err := ToNetAddr(address)
	if err != nil {
		t.Errorf("ToNetAddr() = <nil>, expected %+v (%s)\n", err, address)
	}

	t.Logf("%s, %s", addr.Network(), addr.String())
}
