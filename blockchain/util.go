package blockchain

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
)

func ToHex(num uint64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, num)
	Handle(err)

	return buf.Bytes()
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func CreateDir(path string) error {
	err := os.MkdirAll(path, 0755)
	// dir may already exist
	if os.IsExist(err) {
		return nil
	} else {
		return err
	}
}
