package httpserver

import (
	"encoding/json"
	"errors"
	"sync"

	"github.com/ipoluianov/map_u00_io/utils"
)

type Item struct {
	Address     string `json:"a"`
	DisplayName string `json:"d"`
	DT          string `json:"t"`
	Value       string `json:"v"`
	Signature   string `json:"s"`
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
		bs, _ := json.MarshalIndent(item, "", "  ")
		storage.mtx.Unlock()
		return bs
	}
	storage.mtx.Unlock()
	return nil
}

func SetData(item Item) error {
	if len(item.Address) != 66 {
		return errors.New("code length is not 66")
	}

	if len(item.Signature) != 2+128 {
		return errors.New("signature length is not 130")
	}

	if len(item.Value) > MaxDataSize {
		return errors.New("data size too large")
	}

	verifyResult := utils.VerifySignature(item.Address, []byte(item.DT+item.Value), item.Signature)
	if !verifyResult {
		return errors.New("signature verification failed")
	}

	storage.mtx.Lock()
	if len(storage.items) < 1000 {
		storage.items[item.Address] = &item
	}
	storage.mtx.Unlock()
	return nil
}
