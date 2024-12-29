package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/javimartzs/worker-hub-backend/utils"
)

var revokedTokens = make(map[string]time.Time)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token invalido o nulo",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		if revocationTime, revoked := revokedTokens[tokenString]; revoked && time.Now().Before(revocationTime) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "El Token ha sido revocado",
			})
			c.Abort()
			return
		}

		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("id", (*claims)["id"])
		c.Set("role", (*claims)["role"])
		c.Next()
	}
}

// Funcion para revocar un token
func RevokeToken(tokenString string, duration time.Duration) {
	revokedTokens[tokenString] = time.Now().Add(duration)
}

// Limpieza periodica de tokens revocados
func StartTokenCleanup() {
	go func() {
		for {
			time.Sleep(time.Minute) // Limpiamos los token revocados
			now := time.Now()
			for token, expTime := range revokedTokens {
				if now.After(expTime) {
					delete(revokedTokens, token)
				}
			}
		}
	}()
}
