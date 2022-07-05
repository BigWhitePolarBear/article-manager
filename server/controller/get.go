package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/service"
)

func SearchArticle(c *gin.Context) {
	queries := make(map[service.QueryType]string)
	if title, ok := c.GetQuery("title"); ok {
		queries[service.TitleQuery] = title
	} else {
		c.JSON(http.StatusBadRequest, "Must input title queries!")
	}

	if author, ok := c.GetQuery("author"); ok {
		queries[service.AuthorQuery] = author
	}

	if not, ok := c.GetQuery("not"); ok {
		queries[service.NotQuery] = not
	}

	page := c.DefaultQuery("page", "1")
	if page == "0" {
		c.JSON(http.StatusBadRequest, "There is no 0th page!")
	}
	queries[service.PageQuery] = page

	_, admin := c.Get("admin")
	admin = true
	results, err := service.SearchArticle(queries, admin)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	} else {
		c.JSON(http.StatusOK, results)
	}
}

//func GetAuthor(c *gin.Context) {
//	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"Error": err.Error(),
//		})
//	}
//	name, ok := c.GetQuery("name")
//	if !ok {
//		c.JSON(http.StatusBadRequest, "Please enter name!")
//		return
//	}
//	results, err := service.GetAuthor(name, 50, page*50)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"Error": err.Error(),
//		})
//	} else {
//		c.JSON(http.StatusOK, results)
//	}
//}
//
//func GetTopAuthor(c *gin.Context) {
//	page, err := strconv.Atoi(c.Param("page"))
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"Error": err.Error(),
//		})
//		return
//	}
//	results, err := service.GetTopAuthor(50, 50*page)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{
//			"Error": err.Error(),
//		})
//	} else {
//		c.JSON(http.StatusOK, *results)
//	}
//}
