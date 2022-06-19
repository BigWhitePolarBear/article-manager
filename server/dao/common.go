package dao

// All type in dao package got delay double deletion strategy for their cache.

import (
	jsoniter "github.com/json-iterator/go"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
)
