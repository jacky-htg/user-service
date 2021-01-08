package token

import (
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// MyCustomeClaims struct
type MyCustomeClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

var mySigningKey = []byte(os.Getenv("TOKEN_SALT"))

// ValidateToken func
func ValidateToken(myToken string) (bool, string) {
	token, err := jwt.ParseWithClaims(myToken, &MyCustomeClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(mySigningKey), nil
	})

	if err != nil {
		return false, ""
	}

	claims := token.Claims.(*MyCustomeClaims)
	return token.Valid, claims.Email
}

// ClaimToken func
func ClaimToken(email string) (string, error) {
	claims := MyCustomeClaims{
		email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 5).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret
	return token.SignedString(mySigningKey)
}
