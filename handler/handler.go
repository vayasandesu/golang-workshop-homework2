package handler

import (
	"encoding/json"
	"fmt"
	"goworkshop2/customer"
	"io/ioutil"
	"net/http"
	"time"

	_ "goworkshop2/docs"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type EchoHandler struct {
	Service customer.ServiceInterface
	app     *echo.Echo
	Port    string
}

func (handler *EchoHandler) Start() {
	app := echo.New()
	handler.app = app

	app.GET("/swagger/*any", echoSwagger.WrapHandler)

	app.POST("/user", handler.register())
	app.POST("/user/login", handler.login())
	app.PUT("/user/password", handler.changePassword(), middleware.JWT([]byte("secret")))
	app.PUT("/user/edit", handler.editProfile(), middleware.JWT([]byte("secret")))
	app.GET("/user", handler.getProfile(), middleware.JWT([]byte("secret")))

	app.Logger.Fatal(app.Start(":" + handler.Port))
}

// Register godoc
// @Summary Register
// @Description Register new account
// @Accept  json
// @Produce  json
// @Param email body string true "email"
// @Param password body string true "password"
// @Param name body string true "name of user"
// @Success 200 {string} string "success"
// @Failure 200 {string} string "fail email already exist"
// @Failure 400,404 {object} object httputil.HTTPError
// @Failure 500 {object} object httputil.HTTPError
// @Router /user/register [post]
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

// Login godoc
// @Summary Login
// @Description Login and get token
// @Accept  json
// @Produce  json
// @Param email body string true "email"
// @Param password body string true "password"
// @Success 200 {object} string "{"token" : value}"
// @Failure 401 {string} string "Unauthorized"
// @Failure 400,404 {object} object httputil.HTTPError
// @Failure 500 {object} object httputil.HTTPError
// @Router /user/login [post]
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

		token, err := createToken(data.Email, context)
		if err != nil {
			return err
		}

		return context.JSON(http.StatusOK, token)
	}
}

func (handler *EchoHandler) changePassword() func(context echo.Context) error {
	service := handler.Service

	return func(context echo.Context) error {
		email, err := getTokenEmail(context)
		if err != nil {
			return echo.ErrUnauthorized
		}

		var data = struct {
			OldPassword string `json:"password"`
			NewPassword string `json:"newpassword"`
		}{}

		body, _ := ioutil.ReadAll(context.Request().Body)
		json.Unmarshal(body, &data)

		err = service.ChangePassword(email, data.OldPassword, data.NewPassword)
		if err != nil {
			return err
		}
		return context.JSON(200, "success")
	}
}

func (handler *EchoHandler) getProfile() func(context echo.Context) error {
	service := handler.Service

	return func(context echo.Context) error {
		email, err := getTokenEmail(context)
		if err != nil {
			return echo.ErrUnauthorized
		}

		user, err := service.GetProfile(email)
		if err != nil {
			return err
		}
		return context.JSON(200, user)
	}
}

func (handler *EchoHandler) editProfile() func(context echo.Context) error {
	service := handler.Service

	return func(context echo.Context) error {
		email, err := getTokenEmail(context)
		if err != nil {
			return echo.ErrUnauthorized
		}

		var data = struct {
			Name string `json:"name"`
		}{}

		body, _ := ioutil.ReadAll(context.Request().Body)
		json.Unmarshal(body, &data)

		err = service.UpdateProfile(email, data.Name)
		if err != nil {
			return err
		}
		return context.JSON(200, "success")
	}
}

func getTokenEmail(context echo.Context) (string, error) {
	temp := context.Get("user")
	if temp == nil {
		return "", echo.ErrUnauthorized
	}
	user := temp.(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	return claims["email"].(string), nil
}

func createToken(email string, context echo.Context) (map[string]string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return map[string]string{}, err
	}

	return map[string]string{
		"token": t,
	}, nil
}
