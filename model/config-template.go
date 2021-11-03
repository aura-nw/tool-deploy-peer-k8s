package model

type ConfigNodePersistencePeerTemplate struct {
	Namespace              string
	PeerName               string
	ChainId                string
	PersistencePeerAddress string
	ExternalAddress        string
}

//config template for node-seed.yaml
type ConfigNodeSeedTemplate struct {
	Namespace string
}

//config template for node-peer.yaml
type ConfigNodePeerTemplate struct {
	Namespace        string
	PeerName         string
	SeedNodeID       string
	SeedP2PPort      string
	ChainId          string
	SeedNodeAddress  string
	ClusterIPService string
}
