package utils

func CompareMethodId(methodBS []byte, data []byte) bool {
	if len(data) >= len(methodBS) {
		for i := 0; i < len(methodBS); i++ {
			if data[i] != methodBS[i] {
				return false
			}
		}
	}
	return true
}
