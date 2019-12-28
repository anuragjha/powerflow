package rulesEngine

import (
	//"github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/edgexfoundry/powerflow/traderPowerflow/data"
	"strings"
)

// checking Consume devices and store their requirement(demand) details
func ProcessForDemand() {

	devicesDemandDetails := make(map[string]data.DemandDetails)    // key = deviceId, value = DemandDetails
	devicesEvents := data.InstanceTraderDS.GetEdgeXDevicesEvents() // map[string]map[string]models.Reading
	for deviceName, deviceEvents := range devicesEvents {
		if strings.Contains(deviceName, "Consume") {
			device := data.InstanceTraderDS.GetDeviceFromEdgeXDevices(deviceName)
			demand := data.DemandDetails{
				DeviceName:    deviceName,
				DeviceAddress: device.Service.Addressable.GetBaseURL(),
				DeviceRequire: deviceEvents["require"].Value,
				DeviceBuyRate: deviceEvents["buyRate"].Value,
			}
			devicesDemandDetails[deviceName] = demand
		}
	}

	for _, deviceDemandDetails := range devicesDemandDetails {
		tx := data.NewEnergyRequireTx(deviceDemandDetails.DeviceName,
			deviceDemandDetails.DeviceAddress,
			deviceDemandDetails.DeviceRequire,
			deviceDemandDetails.DeviceBuyRate)

		// todo : send tx to miners

	}
}
