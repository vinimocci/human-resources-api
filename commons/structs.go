package commons

import "time"

type User struct {

	ID 		  			 int64 		`json:"id"`
	Name 				 string 		`json:"name"`
	UserType  			 UserType	`json:"user_type"`
	Email 	  			 string 		`json:"email"`
	Document 	  		 string 		`json:"document"`
	Birthday  			 time.Time 	`json:"birthday"`
	Password			 string 		`json:"password"`
	Address   			*string		`json:"address"`
	AddressComplement   *string 		`json:"addresscomplement"`
	AddressNeighborhood *string		`json:"addressneighborhood"`
	AddressCity			*string		`json:"addresscity"`
	AddressState		*string		`json:"addressstate"`
	AddressZipCode		*string		`json:"addresszipcode"`
	CreatedAt 			 time.Time	`json:"createdat"`
	UpdatedAt 			 time.Time	`json:"updatedat"`
}

type UserType struct {
	ID 			int64 `json:"id"`
	Description string `json:"description"`
}

type Response struct {
	ResponseCode    int
	ResponseMessage string
}
