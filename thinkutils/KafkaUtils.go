package thinkutils

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"strings"
	"sync"
)

type kafkautils struct {
}

var (
	g_lockKafka      sync.Mutex
	g_mapKafkaWriter map[string]*kafka.Writer
)

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

func (this kafkautils) makeWriter(szUrl string, szTopic string) *kafka.Writer {
	defer g_lockKafka.Unlock()
	g_lockKafka.Lock()

	szConn := fmt.Sprintf("%s/%s", szUrl, szTopic)

	pWriter := g_mapKafkaWriter[szConn]
	if nil == pWriter {
		//brokers := strings.Split(kafkaURL, ",")
		pWriter = &kafka.Writer{
			Addr:     kafka.TCP(szUrl),
			Topic:    szTopic,
			Balancer: &kafka.LeastBytes{},
		}

		g_mapKafkaWriter[szConn] = pWriter
		//defer writer.Close()
	}

	return pWriter
}

func (this kafkautils) SendMsg(szUrl string, szTopic string, data []byte) {
	if nil == g_mapKafkaWriter {
		g_mapKafkaWriter = make(map[string]*kafka.Writer)
	}

	szConn := fmt.Sprintf("%s/%s", szUrl, szTopic)
	pWriter := g_mapKafkaWriter[szConn]
	if nil == pWriter {
		pWriter = this.makeWriter(szUrl, szTopic)
	}

	msg := kafka.Message{
		Key:   []byte("1"),
		Value: data,
	}

	//log.Info("%p %p", g_mapKafkaWriter, pWriter)
	err := pWriter.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Error(err.Error())
	}
}
