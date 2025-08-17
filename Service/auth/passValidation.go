package auth

type IPassValidation interface {
	IsValidPassword(password string) bool
}
