package sqlitedb

import (
	"errors"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

var ErrUniqueConstraint = errors.New("violates unique constraint")

func GetSqliteError(err error) error {
	if sqe, ok := err.(*sqlite.Error); ok {
		if sqe.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return ErrUniqueConstraint
		}
		return sqe
	}

	return nil
}
