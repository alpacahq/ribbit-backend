package secret

// Service is the interface to our secret service
type Service interface {
	HashPassword(password string) string
	HashMatchesPassword(hash, password string) bool
	HashRandomPassword() (string, error)
}
