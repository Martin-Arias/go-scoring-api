package auth

type AuthUserData struct {
	ID           string
	Username     string
	PasswordHash string
	IsAdmin      bool
}
