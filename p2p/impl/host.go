package impl

import (
	"sync"

	"net"

	"github.com/invin/kkchain/p2p"
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
		id:   id,
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
func (h *Host) GetConnection(id p2p.ID) (p2p.Conn, error) {
	h.mux.Lock()
	defer h.mux.Unlock()

	pk := string(id.PublicKey)
	conn, ok := h.cMap[pk]

	if !ok {
		return nil, errConnectionNotFound
	}

	return conn, nil
}

func (h *Host) GetAllConnection() map[string]p2p.Conn {
	h.mux.Lock()
	defer h.mux.Unlock()

	return h.cMap
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

func (h *Host) RemoveAllConnection() {
	h.mux.Lock()
	defer h.mux.Unlock()

	for id, _ := range h.cMap {
		delete(h.cMap, id)
	}
}

// ID returns the local ID
func (h *Host) ID() p2p.ID {
	return h.id
}

// Connect connects to remote peer
func (h *Host) Connect(address string) (net.Conn, error) {
	fd, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	return fd, nil
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
