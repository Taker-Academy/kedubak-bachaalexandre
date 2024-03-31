package encoding

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type CustomClaims struct {
	UserID string `json:"userId"`
	jwt.StandardClaims
}

func VerifyJWT(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	if time.Now().Unix() > claims.ExpiresAt {
		return nil, jwt.ValidationError{Inner: err, Errors: jwt.ValidationErrorExpired}
	}
	return claims, nil
}
