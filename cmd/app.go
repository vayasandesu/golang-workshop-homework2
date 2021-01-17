package main

import (
	"goworkshop2/game"
	"goworkshop2/handler"
	"goworkshop2/service"
	"goworkshop2/storage"
)

func main() {
	//1. create game type mongo or local(memory)
	//2. create handler gin/fiber/echo
	//3. Start server

	server := CreateEchoServer(CreateLocalGame())
	server.StartServer()
}

func CreateLocalGame() *service.LocalGame {
	return &service.LocalGame{}
}

func CreateMongoGame() *service.MongoGame {
	storage := storage.MongoDb{
		Config: storage.MongoDbConfiguration{
			ConnectionTimeout:        5,
			ConnectionStringTemplate: "mongodb://%s:%s@%s",
		},
	}
	storage.InitializeResource()

	return &service.MongoGame{
		Resource: storage.Resource,
	}
}

func CreateGinServer(feature game.IFeature) handler.GinServer {
	return handler.GinServer{
		Game: feature,
		Port: "8080",
	}
}

func CreateFiberServer(feature game.IFeature) handler.FiberServer {
	return handler.FiberServer{
		Game: feature,
		Port: "8080",
	}
}

func CreateEchoServer(feature game.IFeature) handler.EchoServer {
	return handler.EchoServer{
		Game: feature,
		Port: "8080",
	}
}
