package main

import (
	customer "goworkshop2/customer"
	handler "goworkshop2/customer/handler"
	"goworkshop2/storage"
)

func main() {
	storage := storage.MongoDb{
		Config: storage.MongoDbConfiguration{
			ConnectionTimeout:        5,
			ConnectionStringTemplate: "mongodb://%s:%s@%s",
		},
	}
	storage.InitializeResource()

	service := customer.CustomerService{
		Resource: storage.Resource,
	}

	server := handler.EchoHandler{
		Service: &service,
		Port:    "8080",
	}

	server.Start()
}
