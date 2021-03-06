package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/peterbourgon/diskv"
	"io"
)

const (
	transformBlockSize = 2 // grouping of chars per directory depth
)

func BlockTransform(s string) []string {
	sliceSize := len(s) / transformBlockSize
	pathSlice := make([]string, sliceSize)
	for i := 0; i < sliceSize; i++ {
		from, to := i*transformBlockSize, (i*transformBlockSize)+transformBlockSize
		pathSlice[i] = s[from:to]
	}
	return pathSlice
}

func main() {
	s := diskv.NewStore("data-dir", BlockTransform, 1024*1024)

	data := []string{
		"I am the very model of a modern Major-General",
		"I've information vegetable, animal, and mineral",
		"I know the kings of England, and I quote the fights historical",
		"From Marathon to Waterloo, in order categorical",
		"I'm very well acquainted, too, with matters mathematical",
		"I understand equations, both the simple and quadratical",
		"About binomial theorem I'm teeming with a lot o' news",
		"With many cheerful facts about the square of the hypotenuse",
	}
	for _, valueStr := range data {
		key, value := md5sum(valueStr), bytes.NewBufferString(valueStr).Bytes()
		s.Write(key, value)
	}

	keyChan, keyCount := s.Keys(), 0
	for key, ok := <-keyChan; ok; key, ok = <-keyChan {
		value, err := s.Read(key)
		if err != nil {
			panic(fmt.Sprintf("key %s had no value", key))
		}
		fmt.Printf("%s: %s\n", key, value)
		keyCount++
	}
	fmt.Printf("%d total keys\n", keyCount)

	// s.Flush() // leave it commented out to see how data is kept on disk
}

func md5sum(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}
