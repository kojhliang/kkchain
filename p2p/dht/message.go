package dht

import (
	"encoding/hex"
)

// NewMessage creates a new message object
func NewMessage(typ Message_Type, key string) *Message {
	m := &Message{
		Type: typ,
		Key:  key,
	}

	return m
}

func PeerIDToPBPeer(p PeerID) *Message_Peer {

	addrs := make([][]byte, 0)
	addrs = append(addrs, []byte(p.Address))
	return &Message_Peer{Id: p.PublicKeyHex(), Addrs: addrs}
}

func PBPeerToPeerID(p Message_Peer) *PeerID {
	pubk, err := hex.DecodeString(p.Id)
	if err != nil || len(p.Addrs) != 1 {
		return nil
	}

	id := CreateID(string(p.Addrs[0]), pubk)
	return &id
}
