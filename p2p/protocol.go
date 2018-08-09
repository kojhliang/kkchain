package p2p

type MsgType int

const (
	msg_dht MsgType = 1
	msg_handshake
	msg_chain
)

var MsgTypeString = map[MsgType]string{
	msg_dht:       "dht",
	msg_handshake: "hanshake",
	msg_chain:     "chain",
}
