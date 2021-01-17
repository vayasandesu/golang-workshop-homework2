package handler

import (
	"encoding/json"
	"goworkshop2/game"

	"github.com/gofiber/fiber/v2"
)

type FiberServer struct {
	Game game.IFeature
	App  *fiber.App
	Port string
}

func (server *FiberServer) StartServer() {
	app := fiber.New()
	server.App = app

	app.Get("/players", server.list())
	app.Post("/players", server.join())
	app.Listen(":" + server.Port)
}

func (server *FiberServer) response(fiber *fiber.Ctx, data interface{}, err error) error {
	if err != nil {
		return err
	} else {
		return fiber.JSON(data)
	}
}

func (server *FiberServer) list() func(fiber *fiber.Ctx) error {
	return func(fiber *fiber.Ctx) error {
		data, err := server.Game.List()
		return server.response(fiber, data, err)
	}
}

func (server *FiberServer) join() func(fiber *fiber.Ctx) error {
	return func(fiber *fiber.Ctx) error {
		var data = struct {
			Name string
		}{
			Name: "",
		}

		body := fiber.Body()
		json.Unmarshal(body, &data)

		character, err := server.Game.Join(data.Name)
		return server.response(fiber, character, err)
	}
}
