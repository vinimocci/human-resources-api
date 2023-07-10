package notifications

import (
	"log"
	"github.com/gin-gonic/gin"
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
	conn, upgradeErr := utils.NetUpgrader.Upgrade(c.Writer, c.Request, nil)
	if upgradeErr != nil {
		log.Println("WebSocket upgrade error:", upgradeErr)
		return
	}

	consumer, consumerErr := utils.CreateKafkaConsumer()
	if consumerErr != nil {
		log.Println("Kafka consumer creation error:", consumerErr)
		return
	}
	defer consumer.Close()

	partitionConsumer, partitionErr := utils.CreateKafkaSinglePartitionConsumer(0, "notifications", consumer)
	if partitionErr != nil {
		log.Println("Kafka partition consumer creation error:", partitionErr)
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