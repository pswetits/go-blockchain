package main

import (
	"fmt"
	"github.com/pswetits/go-blockchain/blockchain"
	"strconv"
)

func main() {
	chain := blockchain.InitBlockChain()

	fmt.Printf("\nAdding blocks (Difficulty=%d)...\n", blockchain.Difficulty)
	chain.AddBlock("First block after genesis")
	chain.AddBlock("Second block after genesis")
	chain.AddBlock("Third block after genesis")

	fmt.Printf("\n\nFull Blockchain:\n\n")

	iter := chain.Iterator()

	for {
		block := iter.Next()
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.NewProof(block)
		fmt.Printf("Valid PoW? %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevHash) == 0 {
			break
		}
	}
}
