package main

import (
	"fmt"
	"github.com/GeertJohan/go.rice"
	"github.com/foolin/gin-template/supports/gorice"
	"os"

	// "github.com/GeertJohan/go.rice"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/voxelbrain/goptions"
	"log"
	"path/filepath"
)

/*
TODO
2019.3.17 添加视频和页面的爬取功能
*/

func MiddleDB(mongo *MongoDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := mongo.GetCollection()
		if err != nil {
			c.Abort()
		} else {
			c.Set("mongo", mongo)
			c.Next()
		}
	}
}

func main() {

	parser := struct {
		Host     string        `goptions:"-h, --host, description='host'"`
		Port     int           `goptions:"-p, --port, description='port'"`
		Dir      string        `goptions:"-d, --dir, description='File directory'"`
		Username string        `goptions:"--user, description='Username'"`
		Password string        `goptions:"--passwd, description='Password'"`
		Help     goptions.Help `goptions:"-h, --help, description='Show this help'"`

		goptions.Verbs
	}{
		Host:     "127.0.0.1",
		Port:     5000,
		Dir:      filepath.Dir(os.Args[0]),
		Username: "admin",
		Password: "admin",
	}
	goptions.ParseAndFail(&parser)

	DirFullPath, err := filepath.Abs(parser.Dir)
	if err != nil {
		log.Fatal(err)
	}

	gin.SetMode(gin.TestMode)

	//db := &MongoDB{URI: parser.URI}

	router := gin.Default()
	//router.Use(MiddleDB(db))

	store := cookie.NewStore([]byte("secret"))
	store.Options(sessions.Options{MaxAge: 3600})
	router.Use(sessions.Sessions("session", store))

	// servers other static files
	staticBox := rice.MustFindBox("static")
	router.StaticFS("/static_", staticBox.HTTPBox())
	// router.Static("/static_", "static")

	//new template engine
	router.HTMLRender = gorice.New(rice.MustFindBox("views"))
	// router.LoadHTMLGlob("views/*")

	// provide file download service
	router.Static("/files", DirFullPath)

	manageRoute(router, DirFullPath, parser.Username, parser.Password)

	err = router.Run(fmt.Sprintf("%s:%d", parser.Host, parser.Port))

	if err != nil {
		log.Fatal(err)
	}
}
