package services

import (
	"time"

	"github.com/gethoopp/go_traver.git/models"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("my_secret_key")

func CreateToken() (string, error) {

	var log models.UserData

	expiredTime := time.Now().Add(10 * time.Minute)

	claims := &models.ClaimsData{
		NamaUser: log.Emailuser,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func verifyToken(tokenstring string) {}
