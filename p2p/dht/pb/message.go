package pb

// NewMessage creates a new message object
func NewMessage(typ Message_MessageType) *Message {
	m := &Message{
		Type: typ,
	}

	return m
}

