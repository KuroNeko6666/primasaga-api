package helper

import (
	"net/http"
	"strings"

	"github.com/KuroNeko6666/prima-api/config"
	"github.com/KuroNeko6666/prima-api/interfaces/response"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func GetTokenAndValidate(ctx *fiber.Ctx) (jwt.MapClaims, error) {

	token := ctx.Get("Authorization")
	if token == "" {
		return nil, ctx.Status(http.StatusUnauthorized).JSON(response.Message{Message: config.RES_MISSING_AUTH})
	}
	token = strings.ReplaceAll(token, "Bearer ", "")

	claims, _ := ExtractClaims(token, config.SECRET_AUTH)

	return claims, nil
}
