package customer

type Feature interface {
	Login(email string, password string) (bool, error)
	ChangePassword(email string, oldPassword string, newPassword string) error
	GetProfile(email string) (User, error)
	Register(email string, password string, name string) error
}
