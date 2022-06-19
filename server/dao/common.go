package dao

// All type in dao package got delay double deletion strategy for their cache.

import (
	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary

	articleCount int64
	authorCount  int64
)

func init() {
	DB.Table("articles").Count(&articleCount)
	DB.Table("authors").Count(&authorCount)
}
