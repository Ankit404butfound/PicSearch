package services

import (
	"fmt"
	"os"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

func GenerateToken(user_id int) (string, error) {

	token_lifespan, err := strconv.Atoi(os.Getenv("TOKEN_HOUR_LIFESPAN"))

	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = user_id
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("API_SECRET")))

}

func ValidateToken(tokenString string) (*jwt.Token, error) {
	fmt.Println("Validating token:", tokenString)
	token := tokenString[len("Bearer "):]
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, nil
		}
		return []byte(os.Getenv("API_SECRET")), nil
	})
}

func ExtractUSerID(tokenString string) (int, error) {
	token, err := ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		fmt.Println("Claims[UserID]:", claims["user_id"])
		uid, err := strconv.Atoi(fmt.Sprintf("%v", claims["user_id"]))
		if err != nil {
			return 0, err
		}
		return uid, nil
	}
	return 0, nil
}
