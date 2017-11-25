package handlers

import (
	"log"
	"net/http"

	"github.com/DeFiXiK/FreshMeat/models"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
)

type ProfileController struct {
	DB *gorm.DB
}

func (pc *ProfileController) GetUserFromContext(c echo.Context) (*models.User, error) {
	ses, _ := session.Get("session", c)
	id := ses.Values["user_id"].(uint)
	user, err := models.GetUserByID(pc.DB, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (pc *ProfileController) GetProfilePage(c echo.Context) error {
	user, err := pc.GetUserFromContext(c)
	if err != nil {
		log.Fatal(err)
	}
	return c.Render(http.StatusOK, "profile.html", map[string]interface{}{
		"user": user,
	})
}

func (pc *ProfileController) UpdateProfile(c echo.Context) error {
	user, err := pc.GetUserFromContext(c)
	if err != nil {
		log.Fatal(err)
	}
	fpassword := c.FormValue("fpassword")
	spassword := c.FormValue("spassword")
	if fpassword != spassword {
		return c.Render(http.StatusBadRequest, "profile.html", map[string]interface{}{
			"error": "Введеные пароли не совпадают",
		})
	}
	user.PasswordHash = models.HashPwd(fpassword)
	user.Save(pc.DB)

	return c.Render(http.StatusOK, "profile.html", map[string]interface{}{
		"message": "Обновление профиля выполнено успешно",
	})
}
