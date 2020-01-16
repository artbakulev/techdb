package models

type PostsRequestQuery struct {
	ThreadID   int64  `json:"id,omitempty"`
	ThreadSlug string `json:"slug,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	Since      int64  `json:"since,omitempty"`
	Sort       string `json:"sort,omitempty"`
	Desc       bool   `json:"desc,omitempty"`
}

type PostsRelatedQuery struct {
	NeedAuthor bool `json:"need_author"`
	NeedForum  bool `json:"need_forum"`
	NeedThread bool `json:"need_thread"`
}
