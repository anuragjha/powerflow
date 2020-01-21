package syncBlockchain

import (
	"sync"

	"github.com/edgexfoundry/powerflow/blockchainPowerflow/data/block"
	"github.com/edgexfoundry/powerflow/blockchainPowerflow/data/blockchain"
	"github.com/edgexfoundry/powerflow/blockchainPowerflow/data/mpt"
	s "github.com/edgexfoundry/powerflow/commonPowerFlow/identity"
)

//SyncBlockchain struct is main - shared - common - datastu
type SyncBlockchain struct {
	bc  blockchain.Blockchain `json:"blockchain"`
	mux sync.Mutex            `json:"mux"`
}

// NewBlockChain func generates a new syncBlockchain
func NewBlockChain() SyncBlockchain {
	return SyncBlockchain{bc: blockchain.NewBlockchain()}
}

//Get func takes height as input and returns list of block at that height
func (sbc *SyncBlockchain) GetLength() int32 {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.Length
}

//Get func takes height as input and returns list of block at that height
func (sbc *SyncBlockchain) Get(height int32) ([]block.Block, bool) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.Get(height)
}

// GetBlock func takes height and hash as parameter and returns a block
func (sbc *SyncBlockchain) GetBlock(height int32, hash string) (block.Block, bool) {

	sbc.mux.Lock()
	defer sbc.mux.Unlock()

	//blks, found := sbc.Get(height)
	blks, found := sbc.bc.Get(height)
	if found == true {
		for _, b := range blks {
			if b.Header.Hash == hash {
				return b, true
			}
		}
	}
	return block.Block{}, false
}

//Insert func inserts a block into blockchain in safe way
func (sbc *SyncBlockchain) Insert(block block.Block) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	sbc.bc.Insert(block)
}

// CheckParentHash func takes a block and checks if parent hash exists and return true or false
func (sbc *SyncBlockchain) CheckParentHash(insertBlock block.Block) bool {
	//sbc.mux.Lock()
	//defer sbc.mux.Unlock() apr 4

	if insertBlock.Header.Height > 1 { // good coz genesis created
		pblocks, found := sbc.Get(insertBlock.Header.Height - 1)
		if found == true {
			for _, pb := range pblocks {
				if pb.Header.Hash == insertBlock.Header.ParentHash {
					//log.Println("Parent Hash found at height :", pb.Header.Height)
					return true
				}
			}
		}
	}
	return false
}

// UpdateEntireBlockChain func takes a json and updates the existing blockchain
func (sbc *SyncBlockchain) UpdateEntireBlockChain(blockChainJson string) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	blockchain.DecodeFromJSON(&sbc.bc, blockChainJson)
}

// BlockChainToJson converts blockchain to json string
func (sbc *SyncBlockchain) BlockChainToJson() (string, error) {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return blockchain.EncodeToJSON(&sbc.bc), nil
}

// GenBlock finc takes in a mpt and returns a block for the node
// takes parentat list[0] in random height
func (sbc *SyncBlockchain) GenBlock(height int32, parentHash string, mpt mpt.MerklePatriciaTrie, nonce string, miner s.PublicIdentity) block.Block {

	var newBlock block.Block
	newBlock.Initial(height, parentHash, mpt, nonce, miner)

	//fmt.Println(" blockHash : ", newBlock.Header.Hash)
	return newBlock
}

// Show func returns blockchain in displayable format
func (sbc *SyncBlockchain) Show() string {
	return sbc.bc.Show()
}

func (sbc *SyncBlockchain) GetLatestBlocks() []block.Block {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.GetLatestBlocks() //blockchain.Chain[blockchain.Length]
}

func (sbc *SyncBlockchain) GetParentBlock(blk block.Block) block.Block {
	sbc.mux.Lock()
	defer sbc.mux.Unlock()
	return sbc.bc.GetParentBlock(blk)
}
