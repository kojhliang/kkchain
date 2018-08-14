package dht

import (
	"encoding/hex"
	"errors"
	"net"
	"strings"

	ma "github.com/multiformats/go-multiaddr"
	manet "github.com/multiformats/go-multiaddr-net"
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

// ToNetAddr converts a Multiaddr string to a net.Addr
// Must be ThinWaist. acceptable protocol stacks are:
// /ip{4,6}/{tcp, udp}
func ToNetAddr(addr string) (net.Addr, error) {
	maddr, err := ma.NewMultiaddr(addr)
	if err != nil {
		return nil, err
	}
	return manet.ToNetAddr(maddr)
}

// FromNetAddr converts a net.Addr type to a Multiaddr.
func FromNetAddr(addr net.Addr) string {
	maddr, err := manet.FromNetAddr(addr)
	if err != nil {
		return ""
	}
	return maddr.String()
}
