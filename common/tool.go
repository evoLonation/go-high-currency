package common

import (
	"hash/fnv"
	"sync"
)

func Block() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
func Hash(key string) int64 {
	h := fnv.New64()
	h.Write([]byte(key))
	return int64(h.Sum64())
}
