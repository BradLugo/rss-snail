package internal

type Dao interface {
	AddUser(id int) error
	AddFeed(userId int, feed Feed) error
	GetFeeds(userId int) ([]Feed, error)
	Close() error
}
