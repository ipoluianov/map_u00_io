package httpserver

func GetData(code string) []byte {
	return []byte("DATA" + code)
}
