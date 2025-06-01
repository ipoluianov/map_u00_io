package httpserver

import (
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"sync"
)

type Item struct {
	Address   []byte `json:"address"`
	Data      []byte `json:"data"`
	Signature []byte `json:"signature"`
}

type Storage struct {
	mtx   sync.Mutex
	items map[string]*Item
}

const (
	MaxDataSize    = 10 * 1024
	MaxHistorySize = 10
)

func NewStorage() *Storage {
	var c Storage
	c.items = make(map[string]*Item)
	return &c
}

var storage *Storage

func init() {
	storage = NewStorage()
}

func GetData(code string) []byte {
	storage.mtx.Lock()
	if item, ok := storage.items[code]; ok {
		bs := item.Data
		storage.mtx.Unlock()
		return bs
	}
	storage.mtx.Unlock()
	return nil
}

func SetData(bs []byte) error {
	if len(bs) < 32+64 {
		return errors.New("data too short")
	}

	address := bs[:32]
	signature := bs[32:96]
	value := bs[96:]

	verifyResult := ed25519.Verify(address, value, signature)
	if !verifyResult {
		return errors.New("invalid signature")
	}

	item := Item{
		Address:   address,
		Data:      value,
		Signature: signature,
	}

	addressHex := "0x" + hex.EncodeToString(address)
	storage.mtx.Lock()
	if len(storage.items) < 1000 {
		storage.items[addressHex] = &item
	}
	storage.mtx.Unlock()
	return nil
}
