package customer

import (
	"encoding/json"
	"goworkshop2/customer"
	"io/ioutil"

	"github.com/labstack/echo/v4"
)

type EchoHandler struct {
	Service customer.Feature
	app     *echo.Echo
	Port    string
}

func (handler *EchoHandler) Start() {
	app := echo.New()
	handler.app = app

	app.POST("/user", handler.register())
	app.GET("/user", handler.getProfile())
	app.POST("/user/login", handler.login())
	app.PUT("/user/password", handler.changePassword())

	app.Start(":" + handler.Port)
}

func (handler *EchoHandler) register() func(context echo.Context) error {
	service := handler.Service
	return func(context echo.Context) error {
		var data customer.User

		body, _ := ioutil.ReadAll(context.Request().Body)
		json.Unmarshal(body, &data)

		err := service.Register(data.Email, data.Password, data.Name)
		if err != nil {
			return err
		} else {
			return context.JSON(200, "success")
		}
	}
}

func (handler *EchoHandler) login() func(context echo.Context) error {
	service := handler.Service

	return func(context echo.Context) error {
		var data = struct {
			email    string
			password string
		}{}

		body, _ := ioutil.ReadAll(context.Request().Body)
		json.Unmarshal(body, &data)

		result, err := service.Login(data.email, data.password)
		if err == nil && result {
			return context.JSON(200, "success")
		} else {
			return err
		}

	}
}

func (handler *EchoHandler) changePassword() func(context echo.Context) error {
	service := handler.Service

	return func(context echo.Context) error {
		var data = struct {
			Email       string `tag:"email"`
			OldPassword string `json:"password"`
			NewPassword string `json:"newpassword"`
		}{}

		body, _ := ioutil.ReadAll(context.Request().Body)
		json.Unmarshal(body, &data)

		err := service.ChangePassword(data.Email, data.OldPassword, data.NewPassword)
		if err != nil {
			return err
		} else {
			return context.JSON(200, "success")
		}
	}
}

func (handler *EchoHandler) getProfile() func(context echo.Context) error {
	service := handler.Service

	return func(context echo.Context) error {
		value := context.Param("email")
		user, err := service.GetProfile(value)
		if err != nil {
			return err
		} else {
			return context.JSON(200, user)
		}
	}
}
