package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Context defines which account type is being used
type Context int

const (
	expirationTimeToken = 240
)

// Auth struct
type Auth struct {
	secret string
}

// RequestAuth is used to request a token
type RequestAuth struct {
	ID    string
	Email string
}

// New returns a new auth service
func New() *Auth {
	randomString := fmt.Sprintf("%s--%s", time.Now().String(), time.Now().String())
	return &Auth{
		secret: randomString,
	}
}

// GetToken returns a new token
func (a *Auth) GetToken(requestAuth *RequestAuth) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    requestAuth.ID,
		"email": requestAuth.Email,
		"exp":   time.Now().Add(time.Hour * expirationTimeToken).Unix(),
	})
	return token.SignedString([]byte(a.secret))
}

// IsValid checks if token is valid
func (a *Auth) IsValid(auhtorization string) bool {
	token, err := jwt.Parse(auhtorization, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error on validating auth token")
		}
		return []byte(a.secret), nil
	})
	if err != nil {
		return false
	}
	return token.Valid
}

// GetClaims transforms the token string into a map with its claims
func GetClaims(auhtorization string) (map[string]string, error) {
	token, _ := jwt.Parse(auhtorization, func(token *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		id := claims["id"].(string)
		email := claims["email"].(string)
		claims := make(map[string]string)
		claims["id"] = id
		claims["email"] = email
		return claims, nil
	}
	return nil, fmt.Errorf("could not get claims")

}
