package handler

import (
	"log"
	"net/http"
	"sort"

	"github.com/KuroNeko6666/prima-api/config"
	"github.com/KuroNeko6666/prima-api/database"
	"github.com/KuroNeko6666/prima-api/database/model"
	"github.com/KuroNeko6666/prima-api/helper"
	"github.com/KuroNeko6666/prima-api/interfaces/response"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func GetUsers(ctx *fiber.Ctx) error {
	var user model.User
	var users []model.User
	var session model.Session
	var resUsers []response.Users

	claims, err := helper.GetTokenAndValidate(ctx)
	if err != nil {
		return err
	}

	// sessionID := ctx.Cookies("session_id")
	search := ctx.Query("search", "")
	limit := ctx.QueryInt("limit", 10)
	offset := ctx.QueryInt("offset", 0)

	// session.ID = sessionID
	session.ID = claims["id"].(string)

	if r := database.DB.Model(&model.Session{}).Select("user_id").Find(&session); r.Error != nil {
		log.Panic(r)
	}

	user.ID = session.UserID

	if r := database.DB.Model(&user).Find(&user); r.Error != nil {
		log.Panic(r)
	}

	r := database.DB.Model(&users).Limit(limit).Offset(offset).
		Preload("Following", func(tx *gorm.DB) *gorm.DB {
			return tx.Where("id = ?", user.ID)
		}).Not("id = ?", user.ID).
		Where("username LIKE ?", "%"+search+"%").
		Find(&users)

	if r.Error != nil {
		log.Panic(r)
	}

	if r.RowsAffected == 0 {
		return ctx.Status(http.StatusOK).JSON(response.Base{
			Message: config.RES_FINE,
			Data:    users,
		})
	}

	for _, item := range users {
		following := item.Following
		var resUser response.Users
		copier.CopyWithOption(&resUser, item, copier.Option{IgnoreEmpty: true})
		if len(following) == 0 {
			resUser.FollowStatus = "follow"
		} else {
			resUser.FollowStatus = "followed"
		}
		resUsers = append(resUsers, resUser)
	}

	sort.Slice(resUsers, func(i, j int) bool {
		return resUsers[i].FollowStatus > resUsers[j].FollowStatus
	})

	return ctx.Status(http.StatusOK).JSON(response.Base{
		Message: config.RES_FINE,
		Data:    resUsers,
	})

}

func GetUser(ctx *fiber.Ctx) error {
	var user model.User
	var users []model.User
	var session model.Session
	var resUsers []response.Users

	claims, err := helper.GetTokenAndValidate(ctx)
	if err != nil {
		return err
	}

	session.ID = claims["id"].(string)
	if r := database.DB.Model(&model.Session{}).Select("user_id").Find(&session); r.Error != nil {
		log.Panic(r)
	}

	user.ID = session.UserID

	if r := database.DB.Model(&user).Find(&user); r.Error != nil {
		log.Panic(r)
	}

	id := ctx.Params("id")
	r := database.DB.Model(&users).Where("id = ?", id).Find(&users)

	if r.Error != nil {
		log.Panic(r)
	}

	if r.RowsAffected == 0 {
		return ctx.Status(http.StatusOK).JSON(response.Base{
			Message: config.RES_FINE,
			Data:    users,
		})
	}

	for _, item := range users {
		following := item.Following
		var resUser response.Users
		copier.CopyWithOption(&resUser, item, copier.Option{IgnoreEmpty: true})
		if len(following) == 0 {
			resUser.FollowStatus = "follow"
		} else {
			resUser.FollowStatus = "followed"
		}
		resUsers = append(resUsers, resUser)
	}

	sort.Slice(resUsers, func(i, j int) bool {
		return resUsers[i].FollowStatus > resUsers[j].FollowStatus
	})

	return ctx.Status(http.StatusOK).JSON(response.Base{
		Message: config.RES_FINE,
		Data:    resUsers,
	})

}

func UpdateUser(ctx *fiber.Ctx) error {
	var user model.User
	var form model.User
	var users []model.User
	var session model.Session

	// checking header and get claims
	claims, err := helper.GetTokenAndValidate(ctx)
	if err != nil {
		return err
	}

	// check user login
	session.ID = claims["id"].(string)
	if r := database.DB.Model(&model.Session{}).Select("user_id").Find(&session); r.Error != nil {
		log.Panic(r)
	}

	// check user data from table
	user.ID = session.UserID
	if r := database.DB.Model(&user).Find(&user); r.Error != nil {
		log.Panic(r)
	}

	// get data form
	if r := ctx.BodyParser(&form); r != nil {
		return ctx.Status(http.StatusBadRequest).JSON(response.Message{
			Message: config.RES_BAD_REQUEST,
		})
	}

	// query
	id := ctx.Params("id")
	r := database.DB.Model(&users).Where("id = ?", id).Updates(&model.User{
		Name:     form.Name,
		Username: form.Username,
	})

	// error for credential at jwt token
	// doesn't same with id at url param
	if r.Error != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(response.Message{Message: config.RES_USER_BAD_CREDETENTIALS})
	}
	if r.RowsAffected == 0 {
		return ctx.Status(http.StatusOK).JSON(response.Base{
			Message: config.RES_FINE,
			Data:    config.RES_BAD_REQUEST,
		})
	}

	// response
	return ctx.Status(http.StatusOK).JSON(response.Base{
		Message: config.RES_FINE,
		Data:    config.RES_USER_UPDATE,
	})

}

func DeleteUser(ctx *fiber.Ctx) error {
	var user model.User
	var users []model.User
	var session model.Session

	// checking header and get claims
	claims, err := helper.GetTokenAndValidate(ctx)
	if err != nil {
		return err
	}

	// check user login
	session.ID = claims["id"].(string)
	if r := database.DB.Model(&model.Session{}).Select("user_id").Find(&session); r.Error != nil {
		log.Panic(r)
	}

	// check user data from table
	user.ID = session.UserID
	if r := database.DB.Model(&user).Find(&user); r.Error != nil {
		log.Panic(r)
	}

	// query
	id := ctx.Params("id")
	r := database.DB.Model(&users).Where("id = ?", id).Delete(&user)

	// error for credential at jwt token
	// doesn't same with id at url param
	if r.Error != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(response.Message{Message: config.RES_USER_BAD_CREDETENTIALS})
	}

	if r.RowsAffected == 0 {
		return ctx.Status(http.StatusOK).JSON(response.Base{
			Message: config.RES_FINE,
			Data:    config.RES_BAD_REQUEST,
		})
	}

	return ctx.Status(http.StatusOK).JSON(response.Base{
		Message: config.RES_FINE,
		Data:    config.RES_USER_DELETE,
	})

}
