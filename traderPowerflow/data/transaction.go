package data

import (
	"encoding/hex"
	"encoding/json"
	"github.com/edgexfoundry/powerflow/traderPowerflow/config"
	"golang.org/x/crypto/sha3"
	"time"
)

// transaction to be made to interface
// currently interface can be "Require", "SellOffer", "BuyOffer", "Buy"
type Tx struct {
	Id              string // Require, SellOffer, ButOffer, Buy
	CreatedAt       int64  // Require, SellOffer, ButOffer, Buy
	TransactionType string // Require, SellOffer, ButOffer, Buy

	DemandDeviceId      string // Require,
	DemandDeviceAddress string // Require,
	DemandDeviceRequire string // Require,
	DemandDeviceBuyRate string // Require,

	SupplyDeviceId       string //
	SupplyDeviceAddress  string
	SupplyDeviceSurplus  string
	SupplyDeviceSellRate string
}

func NewEnergyRequireTx(deviceId string, deviceAddress string, deviceRequire string, deviceBuyRate string) Tx {
	tx := Tx{
		Id:              "",
		CreatedAt:       time.Now().Unix(),
		TransactionType: config.ENERGYTXTYPE_REQUIRE,

		DemandDeviceId:      deviceId,
		DemandDeviceAddress: deviceAddress,
		DemandDeviceRequire: deviceRequire,
		DemandDeviceBuyRate: deviceBuyRate,
	}
	tx.Id = generateTxId(tx)
	return tx
}

func generateTxId(tx Tx) string {
	txJson, _ := json.Marshal(tx)
	int64Id := sha3.Sum512(txJson)
	id := hex.EncodeToString(int64Id[:])
	return id
}
