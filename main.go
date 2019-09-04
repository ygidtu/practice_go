package main

import (
	"fmt"
	"os"

	// "github.com/GeertJohan/go.rice"
	"log"
	"path/filepath"

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
		Version       bool          `goptions:"-v, --version, description='Show version'"`

		goptions.Verbs
	}{
		Host:          "127.0.0.1",
		Port:          5000,
		Dir:           filepath.Dir(os.Args[0]),
		Username:      "admin",
		Password:      "admin",
		DisableDelete: false,
	}

	if parser.Version {
		println("Current version is 2.0.0")
		os.Exit(0)
	}

	goptions.ParseAndFail(&parser)

	DirFullPath, err := filepath.Abs(parser.Dir)
	if err != nil {
		log.Fatal(err)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	manageRoute(router, DirFullPath, parser.Username, parser.Password, parser.DisableDelete)

	err = router.Run(fmt.Sprintf("%s:%d", parser.Host, parser.Port))

	if err != nil {
		log.Fatal(err)
	}
}
