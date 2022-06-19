package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

// Link to local mysql server
func init() {
	var err error

	dsn := "root:zxc05020519@tcp(192.168.200.128:3306)/paper_search_server?" +
		"charset=utf8mb4&interpolateParams=true&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			PrepareStmt: true,
		})
	if err != nil {
		panic(err)
	}
}
