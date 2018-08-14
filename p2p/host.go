package p2p

import "net"

// ConnManager defines interface to manage connections
type ConnManager interface {
	// Add a connection
	AddConnection(id ID, conn Conn) error

	// Get a connection with ID
	GetConnection(id ID) (Conn, error)

	// Get all connection
	GetAllConnection() map[string]Conn

	// Remove a connection
	RemoveConnection(id ID) error

	// Remove all connection
	RemoveAllConnection()
}

// Host defines a host for connections
type Host interface {
	ConnManager

	// Returns ID of local peer
	ID() ID

	// Connect to remote peer
	Connect(address string) (net.Conn, error)

	// Set stream handler
	SetStreamHandler(protocol string, handler StreamHandler) error

	// returns the stream for protocol
	GetStreamHandler(protocol string) (StreamHandler, error)
}
