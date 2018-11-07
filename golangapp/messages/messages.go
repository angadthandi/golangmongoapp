package messages

import (
	"encoding/json"
	"math/rand"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

// Defines our interface for connecting and consuming messages.
type IMessagingClient interface {
	Connect()
	Send(
		exchangeName string,
		exchangeType string,
		publishRoutingKey string,
		replyToRoutingKey string,
		msg []byte,
		MessagesRegistryClient IMessagesRegistry,
		receivedCorrelationId string,

		// should be set to "false" to register message,
		// based on correlationId match found or not
		isReplyMessage bool,
	)
	Receive(
		exchangeName string,
		exchangeType string,
		receiveRoutingKey string, // local app key
		handlerFunc func(
			*MessagingClient,
			IMessagesRegistry,
			json.RawMessage,
			string,
			string,
			bool,
			*mongo.Database,
		),
		MessagesRegistryClient IMessagesRegistry,
		dbRef *mongo.Database,
	)
	Close()
}

// Real implementation, encapsulates a pointer to an amqp.Connection
type MessagingClient struct {
	conn *amqp.Connection
}

func (m *MessagingClient) Connect() {
	// Initialize random seed default value
	// for unique CorrelationId
	rand.Seed(time.Now().UTC().UnixNano())

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
	replyToRoutingKey string,
	msg []byte,
	MessagesRegistryClient IMessagesRegistry,
	receivedCorrelationId string,
	isReplyMessage bool,
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

	correlationId := randomString(32)
	if receivedCorrelationId != "" {
		correlationId = receivedCorrelationId
	}

	err = ch.Publish(
		exchangeName,      // "logs_direct",     // exchange
		publishRoutingKey, //q.Name, // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,

			// Store CorrelationId in a local map before publish.
			// On receive, check for CorrelationId in local map.
			//
			// Work on received response for CorrelationId in local map,
			// then delete CorrelationId from local map
			CorrelationId: correlationId,
			ReplyTo:       replyToRoutingKey,
		})
	log.Debugf(" [x] Sent Message: %s", msg)
	if err != nil {
		log.Errorf("Failed to publish a message: %v", err)
		return
	}

	// register sent message only if NOT a reply
	if !isReplyMessage {
		// register sent message
		log.Debugf(`messages Send: SetCorrelationMapData:
	correlationId: %v, sentToAppName: %v, sentToAppEvent: %v`,
			correlationId, publishRoutingKey, publishRoutingKey)
		MessagesRegistryClient.SetCorrelationMapData(
			correlationId, publishRoutingKey, publishRoutingKey,
		)
	}
}

func (m *MessagingClient) Receive(
	exchangeName string,
	exchangeType string,
	receiveRoutingKey string, // local app key
	handlerFunc func(
		*MessagingClient,
		IMessagesRegistry,
		json.RawMessage,
		string,
		string,
		bool,
		*mongo.Database,
	),
	MessagesRegistryClient IMessagesRegistry,
	dbRef *mongo.Database,
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

	// bind queue to exchnage
	err = ch.QueueBind(
		q.Name,            // queue name
		receiveRoutingKey, // routing key
		exchangeName,      // "logs_direct", // exchange
		false,
		nil,
	)
	if err != nil {
		log.Errorf("Failed to bind to a queue: %v", err)
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
		log.Errorf("Failed to register a consumer: %v", err)
		return
	}

	forever := make(chan bool)

	go m.consumeLoop(msgs, MessagesRegistryClient, dbRef, handlerFunc)

	log.Debugf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (m *MessagingClient) Close() {
	if m.conn != nil {
		m.conn.Close()
	}
}

func (m *MessagingClient) consumeLoop(
	deliveries <-chan amqp.Delivery,
	MessagesRegistryClient IMessagesRegistry,
	dbRef *mongo.Database,
	handlerFunc func(
		*MessagingClient,
		IMessagesRegistry,
		json.RawMessage,
		string,
		string,
		bool,
		*mongo.Database,
	),
) {
	for d := range deliveries {
		// Invoke the handlerFunc func we passed as parameter
		// via handleRefreshEvent
		handleRefreshEvent(d, m, MessagesRegistryClient, dbRef, handlerFunc)

		// Update the data on the service's
		// associated datastore using a local transaction...

		// The 'false' indicates the success of a single delivery, 'true' would
		// mean that this delivery and all prior unacknowledged deliveries on this
		// channel will be acknowledged.
		// d.Ack(false)
	}
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
