package p2p

//type Config struct {
//	NetworkConfig
//	DHTConfig
//}
//
//type NetworkConfig struct {
//	PrivateKeyPath  string
//	MaxPeers        int
//	MaxPendingPeers int
//	BootstrapNodes  []string
//	Listen          string
//
//	SignaturePolicy crypto.SignaturePolicy
//	HashPolicy      crypto.HashPolicy
//}
//
////DefaultConfig return default config.
//func DefaultConfig() *Config {
//	networkConfig := NetworkConfig{
//		PrivateKeyPath:  "",
//		MaxPeers:        50,
//		MaxPendingPeers: 50,
//		BootstrapNodes:  []string{},
//		Listen:          "/ip4/0.0.0.0/tcp/8860",
//	}
//
//	tableConfig := DHTConfig{
//		BucketSize:      16,
//		RoutingTableDir: "",
//		BootstrapNodes:  []string{},
//	}
//
//	return &Config{
//		NetworkConfig: networkConfig,
//		DHTConfig:     tableConfig,
//	}
//}
