package authService

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

/********************************************
Generates a JWT token for every admin issuer.
********************************************/
func GenerateToken(email string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    email,
		ExpiresAt: time.Now().Add(time.Hour * 336).Unix(),
	})

	token, err := claims.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		fmt.Println("Error: ", err)
		return "", err
	}
	return token, nil
}

/**********************************************
Refreshes the JWT token for every admin issuer.
**********************************************/
func GenerateRefreshToken(email string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    email,
		ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
	})

	token, err := claims.SignedString([]byte(os.Getenv("REFRESH")))
	if err != nil {
		fmt.Println("Error: ", err)
		return "", err
	}
	return token, nil
}

/***********************************************************
Checks for the authenticity of the JWT provided by the user.
***********************************************************/
func AuthenticateToken(token string, secretKey string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv(secretKey)), nil
	})
}
