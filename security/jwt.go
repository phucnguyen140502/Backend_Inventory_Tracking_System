package security

import (
	"backend/models"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const SECRET_KEY = "dgrijalvabjkbjlbkbkbkdwederwewgegege"

type JwtCustomClaims struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func GenToken(user models.Users) (string, error) {
	claims := JwtCustomClaims{
		Name:     user.FullName,
		Email:    user.Email,
		Password: user.Password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	results, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}
	return results, nil
}
