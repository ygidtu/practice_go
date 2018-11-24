package controllers

import (
	"log"
	"beego_netdisk/models"
	"html/template"
	"strconv"
	"math/big"
	"os"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

type UserController struct {
	beego.Controller
}

type FileController struct {
	beego.Controller
}

// main page
func (c *MainController) Get() {

	ok := c.GetSession("authorized")

	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	if ok == nil {
		c.Ctx.Redirect(302, "/login")
	} else if true {

		// set page
		c.TplName = "base.html"
		c.Layout = "base.html"
		c.LayoutSections = make(map[string]string)
		c.LayoutSections["PageContent"] = "index.html"

		c.Render()

	} else {
		c.Ctx.Redirect(302, "/login")
	}
}

// login page
func (c *UserController) Login() {
	id := c.GetString("username")
	password := c.GetString("passwd")
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())

	if c.Ctx.Request.Method == "GET" {
		c.TplName = "base.html"
		c.Layout = "base.html"
		c.LayoutSections = make(map[string]string)
		c.LayoutSections["PageContent"] = "login.html"
		c.Render()
	}

	status := models.CheckUser(id, password)

	if status == "Success" {
		c.SetSession("authorized", true)
		c.Ctx.Redirect(302, "/")
	}

	c.Data["json"] = models.LoginError{ status }
	c.ServeJSON()
}

// logout page
func (c *UserController) Logout() {
	ok := c.GetSession("authorized")

	if ok != nil {
		c.DelSession("authorized")
		c.DestroySession()
	}
	c.Ctx.Redirect(302, "/")
}

// disable XSRF protection in FileController
func (c *FileController) Prepare() {
	c.EnableXSRF = false
}

// delete function
func (c *FileController) Delete() {
	file_id := c.Input().Get("file_id")

	if file_id != "" {
		file_id_int, err := strconv.Atoi(file_id)

		if err != nil {
			log.Fatal(err)
			return 
		}

		file := models.GetFile(file_id_int)
		mode, _ := os.Stat(file)
	
		if mode.IsDir() {
			os.RemoveAll(file)
		} else {
			os.Remove(file)
		}
	
		models.DeleteFile(file_id_int)
	}
}

// download
func (c *FileController) Get() {
	file := c.Input().Get("file_id")

	if file == "" {
		c.Abort("404")
	}

	file_id, err := strconv.Atoi(file)

	if err != nil {
		log.Fatal(err)
		return
	}
	
	responseFile := models.GetFile(file_id)

	if responseFile == "" {
		c.Abort("404")
	}

	c.Ctx.Output.Download(responseFile)
}

// retrive file list
func (c *FileController) List() {
	file_id := models.Hash("./Downloads")
	file := c.Input().Get("file_id")

	if file != "" {
		file_id = file
	}

	c.Data["json"] = models.GetDir(file_id)

	c.ServeJSON()
}
