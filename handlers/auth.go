package handlers

import (
	"net/http"

	"github.com/gorilla/sessions"

	"github.com/DeFiXiK/FreshMeat/models"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
)

type AuthController struct {
	DB *gorm.DB
}

func (ac *AuthController) GetRegistrationForm(c echo.Context) error {
	return c.Render(http.StatusOK, "registration.html", nil)
}

func (ac *AuthController) PostRegistrationForm(c echo.Context) error {
	username := c.FormValue("username")
	fpassword := c.FormValue("fpassword")
	spassword := c.FormValue("spassword")
	if fpassword != spassword {
		return c.Render(http.StatusBadRequest, "registration.html", map[string]interface{}{
			"error": "Введеные пароли не совпадают",
		})
	}
	user := models.User{}
	user.Username = username
	user.PasswordHash = models.HashPwd(fpassword)
	user.Create(ac.DB)
	return c.Render(http.StatusOK, "register.html", map[string]interface{}{
		"message": "Аккаунт успешно создан",
	})
}

func (ac *AuthController) GetLoginFrom(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}

func (ac *AuthController) PostLoginForm(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	log.Println("Данные считали")
	user, err := models.GetUserByName(ac.DB, username)
	if err != nil {
		return c.String(http.StatusNotFound, err.Error())
	}
	log.Println("Пользователь есть в системе")
	if !user.CheckPassword(password) {
		return c.String(http.StatusNotFound, "Пароль не подходит")
	}
	log.Println("Авторизация успешна")
	ses, _ := session.Get("session", c)
	ses.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	ses.Values["user_id"] = user.ID
	ses.Save(c.Request(), c.Response())

	return c.Redirect(http.StatusFound, "/")
}

func (ac *AuthController) CheckSessionForAuthorized(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ses, err := session.Get("session", c)
		if err != nil {
			return c.Redirect(http.StatusFound, "/auth/login")
		}
		id, ok := ses.Values["user_id"]
		if !ok {
			return c.Redirect(http.StatusFound, "/auth/login")
		}
		if id == "zero" {
			return c.Redirect(http.StatusFound, "/auth/login")
		}
		_, err = models.GetUserByID(ac.DB, id.(uint))
		if err != nil {
			return c.Redirect(http.StatusFound, "/auth/login")
		}
		return next(c)
	}
}

func (ac *AuthController) CheckSessionForUnauthorized(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ses, err := session.Get("session", c)
		if err != nil {
			return next(c)
		}
		id, ok := ses.Values["user_id"]
		if !ok {
			return next(c)
		}
		if id == "zero" {
			return next(c)
		}
		_, err = models.GetUserByID(ac.DB, id.(uint))
		if err != nil {
			return next(c)
		}
		return c.Redirect(http.StatusFound, "/")
	}
}
