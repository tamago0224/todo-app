package model

import "time"

type Todo struct {
	Id          int64      `json:"id,omitempty"`
	UserId      int64      `json:"user_id,omitempty"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Done        bool       `json:"done"`
	Deadline    *time.Time `json:"deadline"`
	Created     time.Time  `json:"created"`
}
