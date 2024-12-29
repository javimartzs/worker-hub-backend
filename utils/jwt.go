package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/javimartzs/worker-hub-backend/config"
)

var JwtKey = []byte(config.Env.JwtKey)

// Funcion que genera los Json Web Tokens
// ------------------------------------------------------------------
func GenerateJWT(role, id, storeID string) (string, error) {
	claims := jwt.MapClaims{
		"role":     role,
		"id":       id,
		"store_id": storeID, // Añadimos el store_id como parte de los claims
		"exp":      time.Now().Add(time.Hour * 18).Unix(),
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
			return nil, errors.New("the session has expired")
		}
		return nil, errors.New("invalid token")
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Acceder al store_id (si existe)
	if storeID, ok := (*claims)["store_id"].(string); ok && storeID != "" {
		fmt.Println("Store ID:", storeID)
	}

	return claims, nil
}
