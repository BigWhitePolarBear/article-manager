package service

//
//import (
//	"errors"
//	"paperSearchServer/dao"
//	"strconv"
//	"strings"
//)
//
//func UpdateArticle(newArticle Article) (err error) {
//	oldArticle := dao.Article{}
//	dao.DB.Table("works").Where("id = ?", newArticle.ID).Find(&oldArticle)
//	if oldArticle.ID == 0 {
//		return errors.New("please input the id of article")
//	}
//
//	oldArticle.Title, oldArticle.Volume, oldArticle.Month, oldArticle.Year, oldArticle.CdRom, oldArticle.EE, oldArticle.ISBN =
//		newArticle.Title, newArticle.Volume, newArticle.Month, newArticle.Year, newArticle.CdRom, newArticle.EE, newArticle.ISBN
//
//	newAuthorNames := strings.Split(newArticle.Authors, ", ")
//	oldAuthorIDs := stringSlice(strings.Fields(oldArticle.Authors))
//	oldArticle.Authors = ""
//	strWorkID := strconv.Itoa(int(newArticle.ID))
//	newAuthorIDs := stringSlice{}
//
//	for _, newAuthorName := range newAuthorNames {
//		if oldArticle.Authors != "" {
//			oldArticle.Authors += " "
//		}
//		author := dao.Author{}
//		dao.DB.Table("authors").Where("name = ?", newAuthorName).Find(&author)
//		strAuthorID := strconv.Itoa(int(author.ID))
//		if author.ID == 0 {
//			author.Articles = strWorkID
//			author.ArticleCount = 1
//			dao.DB.Table("authors").Create(&author)
//			strAuthorID = strconv.Itoa(int(author.ID))
//			oldArticle.Authors += strAuthorID
//			newAuthorIDs = append(newAuthorIDs, strAuthorID)
//			continue
//		}
//
//		// Check if this author in the original authors.
//		if oldAuthorIDs.contains(strAuthorID) {
//			oldArticle.Authors += strAuthorID
//			newAuthorIDs = append(newAuthorIDs, strAuthorID)
//			continue
//		}
//
//		// Not original authors.
//		author.Articles += " " + strWorkID
//		author.ArticleCount++
//		dao.DB.Table("authors").Save(&author)
//
//		oldArticle.Authors += strAuthorID
//		newAuthorIDs = append(newAuthorIDs, strAuthorID)
//	}
//
//	for _, oldAuthorID := range oldAuthorIDs {
//		if !newAuthorIDs.contains(oldAuthorID) {
//			author := dao.Author{}
//			dao.DB.Table("authors").Where("ID = ?", oldAuthorID).Find(&author)
//			tempWorks := strings.Fields(author.Articles)
//			author.Articles = ""
//			for _, tempWork := range tempWorks {
//				if tempWork != strWorkID {
//					if author.Articles != "" {
//						author.Articles += " "
//					}
//					author.Articles += tempWork
//				}
//			}
//			author.WorkCount--
//			dao.DB.Table("authors").Save(&author)
//		}
//	}
//
//	if newArticle.Journal == "" {
//		oldArticle.JournalID = nil
//	} else {
//		var journalID uint
//		dao.DB.Table("journals").Where("name = ?", newArticle.Journal).Select("id").Find(&journalID)
//		if oldArticle.JournalID == nil || *oldArticle.JournalID != journalID {
//			if journalID != 0 {
//				oldArticle.JournalID = &journalID
//			} else {
//				journal := dao.Journal{Name: newArticle.Journal}
//				dao.DB.Table("journals").Create(&journal)
//				oldArticle.JournalID = &journal.ID
//			}
//		}
//	}
//
//	if newArticle.Publisher == "" {
//		oldArticle.PublisherID = nil
//	} else {
//		var PublisherID uint
//		dao.DB.Table("publishers").Where("name = ?", newArticle.Publisher).Select("id").Find(&PublisherID)
//		if oldArticle.PublisherID == nil || *oldArticle.PublisherID != PublisherID {
//			if PublisherID != 0 {
//				oldArticle.PublisherID = &PublisherID
//			} else {
//				publisher := dao.Publisher{Name: newArticle.Publisher}
//				dao.DB.Table("publishers").Create(&publisher)
//				oldArticle.JournalID = &publisher.ID
//			}
//		}
//	}
//
//	dao.DB.Table("works").Save(&oldArticle)
//}
