package data

import "encoding/json"

type Peer struct {
	IpPort          string `json:"IpPort"`
	PublicKeyHexStr string `json:"publicKeyHexStr"`
}

func PeerFromJson(barr []byte) (Peer, error) {
	p := Peer{}
	err := json.Unmarshal(barr, &p)

	return p, err
}
