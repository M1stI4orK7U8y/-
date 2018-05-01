package main

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// CMD ...
type CMD struct {
	bc *Blockchain
}

func (cmd *CMD) printUsage() {
	fmt.Println("區塊鏈小作品展示")
	fmt.Println("Usage:")
	fmt.Println("  addblock  - 新增區塊")
	fmt.Println("  printchain - 顯示所有區塊 (from best to genesis)")
	fmt.Println("  getblockbyheight - 顯示指定高度的區塊")
	fmt.Println("  getblockbyhash - 顯示指定hash的區塊")
	fmt.Println("  getbestblock - 顯示best block的hash")
	fmt.Println("  getbestheight - 顯示目前高度")
}

func (cmd *CMD) addBlock(data string) {
	cmd.bc.AddBlock(data)
	fmt.Println("成功!")
}

func (cmd *CMD) printChain() {
	bci := cmd.bc.Iterator()

	for {
		block := bci.Next()

		cmd.printBlock(block)

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}
}

func (cmd *CMD) printBlock(block *Block) {
	fmt.Printf("Previous hash: %x\n", block.PrevBlockHash)
	fmt.Printf("Data: %s\n", block.Data)
	fmt.Printf("Height: %d\n", block.Height)
	fmt.Printf("Nonce: %d\n", block.Nonce)
	fmt.Printf("Hash: %x\n", calculateHash(block))
	fmt.Printf("PoW: %s\n", strconv.FormatBool(Validate(block)))
	fmt.Printf("Bits: %d\n", block.Bits)
	fmt.Println()
}

func (cmd *CMD) printBlockByHeight() {
	var height int
	fmt.Print("Target Height:")
	_, err := fmt.Scan(&height)

	// err = nil
	// height = 5

	if err != nil {
		fmt.Println("ERROR: " + err.Error())
	} else if height > cmd.bc.bestHeight {
		fmt.Printf("Blockchain best height is smaller than %d\n", height)
	} else {
		bci := cmd.bc.Iterator()

		block := bci.GetBlockByHeight(height)

		cmd.printBlock(block)
	}
}

func (cmd *CMD) printBestBlock() {
	fmt.Printf("Best Block Hash: %x \n", cmd.bc.bestHash)
}

func (cmd *CMD) printBestBlockHeight() {
	fmt.Printf("Best Block Height: %d \n", cmd.bc.bestHeight)
}

func (cmd *CMD) printBlockByHash() {
	var targetHash string
	fmt.Print("Target Hash:")
	_, err := fmt.Scanf("%s", &targetHash)

	// err = nil
	// targetHash = "000fbf9e8f55b125bb12b617352fabe56ed8e84f3eb66298834e4f5ba0351ad9"

	if err != nil {
		fmt.Printf("ERROR: %s %s \n", targetHash, err.Error())
	} else {
		hashByte, convErr := hex.DecodeString(targetHash)
		if convErr != nil {
			fmt.Printf("ERROR: %s %s \n", hashByte, err.Error())
		}

		bci := cmd.bc.Iterator()

		block := bci.GetBlockByHash(hashByte[:])

		cmd.printBlock(block)
	}
}

// Run processes commands
func (cmd *CMD) Run() {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	for {
		var inputCmd string
		fmt.Print("Enter command: ")
		_, err := fmt.Scan(&inputCmd)
		fmt.Println(inputCmd)

		// err = nil
		// inputCmd = "getblockbyhash"

		if err != nil {
			fmt.Println("ERROR: " + err.Error())
		} else {
			switch inputCmd {
			case "addblock":
				data := "Random: " + strconv.Itoa(r1.Intn(MaxInt))
				cmd.addBlock(data)
			case "printchain":
				cmd.printChain()
			case "getblockbyheight":
				cmd.printBlockByHeight()
			case "getblockbyhash":
				cmd.printBlockByHash()
			case "getbestblock":
				cmd.printBestBlock()
			case "getbestheight":
				cmd.printBestBlockHeight()
			case "exit":
				os.Exit(1)
			default:
				cmd.printUsage()
			}
		}
	}
}
