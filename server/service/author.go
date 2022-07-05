package service

import (
	"context"
	"log"
	"server/dao"
	"strconv"
)

func getAuthor(id uint64, admin bool) (author dao.Author, ok bool) {
	var jsonAuthor []byte
	err := dao.ArticleCache.Get(context.Background(), strconv.FormatUint(id, 10), &jsonAuthor)
	if err == nil {
		// got from cache
		err = json.Unmarshal(jsonAuthor, &author)
		if err != nil {
			log.Println("service/author.go getAuthor error:", err)
			ok = false
		}
		author.Name = getName(id)
	} else {
		// cache missed
		err = dao.DB.Model(&dao.Author{}).Where("id = ?", id).Find(&author).Error
		if err != nil {
			log.Println("service/author.go getAuthor error:", err)
			ok = false
		}
	}

	if !admin {
		author.ID = 0
	}

	return author, true
}

func getName(id uint64) (name string) {
	err := dao.NameCache.Get(context.Background(), strconv.FormatUint(id, 10), &name)
	if err != nil {
		author := dao.Author{}
		err = dao.DB.Model(&dao.Author{}).Where("id = ?", id).
			Select("id", "name").Find(&author).Error
		if err != nil {
			log.Println("service/author.go getName error:", err)
		}
		name = author.Name
	}

	return
}
