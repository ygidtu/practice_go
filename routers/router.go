package routers

import (
	"beego_netdisk/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/login", &controllers.UserController{}, "*:Login")
	beego.Router("/logout", &controllers.UserController{}, "get:Logout")
	beego.Router("/file/delete", &controllers.FileController{}, "post:Delete")
	beego.Router("/file/create", &controllers.FileController{}, "post:Create")
	beego.Router("/file/list", &controllers.FileController{}, "get:List")
	beego.Router("/file/upload", &controllers.FileController{}, "post:Upload")
}
