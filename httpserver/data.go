package httpserver

import (
	"encoding/json"
	"sync"
	"time"
)

type Value struct {
	DT   string `json:"dt"`
	Data string `json:"data"`
}

type Item struct {
	Code    string
	Data    Value
	History []Value
}

type Storage struct {
	mtx  sync.Mutex
	data map[string]*Item
}

const (
	MaxDataSize    = 1024
	MaxHistorySize = 100
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
		return data.Data.Data
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
	dt := time.Now().Format("2006-01-02 15:04:05")
	storage.mtx.Lock()
	if item, ok := storage.data[code]; ok {
		item.History = append(item.History, item.Data)
		item.Data.DT = dt
		item.Data.Data = data

		if len(item.History) > MaxHistorySize {
			item.History = item.History[len(item.History)-MaxHistorySize:]
		}
	} else {
		value := Value{DT: dt, Data: data}
		item := &Item{
			Code:    code,
			Data:    value,
			History: make([]Value, 0),
		}
		item.History = append(item.History, value)
		storage.data[code] = item
	}
	storage.mtx.Unlock()
}
