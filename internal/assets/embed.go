package assets

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Static files embedded into the binary
//
//go:embed web/static
var StaticFiles embed.FS

// Template files embedded into the binary
//
//go:embed web/templates
var TemplateFiles embed.FS

// GetStaticFS returns the embedded static files filesystem
func GetStaticFS() http.FileSystem {
	staticFS, err := fs.Sub(StaticFiles, "web/static")
	if err != nil {
		panic(err)
	}
	return http.FS(staticFS)
}

// GetTemplates returns parsed templates from embedded files
func GetTemplates() *template.Template {
	tmpl, err := template.ParseFS(TemplateFiles, "web/templates/*.html")
	if err != nil {
		panic(err)
	}
	return tmpl
}

// SetupRoutes configures the router to use embedded assets
func SetupRoutes(router *gin.Engine) {
	// Serve static files from embedded filesystem
	router.StaticFS("/static", GetStaticFS())

	// Load templates from embedded filesystem
	router.SetHTMLTemplate(GetTemplates())
}
