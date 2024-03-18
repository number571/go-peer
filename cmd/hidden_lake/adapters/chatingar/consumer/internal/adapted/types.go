package adapted

type sCountDTO struct {
	Post struct {
		CommentCount int `json:"commentCount"`
	} `json:"post"`
}

type sMessagesDTO struct {
	Comments []struct {
		Body string `json:"body"`
	} `json:"comments"`
}
