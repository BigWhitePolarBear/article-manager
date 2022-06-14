package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"paperSearchServer/service"
)

func UpdateArticle(c *gin.Context) {
	newArticle := service.Article{}
	err := c.Bind(&newArticle)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = service.UpdateArticle(newArticle)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	} else {
		c.JSON(http.StatusOK, "Update successfully")
	}
}
