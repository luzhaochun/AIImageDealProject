package main

import (
	"log"

	_ "backend/filters"
	_ "backend/models"
	_ "backend/routers"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	sqlConn, err := beego.AppConfig.String("sqlconn")
	if err != nil {
		log.Fatalf("read sqlconn: %v", err)
	}
	if err := orm.RegisterDataBase("default", "mysql", sqlConn); err != nil {
		log.Fatalf("register database: %v", err)
	}
	if beego.BConfig.RunMode == "dev" {
		orm.Debug = true
	}

	beego.SetStaticPath("/static", "storage")
	beego.Run()
}
