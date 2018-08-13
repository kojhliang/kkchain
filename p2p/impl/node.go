package impl

import (
	"fmt"
	"net"
	"net/url"
	"time"

	"encoding/hex"

	"strconv"

	"github.com/invin/kkchain/p2p"
)

type Node struct {
	IP      string
	TCPPort string
	ID      p2p.ID
}

func (n *Node) String() string {
	u := url.URL{Scheme: "enode"}
	if n.IP == "" {
		u.Host = fmt.Sprintf("%x", n.ID.PublicKey[:])
	} else {
		port, err := strconv.ParseInt(n.TCPPort, 10, 10)
		if err != nil {
			log.Error("failed to parse port:", err)
			port = 0
		}
		addr := net.TCPAddr{IP: net.ParseIP(n.IP), Port: int(port)}
		u.User = url.User(hex.EncodeToString(n.ID.PublicKey))
		u.Host = addr.String()
	}
	return u.String()
}

func (n *Node) Addr() string {
	return n.IP + ":" + n.TCPPort
}

type Peer struct {
	conn     *Connection
	created  time.Time
	lastPing time.Time
	lastPong time.Time
	closed   chan struct{}
	disc     chan DisconnectReason
	ID       p2p.ID
}

func (p *Peer) RemoteAddr() net.Addr {
	return p.conn.RemoteAddr()
}

func (p *Peer) LocalAddr() net.Addr {
	return p.conn.LocalAddr()
}

func (p *Peer) Disconnect(reason DisconnectReason) {
	select {
	case p.disc <- reason:
	case <-p.closed:
		p.conn.Close()
	}
}

func (p *Peer) String() string {
	return fmt.Sprintf("Peer %s %v", hex.EncodeToString(p.ID.PublicKey), p.RemoteAddr())
}

func newPeer(conn *Connection) *Peer {
	p := &Peer{
		conn:    conn,
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
}

func (p *Peer) Info() *PeerInfo {
	info := &PeerInfo{
		ID: hex.EncodeToString(p.ID.PublicKey),
	}
	info.LocalAddress = p.LocalAddr().String()
	info.RemoteAddress = p.RemoteAddr().String()
	return info
}
