package messages

import (
	"log"
	"time"

	"github.com/angadthandi/golangmongoapp/products/config"
	"github.com/streadway/amqp"
)

func Send() {
	conn, err := amqp.Dial(
		"amqp://" +
			config.RabbitMQUsername + ":" +
			config.RabbitMQPassword + "@" +
			config.RabbitMQServiceName +
			config.RabbitMQPort)
	if err != nil {
		log.Printf("Failed to connect to rabbitmq server: %v", err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %v", err)
		return
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Printf("Failed to declare a exchange: %v", err)
		return
	}

	// q, err := ch.QueueDeclare(
	// 	"product", // name
	// 	false,     // durable
	// 	false,     // delete when unused
	// 	false,     // exclusive
	// 	false,     // no-wait
	// 	nil,       // arguments
	// )
	// if err != nil {
	// 	log.Printf("Failed to declare a queue: %v", err)
	// 	return
	// }

	body := "Product Info!"
	err = ch.Publish(
		"logs", // exchange
		"",     //q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	if err != nil {
		log.Printf("Failed to publish a message: %v", err)
		return
	}
}

func Receive() {
	log.Println("ProductApp Receiver Started")

	connStr := "amqp://" +
		config.RabbitMQUsername + ":" +
		config.RabbitMQPassword + "@" +
		config.RabbitMQServiceName +
		config.RabbitMQPort

	conn, err := amqp.Dial(connStr)
	if err != nil {
		log.Printf("Failed to connect to rabbitmq server: %v", err)

		for {
			reconn, err := amqp.Dial(connStr)

			if err == nil {
				conn = reconn
				break
			}

			log.Println(err)
			log.Printf("Trying to reconnect to RabbitMQ at %s\n", connStr)
			time.Sleep(500 * time.Millisecond)
		}
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %v", err)
		return
	}
	defer ch.Close()

	// declare exchange
	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Printf("Failed to declare a exchange: %v", err)
		return
	}

	// declare queue
	q, err := ch.QueueDeclare(
		"",    //"hello", // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Printf("Failed to declare a queue: %v", err)
		return
	}

	// bind queue to exchnage
	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"logs", // exchange
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to declare a queue: %v", err)
		return
	}

	// consume messages on queue bound to exchange
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Printf("Failed to register a consumer: %v", err)
		return
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
