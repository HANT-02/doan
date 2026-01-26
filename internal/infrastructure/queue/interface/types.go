package _interface

import (
	"context"
	"time"
)

type Message struct {
	Id   *string
	Key  interface{}
	Data interface{}
	Meta map[string]interface{}
}

type Handler func(ctx context.Context, message *Message) (Response, error)

type Response int

const (
	Success Response = 0
	Failed  Response = 1
	Retry   Response = 2
)

type RabbitMQConfig struct {
	URL                string
	MaxTryTimes        int
	DeadLetterExchange string
	PrefetchCount      int
	RetryDelay         *time.Duration
	Exchange           *string
	RoutingKey         *string
	Queue              *string
}

type ExchangeKind string

const (
	Fanout ExchangeKind = "fanout"
	Direct ExchangeKind = "direct"
	Topic  ExchangeKind = "topic"
)

type RabbitTopicOption struct {
	ExchangeName *string
	QueueName    *string
	RoutingKey   *string
	Kind         *ExchangeKind
}

type KafkaConsumerConfig struct {
	GroupID         string
	Topic           string
	MaxRetry        int
	MaxMessage      int
	RetryDelay      *time.Duration
	DeadLetterTopic string
	RetryTopic      *string
	CommitInterval  time.Duration
}

type KafkaConfig struct {
	Brokers        []string
	ConsumerConfig KafkaConsumerConfig
}

type KafkaTopicOption struct {
	TopicName         string
	NumPartitions     int
	ReplicationFactor int
}

type TopicOption struct {
	RabbitTopicOption *RabbitTopicOption
	KafkaTopicOption  *KafkaTopicOption
}
