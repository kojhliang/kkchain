package handshake

// NewMessage creates a new message object
func NewMessage(typ Message_Type) *Message {
	m := &Message{
		Type: typ,
	}

	return m
}

// BuildHandshake build message body for handshaking
func BuildHandshake(msg *Message) {
	msg.ProtocolVersion = "handshake/1.0.0"
	// TODO:
}

// TODO: add other routines
