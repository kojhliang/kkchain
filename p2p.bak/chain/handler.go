package chain

import (
	"context"

	"github.com/invin/kkchain/p2p"
)

// chainHandler specifies the signature of functions that handle DHT messages.
type chainHandler func(context.Context, p2p.ID, *Message) (*Message, error)

func (c *Chain) handlerForMsgType(t Message_Type) chainHandler {
	switch t {
	case Message_GET_BLOCK:
		return c.handleGetBlock
	default:
		return nil
	}
}

func (c *Chain) handleGetBlock(ctx context.Context, p p2p.ID, pmes *Message) (_ *Message, err error) {
	// Check result and return corresponding code
	var resp *Message

	// TODO:
	return resp, nil
}
