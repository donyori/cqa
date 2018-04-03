package model

type Comment struct {
	CommentId    int64  `json:"comment_id" bson:"_id" cqadm:"id"`
	BodyHTML     string `json:"body" bson:"body"`
	BodyMarkdown string `json:"body_markdown" bson:"body_markdown"`
	PostId       int64  `json:"post_id" bson:"post_id"`
	PostType     string `json:"post_type" bson:"post_type"`
	Link         string `json:"link" bson:"link"`
	Score        int32  `json:"score" bson:"score"`
}

func NewComment() *Comment {
	return new(Comment)
}
