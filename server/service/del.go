package service

import (
	"errors"
	"server/dao"
	"strconv"
	"strings"
)

func DelWork(id uint) (dao.Article, error) {
	oldWork := dao.Article{}
	worksTable.Where("id = ?", id).Find(&oldWork)
	if oldWork.ID == 0 {
		return oldWork, errors.New("work not found")
	}
	worksTable.Delete(&oldWork)
	strOldWorkID := strconv.Itoa(int(id))

	authorIDs := strings.Fields(oldWork.Authors)
	for _, authorID := range authorIDs {
		author := dao.Author{}
		authorsTable.Where("id = ?", authorID).Find(&author)
		// This author has no works after deleting current work.
		if author.ArticleCount == 1 {
			author.Articles, author.ArticleCount = "", 0
			authorsTable.Save(&author).Delete(&author)
			continue
		}
		author.ArticleCount -= 1

		works := strings.Fields(author.Articles)
		authorWorksBuilder := strings.Builder{}
		for i, work := range works {
			if work == strOldWorkID {
				continue
			}
			if i > 0 {
				authorWorksBuilder.WriteByte(' ')
			}
			authorWorksBuilder.WriteString(work)
		}
		author.Articles = authorWorksBuilder.String()
		authorsTable.Save(&author)
	}
	return oldWork, nil
}
