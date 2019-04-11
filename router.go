package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/mholt/archiver"
	"github.com/pkg/errors"
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
		ctx.HTML(http.StatusOK, "index.html", gin.H{"Mode": "index"})
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

	ctx.HTML(http.StatusOK, "index.html", gin.H{"Error": err.(bool), "Mode": "login"})
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

/*
The only API of mongoDB
*/
//func filesAPI(ctx *gin.Context) {
//	session := sessions.Default(ctx)
//
//	login := session.Get("login")
//
//	if login == nil || !login.(bool) {
//		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Please login first"})
//		return
//	} else {
//		db, ok := ctx.Keys["mongo"].(*MongoDB)
//
//		start, err := strconv.Atoi(ctx.DefaultQuery("start", "0"))
//
//		if err != nil {
//			ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Start format error, %s", err)})
//			return
//		}
//
//		length, err := strconv.Atoi(ctx.DefaultQuery("length", "10"))
//		if err != nil {
//			ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Start format error, %s", err)})
//			return
//		}
//
//		if !ok {
//			ctx.JSON(http.StatusBadRequest, gin.H{"message": ok})
//			return
//		}
//
//		data, err := db.Paginate(start, length)
//
//		if err != nil {
//			ctx.JSON(http.StatusBadRequest, gin.H{"message": err})
//			return
//		}
//
//		ctx.JSON(http.StatusOK, gin.H{
//			"data":   data,
//			"total":  db.Total,
//			"start":  start,
//			"length": length,
//		})
//	}
//}

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
Decode base64 encoded path
*/
func decodePath(path string) (string, error) {

	if path == "" {
		return "", errors.New("Path required")
	}

	pathBytes, err := base64.StdEncoding.DecodeString(path)

	if err != nil {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			return path, nil
		} else {
			return "", errors.New(fmt.Sprintf("Wrong encoding, %s", err))
		}
	}

	path = string(pathBytes)

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return path, nil
	}

	return path, errors.New(fmt.Sprintf("%s not exists", path))
}

/*
Download file
*/
func downloadFile(ctx *gin.Context) {
	path := getVar(ctx, "path")

	log.Println(path)

	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Wrong path, %s", err)})
			return
		}
	} else if stat.IsDir() {
		//ctx.Redirect(http.StatusFound, fmt.Sprintf("/api/compress?path=%s", path))
		return
	}

	contentRange := ctx.GetHeader("Content-Range")

	reader, err := os.Open(path)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": fmt.Sprintf("Wrong path, %s", err)})
		return
	}
	defer reader.Close()

	contentType, _ := GetFileContentType(reader)
	stat, _ = reader.Stat()

	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(reader.Name())),
	}

	var start int64 = 0
	re := regexp.MustCompile(`bytes\s+(%d+)-(%d+)?`)
	if contentRange != "" {
		matched := re.FindStringSubmatch(contentRange)

		if len(matched) > 1 {
			start, err = strconv.ParseInt(matched[0], 10, 64)

			if err != nil {
				start = 0
			}
		}
	}

	_, err = reader.Seek(start, 0)

	if err != nil {
		reader, _ = os.Open(path)
		defer reader.Close()
	}

	ctx.DataFromReader(http.StatusOK, stat.Size()-start, contentType, reader, extraHeaders)
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
				Path:          base64.StdEncoding.EncodeToString([]byte(filepath.Join(tempPath, f.Name()))),
				IsDir:         f.IsDir(),
				Size:          f.Size(),
				Type:          contentType,
				DisableDelete: disableDelete,
			})
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": info, "count": len(info)})
}

/*
Middleware to check path
*/
func pathMiddleware(root string, disabelDelete bool) gin.HandlerFunc {

	return func(c *gin.Context) {

		path := c.Query("path")
		path, _ = decodePath(path)

		c.Set("path", filepath.Join(root, path))
		c.Set("root", root)
		c.Set("delete", disabelDelete)

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

func manageRoute(router *gin.Engine, root, username, password string, disabelDelete bool) {
	router.GET("/", mainPage)
	r := router.Group("/")

	r.Use(loginMiddleware(username, password))
	{
		r.GET("/login", loginGet)
		r.POST("/login", loginPost)
	}

	router.GET("/logout", logout)
	api := router.Group("/api")
	api.Use(pathMiddleware(root, disabelDelete))
	{
		api.GET("/download", downloadFile)
		api.GET("/compress", compressFile)
		api.GET("/delete", deleteFile)
		api.GET("/list", listFile)
		api.Static("/file", root)
	}

}
