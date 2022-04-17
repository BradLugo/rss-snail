package models

import (
	"database/sql"
)

var DB *sql.DB

type User struct {
	UserId int
	Email  string
}

// UserFeeds will consist of UserId's and FeedId's, no need for a duplicate struct

type Feed struct {
	FeedId int
	Url    string
}
