package model

import (
	"time"
)

type Comment struct {
	CommentId    Id         `json:"comment_id" bson:"_id" cqadm:"id"`
	BodyHtml     string     `json:"body" bson:"body"`
	BodyMarkdown string     `json:"body_markdown" bson:"body_markdown"`
	PostId       Id         `json:"post_id" bson:"post_id"`
	PostType     string     `json:"post_type" bson:"post_type"`
	Link         string     `json:"link" bson:"link"`
	Score        int32      `json:"score" bson:"score"`
	CreationDate *time.Time `json:"creation_date" bson:"creation_date"`
}

func NewComment() *Comment {
	return new(Comment)
}
