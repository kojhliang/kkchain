package impl

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	errDuplicateConnection = errors.New("duplicated connection")
	errDuplicateStream     = errors.New("duplicated stream")
	errConnectionNotFound  = errors.New("connection not found")
	errStreamNotFound      = errors.New("stream not found")
	failedNewConnection    = errors.New("failed to new connection")
	failedNewNetwork       = errors.New("failed to new network")
	failedCreateStream     = errors.New("failed to create stream")
	errDuplicateNotifiee   = errors.New("duplicated notifiee")
	errNotifieeNotFound    = errors.New("notifiee not found")
)

const (
	errInvalidMsgCode = iota
	errInvalidMsg
)

var errorToString = map[int]string{
	errInvalidMsgCode: "invalid message code",
	errInvalidMsg:     "invalid message",
}

type DisconnectReason uint

const (
	DiscRequested DisconnectReason = iota
	DiscNetworkError
	DiscProtocolError
	DiscUselessPeer
	DiscTooManyPeers
	DiscAlreadyConnected
	DiscIncompatibleVersion
	DiscInvalidIdentity
	DiscQuitting
	DiscUnexpectedIdentity
	DiscSelf
	DiscReadTimeout
	DiscSubprotocolError = 0x10
)

var discReasonToString = [...]string{
	DiscRequested:           "disconnect requested",
	DiscNetworkError:        "network error",
	DiscProtocolError:       "breach of protocol",
	DiscUselessPeer:         "useless peer",
	DiscTooManyPeers:        "too many peers",
	DiscAlreadyConnected:    "already connected",
	DiscIncompatibleVersion: "incompatible p2p protocol version",
	DiscInvalidIdentity:     "invalid node identity",
	DiscQuitting:            "client quitting",
	DiscUnexpectedIdentity:  "unexpected identity",
	DiscSelf:                "connected to self",
	DiscReadTimeout:         "read timeout",
	DiscSubprotocolError:    "subprotocol error",
}

func (d DisconnectReason) String() string {
	if len(discReasonToString) < int(d) {
		return fmt.Sprintf("unknown disconnect reason %d", d)
	}
	return discReasonToString[d]
}

func (d DisconnectReason) Error() string {
	return d.String()
}
