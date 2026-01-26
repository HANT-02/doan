package rabbitmqtemp

import (
	"context"
	_interface "doan/internal/infrastructure/queue/interface"
	"doan/pkg/logger"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"log"

	"time"
)

type rabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  *_interface.RabbitMQConfig
}

func (q *rabbitMQ) CreateTopic(ctx context.Context, option _interface.TopicOption) error {
	ctxLogger := logger.NewLogger(ctx)
	if option.RabbitTopicOption == nil {
		ctxLogger.Errorf("RabbitMQ topic option is required")
		return fmt.Errorf("RabbitMQ topic option is required")
	}
	exchangeName := option.RabbitTopicOption.ExchangeName
	queueName := option.RabbitTopicOption.QueueName
	routingKey := option.RabbitTopicOption.RoutingKey
	kind := string(_interface.Fanout)
	if exchangeName == nil || *exchangeName == "" {
		ctxLogger.Errorf("Exchange name is required")
		return fmt.Errorf("exchange name is required")
	}
	if queueName == nil || *queueName == "" {
		ctxLogger.Errorf("Queue name is required")
		return fmt.Errorf("queue name is required")
	}
	if routingKey == nil || *routingKey == "" {
		ctxLogger.Errorf("Routing key is required")
		return fmt.Errorf("routing key is required")
	}
	if option.RabbitTopicOption.Kind != nil {
		kind = string(*option.RabbitTopicOption.Kind)
	}
	err := q.setupDeadLetterQueue(ctx, *queueName)
	if err != nil {
		ctxLogger.Errorf("Failed to setup dead letter queue: %v", err)
		return err
	}
	args := amqp.Table{
		"x-dead-letter-exchange":    q.config.DeadLetterExchange,
		"x-dead-letter-routing-key": *queueName + ".dlq",
	}
	err = q.channel.ExchangeDeclare(
		*exchangeName,
		kind,
		true,  // Durable
		false, // Auto-deleted
		false, // Internal
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		ctxLogger.Errorf("Failed to declare exchange: %v", err)
		return err
	}
	_, err = q.channel.QueueDeclare(
		*queueName,
		true,  // Durable
		false, // Delete when unused
		false, // Exclusive
		false, // No-wait
		args,  // Arguments
	)
	if err != nil {
		ctxLogger.Errorf("Failed to declare queue: %v", err)
		return err
	}
	err = q.channel.QueueBind(
		*queueName,
		*routingKey,
		*exchangeName,
		false,
		nil,
	)
	if err != nil {
		ctxLogger.Errorf("Failed to bind exchange: %v", err)
		return err
	}
	return nil
}

func (q *rabbitMQ) Publish(ctx context.Context, topicOption _interface.TopicOption, message *_interface.Message) error {
	ctxLogger := logger.NewLogger(ctx)
	messageId := uuid.New().String()
	messageValue, err := json.Marshal(message.Data)
	if err != nil {
		ctxLogger.Errorf("Failed to marshal message data: %v", err)
		return err
	}
	exchange := ""
	if topicOption.RabbitTopicOption.ExchangeName != nil {
		exchange = *topicOption.RabbitTopicOption.ExchangeName
	}
	routingKey := ""
	if topicOption.RabbitTopicOption.RoutingKey == nil {
		panic("routing key is required")
	}
	routingKey = *topicOption.RabbitTopicOption.RoutingKey
	err = q.channel.Publish(
		exchange,   // Exchange
		routingKey, // Routing key
		false,      // Mandatory
		false,      // Immediate
		amqp.Publishing{
			MessageId:   messageId,
			Headers:     message.Meta,
			ContentType: "application/octet-stream",
			Body:        messageValue,
		},
	)
	return err
}

func (q *rabbitMQ) Consume(ctx context.Context, topicOption _interface.TopicOption, handler _interface.Handler) error {
	ctxLogger := logger.NewLogger(ctx)
	if topicOption.RabbitTopicOption.QueueName == nil {
		panic("queue name is required")
	}
	msgs, err := q.channel.Consume(
		*topicOption.RabbitTopicOption.QueueName,
		"",
		false, // Auto-ack
		false, // Exclusive
		false, // No-local
		false, // No-wait
		nil,   // Args
	)
	if err != nil {
		ctxLogger.Errorf("Failed to consume message: %v", err)
		return err
	}
	go func() {
		for msg := range msgs {
			message := &_interface.Message{
				Meta: msg.Headers,
				Data: msg.Body,
			}
			status, err := handler(ctx, message)
			if err != nil {
				ctxLogger.Errorf("[RabbitMQ] Error handling message: %v", err)
			}

			switch status {
			case _interface.Success:
				if err := msg.Ack(false); err != nil {
					ctxLogger.Errorf("[RabbitMQ] Error acknowledging message: %v", err)
				}
			case _interface.Retry:
				retryCount := int32(0)
				if retryCountRaw, isExist := msg.Headers["retry-count"]; isExist {
					retryCount = retryCountRaw.(int32)
				}
				if retryCount >= int32(q.config.MaxTryTimes) {
					if err := msg.Reject(false); err != nil {
						ctxLogger.Errorf("[RabbitMQ] Error nack message: %v", err)
					}
					continue
				}
				if err := q.retryMessage(msg, *topicOption.RabbitTopicOption.QueueName, retryCount); err != nil {
					ctxLogger.Errorf("[RabbitMQ] Error nack message: %v", err)
				}
			default:
				if err := msg.Ack(false); err != nil {
					ctxLogger.Errorf("[RabbitMQ] Error acknowledging message: %v", err)
				}
			}
		}

	}()

	return nil
}

func (q *rabbitMQ) retryMessage(msg amqp.Delivery, topic string, retryCount int32) error {
	headers := amqp.Table{}
	for k, v := range msg.Headers {
		headers[k] = v
	}
	headers["retry-count"] = retryCount + 1
	err := q.channel.Publish(
		"",
		topic,
		false,
		false,
		amqp.Publishing{
			ContentType: msg.ContentType,
			Body:        msg.Body,
			Headers:     headers,
			Timestamp:   time.Now(),
		},
	)
	if err != nil {
		return err
	}
	if err := msg.Ack(false); err != nil {
		log.Printf("[RabbitMQ] Error acknowledging message: %v", err)
	}

	return nil
}

func (q *rabbitMQ) setupDeadLetterQueue(ctx context.Context, topic string) error {
	// Khai báo Dead-letter Exchange
	err := q.channel.ExchangeDeclare(
		q.config.DeadLetterExchange,
		"direct",
		true,  // Durable
		false, // Auto-deleted
		false, // Internal
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		return err
	}

	// Khai báo Dead-letter Queue
	dlqName := topic + ".dlq"
	_, err = q.channel.QueueDeclare(
		dlqName,
		true,  // Durable
		false, // Delete when unused
		false, // Exclusive
		false, // No-wait
		nil,   // Arguments
	)
	if err != nil {
		return err
	}

	// Bind DLQ với DLX
	err = q.channel.QueueBind(
		dlqName,
		dlqName,
		q.config.DeadLetterExchange,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

func (q *rabbitMQ) Close(ctx context.Context) error {
	var err error
	if q.channel != nil {
		err = q.channel.Close()
	}
	if q.conn != nil {
		connErr := q.conn.Close()
		if connErr != nil && err == nil {
			err = connErr
		}
	}
	return err
}

func New(config *_interface.RabbitMQConfig) (_interface.Queue, error) {
	if config.URL == "" {
		return nil, fmt.Errorf("RabbitMQ URL is required")
	}
	conn, err := amqp.DialConfig(config.URL,
		amqp.Config{
			Heartbeat: 10 * time.Second,
		},
	)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		panic(err)
	}

	return &rabbitMQ{
		conn:    conn,
		channel: channel,
		config:  config,
	}, nil
}
