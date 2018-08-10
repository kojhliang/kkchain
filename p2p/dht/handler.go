package dht

import (
	"context"
	"github.com/invin/kkchain/p2p"
	"github.com/invin/kkchain/p2p/dht/pb"
)

// dhthandler specifies the signature of functions that handle DHT messages.
type dhtHandler func(context.Context, p2p.ID, *pb.Message) (*pb.Message, error)

func (dht *DHT) handlerForMsgType(t pb.Message_Type) dhtHandler {
	switch t {
	case pb.Message_GET_VALUE:
		return dht.handleGetValue
	case pb.Message_PUT_VALUE:
		return dht.handlePutValue
	case pb.Message_FIND_NODE:
		return dht.handleFindPeer
	case pb.Message_FIND_NODE_RESULT:
		return dht.handleFindPeerResult
	case pb.Message_PING:
		return dht.handlePing
	case pb.Message_PONG:
		return dht.handlePong
	default:
		return nil
	}
}

func (dht *DHT) handleGetValue(ctx context.Context, p p2p.ID, pmes *pb.Message) (_ *pb.Message, err error) {
	// TODO: setup response
	return nil, nil
}

func (dht *DHT) handlePutValue(ctx context.Context, p p2p.ID, pmes *pb.Message) (_ *pb.Message, err error) {
	// TODO: setup response
	return nil, nil
}

func (dht *DHT) handleFindPeer(ctx context.Context, p p2p.ID, pmes *pb.Message) (_ *pb.Message, err error) {
	// setup response
	resp := pb.NewMessage(pb.Message_FIND_NODE_RESULT, "")
	resp.CloserPeers = pb.FindCloserPeers(p)

	return resp, nil
}

func (dht *DHT) handleFindPeerResult(ctx context.Context, p p2p.ID, pmes *pb.Message) (_ *pb.Message, err error) {

	// TODO: handle received peers
	// handlePeers(pmesg.CloserPeers)

	return nil, nil
}

func (dht *DHT) handlePing(ctx context.Context, p p2p.ID, pmes *pb.Message) (_ *pb.Message, err error) {
	// setup response
	resp := pb.NewMessage(pb.Message_PONG, "")

	return resp, nil
}

func (dht *DHT) handlePong(ctx context.Context, p p2p.ID, pmes *pb.Message) (_ *pb.Message, err error) {
	// TODO: update connection status
	return nil, nil
}
