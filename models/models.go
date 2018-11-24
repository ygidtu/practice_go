package models

import (
	"github.com/astaxie/beego/orm"
)


type User struct {
	Id int
	Name string
	Password string
	Admin bool
}


type File struct {
	Id int
	FileName string
	FileId string
	Parent string
	IsDir bool
}


func init() {
	// 需要在init中注册定义的model
	orm.RegisterModel(new(File))
	orm.RegisterModel(new(User))
}
