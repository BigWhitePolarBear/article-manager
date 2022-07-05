package controller

//
//import (
//	"github.com/gin-gonic/gin"
//	"net/http"
//	"paperSearchServer/service"
//)
//
//func DelWork(c *gin.Context) {
//	id := struct{ ID uint }{}
//	err := c.Bind(&id)
//
//	if err != nil {
//		c.JSON(http.StatusBadRequest, err.Error())
//		return
//	}
//
//	oldWork, err := service.DelWork(id.ID)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, err.Error())
//	} else {
//		c.JSON(http.StatusOK, oldWork)
//	}
//}
