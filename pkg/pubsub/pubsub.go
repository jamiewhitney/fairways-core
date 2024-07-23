package pubsub

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Pubsub struct {
	Client *sqs.Client
}

type Config struct {
}

func New(ctx context.Context) (*Pubsub, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	client := sqs.NewFromConfig(cfg)

	return &Pubsub{
		Client: client,
	}, nil
}

func (p *Pubsub) ReceiveMessages(ctx context.Context, queueURL string, ch chan *types.Message) error {
	go func() {
		<-ctx.Done()
		close(ch)
	}()

	result, err := p.Client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl: aws.String(queueURL),
	})
	if err != nil {
		return err
	}

	for _, message := range result.Messages {
		ch <- &message
	}

	return nil
}

func (p *Pubsub) PublishMessage(ctx context.Context, queueURL string, message string) error {
	_, err := p.Client.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: aws.String(message),
		QueueUrl:    aws.String(queueURL),
	})
	if err != nil {
		return err
	}

	return nil
}
