package noop

import (
	"context"
	_interface "doan/internal/infrastructure/queue/interface"
	"doan/pkg/logger"
)

type Queue struct{}

func New() _interface.Queue {
	return &Queue{}
}

func (q *Queue) CreateTopic(ctx context.Context, topicOption _interface.TopicOption) error {
	logger.NewLogger(ctx).Info("noop queue CreateTopic called")
	return nil
}

func (q *Queue) Publish(ctx context.Context, topicOption _interface.TopicOption, message *_interface.Message) error {
	logger.NewLogger(ctx).Infof("noop queue Publish called: %+v", message)
	return nil
}

func (q *Queue) Consume(ctx context.Context, topicOption _interface.TopicOption, handler _interface.Handler) error {
	logger.NewLogger(ctx).Info("noop queue Consume called — no-op")
	return nil
}

func (q *Queue) Close(ctx context.Context) error {
	logger.NewLogger(ctx).Info("noop queue Close called")
	return nil
}
