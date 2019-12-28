package data

import "encoding/json"

type RegisterBeat struct {
	//RegisterType string		`json:"RegisterType"`
	Peer Peer `json:"Peer"`
	Hops int  `json:"Hops"`
}

func RegisterBeatFromJson(barr []byte) (RegisterBeat, error) {
	rb := RegisterBeat{
		Peer: Peer{},
		Hops: 0,
	}
	err := json.Unmarshal(barr, &rb)

	return rb, err
}
