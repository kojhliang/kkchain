package dht

import (
	"time"
	"fmt"

)

type DHT struct {
	quitCh chan bool
	table *RoutingTable
	store *PeerStore

}

func NewDHT(dbPath string, self PeerID) (*DHT, error) {
	// If no node database was given, use an in-memory one
	db, err := newPeerStore(dbPath)
	if err != nil {
		return nil, err
	}

	return &DHT{
		quitCh: make(chan bool),
		table: CreateRoutingTable(self),
		store: db,
	}, nil
}

func (dht *DHT) Start() {
	fmt.Println("start sync loop.....")
	go dht.syncLoop()
}

func (dht *DHT) Stop() {
	fmt.Println("stopping sync loop.....")
	dht.quitCh <- true
}

func (dht *DHT) syncLoop()  {

	//TODO: config timer
	syncLoopTicker := time.NewTicker(DefaultSyncTableInterval)
	defer syncLoopTicker.Stop()

	saveTableToStore := time.NewTicker(DefaultSaveTableInterval)
	defer saveTableToStore.Stop()

	for  {
		select {
		case <- dht.quitCh:
			fmt.Println("stopped sync loop")
			return
		case <- syncLoopTicker.C:
			dht.SyncRouteTable()
		case <- saveTableToStore.C:
			dht.saveTableToStore()
		}
	}
}

func (dht *DHT) AddPeer(peer PeerID)  {

	dht.store.Update(&peer)
	peer.addTime = time.Now()
	dht.table.Update(peer)

}

func (dht *DHT) FindPeer(id []byte) {

}

// SyncRouteTable sync route table.
func (dht *DHT) SyncRouteTable() {
	fmt.Println("timer trigger")
	fmt.Printf("table size: %d\n", len(dht.table.GetPeers()))
	//TODO: sync peers

}

// saveTableToStore save peer to db
func (dht *DHT) saveTableToStore() {
	peers := dht.table.GetPeers()
	now := time.Now()
	for _, v := range peers {
		if  now.Sub(v.addTime) > DefaultSeedMinTableTime{
			dht.store.Update(&v)
		}
	}
}