package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	cluster "github.com/bsm/sarama-cluster"
)

func main() {
	kafkaClient, err := cluster.NewClient([]string{"localhost:9092"}, nil)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := kafkaClient.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	consumer, err := cluster.NewConsumerFromClient(kafkaClient, "0", []string{"important"})
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumed := 0
ConsumerLoop:
	for {
		select {
		case msg := <-consumer.Messages():
			s := string(msg.Value[:len(msg.Value)])
			log.Printf("Consumed message time:topic:partition:offset:message %s:%s:%d:%d:%s\n", msg.Timestamp, msg.Topic, msg.Partition, msg.Offset, s)
			consumed++
			for k, v := range consumer.Subscriptions() {
				fmt.Printf("key[%s] value[%v]\n", k, v)
			}
			time.Sleep(1000 * time.Millisecond)
		case <-signals:
			break ConsumerLoop
		}
	}

	log.Printf("Consumed: %d\n", consumed)
}
