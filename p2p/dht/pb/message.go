package pb

import (
	"github.com/invin/kkchain/p2p"
)

// NewMessage creates a new message object
func NewMessage(typ Message_Type) *Message {
	m := &Message{
		Type: typ,
	}

	return m
}

// FindCloserPeers find closer peers
func FindCloserPeers(pid p2p.ID) []*Message_Peer {
	// TODO: 	
	peers := make([]*Message_Peer, 32)
	return peers
}

