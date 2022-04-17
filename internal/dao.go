package internal

type Dao interface {
	GetFeeds(User) ([]Feed, error)
	Close() error
}
