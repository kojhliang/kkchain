package main

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/invin/kkchain/crypto/ed25519"
	"github.com/invin/kkchain/p2p/dht"
	"github.com/invin/kkchain/p2p/impl"
)

func testConfig() *dht.DHTConfig {
	return &dht.DHTConfig{
		BucketSize:      16,
		RoutingTableDir: "/Users/walker/Work/dht.db",
		BootstrapNodes:  []string{},
	}
}

func main() {
	//for i := 0; i < 10; i++ {
	//	kp := ed25519.RandomKeyPair()
	//	fmt.Printf("key: %v, len: %d\n", hex.EncodeToString(kp.PublicKey), len(kp.PublicKey))
	//}

	//key, _ := p2p.GenerateKey(libcrypto.Secp256k1)
	//pbKey, _ := key.GetPublic().Bytes()

	kp := ed25519.RandomKeyPair()
	pbKey := kp.PublicKey
	fmt.Printf("key: %v, len: %d\n", hex.EncodeToString(pbKey), len(pbKey))

	pbKey, _ = hex.DecodeString("c0bc7c08b52dc1df44dcc450068171f7039ea89c6ce9c678908a4f76c5f8d2f4")
	peerA := dht.CreateID("/ip4/127.0.0.1/tcp/8860", pbKey)
	fmt.Printf("peer A: %s\n", peerA.String())

	host := impl.NewHost(peerA.ID)

	kad := dht.NewDHT(testConfig(), host)

	for i := 0; i < 1; i++ {
		//key, _ = p2p.GenerateKey(libcrypto.Secp256k1)
		//pbKey, _ = key.GetPublic().Bytes()
		kp := ed25519.RandomKeyPair()
		pbKey := kp.PublicKey
		peerB := dht.CreateID("/ip4/127.0.0.1/tcp/8861", pbKey)
		fmt.Printf("peer B: %s\n", peerB.String())
		kad.AddPeer(peerB)
	}

	kad.Start()

	time.Sleep(time.Duration(80 * time.Second))

	kad.Stop()
}
