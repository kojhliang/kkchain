package main

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/invin/kkchain/p2p"
	"github.com/invin/kkchain/p2p/dht"
	"github.com/invin/kkchain/p2p/impl"
	libcrypto "github.com/libp2p/go-libp2p-crypto"
)

func testConfig() *dht.DHTConfig {
	return &dht.DHTConfig{
		BucketSize:      16,
		RoutingTableDir: "/Users/walker/Work/dht.db",
		BootstrapNodes:  []string{},
	}
}

func main() {

	key, _ := p2p.GenerateKey(libcrypto.Secp256k1)
	pbKey, _ := key.GetPublic().Bytes()
	fmt.Printf("key: %v, len: %d\n", hex.EncodeToString(pbKey), len(pbKey))

	pbKey, _ = hex.DecodeString("08021221028e3f32ac3a18c594fa2be312c8d417e0d1bcb8f5980d69dfa31974e4715cf373")
	peerA := dht.CreateID("/ip4/127.0.0.1/tcp/8860", pbKey)
	fmt.Printf("peer A: %s\n", peerA.String())

	host := impl.NewHost(peerA.ID)

	kad := dht.NewDHT(testConfig(), host)

	for i := 0; i < 1; i++ {
		key, _ = p2p.GenerateKey(libcrypto.Secp256k1)
		pbKey, _ = key.GetPublic().Bytes()
		peerB := dht.CreateID("/ip4/127.0.0.1/tcp/8861", pbKey)
		fmt.Printf("peer B: %s\n", peerB.String())
		kad.AddPeer(peerB)
	}

	kad.Start()

	time.Sleep(time.Duration(80 * time.Second))

	kad.Stop()
}
