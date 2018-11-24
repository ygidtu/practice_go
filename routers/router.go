package routers

import (
	"beego_netdisk/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
    beego.Router("/login", &controllers.UserController{}, "*:Login")
    beego.Router("/logout", &controllers.UserController{}, "get:Logout")
    // beego.Router("/delete", &controllers.FileController{}, "post:Delete")
    beego.RESTRouter("/file", &controllers.FileController{})
    beego.Router("/list", &controllers.FileController{}, "get:List")
}
