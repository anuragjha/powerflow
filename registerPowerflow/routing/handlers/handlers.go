package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/edgexfoundry/powerflow/registerPowerflow/data"
	"io/ioutil"
	"log"
	"net/http"
)

// runs before everything else
func init() {
	// This function will be executed before everything else.
	fmt.Println("Init Register")
}

// Start handler
func Start(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Register ok"))
}

// Show peers
func Peers(w http.ResponseWriter, r *http.Request) {
	rdsC := data.InstanceRegisterDS.GetRegisterDSCopy()
	rdsCJson, err := json.Marshal(&rdsC)
	if err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(rdsCJson)
}

// Register self to some other Register
func RegisterSelfTo(w http.ResponseWriter, r *http.Request) {
	reqBody, err := readRequestBody(r) // will receive registerBeat for Register peer
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	p, err := data.PeerFromJson(reqBody) // get RegisterBeat struct from reqBody
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// registering itself and adding other register peers in self registerDs
	processRegisterSelfTo(p)
}

// forming register request and sending it to ipPort mentioned in "RegisterSelfTo req body"
// adding peer IpPorts to Register peers data structure
func processRegisterSelfTo(p data.Peer) {
	// forming Register Beat
	url := "http://" + p.IpPort + "/register"
	rb := data.RegisterBeat{
		Peer: data.Peer{
			IpPort: data.InstanceRegisterDS.IpPort,
		},
		Hops: 2,
	}
	rbJson, _ := json.Marshal(&rb)
	// sending register POST request to a register server
	resp, err := http.Post(url, "application/json", bytes.NewReader(rbJson))
	if err != nil {
		log.Println("Could not send register beat to :" + p.IpPort + "\n" + err.Error())
	} else {
		data.InstanceRegisterDS.AddToRegisterPeers(p) // add peer to RegisterPeers, if response received
	}
	respBody, _ := readResponseBody(resp)

	peers := make(map[string]data.Peer)
	_ = json.Unmarshal(respBody, &peers)

	for _, peer := range peers {
		data.InstanceRegisterDS.AddToRegisterPeers(peer)
	}
}

// Register peer in Peers
func Register(w http.ResponseWriter, r *http.Request) {
	reqBody, err := readRequestBody(r) // will receive registerBeat for Register peer
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rb, err := data.RegisterBeatFromJson(reqBody) // get RegisterBeat struct from reqBody
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// if a node receives a forwarded RegistrationBeat request and it does not contain its own address
	if data.InstanceRegisterDS.IpPort != rb.Peer.IpPort {
		data.InstanceRegisterDS.AddToRegisterPeers(rb.Peer) // add peer in request to Register Peers
	}
	// 200 OK response
	w.WriteHeader(http.StatusOK)
	peers := data.InstanceRegisterDS.GetRegisterPeers()
	delete(peers, rb.Peer.IpPort) // deleting peer who sent the request
	peersJson, _ := json.Marshal(&peers)
	w.Write(peersJson) // returning other register peers
	// forwarding register request to other register peers
	go forwardRegisteration("/register", rb)
}

// Register peer in blockchainPeers
func RegisterForBlockchain(w http.ResponseWriter, r *http.Request) {
	reqBody, err := readRequestBody(r) // will receive registerBeat for Register peer
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rb, err := data.RegisterBeatFromJson(reqBody) // get RegisterBeat struct from reqBody
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// if a node receives a forwarded RegistrationBeat request and it does not contain its own address
	if data.InstanceRegisterDS.IpPort != rb.Peer.IpPort {
		data.InstanceRegisterDS.AddToBlockchainPeers(rb.Peer) // add peer in request to Register Peers
	}
	// 200 OK response
	w.WriteHeader(http.StatusOK)
	peers := data.InstanceRegisterDS.GetBlockchainPeers()
	// delete(peers, rb.Peer.IpPort) // deleting peer who sent the request
	peersJson, _ := json.Marshal(&peers)
	w.Write(peersJson) // returning other register peers
	// forwarding register request to other register peers
	go forwardRegisteration("/register/blockchain", rb)
}

// forward registration to all register peers
func forwardRegisteration(endpoint string, rb data.RegisterBeat) {
	missIpPort := rb.Peer.IpPort
	rb.Hops -= 1
	if rb.Hops > 0 {
		barr, _ := json.Marshal(&rb)
		for _, peer := range data.InstanceRegisterDS.GetRegisterPeers() {
			if peer.IpPort != missIpPort {
				url := "http://" + peer.IpPort + endpoint
				fmt.Println("Forwarding to : " + url)
				http.Post(url, "application/json", bytes.NewReader(barr))
			}
		}
	}
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
func readRequestBody(r *http.Request) ([]byte, error) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return []byte{}, err
	}
	defer r.Body.Close()
	return reqBody, nil
}
