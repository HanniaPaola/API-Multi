package infraestructure

import (
	"database/sql"
	"fmt"
)

type TokenRepositoryImpl struct {
	db *sql.DB
}

func NewTokenRepositoryImpl(db *sql.DB) *TokenRepositoryImpl {
	return &TokenRepositoryImpl{db: db}
}

func (r *TokenRepositoryImpl) SaveToken(token string) (int64, error) {
	query := "INSERT INTO tokens (token) VALUES (?)"
	result, err := r.db.Exec(query, token)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func (r *TokenRepositoryImpl) GetAllTokens() ([]string, error) {
	rows, err := r.db.Query("SELECT token FROM tokens")
	if err != nil {
		return nil, fmt.Errorf("error obteniendo tokens de usuarios: %v", err)
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			return nil, fmt.Errorf("error leyendo token de usuario: %v", err)
		}
		tokens = append(tokens, token)
	}

	return tokens, nil
}