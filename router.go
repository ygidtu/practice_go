package main

import (
	"fmt"
	rice "github.com/GeertJohan/go.rice"
	"github.com/foolin/gin-template/supports/gorice"
	"github.com/gin-contrib/sessions/cookie"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mholt/archiver"
)

/*
Get var from url
*/
func getVar(c *gin.Context, varName string) string {
	if temp, ok := c.Get(varName); ok {
		return temp.(string)
	}

	temp := c.Query(varName)

	if temp != "" {
		return temp
	}
	return c.Param(varName)
}

/*
The main page issues, check if logged and so on
*/
func mainPage(ctx *gin.Context) {
	session := sessions.Default(ctx)

	login := session.Get("login")

	if login == nil || !login.(bool) {
		ctx.Redirect(http.StatusFound, "/login")
	} else {
		ctx.HTML(http.StatusOK, "index.html", gin.H{})
	}
}

/*
The login page issues
*/
func loginGet(ctx *gin.Context) {
	session := sessions.Default(ctx)

	if login := session.Get("login"); login != nil && login.(bool) {
		ctx.Redirect(http.StatusFound, "/")
	}

	err := session.Get("error")

	if err == nil {
		err = false
	}

	ctx.HTML(http.StatusOK, "login.html", gin.H{"Error": err.(bool)})
}

/*
The user validation
*/
func loginPost(ctx *gin.Context) {
	session := sessions.Default(ctx)

	username, _ := ctx.Get("username")
	password, _ := ctx.Get("password")

	if ctx.PostForm("username") == username && ctx.PostForm("password") == password {
		session.Set("login", true)
		session.Set("error", false)
		_ = session.Save()

		ctx.Redirect(http.StatusFound, "/")
	} else {
		session.Set("login", false)
		session.Set("error", true)
		_ = session.Save()

		ctx.Redirect(http.StatusFound, "/login")
	}
}

/*
The Logout
*/
func logout(context *gin.Context) {
	session := sessions.Default(context)

	session.Set("login", false)
	err := session.Save()

	log.Println(err)

	context.Redirect(http.StatusFound, "/login")
}

func GetFileContentType(out *os.File) (string, error) {

	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	// Use the net/http package's handy DectectContentType function. Always returns a valid
	// content-type by returning "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return contentType, nil
}

func GetFileContentTypePath(path string) (string, error) {
	if stat, err := os.Stat(path); os.IsNotExist(err) {
		return "", err
	} else {
		if stat.IsDir() {
			return "dir", nil
		} else {

			r, err := os.Open(path)
			if r != nil {
				defer r.Close()
			}
			if err != nil {
				return "", err
			}

			return GetFileContentType(r)
		}
	}
}

/*
Download file
*/
func downloadFile(ctx *gin.Context) {
	path := getVar(ctx, "path")
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Wrong path, %s", err)})
			return
		}
	} else if stat.IsDir() {
		return
	}

	ctx.FileAttachment(path, filepath.Base(path))
}

/*
Compress file
*/
func compressFile(ctx *gin.Context) {
	path := getVar(ctx, "path")
	root := getVar(ctx, "root")

	rootInfo, err1 := os.Stat(root)
	pathInfo, err2 := os.Stat(path)

	if os.IsNotExist(err1) || os.IsNotExist(err2) {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "path not exists"})
		return
	} else if os.SameFile(rootInfo, pathInfo) {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "could not delete root path"})
		return
	}

	go func() {
		err := archiver.Archive([]string{path}, fmt.Sprintf("%s.zip", path))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"Message": err})
		}
	}()

	ctx.JSON(http.StatusOK, gin.H{"Message": fmt.Sprintf("Compress %s in background", filepath.Base(path))})
}

/*
Delete file
*/
func deleteFile(ctx *gin.Context) {

	session := sessions.Default(ctx)

	if login := session.Get("login"); login == nil || !login.(bool) {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "login first"})
		return
	}

	if ctx.MustGet("delete").(bool) {
		return
	}

	path := getVar(ctx, "path")
	if root, ok := ctx.Get("root"); !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "do not have root path"})
		return
	} else {
		rootInfo, err1 := os.Stat(root.(string))
		pathInfo, err2 := os.Stat(path)

		if os.IsNotExist(err1) || os.IsNotExist(err2) {
			ctx.JSON(http.StatusBadRequest, gin.H{"Message": "path not exists"})
			return
		} else if os.SameFile(rootInfo, pathInfo) {
			ctx.JSON(http.StatusBadRequest, gin.H{"Message": "could not delete root path"})
			return
		}
	}

	err := os.RemoveAll(path)

	status := http.StatusOK
	if err != nil {
		status = http.StatusBadRequest
	}

	ctx.JSON(status, gin.H{"Message": err})
}

/*
List files
*/
func listFile(c *gin.Context) {
	session := sessions.Default(c)

	if login := session.Get("login"); login == nil || !login.(bool) {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "login first"})
		return
	}

	disableDelete := c.MustGet("delete").(bool)

	path := getVar(c, "path")

	type Info struct {
		Name          string `json:"name"`
		Path          string `json:"path"`
		IsDir         bool   `json:"is_dir"`
		Size          int64  `json:"size"`
		Type          string `json:"type"`
		DisableDelete bool   `json:"disable_delete"`
	}

	info := make([]*Info, 0, 0)

	// use relative path instead of absolute path
	root := c.MustGet("root").(string)

	if files, err := ioutil.ReadDir(path); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": err})
		return
	} else {
		for _, f := range files {
			tempPath := strings.Replace(path, root, "", -1)

			contentType, err := GetFileContentTypePath(filepath.Join(path, f.Name()))

			if err != nil {
				contentType = err.Error()
			}

			info = append(info, &Info{
				Name:          f.Name(),
				Path:          filepath.Join(tempPath, f.Name()),
				IsDir:         f.IsDir(),
				Size:          f.Size(),
				Type:          contentType,
				DisableDelete: disableDelete,
			})
		}
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"data":    info,
			"count":   len(info),
			"current": strings.Replace(path, root, "", -1),
		},
	)
}

/*
Middleware to check path
*/
func pathMiddleware(root string, disableDelete bool) gin.HandlerFunc {

	return func(c *gin.Context) {

		path := c.Query("path")

		c.Set("path", filepath.Join(root, path))
		c.Set("root", root)
		c.Set("delete", disableDelete)

		c.Next()
	}
}

/*
Middleware to set username and password
*/
func loginMiddleware(username, password string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("username", username)
		c.Set("password", password)
		c.Next()
	}
}

func manageRoute(router *gin.Engine, root, username, password string, disableDelete bool) {

	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{MaxAge: 3600})
	router.Use(sessions.Sessions("session", store))

	// servers other static files
	staticBox := rice.MustFindBox("./views/static")
	router.StaticFS("/static", staticBox.HTTPBox())

	router.Static("/api/preview", root)

	//new template engine
	router.HTMLRender = gorice.New(rice.MustFindBox("views"))

	router.GET("/", mainPage)
	r := router.Group("/")

	r.Use(loginMiddleware(username, password))
	{
		r.GET("/login", loginGet)
		r.POST("/login", loginPost)
	}

	router.GET("/logout", logout)
	api := router.Group("/api")
	api.Use(pathMiddleware(root, disableDelete))
	{
		api.GET("/download", downloadFile)
		api.GET("/compress", compressFile)
		api.GET("/delete", deleteFile)
		api.GET("/list", listFile)
	}

}
