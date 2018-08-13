package dht

import (
	"context"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"
	"github.com/invin/kkchain/p2p"
	"github.com/invin/kkchain/p2p/dht/pb"
	"github.com/libp2p/go-libp2p-crypto"
)

const (
	protocolDHT = "/kkchain/p2p/dht/1.0.0"
)

// DHT implements a Distributed Hash Table for p2p
type DHT struct {
	// self
	host p2p.Host

	quitCh         chan bool
	table          *RoutingTable
	store          *PeerStore
	config         *p2p.Config
	self           PeerID
	selfPrivateKey crypto.PrivKey
	recvCh         chan interface{}
}

func (dht *DHT) GetRecvchan() chan interface{} {
	return dht.recvCh
}

// NewDHT creates a new DHT object with the given peer as as the 'local' host
func NewDHT(config *p2p.Config) *DHT {

	// If no node database was given, use an in-memory one
	db, err := newPeerStore(config.RoutingTableDir)
	if err != nil {
		return nil
	}

	self, privateKey, err := selfPeerID(config)
	if err != nil {
		return nil
	}

	dht := &DHT{
		quitCh:         make(chan bool),
		config:         config,
		selfPrivateKey: privateKey,
		self:           *self,
		table:          CreateRoutingTable(*self),
		store:          db,
		recvCh:         make(chan interface{}),
	}

	//initNetwork(config, dht)

	if err := dht.host.SetStreamHandler(protocolDHT, dht.handleNewStream); err != nil {
		panic(err)
	}

	return dht
}

//initNetwork init network and related config before startup
//func initNetwork(config *p2p.Config, dht *DHT) {
//	//TODO: new server & and set to dht
//
//	//dht.host = impl.NewHost(dht.self.ID)
//}

// selfPeerID get local peer info
func selfPeerID(config *p2p.Config) (*PeerID, crypto.PrivKey, error) {

	peerKey, err := LoadNetworkKeyFromFileOrCreateNew(config.PrivateKeyPath)
	if err != nil {
		return nil, nil, err
	}

	pubk, err := peerKey.GetPublic().Bytes()
	if err != nil {
		return nil, nil, err
	}
	id := CreateID(config.Listen, pubk)

	return &id, peerKey, nil
}

func (dht *DHT) Self() PeerID {
	return dht.self
}

func (dht *DHT) SetHost(host p2p.Host) {
	dht.host = host
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
