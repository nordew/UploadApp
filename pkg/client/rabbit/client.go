package rabbit

import "github.com/streadway/amqp"

func NewRabbitClient(url string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
