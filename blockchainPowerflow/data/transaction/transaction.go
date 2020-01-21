package transaction

import (
	"encoding/hex"
	"encoding/json"
	"github.com/edgexfoundry/powerflow/commonPowerFlow/identity"
	"golang.org/x/crypto/sha3"
	"log"
	"strconv"
	"time"
	//s "../p5security"
	//"../client"
	//s "../identity"
	//"github.com/edgexfoundry/powerflow/commonPowerFlow/identity"
)

type Transaction struct {
	Id        string                  `json:"id"`
	From      identity.PublicIdentity `json:"from"`
	To        identity.PublicIdentity `json:"to"` //if To is empty then its a borrowing tx
	ToTxId    string                  `json:"toTxId"`
	Tokens    float64                 `json:"tokens"`
	Fees      float64                 `json:"fees"`
	TxType    string                  `json:"txtype"`
	Timestamp time.Time               `json:"timestamp"`
}

func NewTransaction(from identity.PublicIdentity, to identity.PublicIdentity, toTxId string, tokens float64, fees float64, txType string) Transaction {
	tx := Transaction{
		From:      from,
		To:        to,
		ToTxId:    toTxId,
		Tokens:    tokens,
		Fees:      fees,
		TxType:    txType,
		Timestamp: time.Now(),
	}

	tx.Id = tx.genId()

	return tx
}

func (tx *Transaction) genId() string {
	str := tx.From.PublicIdentityToJson() +
		tx.To.PublicIdentityToJson() +
		tx.ToTxId +
		strconv.FormatFloat(float64(tx.Tokens), 'f', -1, 64) +
		strconv.FormatFloat(float64(tx.Fees), 'f', -1, 64) +
		tx.TxType +
		tx.Timestamp.String()
	sum := sha3.Sum256([]byte(str))
	return hex.EncodeToString(sum[:])
}

func (tx *Transaction) Show() string {
	str := "\ntx id :" + tx.Id +
		"\ntx From :" + tx.From.PublicIdentityToJson() +
		"\ntx To :" + tx.To.PublicIdentityToJson() +
		"\ntx ToTxId :" + tx.ToTxId +
		"\ntx Tokens :" + strconv.FormatFloat(float64(tx.Tokens), 'f', -1, 64) +
		"\ntx Fees :" + strconv.FormatFloat(float64(tx.Fees), 'f', -1, 64) +
		"\ntx Type :" + tx.TxType +
		"\ntx Time :" + tx.Timestamp.String() + "\n"
	return str
}

func (tx *Transaction) CreateTxSig(fromCid identity.Identity) []byte {
	return identity.GenerateSignature(fromCid.PrivateKeyHexStr, tx.TransactionToJsonByteArray())
	//old - return fromCid.GenSignature(tx.TransactionToJsonByteArray())
}

func (tx *Transaction) CreateTxSigForMiner(fromId identity.Identity) []byte {
	return identity.GenerateSignature(fromId.PrivateKeyHexStr, tx.TransactionToJsonByteArray())
	// old - return fromId.GenSignature(tx.TransactionToJsonByteArray())
}

func VerifyTxSig(fromPid identity.PublicIdentity, tx Transaction, txSig []byte) bool {
	return identity.VerifySignature(tx.TransactionToJsonByteArray(), txSig, fromPid.GetPublicKeyBytes())
	// old - return s.VerifySingature(fromPid.PublicKey, tx.TransactionToJsonByteArray(), txSig)
}

func (tx *Transaction) TransactionToJsonByteArray() []byte {
	txJson, err := json.Marshal(tx)
	if err != nil {
		log.Println("in TransactionToJsonByteArray : Error in marshalling Tx : ", err)
	}

	return txJson
}

func (tx *Transaction) TransactionToJson() string {
	txJson, err := json.Marshal(tx)
	if err != nil {
		log.Println("in TransactionToJsonByteArray : Error in marshalling Tx : ", err)
	}

	return string(txJson)
}

func JsonToTransaction(txJson string) Transaction {
	tx := Transaction{}
	err := json.Unmarshal([]byte(txJson), &tx)
	if err != nil {
		log.Println("Error in unmarshalling Transaction, err - ", err)
		log.Println("String given to unmarshall Transaction, ================> \n ", txJson, "\nxxxxxxxxxxxxxxxxxxxxxxxxxxx\n")
	}

	return tx
}
