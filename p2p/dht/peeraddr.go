package dht

import (
	"encoding/hex"
	"errors"
	"strings"
)

type peeraddr = string

var (
	errStrAddressEmpty   = errors.New("address was empty")
	errStrInvalidAddress = errors.New("invalid address")
)

func ParsePeerAddr(addr peeraddr) (*PeerID, error) {

	addr = strings.TrimSpace(addr)
	if len(addr) == 0 {
		return nil, errStrAddressEmpty
	}

	subAddrs := strings.Split(addr, "@")
	if len(subAddrs) != 2 {
		return nil, errStrInvalidAddress
	}

	pubk := subAddrs[0]
	url := subAddrs[1]
	pbPubk, err := hex.DecodeString(pubk)
	if err != nil {
		return nil, err
	}

	id := CreateID(url, pbPubk)
	return &id, nil

}

func FormatPeerAddr(id PeerID) peeraddr {

	subAddrs := []string{id.PublicKeyHex(), id.ID.Address}
	return strings.Join(subAddrs, "@")
}
