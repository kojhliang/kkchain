package dht

import (
	"testing"

	"encoding/hex"
	"fmt"
	crypto "github.com/libp2p/go-libp2p-crypto"
)

func TestFormatAddr(t *testing.T) {
	t.Parallel()

	key, _ := GenerateKey(crypto.Secp256k1)
	pbKey, _ := key.GetPublic().Bytes()

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

	address := peeraddr("08021221022c8c52f66009dd50396dc49876742d37b26ddb5cb83006ee417112cf2da7edd7@/ip4/127.0.0.1/tcp/8860")

	id, err := ParsePeerAddr(address)
	if err != nil {
		t.Errorf("ParseAddress() = <nil>, expected %+v (%s)\n", err, address)
	}

	t.Logf("%s\n", id.String())
}
