package adapted

type sCountDTO struct {
	Post struct {
		CommentCount int `json:"commentCount"`
	} `json:"post"`
}

type sMessagesDTO struct {
	Comments []sCommentBlock `json:"comments"`
}

type sCommentBlock struct {
	Body      string `json:"body"`
	Timestamp string `json:"timestamp"`
}
