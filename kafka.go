package randomchat

import (
	"time"

	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
)

var (
	logger = watermill.NewStdLogger(
		false,
		false,
	)
)

func NewKafkaPublisher() (message.Publisher, error) {
	kafkaPublisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   getServerAddrs(getenv("KAFKA_ADDRS", "localhost:9092")),
			Marshaler: kafka.DefaultMarshaler{},
		},
		logger,
	)
	if err != nil {
		return nil, err
	}

	return kafkaPublisher, nil
}

func NewKafkaSubscriber() (message.Subscriber, error) {
	config := sarama.NewConfig()
	config.Consumer.Fetch.Default = 1024 * 1024
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second

	kafkaSubscriber, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:       getServerAddrs(getenv("KAFKA_ADDRS", "localhost:9092")),
			Unmarshaler:   kafka.DefaultMarshaler{},
			ConsumerGroup: watermill.NewUUID(),
			InitializeTopicDetails: &sarama.TopicDetail{
				NumPartitions:     1,
				ReplicationFactor: 2,
			},
			OverwriteSaramaConfig: config,
		},
		logger,
	)
	if err != nil {
		return nil, err
	}

	return kafkaSubscriber, nil
}

func NewMessageRouter() (*message.Router, error) {
	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		return nil, err
	}

	router.AddMiddleware(
		middleware.CorrelationID,
		middleware.Timeout(time.Second*15),
		middleware.Recoverer,
	)
	return router, nil
}
