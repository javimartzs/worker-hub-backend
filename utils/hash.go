package utils

import "golang.org/x/crypto/bcrypt"

// Funcion para hashear las contraseñas de los usuarios
// ------------------------------------------------------------------
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// Funcion para comparar contraseñas en string con su hash
// ------------------------------------------------------------------
func CheckPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(hash)) == nil
}
