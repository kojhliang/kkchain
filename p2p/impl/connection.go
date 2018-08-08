package impl

import (
	"net"
	"sync"
	"io"
	"encoding/binary"
	"bufio"

	"github.com/invin/kkchain/p2p"
	"github.com/invin/kkchain/p2p/protobuf"
	"github.com/invin/kkchain/crypto"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"github.com/golang/glog"
)

var (
	errEmptyMessage = errors.New("empty message")
)

// Connection represents a connection to remote peer
type Connection struct {
	conn net.Conn
	n	p2p.Network
	mux sync.Mutex
	remotePeer p2p.ID
}

// NewConnection creates a new connection object
func NewConnection(conn net.Conn, n p2p.Network) *Connection {
	return &Connection{
		conn: conn,
		n: n,
	}
}

// Close closes current connection
func (c *Connection) Close() error {
	return c.conn.Close()
}

// GetPeerID gets remote peer's ID
func (c *Connection) GetPeerID() p2p.ID {
	return c.remotePeer
}

// PrepareMessage prepares the message before sending
func (c *Connection) PrepareMessage(message proto.Message) (*protobuf.Message, error) {
	if message == nil {
		return nil, errEmptyMessage
	}

	raw, err := types.MarshalAny(message)
	if err != nil {
		return nil, err
	}
	// TODO: set localID
	id := protobuf.ID(c.n.ID())

	// TODO: 
	signature, err := c.n.Sign(
		SerializeMessage(&id, raw.Value),
	)

	// var signature []byte

	if err != nil {
		return nil, err
	}

	msg := &protobuf.Message{
		Message:   raw,
		Sender:    &id,
		Signature: signature,
	}
	
	return msg, nil
}

// WriteMessage writes pb message to the connection
func (c *Connection) WriteMessage(msg *protobuf.Message) error {
	return c.write(c.conn, msg, &c.mux)
}

func (c *Connection) write(w io.Writer, message *protobuf.Message, mux *sync.Mutex) error {
	// Marshal message before writing
	bytes, err := proto.Marshal(message)
	if err != nil {
		return errors.Wrap(err, "failed to marshal message")
	}
	// Allocate space for storing message
	msgLen := uint32(len(bytes))
	buffer := make([]byte, 4+msgLen)
	binary.BigEndian.PutUint32(buffer, msgLen)

	buffer = append(buffer, bytes...)
	totalSize := len(buffer)

	mux.Lock()
	defer mux.Unlock()

	// Flush buffer if any
	bw, isBuffered := w.(*bufio.Writer)
	if isBuffered && (bw.Buffered() > 0) && (bw.Available() < totalSize) {
		if err := bw.Flush(); err != nil {
			return err
		}
	}

	// Loop to write the entire buffer
	bytesWritten, totalBytesWritten := 0, 0
	
	for totalBytesWritten < totalSize {
		bytesWritten, err = w.Write(buffer[totalBytesWritten:])
		if err != nil {
			glog.Errorf("stream: failed to write entire buffer, err: %+v\n", err)
			break
		}
		totalBytesWritten += bytesWritten
	}

	if err != nil {
		return errors.Wrap(err, "stream: failed to write to socket")
	}

	return nil
}

// ReadMessage reads a pb message from connection
func (c *Connection) ReadMessage() (*protobuf.Message, error) {
	var err error

	// Firstly, read the length of following message
	buffer := make([]byte, 4)
	if err = c.readFull(buffer); err != nil {
		return nil, errors.Wrap(err, "failed to read header of message")
	}

	// Check message size
	msgLen := binary.BigEndian.Uint32(buffer)

	if msgLen == 0 {
		return nil, errEmptyMsg
	}

	// Message size is limited to 4MB
	if msgLen > (4*1024*1024) {
		return nil, errors.Errorf("message length exceeds the max size")
	}

	// Then, read the message body
	buffer = make([]byte, msgLen)
	if err = c.readFull(buffer); err != nil {
		return nil, errors.Wrap(err, "failed to read body of message")
	}

	msg := &protobuf.Message{}

	if err = proto.Unmarshal(buffer, msg); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal message")
	}

	// Sanity checking the message
	if msg.Message == nil || msg.Sender == nil || msg.Sender.PublicKey == nil || len(msg.Sender.Address) == 0 || msg.Signature == nil {
		return nil, errInvalidMessage
	}

	// Verify signature of the message
	if !crypto.Verify(
		c.n.Conf().SignaturePolicy,
		c.n.Conf().HashPolicy,
		msg.Sender.PublicKey,
		SerializeMessage(msg.Sender, msg.Message.Value),
		msg.Signature,
	) {
		return nil, errVerifySign
	}

	// FIXME: set remote peer before handshaking
	c.remotePeer = p2p.ID(*msg.Sender)

	return msg, nil
}

// ReadFull reads the full content
func (c *Connection) readFull(buffer []byte) error {
	var err error

	len := len(buffer)

	bytesRead, totalBytesRead := 0, 0
	for totalBytesRead < len {
		if bytesRead, err = c.conn.Read(buffer[totalBytesRead:]); err != nil {
			break
		}

		totalBytesRead += bytesRead
	}

	return err
}
