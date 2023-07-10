package utils

import (
	"os"
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
	config, tomlErr := GetCurrentEnvironment(os.Getenv("ENVIRONMENT"))
    if tomlErr != nil {
		panic("error loading config file")
    }
	
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{config.Get("kafka.host").(string)}, saramaConfig)
	if err != nil {
		return nil, err
	}

	return consumer, nil
}

func CreateKafkaSinglePartitionConsumer(partition int32, topic string, consumer sarama.Consumer)(sarama.PartitionConsumer, error){
	partitionConsumer, partitionErr := consumer.ConsumePartition(topic, partition, sarama.OffsetNewest)
	if partitionErr != nil {
		return nil, partitionErr
	}

	return partitionConsumer, nil
}