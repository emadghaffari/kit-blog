package model

// User model
type User struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}
