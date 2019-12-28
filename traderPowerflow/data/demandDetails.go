package data

// struct to hold consumer demand(require) details // to be used as non shared variables
type DemandDetails struct {
	DeviceName    string
	DeviceAddress string // base url => http://Ip:port/

	DeviceRequire string
	DeviceBuyRate string
}

//type DevicesDemandDetails struct {
//
//}
