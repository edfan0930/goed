package rabbit

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"git.cchntek.com/Cypress/sts/module/rabbit/consume"

	"github.com/streadway/amqp"
)

type (
	Body struct {
		A int
	}
)

func Send() {
	fmt.Println("in send")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	fmt.Println("conn", conn)
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"ststest", // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := consume.Order{
		OrderID:    "orderID",
		GameCode:   "gamecode",
		GameHall:   "gamehall",
		GameType:   "gametype",
		Platform:   "web",
		Account:    "test01",
		Owner:      "ownerID",
		Parent:     "parentID",
		GameToken:  "gametoken",
		Wins:       10,
		Bets:       5,
		RoundID:    "roundid",
		IndexID:    "indexid",
		CreateTime: time.Now(),
	}
	body.OrderID = "ooxx"
	var a int
	for {

		a++
		playerID := strconv.Itoa(a)
		body.PlayerID = playerID
		//body.A++
		bod, _ := json.Marshal(body)
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        bod,
			})
		log.Printf(" [x] Sent %s", body)
		failOnError(err, "Failed to publish a message")
	}

}
