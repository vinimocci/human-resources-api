package auth

import (
	"database/sql"
	"fmt"
	"human-resources-api/commons/structs"

	commons "human-resources-api/commons"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db}
}

func (r *repository) verifyIfEmailExists(email string)(bool, error){
	query := `SELECT * FROM users.users usr WHERE usr.email = ?`

	transaction, trsErr := r.db.Begin()
	if trsErr != nil {
		return false, trsErr
	}

	verifyEmailStmt, stmtErr := transaction.Prepare(query)
	if stmtErr != nil {
		return false, trsErr
	}
	defer verifyEmailStmt.Close()

	rows, rstErr := verifyEmailStmt.Query(email)
	if rstErr != nil {
		return false, rstErr
	}

	var totalResults int64 = 0
	
	for rows.Next() {
		totalResults = commons.HasResults
	}

	if totalResults == commons.EmptyEmailResult {
		return false, nil
	}

	return true, nil
}

func (r *repository) VerifyIfPasswordMatches(email, password string)(*structs.UserInfo, error){
	query := `SELECT * FROM users.users usr WHERE usr.email = ? AND  usr.password = ?`

	transaction, trsErr := r.db.Begin()
	if trsErr != nil {
		return nil, trsErr
	}

	verifyEmailStmt, stmtErr := transaction.Prepare(query)
	if stmtErr != nil {
		return nil, trsErr
	}
	defer verifyEmailStmt.Close()

	rows, rstErr := verifyEmailStmt.Query(email, password)
	if rstErr != nil {
		return nil, rstErr
	}

	var totalResults int64 = 0

	var currentUser *structs.UserInfo
	
	for rows.Next() {
		totalResults = commons.HasResults

		if err := rows.Scan(
			&currentUser.ID,
			&currentUser.Name,
		); err != nil {
			return nil, err
		}
	}

	if totalResults == commons.EmptyEmailResult {
		return nil, fmt.Errorf("user password don't match")
	}

	return currentUser, nil
}