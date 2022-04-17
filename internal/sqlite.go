package internal

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteDao struct {
	db *sqlx.DB
}

func NewSqliteDao(dsn string) (*SqliteDao, error) {
	s, err := sqlx.Connect("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	driver, err := sqlite3.WithInstance(s.DB, &sqlite3.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file:///Users/brad/Development/BradLugo/rss-snail/migrations",
		"sqlite3", driver)
	if err != nil {
		return nil, err
	}

	err = m.Up()
	if err != nil {
		s.Close()
		return nil, err
	}

	sd := &SqliteDao{
		db: s,
	}
	return sd, nil
}

func (sd SqliteDao) GetFeeds(User) ([]Feed, error) {
	feeds := []Feed{}
	err := sd.db.Select(&feeds, "SELECT * FROM person ORDER BY first_name ASC")
	if err != nil {
		return nil, err
	}
	return feeds, nil
}

func (sd SqliteDao) Close() error {
	return sd.db.Close()
}
