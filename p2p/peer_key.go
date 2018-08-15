package p2p

import (
	"encoding/base64"
	"io/ioutil"
	"os"

	"github.com/invin/kkchain/crypto"
	"github.com/invin/kkchain/crypto/ed25519"
)

// LoadNodeKeyFromFile load node priv key from file.
func LoadNodeKeyFromFile(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return UnmarshalNodeKey(string(data))
}

// LoadNodeKeyFromFileOrCreateNew load node priv key from file or create new one.
func LoadNodeKeyFromFileOrCreateNew(path string) (*crypto.KeyPair, error) {
	create := false
	if path == "" {
		path = "node.key"
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		create = true
	}

	if create {
		keys := ed25519.RandomKeyPair()
		privkey := MarshalNodeKey(keys.PrivateKey)
		ioutil.WriteFile(path, []byte(privkey), os.ModePerm)
		return keys, nil
	} else {
		privateKey, err := LoadNodeKeyFromFile(path)
		if err != nil {
			return nil, err
		}
		publicKey, _ := ed25519.New().PrivateToPublic(privateKey)

		return &crypto.KeyPair{
			PublicKey:  publicKey,
			PrivateKey: privateKey,
		}, nil
	}

}

// UnmarshalNodeKey unmarshal node key.
func UnmarshalNodeKey(data string) ([]byte, error) {
	binaryData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	return binaryData, nil
}

// MarshalNodeKey marshal node key.
func MarshalNodeKey(key []byte) string {
	return base64.StdEncoding.EncodeToString(key)
}
