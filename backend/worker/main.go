package main

import (
	"context"
	"log"
	"os"
	"time"

	"backend/models"
	"backend/pkg/redisx"

	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	_ = os.Chdir(findBackendRoot())

	sqlConn, _ := beego.AppConfig.String("sqlconn")
	if err := orm.RegisterDataBase("default", "mysql", sqlConn); err != nil {
		log.Fatal(err)
	}
	_ = models.Tenant{} // ensure models package init (RegisterModel)

	ctx := context.Background()
	if err := redisx.Ping(ctx); err != nil {
		log.Fatal("redis:", err)
	}
	log.Println("imagedeal worker started")

	for {
		msg, err := redisx.Dequeue(ctx, redisx.QueueParseCSV, 3*time.Second)
		if err == nil && msg != nil {
			log.Println("parse_csv", msg.TenantId)
			if err := handleParseCSV(ctx, msg.TenantId, msg.Payload); err != nil {
				log.Println("parse_csv err:", err)
			}
			continue
		}
		msg, err = redisx.Dequeue(ctx, redisx.QueueDownload, 3*time.Second)
		if err == nil && msg != nil {
			log.Println("download", string(msg.Payload))
			if err := handleDownload(ctx, msg.TenantId, msg.Payload); err != nil {
				log.Println("download err:", err)
			}
			continue
		}
	}
}

func findBackendRoot() string {
	wd, _ := os.Getwd()
	if _, err := os.Stat(wd + "/conf/app.conf"); err == nil {
		return wd
	}
	return wd + "/backend"
}
