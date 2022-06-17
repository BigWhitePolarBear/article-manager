package main

import (
	"fmt"
	"strconv"
	"strings"
)

type InvertedIndex map[uint64]int

func (s InvertedIndex) Add(num uint64) {
	s[num]++
}

func (s InvertedIndex) AddSlice(nums []uint64) {
	for _, num := range nums {
		s.Add(num)
	}
}

func (s InvertedIndex) Get(num uint64) (cnt int) {
	cnt = s[num]
	return
}

// Serialize there should not be error
func (s InvertedIndex) Serialize() string {
	builder := strings.Builder{}

	for key := range s {
		builder.WriteString(fmt.Sprintf("%x,%x ", key, s[key]))
	}

	return builder.String()
}

func (s InvertedIndex) UnSerialize(str string) (err error) {
	temp := strings.Fields(str)

	for _, t := range temp {
		_t := strings.Split(t, ",")
		id, err := strconv.ParseUint(_t[0], 16, 64)
		if err != nil {
			return err
		}
		cnt, err := strconv.Atoi(_t[1])
		if err != nil {
			return err
		}
		s[id] = cnt
	}

	return nil
}
