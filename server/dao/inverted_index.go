package dao

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// InvertedIndex is a mapping of id and counts belonging to a word.
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

// Serialize there should not be error.
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

// Union return the union set of inverted-indexes in the slice.
func Union(InvertedIndexes []InvertedIndex) (indexes []uint64, err error) {
	n := len(InvertedIndexes)
	if n == 0 {
		return nil, errors.New("the input of Union func is empty")
	}

	firstInvertedIndex := InvertedIndexes[0]
Loop:
	for word := range firstInvertedIndex {
		for i := 1; i < n; i++ {
			_, exist := InvertedIndexes[i][word]
			if !exist {
				continue Loop
			}
		}
		indexes = append(indexes, word)
	}

	return
}
