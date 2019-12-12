package rabbit

import (
	"github.com/streadway/amqp"
	"log"
	"sync"
	"time"
)

type (
	Receiver struct {
		Receive <-chan amqp.Delivery
		Data    [][]byte
		Payload chan [][]byte
		mux     sync.Mutex
	}
)

//NewReceiver 接收queue交付的message , 交付整理過的message
//@d queue交付訊息
func NewReceiver(d <-chan amqp.Delivery) *Receiver {
	return &Receiver{
		Receive: d,
		Payload: make(chan [][]byte),
	}
}

//Run 將資料整理成多筆
//@collection 採集數量
func (r *Receiver) Run(collection uint64) {
	go func() {
		/* var mutex sync.Mutex
		ticker := time.NewTicker(5 * time.Second)
		select {
		case <-ticker.C:
			mutex.Lock()
			data := r.Data
			r.Data = nil
			mutex.Unlock()
			r.Delivery <- data
		case d := <-r.Receive:
			r.Data = append(r.Data, d.Body)
			d.Ack(true)
		} */
		go func() {
			ticker := time.NewTicker(60 * time.Second)
			<-ticker.C
			r.mux.Lock()
			r.Payload <- r.Data
			r.Data = nil
			r.mux.Unlock()

		}()

		for d := range r.Receive {
			r.Data = append(r.Data, d.Body)
			d.Ack(false)
			if d.DeliveryTag%collection == 0 {
				r.Payload <- r.Data
				r.Data = nil
			}
		}
	}()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
