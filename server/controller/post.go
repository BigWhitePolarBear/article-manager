package controller

//
//import (
//	"github.com/gin-gonic/gin"
//	"net/http"
//	"paperSearchServer/dao"
//	"strconv"
//	"strings"
//)
//
//type work struct {
//	Title     string
//	Authors   string
//	Journal   string
//	Volume    string
//	Month     string
//	Year      *uint16
//	CdRom     string
//	EE        string
//	Publisher string
//	ISBN      string
//}
//
//func Work(c *gin.Context) {
//	tempWork := work{}
//	err := c.Bind(&tempWork)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, err.Error())
//		return
//	}
//
//	newWork := dao.Work{}
//	newWork.Title, newWork.Volume, newWork.Year, newWork.Month, newWork.CdRom, newWork.EE, newWork.ISBN =
//		tempWork.Title, tempWork.Volume, tempWork.Year, tempWork.Month, tempWork.CdRom, tempWork.EE, tempWork.ISBN
//
//	var journalID, publisherID uint
//	if tempWork.Journal != "" {
//		dao.DB.Table("journals").Where("name = ?", tempWork.Journal).Select("id").Find(&journalID)
//		if journalID != 0 {
//			newWork.JournalID = &journalID
//		} else {
//			newJournal := dao.Journal{Name: tempWork.Journal}
//			dao.DB.Table("journals").Create(&newJournal)
//			newWork.JournalID = &newJournal.ID
//		}
//	}
//	if tempWork.Publisher != "" {
//		dao.DB.Table("publishers").Where("name = ?", tempWork.Publisher).Select("id").Find(&publisherID)
//		if publisherID != 0 {
//			newWork.PublisherID = &publisherID
//		} else {
//			newPublisher := dao.Publisher{Name: tempWork.Publisher}
//			dao.DB.Table("publishers").Create(&newPublisher)
//			newWork.PublisherID = &newPublisher.ID
//		}
//	}
//
//	dao.DB.Table("works").Create(&newWork)
//	newWorkID := strconv.Itoa(int(newWork.ID))
//
//	authorNames := strings.Split(tempWork.Authors, ", ")
//	for _, authorName := range authorNames {
//		author := dao.Author{}
//		dao.DB.Table("authors").Where("name = ?", authorName).Find(&author)
//		if newWork.Authors != "" {
//			newWork.Authors += " "
//		}
//		if author.ID != 0 {
//			newWork.Authors = newWork.Authors + strconv.Itoa(int(author.ID))
//			author.Articles += " " + newWorkID
//			author.WorkCount++
//			dao.DB.Table("authors").Save(&author)
//		} else {
//			author.Name, author.Articles, author.WorkCount = authorName, newWorkID, 1
//			dao.DB.Table("authors").Create(&author)
//			newWork.Authors += strconv.Itoa(int(author.ID))
//		}
//	}
//
//	dao.DB.Table("works").Where("id = ?", newWork.ID).Update("authors", newWork.Authors)
//
//	c.JSON(http.StatusOK, newWork)
//}
