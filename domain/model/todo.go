package model

type Todo struct {
	Id          int64  `json:"id,omitempty"`
	UserId      int64  `json:"user_id,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
