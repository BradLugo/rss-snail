package internal

type User struct {
	id    int
	email string
	feeds []Feed
}

type Feed struct {
	id   int
	name string
	url  string
}
