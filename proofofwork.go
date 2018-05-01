package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

var (
	maxNonce = math.MaxInt64
)

// 24 -> 12, make it easy 000...
const targetBits = 12

// NewProofOfWork builds and returns a ProofOfWork
func getTarget(Bits int) *big.Int {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-Bits))

	return target
}

func prepareData(block *Block) []byte {
	data := bytes.Join(
		[][]byte{
			block.PrevBlockHash,
			block.Data,
			Int64ToHex(block.Timestamp),
			Int64ToHex(int64(block.Bits)),
			Int64ToHex(int64(block.Nonce)),
		},
		[]byte{},
	)

	return data
}

func calculateHash(block *Block) []byte {
	data := prepareData(block)
	hash := sha256.Sum256(data)

	return hash[:]
}

// getNonceAndRetHash : proof of work: get nonce and hash
func getNonceAndRetHash(block *Block) (int, []byte) {
	var hashInt big.Int
	var hash []byte

	target := getTarget(block.Bits)

	fmt.Printf("計算區塊hash, 區塊內容: \"%s\"\n", block.Data)
	for block.Nonce < maxNonce {
		hash = calculateHash(block)

		fmt.Printf("HASH: \r%x", hash)

		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(target) == -1 {
			break
		} else {
			block.Nonce++
		}
	}

	fmt.Println("成功!")
	fmt.Print("\n\n")

	return block.Nonce, hash[:]
}

// Validate validates block's PoW
func Validate(block *Block) bool {
	var hashInt big.Int

	data := prepareData(block)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(getTarget(block.Bits)) == -1

	return isValid
}
