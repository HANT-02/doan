package queue_temp

import (
	"doan/internal/infrastructure/queue/interface"
	kafka "doan/internal/infrastructure/queue/kafka"
	rabbitmq "doan/internal/infrastructure/queue/rabbitmq"
	"time"
)

type Type string

const (
	RabbitMQ Type = "rabbitmq"
	Kafka    Type = "kafka"
)

type ConfigFactory struct {
	Type           Type
	RabbitMQConfig *_interface.RabbitMQConfig
	KafkaConfig    *_interface.KafkaConfig
}

func MappingConfig(config ConsumerConfig) ConfigFactory {
	driver := Type(config.Driver)
	switch driver {
	case RabbitMQ:
		if config.RabbitMQ == nil {
			panic("invalid rabbitmq config")
		}
		var retryDelay *time.Duration
		if config.RabbitMQ.RetryDelay != "" {
			parseRetryDelay, err := time.ParseDuration(config.RabbitMQ.RetryDelay)
			if err != nil {
				panic(err)
			}
			retryDelay = &parseRetryDelay
		}
		return ConfigFactory{
			Type: RabbitMQ,
			RabbitMQConfig: &_interface.RabbitMQConfig{
				URL:                config.RabbitMQ.URL,
				MaxTryTimes:        config.RabbitMQ.MaxTryTimes,
				DeadLetterExchange: config.RabbitMQ.DeadLetterExchange,
				PrefetchCount:      config.RabbitMQ.PrefetchCount,
				RetryDelay:         retryDelay,
				Exchange:           config.RabbitMQ.Exchange,
				RoutingKey:         config.RabbitMQ.RoutingKey,
			},
		}
	case Kafka:
		if config.Kafka == nil {
			panic("invalid kafka config")
		}
		var retryDelay *time.Duration
		if config.Kafka.RetryDelay != "" {
			parseRetryDelay, err := time.ParseDuration(config.Kafka.RetryDelay)
			if err != nil {
				panic(err)
			}
			retryDelay = &parseRetryDelay
		}
		var commitInterval time.Duration
		if config.Kafka.CommitInterval != "" {
			parseCommitInterval, err := time.ParseDuration(config.Kafka.CommitInterval)
			if err != nil {
				panic(err)
			}
			commitInterval = parseCommitInterval
		}
		return ConfigFactory{
			Type: Kafka,
			KafkaConfig: &_interface.KafkaConfig{
				Brokers: config.Kafka.Brokers,
				ConsumerConfig: _interface.KafkaConsumerConfig{
					GroupID:         config.Kafka.GroupID,
					MaxRetry:        config.Kafka.MaxRetry,
					MaxMessage:      config.Kafka.MaxMessage,
					RetryDelay:      retryDelay,
					DeadLetterTopic: config.Kafka.DeadLetterTopic,
					CommitInterval:  commitInterval,
				},
			},
		}
	default:
		panic("invalid queue config")
	}
}

func MappingTopicOption(config ConsumerConfig) _interface.TopicOption {
	driver := Type(config.Driver)
	switch driver {
	case RabbitMQ:
		if config.RabbitMQ == nil {
			panic("invalid rabbitmq config")
		}
		return _interface.TopicOption{
			RabbitTopicOption: &_interface.RabbitTopicOption{
				ExchangeName: config.RabbitMQ.Exchange,
				RoutingKey:   config.RabbitMQ.RoutingKey,
				QueueName:    config.RabbitMQ.Queue,
			},
		}
	case Kafka:
		if config.Kafka == nil {
			panic("invalid kafka config")
		}
		return _interface.TopicOption{
			KafkaTopicOption: &_interface.KafkaTopicOption{
				TopicName: config.Kafka.Topic,
			},
		}
	default:
		panic("invalid queue config")
	}
}

func New(configFactory ConfigFactory) (_interface.Queue, error) {
	switch configFactory.Type {
	case RabbitMQ:
		return rabbitmq.New(configFactory.RabbitMQConfig)
	case Kafka:
		return kafka.New(configFactory.KafkaConfig)
	default:
		panic("invalid queue type")
	}
}
