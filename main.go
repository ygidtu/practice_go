package main

import (
	"beego_netdisk/models"
	_ "beego_netdisk/routers"
	"flag"
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	orm.RegisterDriver("sqlite", orm.DRSqlite)
	orm.RegisterDataBase("default", "sqlite3", "./data.db")
	orm.RunSyncdb("default", false, true)
}

func config() {
	beego.BConfig.AppName = "netdisk_beego"
	beego.BConfig.Listen.HTTPAddr = "127.0.0.1"
	beego.BConfig.Listen.HTTPPort = 5000
	beego.BConfig.WebConfig.StaticDir["/static_"] = "static"
	beego.BConfig.WebConfig.AutoRender = false
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.EnableXSRF = true
	beego.BConfig.WebConfig.XSRFKey = "download"
	beego.BConfig.WebConfig.XSRFExpire = 3600

	// models.UpdateDatabase()
}

func main() {
	user := flag.String("u", "user", "Username")
	password := flag.String("p", "passwd", "Password")

	if len(os.Args) > 1 {
		flag.Parse()
		models.CreateUser(*user, *password)	
	} else {
		config()
		beego.Run()
	}
}
