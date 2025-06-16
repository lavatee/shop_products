package service

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type RabbitMQProducer struct {
	Channel *amqp.Channel
}

func ConnectRabbitMQ(host string, port string, user string, password string) (*amqp.Channel, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s", user, password, host, port))
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	return ch, err
}

func NewRabbitMQProducer(ch *amqp.Channel) *RabbitMQProducer {
	return &RabbitMQProducer{Channel: ch}
}

type MQMessage struct {
	MessageData interface{} `json:"eventInfo"`
	EventId     string      `json:"eventId"`
}

func (p *RabbitMQProducer) SendMessage(queue string, message interface{}) error {
	mqQueue, err := p.Channel.QueueDeclare(queue, false, false, false, false, nil)
	if err != nil {
		return nil
	}
	eventId, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}
	messageBody := MQMessage{
		MessageData: message,
		EventId:     eventId,
	}
	jsonBody, err := json.Marshal(messageBody)
	if err != nil {
		return err
	}
	for {
		if err := p.Channel.Publish("", mqQueue.Name, false, false, amqp.Publishing{ContentType: "application/json", Body: jsonBody}); err != nil {
			continue
		}
		logrus.Infof("Message with EventId \"%s\" was published", eventId)
		break
	}
	return nil
}

func (p *RabbitMQProducer) GetConfirmersAmount(queue string) (int, error) {
	queueInfo, err := p.Channel.QueueInspect(queue)
	if err != nil {
		return 0, err
	}
	return queueInfo.Consumers, nil
}
