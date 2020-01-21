package identity

//package main

import (
	"crypto"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	goEthCrypto "github.com/ethereum/go-ethereum/crypto"
)

type Identity struct {
	PrivateKey       *ecdsa.PrivateKey `json:"privateKey"`
	PrivateKeyBytes  []byte            `json:"privateKeyBytes"`
	PrivateKeyHexStr string            `json:"privateKeyHexStr"`
	PublicKey        crypto.PublicKey  `json:"publicKey"`
	PublicKeyECDSA   *ecdsa.PublicKey  `json:"publicKeyECDSA"`
	PublicKeyBytes   []byte            `json:"publicKeyBytes"`
	PublicKeyHexStr  string            `json:"publicKeyHexStr"`
	Address          string            `json:"address"`

	Label  string `json:"label"`
	IpPort string `json:"ipPort"`
}

// taken from = https://goethereumbook.org/wallet-generate/ => generate_wallet.go
// have renamed to generateIdentity
func generateIdentity() (*ecdsa.PrivateKey, []byte, string, crypto.PublicKey, *ecdsa.PublicKey, []byte, string, string) {
	// go-ethereum crypto package, method for generating a random private key
	privateKey, err := goEthCrypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	// Convert it to bytes by importing the golang crypto/ecdsa package and using the FromECDSA method
	privateKeyBytes := goEthCrypto.FromECDSA(privateKey)
	// Convert it to a hexadecimal string by using the go-ethereum hexutil package which provides the
	//Encode method which takes a byte slice. Then we strip of the 0x after it's hex encoded.
	privateKeyHex := hexutil.Encode(privateKeyBytes)[2:]
	fmt.Println(privateKeyHex) // 0xfad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19
	// This is the private key which is used for signing transactions and is to be treated like a password
	//and never be shared, since who ever is in possesion of it will have access to all your funds.
	//// ////
	// Public key is derived from the private key, go-ethereum's crypto private key has a Public method that will return the public key
	publicKey := privateKey.Public()
	// Converting it to hex is a similar process that we went through with the private key.
	//We strip off the 0x and the first 2 characters 04 which is always the EC prefix and is not required.
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	publicKeyBytes := goEthCrypto.FromECDSAPub(publicKeyECDSA)
	publicKeyHexStr := hexutil.Encode(publicKeyBytes)[4:]
	fmt.Println(publicKeyHexStr) // 0x049a7df67f79246283fdc93af76d4f8cdd62c4886e8cd870944e817dd0b97934fdd7719d0810951e03418205868a5c1b40b192451367f28e0088dd75e15de40c05
	// go-ethereum crypto package has a PubkeyToAddress method which accepts an ECDSA public key, and returns the public address
	// we have the public key we can easily generate the public address which is what we are used to seeing
	address := goEthCrypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	// The public address is simply the Keccak-256 hash of the public key, and then we take the last 40 characters (20 bytes) and prefix it with 0x
	fmt.Println(address) // 0x96216849c49358B10257cb55b28eA603c874b05E
	// added return
	return privateKey, privateKeyBytes, privateKeyHex, publicKey, publicKeyECDSA, publicKeyBytes, publicKeyHexStr, address
}

// taken from = https://goethereumbook.org/signature-generate/
// have renamed to generateSignature
func GenerateSignature(privateKeyHexStr string, data []byte) []byte {
	// Components for generating a signature are: the signers private key, and the hash of the data that will be signed
	// load private key
	privateKey, err := goEthCrypto.HexToECDSA(privateKeyHexStr /*"fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19"*/)
	if err != nil {
		log.Fatal(err)
	}
	// take the Keccak-256 of the data that we wish to sign
	hash := goEthCrypto.Keccak256Hash(data)
	fmt.Println(hash.Hex()) // 0x1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8
	// sign the hash with our private, which gives us the signature
	signature, err := goEthCrypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	// printing signature in byte array format to hex string
	fmt.Println(hexutil.Encode(signature)) // 0x789a80053e4927d0a898db8e065e948f5cf086e32f9ccaa54c1908e22ac430c62621578113ddbb62d509bf6049b8fb544ab06d36f916685a2eb8e57ffadde02301
	return signature
}

// taken from = https://goethereumbook.org/signature-verify/
// have renamed to generateSignature
func VerifySignature(data []byte, signature []byte, publicKeyBytesToVerify []byte) bool {
	hash := goEthCrypto.Keccak256Hash(data)
	fmt.Println(hash.Hex()) // 0x1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8
	//  VerifySignature function which takes in the signature, hash of the original data,
	//and the public key in bytes format. It returns a boolean which will be true if the public key
	//matches the signature's signer. An important gotcha is that we must first remove the last byte
	//of the signture because it's the ECDSA recover ID which must not be included.
	signatureNoRecoverID := signature[:len(signature)-1] // remove recovery id
	verified := goEthCrypto.VerifySignature(publicKeyBytesToVerify, hash.Bytes(), signatureNoRecoverID)
	//fmt.Println(verified) // true
	return verified
}

func SavePrivateKeyHexStrToFile(file string, key *ecdsa.PrivateKey) {
	err := goEthCrypto.SaveECDSA(file, key)
	if err != nil {
		fmt.Println("Error in SavePrivateKeyHexStrToFile : ", err)
	}
}

func LoadPrivateKeyECDSAFromFile(file string) (*ecdsa.PrivateKey, error) {
	return goEthCrypto.LoadECDSA(file)
}

// Generate New Identity
func NewIdentity(label string, ipPort string, privateKeyFile string, identiyFile string) Identity {
	privateKey, privateKeyBytes, privateKeyHex, publicKey, publicKeyECDSA, publicKeyBytes, publicKeyHexStr, address := generateIdentity()
	id := Identity{
		PrivateKey:       privateKey,
		PrivateKeyBytes:  privateKeyBytes,
		PrivateKeyHexStr: privateKeyHex,
		PublicKey:        publicKey,
		PublicKeyECDSA:   publicKeyECDSA,
		PublicKeyBytes:   publicKeyBytes,
		PublicKeyHexStr:  publicKeyHexStr,
		Address:          address,
		Label:            label,
		IpPort:           ipPort,
	}

	idJson, _ := json.Marshal(&id)
	fmt.Println("idJson")
	fmt.Println(string(idJson))

	data := []byte("hello")
	//signature := GenerateSignature(id.PrivateKeyHexStr /*"fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19"*/, data)
	signature := GenerateSignature(id.PrivateKeyHexStr, data)
	fmt.Println(VerifySignature(data, signature, id.PublicKeyBytes))

	// saving privateKeyHexStr to file
	SavePrivateKeyHexStrToFile(privateKeyFile, id.PrivateKey)
	// saving identity in json to a file
	_ = ioutil.WriteFile(identiyFile, idJson, 0644)

	return id
}

// Get Identity from File
func OldIdentity(label string, ipPort string, privateKeyFile string) Identity {
	privateKey, privateKeyBytes, privateKeyHex, publicKey, publicKeyECDSA, publicKeyBytes, publicKeyHexStr, address := loadIdentity(privateKeyFile)
	id := Identity{
		PrivateKey:       privateKey,
		PrivateKeyBytes:  privateKeyBytes,
		PrivateKeyHexStr: privateKeyHex,
		PublicKey:        publicKey,
		PublicKeyECDSA:   publicKeyECDSA,
		PublicKeyBytes:   publicKeyBytes,
		PublicKeyHexStr:  publicKeyHexStr,
		Address:          address,
		Label:            label,
		IpPort:           ipPort,
	}

	idJson, _ := json.Marshal(&id)
	fmt.Println("idJson")
	fmt.Println(string(idJson))

	data := []byte("hello")
	signature := GenerateSignature(id.PrivateKeyHexStr /*"fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19"*/, data)
	fmt.Println(VerifySignature(data, signature, id.PublicKeyBytes))

	//SavePrivateKeyHexStrToFile("privateKeyFile.txt", id.PrivateKey)

	return id
}

//////

func loadIdentity(privateKeyFile string) (*ecdsa.PrivateKey, []byte, string, crypto.PublicKey, *ecdsa.PublicKey, []byte, string, string) {
	privateKey, err := LoadPrivateKeyECDSAFromFile(privateKeyFile)
	if err != nil {
		log.Fatal("Cannot load private key :", err)
	}
	// Convert it to bytes by importing the golang crypto/ecdsa package and using the FromECDSA method
	privateKeyBytes := goEthCrypto.FromECDSA(privateKey)
	// Convert it to a hexadecimal string by using the go-ethereum hexutil package which provides the
	//Encode method which takes a byte slice. Then we strip of the 0x after it's hex encoded.
	privateKeyHex := hexutil.Encode(privateKeyBytes)[2:]
	fmt.Println(privateKeyHex) // 0xfad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19
	// This is the private key which is used for signing transactions and is to be treated like a password
	//and never be shared, since who ever is in possesion of it will have access to all your funds.
	//// ////
	// Public key is derived from the private key, go-ethereum's crypto private key has a Public method that will return the public key
	publicKey := privateKey.Public()
	// Converting it to hex is a similar process that we went through with the private key.
	//We strip off the 0x and the first 2 characters 04 which is always the EC prefix and is not required.
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	publicKeyBytes := goEthCrypto.FromECDSAPub(publicKeyECDSA)
	publicKeyHexStr := hexutil.Encode(publicKeyBytes)[4:]
	fmt.Println(publicKeyHexStr) // 0x049a7df67f79246283fdc93af76d4f8cdd62c4886e8cd870944e817dd0b97934fdd7719d0810951e03418205868a5c1b40b192451367f28e0088dd75e15de40c05
	// go-ethereum crypto package has a PubkeyToAddress method which accepts an ECDSA public key, and returns the public address
	// we have the public key we can easily generate the public address which is what we are used to seeing
	address := goEthCrypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	// The public address is simply the Keccak-256 hash of the public key, and then we take the last 40 characters (20 bytes) and prefix it with 0x
	fmt.Println(address) // 0x96216849c49358B10257cb55b28eA603c874b05E
	// added return
	return privateKey, privateKeyBytes, privateKeyHex, publicKey, publicKeyECDSA, publicKeyBytes, publicKeyHexStr, address
}

func LoadIdentityFromFile(filePath string) Identity {
	file, _ := ioutil.ReadFile(filePath)

	data := Identity{}

	_ = json.Unmarshal([]byte(file), &data)

	return data

}
