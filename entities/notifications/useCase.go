package notifications

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/Shopify/sarama"
	"human-resources-api/utils"
	"github.com/gorilla/websocket"

)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}


func (s *service) GetNotifications(c *gin.Context) {
	conn, err := utils.NetUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Println("Kafka consumer creation error:", err)
		return
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition("notifications", 0, sarama.OffsetNewest)
	if err != nil {
		log.Println("Kafka partition consumer creation error:", err)
		return
	}
	defer partitionConsumer.Close()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			err := conn.WriteMessage(websocket.TextMessage, msg.Value)
			if err != nil {
				log.Println("WebSocket write error:", err)
				return
			}
		case err := <-partitionConsumer.Errors():
			log.Println("Kafka consumer error:", err.Err)
			return
		}
	}
}