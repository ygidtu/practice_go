package main

import (
	"fmt"
	"os"

	"github.com/GeertJohan/go.rice"
	"github.com/foolin/gin-template/supports/gorice"

	// "github.com/GeertJohan/go.rice"
	"log"
	"path/filepath"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/voxelbrain/goptions"
)

func main() {

	parser := struct {
		Host          string        `goptions:"--host, description='host'"`
		Port          int           `goptions:"--port, description='port'"`
		Dir           string        `goptions:"--dir, description='File directory'"`
		Username      string        `goptions:"--user, description='Username'"`
		Password      string        `goptions:"--passwd, description='Password'"`
		DisableDelete bool          `goptions:"--disable-delete, description='Disable delete button'"`
		Help          goptions.Help `goptions:"-h, --help, description='Show this help'"`

		goptions.Verbs
	}{
		Host:          "127.0.0.1",
		Port:          5000,
		Dir:           filepath.Dir(os.Args[0]),
		Username:      "admin",
		Password:      "admin",
		DisableDelete: false,
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

	manageRoute(router, DirFullPath, parser.Username, parser.Password, parser.DisableDelete)

	err = router.Run(fmt.Sprintf("%s:%d", parser.Host, parser.Port))

	if err != nil {
		log.Fatal(err)
	}
}
