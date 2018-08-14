package dht

import (
	"container/list"
	"fmt"
	"sort"
	"sync"
)

// BucketSize defines the NodeID, Key, and routing table data structures.
const BucketSize = 16

// RoutingTable contains one bucket list for lookups.
type RoutingTable struct {
	// Current node's ID.
	self PeerID

	buckets []*Bucket
}

// Bucket holds a list of contacts of this node.
type Bucket struct {
	*list.List
	mutex *sync.RWMutex
}

// NewBucket is a Factory method of Bucket, contains an empty list.
func NewBucket() *Bucket {
	return &Bucket{
		List:  list.New(),
		mutex: &sync.RWMutex{},
	}
}

// CreateRoutingTable is a Factory method of RoutingTable containing empty buckets.
func CreateRoutingTable(pid PeerID) *RoutingTable {
	table := &RoutingTable{
		self:    pid,
		buckets: make([]*Bucket, len(pid.Hash)*8),
	}
	for i := 0; i < len(pid.Hash)*8; i++ {
		table.buckets[i] = NewBucket()
	}

	table.Update(pid)

	return table
}

// Self returns the ID of the node hosting the current routing table instance.
func (t *RoutingTable) Self() PeerID {
	return t.self
}

// Update moves a peer to the front of a bucket in the routing table.
func (t *RoutingTable) Update(target PeerID) {
	if len(t.self.PublicKey) != len(target.PublicKey) {
		return
	}

	bucketID := target.Xor(t.self).PrefixLen()
	bucket := t.Bucket(bucketID)

	var element *list.Element

	// Find current node in bucket.
	bucket.mutex.Lock()

	for e := bucket.Front(); e != nil; e = e.Next() {
		if e.Value.(PeerID).Equals(target) {
			element = e
			break
		}
	}

	if element == nil {
		// Populate bucket if its not full.
		if bucket.Len() <= BucketSize {
			bucket.PushFront(target)
		}
	} else {
		bucket.MoveToFront(element)
	}

	bucket.mutex.Unlock()
}

// GetPeers returns a randomly-ordered, unique list of all peers within the routing network (excluding itself).
func (t *RoutingTable) GetPeers() (peers []PeerID) {
	visited := make(map[string]struct{})
	visited[t.self.PublicKeyHex()] = struct{}{}

	for _, bucket := range t.buckets {
		bucket.mutex.RLock()

		for e := bucket.Front(); e != nil; e = e.Next() {
			id := e.Value.(PeerID)
			if _, seen := visited[id.PublicKeyHex()]; !seen {
				peers = append(peers, id)
				visited[id.PublicKeyHex()] = struct{}{}
			}
		}

		bucket.mutex.RUnlock()
	}

	return
}

// GetPeerAddresses returns a unique list of all peer addresses within the routing network.
func (t *RoutingTable) GetPeerAddresses() (peers []string) {
	visited := make(map[string]struct{})
	visited[t.self.PublicKeyHex()] = struct{}{}

	for _, bucket := range t.buckets {
		bucket.mutex.RLock()

		for e := bucket.Front(); e != nil; e = e.Next() {
			id := e.Value.(PeerID)
			if _, seen := visited[id.PublicKeyHex()]; !seen {
				peers = append(peers, id.Address)
				visited[id.PublicKeyHex()] = struct{}{}
			}
		}

		bucket.mutex.RUnlock()
	}

	return
}

// RemovePeer removes a peer from the routing table with O(bucket_size) time complexity.
func (t *RoutingTable) RemovePeer(target PeerID) bool {
	bucketID := target.Xor(t.self).PrefixLen()
	bucket := t.Bucket(bucketID)

	bucket.mutex.Lock()

	for e := bucket.Front(); e != nil; e = e.Next() {
		if e.Value.(PeerID).Equals(target) {
			bucket.Remove(e)

			bucket.mutex.Unlock()
			return true
		}
	}

	bucket.mutex.Unlock()

	return false
}

// PeerExists checks if a peer exists in the routing table with O(bucket_size) time complexity.
func (t *RoutingTable) PeerExists(target PeerID) bool {
	bucketID := target.Xor(t.self).PrefixLen()
	bucket := t.Bucket(bucketID)

	bucket.mutex.Lock()

	defer bucket.mutex.Unlock()

	for e := bucket.Front(); e != nil; e = e.Next() {
		if e.Value.(PeerID).Equals(target) {
			return true
		}
	}

	return false
}

// FindClosestPeers returns a list of k(count) peers with smallest XOR distance.
func (t *RoutingTable) FindClosestPeers(target PeerID, count int) (peers []PeerID) {
	if len(t.self.PublicKey) != len(target.PublicKey) {
		return []PeerID{}
	}

	bucketID := target.Xor(t.self).PrefixLen()
	bucket := t.Bucket(bucketID)

	bucket.mutex.RLock()

	for e := bucket.Front(); e != nil; e = e.Next() {
		peers = append(peers, e.Value.(PeerID))
	}

	bucket.mutex.RUnlock()

	for i := 1; len(peers) < count && (bucketID-i >= 0 || bucketID+i < len(t.self.PublicKey)*8); i++ {
		if bucketID-i >= 0 {
			other := t.Bucket(bucketID - i)
			other.mutex.RLock()
			for e := other.Front(); e != nil; e = e.Next() {
				peers = append(peers, e.Value.(PeerID))
			}
			other.mutex.RUnlock()
		}

		if bucketID+i < len(t.self.PublicKey)*8 {
			other := t.Bucket(bucketID + i)
			other.mutex.RLock()
			for e := other.Front(); e != nil; e = e.Next() {
				peers = append(peers, e.Value.(PeerID))
			}
			other.mutex.RUnlock()
		}
	}

	// Sort peers by XOR distance.
	sort.Slice(peers, func(i, j int) bool {
		left := peers[i].Xor(target)
		right := peers[j].Xor(target)
		return left.Less(right)
	})

	if len(peers) > count {
		peers = peers[:count]
	}

	return peers
}

// Bucket returns a specific Bucket by ID.
func (t *RoutingTable) Bucket(id int) *Bucket {
	if id >= 0 && id < len(t.buckets) {
		return t.buckets[id]
	}
	return nil
}

func (t *RoutingTable) printTable() {
	fmt.Printf("self = %s\n", t.self)
	for idx, bucket := range t.buckets {
		bucket.mutex.RLock()

		for e := bucket.Front(); e != nil; e = e.Next() {
			id := e.Value.(PeerID)
			fmt.Printf("distance: %d, %s\n", idx, id)
		}

		bucket.mutex.RUnlock()
	}
}
