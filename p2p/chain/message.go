package chain

// NewMessage creates a new message object
func NewMessage(typ Message_Type) *Message {
	m := &Message{
		Type: typ,
	}

	return m
}

// TODO: add other routines
