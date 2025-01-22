package utils

import (
	"regexp"
	"unicode"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// ValidateEmail verifica se o email é válido
func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// ValidatePassword verifica se a senha é válida
func ValidatePassword(password string) (bool, string) {
	if len(password) < 8 {
		return false, "a senha deve ter no mínimo 8 caracteres"
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return false, "a senha deve conter pelo menos uma letra maiúscula"
	}
	if !hasLower {
		return false, "a senha deve conter pelo menos uma letra minúscula"
	}
	if !hasNumber {
		return false, "a senha deve conter pelo menos um número"
	}
	if !hasSpecial {
		return false, "a senha deve conter pelo menos um caractere especial"
	}

	return true, ""
}

// ValidateUsername verifica se o username é válido
func ValidateUsername(username string) (bool, string) {
	if len(username) < 3 {
		return false, "o nome de usuário deve ter no mínimo 3 caracteres"
	}
	if len(username) > 30 {
		return false, "o nome de usuário deve ter no máximo 30 caracteres"
	}

	validUsername := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !validUsername.MatchString(username) {
		return false, "username deve conter apenas letras, números, underscores e hifens"
	}

	return true, ""
}
