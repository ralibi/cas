package cas

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestShouldReadAuthenticationResponses(t *testing.T) {
	db, mock, err := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer sqlxDB.Close()
	defer db.Close()

	columns := []string{"id", "username", "proxy_granting_ticket", "is_new_login", "is_remembered_login"}
	rows := sqlmock.NewRows(columns).AddRow("ticket-id-01", "username-01", "proxy_granting_ticket-01", false, false)

	mock.ExpectQuery("SELECT .+ FROM authentication_responses .+").WithArgs("ticket-id-01").WillReturnRows(rows)

	dbStore := &DbStore{db: sqlxDB}

	// now we execute our method
	if _, err := dbStore.Read("ticket-id-01"); err != nil {
		t.Errorf("error was not expected while reading from DbStore: %s", err)
	}
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestShouldWriteAuthenticationResponses(t *testing.T) {
	db, mock, err := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer sqlxDB.Close()
	defer db.Close()

	mock.ExpectExec(`INSERT INTO authentication_responses \(id, username, proxy_granting_ticket, authentication_date, is_new_login, is_remembered_login\) VALUES \(.*, .*, .*, .*, .*, .*\)`).WithArgs("ticket-id-01", "user-01", "pgt-01", sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))

	dbStore := &DbStore{db: sqlxDB}
	ticket := &AuthenticationResponse{User: "user-01", ProxyGrantingTicket: "pgt-01"}

	// now we execute our method
	if err := dbStore.Write("ticket-id-01", ticket); err != nil {
		t.Errorf("error was not expected while writing from DbStore: %s", err)
	}
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestShouldDeleteAuthenticationResponses(t *testing.T) {
	db, mock, err := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer sqlxDB.Close()
	defer db.Close()

	mock.ExpectExec(`UPDATE authentication_responses SET deleted_at = NOW\(\) WHERE id = .+ AND deleted_at IS NULL`).WithArgs("ticket-id-01").WillReturnResult(sqlmock.NewResult(1, 1))

	dbStore := &DbStore{db: sqlxDB}

	// now we execute our method
	if err := dbStore.Delete("ticket-id-01"); err != nil {
		t.Errorf("error was not expected while deleting from DbStore: %s", err)
	}
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestShouldClearAuthenticationResponses(t *testing.T) {
	db, mock, err := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer sqlxDB.Close()
	defer db.Close()

	mock.ExpectExec(`UPDATE authentication_responses SET deleted_at = NOW\(\) WHERE deleted_at IS NULL`).WillReturnResult(sqlmock.NewResult(1, 1))

	dbStore := &DbStore{db: sqlxDB}

	// now we execute our method
	if err := dbStore.Clear(); err != nil {
		t.Errorf("error was not expected while clearing from DbStore: %s", err)
	}
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
