package dht

import (
	"context"
	"encoding/hex"

	"github.com/invin/kkchain/p2p"
)

// dhthandler specifies the signature of functions that handle DHT messages.
type dhtHandler func(context.Context, p2p.ID, *Message) (*Message, error)

func (dht *DHT) handlerForMsgType(t Message_Type) dhtHandler {
	switch t {
	case Message_GET_VALUE:
		return dht.handleGetValue
	case Message_PUT_VALUE:
		return dht.handlePutValue
	case Message_FIND_NODE:
		return dht.handleFindPeer
	case Message_FIND_NODE_RESULT:
		return dht.handleFindPeerResult
	case Message_PING:
		return dht.handlePing
	case Message_PONG:
		return dht.handlePong
	default:
		return nil
	}
}

func (dht *DHT) handleGetValue(ctx context.Context, p p2p.ID, pmes *Message) (_ *Message, err error) {
	// TODO: setup response
	return nil, nil
}

func (dht *DHT) handlePutValue(ctx context.Context, p p2p.ID, pmes *Message) (_ *Message, err error) {
	// TODO: setup response
	return nil, nil
}

func (dht *DHT) handleFindPeer(ctx context.Context, p p2p.ID, pmes *Message) (_ *Message, err error) {
	// setup response
	resp := NewMessage(Message_FIND_NODE_RESULT, "")

	target, err := hex.DecodeString(pmes.Key)
	if err != nil {
		return nil, err
	}

	peer := PeerID{Hash: target}
	closest := dht.table.FindClosestPeers(peer, BucketSize)

	for _, p := range closest {
		resp.CloserPeers = append(resp.CloserPeers, PeerIDToPBPeer(p))
	}

	return resp, nil
}

func (dht *DHT) handleFindPeerResult(ctx context.Context, p p2p.ID, pmes *Message) (_ *Message, err error) {

	pbPeers := pmes.CloserPeers
	for _, p := range pbPeers {
		peer := PBPeerToPeerID(*p)
		dht.table.Update(*peer)
	}

	return nil, nil
}

func (dht *DHT) handlePing(ctx context.Context, p p2p.ID, pmes *Message) (_ *Message, err error) {
	// setup response
	resp := NewMessage(Message_PONG, "")

	return resp, nil
}

func (dht *DHT) handlePong(ctx context.Context, p p2p.ID, pmes *Message) (_ *Message, err error) {
	// TODO: update connection status
	return nil, nil
}
