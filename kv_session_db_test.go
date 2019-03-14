package cas

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestShouldReadKvSessionDb(t *testing.T) {
	db, mock, err := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer sqlxDB.Close()
	defer db.Close()

	columns := []string{"value"}
	rows := sqlmock.NewRows(columns).AddRow("ticket-id-01")

	mock.ExpectQuery("SELECT .+ FROM key_value_sessions .+").WithArgs("cookie-id-01").WillReturnRows(rows)

	kvSession := &KvSessionDb{db: sqlxDB}

	// now we execute our method
	if _, ok := kvSession.Read("cookie-id-01"); ok != true {
		t.Errorf("error was not expected while reading from KvSessionDb: %v", ok)
	}
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestShouldWriteKvSessionDb(t *testing.T) {
	db, mock, err := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer sqlxDB.Close()
	defer db.Close()

	mock.ExpectExec(`INSERT INTO key_value_sessions \(key, value\) VALUES \(.*, .*\)`).WithArgs("cookie-id-01", "ticket-id-01").WillReturnResult(sqlmock.NewResult(1, 1))

	kvSession := &KvSessionDb{db: sqlxDB}

	// now we execute our method
	kvSession.Write("cookie-id-01", "ticket-id-01")
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestShouldDeleteKvSessionDb(t *testing.T) {
	db, mock, err := sqlmock.New()
	sqlxDB := sqlx.NewDb(db, "postgres")

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer sqlxDB.Close()
	defer db.Close()

	mock.ExpectExec(`DELETE FROM key_value_sessions WHERE key = .+`).WithArgs("cookie-id-01").WillReturnResult(sqlmock.NewResult(1, 1))

	kvSession := &KvSessionDb{db: sqlxDB}

	// now we execute our method
	kvSession.Delete("cookie-id-01")
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
