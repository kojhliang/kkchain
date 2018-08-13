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
	sec := flag.Int("s", 60, "")
	flag.Parse()

	config := p2p.Config{
		SignaturePolicy: ed25519.New(),
		HashPolicy:      blake2b.New(),
	}
	net := impl.NewNetwork("127.0.0.1:"+*port, config)
	seconds := time.Duration(*sec)

	peerPort := "9999"
	if *port == "9999" {
		peerPort = "9998"
	}

	node := &impl.Node{
		"127.0.0.1",
		peerPort,
		p2p.ID{},
	}

	net.BootstrapNodes = []*impl.Node{node}

	err := net.Start()
	if err != nil {
		fmt.Println("failed to start server:", err)
	}

	time.Sleep(time.Second * seconds)
}
