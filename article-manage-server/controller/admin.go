package controller

import (
	"github.com/gin-gonic/gin"
)

func LoginMidWare() gin.HandlerFunc {
	accounts := gin.Accounts{"admin": "123"}
	return gin.BasicAuth(accounts)
}

func AdminMark(c *gin.Context) {
	c.Keys["admin"] = struct{}{}
	c.Next()
}
