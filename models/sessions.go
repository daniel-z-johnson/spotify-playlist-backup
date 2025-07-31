package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
)

type Sessions struct {
	Db *sql.DB
}

func (s *Sessions) WriteValue(session string, key string, value string) error {
	hashSessionBytes := sha256.Sum256([]byte(session))
	hashHex := hex.EncodeToString(hashSessionBytes[:])
	_, err := s.Db.Exec("INSERT INTO sessions (sessions_hash, key, value) VALUES (?, ?, ?)", hashHex, key, value)
	if err != nil {
		return err
	}
	return nil
}

func (s *Sessions) ReadValue(session string, key string) (string, error) {
	hashSessionBytes := sha256.Sum256([]byte(session))
	hashHex := hex.EncodeToString(hashSessionBytes[:])
	var value string
	err := s.Db.QueryRow("SELECT value FROM sessions WHERE sessions_hash = ? AND key = ?", hashHex, key).Scan(&value)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil // No value found for this key
		}
		return "", err // Other error
	}
	return value, nil
}
