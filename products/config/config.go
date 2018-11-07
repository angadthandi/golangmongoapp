package config

const (
	ServerPort = ":9003" // http server port

	// MongoDB
	MongoDBPort        = ":27018"
	MongoDBServiceName = "products_mongodbservice"
	MongoDBUsername    = "root"
	MongoDBPassword    = "password"

	// RabbitMQ Exchange
	ExchangeName = "golangmongoapp_direct"
	ExchangeType = "direct"

	// Products Routing Key
	ProductsRoutingKey = "product"

	// RabbitMQ Goapp Routing Key
	GoappRoutingKey = "golangapp"
)
