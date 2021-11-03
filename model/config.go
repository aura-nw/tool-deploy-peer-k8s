package model

// config template from app.env
type Config struct {
	Namespace                   string `mapstructure:"NAMESPACE"`
	PathTemplate                string `mapstructure:"PATH_TEMPLATE"`
	TemplateNodeSeed            string `mapstructure:"TEMPLATE_NODE_SEED"`
	TemplateNodePeer            string `mapstructure:"TEMPLATE_NODE_PEER"`
	TemplateNodePeerSVC         string `mapstructure:"TEMPLATE_NODE_PEER_SVC"`
	TemplateNodePersistencePeer string `mapstructure:"TEMPLATE_NODE_PERSISTENCE_PEER"`
	PathOutput                  string `mapstructure:"PATH_OUTPUT"`
	SeedSvc                     string `mapstructure:"SEED_SVC"`
	SeedSvcPortP2P              string `mapstructure:"SEED_SVC_PORT_P2P"`
	RunLocal                    bool   `mapstructure:"RUN_LOCAL"`
	ChainId                     string `mapstructure:"CHAIN_ID"`

	SeedNodeAddress string `mapstructure:"SEED_NODE_ADDRESS"`
	SeedHomePath    string `mapstructure:"SEED_HOME_PATH"`
	SeedBinaryPath  string `mapstructure:"SEED_BINARY_PATH"`

	PeerBinaryPath string `mapstructure:"PEER_BINARY_PATH"`
	PeerHomePath   string `mapstructure:"PEER_HOME_PATH"`
	PeerPvcPath    string `mapstructure:"PEER_PVC_PATH"`
}
