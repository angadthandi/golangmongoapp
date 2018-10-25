package config

const (
	ServerPort = ":9003" // http server port

	// RabbitMQ
	RabbitMQPort        = ":5672"
	RabbitMQServiceName = "golangrabbitmq"
	RabbitMQUsername    = "guest"
	RabbitMQPassword    = "guest"

	// RabbitMQ Exchange
	ExchangeName = "golangmongoapp_direct"
	ExchangeType = "direct"

	// Products Publish Routing Key
	ProductsPublishRoutingKey = "product"

	// RabbitMQ Goapp Routing Key
	GoappRoutingKey = "goapp"
)
