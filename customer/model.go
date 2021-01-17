package customer

type User struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name" bson:"name"`
	Email    string `tag:"email"`
	Password string `json:"password" bson:"password"`
}
