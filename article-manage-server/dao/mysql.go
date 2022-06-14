package dao

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

var (
	DB *gorm.DB
)

// Link to local mysql server
func init() {
	var err error

	master := "root:zxc05020519@tcp(192.168.200.128:3306)/paper_search_server?charset=utf8&interpolateParams=true&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(master),
		&gorm.Config{
			PrepareStmt: true,
		})
	if err != nil {
		panic(err)
	}

	slave := "root:zxc05020519@tcp(192.168.200.128:43306)/paper_search_server?charset=utf8&interpolateParams=true&parseTime=True&loc=Local"
	err = DB.Use(dbresolver.Register(dbresolver.Config{
		Replicas: []gorm.Dialector{mysql.Open(slave)},
		Policy:   dbresolver.RandomPolicy{},
	}))
	if err != nil {
		panic(err)
	}
}
