package impl

import (
	"github.com/invin/kkchain/p2p"
)

type Node struct {
	IP      string
	TCPPort string
	ID      p2p.ID
}

func (n *Node) Addr() string {
	return n.IP + ":" + n.TCPPort
}
