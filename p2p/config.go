package p2p

import (
	"time"

	"github.com/invin/kkchain/crypto"
)

const (
	DefaultSyncTableInterval   = 10 * time.Second
	DefaultSaveTableInterval   = 2 * time.Minute
	DefaultSeedMinTableTime    = 6 * time.Minute
	DefaultMaxPeersCountToSync = 6
)

type Config struct {
	NetworkConfig
	TableConfig
}

type NetworkConfig struct {
	PrivateKeyPath  string
	MaxPeers        int
	MaxPendingPeers int
	BootstrapNodes  []string
	Listen          string

	SignaturePolicy crypto.SignaturePolicy
	HashPolicy      crypto.HashPolicy
}

type TableConfig struct {
	BucketSize      int
	RoutingTableDir string
}

//DefaultConfig return default config.
func DefaultConfig() *Config {
	networkConfig := NetworkConfig{
		PrivateKeyPath:  "",
		MaxPeers:        50,
		MaxPendingPeers: 50,
		BootstrapNodes:  []string{},
		Listen:          "/ip4/0.0.0.0/tcp/8860",
	}

	tableConfig := TableConfig{
		BucketSize:      16,
		RoutingTableDir: "",
	}

	return &Config{
		NetworkConfig: networkConfig,
		TableConfig:   tableConfig,
	}
}
