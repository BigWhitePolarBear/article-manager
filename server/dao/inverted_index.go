package dao

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// InvertedIndex is a mapping of id and counts belonging to a Word.
type InvertedIndex map[uint64]uint64

func (s InvertedIndex) Add(num uint64) {
	s[num]++
}

func (s InvertedIndex) AddSlice(nums []uint64) {
	for _, num := range nums {
		s.Add(num)
	}
}

func (s InvertedIndex) Get(num uint64) (cnt uint64) {
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

func (s *InvertedIndex) UnSerialize(str string) (err error) {
	*s = InvertedIndex{}

	temp := strings.Fields(str)

	for _, t := range temp {
		_t := strings.Split(t, ",")
		id, err := strconv.ParseUint(_t[0], 16, 64)
		if err != nil {
			return err
		}
		cnt, err := strconv.ParseUint(_t[1], 16, 64)
		if err != nil {
			return err
		}
		(*s)[id] = cnt
	}

	return nil
}

// Intersection return the union set of inverted-IDs in the slice.
func Intersection(InvertedIndexes []InvertedIndex) (indexes []uint64, err error) {
	n := len(InvertedIndexes)
	if n == 0 {
		return nil, errors.New("the input of Intersection func is empty")
	}

	// Shorter index first.
	sort.Slice(InvertedIndexes, func(i, j int) bool {
		return len(InvertedIndexes[i]) < len(InvertedIndexes[j])
	})

Loop:
	for id := range InvertedIndexes[0] {
		for i := 1; i < n; i++ {
			_, exist := InvertedIndexes[i][id]
			if !exist {
				continue Loop
			}
		}
		indexes = append(indexes, id)
	}

	return
}

func Union(InvertedIndexes []InvertedIndex) (invertedIndex InvertedIndex, err error) {
	n := len(InvertedIndexes)
	if n == 0 {
		return InvertedIndex{}, errors.New("the input of Union func is empty")
	}

	invertedIndex = map[uint64]uint64{}

	for i := range InvertedIndexes {
		for id := range InvertedIndexes[i] {
			invertedIndex[id] = 1
		}
	}

	return
}
