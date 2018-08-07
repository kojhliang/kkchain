package p2p

import (
	"github.com/invin/kkchain/p2p/discovery"
	"time"
)

const (
	DefaultSyncTableInterval = 30 * time.Second
	DefaultSaveTableInterval = 2 * time.Minute
	DefaultSeedMinTableTime  = 6 * time.Minute
)

type Config struct {
	BucketSize           int
	RoutingTableDir      string
	BootNodes            []string
	PrivateKeyPath       string
}

//DefaultConfig return default config.
func DefaultConfig() *Config {
	return &Config{
		BucketSize : discovery.BucketSize,
		RoutingTableDir : "",
		BootNodes: []string{},
		PrivateKeyPath: "",
	}
}
