package authService

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)



func GenerateToken(email string)(string,error){
claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
	Issuer: email,
	ExpiresAt: time.Now().Add(time.Hour * 8).Unix(),
})

token, err := claims.SignedString([]byte(os.Getenv("SECRET")))
if err != nil {
	fmt.Println("Error: ", err)
	return "", err
}
return token, nil
}

func GenerateRefreshToken(email string)(string,error){
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer: email,
		ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
	})
	
	token, err := claims.SignedString([]byte(os.Getenv("REFRESH")))
	if err != nil {
		fmt.Println("Error: ", err)
		return "", err
	}
	return token, nil
}