package user

import (
	"fmt"
	"time"
	"strconv"
	"github.com/gin-gonic/gin"

	commons "human-resources-api/commons"
	structs "human-resources-api/commons/structs"
)

type service struct {
	repo Repository
}

func NewService (repo Repository) Service {
	return &service{repo}
}

func (s *service) PostUser(context *gin.Context) (bool, error) {

	user, err := validateUserData(context)
	if err != nil {
		return false, err
	}
	
	result, resultErr := s.repo.PostUser(context, user)
	if resultErr != nil {
		return false, resultErr
	}

	return result, nil
}

func validateUserData (context *gin.Context) (*structs.User, error) {
	user := &structs.User{}

	var userType int64

	if context.Request.PostFormValue("usertype") != commons.EmptyResult {
		convertedUserType, usrTypeErr := strconv.ParseInt(context.Request.PostFormValue("usertype"), 10, 64) 
		if usrTypeErr != nil {
			return nil, fmt.Errorf("error parsing user type to int64 type")
		}else{
			userType = convertedUserType
		}

		user.UserType.ID = userType
	}else {
		return nil, fmt.Errorf("missing user type from user")
	}

	if context.Request.PostFormValue("name") != commons.EmptyResult {
		user.Name = context.Request.PostFormValue("name")
	}else {
		return nil, fmt.Errorf("missing name from user")
	}

	if context.Request.PostFormValue("document") != commons.EmptyResult {
		user.Document = context.Request.PostFormValue("document")
	}else {
		return nil, fmt.Errorf("missing document from user")
	}

	if context.Request.PostFormValue("email") != commons.EmptyResult {
		user.Email = context.Request.PostFormValue("email")
	}else {
		return nil, fmt.Errorf("missing email from user")
	}

	if context.Request.PostFormValue("password") != commons.EmptyResult {
		user.Password = context.Request.PostFormValue("password")
	}else {
		return nil, fmt.Errorf("missing password from user")
	}

	if userType == commons.CANDIDATE  {
		if context.Request.PostFormValue("birthday") != commons.EmptyResult {
			parsedUserBirthday, err := time.Parse("01-02-2006 00:00:00", context.Request.PostFormValue("birthday")) 
			if err != nil {
				return nil, err
			}
	
			user.Birthday = parsedUserBirthday
		}else {
			return nil, fmt.Errorf("missing birthday from user")
		}
	}

	if userType == commons.COMPANY {

		if context.Request.PostFormValue("address") != commons.EmptyResult {
			reqAddres := context.Request.PostFormValue("address")
			user.Address = reqAddres
		}else {
			return nil, fmt.Errorf("missing main address from company")
		}

		if context.Request.PostFormValue("addresscomplement") != commons.EmptyResult {
			reqAddressComplement := context.Request.PostFormValue("addresscomplement")
			user.AddressComplement = reqAddressComplement
		}

		if context.Request.PostFormValue("addressneighborhood") != commons.EmptyResult {
			reqAddressNeighborhood := context.Request.PostFormValue("addressneighborhood")
			user.AddressNeighborhood = reqAddressNeighborhood
		}else {
			return nil, fmt.Errorf("missing address neighborhood from company")
		}

		if context.Request.PostFormValue("addresscity") != commons.EmptyResult {
			reqAddressCity := context.Request.PostFormValue("addresscity")
			user.AddressCity = reqAddressCity
		}else {
			return nil, fmt.Errorf("missing address city from company")
		}

		if context.Request.PostFormValue("addressstate") != commons.EmptyResult {
			reqAddressState := context.Request.PostFormValue("addressstate")
			user.AddressState = reqAddressState
		}else {
			return nil, fmt.Errorf("missing address state from company")
		}

		if context.Request.PostFormValue("addresszipcode") != commons.EmptyResult {
			reqAddressZipCode := context.Request.PostFormValue("addresszipcode")
			user.AddressZipCode = reqAddressZipCode
		}else {
			return nil, fmt.Errorf("missing address zip code from company")
		}
		
	}

	return user, nil
}

func (s *service) GetUserInfoByID (context *gin.Context) (*structs.UserInfo, error) {
	var userID int64

	userIdFromReq := context.Param("id")

	if userIdFromReq != commons.EmptyResult {
		convertedUserID, convErr := strconv.ParseInt(userIdFromReq, 10, 64) 
		if convErr != nil {
			return nil, fmt.Errorf("error parsing user ID to int64 type")
		}

		userID = convertedUserID
	}else {
		return nil, fmt.Errorf("missing user id")
	}

	userData, usrDataErr := s.repo.GetUserInfoByID(userID) 
	if usrDataErr != nil {
		return nil, usrDataErr
	}

	return userData, nil
}