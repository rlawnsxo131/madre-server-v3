package account

type UserUseCase interface {
	Create(u *User) (string, error)
	FindOneById(id string) (*User, error)
	FindOneByUsername(username string) (*User, error)
}
