package chain 

import (
	"context"

	"github.com/invin/kkchain/p2p/chain/pb"
	"github.com/invin/kkchain/p2p"
)

// chainHandler specifies the signature of functions that handle DHT messages.
type chainHandler func(context.Context, p2p.ID, *pb.Message) (*pb.Message, error)

func (c *Chain) handlerForMsgType(t pb.Message_Type) chainHandler {
	switch t {
	case pb.Message_GET_BLOCK:
		return c.handleGetBlock
	default:
		return nil
	}
}

func (c *Chain) handleGetBlock(ctx context.Context, p p2p.ID, pmes *pb.Message) (_ *pb.Message, err error) {
	// Check result and return corresponding code
	var resp *pb.Message

	// TODO: 
	return resp, nil
}

