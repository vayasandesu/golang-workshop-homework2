package handler

import (
	"encoding/json"
	"goworkshop2/game"
	"io/ioutil"

	"github.com/gin-gonic/gin"
)

type GinServer struct {
	Game game.IFeature
	App  *gin.Engine
	Port string
}

func (server *GinServer) StartServer() {
	app := gin.New()
	server.App = app

	app.GET("/players", server.list())
	app.POST("/players", server.join())
	app.Run(":" + server.Port)
}

func (server *GinServer) response(context *gin.Context, data interface{}, err error) {
	if err != nil {
		context.Error(err)
	} else {
		context.JSON(200, data)
	}
}

func (server *GinServer) list() func(context *gin.Context) {
	return func(context *gin.Context) {
		data, err := server.Game.List()
		server.response(context, data, err)
	}
}

func (server *GinServer) join() func(context *gin.Context) {
	return func(context *gin.Context) {
		var data = struct {
			Name string
		}{
			Name: "",
		}

		body, _ := ioutil.ReadAll(context.Request.Body)
		json.Unmarshal(body, &data)

		character, err := server.Game.Join(data.Name)
		server.response(context, character, err)
	}
}
