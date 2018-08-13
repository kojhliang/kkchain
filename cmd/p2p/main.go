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

func main() {

	key, _ := p2p.GenerateKey(libcrypto.Secp256k1)
	pbKey, _ := key.GetPublic().Bytes()
	fmt.Printf("key: %v, len: %d\n", hex.EncodeToString(pbKey), len(pbKey))

	peerA := dht.CreateID("/ip4/127.0.0.1/tcp/8860", pbKey)
	fmt.Printf("peer A: %s\n", peerA.String())

	host := impl.NewHost(peerA.ID)

	kad := dht.NewDHT(dht.DefaultConfig(), host)

	key, _ = p2p.GenerateKey(libcrypto.Secp256k1)
	pbKey, _ = key.GetPublic().Bytes()
	peerB := dht.CreateID("/ip4/127.0.0.1/tcp/8861", pbKey)
	fmt.Printf("peer B: %s\n", peerB.String())

	kad.AddPeer(peerB)

	kad.Start()

	time.Sleep(time.Duration(30 * time.Second))

	kad.Stop()
}
