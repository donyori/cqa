package stackexchange

import (
	"time"

	"github.com/donyori/cqa/data/model"
)

type ResponseWrapper struct {
	Backoff        int    `json:"backoff"`
	ErrorId        int    `json:"error_id"`
	ErrorMessage   string `json:"error_message"`
	ErrorName      string `json:"error_name"`
	HasMore        bool   `json:"has_more"`
	Page           int    `json:"page"`
	PageSize       int    `json:"page_size"`
	QuotaMax       int    `json:"quota_max"`
	QuotaRemaining int    `json:"quota_remaining"`
}

type Comment struct {
	CommentId    int64  `json:"comment_id"`
	BodyHtml     string `json:"body"`
	BodyMarkdown string `json:"body_markdown"`
	PostId       int64  `json:"post_id"`
	PostType     string `json:"post_type"`
	Link         string `json:"link"`
	Score        int32  `json:"score"`
	CreationDate int64  `json:"creation_date"`
}

type Answer struct {
	AnswerId         int64      `json:"answer_id"`
	BodyHtml         string     `json:"body"`
	BodyMarkdown     string     `json:"body_markdown"`
	IsAccepted       bool       `json:"is_accepted"`
	QuestionId       int64      `json:"question_id"`
	Comments         []*Comment `json:"comments"`
	Tags             []string   `json:"tags"`
	Link             string     `json:"link"`
	Score            int32      `json:"score"`
	CreationDate     int64      `json:"creation_date"`
	LastActivityDate int64      `json:"last_activity_date"`
	LastEditDate     int64      `json:"last_edit_date"`
}

type Question struct {
	QuestionId       int64      `json:"question_id"`
	Title            string     `json:"title"`
	BodyHtml         string     `json:"body"`
	BodyMarkdown     string     `json:"body_markdown"`
	IsAnswered       bool       `json:"is_answered"`
	AcceptedAnswerId int64      `json:"accepted_answer_id"`
	Answers          []*Answer  `json:"answers"`
	Comments         []*Comment `json:"comments"`
	Tags             []string   `json:"tags"`
	Link             string     `json:"link"`
	Score            int32      `json:"score"`
	ViewCount        int64      `json:"view_count"`
	CreationDate     int64      `json:"creation_date"`
	LastActivityDate int64      `json:"last_activity_date"`
	LastEditDate     int64      `json:"last_edit_date"`
}

type QuestionsResponse struct {
	ResponseWrapper `json:",inline"`
	Items           []*Question `json:"items"`
}

func (c *Comment) ToDataModel() *model.Comment {
	if c == nil {
		return nil
	}
	res := model.NewComment()
	res.CommentId = model.Id(c.CommentId)
	res.BodyHtml = c.BodyHtml
	res.BodyMarkdown = c.BodyMarkdown
	res.PostId = model.Id(c.PostId)
	res.PostType = c.PostType
	res.Link = c.Link
	res.Score = c.Score
	cd := time.Unix(c.CreationDate, 0)
	res.CreationDate = &cd
	return res
}

func (a *Answer) ToDataModel() *model.Answer {
	if a == nil {
		return nil
	}
	res := model.NewAnswer()
	res.AnswerId = model.Id(a.AnswerId)
	res.BodyHtml = a.BodyHtml
	res.BodyMarkdown = a.BodyMarkdown
	res.IsAccepted = a.IsAccepted
	res.QuestionId = model.Id(a.QuestionId)
	res.Tags = a.Tags
	res.Link = a.Link
	res.Score = a.Score
	cd := time.Unix(a.CreationDate, 0)
	res.CreationDate = &cd
	lad := time.Unix(a.LastActivityDate, 0)
	res.LastActivityDate = &lad
	if a.LastEditDate != 0 {
		led := time.Unix(a.LastEditDate, 0)
		res.LastEditDate = &led
	} else {
		res.LastEditDate = nil
	}
	if a.Comments != nil {
		res.Comments = make([]*model.Comment, 0, len(a.Comments))
		for _, c := range a.Comments {
			res.Comments = append(res.Comments, c.ToDataModel())
		}
	} else {
		res.Comments = nil
	}
	return res
}

func (q *Question) ToDataModel() *model.Question {
	if q == nil {
		return nil
	}
	res := model.NewQuestion()
	res.QuestionId = model.Id(q.QuestionId)
	res.Title = q.Title
	res.BodyHtml = q.BodyHtml
	res.BodyMarkdown = q.BodyMarkdown
	res.IsAnswered = q.IsAnswered
	res.AcceptedAnswerId = model.Id(q.AcceptedAnswerId)
	res.Tags = q.Tags
	res.Link = q.Link
	res.Score = q.Score
	res.ViewCount = q.ViewCount
	cd := time.Unix(q.CreationDate, 0)
	res.CreationDate = &cd
	lad := time.Unix(q.LastActivityDate, 0)
	res.LastActivityDate = &lad
	if q.LastEditDate != 0 {
		led := time.Unix(q.LastEditDate, 0)
		res.LastEditDate = &led
	} else {
		res.LastEditDate = nil
	}
	if q.Comments != nil {
		res.Comments = make([]*model.Comment, 0, len(q.Comments))
		for _, c := range q.Comments {
			res.Comments = append(res.Comments, c.ToDataModel())
		}
	} else {
		res.Comments = nil
	}
	if q.Answers != nil {
		res.Answers = make([]*model.Answer, 0, len(q.Answers))
		for _, a := range q.Answers {
			res.Answers = append(res.Answers, a.ToDataModel())
		}
	} else {
		res.Answers = nil
	}
	return res
}

func (qr *QuestionsResponse) ExtractItems() []*model.Question {
	if qr == nil || qr.Items == nil {
		return nil
	}
	res := make([]*model.Question, 0, len(qr.Items))
	for _, q := range qr.Items {
		res = append(res, q.ToDataModel())
	}
	return res
}
