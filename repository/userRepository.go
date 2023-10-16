package repository

type userRepository struct {
}

type UserReaderWriter interface {
	GetUser(string username) &User
}
