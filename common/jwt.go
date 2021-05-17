package common

import (
	"gin.vue.demo/ginessential/model"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtkey = []byte("a_secret_crect")

type Claims struct {
	UserId uint
	jwt.StandardClaims
}

func ReleaseToken(user model.User)  (string, error){
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserId : user.ID,
		StandardClaims : jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt: time.Now().Unix(),
			Issuer: "sweeteriri",
			Subject: "user token",
		},

	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtkey)

	if err != nil {
		return "", err
	}

	return tokenString, nil

}