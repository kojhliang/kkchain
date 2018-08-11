package impl

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/invin/kkchain/p2p"
)

type Dialer interface {
	Dial(*Node) (net.Conn, error)
}

type TCPDialer struct {
	*net.Dialer
}

func (t TCPDialer) Dial(dest *Node) (net.Conn, error) {
	port, err := strconv.ParseInt(dest.TCPPort, 10, 10)
	if err != nil {
		return nil, err
	}
	addr := &net.TCPAddr{IP: net.ParseIP(dest.IP), Port: int(port)}
	return t.Dialer.Dial("tcp", addr.String())
}

// dialstate schedules dials
type dialstate struct {
	dialing   map[p2p.ID]*dialTask
	task      map[p2p.ID]*dialTask
	bootnodes []*Node
}

type dialTask struct {
	flag connFlag
	dest *Node
}

func newDialState(bootnodes []*Node) *dialstate {
	s := &dialstate{
		task:      make(map[p2p.ID]*dialTask),
		bootnodes: make([]*Node, len(bootnodes)),
	}
	copy(s.bootnodes, bootnodes)
	for _, node := range bootnodes {
		dialtask := &dialTask{
			outboundConn,
			node,
		}
		s.task[node.ID] = dialtask
	}
	return s
}

var (
	errSelf             = errors.New("is self")
	errAlreadyDialing   = errors.New("already dialing")
	errAlreadyConnected = errors.New("already connected")
	errRecentlyDialed   = errors.New("recently dialed")
	errNotWhitelisted   = errors.New("not contained in netrestrict whitelist")
)

func (s *dialstate) checkDial(n *Node) error {
	_, dialing := s.dialing[n.ID]
	switch {
	case dialing:
		return errAlreadyDialing

		// TODO: other case
	}
	return nil
}

func (t *dialTask) Do(n *Network) {
	if t.dest.IP == "" {
		if !t.resolve(n) {
			return
		}
	}
	err := t.dial(n, t.dest)
	if err != nil {
		log.Error("failed to dial:", err)
	}
}

func (t *dialTask) resolve(n *Network) bool {
	var resolved *Node
	peer := n.peers[t.dest.ID]
	if peer != nil {
		addrStr := peer.RemoteAddr().String()
		addrArr := strings.Split(addrStr, ":")

		resolved = &Node{
			addrArr[0],
			addrArr[1],
			peer.ID,
		}
	}
	t.dest = resolved
	return true
}

func (t *dialTask) dial(n *Network, dest *Node) error {
	return n.host.Connect(dest.Addr())
}

func (t *dialTask) String() string {
	return fmt.Sprintf("%v %x %v:%d", t.flag, t.dest.ID.PublicKey[:8], t.dest.IP, t.dest.TCPPort)
}
