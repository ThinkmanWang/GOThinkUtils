package thinkutils

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type kafkautils struct {
}

var (
	// g_mapKafkaWriter: key "url/topic" -> *kafka.Writer。
	// 用 sync.Map 保证并发安全，消除此前"未加锁读 + 加锁写"同一 map 导致的
	// fatal error: concurrent map read and map write（冷启动 + 突发高并发时会崩整个进程）。
	g_mapKafkaWriter sync.Map
)

type OnMsgCallback func(message kafka.Message)

func (this kafkautils) StartConsumer(szUrl string, szTopic string, szGroupId string, callback OnMsgCallback) {
	go func() {
		brokers := strings.Split(szUrl, ",")
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			GroupID:  szGroupId,
			Topic:    szTopic,
			MinBytes: 1,    // 10KB
			MaxBytes: 10e6, // 10MB
			MaxWait:  1 * time.Second,
		})

		defer reader.Close()

		for {
			m, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Error(err.Error())
				continue
			}

			if callback != nil {
				go callback(m)
			}
			//fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		}
	}()
}

func (this kafkautils) getWriter(szUrl, szTopic string) *kafka.Writer {
	szConn := fmt.Sprintf("%s/%s", szUrl, szTopic)

	// 快路径：已存在直接返回（sync.Map 并发读安全）。
	if v, ok := g_mapKafkaWriter.Load(szConn); ok {
		return v.(*kafka.Writer)
	}

	// 慢路径：新建 writer，用 LoadOrStore 保证同一 url/topic 只有一个 writer 生效。
	// 并发下偶尔多建出来的会被丢弃：未写入过的 *kafka.Writer 不持有连接/goroutine，可安全 GC。
	lstUrl := strings.Split(szUrl, ",")
	pWriter := &kafka.Writer{
		Addr:     kafka.TCP(lstUrl...),
		Topic:    szTopic,
		Balancer: &kafka.LeastBytes{},
		Async:    true,
	}

	pActual, _ := g_mapKafkaWriter.LoadOrStore(szConn, pWriter)
	return pActual.(*kafka.Writer)
}

//func (this kafkautils) makeSingleWriter(szUrl, szTopic string) *kafka.Writer {
//
//	//brokers := strings.Split(kafkaURL, ",")
//	lstUrl := strings.Split(szUrl, ",")
//	pWriter := &kafka.Writer{
//		Addr:     kafka.TCP(lstUrl...),
//		Topic:    szTopic,
//		Balancer: &kafka.LeastBytes{},
//		Async:    true,
//	}
//
//	return pWriter
//}

func (this kafkautils) SendMsg(szUrl string, szTopic string, data []byte) {
	this.SendMsgPlus(szUrl, szTopic, "", data)
}

func (this kafkautils) SendMsgPlus(szUrl string, szTopic string, szKey string, data []byte) {
	pWriter := this.getWriter(szUrl, szTopic)

	msg := kafka.Message{
		Value: data,
	}
	if false == StringUtils.IsEmpty(szKey) {
		msg.Key = []byte(szKey)
	}

	// writer 为 Async 模式，WriteMessages 仅把消息放入内存批次队列后立即返回、不阻塞，
	// 因此无需再为每条消息单独 spawn goroutine（旧实现每消息一个 goroutine，突发时无上限暴涨）。
	if err := pWriter.WriteMessages(context.Background(), msg); err != nil {
		log.Error(err.Error())
	}
}

//func (this kafkautils) SendMsgPlus(szUrl string, szTopic string, data []byte) {
//	go func(szUrl string, szTopic string, data []byte) {
//		pWriter := this.makeSingleWriter(szUrl, szTopic)
//		defer pWriter.Close()
//
//		msg := kafka.Message{
//			//Key:   []byte("1"),
//			Value: data,
//		}
//
//		//log.Info("%p %p", g_mapKafkaWriter, pWriter)
//		err := pWriter.WriteMessages(context.Background(), msg)
//		if err != nil {
//			log.Error(err.Error())
//		}
//	}(szUrl, szTopic, data)
//}
