package model

type User struct {
	Id       int64  `json:"id,omitempty"`
	Name     string `json:"name"`
	Password string `json:"password"`
}
