package controllers

import (
	"beego_netdisk/models"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strconv"

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
		current := c.Input().Get("page")

		if current == "" {
			current = "0"
		}

		c.Data["current"] = current
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

	c.Data["Error"] = status
	c.TplName = "base.html"
	c.Layout = "base.html"
	c.LayoutSections = make(map[string]string)
	c.LayoutSections["PageContent"] = "login.html"
	c.Render()
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
	fileId := c.Input().Get("file_id")

	if fileId != "" {

		file := models.GetFile(fileId)
		mode, _ := os.Stat(file)

		if mode.IsDir() {
			os.RemoveAll(file)
		} else {
			os.Remove(file)
		}

		models.DeleteFile(fileId)
	}
}

// download
func (c *FileController) Get() {
	file := c.Input().Get("file_id")

	if file == "" {
		c.Abort("404")
	}

	responseFile := models.GetFile(file)

	if responseFile == "" {
		c.Abort("404")
	}

	c.Ctx.Output.Download(responseFile)
}

// create new directory
func (c *FileController) Create() {
	fileName := c.GetString("file")
	directory := c.GetString("dir")

	if directory == "" {
		directory = "./Downloads"
	}

	directory = models.GetFile(directory)

	fullPath := filepath.Join(directory, fileName)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		// path/to/whatever does not exist
		os.MkdirAll(fullPath, os.ModePerm)
	}

	go models.InsertRecord(fullPath)
}

// upload file
func (c *FileController) Upload() {
	directory := c.GetString("dir")

	f, h, err := c.GetFile("uploadname")
	if err != nil {
		log.Fatal("getfile err ", err)
	}
	defer f.Close()

	fullPath := filepath.Join(models.GetFile(directory), h.Filename)
	c.SaveToFile("uploadname", fullPath)

	go models.InsertRecord(fullPath)

	c.Ctx.Redirect(302, beego.URLFor("MainController.Get", ":page", directory))
}

// retrive file list
func (c *FileController) List() {
	fileId := strconv.Itoa(models.Hash("./Downloads"))
	file := c.Input().Get("file_id")

	if file != "" {
		fileId = file
	}

	c.Data["json"] = models.GetDir(fileId)

	c.ServeJSON()
}
