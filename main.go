package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/invin/kkchain/crypto/blake2b"
	"github.com/invin/kkchain/crypto/ed25519"
	"github.com/invin/kkchain/p2p"
	"github.com/invin/kkchain/p2p/impl"
)

func main() {
	port := flag.String("p", "9999", "")
	sec := flag.Int("s", 120, "")
	flag.Parse()

	config := p2p.Config{
		SignaturePolicy: ed25519.New(),
		HashPolicy:      blake2b.New(),
	}

	listen := "/ip4/127.0.0.1/tcp/" + *port

	net := impl.NewNetwork(listen, config)
	seconds := time.Duration(*sec)

	peerPort := "9999"
	if *port == "9999" {
		peerPort = "9998"
		node := "/ip4/127.0.0.1/tcp/" + peerPort
		node = "c0bc7c08b52dc1df44dcc450068171f7039ea89c6ce9c678908a4f76c5f8d2f4" + "@" + node
		net.BootstrapNodes = []string{node}
	}

	err := net.Start()
	if err != nil {
		fmt.Println("failed to start server:", err)
	}

	time.Sleep(time.Second * seconds)
}
