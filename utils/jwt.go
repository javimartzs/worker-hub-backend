package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/javimartzs/worker-hub-backend/config"
)

var JwtKey = []byte(config.Env.JwtKey)

// Funcion que genera los Json Web Tokens
// ------------------------------------------------------------------
func GenerateJWT(id, role string) (string, error) {

	claims := jwt.MapClaims{
		"id":   id,
		"role": role,
		"exp":  time.Now().Add(time.Hour * 18).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JwtKey)
}

// Funcion para validar el token y retornar sus claims
// ------------------------------------------------------------------
func ValidateJWT(tokenString string) (*jwt.MapClaims, error) {

	claims := &jwt.MapClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("la sesión a expirado")
		}
		return nil, errors.New("token invalido")
	}

	if !token.Valid {
		return nil, errors.New("token invalido")
	}

	return claims, nil
}
