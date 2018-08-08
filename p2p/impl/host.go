package impl 

import (
	"errors"
	"sync"

	"github.com/invin/kkchain/p2p"
)

var (
	errDuplicateConnection = errors.New("duplicated connection")
	errDuplicateStream = errors.New("duplicated stream")
	errConnectionNotFound = errors.New("connection not found")
	errStreamNotFound = errors.New("stream not found")
)

// Host defines a host for connections
type Host struct {
	id p2p.ID

	// connection map
	cMap map[string]p2p.Conn

	// stream handler map
	sMap map[string]p2p.StreamHandler

	// mutex to sync access
	mux sync.Mutex
}

// NewHost creates a new host object
func NewHost(id p2p.ID) *Host {
	return &Host{
		id: id,
		cMap: make(map[string]p2p.Conn),
		sMap: make(map[string]p2p.StreamHandler),
	}
}

// AddConnection a connection
func (h *Host) AddConnection(id p2p.ID, conn p2p.Conn) error {
	h.mux.Lock()
	defer h.mux.Unlock()

	pk := string(id.PublicKey)
	_, found := h.cMap[pk]

	if found {
		return errDuplicateConnection
	}

	h.cMap[pk] = conn
	return nil
}

// GetConnection get a connection with ID
func (h *Host) GetConnection(id p2p.ID)(p2p.Conn, error) {
	h.mux.Lock()
	defer h.mux.Unlock()

	pk := string(id.PublicKey)
	conn, ok := h.cMap[pk]

	if !ok {
		return nil, errConnectionNotFound
	}

	return conn, nil	
}

// RemoveConnection removes a connection
func (h *Host) RemoveConnection(id p2p.ID) error {
	h.mux.Lock()
	defer h.mux.Unlock()

	pk := string(id.PublicKey)
	_, ok := h.cMap[pk]

	if !ok {
		return errConnectionNotFound
	}

	delete(h.cMap, pk)

	return nil	
}

// ID returns the local ID
func (h *Host) ID() p2p.ID {
	return h.id
}

// Connect connects to remote peer
func (h *Host) Connect(address string) error {
	// TODO: connect to remote peer if not exist
	
	return nil
}

// SetStreamHandler sets handler for some a stream
func (h *Host) SetStreamHandler(protocol string, handler p2p.StreamHandler) error {
	h.mux.Lock()
	defer h.mux.Unlock()

	_, found := h.sMap[protocol]
	if found {
		return errDuplicateStream
	}

	h.sMap[protocol] = handler
	return nil
}

// GetStreamHandler get stream handler
func (h *Host) GetStreamHandler(protocol string) (p2p.StreamHandler, error) {
	h.mux.Lock()
	defer h.mux.Unlock()

	handler, ok := h.sMap[protocol]

	if !ok {
		return nil, errStreamNotFound
	}

	return handler, nil
}
