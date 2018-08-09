package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"flag"
	"time"
	"net"
	"strconv"
	"github.com/invin/kkchain/p2p"
)

func main() {
	port := flag.String("p", "9999", "")
	sec := flag.Int("s", 60, "")
	flag.Parse()
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		fmt.Println("failed to generate privkey:", err)
	}

	seconds := time.Duration(*sec)

	peerPort, err := strconv.ParseInt(*port, 10, 16)
	if err != nil {
		fmt.Println("failed to parse int:", err)
	}
	node := &p2p.Node{
		net.ParseIP("127.0.0.1"),
		uint16(peerPort + 1),
		p2p.NodeID{},
	}

	server := &p2p.Server{
		Config: p2p.Config{
			MaxPeers:       10,
			ListenAddr:     "127.0.0.1:" + *port,
			PrivateKey:     privKey,
			BootstrapNodes: []*p2p.Node{node},
		},
		SendMsg: make(map[string]interface{}),
		RecvMsg: make(map[string]interface{}),
	}

	err = server.Start()
	if err != nil {
		fmt.Println("failed to start server:", err)
	}

	time.Sleep(time.Second * seconds)
}
