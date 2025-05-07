package httpserver

import "sync"

type Storage struct {
	mtx  sync.Mutex
	data map[string][]byte
}

const (
	MaxDataSize = 1024
)

func NewStorage() *Storage {
	var c Storage
	c.data = make(map[string][]byte)
	return &c
}

var storage *Storage

func init() {
	storage = NewStorage()
}

func GetData(code string) []byte {
	storage.mtx.Lock()
	if data, ok := storage.data[code]; ok {
		storage.mtx.Unlock()
		return data
	}
	storage.mtx.Unlock()
	return nil
}

func SetData(code string, data []byte) {
	if len(data) > MaxDataSize {
		return
	}
	storage.mtx.Lock()
	storage.data[code] = data
	storage.mtx.Unlock()
}
