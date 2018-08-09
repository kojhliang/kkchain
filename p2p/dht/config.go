package dht

import (

	"time"

)

const (
	DefaultSyncTableInterval = 10 * time.Second
	DefaultSaveTableInterval = 2 * time.Minute
	DefaultSeedMinTableTime  = 6 * time.Minute
	DefaultMaxPeersCountToSync = 6
)

type Config struct {
	BucketSize           int
	RoutingTableDir      string
	BootNodes            []peeraddr
	PrivateKeyPath       string
	localhost	string
}

//DefaultConfig return default config.
func DefaultConfig() *Config {
	return &Config{
		BucketSize : BucketSize,
		RoutingTableDir : "",
		BootNodes: []string{},
		PrivateKeyPath: "",
		localhost: "/ip4/0.0.0.0/tcp/8860",
	}
}
