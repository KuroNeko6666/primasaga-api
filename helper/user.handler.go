package helper

import (
	"github.com/KuroNeko6666/prima-api/database"
	"github.com/KuroNeko6666/prima-api/database/model"
)

func GetUserLoggedIn(user *model.User, sessionID string) error {
	var session model.Session

	if r := database.DB.Model(&model.Session{}).Select("user_id").Find(&session); r.Error != nil {
		return r.Error
	}

	if r := database.DB.Model(&session).Association("User").Find(&user); r != nil {
		return r
	}

	return nil
}
