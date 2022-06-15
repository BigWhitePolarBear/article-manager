package main

import (
	"fmt"
	"strconv"
	"strings"
)

type set map[uint64]struct{}

func (s set) put(num uint64) {
	s[num] = struct{}{}
}

func (s set) putSlice(nums []uint64) {
	for _, num := range nums {
		s.put(num)
	}
}

func (s set) exist(num uint64) (ok bool) {
	_, ok = s[num]
	return
}

// there should not be error
func (s set) serialize() string {
	builder := strings.Builder{}

	for key := range s {
		builder.WriteString(fmt.Sprintf("%x ", key))
	}

	return builder.String()
}

func (s set) unSerialize(str string) (err error) {
	indexes := strings.Fields(str)

	for _, index := range indexes {
		id, err := strconv.ParseUint(index, 16, 64)
		if err != nil {
			return err
		}
		s[id] = struct{}{}
	}

	return nil
}
