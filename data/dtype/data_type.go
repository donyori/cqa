package dtype

type Vector32 struct {
	Data   []float32 `json:"data" bson:"data"`
	L2Norm float32   `json:"l2_norm" bson:"l2_norm"`
}

type Comment struct {
	CommentId    int    `json:"comment_id" bson:"_id"`
	BodyHTML     string `json:"body" bson:"body"`
	BodyMarkdown string `json:"body_markdown" bson:"body_markdown"`
	PostId       int    `json:"post_id" bson:"post_id"`
	PostType     string `json:"post_type" bson:"post_type"`
	Link         string `json:"link" bson:"link"`
	Score        int    `json:"score" bson:"score"`
}

type Answer struct {
	AnswerId     int        `json:"answer_id" bson:"_id"`
	BodyHTML     string     `json:"body" bson:"body"`
	BodyMarkdown string     `json:"body_markdown" bson:"body_markdown"`
	IsAccepted   bool       `json:"is_accepted" bson:"is_accepted"`
	QuestionId   int        `json:"question_id" bson:"question_id"`
	Comments     []*Comment `json:"comments" bson:"comments"`
	Tags         []string   `json:"tags" bson:"tags"`
	Link         string     `json:"link" bson:"link"`
	Score        int        `json:"score" bson:"score"`
}

type Question struct {
	QuestionId       int        `json:"question_id" bson:"_id"`
	Title            string     `json:"title" bson:"title"`
	BodyHTML         string     `json:"body" bson:"body"`
	BodyMarkdown     string     `json:"body_markdown" bson:"body_markdown"`
	IsAnswered       bool       `json:"is_answered" bson:"is_answered"`
	AcceptedAnswerId int        `json:"accepted_answer_id" bson:"accepted_answer_id"`
	Answers          []*Answer  `json:"answers" bson:"answers"`
	Comments         []*Comment `json:"comments" bson:"comments"`
	Tags             []string   `json:"tags" bson:"tags"`
	Link             string     `json:"link" bson:"link"`
	Score            int        `json:"score" bson:"score"`
	ViewCount        int        `json:"view_count" bson:"view_count"`
}

type QuestionVector struct {
	QuestionId  int       `json:"question_id" bson:"_id"`
	TitleVector *Vector32 `json:"title_vector" bson:"title_vector"`
}
