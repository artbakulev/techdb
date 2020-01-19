package models

import "time"

type Post struct {
	ID       int64     `json:"id,omitempty"`
	Author   string    `json:"author"`
	Created  time.Time `json:"created,omitempty"`
	Forum    string    `json:"forum,omitempty"`
	IsEdited bool      `json:"isEdited,omitempty"`
	Message  string    `json:"message"`
	Parent   int64     `json:"parent,omitempty"`
	Thread   int64     `json:"thread,omitempty"`
	Path     []int64   `json:"-"`
}

type PostFull struct {
	Author *User   `json:"author,omitempty"`
	Forum  *Forum  `json:"forum,omitempty"`
	Post   *Post   `json:"post,omitempty"`
	Thread *Thread `json:"thread,omitempty"`
}

type PostNew struct {
	Author  string `json:"author"`
	Message string `json:"message"`
	Parent  int64  `json:"parent,omitempty"`
}

type PostUpdate struct {
	Message string `json:"message"`
}

type Posts []Post
