package p2p

import (
	"github.com/invin/kkchain/crypto"
)

// Config defines configurations for a basic p2p instance
type Config struct {
	signaturePolicy	crypto.SignaturePolicy
	hashPolicy crypto.HashPolicy	
}

// P2P represents the whole stack of p2p communication between peers
type P2P struct {
	conf Config
}

// New creates a new P2P instance with the specified configuration
func New(conf Config) *P2P {
	return &P2P{
		conf: conf,
	}
}
