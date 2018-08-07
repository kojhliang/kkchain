package p2p

import (
	"net"

	"github.com/invin/kkchain/crypto"
)

// Config defines configurations for a basic p2p instance
type Config struct {
	SignaturePolicy	crypto.SignaturePolicy
	HashPolicy crypto.HashPolicy	
}

// P2P defines the interface for p2p communication
type P2P interface {
	// Start kicks off the p2p stack
	Start() error

	// Get configurations
	Conf() Config

	// Stop the p2p stack
	Stop()

	// Accept connection
	Accept(incoming net.Conn)

	// returns ID of local peer
	ID() ID
	
}
