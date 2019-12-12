package rabbit

import (
	"fmt"

	"github.com/streadway/amqp"
)

type (
	//Rabbbit ...
	Rabbit struct {
		Host          string
		Account       string
		Password      string
		QueueName     string
		WorkOn        bool
		PrefetchCount int
		Conn          *amqp.Connection
		Channel       *amqp.Channel
		Que           amqp.Queue
	}
)

//NewRabbit new instance
func NewRabbit(host, account, password string) *Rabbit {
	return &Rabbit{
		Host:     host,
		Account:  account,
		Password: password,
		WorkOn:   true,
	}
}

//Dial get connection
func (r *Rabbit) Dial() (err error) {
	host := fmt.Sprintf("amqp://%v:%v@%v/", r.Account, r.Password, r.Host)
	r.Conn, err = amqp.Dial(host)
	return
}

//NewChannel ...
func (r *Rabbit) NewChannel() (err error) {
	r.Channel, err = r.Conn.Channel()
	return
}

//QueueDeclare ...
func (r *Rabbit) QueueDeclare(queueName string) (err error) {
	r.Que, err = r.Channel.QueueDeclare(
		r.QueueName, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	return
}

//NewConsumer ...
func (r *Rabbit) NewConsumer() (<-chan amqp.Delivery, error) {

	err := r.Channel.Qos(
		r.PrefetchCount, //prefetch count
		0,               //prefetch size
		false,           //global
	)
	if err != nil {
		return nil, err
	}

	return r.Channel.Consume(
		r.QueueName, // queue
		"",          // consumer
		false,       // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
}

func (r *Rabbit) SetQueueName(qName string) *Rabbit {
	r.QueueName = qName
	return r
}

func (r *Rabbit) SetPrefetchCount(count int) *Rabbit {
	r.PrefetchCount = count
	return r
}
