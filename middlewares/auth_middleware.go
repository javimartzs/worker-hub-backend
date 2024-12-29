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

		// Leemos el encabezado Authorization
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Missing or invalid token",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Verificamos si el token ha sido revocado
		if revocationTime, revoked := revokedTokens[tokenString]; revoked && time.Now().Before(revocationTime) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token has been revoked",
			})
			c.Abort()
			return
		}

		// Validamos el token
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		// Pasamos los datos del usuario al contexto
		c.Set("id", (*claims)["id"])
		c.Set("role", (*claims)["role"])
		if storeID, ok := (*claims)["store_id"].(string); ok && storeID != "" {
			c.Set("store_id", storeID)
		}
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
