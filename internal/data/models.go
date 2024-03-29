package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Movies      MovieModel
	Tokens      TokenModel
	Users       UserModel
	Permissions PermissionModel // Add a new Permissions field.
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies:      MovieModel{DB: db},
		Permissions: PermissionModel{DB: db}, // Initialize a new PermissionModel instance
		Tokens:      TokenModel{DB: db},
		Users:       UserModel{DB: db},
	}
}
