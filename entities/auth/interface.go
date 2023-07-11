package auth

import (
	"github.com/gin-gonic/gin"
	
	structs "web-solutions-api/commons/structs"
)
type Service interface {
	verifyIfEmailExists (email string) (bool, error)
	SignIn (context *gin.Context) (*structs.UserInfo, error)
}
type Repository interface {
	verifyIfEmailExists (email string)(bool, error)
	VerifyIfPasswordMatches (email, password string)(*structs.UserInfo, error)
}