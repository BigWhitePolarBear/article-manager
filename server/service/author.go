package service

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/singleflight"
	"log"
	"server/dao"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	authorGroup singleflight.Group
	nameGroup   singleflight.Group
)

// If ok == false, err had been logged in this func.
func getAuthor(id uint64, admin bool) (author dao.Author, ok bool) {
	sID := strconv.FormatUint(id, 10)

	var jsonAuthor []byte
	err := dao.AuthorCache.Get(context.Background(), sID, &jsonAuthor)
	if err == nil {
		// got from cache
		err = json.Unmarshal(jsonAuthor, &author)
		if err != nil {
			log.Println("service/author.go getAuthor error:", err)
			return author, false
		}

		author.Name = getName(id)
	} else {
		// cache missed
		_author, _err, _ := authorGroup.Do(sID, func() (interface{}, error) {
			go func() {
				time.Sleep(200 * time.Millisecond)
				authorGroup.Forget(sID)
			}()

			_author := dao.Author{}
			err = dao.DB.Model(&dao.Author{}).Where("id = ?", id).Find(&_author).Error
			return _author, err
		})
		err = _err
		if err != nil {
			log.Println("service/author.go getAuthor error:", err)
			return author, false
		}

		author = _author.(dao.Author)
	}

	articleIDs := getAuthorArticle(id)
	author.Articles = make([]string, len(articleIDs))
	for i := range articleIDs {
		author.Articles[i] = getTitle(articleIDs[i])
	}

	if !admin {
		author.ID = 0
	}

	return author, true
}

func getName(id uint64) (name string) {
	sID := strconv.FormatUint(id, 10)

	err := dao.NameCache.Get(context.Background(), sID, &name)
	if err != nil {
		// cache missed
		_name, _err, _ := nameGroup.Do(sID, func() (interface{}, error) {
			go func() {
				time.Sleep(200 * time.Millisecond)
				nameGroup.Forget(sID)
			}()

			_author := dao.Author{}
			err = dao.DB.Model(&dao.Author{}).Where("id = ?", id).
				Select("id", "name").Find(&_author).Error
			return _author.Name, err
		})
		err = _err
		if err != nil {
			log.Println("service/author.go getName error:", err)
			return
		}

		name = _name.(string)
	}

	return
}

func getCachedAuthorRes(nameWords []string, page int, admin bool) (authors []dao.Author, ok bool) {
	builder := strings.Builder{}
	sort.Strings(nameWords)
	for _, word := range nameWords {
		builder.WriteString(word)
	}
	basicKey := builder.String()

	sIDs, err := dao.AuthorResRDB.LRange(context.Background(),
		fmt.Sprintf("%s:%d", basicKey, page), 0, -1).Result()
	if err != nil {
		return nil, false
	}

	exist := dao.AuthorResRDB.Exists(context.Background(), basicKey+":1").Val()
	if exist == 0 {
		return nil, false
	}

	for _, sID := range sIDs {
		id, _ := strconv.ParseUint(sID, 10, 64)
		author, ok := getAuthor(id, admin)
		if ok {
			authors = append(authors, author)
		}
	}

	return authors, true
}

func cacheAuthorRes(nameWords []string, indexScores []IndexScore) {
	builder := strings.Builder{}
	sort.Strings(nameWords)
	for _, word := range nameWords {
		builder.WriteString(word)
	}
	basicKey := builder.String()

	for page := 1; page <= (len(indexScores)+49)/50; page++ {
		key := fmt.Sprintf("%s:%d", basicKey, page)
		for i := (page - 1) * 50; i < page*50 && i < len(indexScores); i++ {
			// store 10 base int for redis
			err := dao.AuthorResRDB.RPush(context.Background(), key,
				strconv.FormatUint(indexScores[i].ID, 10)).Err()
			if err != nil {
				log.Println("service/author.go cacheAuthorRes error:", err)
			}
		}

		// Retry 5 times.
		var expireErr error
		expireErr = errors.New("")
		for i := 0; expireErr != nil && i < 5; i++ {
			expireErr = dao.AuthorResRDB.Expire(context.Background(), key, time.Minute).Err()
		}
		if expireErr != nil {
			log.Println("service/author.go cacheAuthorRes error:", expireErr)
		}
	}

	return
}
