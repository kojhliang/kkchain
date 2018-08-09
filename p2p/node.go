package p2p

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"fmt"
	"math/big"
	"net"
	"net/url"
	"time"
)

type NodeID [64]byte

func (n NodeID) String() string {
	return fmt.Sprintf("%x", n[:])
}

type Node struct {
	IP  net.IP
	TCP uint16
	ID  NodeID
}

func (n *Node) String() string {
	u := url.URL{Scheme: "enode"}
	if n.IP == nil {
		u.Host = fmt.Sprintf("%x", n.ID[:])
	} else {
		addr := net.TCPAddr{IP: n.IP, Port: int(n.TCP)}
		u.User = url.User(fmt.Sprintf("%x", n.ID[:]))
		u.Host = addr.String()
	}
	return u.String()
}

// from PubkeyKey to NodeID.
func PubkeyID(pub *ecdsa.PublicKey) NodeID {
	var id NodeID
	pbytes := elliptic.Marshal(pub.Curve, pub.X, pub.Y)
	if len(pbytes)-1 != len(id) {
		panic(fmt.Errorf("need %d bit pubkey, got %d bits", (len(id)+1)*8, len(pbytes)))
	}
	copy(id[:], pbytes[1:])
	return id
}

// from NodeID to PubkeyKey
func (id NodeID) Pubkey() (*ecdsa.PublicKey, error) {
	p := &ecdsa.PublicKey{Curve: elliptic.P256(), X: new(big.Int), Y: new(big.Int)}
	half := len(id) / 2
	p.X.SetBytes(id[:half])
	p.Y.SetBytes(id[half:])
	if !p.Curve.IsOnCurve(p.X, p.Y) {
		return nil, errors.New("id is invalid secp256k1 curve point")
	}
	return p, nil
}

type Peer struct {
	rw      *conn
	created time.Time
	closed  chan struct{}
	disc    chan DisconnectReason
}

func (p *Peer) ID() NodeID {
	return p.rw.id
}

func (p *Peer) RemoteAddr() net.Addr {
	return p.rw.fd.RemoteAddr()
}

func (p *Peer) LocalAddr() net.Addr {
	return p.rw.fd.LocalAddr()
}

func (p *Peer) Disconnect(reason DisconnectReason) {
	select {
	case p.disc <- reason:
	case <-p.closed:
	}
}

func (p *Peer) String() string {
	return fmt.Sprintf("Peer %x %v", p.rw.id[:8], p.RemoteAddr())
}

func (p *Peer) Inbound() bool {
	return p.rw.flags&inboundConn != 0
}

func newPeer(conn *conn) *Peer {
	p := &Peer{
		rw:      conn,
		created: time.Now(),
		disc:    make(chan DisconnectReason),
		closed:  make(chan struct{}),
	}
	return p
}

type PeerInfo struct {
	ID            string
	LocalAddress  string
	RemoteAddress string
	Inbound       bool
	Static        bool
}

func (p *Peer) Info() *PeerInfo {
	info := &PeerInfo{
		ID: p.ID().String(),
	}
	info.LocalAddress = p.LocalAddr().String()
	info.RemoteAddress = p.RemoteAddr().String()
	info.Inbound = p.rw.isThisFlag(inboundConn)
	info.Static = p.rw.isThisFlag(staticDialedConn)
	return info
}
