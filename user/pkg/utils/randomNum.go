package utils

import "math/rand"

const letterBytes = "ABCDEF0123456789"

func Random16(n int) string {
	res := make([]byte, n)
	for i := 0; i < n; i++ {
		res[i] = letterBytes[rand.Int31n(int32(16))]
	}
	return string(res)
}
