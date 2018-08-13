package impl

import (
	"bytes"
	"testing"

	"github.com/invin/kkchain/crypto/ed25519"
	"github.com/invin/kkchain/p2p"
	"github.com/invin/kkchain/p2p/protobuf"
)

func TestSerialization(t *testing.T) {
	// Create keypair
	keys := ed25519.RandomKeyPair()
	address := "localhost:12345"
	id := p2p.CreateID(address, keys.PublicKey)

	pid := protobuf.ID(id)
	message := "message1"

	out := SerializeMessage(&pid, []byte(message))
	t.Log("result =", out)

	rid, rmsg := DeserializeMessage(out)

	if rid.Address != address {
		t.Errorf("expected %v, actual %v", address, rid.Address)
	}

	if bytes.Compare(rid.PublicKey, keys.PublicKey) != 0 {
		t.Errorf("expected %v, actual %v", keys.PublicKey, rid.PublicKey)
	}

	if string(rmsg) != message {
		t.Errorf("expected %v, actual %v", message, rmsg)
	}

	t.Log(rmsg)
}
