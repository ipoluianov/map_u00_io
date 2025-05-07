package httpserver

import (
	"encoding/json"
	"sync"
)

type Item struct {
	Code    string
	Data    string
	History []string
}

type Storage struct {
	mtx  sync.Mutex
	data map[string]*Item
}

const (
	MaxDataSize = 1024
)

func NewStorage() *Storage {
	var c Storage
	c.data = make(map[string]*Item)
	return &c
}

var storage *Storage

func init() {
	storage = NewStorage()
}

func GetData(code string) string {
	storage.mtx.Lock()
	if data, ok := storage.data[code]; ok {
		storage.mtx.Unlock()
		return data.Data
	}
	storage.mtx.Unlock()
	return ""
}

func GetHistory(code string) []byte {
	result := make([]byte, 0)
	storage.mtx.Lock()
	if data, ok := storage.data[code]; ok {
		result, _ = json.MarshalIndent(data.History, "", "  ")
	}
	storage.mtx.Unlock()
	return result
}

func SetData(code string, data string) {
	if len(data) > MaxDataSize {
		return
	}
	storage.mtx.Lock()
	if item, ok := storage.data[code]; ok {
		item.History = append(item.History, item.Data)
		item.Data = data
	} else {
		item := &Item{
			Code:    code,
			Data:    data,
			History: make([]string, 0),
		}
		item.History = append(item.History, data)
		storage.data[code] = item
	}
	storage.mtx.Unlock()
}
