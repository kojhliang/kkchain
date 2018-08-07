package p2p

import (
	"github.com/invin/kkchain/p2p/protobuf"
	"github.com/gogo/protobuf/proto"
)

// Conn wraps connection related operations, such as reading and writing 
// messages
type Conn interface {
	ReadMessage()(*protobuf.Message, error)
	WriteMessage(*protobuf.Message) error
	PrepareMessage(message proto.Message) (*protobuf.Message, error)
	GetPeerID() ID
}