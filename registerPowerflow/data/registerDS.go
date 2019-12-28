package data

import (
	"fmt"
	"sync"
)

// Stores all states of trader
type registerDS struct {
	// ip and port of register service
	IpPort string
	// peers collection
	registerPeers   map[string]Peer
	blockchainPeers map[string]Peer

	// locks for variables defined above
	registerPeersMux   sync.Mutex
	blockchainPeersMux sync.Mutex
}

type RegisterDSCopy struct {
	// ip and port of register service
	IpPort string
	// peers collection
	RegisterPeers   map[string]Peer
	BlockchainPeers map[string]Peer
}

// Single instance of traderDS
var InstanceRegisterDS *registerDS

// Allowing "run only once" behaviour
var onceDataStore sync.Once

// Method to initialize Single instance of traderDS
func SetDataStore(ipPort string) *registerDS {
	onceDataStore.Do(func() {
		InstanceRegisterDS = &registerDS{
			IpPort:          ipPort,
			registerPeers:   make(map[string]Peer),
			blockchainPeers: make(map[string]Peer),
		}
	})
	fmt.Println(InstanceRegisterDS)
	return InstanceRegisterDS
}

func (rds *registerDS) GetRegisterDSCopy() RegisterDSCopy {
	rdsc := RegisterDSCopy{
		IpPort:          InstanceRegisterDS.IpPort,
		RegisterPeers:   InstanceRegisterDS.GetRegisterPeers(),
		BlockchainPeers: InstanceRegisterDS.GetBlockchainPeers(),
	}
	return rdsc
}

func (rds *registerDS) AddToRegisterPeers(p Peer) {
	rds.registerPeersMux.Lock()
	defer rds.registerPeersMux.Unlock()

	rds.registerPeers[p.IpPort] = p
}

func (rds *registerDS) GetRegisterPeers() map[string]Peer {
	rds.registerPeersMux.Lock()
	defer rds.registerPeersMux.Unlock()

	copyOfPeers := make(map[string]Peer)
	for key, value := range rds.registerPeers {
		copyOfPeers[key] = value
	}
	return copyOfPeers
}

func (rds *registerDS) AddToBlockchainPeers(p Peer) {
	rds.blockchainPeersMux.Lock()
	defer rds.blockchainPeersMux.Unlock()

	rds.blockchainPeers[p.IpPort] = p
}

func (rds *registerDS) GetBlockchainPeers() map[string]Peer {
	rds.registerPeersMux.Lock()
	defer rds.registerPeersMux.Unlock()

	copyOfPeers := make(map[string]Peer)
	for key, value := range rds.blockchainPeers {
		copyOfPeers[key] = value
	}
	return copyOfPeers
}
