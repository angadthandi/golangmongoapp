package config

const (
	ServerPort = ":9002" // http server port

	// MongoDB
	MongoDBPort        = ":27016"
	MongoDBServiceName = "mongodbservice"
	MongoDBUsername    = "root"
	MongoDBPassword    = "password"

	// RabbitMQ Exchange
	ExchangeName = "golangmongoapp_direct"
	ExchangeType = "direct"

	// Goapp Routing Key
	GoappRoutingKey = "goapp"

	// RabbitMQ Products Routing Key
	ProductsRoutingKey = "product"
)
