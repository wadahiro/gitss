package main

import (
	"log"

	"html/template"
	"net/http"
	"strings"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-gonic/contrib/renders/multitemplate"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/nu7hatch/gouuid"
	"github.com/wadahiro/gitss/server/controller"
	"github.com/wadahiro/gitss/server/indexer"
)

func initRouter(indexer indexer.Indexer, port string, debugMode bool, gitDataDir string) {
	if !debugMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(gin.Recovery())

	r.Use(func(c *gin.Context) {
		c.Set("indexer", indexer)
	})

	r.Use(func(c *gin.Context) {
		c.Set("gitDataDir", gitDataDir)
	})

	// r.LoadHTMLGlob("server/templates/*")
	r.HTMLRender = loadTemplates("react.html")

	// We can't use router.Static method to use '/' for static files.
	// see https://github.com/gin-gonic/gin/issues/75
	// r.StaticFS("/", assetFS())

	r.Use(func(c *gin.Context) {
		id, _ := uuid.NewV4()
		c.Set("uuid", id)
	})

	apiPrefix := "/api/v1/"

	// add API routes
	r.GET(apiPrefix+"version", func(c *gin.Context) {
		c.JSON(200, map[string]string{
			"version":    Version,
			"commitHash": CommitHash,
		})
	})

	r.GET(apiPrefix+"search", controller.SearchIndex)

	// react server-side rendering
	// react := NewReact(
	// 	"assets/js/bundle.js",
	// 	debugMode,
	// 	r,
	// )
	// r.GET("/", react.Handle)
	// r.GET("/issues", react.Handle)

	r.Use(static.Serve("/", BinaryFileSystem("assets")))

	r.Run(":" + port)
}

func loadTemplates(list ...string) multitemplate.Render {
	r := multitemplate.New()

	for _, x := range list {
		templateString, err := Asset("server/templates/" + x)
		if err != nil {
			log.Fatal(err)
		}

		tmplMessage, err := template.New(x).Parse(string(templateString))
		if err != nil {
			log.Fatal(err)
		}

		r.Add(x, tmplMessage)
	}

	return r
}

type binaryFileSystem struct {
	fs http.FileSystem
}

func (b *binaryFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *binaryFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

func BinaryFileSystem(root string) *binaryFileSystem {
	fs := &assetfs.AssetFS{Asset, AssetDir, root}
	return &binaryFileSystem{
		fs,
	}
}
