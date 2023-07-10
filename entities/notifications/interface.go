package notifications

import (
	"github.com/gin-gonic/gin"	
)
type Service interface {
	GetNotifications (context *gin.Context) 
}

type Repository interface {
}