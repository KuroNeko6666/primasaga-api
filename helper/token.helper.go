package helper

import (
	"errors"
	"time"

	"github.com/KuroNeko6666/prima-api/database/model"
	"github.com/golang-jwt/jwt/v4"
)

func ExtractClaims(tokenStr string, secretStr string) (jwt.MapClaims, error) {
	secret := []byte(secretStr)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, errors.New("error extract token")
	}
}

func GenerateToken(user model.User, secret string) (string, error) {

	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := rawToken.SignedString([]byte(secret))

	if err != nil {
		return "", err
	}

	return token, err

}
