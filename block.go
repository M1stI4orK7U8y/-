package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

// Block a block
type Block struct {
	Timestamp     int64
	Height        int
	Data          []byte
	PrevBlockHash []byte
	Nonce         int
	Bits          int
}

// Serialize the block
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

// NewBlock : creates and returns a new Block anh hash
func NewBlock(data string, prevBlockHash []byte, height int) (*Block, []byte) {
	block := &Block{time.Now().Unix(), height, []byte(data), prevBlockHash, 0, targetBits}
	nonce, hash := getNonceAndRetHash(block)

	block.Nonce = nonce

	return block, hash
}

// NewGenesisBlock : creates and returns genesis Block and hash
func NewGenesisBlock() (*Block, []byte) {
	retval, hash := NewBlock("Genesis Block", []byte{}, 0)
	return retval, hash
}

// Deserialize a block
func Deserialize(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}
