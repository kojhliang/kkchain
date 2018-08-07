package impl 

import (
	"errors"

	"github.com/invin/kkchain/p2p"
)

var (
	errDuplicateStream = errors.New("duplicated stream")
	errStreamNotFound = errors.New("stream not found")
)

// Host defines a host for connections
type Host struct {
	address string
	// TODO: add connection manager
	m map[string]p2p.StreamHandler
}

// NewHost creates a new host object
func NewHost(address string) *Host {
	return &Host{
		address: address,
		m: make(map[string]p2p.StreamHandler),
	}
}

// SetStreamHandler sets handler for some a stream
func (h *Host) SetStreamHandler(protocol string, handler p2p.StreamHandler) error {
	_, found := h.m[protocol]
	if found {
		return errDuplicateStream
	}

	h.m[protocol] = handler
	return nil
}

// GetStreamHandler get stream handler
func (h *Host) GetStreamHandler(protocol string) (p2p.StreamHandler, error) {
	handler, ok := h.m[protocol]
	if !ok {
		return nil, errStreamNotFound
	}

	return handler, nil
}
