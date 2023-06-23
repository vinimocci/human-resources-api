package user

import (
	"fmt"
	"context"
	"database/sql"

	commons "human-resources-api/commons"
	structs "human-resources-api/commons/structs"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db}
}

func (r *repository) PostUser (context context.Context, user *structs.User) (bool, error) {
	query := `INSERT INTO users.users
	(name, userType, email, document, birthday, password, address, addressComplement, addressNeighborhood, addressCity, addressState, addressZipCode, createdAt, updatedAt)
	VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
	`

	transaction, trsnErr := r.db.Begin()
	if trsnErr != nil {
		return false, trsnErr
	}
	
	postUserStmt, stmtErr := transaction.PrepareContext(context, query)
	if stmtErr != nil{
		return false, stmtErr
	}
	defer postUserStmt.Close()

	var birthdayDate  sql.NullTime

	if user.Birthday.IsZero() {
		birthdayDate.Valid = false
	} else {
		birthdayDate.Time = user.Birthday
		birthdayDate.Valid = true
	}

	_, rstErr := postUserStmt.Exec(
		&user.Name,
		&user.UserType.ID,
		&user.Email,
		&user.Document,
		&birthdayDate,
		&user.Password,
		&user.Address,
		&user.AddressComplement,
		&user.AddressNeighborhood,
		&user.AddressCity,
		&user.AddressState,
		&user.AddressZipCode,
	)
	if rstErr != nil{
		transaction.Rollback()
		return false, rstErr
	}

	transaction.Commit()

	return true, nil
}

func (r *repository) GetUserInfoByID (userID int64) (*structs.UserInfo, error) {

	query := `
			SELECT
			usr.name,
			usr.email,
			usr.document,
			usr.birthday,
			usr.address,
			usr.addressComplement,
			usr.addressNeighborhood,
			usr.addressCity,
			usr.addressState,
			usr.addressZipCode
		FROM
			users.users usr
		WHERE
			usr.id = ?;
	`

	findUsrStmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer findUsrStmt.Close()

	rows, err := findUsrStmt.Query(userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var hasResult bool = false
	user := &structs.UserInfo{}

	for rows.Next() {
		hasResult = commons.HasResult

		if err := rows.Scan(
			&user.Name,
			&user.Email,
			&user.Document,
			&user.Birthday,
			&user.Address,
			&user.AddressComplement,
			&user.AddressNeighborhood,
			&user.AddressCity,
			&user.AddressState,
			&user.AddressZipCode,
		); err != nil {
			return nil, err
		}
	}

	if !hasResult {
		return nil, fmt.Errorf("user not found with given ID")
	}

	return user, nil
}

func (r *repository) UpdateUser (user *structs.User) (bool, error) {
	query := `
	UPDATE
		users.users
	SET
		name = ?,
		password = ?,
		birthday = ?,
		address = ?,
		addressComplement = ?,
		addressNeighborhood = ?,
		addressCity = ?,
		addressState = ?,
		addressZipCode = ?
	WHERE
		id = ?;
	`

	transaction, trsnErr := r.db.Begin()
	if trsnErr != nil {
		return false, trsnErr
	}
	
	postUserStmt, stmtErr := transaction.Prepare(query)
	if stmtErr != nil{
		return false, stmtErr
	}
	defer postUserStmt.Close()

	var birthdayDate  sql.NullTime

	if user.Birthday.IsZero() {
		birthdayDate.Valid = false
	} else {
		birthdayDate.Time = user.Birthday
		birthdayDate.Valid = true
	}

	_, rstErr := postUserStmt.Exec(
		&user.Name,
		&user.Password,
		&birthdayDate,
		&user.Address,
		&user.AddressComplement,
		&user.AddressNeighborhood,
		&user.AddressCity,
		&user.AddressState,
		&user.AddressZipCode,
		&user.ID,
	)
	if rstErr != nil{
		transaction.Rollback()
		return false, rstErr
	}

	transaction.Commit()

	return true, nil
}

func (r *repository) GetUserTypeByID (userID int64) (int64, error) {

	query := `
		SELECT
			usr.userType
		FROM
			users.users usr
		WHERE
			usr.id = ?;
	`

	findUserTypeStmt, err := r.db.Prepare(query)
	if err != nil {
		return 0, err
	}
	defer findUserTypeStmt.Close()

	rows, err := findUserTypeStmt.Query(userID)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var hasResult bool = false
	userType := &structs.UserType{}

	for rows.Next() {
		hasResult = commons.HasResult

		if err := rows.Scan(
			&userType.ID,
		); err != nil {
			return 0, err
		}
	}

	if !hasResult {
		return 0, fmt.Errorf("user not found with given ID")
	}

	return userType.ID, nil
}