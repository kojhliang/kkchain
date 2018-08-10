package dht

import (
	"encoding/json"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/storage"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type PeerStore struct {
	db          *leveldb.DB
	enableBatch bool
	opts        []operation
	mutex       sync.Mutex
}

type opType int

const (
	opTypePut opType = 1
	opTypeDel opType = 2
)

type operation struct {
	opType
	key   []byte
	value []byte
}

// newPeerStore creates a new peer store for storing and retrieving infos about
// known peers in the network. If no path is given, an in-memory, temporary
// database is constructed.
func newPeerStore(path string) (*PeerStore, error) {
	if path == "" {
		return newMemoryPeerStore()
	}
	return newPersistentPeerStore(path)
}

// newMemoryPeerStore creates a new in-memory peer store without a persistent
// backend.
func newMemoryPeerStore() (*PeerStore, error) {
	db, err := leveldb.Open(storage.NewMemStorage(), nil)
	if err != nil {
		return nil, err
	}
	return &PeerStore{
		db:          db,
		enableBatch: false,
	}, nil
}

// newPersistentPeerStore creates/opens a leveldb backed persistent peer store
func newPersistentPeerStore(path string) (*PeerStore, error) {
	opts := &opt.Options{OpenFilesCacheCapacity: 5}
	db, err := leveldb.OpenFile(path, opts)
	if _, iscorrupted := err.(*errors.ErrCorrupted); iscorrupted {
		db, err = leveldb.RecoverFile(path, nil)
	}
	if err != nil {
		return nil, err
	}

	return &PeerStore{
		db:          db,
		enableBatch: false,
	}, nil
}

// FindPeerByPublicKey
func (ps *PeerStore) FindPeerByPublicKey(pubk []byte) *PeerID {
	id := GetIDFromPublicKey(pubk)
	return ps.FindPeerByID(id)
}

// FindPeerByID
func (ps *PeerStore) FindPeerByID(id []byte) *PeerID {

	blob, err := ps.db.Get(id, nil)
	if err != nil {
		return nil
	}
	peer := new(PeerID)

	err = json.Unmarshal(blob, &peer)
	if err != nil {
		return nil
	}

	return peer
}

// Update put key-value to levelDB
func (ps *PeerStore) Update(id *PeerID) error {

	blob, err := json.Marshal(id)
	if err != nil {
		return err
	}

	if ps.enableBatch {
		ps.mutex.Lock()
		defer ps.mutex.Unlock()

		ps.opts = append(ps.opts, operation{opTypePut, id.Hash, blob})
		return nil
	}

	return ps.db.Put(id.Hash, blob, nil)
}

// Delete delete the specify key
func (ps *PeerStore) Delete(id *PeerID) error {

	if ps.enableBatch {
		ps.mutex.Lock()
		defer ps.mutex.Unlock()

		ps.opts = append(ps.opts, operation{opTypeDel, id.Hash, nil})
		return nil
	}

	deleter := ps.db.NewIterator(util.BytesPrefix(id.Hash), nil)
	for deleter.Next() {
		if err := ps.db.Delete(deleter.Key(), nil); err != nil {
			return err
		}
	}
	return nil
}

// Close levelDB
func (ps *PeerStore) Close() error {
	return ps.db.Close()
}

// Flush write and flush pending batch write.
func (ps *PeerStore) Flush() error {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	if !ps.enableBatch {
		return nil
	}

	batch := new(leveldb.Batch)
	for _, opt := range ps.opts {
		switch opt.opType {
		case opTypePut:
			batch.Put(opt.key, opt.value)
		case opTypeDel:
			batch.Delete(opt.key)
		}
	}
	ps.opts = nil

	return ps.db.Write(batch, nil)
}

// EnableBatch enable batch write.
func (ps *PeerStore) EnableBatch() {
	ps.enableBatch = true
}

// DisableBatch disable batch write.
func (ps *PeerStore) DisableBatch() {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	ps.opts = nil
	ps.enableBatch = false
}
