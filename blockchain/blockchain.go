package blockchain

import (
	"fmt"
	badger "github.com/dgraph-io/badger/v2"
)

const (
	dbPath = "./db/blocks"

	// DB keys
	kLastHash = "lh"
)

type BlockChain struct {
	LastHash []byte
	Database *badger.DB
}

type BlockChainIterator struct {
	CurrentHash []byte
	Database    *badger.DB
}

func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte

	// get previous hash
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(kLastHash))
		Handle(err)
		lastHash, err = item.ValueCopy(nil)
		return err
	})
	Handle(err)

	// create new block
	block := CreateBlock(data, lastHash)

	// load block to db
	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(block.Hash, block.Serialize())
		Handle(err)

		// update LastHash
		err = txn.Set([]byte(kLastHash), block.Hash)
		chain.LastHash = block.Hash
		return err
	})
	Handle(err)
}

func InitBlockChain() *BlockChain {
	var lastHash []byte

	fmt.Println("Initializing BadgerDB...")
	err := CreateDir(dbPath)
	Handle(err)

	opts := badger.DefaultOptions(dbPath)
	db, err := badger.Open(opts)
	Handle(err)

	// create & store genesis block, or fetch lastHash
	err = db.Update(func(txn *badger.Txn) error {
		// check for presence of blockchain
		if _, err := txn.Get([]byte(kLastHash)); err == badger.ErrKeyNotFound {
			fmt.Println("No existing blockchain found. Initializing Genesis.")
			// create genesis block, serialize, and load to db
			genesis := genesis()
			fmt.Printf("Genesis proved. Hash: %x\n", genesis.Hash)

			// load genesis
			err = txn.Set(genesis.Hash, genesis.Serialize())
			Handle(err)

			// load lastHash
			err = txn.Set([]byte(kLastHash), genesis.Hash)
			lastHash = genesis.Hash

			return err
		} else {
			// found existing chain
			item, err := txn.Get([]byte(kLastHash))
			Handle(err)

			// set lastHash
			lastHash, err = item.ValueCopy(nil)
			fmt.Printf("Found existing Blockchain (LastHash: %x).\n", lastHash)
			return err
		}
	})
	Handle(err)

	// create in-mem blockchain
	blockchain := BlockChain{lastHash, db}
	return &blockchain
}

func (chain *BlockChain) Iterator() *BlockChainIterator {
	iter := &BlockChainIterator{chain.LastHash, chain.Database}
	return iter
}

func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		// get current block
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)
		// decode block
		encodedBlock, err := item.ValueCopy(nil)
		block = Deserialize(encodedBlock)

		return err
	})
	Handle(err)

	// update iterator
	iter.CurrentHash = block.PrevHash

	return block
}
