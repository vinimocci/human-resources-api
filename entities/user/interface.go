package user

import (
	"context"
	"github.com/gin-gonic/gin"

	structs "human-resources-api/commons/structs"
)


type Service interface {
	PostUser(context *gin.Context) (bool, error)
	GetUserInfoByID(context *gin.Context) (*structs.UserInfo, error)
}

type Repository interface {
	GetUserInfoByID(userID int64) (*structs.UserInfo, error)
	PostUser (context context.Context, user *structs.User) (bool, error)
}