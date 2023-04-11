package handler

import (
	"log"
	"net/http"

	"github.com/KuroNeko6666/prima-api/config"
	"github.com/KuroNeko6666/prima-api/database"
	"github.com/KuroNeko6666/prima-api/database/model"
	"github.com/KuroNeko6666/prima-api/helper"
	"github.com/KuroNeko6666/prima-api/interfaces/form"
	"github.com/KuroNeko6666/prima-api/interfaces/response"
	"github.com/gofiber/fiber/v2"
)

func Follow(ctx *fiber.Ctx) error {
	var following model.User
	var form form.Follow
	var session model.Session
	var follower model.Follower
	// sessionID := ctx.Cookies("session_id")

	claims, err := helper.GetTokenAndValidate(ctx)
	if err != nil {
		return err
	}

	if r := ctx.BodyParser(&form); r != nil {
		return ctx.Status(http.StatusBadRequest).JSON(response.Message{
			Message: config.RES_BAD_REQUEST,
		})
	}

	// session.ID = sessionID
	session.ID = claims["id"].(string)
	following.ID = form.FollowingID

	if r := database.DB.Model(&model.Session{}).Select("user_id").Find(&session); r.Error != nil {
		log.Panic(r)
	}

	follower.FollowerID = session.UserID
	follower.FollowingID = following.ID

	if r := database.DB.Model(&follower).Find(&follower); r.RowsAffected > 0 {
		if er := database.DB.Model(&follower).Delete(&follower); er.Error != nil {
			return ctx.Status(http.StatusConflict).JSON(response.Message{
				Message: er.Error.Error(),
			})
		}
		return ctx.Status(http.StatusOK).JSON(response.Message{
			Message: config.RES_UNFOLLOWING,
		})
	}

	if r := database.DB.Model(&follower).Create(&follower); r.Error != nil {
		return ctx.Status(http.StatusConflict).JSON(response.Message{
			Message: r.Error.Error(),
		})
	}

	return ctx.Status(http.StatusOK).JSON(response.Message{
		Message: config.RES_FOLLOWING,
	})

}

func GetFollowers(ctx *fiber.Ctx) error {
	var users []model.User
	var followers []model.Follower
	userID := ctx.Query("user_id")
	role := ctx.Query("get_type", "follower")
	limit := ctx.QueryInt("limit", 10)
	offset := ctx.QueryInt("offset", 0)
	where, association := helper.FollowRole(role)

	if r := database.DB.Model(&model.Follower{}).Where(where, userID).Find(&followers); r.Error != nil {
		return ctx.Status(http.StatusConflict).JSON(response.Message{
			Message: config.RES_UNFOLLOWING,
		})
	}

	if r := database.DB.Model(&followers).Limit(limit).Offset(offset).Association(association).Find(&users); r != nil {
		return ctx.Status(http.StatusConflict).JSON(response.Message{
			Message: config.RES_UNFOLLOWING,
		})
	}

	return ctx.Status(http.StatusOK).JSON(response.Base{
		Message: config.RES_FINE,
		Data:    users,
	})

}
