package utils

import (
	"math/rand"
	"sync"
	"time"
)

var letters = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

var rngPool = sync.Pool{
	New: func() any {
		return rand.New(rand.NewSource(time.Now().UnixNano()))
	},
}

func RandomString(length int) string {
	r := rngPool.Get().(*rand.Rand)
	defer rngPool.Put(r)

	b := make([]byte, length)
	for i := range b {
		b[i] = letters[r.Intn(len(letters))]
	}
	return string(b)
}
