package auth

import (
	"context"
	"database/sql"
	"human-resources-api/commons/structs"

	commons "human-resources-api/commons"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repository{db}
}

func (r * repository) SignIn(context context.Context, auth *structs.AuthUser)(bool, error){
	query := ` SELECT * FROM users.users(email, password)`
	transaction, trsErr := r.db.Begin()
	if trsErr != nil {
		return false, trsErr
	}

	signInUserStmt, stmtErr := transaction.PrepareContext(context, query)
	if stmtErr != nil{
		return false, stmtErr
	}
	defer signInUserStmt.Close()
	_,rstErr := signInUserStmt.Exec(
		auth.Email,
		auth.Password,
	)
	if rstErr != nil{
		transaction.Rollback()
		return false, rstErr
	}
	transaction.Commit()
	return true, nil
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
