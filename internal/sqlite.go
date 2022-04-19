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
		"file://migrations",
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

func (sd SqliteDao) AddUser(user User) error {
	_, err := sd.db.NamedExec("INSERT INTO user (email) VALUES (:email)", user)
	return err
}

func (sd SqliteDao) AddFeed(userId int, feed Feed) error {
	_, err := sd.db.NamedExec("INSERT INTO ")
	return err
}

func (sd SqliteDao) GetFeeds(userId int) ([]Feed, error) {
	var feeds []Feed
	err := sd.db.Select(&feeds, "SELECT title, url FROM feed WHERE id = () ORDER BY title", feeds)
	if err != nil {
		return nil, err
	}
	return feeds, nil
}

func (sd SqliteDao) Close() error {
	return sd.db.Close()
}
