package auth

import (
	"fmt"
	"github.com/gin-gonic/gin"

	commons "web-solutions-api/commons"
	structs "web-solutions-api/commons/structs"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) SignIn (context *gin.Context)(*structs.UserInfo, error){
	userInfo, usrInfoErr := validateUserLoginData(context)
	if usrInfoErr != nil {
		return nil, usrInfoErr
	}

	hasEmailOnSytem, emailErr := s.verifyIfEmailExists(userInfo.Email)
	if emailErr != nil{
		return nil, emailErr
	}

	if !hasEmailOnSytem {
		return nil, emailErr
	}

	userData, PassMatchErr := s.repo.VerifyIfPasswordMatches(userInfo.Email, userInfo.Password)
	if PassMatchErr != nil{
		return nil, PassMatchErr
	}

	return userData, nil
}

func validateUserLoginData (context *gin.Context) (*structs.AuthUser, error) {
	userLoginInfo := &structs.AuthUser{}

	if context.Request.PostFormValue("email")!= commons.EmptyResult {
		userLoginInfo.Email = context.Request.PostFormValue("email")
	}else{
		return nil, fmt.Errorf("missing user email")
	}

	if context.Request.PostFormValue("password")!= commons.EmptyResult{
		userLoginInfo.Password = context.Request.PostFormValue("password")
	}else{
		return nil, fmt.Errorf("missing user password")
	}

	return userLoginInfo, nil
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