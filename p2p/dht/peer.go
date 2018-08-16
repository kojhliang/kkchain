package dht

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/bits"
	"time"

	"github.com/invin/kkchain/p2p"
)

type PeerID struct {
	p2p.ID
	Hash    []byte
	addTime time.Time
}

// CreateID is a factory function creating ID.
func CreateID(address string, publicKey []byte) PeerID {
	id := GetIDFromPublicKey(publicKey)

	return PeerID{ID: p2p.CreateID(address, publicKey), Hash: id}
}

// String returns the identity address and public key.
func (p PeerID) String() string {
	return fmt.Sprintf("ID{Address: %s, PublicKey: %s, id: %s}", p.Address, hex.EncodeToString(p.PublicKey), hex.EncodeToString(p.Hash))
}

// Equals determines if two peer IDs are equal to each other based on the contents of their public keys.
func (p PeerID) Equals(other PeerID) bool {
	return bytes.Equal(p.PublicKey, other.PublicKey) || bytes.Equal(p.Hash, other.Hash)
}

// Less determines if this peer ID's public key is less than other ID's public key.
func (p PeerID) Less(other interface{}) bool {
	if other, is := other.(PeerID); is {
		return bytes.Compare(p.Hash, other.Hash) == -1
	}
	return false
}

// PublicKeyHex generates a hex-encoded string of public key of this given peer ID.
func (p PeerID) PublicKeyHex() string {
	return hex.EncodeToString(p.PublicKey)
}

// Xor performs XOR (^) over another peer ID's public key.
func (p PeerID) Xor(other PeerID) PeerID {
	result := make([]byte, len(p.Hash))

	for i := 0; i < len(p.Hash) && i < len(other.Hash); i++ {
		result[i] = p.Hash[i] ^ other.Hash[i]
	}
	return PeerID{ID: p.ID, Hash: result}
}

// PrefixLen returns the number of prefixed zeros in a peer ID.
func (p PeerID) PrefixLen() int {
	for i, b := range p.Hash {
		if b != 0 {
			return i*8 + bits.LeadingZeros8(uint8(b))
		}
	}
	return len(p.Hash)*8 - 1
}

func GetIDFromPublicKey(publicKey []byte) []byte {

	h := sha256.New()
	h.Write(publicKey)
	return h.Sum(nil)
}

func (p PeerID) HashHex() string {
	return hex.EncodeToString(p.Hash)
}
