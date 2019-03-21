package cas

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

// DbStore implements the TicketStore interface storing ticket data in memory.
type DbStore struct {
	db *sqlx.DB
}

const readQuery = "SELECT id, username, proxy_granting_ticket, is_new_login, is_remembered_login FROM authentication_responses WHERE id = $1"
const writeQuery = "INSERT INTO authentication_responses (id, username, proxy_granting_ticket, authentication_date, is_new_login, is_remembered_login) VALUES ($1, $2, $3, $4, $5, $6)"
const deleteQuery = "DELETE FROM authentication_responses WHERE username IN (SELECT username FROM authentication_responses WHERE id = $1)"
const clearQuery = "DELETE FROM authentication_responses WHERE 1 = 1"

// Read returns the AuthenticationResponse for a ticket
func (s *DbStore) Read(id string) (*AuthenticationResponse, error) {
	rows, err := s.db.Query(readQuery, id)
	if err != nil {
		return nil, ErrInvalidTicket
	}

	for rows.Next() {
		var id, user, proxy_granting_ticket string
		var is_new_login, is_remembered_login sql.NullBool
		err := rows.Scan(&id, &user, &proxy_granting_ticket, &is_new_login, &is_remembered_login)
		if err != nil {
			return nil, ErrInvalidTicket
		}

		return &AuthenticationResponse{
			User:                user,
			ProxyGrantingTicket: proxy_granting_ticket,
		}, nil
	}

	return nil, ErrInvalidTicket
}

// Write stores the AuthenticationResponse for a ticket
func (s *DbStore) Write(id string, ticket *AuthenticationResponse) error {
	_, err := s.db.Exec(writeQuery, id, ticket.User, ticket.ProxyGrantingTicket, ticket.AuthenticationDate, ticket.IsNewLogin, ticket.IsRememberedLogin)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes the AuthenticationResponse for a ticket
func (s *DbStore) Delete(id string) error {
	_, err := s.db.Exec(deleteQuery, id)
	if err != nil {
		return err
	}
	return nil
}

// Clear removes all ticket data
func (s *DbStore) Clear() error {
	_, err := s.db.Exec(clearQuery)
	if err != nil {
		return err
	}
	return nil
}
