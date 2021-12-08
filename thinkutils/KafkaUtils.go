package thinkutils

import (
	"context"
	"github.com/segmentio/kafka-go"
	"strings"
)

type kafkautils struct {
}

type OnMsgCallback func(message kafka.Message)

func (this kafkautils) StartConsumer(szUrl string, szTopic string, szGroupId string, callback OnMsgCallback) {
	go func() {
		brokers := strings.Split(szUrl, ",")
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			GroupID:  szGroupId,
			Topic:    szTopic,
			MinBytes: 10e3, // 10KB
			MaxBytes: 10e6, // 10MB
		})

		defer reader.Close()

		for {
			m, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Error(err.Error())
				continue
			}

			if callback != nil {
				callback(m)
			}
			//fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		}
	}()
}

func (this kafkautils) SendMsg(szUrl string, szTopic string, data []byte) {
	//brokers := strings.Split(kafkaURL, ",")
	writer := &kafka.Writer{
		Addr:     kafka.TCP(szUrl),
		Topic:    szTopic,
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	msg := kafka.Message{
		Key:   []byte("1"),
		Value: data,
	}

	err := writer.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Error(err.Error())
	}
}
