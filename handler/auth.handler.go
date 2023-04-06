package handler

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/KuroNeko6666/prima-api/config"
	"github.com/KuroNeko6666/prima-api/database"
	"github.com/KuroNeko6666/prima-api/database/model"
	"github.com/KuroNeko6666/prima-api/helper"
	"github.com/KuroNeko6666/prima-api/interfaces/form"
	"github.com/KuroNeko6666/prima-api/interfaces/response"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
)

func Login(ctx *fiber.Ctx) error {
	var form form.Login
	var user model.User
	var session model.Session

	ipAddress := ctx.IP()
	userAggent := ctx.GetReqHeaders()["User-Agent"]

	if r := ctx.BodyParser(&form); r != nil {
		return ctx.Status(http.StatusBadRequest).JSON(response.Message{
			Message: config.RES_BAD_REQUEST,
		})
	}

	if r := helper.Validate().Struct(&form); r != nil {
		return ctx.Status(http.StatusBadRequest).JSON(response.Message{
			Message: config.RES_BAD_REQUEST,
		})
	}

	db := database.DB.Model(&model.User{}).
		Preload("Credential").Preload("Followers").Preload("Following").
		Where("username = ?", form.Username).Or("email = ?", form.Username).Find(&user)

	if db.Error != nil {
		return ctx.Status(http.StatusConflict).JSON(response.Message{
			Message: db.Error.Error(),
		})
	}

	if db.RowsAffected == 0 {
		return ctx.Status(http.StatusUnauthorized).JSON(response.Message{
			Message: config.RES_LOGIN_FAIL,
		})
	}

	match := helper.CompareHash(form.Password, user.Credential.Password)

	if !match {
		return ctx.Status(http.StatusUnauthorized).JSON(response.Message{
			Message: config.RES_LOGIN_FAIL,
		})
	}

	// if !user.EmailVerified {
	// 	return ctx.Status(http.StatusUnauthorized).JSON(response.Message{
	// 		Message: config.RES_LOGIN_FAIL,
	// 	})
	// }

	db = database.DB.Model(&model.Session{}).Where("ip_address = ?", ipAddress).Where("user_agent = ?", userAggent).Find(&session)

	if db.Error != nil {
		return ctx.Status(http.StatusConflict).JSON(response.Message{
			Message: db.Error.Error(),
		})
	}

	if db.RowsAffected > 0 {
		session.ExpiredAt = time.Now().Add(time.Hour * 72)
		db := db.Model(&model.Session{}).Updates(&session)
		if db.Error != nil {
			return ctx.Status(http.StatusConflict).JSON(response.Message{
				Message: db.Error.Error(),
			})
		}

		ctx.Cookie(&fiber.Cookie{
			Name:  "session_id",
			Value: session.ID,
		})

		return ctx.Status(http.StatusOK).JSON(response.Base{
			Message: "has logged in",
			Data:    user,
		})
	}

	session.UserID = user.ID
	session.IpAddress = ipAddress
	session.UserAgent = userAggent

	sess := database.DB.Model(&model.Session{}).Create(&session)

	if sess.Error != nil {
		return ctx.Status(http.StatusConflict).JSON(response.Message{
			Message: sess.Error.Error(),
		})
	}

	ctx.Cookie(&fiber.Cookie{
		Name:  "session_id",
		Value: session.ID,
	})

	return ctx.Status(http.StatusOK).JSON(response.Base{
		Message: config.RES_LOGIN,
		Data:    user,
	})

}

func Register(ctx *fiber.Ctx) error {
	var form form.Register
	var user model.User
	var credential model.UserCredential

	if r := ctx.BodyParser(&form); r != nil {
		return ctx.Status(http.StatusBadRequest).JSON(response.Message{
			Message: config.RES_BAD_REQUEST,
		})
	}

	if r := helper.Validate().Struct(&form); r != nil {
		return ctx.Status(http.StatusBadRequest).JSON(response.Message{
			Message: config.RES_BAD_REQUEST,
		})
	}

	if r := copier.CopyWithOption(&user, &form, copier.Option{IgnoreEmpty: true}); r != nil {
		log.Panic(r.Error())
	}

	hash, r := helper.CreateHash(form.Password)

	if r != nil {
		log.Panic(r)
	}

	credential.Password = hash
	credential.UserID = user.ID
	user.Credential = credential

	db := database.DB.Model(&model.User{}).Create(&user)

	if db.Error != nil {
		return ctx.Status(http.StatusConflict).JSON(response.Message{
			Message: db.Error.Error(),
		})
	}

	token, r := helper.GenerateToken(user, config.SECRET_VERIFIY)

	if r != nil {
		log.Panic(r.Error())
	}

	if r := helper.SendEmailVerification(user.Email, token); r != nil {
		log.Panic(r.Error())
	}

	return ctx.Status(http.StatusOK).JSON(response.Message{
		Message: config.RES_REGISTER,
	})

}

func Logout(ctx *fiber.Ctx) error {
	var session model.Session

	expiredAt := time.Now().Add(time.Second * -1)
	sessionID := ctx.Cookies("session_id", "")

	db := database.DB.Model(&model.Session{}).Where("id = ?", sessionID).Delete(&session)

	if db.Error != nil {
		return ctx.Status(http.StatusConflict).JSON(response.Message{
			Message: db.Error.Error(),
		})
	}

	ctx.Cookie(&fiber.Cookie{
		Name:    "session_id",
		Value:   sessionID,
		Expires: expiredAt,
	})

	return ctx.Status(http.StatusOK).JSON(response.Message{
		Message: config.RES_LOGOUT,
	})
}

func EmailVerififation(ctx *fiber.Ctx) error {
	var user model.User
	token := ctx.Query("token")

	claims, r := helper.ExtractClaims(token, config.SECRET_VERIFIY)

	if r != nil {
		log.Panic(r.Error())
	}

	db := database.DB.Model(&model.User{}).Where("id = ?", claims["id"]).Find(&user)

	if db.Error != nil {
		return ctx.Status(http.StatusConflict).JSON(response.Message{
			Message: db.Error.Error(),
		})
	}

	if db.RowsAffected == 0 {
		return ctx.Status(http.StatusUnauthorized).JSON(response.Message{
			Message: config.RES_INVALID_TOKEN,
		})
	}

	if user.EmailVerified {
		return ctx.Status(http.StatusUnauthorized).JSON(response.Message{
			Message: config.RES_INVALID_TOKEN,
		})
	}

	user.EmailVerified = true

	db = db.Updates(&user)

	if db.Error != nil {
		return ctx.Status(http.StatusConflict).JSON(response.Message{
			Message: db.Error.Error(),
		})
	}

	if db.RowsAffected == 0 {
		return ctx.Status(http.StatusUnauthorized).JSON(response.Message{
			Message: config.RES_INVALID_TOKEN,
		})
	}

	return ctx.Status(http.StatusOK).JSON(response.Message{
		Message: config.RES_VERIFIED,
	})
}

func Check(ctx *fiber.Ctx) error {
	var session model.Session
	sessionID := ctx.Cookies("session_id", "")

	db := database.DB.Model(&model.Session{}).Where("id", sessionID).Find(&session)

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
		db := db.Delete(&session)

		if db.Error != nil {
			return ctx.Status(http.StatusConflict).JSON(response.Message{
				Message: db.Error.Error(),
			})
		}

		ctx.Cookie(&fiber.Cookie{
			Name:    "session_id",
			Value:   session.UserID,
			Expires: time.Now().Add(time.Second - 2000000000),
		})

		return ctx.Status(http.StatusUnauthorized).JSON(response.Message{
			Message: config.RES_UNAUTHORIZED,
		})
	}

	return ctx.Status(http.StatusOK).JSON(response.Message{
		Message: config.RES_FINE,
	})

}
