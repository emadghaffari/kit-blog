package model

// Post struct
type Post struct {
	ID          string
	Token       *string
	Title       string
	Slug        string
	Description string
	Body        string
	Header      string
	CreatedAT   *string
}
