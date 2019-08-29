package jaeger

import (
	"github.com/opentracing/opentracing-go"
	"github.com/streadway/amqp"
)

// InjectRabbitMQHeaders injects span into message header
func InjectRabbitMQHeaders(span opentracing.Span, headers amqp.Table) error {
	carrier := rabbitmqHeadersCarrier(headers)
	return span.Tracer().Inject(span.Context(), opentracing.TextMap, carrier)
}

// ExtractRabbitMQHeaders extracts span from message header
func ExtractRabbitMQHeaders(headers amqp.Table) (opentracing.SpanContext, error) {
	carrier := rabbitmqHeadersCarrier(headers)
	return opentracing.GlobalTracer().Extract(opentracing.TextMap, carrier)
}

type rabbitmqHeadersCarrier map[string]interface{}

// ForeachKey conforms to the TextMapReader interface.
func (c rabbitmqHeadersCarrier) ForeachKey(handler func(key, val string) error) error {
	for k, val := range c {
		v, ok := val.(string)
		if !ok {
			continue
		}
		if err := handler(k, v); err != nil {
			return err
		}
	}
	return nil
}

// Set implements Set() of opentracing.TextMapWriter.
func (c rabbitmqHeadersCarrier) Set(key, val string) {
	c[key] = val
}
