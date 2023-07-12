package notifications

import (
	"fmt"
	"log"
	"web-solutions-api/utils"
	"github.com/gin-gonic/gin"
	"github.com/Shopify/sarama"
	"github.com/gorilla/websocket"

	commons "web-solutions-api/commons"
	structs "web-solutions-api/commons/structs"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func validateNotificationData (context *gin.Context) (*structs.Notification, error) {
	notification := &structs.Notification{}

	if context.Request.PostFormValue("notification") != commons.EmptyResult {
		notification.Description = context.Request.PostFormValue("notification")
	}else{
		return nil, fmt.Errorf("missing notification description")
	}

	if context.Request.PostFormValue("topic") != commons.EmptyResult {
		notification.Topic = context.Request.PostFormValue("topic")
	}else{
		return nil, fmt.Errorf("missing notification topic")
	}

	return notification, nil
}

func (s *service) GetNotificationMessages (c *gin.Context) {
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

func (s *service) RegisterNotificationMessage (context *gin.Context, producer sarama.SyncProducer) (bool, error){
	notification, notErr := validateNotificationData(context)
	if notErr != nil {
		return false, notErr
	}

	kafkaMessage := &sarama.ProducerMessage {
		Topic: notification.Topic,
		Value: sarama.StringEncoder(notification.Description),
	}

	_,_, prdErr := producer.SendMessage(kafkaMessage)
	if prdErr != nil {
		return false, prdErr
	}
	
	return true, nil
}