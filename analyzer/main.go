package main

import (
	"github.com/bits-and-blooms/bloom/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var (
	wg sync.WaitGroup

	DB *gorm.DB

	// ID filter use hexadecimal string number.
	articleIDFilter *bloom.BloomFilter
	authorIDFilter  *bloom.BloomFilter

	articleWordFilter *bloom.BloomFilter
	authorWordFilter  *bloom.BloomFilter
)

func main() {
	var err error
	DB, err = gorm.Open(mysql.Open("root:zxc05020519@tcp(localhost:3306)/"+
		"article_manager?charset=utf8mb4&interpolateParams=true&parseTime=True&loc=Local"),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		})
	if err != nil {
		panic(err)
	}

	articleIDFilter = bloom.NewWithEstimates(1e7, 0.05)
	articleWordFilter = bloom.NewWithEstimates(2e6, 0.05)
	authorIDFilter = bloom.NewWithEstimates(1e7, 0.05)
	authorWordFilter = bloom.NewWithEstimates(2e6, 0.05)

	wg.Add(4)

	go articleIDFilterAndWordCntLoader()
	go authorIDFilterAndWordCntLoader()
	go articleWordFilterLoader()
	go authorWordFilterLoader()

	wg.Wait()
}
