package models

type Forum struct {
	Posts   int64  `json:"posts,omitempty"`
	Slug    string `json:"slug"`
	Threads int32  `json:"threads,omitempty"`
	Title   string `json:"title"`
	User    string `json:"user"`
}

type ForumNew struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
	User  string `json:"user"`
}
