package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/KuroNeko6666/prima-api/config"
	"github.com/KuroNeko6666/prima-api/database"
	"github.com/KuroNeko6666/prima-api/database/model"
	"github.com/KuroNeko6666/prima-api/helper"
	"github.com/KuroNeko6666/prima-api/interfaces/response"
	"github.com/gofiber/fiber/v2"
)

func Auth(ctx *fiber.Ctx) error {
	var session model.Session
	token := ctx.Get("Authorization")
	if token == "" {
		return ctx.Status(http.StatusUnauthorized).JSON(response.Message{
			Message: config.RES_MISSING_AUTH,
		})
	}

	token = strings.ReplaceAll(token, "Bearer ", "")

	claims, err := helper.ExtractClaims(token, config.SECRET_AUTH)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(response.Message{
			Message: err.Error(),
		})
	}
	sessionID := claims["id"].(string)

	// sessionID := ctx.Cookies("session_id", "")

	db := database.DB.Model(&model.Session{}).Preload("User").Where("id", sessionID).Find(&session)

	if db.Error != nil {
		return ctx.Status(http.StatusConflict).JSON(response.Message{
			Message: db.Error.Error(),
		})
	}

	if db.RowsAffected == 0 {
		return ctx.Status(http.StatusUnauthorized).JSON(response.Message{
			Message: config.RES_UNAUTHORIZED,
		})
	}

	remain := time.Until(session.ExpiredAt).String()

	if strings.Contains(remain, "-") {
		fmt.Println("sampe delete session")
		db := db.Delete(&session)

		if db.Error != nil {
			return ctx.Status(http.StatusConflict).JSON(response.Message{
				Message: db.Error.Error(),
			})
		}

		// ctx.Cookie(&fiber.Cookie{
		// 	Name:    "session_id",
		// 	Value:   session.UserID,
		// 	Expires: time.Now().Add(time.Second - 2000000000),
		// })

		return ctx.Status(http.StatusUnauthorized).JSON(response.Message{
			Message: config.RES_UNAUTHORIZED,
		})
	}

	ctx.Locals("user", session.User)

	return ctx.Next()

}
