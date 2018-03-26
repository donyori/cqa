package migrate

import (
	"log"
	"sync"

	"bupt.cn/cqa/data/access"
	odtype "bupt.cn/cqa/data/data_type"
	"github.com/donyori/cqa/data/db"
	"github.com/donyori/cqa/data/dtype"
)

func MigrateQuestions(goroutineNumber int) {
	log.Println("Start migration, goroutine number:", goroutineNumber)
	out, res, _, err := access.AllScanQuestions(goroutineNumber, nil)
	if err != nil {
		log.Fatalln(err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(goroutineNumber)
	for i := 0; i < goroutineNumber; i++ {
		go func(no int) {
			defer wg.Done()
			qa, e := db.NewQuestionAccessor()
			if e != nil {
				log.Fatalln(no, e)
				return
			}
			e = qa.Connect()
			if e != nil {
				log.Fatalln(no, e)
				return
			}
			defer qa.Close()
			for oldQ := range out {
				newQ := CopyMgoQuestion(oldQ)
				_, e = qa.Save(newQ)
				if e != nil {
					log.Fatalln(no, e)
				}
			}
		}(i)
	}
	wg.Wait()
	err = <-res
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("Finish")
	}
}

func CopyComment(src *odtype.StackExchangeAPIComment) *dtype.Comment {
	newC := &dtype.Comment{
		CommentId:    src.CommentID,
		BodyHTML:     src.BodyHTML,
		BodyMarkdown: src.BodyMarkdown,
		PostId:       src.PostID,
		PostType:     src.PostType,
		Link:         src.Link,
		Score:        src.Score,
	}
	return newC
}

func CopyAnswer(src *odtype.StackExchangeAPIAnswer) *dtype.Answer {
	newA := &dtype.Answer{
		AnswerId:     src.AnswerID,
		BodyHTML:     src.BodyHTML,
		BodyMarkdown: src.BodyMarkdown,
		IsAccepted:   src.IsAccepted,
		QuestionId:   src.QuestionID,
		Tags:         src.Tags,
		Link:         src.Link,
		Score:        src.Score,
	}
	for _, c := range src.Comments {
		newC := CopyComment(c)
		newA.Comments = append(newA.Comments, newC)
	}
	return newA
}

func CopyMgoQuestion(src *odtype.MgoQuestion) *dtype.Question {
	newQ := &dtype.Question{
		QuestionId:       src.ID,
		Title:            src.Title,
		BodyHTML:         src.BodyHTML,
		BodyMarkdown:     src.BodyMarkdown,
		IsAnswered:       src.IsAnswered,
		AcceptedAnswerId: src.AcceptedAnswerID,
		Tags:             src.Tags,
		Link:             src.Link,
		Score:            src.Score,
		ViewCount:        src.ViewCount,
	}
	for _, a := range src.Answers {
		newA := CopyAnswer(a)
		newQ.Answers = append(newQ.Answers, newA)
	}
	for _, c := range src.Comments {
		newC := CopyComment(c)
		newQ.Comments = append(newQ.Comments, newC)
	}
	return newQ
}
