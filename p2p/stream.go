package p2p

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
	"sync"

	"github.com/invin/kkchain/crypto"
	"github.com/invin/kkchain/p2p/protobuf"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

var (
	errEmptyMsg			= errors.New("received an empty message from peer")
	errInvalidMessage	= errors.New("invalid message")
	errVerifySign		= errors.New("invalid signature")
)

// Stream defines interface for reading and writing pb messages
type Stream interface {
	Write(w io.Writer, message *protobuf.Message, mux *sync.Mutex) error		// Write pb message with a mutex
	Read(conn net.Conn) (*protobuf.Message, error)								// Read pb message from connection
}

// MyStream implements the Stream interface and provides basic routiens 
// for sending and writing pb messages
type MyStream struct {
	p2p *P2P
}

// NewStream creates a new stream
func NewStream(p2p *P2P) *MyStream {
	return &MyStream{ p2p: p2p }
}

// Write writes a pb message to underlying writer
func (ms *MyStream) Write(w io.Writer, message *protobuf.Message, mux *sync.Mutex) error {
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

// Read reads a pb message from connection
func (ms *MyStream) Read(conn net.Conn) (*protobuf.Message, error) {
	var err error

	// Firstly, read the length of following message
	buffer := make([]byte, 4)
	if err = ms.ReadFull(conn, buffer); err != nil {
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
	if err = ms.ReadFull(conn, buffer); err != nil {
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
		ms.p2p.conf.signaturePolicy,
		ms.p2p.conf.hashPolicy,
		msg.Sender.PublicKey,
		SerializeMessage(msg.Sender, msg.Message.Value),
		msg.Signature,
	) {
		return nil, errVerifySign
	}

	return msg, nil
}

// ReadFull reads the full content
func (ms *MyStream) ReadFull(conn net.Conn, buffer []byte) error {
	var err error

	len := len(buffer)

	bytesRead, totalBytesRead := 0, 0
	for totalBytesRead < len {
		if bytesRead, err = conn.Read(buffer[totalBytesRead:]); err != nil {
			break
		}

		totalBytesRead += bytesRead
	}

	return err
}




