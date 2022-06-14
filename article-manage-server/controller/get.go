package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"paperSearchServer/service"
	"strconv"
)

func GetWork(c *gin.Context) {
	query := make(map[service.QueryType]string)
	if title, ok := c.GetQuery("title"); ok {
		query[service.TitleQuery] = title
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
	page, err := strconv.Atoi(c.DefaultQuery("page", "0"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}
	results, err := service.GetWork(query, 50, page*50, admin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
	} else {
		c.JSON(http.StatusOK, *results)
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
