package handler

import (
	"encoding/json"
	"goworkshop2/game"
	"io/ioutil"

	"github.com/labstack/echo/v4"
)

type EchoServer struct {
	Game game.IFeature
	App  *echo.Echo
	Port string
}

func (server *EchoServer) StartServer() {
	app := echo.New()
	server.App = app

	app.GET("/players", server.list())
	app.POST("/players", server.join())
	app.Start(":" + server.Port)
}

func (server *EchoServer) response(context echo.Context, data interface{}, err error) error {
	if err != nil {
		context.Error(err)
		return err
	} else {
		return context.JSON(200, data)
	}
}

func (server *EchoServer) list() func(context echo.Context) error {
	return func(context echo.Context) error {
		data, err := server.Game.List()
		return server.response(context, data, err)
	}
}

func (server *EchoServer) join() func(context echo.Context) error {
	return func(context echo.Context) error {
		var data = struct {
			Name string
		}{
			Name: "",
		}

		body, _ := ioutil.ReadAll(context.Request().Body)
		json.Unmarshal(body, &data)

		character, err := server.Game.Join(data.Name)
		return server.response(context, character, err)
	}
}
