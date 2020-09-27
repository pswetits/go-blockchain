package main

import (
	"fmt"
	"github.com/pswetits/go-blockchain/blockchain"
	"strconv"
)

func main() {
	fmt.Printf("Difficulty setting: %d\n", blockchain.Difficulty)
	fmt.Printf("Adding blocks...\n\n")

	chain := blockchain.InitBlockChain()
	chain.AddBlock("First block after genesis")
	chain.AddBlock("Second block after genesis")
	chain.AddBlock("Third block after genesis")

	fmt.Printf("\n\nFull Blockchain:\n\n")
	for _, block := range chain.Blocks {
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("Valid PoW? %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
