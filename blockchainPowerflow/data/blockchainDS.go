package data

import (
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/edgexfoundry/powerflow/blockchainPowerflow/data/peerList"
	"github.com/edgexfoundry/powerflow/blockchainPowerflow/data/syncBlockchain"
	"github.com/edgexfoundry/powerflow/traderPowerflow/config"

	//"github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/edgexfoundry/powerflow/commonPowerFlow/identity"
	//"github.com/edgexfoundry/powerflow/traderPowerflow/config"
	"log"
	"os"
	"sync"
)

//package data
//
//import (
//"fmt"
//"github.com/edgexfoundry/go-mod-core-contracts/models"
//"github.com/edgexfoundry/powerflow/commonPowerFlow/identity"
//"github.com/edgexfoundry/powerflow/traderPowerflow/config"
//"log"
//"os"
//"sync"
//)

// Stores all states of blockchain
type blockchainDS struct {

	// ip:port of blockchain service and registerer
	IpPort, RegisterIpPort string
	//blockchain Identity - used as Id for blockchain
	BlockchainId identity.Identity
	//Peers hold peers detail and manage peers, is safe for distributed use
	Peers peerList.PeerList
	//Transaction Pool
	TxPool tkn.TransactionPool

	// SBC is safe for distributed use
	SBC syncBlockchain.SyncBlockchain

	// locks for variables defined above
	traderIdMux           sync.Mutex
	edgexDevicesMux       sync.Mutex
	edgexDevicesEventsMux sync.Mutex
}

// Single instance of traderDS
var InstanceTraderDS *traderDS

// Allowing "run only once" behaviour
var onceDataStore sync.Once

// Method to initialize Single instance of traderDS
func SetDataStore(label string, ipPort string) *traderDS {
	onceDataStore.Do(func() {
		// setting up traderId for InstanceTraderDS
		var traderId identity.Identity
		if label == "Anon" { // todo : make logic streamlined
			if _, err := os.Stat(config.PRIVATEKEYFILE); err == nil { // here, checking if old identity has to be loaded
				log.Println("Loading old key ...")
				if _, err := os.Stat(config.IDENTITYFILE); err == nil {
					i := identity.LoadIdentityFromFile(config.IDENTITYFILE) // here, loading whole Identity just to get label of stored id
					label = i.Label
				}
				traderId = identity.OldIdentity(label, ipPort, config.PRIVATEKEYFILE)
			}
		} else {
			traderId = identity.NewIdentity(label, ipPort, config.PRIVATEKEYFILE, config.IDENTITYFILE)
		}
		// setting up InstanceTraderDS
		InstanceTraderDS = &traderDS{
			TraderId:           traderId,
			edgexDevices:       make(map[string]models.Device),
			edgexDevicesEvents: make(map[string]map[string]models.Reading),
		}
	})
	fmt.Println(InstanceTraderDS)
	return InstanceTraderDS
}

//func GetDataStore() *traderDS {
//	return InstanceTraderDS
//}

// Add device to variable edgexDevices
func (tds *traderDS) AddToEdgeXDevices(d models.Device) {
	tds.edgexDevicesMux.Lock()
	defer tds.edgexDevicesMux.Unlock()
	tds.edgexDevices[d.Id] = d
}

// Getting a copy of variable edgexDevices
func (tds *traderDS) GetEdgeXDevices() map[string]models.Device {
	tds.edgexDevicesMux.Lock()
	defer tds.edgexDevicesMux.Unlock()
	copyOfedgexDevices := tds.edgexDevices
	return copyOfedgexDevices
}

// Adding device-event-reading from device-event to variable edgexDevicesEvents
func (tds *traderDS) AddToEdgeXDevicesEvents(e models.Event) {
	tds.edgexDevicesEventsMux.Lock()
	defer tds.edgexDevicesEventsMux.Unlock()
	if tds.edgexDevicesEvents[e.Readings[0].Device] == nil {
		tds.edgexDevicesEvents[e.Readings[0].Device] = make(map[string]models.Reading)
	}

	deviceReadings := tds.edgexDevicesEvents[e.Readings[0].Device]
	deviceReadings[e.Readings[0].Name] = e.Readings[0]

	tds.edgexDevicesEvents[e.Readings[0].Device] = deviceReadings
}

// Getting a copy of variable edgexDevicesEvents
func (tds *traderDS) GetEdgeXDevicesEvents() map[string]map[string]models.Reading {
	tds.edgexDevicesEventsMux.Lock()
	defer tds.edgexDevicesEventsMux.Unlock()
	copyOfEdgeXDevicesEvents := tds.edgexDevicesEvents
	return copyOfEdgeXDevicesEvents
}

func (tds *traderDS) GetDeviceFromEdgeXDevices(deviceName string) models.Device {
	tds.edgexDevicesMux.Lock()
	defer tds.edgexDevicesMux.Unlock()
	return tds.edgexDevices[deviceName]
}
