package messages

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// Defines our interface for connecting and consuming messages.
type IMessagingClient interface {
	Connect()
	Send(exchangeName string, exchangeType string, publishRoutingKey string, msg []byte)
	Receive(exchangeName string, exchangeType string, sbindRoutingKeys []string)
	Close()
}

// Real implementation, encapsulates a pointer to an amqp.Connection
type MessagingClient struct {
	conn *amqp.Connection
}

func (m *MessagingClient) Connect() {
	connStr := "amqp://" +
		RabbitMQUsername + ":" +
		RabbitMQPassword + "@" +
		RabbitMQServiceName +
		RabbitMQPort

	var err error
	m.conn, err = amqp.Dial(connStr)
	if err != nil {
		log.Errorf("Failed to connect to rabbitmq server: %v", err)

		for {
			reconn, err := amqp.Dial(connStr)

			if err == nil {
				m.conn = reconn
				break
			}

			log.Debugf("AMQP connection error: %v\n", err)
			log.Debugf("Trying to reconnect to RabbitMQ at %s\n", connStr)
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func (m *MessagingClient) Send(
	exchangeName string,
	exchangeType string,
	publishRoutingKey string,
	msg []byte,
) {
	ch, err := m.conn.Channel()
	if err != nil {
		log.Errorf("Failed to open a channel: %v", err)
		return
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		exchangeName, // "logs_direct", // name
		exchangeType, // "direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Errorf("Failed to declare a exchange: %v", err)
		return
	}

	err = ch.Publish(
		exchangeName,      // "logs_direct",     // exchange
		publishRoutingKey, //q.Name, // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})
	log.Debugf(" [x] Sent Message: %s", msg)
	if err != nil {
		log.Errorf("Failed to publish a message: %v", err)
		return
	}
}

func (m *MessagingClient) Receive(
	exchangeName string,
	exchangeType string,
	sbindRoutingKeys []string,
) {
	log.Debugf("Receiver %v", "Started")

	ch, err := m.conn.Channel()
	if err != nil {
		log.Errorf("Failed to open a channel: %v", err)
		return
	}
	defer ch.Close()

	// declare exchange
	err = ch.ExchangeDeclare(
		exchangeName, // "logs_direct", // name
		exchangeType, // "direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Errorf("Failed to declare a exchange: %v", err)
		return
	}

	// declare queue
	q, err := ch.QueueDeclare(
		"",    //"product", // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Errorf("Failed to declare a queue: %v", err)
		return
	}

	for _, rKey := range sbindRoutingKeys {
		// bind queue to exchnage
		err = ch.QueueBind(
			q.Name,       // queue name
			rKey,         // routing key
			exchangeName, // "logs_direct", // exchange
			false,
			nil,
		)
		if err != nil {
			log.Errorf("Failed to declare a queue: %v", err)
			return
		}
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
		log.Errorf("Failed to register a consumer: %v", err)
		return
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Debugf("Received a message: %s", d.Body)
		}
	}()

	log.Debugf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (m *MessagingClient) Close() {
	if m.conn != nil {
		m.conn.Close()
	}
}
