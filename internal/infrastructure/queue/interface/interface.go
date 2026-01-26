package _interface

import (
	"context"
)

type Queue interface {
	CreateTopic(ctx context.Context, topicOption TopicOption) error
	Publish(ctx context.Context, topicOption TopicOption, message *Message) error
	Consume(ctx context.Context, topicOption TopicOption, handler Handler) error
	Close(ctx context.Context) error
}
