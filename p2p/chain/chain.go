package chain

import (
	"context"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"
	"github.com/invin/kkchain/p2p"
)

const (
	protocolChain = "/kkchain/p2p/chain/1.0.0"
)

// Chain implements protocol for chain related messages
type Chain struct {
	// self
	host p2p.Host
}

// NewChain creates a new Chain object
func NewChain(host p2p.Host) *Chain {
	c := &Chain{
		host: host,
	}

	if err := host.SetStreamHandler(protocolChain, c.handleNewStream); err != nil {
		panic(err)
	}

	return c
}

// handleNewStream handles messages within the stream
func (c *Chain) handleNewStream(s p2p.Stream, msg proto.Message) {
	// check message type
	switch message := msg.(type) {
	case *Message:
		c.handleMessage(s, message)
	default:
		s.Reset()
		glog.Errorf("unexpected message: %v", msg)
	}
}

// handleMessage handles messsage
func (c *Chain) handleMessage(s p2p.Stream, msg *Message) {
	// get handler
	handler := c.handlerForMsgType(msg.GetType())
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
}
