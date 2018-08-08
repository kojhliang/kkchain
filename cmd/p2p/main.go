package main

import (
	"fmt"
	"encoding/hex"
	"time"

	"github.com/invin/kkchain/p2p/dht"

	libcrypto "github.com/libp2p/go-libp2p-crypto"

)

func main()  {

	key, _ := dht.GenerateKey(libcrypto.Secp256k1)
	pbKey,_ := key.GetPublic().Bytes()
	fmt.Printf("key: %v, len: %d\n", hex.EncodeToString(pbKey), len(pbKey))

	peerA := dht.CreateID("111", pbKey)
	fmt.Printf("peer A: %s\n", peerA.String())

	kad,_ := dht.NewDHT("", peerA)

	key, _ = dht.GenerateKey(libcrypto.Secp256k1)
	pbKey,_ = key.GetPublic().Bytes()
	peerB := dht.CreateID("222", pbKey)
	fmt.Printf("peer B: %s\n", peerB.String())

	kad.AddPeer(peerB)

	kad.Start()

	time.Sleep(time.Duration(30*time.Second))

	kad.Stop()

}
