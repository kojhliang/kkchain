package p2p

import (
	"bytes"
	"encoding/binary"
	"testing"
)

var (
	publicKey1 = []byte("12345678901234567890123456789012")
	publicKey2 = []byte("12345678901234567890123456789011")
	publicKey3 = []byte("12345678901234567890123456789013")
	address    = "localhost:12345"

	id1 = CreateID(address, publicKey1)
	id2 = CreateID(address, publicKey2)
	id3 = CreateID(address, publicKey3)
)

func TestCreateID(t *testing.T) {
	t.Parallel()

	if !bytes.Equal(id1.PublicKey, publicKey1) {
		t.Errorf("PublicKey = %s, want %s", id1.PublicKey, publicKey1)
	}
	if id1.Address != address {
		t.Errorf("Address = %s, want %s", id1.Address, address)
	}
}

func TestString(t *testing.T) {
	t.Parallel()

	want := "ID{Address: localhost:12345, PublicKey: [49 50 51 52 53 54 55 56 57 48 49 50 51 52 53 54 55 56 57 48 49 50 51 52 53 54 55 56 57 48 49 50]}"

	if id1.String() != want {
		t.Errorf("String() = %s, want %s", id1.String(), want)
	}
}

func TestEquals(t *testing.T) {
	t.Parallel()

	if !id1.Equals(CreateID(address, publicKey1)) {
		t.Errorf("Equals() = %s, want %s", id1.PublicKeyHex(), id2.PublicKeyHex())
	}
}

func TestLess(t *testing.T) {
	t.Parallel()

	if id1.Less(id2) {
		t.Errorf("'%s'.Less(%s) should be true", id1.PublicKeyHex(), id2.PublicKeyHex())
	}

	if !id2.Less(id1) {
		t.Errorf("'%s'.Less(%s) should be false", id2.PublicKeyHex(), id1.PublicKeyHex())
	}

	if !id1.Less(id3) {
		t.Errorf("'%s'.Less(%s) should be false", id1.PublicKeyHex(), id3.PublicKeyHex())
	}
}

func TestPublicKeyHex(t *testing.T) {
	t.Parallel()

	want := "3132333435363738393031323334353637383930313233343536373839303132"
	if id1.PublicKeyHex() != want {
		t.Errorf("PublicKeyHex() = %s, want %s", id1.PublicKeyHex(), want)
	}
}

func TestXor(t *testing.T) {
	t.Parallel()

	xor := CreateID(
		address,
		[]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
	)

	result := id1.Xor(id3)

	if !xor.Equals(result) {
		t.Errorf("Xor() = %v, want %v", xor, result)
	}
}

func TestPrefixLen(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		publicKey uint32
		expected  int
	}{
		{1, 7},
		{2, 6},
		{4, 5},
		{8, 4},
		{16, 3},
		{32, 2},
		{64, 1},
	}
	for _, tt := range testCases {
		publicKey := make([]byte, 4)
		binary.LittleEndian.PutUint32(publicKey, tt.publicKey)
		id := CreateID(address, publicKey)
		if id.PrefixLen() != tt.expected {
			t.Errorf("PrefixLen() expected: %d, value: %d", tt.expected, id.PrefixLen())
		}
	}
}
