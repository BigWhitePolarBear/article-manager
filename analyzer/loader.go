package main

import (
	"log"
	"strconv"
)

func articleIDFilterAndWordCntLoader() {
	defer wg.Done()

	var (
		avgWordCnt float64

		i uint64
		// found
		j float64

		maxID uint64
	)

	err := DB.Table("articles").Select("max(id)").Find(&maxID).Error
	if err != nil {
		panic(err)
	}

	for ; i <= maxID; i++ {
		receiver := ArticleWordCount{}
		err = DB.Model(&ArticleWordCount{}).Where("id = ?", i).Find(&receiver).Error
		if err != nil {
			continue
		}

		j++
		avgWordCnt = (avgWordCnt*(j-1.) + float64(receiver.Count)) / j
		articleIDFilter.AddString(strconv.FormatUint(receiver.ID, 16))
	}

	err = DB.Model(&Variable{}).Create(&Variable{Key: "ArticleAvgWordCnt",
		Value: strconv.FormatFloat(avgWordCnt, 'g', -1, 64)}).Error
	if err != nil {
		log.Println("loader.go articleIDFilterAndWordCntLoader error:", err)
	}

	jsonFilter, err := articleIDFilter.MarshalJSON()
	if err != nil {
		log.Println("loader.go articleIDFilterAndWordCntLoader error:", err)
	} else {
		err = DB.Model(&Variable{}).Create(&Variable{Key: "ArticleIDFilter",
			Value: string(jsonFilter)}).Error
		if err != nil {
			log.Println("loader.go articleIDFilterAndWordCntLoader error:", err)
		}
	}
}

func authorIDFilterAndWordCntLoader() {
	defer wg.Done()

	var (
		avgWordCnt float64

		i uint64
		// found
		j float64

		maxID uint64
	)

	err := DB.Table("authors").Select("max(id)").Find(&maxID).Error
	if err != nil {
		panic(err)
	}

	for ; i <= maxID; i++ {
		receiver := AuthorWordCount{}
		err = DB.Model(&AuthorWordCount{}).Where("id = ?", i).Find(&receiver).Error
		if err != nil {
			continue
		}

		j++
		avgWordCnt = (avgWordCnt*(j-1.) + float64(receiver.Count)) / j
		authorIDFilter.AddString(strconv.FormatUint(receiver.ID, 16))
	}

	err = DB.Model(&Variable{}).Create(&Variable{Key: "AuthorAvgWordCnt",
		Value: strconv.FormatFloat(avgWordCnt, 'g', -1, 64)}).Error
	if err != nil {
		log.Println("loader.go authorIDFilterAndWordCntLoader error:", err)
	}

	jsonFilter, err := authorIDFilter.MarshalJSON()
	if err != nil {
		log.Println("loader.go authorIDFilterAndWordCntLoader error:", err)
	} else {
		err = DB.Model(&Variable{}).Create(&Variable{Key: "AuthorIDFilter",
			Value: string(jsonFilter)}).Error
		if err != nil {
			log.Println("loader.go authorIDFilterAndWordCntLoader error:", err)
		}
	}
}

func articleWordFilterLoader() {
	defer wg.Done()

	rows, err := DB.Model(&WordToArticle{}).Select("word").Rows()
	if err != nil {
		panic(err)
	}

	var word string

	for rows.Next() {
		err = rows.Scan(&word)
		if err != nil {
			log.Println("loader.go articleWordFilterLoader error:", err)
			continue
		}

		articleWordFilter.AddString(word)
	}

	jsonFilter, err := articleWordFilter.MarshalJSON()
	if err != nil {
		log.Println("loader.go articleWordFilterLoader error:", err)
	} else {
		err = DB.Model(&Variable{}).Create(&Variable{Key: "ArticleWordFilter",
			Value: string(jsonFilter)}).Error
		if err != nil {
			log.Println("loader.go articleWordFilterLoader error:", err)
		}
	}
}

func authorWordFilterLoader() {
	defer wg.Done()

	rows, err := DB.Model(&WordToAuthor{}).Select("word").Rows()
	if err != nil {
		panic(err)
	}

	var word string

	for rows.Next() {
		err = rows.Scan(&word)
		if err != nil {
			log.Println("loader.go authorWordFilterLoader error:", err)
			continue
		}

		authorWordFilter.AddString(word)
	}

	jsonFilter, err := authorWordFilter.MarshalJSON()
	if err != nil {
		log.Println("loader.go authorWordFilterLoader error:", err)
	} else {
		err = DB.Model(&Variable{}).Create(&Variable{Key: "AuthorWordFilter",
			Value: string(jsonFilter)}).Error
		if err != nil {
			log.Println("loader.go authorWordFilterLoader error:", err)
		}
	}
}
