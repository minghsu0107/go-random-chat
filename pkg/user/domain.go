package user

type User struct {
	ID       uint64
	Email    string
	Name     string
	AuthType AuthType
}

type AuthType string

const (
	LocalAuth  AuthType = "local"
	GoogleAuth AuthType = "google"
)
