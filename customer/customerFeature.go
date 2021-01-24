package customer

type CustomerFeature interface {
	Login(email string, password string) (bool, error)
	ChangePassword(email string, oldPassword string, newPassword string) error
	GetProfile(email string) (User, error)
	UpdateProfile(email string, name string) error
	Register(email string, password string, name string) error
}
