package p2p

// Host defines a host for connections
type Host interface {
	// ID returns the peer ID associated with this host
	// ID() string

	// Returns the address of the host
	// Address() string

	// set stream handler
	SetStreamHandler(protocol string, handler StreamHandler) error

	// returns the stream for protocol
	GetStreamHandler(protocol string) (StreamHandler, error)
}