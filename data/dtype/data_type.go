package dtype

type Comment struct {
	CommentID    int    `json:"comment_id" bson:"_id"`
	BodyHTML     string `json:"body" bson:"body"`
	BodyMarkdown string `json:"body_markdown" bson:"body_markdown"`
	PostID       int    `json:"post_id" bson:"post_id"`
	PostType     string `json:"post_type" bson:"post_type"`
	Link         string `json:"link" bson:"link"`
	Score        int    `json:"score" bson:"score"`
}

type Answer struct {
	AnswerID     int        `json:"answer_id" bson:"_id"`
	BodyHTML     string     `json:"body" bson:"body"`
	BodyMarkdown string     `json:"body_markdown" bson:"body_markdown"`
	IsAccepted   bool       `json:"is_accepted" bson:"is_accepted"`
	QuestionID   int        `json:"question_id" bson:"question_id"`
	Comments     []*Comment `json:"comments" bson:"comments"`
	Tags         []string   `json:"tags" bson:"tags"`
	Link         string     `json:"link" bson:"link"`
	Score        int        `json:"score" bson:"score"`
}

type Question struct {
	QuestionID       int        `json:"question_id" bson:"_id"`
	Title            string     `json:"title" bson:"title"`
	BodyHTML         string     `json:"body" bson:"body"`
	BodyMarkdown     string     `json:"body_markdown" bson:"body_markdown"`
	IsAnswered       bool       `json:"is_answered" bson:"is_answered"`
	AcceptedAnswerID int        `json:"accepted_answer_id" bson:"accepted_answer_id"`
	Answers          []*Answer  `json:"answers" bson:"answers"`
	Comments         []*Comment `json:"comments" bson:"comments"`
	Tags             []string   `json:"tags" bson:"tags"`
	Link             string     `json:"link" bson:"link"`
	Score            int        `json:"score" bson:"score"`
	ViewCount        int        `json:"view_count" bson:"view_count"`
}
