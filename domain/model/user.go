package model

import "time"

type User struct {
	Id       int64     `json:"id,omitempty"`
	Name     string    `json:"name"`
	Password string    `json:"password"`
	Created  time.Time `json:"created"`
}
