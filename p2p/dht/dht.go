package dht

import (
	"context"

	"github.com/invin/kkchain/p2p"
	"github.com/invin/kkchain/p2p/dht/pb"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"
)

const (
	protocolDHT = "/kkchain/p2p/dht/1.0.0"
)

// DHT implements a Distributed Hash Table for p2p
type DHT struct {
	// self
	host p2p.Host

	quitCh chan bool
	table *RoutingTable
	store *PeerStore
	seedPeers []PeerID
}

// NewDHT creates a new DHT object with the given peer as as the 'local' host
func NewDHT(host p2p.Host, dbPath string, self PeerID) *DHT {

	// If no node database was given, use an in-memory one
	db, err := newPeerStore(dbPath)
	if err != nil {
		return nil
	}

	dht := &DHT {
		host: host,
		quitCh: make(chan bool),
		table: CreateRoutingTable(self),
		store: db,
	}

	if err := host.SetStreamHandler(protocolDHT, dht.handleNewStream); err != nil {
		panic(err)
	}

	return dht
}



// handleNewStream handles messages within the stream
func (dht *DHT) handleNewStream(s p2p.Stream, msg proto.Message) {
	// check message type
	switch message := msg.(type) {
	case *pb.Message:
		dht.handleMessage(s, message)
	default:
		s.Reset()
		glog.Errorf("unexpected message: %v", msg)
	}	
}

// handleMessage handles messsage 
func (dht *DHT) handleMessage(s p2p.Stream, msg *pb.Message) {
	// get handler
	handler := dht.handlerForMsgType(msg.GetType())
	if handler == nil {
		s.Reset()
		glog.Errorf("unknown message type: %v", msg.GetType())
		return	
	}

	// dispatch handler
	// TODO: get context and peer id
	ctx := context.Background()
	pid := s.RemotePeer()

	rpmes, err := handler(ctx, pid, msg)

	// if nil response, return it before serializing
	if rpmes == nil {
		glog.Warning("got back nil response from request")
		return
	}

	// send out response msg
	if err = s.Write(rpmes); err != nil {
		s.Reset()
		glog.Errorf("send response error: %s", err)
		return
	}
}