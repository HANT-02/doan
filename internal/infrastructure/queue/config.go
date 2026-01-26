package queue_temp

type ConsumerConfig struct {
	Driver   string          `json:"driver" yaml:"driver"`
	RabbitMQ *RabbitMQConfig `json:"rabbitmq" yaml:"rabbitmq"`
	Kafka    *KafkaConfig    `json:"kafka" yaml:"kafka"`
}

type KafkaConfig struct {
	Brokers         []string `json:"brokers" yaml:"brokers"`
	GroupID         string   `json:"groupID" yaml:"groupID"`
	MaxRetry        int      `json:"maxRetry" yaml:"maxRetry"`
	MaxMessage      int      `json:"maxMessage" yaml:"maxMessage"`
	RetryDelay      string   `json:"retryDelay" yaml:"retryDelay"`
	CommitInterval  string   `json:"commitInterval" yaml:"commitInterval"`
	Topic           string   `json:"topic" yaml:"topic"`
	DeadLetterTopic string   `json:"deadLetterTopic" yaml:"deadLetterTopic"`
}

type RabbitMQConfig struct {
	URL                string  `json:"url" yaml:"url"`
	MaxTryTimes        int     `json:"maxTryTimes" yaml:"maxTryTimes"`
	DeadLetterExchange string  `json:"deadLetterExchange" yaml:"deadLetterExchange"`
	PrefetchCount      int     `json:"prefetchCount" yaml:"prefetchCount"`
	RetryDelay         string  `json:"retryDelay" yaml:"retryDelay"`
	Exchange           *string `json:"exchange" yaml:"exchange"`
	Queue              *string `json:"queue" yaml:"queue"`
	RoutingKey         *string `json:"routingKey" yaml:"routingKey"`
}
