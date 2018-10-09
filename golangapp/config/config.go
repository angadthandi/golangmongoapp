package config

const (
	ServerPort = ":9002" // http server port

	// MongoDB
	MongoDBPort        = ":27016"
	MongoDBServiceName = "mongodbservice"
	MongoDBUsername    = "root"
	MongoDBPassword    = "password"

	// RabbitMQ
	RabbitMQPort        = ":5672"
	RabbitMQServiceName = "golangrabbitmq"
	RabbitMQUsername    = "guest"
	RabbitMQPassword    = "guest"
)
