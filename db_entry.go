package odbc

import "time"

type DBEntry struct {
	InsertedAt time.Time `db:"inserted_at" json:"inserted_at"`
}
