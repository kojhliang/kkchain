package impl

import (
	"sync"

	"net"

	"github.com/invin/kkchain/p2p"
	"github.com/invin/kkchain/p2p/dht"
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

	// notification listeners
	nMap map[p2p.Notifiee]struct{}

	// notifier mux
	notifyMux sync.Mutex
}

// NewHost creates a new host object
func NewHost(id p2p.ID) *Host {
	return &Host{
		id:   id,
		cMap: make(map[string]p2p.Conn),
		sMap: make(map[string]p2p.StreamHandler),
		nMap: make(map[p2p.Notifiee]struct{}),
	}
}

// Register registers a notifier
func (h *Host) Register(n p2p.Notifiee) error {
	h.notifyMux.Lock()
	defer h.notifyMux.Unlock()

	_, found := h.nMap[n]
	if found {
		return errDuplicateNotifiee
	}

	h.nMap[n] = struct{}{}

	return nil
}

// Revoke revokes a notifier
func (h *Host) Revoke(n p2p.Notifiee) error {
	h.notifyMux.Lock()
	defer h.notifyMux.Unlock()

	_, found := h.nMap[n]
	if !found {
		return errNotifieeNotFound
	}

	delete(h.nMap, n)

	return nil
}

func (h *Host) NotifyAll(notification func(n p2p.Notifiee)) {
	h.notifyMux.Lock()
	defer h.notifyMux.Unlock()

	var wg sync.WaitGroup
	for n := range h.nMap {
		wg.Add(1)
		go func(n p2p.Notifiee) {
			defer wg.Done()
			notification(n)
		}(n)
	}

	wg.Wait()
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

	// notify a new conn
	h.NotifyAll(func(n p2p.Notifiee) {
		n.Connected(conn)
	})
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
	conn, ok := h.cMap[pk]

	if !ok {
		return errConnectionNotFound
	}
	conn.Close()
	delete(h.cMap, pk)

	return nil
}

func (h *Host) RemoveAllConnection() {
	h.mux.Lock()
	defer h.mux.Unlock()

	for id, conn := range h.cMap {
		conn.Close()
		delete(h.cMap, id)
	}
}

// ID returns the local ID
func (h *Host) ID() p2p.ID {
	return h.id
}

// Connect connects to remote peer
func (h *Host) Connect(address string, network p2p.Network) (p2p.Conn, error) {
	addr, err := dht.ToNetAddr(address)
	if err != nil {
		return nil, err
	}

	fd, err := net.Dial(addr.Network(), addr.String())
	if err != nil {
		return nil, err
	}

	conn := NewConnection(fd, network, h)
	if conn != nil {
		return conn, nil
	}

	return nil, failedNewConnection
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
