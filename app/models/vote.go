package models

type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int8   `json:"voice"`
	Thread   int64  `json:"thread"`
}
