package endpoint

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type RabbitMQConsumer struct {
	Channel *amqp.Channel
}

func NewRabbitMQConsumer(ch *amqp.Channel) *RabbitMQConsumer {
	return &RabbitMQConsumer{Channel: ch}
}

func (c *RabbitMQConsumer) ConsumeQueue(queue string, handler func([]byte) error) error {
	mqQueue, err := c.Channel.QueueDeclare(queue, false, false, false, false, nil)
	if err != nil {
		return err
	}
	var forever chan struct{}
	messages, err := c.Channel.Consume(mqQueue.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}
	go func() {
		for {
			for m := range messages {
				if err := handler(m.Body); err != nil {
					logrus.Errorf("Error while consuming '%s' queue: %s", queue, err.Error())
				}
			}
		}
	}()
	<-forever
	return nil
}
