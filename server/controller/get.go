package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/service"
	"strconv"
)

func SearchArticle(c *gin.Context) {
	query := make(map[service.QueryType]string)
	if title, ok := c.GetQuery("title"); ok {
		query[service.TitleQuery] = title
	} else {
		c.JSON(http.StatusBadRequest, "Must input title queries!")
	}
	if year, ok := c.GetQuery("year"); ok {
		query[service.YearQuery] = year
	}
	if authors, ok := c.GetQuery("authors"); ok {
		query[service.AuthorQuery] = authors
	}

	admin := false
	if c.Keys["admin"] != nil {
		admin = true
	}
	results, err := service.SearchArticle(query, admin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, results)
	}
}

func GetAuthor(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
	}
	name, ok := c.GetQuery("name")
	if !ok {
		c.JSON(http.StatusBadRequest, "Please enter name!")
		return
	}
	results, err := service.GetAuthor(name, 50, page*50)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, results)
	}
}

func GetTopAuthor(c *gin.Context) {
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}
	results, err := service.GetTopAuthor(50, 50*page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, *results)
	}
}
