package handler

import (
	"goworkshop2/customer"
)

type LocalCustomerFeature struct {
	Collections []customer.User
}

func (feature *LocalCustomerFeature) Login(email string, password string) (bool, error) {
	return true, nil
}

func (feature *LocalCustomerFeature) ChangePassword(email string, oldPassword string, newPassword string) error {
	return nil
}

func (feature *LocalCustomerFeature) GetProfile(email string) (customer.User, error) {
	user := customer.User{
		Email: email,
		Name:  email,
	}

	return user, nil
}

func (feature *LocalCustomerFeature) UpdateProfile(email string, name string) error {
	return nil
}

func (feature *LocalCustomerFeature) Register(email string, password string, name string) error {
	return nil
}
