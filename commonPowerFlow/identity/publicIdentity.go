package identity

import (
	"crypto"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common/hexutil"
	goEthCrypto "github.com/ethereum/go-ethereum/crypto"
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

func (pid *PublicIdentity) GetPublicKeyBytes() []byte {
	// Converting it to hex is a similar process that we went through with the private key.
	//We strip off the 0x and the first 2 characters 04 which is always the EC prefix and is not required.
	publicKeyECDSA, ok := pid.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	publicKeyBytes := goEthCrypto.FromECDSAPub(publicKeyECDSA)
	publicKeyHexStr := hexutil.Encode(publicKeyBytes)[4:]
	fmt.Println("GetPublicKeyHexStr : " + publicKeyHexStr)
	return publicKeyBytes
}
