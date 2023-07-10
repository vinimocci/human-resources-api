package utils

import (
	"net/http"
	"github.com/Shopify/sarama"
	"github.com/gorilla/websocket"
)

var NetUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func CreateKafkaConsumer()(sarama.Consumer, error){
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

func CreateKafkaPartitionConsumer(consumer sarama.Consumer)(sarama.PartitionConsumer, error){
	partitionConsumer, partitionErr := consumer.ConsumePartition("notifications", 0, sarama.OffsetNewest)
	if partitionErr != nil {
		return nil, partitionErr
	}

	return partitionConsumer, nil
}