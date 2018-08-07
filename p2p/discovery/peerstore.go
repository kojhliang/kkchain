package discovery

import (
	"sync"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/storage"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/util"
	"encoding/json"
)

type PeerStore struct {
	db	*leveldb.DB
	self	PeerID
	mutex	sync.Mutex
}

// newPeerStore creates a new peer store for storing and retrieving infos about
// known peers in the network. If no path is given, an in-memory, temporary
// database is constructed.
func newPeerStore(path string, self PeerID) (*PeerStore, error) {
	if path == "" {
		return newMemoryPeerStore(self)
	}
	return newPersistentPeerStore(path, self)
}

// newMemoryPeerStore creates a new in-memory peer store without a persistent
// backend.
func newMemoryPeerStore(self PeerID) (*PeerStore, error) {
	db, err := leveldb.Open(storage.NewMemStorage(), nil)
	if err != nil {
		return nil, err
	}
	return &PeerStore{
		db:  db,
		self: self,
	}, nil
}

// newPersistentPeerStore creates/opens a leveldb backed persistent peer store
func newPersistentPeerStore(path string, self PeerID) (*PeerStore, error) {
	opts := &opt.Options{OpenFilesCacheCapacity: 5}
	db, err := leveldb.OpenFile(path, opts)
	if _, iscorrupted := err.(*errors.ErrCorrupted); iscorrupted {
		db, err = leveldb.RecoverFile(path, nil)
	}
	if err != nil {
		return nil, err
	}

	return &PeerStore{
		db:  db,
		self: self,
	}, nil
}

func (ps *PeerStore) FindPeerByPublicKey(pubk []byte) (*PeerID) {
	id := GetIDFromPublicKey(pubk)
	return ps.FindPeerByID(id)
}

func (ps *PeerStore) FindPeerByID(id []byte) (*PeerID) {

	blob, err := ps.db.Get(id, nil)
	if err != nil {
		return nil
	}
	peer := new(PeerID)

	//TODO: protobuf unmarshal PeerID
	err = json.Unmarshal(blob, &peer)
	if err != nil {
		return  nil
	}

	return peer
}

func (ps *PeerStore)Update(id *PeerID) error {
	if id.ID == nil {
		GetIDFromPublicKey(id.PublicKey)
	}

	//TODO: protobuf marshal PeerID
	blob, err := json.Marshal(id)
	if err != nil {
		return err
	}

	return ps.db.Put(id.ID, blob, nil)
}

func (ps *PeerStore)Delete(id *PeerID) error {
	if id.ID == nil {
		id.ID = GetIDFromPublicKey(id.PublicKey)
	}

	deleter := ps.db.NewIterator(util.BytesPrefix(id.ID), nil)
	for deleter.Next() {
		if err := ps.db.Delete(deleter.Key(), nil); err != nil {
			return err
		}
	}
	return nil
}

