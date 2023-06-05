package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"
	
	structs "human-resources-api/commons/structs"
)

type service struct {
	repo Repository
}


func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service)SignIn (c *gin.Context)(bool, error){
	auth := &structs.User{}

	if c.Request.PostFormValue("email")!= "" {
		auth.Email = c.Request.PostFormValue("email")
	}else{
		return false, fmt.Errorf("missing user email")
	}

	if c.Request.PostFormValue("password")!= ""{
		auth.Password = c.Request.PostFormValue("password doesn't match")
	}else{
		return false, fmt.Errorf("missing user password")
	}

	hasEmailOnSytem, emailErr := s.verifyIfEmailExists(auth.Email)
	if emailErr != nil{
		return false, emailErr
	}

	if !hasEmailOnSytem {
		return false, emailErr
	}

	// verificar a logica da senha, se bate com o email cadastrado e senha antes inserida

	return true, nil
}

func (s *service) verifyIfEmailExists (email string) (bool, error){
	hasResult, rstErr := s.repo.verifyIfEmailExists(email)
	if rstErr != nil {
		return false, rstErr
	}

	if !hasResult {
		return false, fmt.Errorf("user email not found on system")
	}

	return true, nil
}