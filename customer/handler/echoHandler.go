package customer

import (
	"encoding/json"
	"fmt"
	"goworkshop2/customer"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type EchoHandler struct {
	Service customer.Feature
	app     *echo.Echo
	Port    string
}

func (handler *EchoHandler) Start() {
	app := echo.New()
	handler.app = app

	// app.Use(middleware.Logger())
	// app.Use(middleware.Recover())

	app.POST("/user", handler.register())
	app.POST("/user/login", handler.login())
	app.PUT("/user/password", handler.changePassword(), middleware.JWT([]byte("secret")))
	app.PUT("/user/edit", handler.editProfile(), middleware.JWT([]byte("secret")))
	app.GET("/user", handler.getProfile(), middleware.JWT([]byte("secret")))

	app.Logger.Fatal(app.Start(":" + handler.Port))
}

func (handler *EchoHandler) register() func(context echo.Context) error {
	service := handler.Service
	return func(context echo.Context) error {
		var data customer.User

		body, _ := ioutil.ReadAll(context.Request().Body)
		json.Unmarshal(body, &data)

		err := service.Register(data.Email, data.Password, data.Name)
		if err != nil {
			return context.JSON(200, "fail email already exist")
		}
		return context.JSON(200, "success")
	}
}

func (handler *EchoHandler) login() func(context echo.Context) error {
	service := handler.Service

	return func(context echo.Context) error {
		var data = struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{}

		body, _ := ioutil.ReadAll(context.Request().Body)
		json.Unmarshal(body, &data)
		fmt.Println("login account : ", data.Email, " ,password : ", data.Password)

		result, err := service.Login(data.Email, data.Password)
		fmt.Println("login result ", result, " || error ", err)
		if err != nil || !result {
			return echo.ErrUnauthorized
		}

		return createToken(data.Email, context)
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
		}
		return context.JSON(200, "success")
	}
}

func (handler *EchoHandler) getProfile() func(context echo.Context) error {
	service := handler.Service

	return func(context echo.Context) error {
		value := context.Param("email")
		user, err := service.GetProfile(value)
		if err != nil {
			return err
		}
		return context.JSON(200, user)
	}
}

func (handler *EchoHandler) editProfile() func(context echo.Context) error {
	service := handler.Service

	return func(context echo.Context) error {
		var data = struct {
			Email string `tag:"email"`
			Name  string `json:"name"`
		}{}

		body, _ := ioutil.ReadAll(context.Request().Body)
		json.Unmarshal(body, &data)

		err := service.UpdateProfile(data.Email, data.Name)
		if err != nil {
			return err
		}
		return context.JSON(200, "success")
	}
}

func createToken(email string, context echo.Context) error {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = email
	claims["admin"] = false
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return context.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}
