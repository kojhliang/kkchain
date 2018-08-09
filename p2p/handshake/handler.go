package handshake

import (
	"context"

	"github.com/invin/kkchain/p2p"
	"github.com/invin/kkchain/p2p/handshake/pb"
)

// handshakeHandler specifies the signature of functions that handle DHT messages.
type handshakeHandler func(context.Context, p2p.ID, *pb.Message) (*pb.Message, error)

func (hs *Handshake) handlerForMsgType(t pb.Message_Type) handshakeHandler {
	switch t {
	case pb.Message_HELLO:
		return hs.handleHello
	case pb.Message_HELLO_OK:
		return hs.handleHelloOK
	case pb.Message_HELLO_ERROR:
		return hs.handleHelloError
	default:
		return nil
	}
}

func (hs *Handshake) handleHello(ctx context.Context, p p2p.ID, pmes *pb.Message) (_ *pb.Message, err error) {
	// Check result and return corresponding code
	ok := true
	var resp *pb.Message

	if ok {
		resp = pb.NewMessage(pb.Message_HELLO_OK)
		// TODO: fill with local info to finish handshaking
		pb.BuildHandshake(resp)
	} else {
		resp = pb.NewMessage(pb.Message_HELLO_ERROR)
		// TODO: set error info
		// resp.Error.code = XXX
		// resp.Error.desc = "balabla"
	}

	return resp, nil
}

func (hs *Handshake) handleHelloOK(ctx context.Context, p p2p.ID, pmes *pb.Message) (_ *pb.Message, err error) {

	// TODO: handle received handshake info

	return nil, nil
}

func (hs *Handshake) handleHelloError(ctx context.Context, p p2p.ID, pmes *pb.Message) (_ *pb.Message, err error) {

	// TODO: handle error

	return nil, nil
}
