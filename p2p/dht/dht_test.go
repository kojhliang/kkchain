package dht

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/invin/kkchain/p2p"
	"github.com/invin/kkchain/p2p/impl"
	libcrypto "github.com/libp2p/go-libp2p-crypto"
)

func TestNewDHT(t *testing.T) {
	t.Parallel()

	key, _ := p2p.GenerateKey(libcrypto.Secp256k1)
	pbKey, _ := key.GetPublic().Bytes()
	fmt.Printf("key: %v, len: %d\n", hex.EncodeToString(pbKey), len(pbKey))

	peerA := CreateID("/ip4/127.0.0.1/tcp/8860", pbKey)
	fmt.Printf("peer A: %s\n", peerA.String())

	host := impl.NewHost(peerA.ID)

	dht := NewDHT(DefaultConfig(), host)

	key, _ = p2p.GenerateKey(libcrypto.Secp256k1)
	pbKey, _ = key.GetPublic().Bytes()
	peerB := CreateID("/ip4/127.0.0.1/tcp/8861", pbKey)
	fmt.Printf("peer B: %s\n", peerB.String())

	dht.AddPeer(peerB)

	dht.Start()

	time.Sleep(time.Duration(30 * time.Second))

	dht.Stop()
}
