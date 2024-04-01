package events

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaConfig struct {
	Servers      string
	Topic        string
	SaslUserName string
	SaslPassword string
}

type EventsKafkaMq struct {
	producer *kafka.Producer
	consumer *kafka.Consumer

	topic string
	key   string
}

func NewKafkaEventProducer(kfc *KafkaConfig) (*EventsKafkaMq, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kfc.Servers,
		"batch.size":        1674380,
		"compression.type":  "gzip",
		"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "SCRAM-SHA-256",
		"sasl.username":     kfc.SaslUserName,
		"sasl.password":     kfc.SaslPassword,
	})
	if err != nil {
		return nil, err
	}

	return &EventsKafkaMq{
		producer: p,
		topic:    kfc.Topic,
		key:      "events",
	}, nil
}

// Replace with WorkerPool, maybe
func (mq *EventsKafkaMq) PushEvent(ctx context.Context, evs ...*Event) error {
	log.Info().Msgf("sending %d messages", len(evs))

	for _, ev := range evs {
		value, err := json.Marshal(ev)
		if err != nil {
			return err
		}

		err = mq.producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &mq.topic,
				Partition: kafka.PartitionAny,
			},
			Value:   value,
			Headers: []kafka.Header{{Key: mq.key}},
		}, nil)
		if err != nil {
			log.Error().Err(err).Msg("")
			return err
		}
	}

	// for mq.producer.Flush(10000) > 0 {
	// 	log.Info().Msg("Still waiting to flush outstanding messages")
	// }

	return nil
}

type EventProducerResult struct {
	Err     error
	Message string
}

func (mq *EventsKafkaMq) GetResults(ctx context.Context) chan EventProducerResult {
	resultsChan := make(chan EventProducerResult)

	go func() {
		defer close(resultsChan)

		for {
			select {
			case <-ctx.Done():
				return
			case e := <-mq.producer.Events():
				message := e.(*kafka.Message)
				result := EventProducerResult{
					Err: message.TopicPartition.Error,
				}
				var msg string

				if message.TopicPartition.Error != nil {
					msg = fmt.Sprintf("failed to deliver message: %v\n", message.TopicPartition)
				} else {
					msg = fmt.Sprintf("delivered to topic %s [%d] at offset %v\n",
						*message.TopicPartition.Topic,
						message.TopicPartition.Partition,
						message.TopicPartition.Offset)
				}

				result.Message = msg
				resultsChan <- result
			}
		}
	}()

	return resultsChan
}

func NewKafkaEventConsumer(kfc *KafkaConfig) (*EventsKafkaMq, error) {
	p, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kfc.Servers,
		"batch.size":        1674380,
		"compression.type":  "gzip",
		"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "SCRAM-SHA-256",
		"sasl.username":     kfc.SaslUserName,
		"sasl.password":     kfc.SaslPassword,
		"group.id":          fmt.Sprintf("%s-group", kfc.Topic),
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	return &EventsKafkaMq{
		consumer: p,
		topic:    kfc.Topic,
		key:      "events",
	}, nil
}

type EventsConsumerResultChan struct {
	Message []byte
	Topic   string
	Err     error
}

func (e *EventsKafkaMq) Consume(ctx context.Context, interval time.Duration) (chan EventsConsumerResultChan, error) {
	if err := e.consumer.SubscribeTopics([]string{e.topic}, nil); err != nil {
		return nil, err
	}

	resultsChan := make(chan EventsConsumerResultChan)

	go func() {
		defer close(resultsChan)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				message, err := e.consumer.ReadMessage(interval)
				if err != nil {
					if kafkaErr, ok := err.(kafka.Error); ok {
						if kafkaErr.Code() == kafka.ErrTimedOut {
							continue
						}
					}

					log.Error().Err(err).Msg("failed to read messages")
					resultsChan <- EventsConsumerResultChan{
						Err: err,
					}

					continue
				}

				resultsChan <- EventsConsumerResultChan{
					Message: message.Value,
					Topic:   *message.TopicPartition.Topic,
					Err:     nil,
				}
			}
		}
	}()

	return resultsChan, nil
}
