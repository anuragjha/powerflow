package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/edgexfoundry/powerflow/traderPowerflow/config"
	"github.com/edgexfoundry/powerflow/traderPowerflow/data"
	"io/ioutil"
	"net/http"
)

// runs before everything else
func init() {
	// This function will be executed before everything else.
	fmt.Println("Init Trader")
}

// Start handler
func Start(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok ok"))
}

// Getting device(s) from edgex instance
func GetEdgeXDevices(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(config.DEVICELISTAPI)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	respBody, err := readResponseBody(resp)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	ds := make([]models.Device, 0)
	err = json.Unmarshal(respBody, &ds)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// update device list var in traderDS
	for _, d := range ds {
		data.InstanceTraderDS.AddToEdgeXDevices(d)
	}

	w.WriteHeader(http.StatusOK)
	dj, _ := json.Marshal(data.InstanceTraderDS.GetEdgeXDevices())
	w.Write(dj)
}

// Getting events for device(s) from edgex instance
func GetEdgeXDevicesReadings(w http.ResponseWriter, r *http.Request) {
	for _, d := range data.InstanceTraderDS.GetEdgeXDevices() {
		uri := config.EVENTSBYDEVICEAPI + "/" + d.Name + "/10"
		resp, err := http.Get(uri)
		if err != nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		respBody, err := readResponseBody(resp)
		if err != nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		events := make([]models.Event, 0)
		err = json.Unmarshal(respBody, &events)
		if err != nil {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		for _, event := range events {
			// commented lines are -> this is happening in => AddToEdgeXDevicesEvents(event)
			//for _, reading := range event.Readings {
			//key = [reading.Name], value = reading.Value
			data.InstanceTraderDS.AddToEdgeXDevicesEvents(event)
			//}
		}
	}

	w.WriteHeader(http.StatusOK)
	dj, _ := json.Marshal(data.InstanceTraderDS.GetEdgeXDevicesEvents())
	w.Write(dj)
}

// Getting trader Id
func GetTraderId(w http.ResponseWriter, r *http.Request) {
	j, _ := json.Marshal(data.InstanceTraderDS.TraderId)
	_, _ = w.Write(j)
}

// read response body
func readResponseBody(resp *http.Response) ([]byte, error) {
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, errors.New("cannot read request body")
	}
	defer resp.Body.Close()
	return respBody, nil
}

// read request body
func readRequestBody(r *http.Request) (string, error) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", errors.New("cannot read request body")
	}
	defer r.Body.Close()
	return string(reqBody), nil
}
