package main

import (
	"strconv"
	"strings"
)

func record() {
	defer wg.Done()

	for {
		select {
		case article := <-match2record:
			authors := strings.Split(article.Authors, ", ")
			article.Authors = ""
			DB.Model(Article{}).Create(&article)
			authorBuilder := strings.Builder{}
			for _, author := range authors {
				if author == "" {
					continue
				}
				tempAuthor := Author{}
				DB.Model(&Author{}).Order("name").Where("name = ?", author).Find(&tempAuthor)
				if tempAuthor.ID == 0 {
					tempAuthor.Name = author
					tempAuthor.Articles = strconv.FormatUint(article.ID, 10)
					tempAuthor.ArticleCount = 1
					DB.Model(&Author{}).Create(&tempAuthor)
				} else {
					tempAuthor.Articles += " " + strconv.FormatUint(article.ID, 10)
					tempAuthor.ArticleCount++
					DB.Model(&Author{}).Where("id = ?", tempAuthor.ID).Save(&tempAuthor)
				}

				// Save the authors' id instead of name
				authorBuilder.WriteString(strconv.FormatUint(tempAuthor.ID, 10))
				authorBuilder.WriteByte(' ')
			}
			article.Authors = authorBuilder.String()
			DB.Model(&Article{}).Where("id = ?", article.ID).Update("authors", article.Authors)

		case <-matchOK:
			return
		}
	}
}
