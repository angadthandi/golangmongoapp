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

	// Goapp Publish Routing Key
	GoappPublishRoutingKey = "goapp"

	// RabbitMQ Products Routing Key
	ProductsRoutingKey = "product"
)
