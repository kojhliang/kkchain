package dht

import (
	"fmt"
	"bytes"
	"encoding/hex"
	"math/bits"
	"crypto/sha256"
	"time"
)

type PeerID struct {
	Address string

	PublicKey []byte
	ID []byte
	addTime time.Time
}

// CreateID is a factory function creating ID.
func CreateID(address string, publicKey []byte) PeerID {
	id := GetIDFromPublicKey(publicKey)
	return PeerID{Address: address, PublicKey: publicKey, ID: id}
}

// String returns the identity address and public key.
func (p PeerID) String() string {
	return fmt.Sprintf("ID{Address: %s, PublicKey: %s, id: %s}", p.Address, hex.EncodeToString(p.PublicKey), hex.EncodeToString(p.ID))
}

// Equals determines if two peer IDs are equal to each other based on the contents of their public keys.
func (p PeerID) Equals(other PeerID) bool {
	return bytes.Equal(p.PublicKey, other.PublicKey) || bytes.Equal(p.ID, other.ID)
}

// Less determines if this peer ID's public key is less than other ID's public key.
func (p PeerID) Less(other interface{}) bool {
	if other, is := other.(PeerID); is {
		return bytes.Compare(p.ID, other.ID) == -1
	}
	return false
}

// PublicKeyHex generates a hex-encoded string of public key of this given peer ID.
func (p PeerID) PublicKeyHex() string {
	return hex.EncodeToString(p.PublicKey)
}

// Xor performs XOR (^) over another peer ID's public key.
func (p PeerID) Xor(other PeerID) PeerID {
	result := make([]byte, len(p.ID))

	for i := 0; i < len(p.ID) && i < len(other.ID); i++ {
		result[i] = p.ID[i] ^ other.ID[i]
	}
	return PeerID{Address: p.Address, PublicKey: p.PublicKey, ID: result}
}

// PrefixLen returns the number of prefixed zeros in a peer ID.
func (p PeerID) PrefixLen() int {
	for i, b := range p.ID {
		if b != 0 {
			return i*8 + bits.LeadingZeros8(uint8(b))
		}
	}
	return len(p.ID)*8 - 1
}

func GetIDFromPublicKey(publicKey []byte) []byte {
	return sha256.New().Sum(publicKey)
}