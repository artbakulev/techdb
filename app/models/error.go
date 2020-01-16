package models

type Error struct {
	StatusCode int8   `json:"status_code"`
	Message    string `json:"message"`
}
