package dht

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

//type DHT struct {
//	quitCh chan bool
//	table *RoutingTable
//	store *PeerStore
//	seedPeers []PeerID
//
//}
//
//func NewDHT(dbPath string, self PeerID) (*DHT, error) {
//	// If no node database was given, use an in-memory one
//	db, err := newPeerStore(dbPath)
//	if err != nil {
//		return nil, err
//	}
//
//	return &DHT{
//		quitCh: make(chan bool),
//		table: CreateRoutingTable(self),
//		store: db,
//	}, nil
//}

func (dht *DHT) Start() {
	fmt.Println("start sync loop.....")
	go dht.syncLoop()
	go dht.waitReceive()
}

func (dht *DHT) Stop() {
	fmt.Println("stopping sync loop.....")
	dht.quitCh <- true
}

func (dht *DHT) syncLoop() {

	//TODO: config timer
	syncLoopTicker := time.NewTicker(DefaultSyncTableInterval)
	defer syncLoopTicker.Stop()

	saveTableToStore := time.NewTicker(DefaultSaveTableInterval)
	defer saveTableToStore.Stop()

	for {
		select {
		case <-dht.quitCh:
			fmt.Println("stopped sync loop")
			return
		case <-syncLoopTicker.C:
			dht.SyncRouteTable()
		case <-saveTableToStore.C:
			dht.saveTableToStore()
		}
	}
}

func (dht *DHT) waitReceive() {
	for {
		select {
		case msg := <-dht.recvCh:
			fmt.Println(msg)
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

	if peer.Equals(dht.self) {
		return
	}

	conn, err := dht.host.GetConnection(peer.ID)
	if err != nil {
		return
	}

	//TODO: dial remote peer???
	if conn == nil {
		conn, err = dht.host.Connect(peer.ID.Address)
		if err != nil {
			return
		}
	}

	//send find neighbours request to peer
	pmes := NewMessage(Message_FIND_NODE, hex.EncodeToString(target))
	dht.host.SendMsg(conn, protocolDHT, pmes)
}

// RandomTargetID generate random peer id for query target
func RandomTargetID() []byte {
	id := make([]byte, 32)
	rand.Read(id)
	return sha256.New().Sum(id)
}

// SyncRouteTable sync route table.
func (dht *DHT) SyncRouteTable() {
	fmt.Println("timer trigger")
	fmt.Printf("table size: %d\n", len(dht.table.GetPeers()))

	target := RandomTargetID()
	syncedPeers := make(map[string]bool)

	// sync with seed nodes.
	for _, addr := range dht.config.BootstrapNodes {
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
