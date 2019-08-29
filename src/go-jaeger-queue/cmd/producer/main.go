package main

import (
	"context"
	"fmt"
	"go-jaeger-queue/utils/jaeger"
	"os"
	"os/signal"
	"time"

	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go/ext"
)

func main() {
	_, err := jaeger.Init("localhost", 5775, "producer")
	if err != nil {
		fmt.Println(err)
		return
	}

	var (
		brokers = []string{"localhost:9092"}
		topic   = "test"
	)
	config := sarama.NewConfig()
	config.Version = sarama.V2_3_0_0
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	asyncProducer, err := sarama.NewAsyncProducer(brokers, config)
	if err == nil {
		go func() {
			for {
				select {
				case producerErr := <-asyncProducer.Errors():
					fmt.Println(fmt.Sprintf("Failed to send to topic %v with error %v", producerErr.Msg.Topic, producerErr.Err))
					asyncProducer.Input() <- producerErr.Msg
				case message := <-asyncProducer.Successes():
					fmt.Println(fmt.Sprintf("Sent key=%v value=%v to topic=%v", message.Key, message.Value, message.Topic))
				}
			}
		}()
	}

	// Create a new span
	span := jaeger.Start(context.Background(), ">Send", ext.SpanKindProducer)

	message := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder("1"),
		Value: sarama.StringEncoder("2"),
	}

	// Inject span into message header
	if err := jaeger.InjectKafkaHeaders(span, &message.Headers); err != nil {
		fmt.Println(fmt.Sprintf("%v-Failed to inject span with error %v", message.Key, err))
	}

	asyncProducer.Input() <- message

	// TODO: continue doing something after sending queue successfully
	time.Sleep(time.Second)

	// Finish span
	jaeger.Finish(span, err)

	signals := make(chan os.Signal, 1)
	shutdown := make(chan bool, 1)

	signal.Notify(signals, os.Interrupt)
	go func() {
		<-signals

		// TODO: Release resources
		asyncProducer.Close()

		shutdown <- true
	}()
	<-shutdown
}
