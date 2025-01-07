package services

import (
	"errors"
	"strings"
	"time"

	"github.com/javimartzs/worker-hub-backend/middlewares"
	"github.com/javimartzs/worker-hub-backend/repositories"
	"github.com/javimartzs/worker-hub-backend/utils"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// Login - Panel de admininstrador
// --------------------------------------------------------------------
func (s *AuthService) LoginAdmin(username, password string) (string, error) {

	// Buscamos el usuario por el username
	user, err := s.userRepo.FindUserByUsername(nil, username)
	if err != nil || user.Role != "admin" {
		return "", errors.New("credenciales del admin invalidas")
	}

	// Comprobamos que el password sea correcto
	if !utils.CheckPassword(user.Password, password) {
		return "", errors.New("credenciales del admin invalidas")
	}

	// Generamos el token JWT
	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return "", errors.New("error al generar el token de auth")
	}

	return token, nil
}

// Login - Panel de tiendas
// --------------------------------------------------------------------
func (s *AdminService) LoginStore(username, password string) (string, error) {

	// Buscamos el usuario por el username
	user, err := s.userRepo.FindUserByUsername(nil, username)
	if err != nil || user.Role != "store" {
		return "", errors.New("credenciales del usuario invalidas ")
	}

	// Comprobamos que el password sea correcto
	if !utils.CheckPassword(user.Password, password) {
		return "", errors.New("credenciales del usuario invalidas")
	}

	// Generamos el token JWT
	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return "", errors.New("error al generar el token de auth")
	}

	return token, nil
}

// Logout
// --------------------------------------------------------------------
func (s *AuthService) LogoutAdmin(token string) error {
	trimmedToken := strings.TrimPrefix(token, "Bearer ")
	middlewares.RevokeToken(trimmedToken, 18*time.Hour)
	return nil
}
