package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/javimartzs/worker-hub-backend/services"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Handler para el login del panel de administrador
// --------------------------------------------------------------------
func (h *AuthHandler) LoginAdmin(c *gin.Context) {

	// Recogemos los datos del body
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// Decodificamos el body
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request",
		})
		return
	}

	// Llamamos al servicio de autenticación
	token, err := h.authService.LoginAdmin(request.Username, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Devolvemos el token
	c.JSON(http.StatusOK, gin.H{"token": token})
}
