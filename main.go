package main

import (
	"html/template"
	"io"
	"log"
	"net/http"

	"github.com/DeFiXiK/FreshMeat/handlers"
	"github.com/DeFiXiK/FreshMeat/models"
	"github.com/gorilla/sessions"

	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// Render renders a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("./templates/*.html")),
	}
	e.Renderer = renderer

	db, err := gorm.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&models.User{}, &models.Storage{}, &models.Record{})

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	authCtl := handlers.AuthController{
		DB: db,
	}

	unauthorizedGroup := e.Group("/auth")
	authorizedGroup := e.Group("")

	authorizedGroup.Use(authCtl.CheckSessionForAuthorized)
	unauthorizedGroup.Use(authCtl.CheckSessionForUnauthorized)

	authorizedGroup.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world")
	})

	unauthorizedGroup.GET("/login", authCtl.GetLoginFrom)
	unauthorizedGroup.POST("/login", authCtl.PostLoginForm)
	e.Logger.Fatal(e.Start(":5000"))
}
