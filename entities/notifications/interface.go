package notifications

import (
	"github.com/gin-gonic/gin"	
	"github.com/Shopify/sarama"
)
type Service interface {
	GetNotificationMessages (context *gin.Context)
	RegisterNotificationMessage (context *gin.Context, producer sarama.SyncProducer) (bool, error) 
}

type Repository interface {
}