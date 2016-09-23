package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/AlbinOS/kafka-go-test/models"
	cluster "github.com/bsm/sarama-cluster"
	avro "github.com/elodina/go-avro"
)

func main() {
	config := cluster.NewConfig()
	config.Group.Return.Notifications = true
	config.Consumer.Return.Errors = true
	kafkaClient, err := cluster.NewClient([]string{"localhost:9092"}, config)
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
			consumed++

			log.Printf("Consumed message time:topic:partition:offset:message %s:%s:%d:%d:%v:%d\n", msg.Timestamp, msg.Topic, msg.Partition, msg.Offset, msg.Value, len(msg.Value))

			// Create a new TestRecord to decode data into
			decodedRecord := new(models.Device)
			// Avro stuff
			reader := avro.NewSpecificDatumReader()
			// SetSchema must be called before calling Read
			reader.SetSchema(decodedRecord.Schema())

			// Create a new Decoder with a given buffer
			decoder := avro.NewBinaryDecoder(msg.Value)

			// Read data into a given record with a given Decoder.
			err = reader.Read(decodedRecord, decoder)
			if err != nil {
				panic(err)
			}

			log.Println(decodedRecord)
		case notification := <-consumer.Notifications():
			for k, v := range notification.Claimed {
				fmt.Printf("Claimed : Topic %s - partition %v\n", k, v)
			}
			for k, v := range notification.Current {
				fmt.Printf("Current : Topic %s - partition %v\n", k, v)
			}
			for k, v := range notification.Released {
				fmt.Printf("Released : Topic %s - partition %v\n", k, v)
			}
		case error := <-consumer.Errors():
			fmt.Printf("Error : %s\n", error)
		case <-signals:
			break ConsumerLoop
		}
	}

	log.Printf("Consumed: %d\n", consumed)
}
