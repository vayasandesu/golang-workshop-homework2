package main

import (
	"goworkshop2/customer"
	"goworkshop2/handler"
	"goworkshop2/storage"
)

// @title Goworkshop example api
// @version 1.0
// @description sample swagger with some api XD
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @BasePath /
func main() {
	config := storage.MongoDbConfiguration{
		ConnectionTimeout:        5,
		ConnectionStringTemplate: "mongodb://%s:%s@%s",
	}
	db, _ := storage.CreateDatabase(&config)

	service := customer.MongoCustomerService{
		Resource: db.Resource,
	}

	server := handler.EchoHandler{
		Port:    "8080",
		Service: &service,
	}

	server.Start()
}
