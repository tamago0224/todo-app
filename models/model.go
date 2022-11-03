package models

type Todo struct {
	Id          int64  `json:"id,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description"`
}
