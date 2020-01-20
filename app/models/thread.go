package models

import "time"

type Thread struct {
	ID      int32     `json:"id,omitempty"`
	Slug    string    `json:"slug,omitempty"`
	Author  string    `json:"author"`
	Created time.Time `json:"created,omitempty"`
	Forum   string    `json:"forum,omitempty"`
	Message string    `json:"message"`
	Title   string    `json:"title"`
	Votes   int32     `json:"votes,omitempty"`
}

type ThreadUpdate struct {
	ID      int32  `json:"id,omitempty"`
	Slug    string `json:"slug,omitempty"`
	Title   string `json:"title,omitempty"`
	Message string `json:"message,omitempty"`
}

//easyjson:json
type Threads []Thread
