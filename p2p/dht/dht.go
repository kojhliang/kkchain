package dht

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"encoding/json"

	"sync"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"
	"github.com/invin/kkchain/p2p"
)

const (
	protocolDHT = "/kkchain/p2p/dht/1.0.0"
)

const (
	DefaultSyncTableInterval   = 20 * time.Second
	DefaultSaveTableInterval   = 1 * time.Minute
	DefaultSeedMinTableTime    = 50 * time.Second
	DefaultMaxPeersCountToSync = 6
)

type DHTConfig struct {
	BucketSize      int
	RoutingTableDir string
}

type PingPongService struct {
	mutex      *sync.RWMutex
	stopCh     map[string]chan interface{}
	pingpongAt map[string]time.Time
}

func newPingPongService() *PingPongService {
	time.Now()
	return &PingPongService{
		mutex:      &sync.RWMutex{},
		stopCh:     make(map[string]chan interface{}),
		pingpongAt: make(map[string]time.Time),
	}
}

// DHT implements a Distributed Hash Table for p2p
type DHT struct {
	// self
	host    p2p.Host
	network p2p.Network

	quitCh         chan bool
	table          *RoutingTable
	store          *PeerStore
	config         *DHTConfig
	self           PeerID
	BootstrapNodes []string
	pingpong       *PingPongService
}

func DefaultConfig() *DHTConfig {
	return &DHTConfig{
		BucketSize:      BucketSize,
		RoutingTableDir: "",
	}
}

// NewDHT creates a new DHT object with the given peer as as the 'local' host
func NewDHT(config *DHTConfig, network p2p.Network, host p2p.Host) *DHT {

	// If no node database was given, use an in-memory one
	db, err := newPeerStore(config.RoutingTableDir)
	if err != nil {
		return nil
	}

	self := CreateID(host.ID().Address, host.ID().PublicKey)

	dht := &DHT{
		quitCh:   make(chan bool),
		config:   config,
		self:     self,
		table:    CreateRoutingTable(self),
		store:    db,
		pingpong: newPingPongService(),
	}

	dht.host = host
	dht.network = network

	if err := dht.host.SetStreamHandler(protocolDHT, dht.handleNewStream); err != nil {
		panic(err)
	}

	host.Register(dht)

	return dht
}

func (dht *DHT) Self() PeerID {
	return dht.self
}

// handleNewStream handles messages within the stream
func (dht *DHT) handleNewStream(s p2p.Stream, msg proto.Message) {
	// check message type
	switch message := msg.(type) {
	case *Message:
		go dht.handleMessage(s, message)
	default:
		s.Reset()
		glog.Errorf("unexpected message: %v", msg)
	}
}

// handleMessage handles messsage
func (dht *DHT) handleMessage(s p2p.Stream, msg *Message) {
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

	fmt.Printf("dht handle %d success and send resp to: %s, protocol: %s, conn: %v", msg.Type, pid, s.Protocol(), s.Conn())
}

func (dht *DHT) Start() {
	dht.loadBootstrapNodes()

	//load table from db
	dht.loadTableFromDB()

	fmt.Println("start sync loop.....")
	go dht.syncLoop()
	go dht.checkPingPong()
}

func (dht *DHT) Stop() {
	fmt.Println("stopping sync loop.....")
	dht.quitCh <- true

}

func (dht *DHT) syncLoop() {

	dht.table.printTable()

	//first sync
	dht.SyncRouteTable()

	//TODO: config timer
	syncLoopTicker := time.NewTicker(DefaultSyncTableInterval)
	defer syncLoopTicker.Stop()

	saveTableToStore := time.NewTicker(DefaultSaveTableInterval)
	defer saveTableToStore.Stop()

	for {
		select {
		case <-dht.quitCh:
			fmt.Println("stopped sync loop")
			dht.store.Close()
			return
		case <-syncLoopTicker.C:
			dht.SyncRouteTable()
		case <-saveTableToStore.C:
			dht.saveTableToStore()
		}
	}
}

func (dht *DHT) AddPeer(peer PeerID) {

	//dht.store.Update(&peer)
	peer.addTime = time.Now()
	dht.table.Update(peer)

}

func (dht *DHT) RemovePeer(peer PeerID) {

	dht.store.Delete(&peer)
	dht.table.RemovePeer(peer)

}

//FindTargetNeighbours searches target's neighbours from given PeerID
func (dht *DHT) FindTargetNeighbours(target []byte, peer PeerID) {
	fmt.Printf("FindTargetNeighbours from %s, target: %s\n", peer, hex.EncodeToString(target))
	if peer.Equals(dht.self) {
		return
	}

	conn, err := dht.host.GetConnection(peer.ID)
	//TODO: dial remote peer???
	if conn == nil {
		fmt.Println("FindTargetNeighbours dail remote peer")
		conn, err = dht.host.Connect(peer.ID.Address, dht.network)
		if err != nil {
			return
		}

		dht.host.AddConnection(peer.ID, conn)

		*dht.network.GetConnChan() <- conn
	}

	//send find neighbours request to peer
	pmes := NewMessage(Message_FIND_NODE, hex.EncodeToString(target))
	stream, err := dht.network.CreateStream(conn, protocolDHT)
	if err != nil {
		return
	}
	stream.Write(pmes)
}

// RandomTargetID generate random peer id for query target
func RandomTargetID() []byte {
	id := make([]byte, 32)
	rand.Read(id)

	h := sha256.New()
	h.Write(id)
	return h.Sum(nil)
}

// SyncRouteTable sync route table.
func (dht *DHT) SyncRouteTable() {
	fmt.Println("timer trigger")
	dht.table.printTable()

	target := RandomTargetID()
	syncedPeers := make(map[string]bool)

	// sync with seed nodes.
	for _, addr := range dht.network.Bootstraps() {
		pid, err := ParsePeerAddr(addr)
		if err != nil {
			continue
		}

		dht.FindTargetNeighbours(target, *pid)
		syncedPeers[hex.EncodeToString(pid.Hash)] = true
	}

	// random peer selection.
	peers := dht.table.GetPeers()
	peersCount := len(peers)
	fmt.Printf("peersCount=%d\n", peersCount)
	if peersCount <= 1 {
		return
	}

	peersCountToSync := DefaultMaxPeersCountToSync
	if peersCount < peersCountToSync {
		peersCountToSync = peersCount
	}

	for i := 0; i < peersCountToSync; i++ {
		pid := peers[i]
		if syncedPeers[hex.EncodeToString(pid.Hash)] == false {
			dht.FindTargetNeighbours(target, pid)
			syncedPeers[hex.EncodeToString(pid.Hash)] = true
		}
	}
}

// saveTableToStore save peer to db
func (dht *DHT) saveTableToStore() {
	peers := dht.table.GetPeers()
	now := time.Now()
	for _, v := range peers {
		if now.Sub(v.addTime) > DefaultSeedMinTableTime {
			dht.store.Update(&v)
		}
	}
}

func (dht *DHT) loadBootstrapNodes() {
	for _, addr := range dht.network.Bootstraps() {
		peer, err := ParsePeerAddr(addr)
		if err != nil {
			continue
		}

		dht.table.Update(*peer)
	}
}

func (dht *DHT) loadTableFromDB() {
	it := dht.store.db.NewIterator(nil, nil)
	for end := false; !end; end = !it.Next() {
		peer := new(PeerID)
		err := json.Unmarshal(it.Value(), &peer)
		if err != nil {
			continue
		}
		dht.table.Update(*peer)
	}

}

// Connected is called when new connection is established
func (dht *DHT) Connected(c p2p.Conn) {
	fmt.Println("在dht中获取通知：connected")
	id := c.RemotePeer()
	peerID := CreateID(id.Address, id.PublicKey)
	dht.AddPeer(peerID)
	go dht.ping(c)
}

// Disconnected is called when the connection is closed
func (dht *DHT) Disconnected(c p2p.Conn) {
	fmt.Println("disconnect")
}

// OpenedStream is called when new stream is opened
func (dht *DHT) OpenedStream(s p2p.Stream) {

}

// ClosedStream is called when the stream is closed
func (dht *DHT) ClosedStream(s p2p.Stream) {

}

func (dht *DHT) ping(c p2p.Conn) {

	pingTicker := time.NewTicker(10 * time.Second)
	defer pingTicker.Stop()

	stop := make(chan interface{})
	dht.pingpong.mutex.Lock()
	peer := CreateID(c.RemotePeer().Address, c.RemotePeer().PublicKey).HashHex()
	if dht.pingpong.stopCh[peer] == nil {
		dht.pingpong.stopCh[peer] = stop
	} else {
		dht.pingpong.mutex.Unlock()
		return
	}
	dht.pingpong.mutex.Unlock()
	for {
		select {
		case <-stop:
			delete(dht.pingpong.stopCh, peer)
			return
		case <-pingTicker.C:
			fmt.Printf("sending ping to %s\n", c.RemotePeer().Address)
			pmes := NewMessage(Message_PING, "")
			stream, err := dht.network.CreateStream(c, protocolDHT)
			if err != nil {
				delete(dht.pingpong.stopCh, peer)
				return
			}
			err = stream.Write(pmes)
			if err != nil {
				dht.pingpong.stopCh[peer] = nil
				delete(dht.pingpong.stopCh, peer)
				return
			}
			dht.pingpong.pingpongAt[peer] = time.Now()
		}
	}

}

func (dht *DHT) checkPingPong() {
	checkTicker := time.NewTicker(30 * time.Second)
	defer checkTicker.Stop()

	for {
		select {
		case <-checkTicker.C:
			for p, t := range dht.pingpong.pingpongAt {
				if time.Now().Sub(t) > 60*time.Second {
					dht.pingpong.stopCh[p] <- new(interface{})
					hash, _ := hex.DecodeString(p)
					peer := PeerID{Hash: hash}
					dht.RemovePeer(peer)
				}
			}

		}
	}
}
