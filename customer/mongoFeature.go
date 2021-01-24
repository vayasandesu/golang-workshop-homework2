package customer

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerService struct {
	Resource *mongo.Database
}

func initContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	return ctx, cancel
}

func (feature *CustomerService) Login(email string, password string) (bool, error) {
	user, err := feature.getUser(email)
	if err == nil {
		return (user.Password == password), err
	} else {
		return false, err
	}
}

func (feature *CustomerService) ChangePassword(email string, oldPassword string, newPassword string) error {
	user, err := feature.getUser(email)
	if err != nil {
		return err
	}

	if user.Password != oldPassword {
		return errors.New("password not match")
	}

	ctx, _ := initContext()
	collection := feature.Resource.Collection("Users")
	update := bson.M{
		"$set": bson.M{
			"password": newPassword,
		},
	}

	_, err = collection.UpdateOne(ctx,
		bson.M{"email": user.Email},
		update)

	return err
}

func (feature *CustomerService) GetProfile(email string) (User, error) {
	user, err := feature.getUser(email)
	output := User{
		Email: user.Email,
		Name:  user.Name,
	}

	return output, err
}

func (feature *CustomerService) UpdateProfile(email string, name string) error {
	user, err := feature.getUser(email)
	if err != nil {
		return err
	}

	ctx, _ := initContext()
	collection := feature.Resource.Collection("Users")
	update := bson.M{
		"$set": bson.M{
			"name": name,
		},
	}
	_, err = collection.UpdateOne(ctx,
		bson.M{"email": user.Email},
		update,
	)

	return err
}

func (feature *CustomerService) Register(email string, password string, name string) error {
	if feature.isExist(email) {
		return errors.New("Email already exist")
	}

	ctx, _ := initContext()
	collection := feature.Resource.Collection("Users")
	user := User{
		Email:    email,
		Name:     name,
		Password: password,
	}

	_, err := collection.InsertOne(ctx, user)

	return err
}

func (feature *CustomerService) isExist(email string) bool {
	user, err := feature.getUser(email)
	if err == nil && user.Email == email {
		return true
	} else {
		return false
	}
}

func (feature *CustomerService) getUser(email string) (User, error) {
	var data User
	ctx, _ := initContext()
	collection := feature.Resource.Collection("Users")
	c := collection.FindOne(ctx, bson.M{"email": email})
	c.Decode(&data)

	return data, c.Err()
}
