package jaeger

import (
	"github.com/Shopify/sarama"
	"github.com/opentracing/opentracing-go"
)

// InjectKafkaHeaders injects span into message header
func InjectKafkaHeaders(span opentracing.Span, headers *[]sarama.RecordHeader) error {
	carrier := (*kafkaHeaderCarrierWriter)(headers)
	return span.Tracer().Inject(span.Context(), opentracing.TextMap, carrier)
}

// ExtractKafkaHeaders extracts span from message header
func ExtractKafkaHeaders(headers []*sarama.RecordHeader) (opentracing.SpanContext, error) {
	carrier := kafkaHeadersCarrierReader(headers)
	return opentracing.GlobalTracer().Extract(opentracing.TextMap, carrier)
}

type kafkaHeadersCarrierReader []*sarama.RecordHeader
type kafkaHeaderCarrierWriter []sarama.RecordHeader

// ForeachKey conforms to the TextMapReader interface.
func (c kafkaHeadersCarrierReader) ForeachKey(handler func(key, val string) error) error {
	for _, val := range c {
		if err := handler(string(val.Key), string(val.Value)); err != nil {
			return err
		}
	}
	return nil
}

// Set implements Set() of opentracing.TextMapWriter.
func (c *kafkaHeaderCarrierWriter) Set(key, val string) {
	headers := (*[]sarama.RecordHeader)(c)
	header := sarama.RecordHeader{
		Key:   []byte(key),
		Value: []byte(val),
	}
	*headers = append(*headers, header)
}
