package customer

type User struct {
	Name     string `json:"name" bson:"name"`
	Email    string `bson:"email" json:"email" tag:"email"`
	Password string `json:"password" bson:"password"`
}
