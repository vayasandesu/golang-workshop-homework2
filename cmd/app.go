package main

import (
	customer "goworkshop2/customer"
	handler "goworkshop2/customer/handler"
	"goworkshop2/storage"
)

func main() {
	config := storage.MongoDbConfiguration{
		ConnectionTimeout:        5,
		ConnectionStringTemplate: "mongodb://%s:%s@%s",
	}
	db, _ := storage.CreateDatabase(&config)

	service := customer.MongoCustomer{
		Resource: db.Resource,
	}

	server := handler.EchoHandler{
		Port:    "8080",
		Service: &service,
	}

	server.Start()
}
