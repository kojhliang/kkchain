package main

import (
	"github.com/invin/kkchain/p2p/discovery"
	"fmt"
	"encoding/hex"
	"time"

	libcrypto "github.com/libp2p/go-libp2p-crypto"
)

func main()  {

	key, _ := discovery.GenerateKey(libcrypto.Secp256k1)
	pbKey,_ := key.GetPublic().Bytes()
	fmt.Printf("key: %v, len: %d\n", hex.EncodeToString(pbKey), len(pbKey))

	peerA := discovery.CreateID("111", pbKey)
	fmt.Printf("peer A: %s\n", peerA.String())

	dht,_ := discovery.NewDHT("", peerA)

	key, _ = discovery.GenerateKey(libcrypto.Secp256k1)
	pbKey,_ = key.GetPublic().Bytes()
	peerB := discovery.CreateID("222", pbKey)
	fmt.Printf("peer B: %s\n", peerB.String())

	dht.AddPeer(peerB)

	dht.Start()

	time.Sleep(time.Duration(30*time.Second))

	dht.Stop()

}
