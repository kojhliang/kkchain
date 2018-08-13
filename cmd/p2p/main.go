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

	key, _ := dht.GenerateKey(libcrypto.Secp256k1)
	pbKey, _ := key.GetPublic().Bytes()
	fmt.Printf("key: %v, len: %d\n", hex.EncodeToString(pbKey), len(pbKey))

	peerA := dht.CreateID("/ip4/127.0.0.1/tcp/8860", pbKey)
	fmt.Printf("peer A: %s\n", peerA.String())

	kad := dht.NewDHT(p2p.DefaultConfig())

	kad.SetHost(impl.NewHost(kad.Self().ID))

	key, _ = dht.GenerateKey(libcrypto.Secp256k1)
	pbKey, _ = key.GetPublic().Bytes()
	peerB := dht.CreateID("/ip4/127.0.0.1/tcp/8861", pbKey)
	fmt.Printf("peer B: %s\n", peerB.String())

	kad.AddPeer(peerB)

	kad.Start()

	time.Sleep(time.Duration(30 * time.Second))

	kad.Stop()

}
