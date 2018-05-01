package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

const dbFile = "blockchainDemo.db"
const blocksBucket = "blocks"

// Blockchain : sequence of Blocks
type Blockchain struct {
	bestHash   []byte
	bestHeight int
	db         *bolt.DB
}

// BlockchainIterator : go through a blockchain
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

// AddBlock : add a block into blockchain
func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte
	var block *Block
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if len(bc.bestHash) == 0 {
			lastHash = b.Get([]byte("l"))
		} else {
			lastHash = bc.bestHash
		}

		encodedBlock := b.Get(lastHash)
		block = Deserialize(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	newBlock, lastHash := NewBlock(data, lastHash, block.Height+1)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(lastHash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("l"), lastHash)
		if err != nil {
			log.Panic(err)
		}

		bc.bestHash = lastHash
		bc.bestHeight = newBlock.Height

		return nil
	})
}

// Iterator ...
func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.bestHash, bc.db}

	return bci
}

// Next : return the previous block of current block (hash)
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = Deserialize(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevBlockHash

	return block
}

// GetBlockByHash : get block by hash
func (i *BlockchainIterator) GetBlockByHash(targetHash []byte) *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(targetHash)
		block = Deserialize(encodedBlock)

		return nil
	})

	if err != nil {
		fmt.Printf("No such bloch with hash %x\n", targetHash)
	}

	return block
}

// GetBlockByHeight : get block by height
func (i *BlockchainIterator) GetBlockByHeight(targetHeight int) *Block {
	var block *Block

	for {
		block = i.Next()

		if block.Height == targetHeight {
			break
		}

		if block.Height == 0 {
			break
		}

	}

	return block
}

// NewBlockchain : create new blockchain
func NewBlockchain() *Blockchain {
	var bestHash []byte
	var bestHeight int
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			fmt.Println("區塊鏈不存在, 自動產生區塊鏈(genesis block)")
			genesis, hash := NewGenesisBlock()

			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}

			err = b.Put(hash, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}

			err = b.Put([]byte("l"), hash)
			if err != nil {
				log.Panic(err)
			}
			bestHash = hash
			bestHeight = genesis.Height
		} else {
			lastHash := b.Get([]byte("l"))
			encodedBlock := b.Get(lastHash)
			block := Deserialize(encodedBlock)
			bestHash = lastHash
			bestHeight = block.Height

		}

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{bestHash, bestHeight, db}

	return &bc
}
