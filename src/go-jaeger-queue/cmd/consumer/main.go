package main

import (
	"context"
	"fmt"
	"go-jaeger-queue/utils/jaeger"
	"os"
	"os/signal"
	"time"

	"github.com/opentracing/opentracing-go"

	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go/ext"
)

func main() {
	_, err := jaeger.Init("localhost", 5775, "consumer")
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
	config.Consumer.Return.Errors = true
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		fmt.Println(err)
		return
	}
	partitions, err := consumer.Partitions(topic)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, partition := range partitions {
		partitionConsumer, err := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
		if nil != err {
			fmt.Println(err)
			return
		}
		go func(topic string, partitionConsumer sarama.PartitionConsumer) {
			for {

				select {
				case consumerError := <-partitionConsumer.Errors():
					fmt.Println(fmt.Sprintf("Failed to receive from %v with error %v", topic, consumerError))

				case msg := <-partitionConsumer.Messages():
					fmt.Println(fmt.Sprintf("Received message with key=%v value=%v topic=%v", string(msg.Key), string(msg.Value), msg.Topic))

					var span opentracing.Span
					// Extract span from message header
					spanCtx, err := jaeger.ExtractKafkaHeaders(msg.Headers)
					if err != nil {
						fmt.Println(fmt.Sprintf("Failed to extract span context with error %v", err))

						// Start a new span
						span = jaeger.Start(context.Background(), ">Receive", ext.SpanKindConsumer)
					} else {
						// Continue with injected spance
						span = jaeger.Continue(spanCtx, ">Receive", ext.SpanKindConsumer)
					}

					// TODO: continue doing something after receving successfully
					time.Sleep(time.Second)

					// Finish span
					jaeger.Finish(span, err)
				}
			}
		}(topic, partitionConsumer)
	}

	signals := make(chan os.Signal, 1)
	shutdown := make(chan bool, 1)

	signal.Notify(signals, os.Interrupt)
	go func() {
		<-signals

		// TODO: Release resources
		_ = consumer.Close()

		shutdown <- true
	}()
	<-shutdown
}
