package cas

import (
	"github.com/jmoiron/sqlx"
)

// KvSessionDb
type KvSessionDb struct {
	db *sqlx.DB
}

const readSessionQuery = "SELECT value FROM key_value_sessions WHERE key = $1"
const writeSessionQuery = `INSERT INTO key_value_sessions (key, value) VALUES ($1, $2)`
const deleteSessionQuery = "DELETE FROM key_value_sessions WHERE key = $1"

func (s *KvSessionDb) Read(key string) (string, bool) {
	rows, err := s.db.Query(readSessionQuery, key)
	if err != nil {
		return "", false
	}

	for rows.Next() {
		var value string
		err := rows.Scan(&value)
		if err != nil {
			return "", false
		}

		return value, true
	}

	return "", false
}

func (s *KvSessionDb) Write(key, value string) {
	s.db.Exec(writeSessionQuery, key, value)
}

func (s *KvSessionDb) Delete(key string) {
	s.db.Exec(deleteSessionQuery, key)
}
