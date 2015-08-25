package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
)

var Orm *xorm.Engine
var (
	DB_TYPE = "mysql"
	DB_NAME = "root:root@/crawler?charset=utf8"
)

type Info struct {
	Id    int64
	Title string `xorm: "varchar(50)"`
	Url   string `xorm: "varchar(500)"`
}

func init() {
	var err error
	Orm, err = xorm.NewEngine(DB_TYPE, DB_NAME)
	checkerror(err, "create engine")

	err = Orm.CreateTables(new(Info))
	checkerror(err, "create tables")
}

func Add(title, url string) error {
	_, err := Orm.Insert(&Info{Title: title, Url: url})
	checkerror(err, "insert a item")
	return nil
}

func checkerror(err error, text string) {
	if err != nil {
		log.Fatalf("Failed to "+text+" : %v", err)
	}
}
