package diskv

import (
	"fmt"
	"math/rand"
	"testing"
)

func shuffle(keys []string) {
	ints := rand.Perm(len(keys))
	for i, _ := range keys {
		keys[i], keys[ints[i]] = keys[ints[i]], keys[i]
	}
}

func genValue(size int) []byte {
	v := make([]byte, size)
	for i := 0; i < size; i++ {
		v[i] = uint8((rand.Int() % 26) + 97) // a-z
	}
	return v
}

const (
	KEY_COUNT = 1000
)

func genKeys() []string {
	keys := make([]string, KEY_COUNT)
	for i := 0; i < KEY_COUNT; i++ {
		keys[i] = fmt.Sprintf("%d", i)
	}
	return keys
}

func (s *Store) load(keys []string, v []byte) {
	for _, k := range keys {
		s.Write(k, v)
	}
}

func benchRead(b *testing.B, size, cachesz int) {
	b.StopTimer()
	s := NewStore("speed-test", dumbXf, uint(cachesz))
	defer s.Flush()
	keys := genKeys()
	value := genValue(size)
	s.load(keys, value)
	shuffle(keys)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.Read(keys[i%len(keys)])
	}
	b.StopTimer()
}

func benchWrite(b *testing.B, size int, withIndex bool) {
	b.StopTimer()
	type Writeable interface {
		Write(k string, v []byte) error
		Flush() error
	}
	var s Writeable = nil
	if withIndex {
		s = NewOrderedStore("speed-test", dumbXf, 0)
	} else {
		s = NewStore("speed-test", dumbXf, 0)
	}
	defer s.Flush()
	keys := genKeys()
	value := genValue(size)
	shuffle(keys)
	b.SetBytes(int64(size))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.Write(keys[i%len(keys)], value)
	}
	b.StopTimer()
}

func BenchmarkWrite_32B_NoIndex(b *testing.B) {
	benchWrite(b, 32, false)
}

func BenchmarkWrite_1KB_NoIndex(b *testing.B) {
	benchWrite(b, 1024, false)
}

func BenchmarkWrite_4KB_NoIndex(b *testing.B) {
	benchWrite(b, 4096, false)
}

func BenchmarkWrite_10KB_NoIndex(b *testing.B) {
	benchWrite(b, 10240, false)
}

func BenchmarkWrite_32B_WithIndex(b *testing.B) {
	benchWrite(b, 32, true)
}

func BenchmarkWrite_1KB_WithIndex(b *testing.B) {
	benchWrite(b, 1024, true)
}

func BenchmarkWrite_4KB_WithIndex(b *testing.B) {
	benchWrite(b, 4096, true)
}

func BenchmarkWrite_10KB_WithIndex(b *testing.B) {
	benchWrite(b, 10240, true)
}

func BenchmarkRead_32B_NoCache(b *testing.B) {
	benchRead(b, 32, 0)
}

func BenchmarkRead_1KB_NoCache(b *testing.B) {
	benchRead(b, 1024, 0)
}

func BenchmarkRead_4KB_NoCache(b *testing.B) {
	benchRead(b, 4096, 0)
}

func BenchmarkRead_10KB_NoCache(b *testing.B) {
	benchRead(b, 10240, 0)
}

func BenchmarkRead_32B_WithCache(b *testing.B) {
	benchRead(b, 32, KEY_COUNT*32*2)
}

func BenchmarkRead_1KB_WithCache(b *testing.B) {
	benchRead(b, 1024, KEY_COUNT*1024*2)
}

func BenchmarkRead_4KB_WithCache(b *testing.B) {
	benchRead(b, 4096, KEY_COUNT*4096*2)
}

func BenchmarkRead_10KB_WithCache(b *testing.B) {
	benchRead(b, 10240, KEY_COUNT*4096*2)
}
