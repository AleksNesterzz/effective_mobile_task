package cerrors

import (
	"errors"
)

var (
	ErrDbConnect = errors.New("error connecting to database")
	ErrMigration = errors.New("error during migration")
	ErrLoadEnv   = errors.New("error loading .env file")
)
