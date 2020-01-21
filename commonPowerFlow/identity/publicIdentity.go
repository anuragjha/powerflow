package identity

import (
	"crypto"
	"encoding/json"
	"log"
)

type PublicIdentity struct {
	PublicKey crypto.PublicKey `json:"publicKey"`
	//HashForKey 	hash.Hash 		`json="hashForKey"`
	Label string `json:"label"`
}

func (id *Identity) GetMyPublicIdentity() PublicIdentity {
	pid := PublicIdentity{}
	pid.PublicKey = id.PublicKey
	//pid.HashForKey = id.HashForKey
	pid.Label = id.Label

	return pid
}

func (pid *PublicIdentity) PublicIdentityToJson() string {
	jsonBytes, err := json.Marshal(&pid)
	if err != nil {
		log.Println("Error in marshalling publicIdentity, err - ", err)
		return "{}"
	}
	return string(jsonBytes)
}

func JsonToPublicIdentity(str string) PublicIdentity {
	pid := PublicIdentity{}
	if len(str) > 0 {
		err := json.Unmarshal([]byte(str), &pid)
		if err != nil {
			log.Println("Error in Unmarshalling publicIdentity, err - ", err)
		}
	}
	return pid
}
