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
<<<<<<< HEAD
	Permissions PermissionModel // Add a new Permissions field.
	Tokens      TokenModel
	Users       UserModel
=======
	Tokens      TokenModel
	Users       UserModel
	Permissions PermissionModel // Add a new Permissions field.
>>>>>>> 8538e24fce30ecbdc2156c185f70e0f24b335c65
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies:      MovieModel{DB: db},
<<<<<<< HEAD
		Permissions: PermissionModel{DB: db}, // Initialize a new PermissionModel instance.
=======
		Permissions: PermissionModel{DB: db}, // Initialize a new PermissionModel instance
>>>>>>> 8538e24fce30ecbdc2156c185f70e0f24b335c65
		Tokens:      TokenModel{DB: db},
		Users:       UserModel{DB: db},
	}
}
