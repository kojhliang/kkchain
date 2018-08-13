package handshake

import (
	"context"

	"github.com/invin/kkchain/p2p"
	"github.com/pkg/errors"
)

// handshakeHandler specifies the signature of functions that handle DHT messages.
type handshakeHandler func(context.Context, p2p.ID, *Message) (*Message, error)

func (hs *Handshake) handlerForMsgType(t Message_Type) handshakeHandler {
	switch t {
	case Message_HELLO:
		return hs.handleHello
	case Message_HELLO_OK:
		return hs.handleHelloOK
	case Message_HELLO_ERROR:
		return hs.handleHelloError
	default:
		return nil
	}
}

func (hs *Handshake) handleHello(ctx context.Context, p p2p.ID, pmes *Message) (_ *Message, err error) {
	// Check result and return corresponding code
	ok := true
	var resp *Message

	if ok {
		resp = NewMessage(Message_HELLO_OK)
		// TODO: fill with local info to finish handshaking
		BuildHandshake(resp)
		return resp, nil
	}
	resp = NewMessage(Message_HELLO_ERROR)
	// TODO: set error info
	// resp.Error.code = XXX
	// resp.Error.desc = "balabla"

	return resp, errors.New("hello error")
}

func (hs *Handshake) handleHelloOK(ctx context.Context, p p2p.ID, pmes *Message) (_ *Message, err error) {

	// TODO: handle received handshake info

	return nil, nil
}

func (hs *Handshake) handleHelloError(ctx context.Context, p p2p.ID, pmes *Message) (_ *Message, err error) {

	// TODO: handle error

	return nil, nil
}
