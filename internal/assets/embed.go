package assets

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// All dist files embedded into the binary
//
//go:embed dist
var DistFiles embed.FS // Exported for direct access in main.go

// GetDistFS returns the embedded dist files filesystem
func GetDistFS() http.FileSystem {
	distFS, err := fs.Sub(DistFiles, "dist")
	if err != nil {
		panic(err)
	}
	return http.FS(distFS)
}

// GetAssetsFS returns the embedded assets folder filesystem
func GetAssetsFS() http.FileSystem {
	assetsFS, err := fs.Sub(DistFiles, "dist/assets")
	if err != nil {
		panic(err)
	}
	return http.FS(assetsFS)
}

// SetupRoutes configures the router to use embedded assets
func SetupRoutes(router *gin.Engine) {
	// Serve root index.html for the React app
	router.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		content, err := fs.ReadFile(DistFiles, "dist/index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Error reading index.html")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", content)
	})

	// Directly serve specific static files
	router.StaticFileFS("/index.html", "index.html", GetDistFS())
	router.StaticFileFS("/vite.svg", "vite.svg", GetDistFS())

	// Serve all assets files for the React build
	router.StaticFS("/assets", GetAssetsFS())
}
