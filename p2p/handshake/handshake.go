package handshake

import (
	"context"
	"fmt"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"
	"github.com/invin/kkchain/p2p"
)

const (
	protocolHandshake = "/kkchain/p2p/handshake/1.0.0"
)

// Handshake implements protocol for handshaking
type Handshake struct {
	// self
	host p2p.Host
}

// NewHandshake creates a new Handshake object with the given peer as as the 'local' host
func NewHandshake(host p2p.Host) *Handshake {
	hs := &Handshake{
		host: host,
	}

	if err := host.SetStreamHandler(protocolHandshake, hs.handleNewStream); err != nil {
		panic(err)
	}

	host.Register(hs)

	return hs
}

// handleNewStream handles messages within the stream
func (hs *Handshake) handleNewStream(s p2p.Stream, msg proto.Message) {
	// check message type
	switch message := msg.(type) {
	case *Message:
		hs.handleMessage(s, message)
	default:
		s.Reset()
		glog.Errorf("unexpected message: %v", msg)
	}
}

// handleMessage handles messsage
func (hs *Handshake) handleMessage(s p2p.Stream, msg *Message) {
	// get handler
	handler := hs.handlerForMsgType(msg.GetType())
	if handler == nil {
		s.Reset()
		glog.Errorf("unknown message type: %v", msg.GetType())
		return
	}

	// dispatch handler
	// TODO: get context and peer id
	ctx := context.Background()
	pid := s.RemotePeer()

	// successfully recv conn, add it
	hs.host.AddConnection(pid, s.Conn())

	rpmes, err := handler(ctx, pid, msg)

	// if nil response, return it before serializing
	if rpmes == nil {
		glog.Warning("got back nil response from request")
		return
	}

	// hello error resp will return err,so reset conn and remove it
	if err != nil {
		hs.host.RemoveConnection(pid)
		s.Reset()
	}

	// send out response msg
	if err = s.Write(rpmes); err != nil {
		glog.Errorf("send response error: %s", err)
		return
	}
}

// Connected is called when new connection is established
func (hs *Handshake) Connected(c p2p.Conn) {
	fmt.Println("connected")
}

// Disconnected is called when the connection is closed
func (hs *Handshake) Disconnected(c p2p.Conn) {
	fmt.Println("disconnect")
}

// OpenedStream is called when new stream is opened
func (hs *Handshake) OpenedStream(s p2p.Stream) {

}

// ClosedStream is called when the stream is closed
func (hs *Handshake) ClosedStream(s p2p.Stream) {

}
