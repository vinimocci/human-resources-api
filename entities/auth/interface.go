package auth

import (
	"context"
	"github.com/gin-gonic/gin"
	
	structs "human-resources-api/commons/structs"
)

//Regras de Negócios
type Service interface {
	SignIn(context *gin.Context) (bool, error)
	verifyIfEmailExists (email string) (bool, error)
}
//Repositórios
type Repository interface {
	verifyIfEmailExists (email string)(bool, error)
	SignIn (context context.Context, user *structs.AuthUser) (bool, error)
}
