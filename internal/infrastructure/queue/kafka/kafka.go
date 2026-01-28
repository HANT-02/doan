package kafkatemp

import (
	"context"
	_interface "doan/internal/infrastructure/queue/interface"
	"doan/pkg/logger"
	xerror "doan/pkg/x-error"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"sync"
	"time"
)

type kafkaQueue struct {
	config      *_interface.KafkaConfig
	writer      *kafka.Writer
	readers     map[string]*kafka.Reader
	cancelFuncs map[string]context.CancelFunc
	mu          sync.Mutex
}

func New(config *_interface.KafkaConfig) (_interface.Queue, error) {
	if len(config.Brokers) == 0 {
		return nil, xerror.NewError(xerror.KafkaBrokersRequired)
	}
	address := kafka.TCP(config.Brokers...)

	writer := &kafka.Writer{
		Addr:         address,
		Balancer:     &kafka.Hash{},
		Async:        true,
		RequiredAcks: kafka.RequireAll,
	}

	return &kafkaQueue{
		config:      config,
		writer:      writer,
		readers:     make(map[string]*kafka.Reader),
		cancelFuncs: make(map[string]context.CancelFunc),
	}, nil
}

func (q *kafkaQueue) CreateTopic(ctx context.Context, topicOption _interface.TopicOption) error {
	ctxLogger := logger.NewLogger(ctx)
	if topicOption.KafkaTopicOption == nil {
		return xerror.NewError(xerror.KafkaTopicOptionRequired)
	}
	if topicOption.KafkaTopicOption.TopicName == "" {
		return xerror.NewError(xerror.KafkaTopicNameRequired)
	}
	if topicOption.KafkaTopicOption.NumPartitions <= 0 {
		return xerror.NewError(xerror.KafkaTopicNumPartitionsInvalid)
	}
	if topicOption.KafkaTopicOption.ReplicationFactor <= 0 {
		return xerror.NewError(xerror.KafkaTopicReplicationFactorInvalid)
	}
	client := &kafka.Client{
		Addr:    kafka.TCP(q.config.Brokers...),
		Timeout: 10 * time.Second,
	}
	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topicOption.KafkaTopicOption.TopicName,
			NumPartitions:     topicOption.KafkaTopicOption.NumPartitions,
			ReplicationFactor: topicOption.KafkaTopicOption.ReplicationFactor,
		},
	}
	resp, err := client.CreateTopics(ctx, &kafka.CreateTopicsRequest{
		Topics: topicConfigs,
	})
	if err != nil {
		ctxLogger.Errorf("[Kafka] Error creating topic: %v", err)
		return xerror.NewError(xerror.KafkaTopicCreateFailed)
	}
	failedTopics := make([]string, 0)
	for failedTopic, rErr := range resp.Errors {
		failedTopics = append(failedTopics, failedTopic)
		ctxLogger.Errorf("[Kafka] Failed to create topic '%s': %v", failedTopic, rErr)
	}
	if len(failedTopics) > 0 {
		return xerror.NewError(xerror.KafkaTopicCreateFailed)
	}
	ctxLogger.Infof(
		"[Kafka] Topic '%s' created successfully with %d partitions and replication factor %d",
		topicOption.KafkaTopicOption.TopicName,
		topicOption.KafkaTopicOption.NumPartitions,
		topicOption.KafkaTopicOption.ReplicationFactor,
	)
	return nil
}

func (q *kafkaQueue) Publish(ctx context.Context, topicOption _interface.TopicOption, message *_interface.Message) error {
	ctxLogger := logger.NewLogger(ctx)
	var headers []kafka.Header
	for k, v := range message.Meta {
		headers = append(headers, kafka.Header{
			Key:   k,
			Value: []byte(fmt.Sprintf("%v", v)),
		})
	}
	msgId := uuid.NewString()
	if message.Id != nil {
		msgId = *message.Id
	}
	headers = append(headers, kafka.Header{
		Key:   "MessageId",
		Value: []byte(msgId),
	})
	key, err := json.Marshal(message.Key)
	if err != nil {
		ctxLogger.Errorf("[Kafka] Error marshalling message key: %v", err)
		return xerror.NewError(xerror.KafkaPublishFailed)
	}
	value, err := json.Marshal(message.Data)
	if err != nil {
		ctxLogger.Errorf("[Kafka] Error marshalling message data: %v", err)
		return xerror.NewError(xerror.KafkaPublishFailed)
	}
	if topicOption.KafkaTopicOption.TopicName == "" {
		return xerror.NewError(xerror.KafkaTopicNameRequired)
	}
	topic := topicOption.KafkaTopicOption.TopicName
	msg := &kafka.Message{
		Topic:      topic,
		Key:        key,
		Value:      value,
		Headers:    headers,
		WriterData: nil,
		Time:       time.Time{},
	}
	err = q.writer.WriteMessages(ctx, *msg)
	if err != nil {
		ctxLogger.Errorf("[Kafka] Error publishing message: %v", err)
		return xerror.NewError(xerror.KafkaPublishFailed)
	}
	ctxLogger.Infof("[Kafka] Message published to topic %s - id: %s", topic, msgId)
	return nil
}
func (q *kafkaQueue) Consume(ctx context.Context, topicOption _interface.TopicOption, handler _interface.Handler) error {
	ctxLogger := logger.NewLogger(ctx)
	if topicOption.KafkaTopicOption.TopicName == "" {
		return xerror.NewError(xerror.KafkaTopicNameRequired)
	}
	topic := topicOption.KafkaTopicOption.TopicName
	q.mu.Lock()
	if _, exists := q.readers[topic]; exists {
		q.mu.Unlock()
		return xerror.NewError(xerror.KafkaTopicAlreadyConsumed)
	}

	readerConfig := kafka.ReaderConfig{
		Brokers:        q.config.Brokers,
		GroupID:        q.config.ConsumerConfig.GroupID,
		Topic:          topic,
		CommitInterval: q.config.ConsumerConfig.CommitInterval,
	}

	reader := kafka.NewReader(readerConfig)
	ctx, cancel := context.WithCancel(ctx)
	q.readers[topic] = reader
	q.cancelFuncs[topic] = cancel
	q.mu.Unlock()
	ctxLogger.Infof("[Kafka] Bắt đầu consumer cho topic: %s", topic)

	// Chọn xử lý tin nhắn theo kiểu batch hoặc từng tin nhắn
	if q.config.ConsumerConfig.MaxMessage > 1 {
		go q.consumeBatch(ctx, reader, handler)
	} else {
		go q.consumeSingle(ctx, reader, handler)
	}

	return nil
}

// Hàm xử lý từng tin nhắn một
func (q *kafkaQueue) consumeSingle(ctx context.Context, reader *kafka.Reader, handler _interface.Handler) {
	defer func(reader *kafka.Reader) {
		err := reader.Close()
		if err != nil {
			ctxLogger := logger.NewLogger(ctx)
			ctxLogger.Errorf("[Kafka] Lỗi khi đóng reader cho topic %s: %v", reader.Config().Topic, err)
		}
	}(reader)
	ctxLogger := logger.NewLogger(ctx)
	for {
		m, err := reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				ctxLogger.Warnf("[Kafka] Dừng consumer cho topic: %s", reader.Config().Topic)
				return
			}
			ctxLogger.Errorf("[Kafka] Lỗi khi fetch tin nhắn từ topic %s: %v", reader.Config().Topic, err)
			continue
		}

		// Xử lý logic cho từng tin nhắn
		q.handleMessage(ctx, reader, handler, m)
	}
}

// Hàm xử lý batch tin nhắn
func (q *kafkaQueue) consumeBatch(ctx context.Context, reader *kafka.Reader, handler _interface.Handler) {
	defer func(reader *kafka.Reader) {
		err := reader.Close()
		if err != nil {
			ctxLogger := logger.NewLogger(ctx)
			ctxLogger.Errorf("[Kafka] Lỗi khi đóng reader cho topic %s: %v", reader.Config().Topic, err)
		}
	}(reader)
	ctxLogger := logger.NewLogger(ctx)
	for {
		var messages []kafka.Message

		for i := 0; i < q.config.ConsumerConfig.MaxMessage; i++ {
			// Tạo context có timeout để giới hạn thời gian chờ
			fetchCtx, cancel := context.WithTimeout(ctx, 1500*time.Millisecond)
			defer cancel() // Hủy context sau khi sử dụng

			m, err := reader.FetchMessage(fetchCtx)
			if err != nil {
				if errors.Is(err, context.Canceled) {
					ctxLogger.Warnf("[Kafka] Dừng consumer cho topic: %s", reader.Config().Topic)
					return
				}

				// Kiểm tra nếu lỗi là do timeout, bỏ qua và tiếp tục
				if errors.Is(err, context.DeadlineExceeded) {
					break
				}

				// Nếu gặp lỗi tạm thời, bỏ qua và tiếp tục
				ctxLogger.Warnf("[Kafka] Lỗi khi fetch tin nhắn từ topic %s: %v", reader.Config().Topic, err)
				break
			}
			messages = append(messages, m)

			// Nếu đã đạt đến số lượng tin nhắn tối đa cho batch, thì dừng vòng lặp
			if len(messages) >= q.config.ConsumerConfig.MaxMessage {
				break
			}
		}

		// Xử lý batch tin nhắn song song
		if len(messages) > 0 {
			ctxLogger.Infof("[Kafka] Xử lý batch tin nhắn từ topic %s", reader.Config().Topic)
			var wg sync.WaitGroup
			wg.Add(len(messages))

			for _, m := range messages {
				go func(msg kafka.Message) {
					defer wg.Done()
					defer func() {
						if r := recover(); r != nil {

							ctxLogger.Errorf("[Kafka] Panic khi xử lý tin nhắn từ topic %s: %v", reader.Config().Topic, r)

							// Đẩy tin nhắn vào DLQ khi xảy ra panic
							q.sendToDeadLetterTopic(ctx, msg)
							// Commit message
							if err := reader.CommitMessages(ctx, msg); err != nil {
								ctxLogger.Errorf("[Kafka] Lỗi khi commit tin nhắn từ topic %s: %v", reader.Config().Topic, err)
							}
						}
					}()

					// Xử lý tin nhắn
					q.handleMessage(ctx, reader, handler, msg)
				}(m)
			}

			wg.Wait()
		}
	}
}
func (q *kafkaQueue) handleMessage(ctx context.Context, reader *kafka.Reader, handler _interface.Handler, m kafka.Message) {
	ctxLogger := logger.NewLogger(ctx)
	msgId := getMessageId(extractHeaders(m.Headers))
	message := &_interface.Message{
		Key:  m.Key,
		Data: m.Value,
		Meta: extractHeaders(m.Headers),
	}

	retryCount := 0
	for retryCount < q.config.ConsumerConfig.MaxRetry {
		status, err := handler(ctx, message)

		if err != nil {
			ctxLogger.Errorf("[Kafka] Lỗi khi xử lý tin nhắn từ topic %s - id: %s: %v", reader.Config().Topic, msgId, err)
		}

		if status == _interface.Success {
			ctxLogger.Infof("[Kafka] Xử lý tin nhắn thành công từ topic %s - id: %s", reader.Config().Topic, msgId)
			if err := reader.CommitMessages(ctx, m); err != nil {
				ctxLogger.Errorf("[Kafka] Lỗi khi commit tin nhắn từ topic %s - id: %s: %v", reader.Config().Topic, msgId, err)
			}
			break
		} else if status == _interface.Retry {
			retryCount++
			ctxLogger.Warnf("[Kafka] Retry lần thứ %d cho tin nhắn từ topic %s - id: %s", retryCount, reader.Config().Topic, msgId)
		} else if status == _interface.Failed {
			ctxLogger.Warnf("[Kafka] Tin nhắn bị từ chối, gửi vào DLQ cho topic %s - id: %s", reader.Config().Topic, msgId)
			q.sendToDeadLetterTopic(ctx, m)
			if err := reader.CommitMessages(ctx, m); err != nil {
				ctxLogger.Errorf("[Kafka] Lỗi khi commit tin nhắn từ topic %s - id: %s: %v", reader.Config().Topic, msgId, err)
			}
			break
		}
	}

	// Nếu đã retry tối đa số lần và vẫn thất bại, đưa vào DLQ
	if retryCount >= q.config.ConsumerConfig.MaxRetry {
		ctxLogger.Warnf("[Kafka] Đã đạt số lần retry tối đa, gửi tin nhắn vào DLQ cho topic %s - id: %s", reader.Config().Topic, msgId)
		q.sendToDeadLetterTopic(ctx, m)
		if err := reader.CommitMessages(ctx, m); err != nil {
			ctxLogger.Errorf("[Kafka] Lỗi khi commit tin nhắn từ topic %s - id: %s: %v", reader.Config().Topic, msgId, err)
		}
	}
}
func (q *kafkaQueue) sendToDeadLetterTopic(ctx context.Context, msg kafka.Message) {
	ctxLogger := logger.NewLogger(ctx)
	msgId := getMessageId(extractHeaders(msg.Headers))
	if q.config.ConsumerConfig.DeadLetterTopic == "" {
		ctxLogger.Warnf("[Kafka] Dead-letter topic not configured, message lost for topic %s - id: %s: %v", msg.Topic, msgId)
		return
	}
	var headers []kafka.Header
	for _, h := range msg.Headers {
		headers = append(headers, h)
	}
	dlqMessage := kafka.Message{
		Topic:   q.config.ConsumerConfig.DeadLetterTopic,
		Key:     msg.Key,
		Value:   msg.Value,
		Headers: headers,
		Time:    time.Now(),
	}
	err := q.writer.WriteMessages(ctx, dlqMessage)
	if err != nil {
		ctxLogger.Errorf("[Kafka] Error sending message to DLQ for topic %s - id: %s: %v", msg.Topic, msgId, err)
	} else {
		ctxLogger.Warnf("[Kafka] Message sent to DLQ %s - id: %s", q.config.ConsumerConfig.DeadLetterTopic, msgId)
	}
}

func (q *kafkaQueue) publishRetryMessage(ctx context.Context, message *_interface.Message) error {
	ctxLogger := logger.NewLogger(ctx)
	msgId := getMessageId(message.Meta)
	if q.config.ConsumerConfig.RetryTopic == nil {
		ctxLogger.Warnf("[Kafka] Retry topic not configured, message send to DLQ for topic %s - id: %s", q.config.ConsumerConfig.DeadLetterTopic, msgId)
		// Gửi message vào DLQ
		msg := convertMessageToKafkaMessage(message)
		q.sendToDeadLetterTopic(ctx, *msg)
		return nil
	}
	msg := convertMessageToKafkaMessage(message)
	msg.Topic = *q.config.ConsumerConfig.RetryTopic
	err := q.writer.WriteMessages(ctx, *msg)
	if err != nil {
		ctxLogger.Errorf("[Kafka] Error publishing message to retry topic %s - id: %s", *q.config.ConsumerConfig.RetryTopic, msgId)
		return xerror.NewError(xerror.KafkaPublishFailed)
	}
	ctxLogger.Warnf("[Kafka] Message sent to retry topic %s - id: %s", *q.config.ConsumerConfig.RetryTopic, msgId)
	return nil
}

func (q *kafkaQueue) Close(ctx context.Context) error {
	ctxLogger := logger.NewLogger(ctx)
	var err error
	if q.writer != nil {
		if err = q.writer.Close(); err != nil {
			ctxLogger.Errorf("[Kafka] Error closing writer, %v", err)
		}
	}

	q.mu.Lock()
	for topic, reader := range q.readers {
		cancel := q.cancelFuncs[topic]
		cancel()
		if readerErr := reader.Close(); readerErr != nil && err == nil {
			ctxLogger.Errorf("[Kafka] Error closing reader, %v", readerErr)
			err = readerErr
		}
		delete(q.readers, topic)
		delete(q.cancelFuncs, topic)
	}
	q.mu.Unlock()
	return err
}

func extractHeaders(kafkaHeaders []kafka.Header) map[string]interface{} {
	headers := make(map[string]interface{})
	for _, h := range kafkaHeaders {
		headers[h.Key] = string(h.Value)
	}
	return headers
}

func getMessageId(headers map[string]interface{}) string {
	if val, ok := headers["MessageId"]; ok {
		return fmt.Sprintf("%v", val)
	}
	return ""
}

func convertMessageToKafkaMessage(message *_interface.Message) *kafka.Message {
	var headers []kafka.Header
	for k, v := range message.Meta {
		headers = append(headers, kafka.Header{
			Key:   k,
			Value: []byte(fmt.Sprintf("%v", v)),
		})
	}

	msg := &kafka.Message{
		Key:     []byte(fmt.Sprintf("%s", message.Key)),
		Value:   []byte(fmt.Sprintf("%s", message.Data)),
		Headers: headers,
		Time:    time.Now(),
	}
	return msg
}
