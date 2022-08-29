package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server/service"
	"strconv"
)

func Del(c *gin.Context) {
	_id, ok := c.GetQuery("id")
	if !ok {
		c.JSON(http.StatusBadRequest, "Please input article id to delete it.")
		return
	}
	id, err := strconv.ParseUint(_id, 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, "Please input legal article id.")
		return
	}

	oldWork, err := service.DelArticle(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
	} else {
		c.JSON(http.StatusOK, oldWork)
	}
}
