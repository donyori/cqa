package classifier

import (
	"log"
	"sync"

	"gopkg.in/mgo.v2/bson"

	"github.com/donyori/cqa/data/db"
	"github.com/donyori/cqa/data/db/mongodb"
	"github.com/donyori/cqa/data/db/wrapper"
	"github.com/donyori/cqa/data/model"
	"github.com/donyori/cqa/qa/qc"
)

func ClassifyByTag() (err error) {
	goroutineNumber := GlobalSettings.GoroutineNumber
	logStep := GlobalSettings.LogStep
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()
	if goroutineNumber <= 0 {
		err = ErrNonPositiveGoroutineNumber
		return err
	}
	log.Println("Start classifying by tag. goroutine number:", goroutineNumber)
	session, err := db.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	log.Println("*** Connect to database successfully.")
	accessor, err := db.NewAccessor(session)
	if err != nil {
		return err
	}
	qa, err := wrapper.NewQuestionAccessor(accessor)
	if err != nil {
		return err
	}
	params := mongodb.NewQueryParams()
	params.Selector = bson.M{
		"_id": 1, "tags": 1,
	}
	quitC := make(chan struct{}, 1)
	defer close(quitC)
	outC, resC, err := qa.Scan(params, uint32(goroutineNumber), quitC)
	if err != nil {
		return err
	}
	log.Println("*** Start classifying.")
	var wg sync.WaitGroup
	wg.Add(goroutineNumber)
	for i := 0; i < goroutineNumber; i++ {
		go func(number int) {
			defer wg.Done()
			var e error
			defer func() {
				if e != nil {
					log.Printf("*** Error occurs on %v: %v\n", number, e)
				}
			}()
			sess, e := db.NewSession()
			if e != nil {
				return
			}
			defer sess.Close()
			acc, e := db.NewAccessor(sess)
			if e != nil {
				return
			}
			qca, e := wrapper.NewQuestionClassificationAccessor(acc)
			if e != nil {
				return
			}
			count := 0
			for question := range outC {
				if count%logStep == 0 {
					log.Printf("*** Goroutine %v has classified %d questions.",
						number, count)
				}
				labels := qc.GetLabelStringsByTags(question.Tags)
				var qc *model.QuestionClassification
				qc, e = qca.FetchOneById(question.QuestionId)
				if e != nil {
					return
				}
				if qc == nil {
					qc = model.NewQuestionClassification()
					qc.ClassificationByNn = nil
				}
				qc.QuestionId = question.QuestionId
				qc.ClassificationByTag = labels
				_, e = qca.SaveOne(nil, qc)
				if e != nil {
					return
				}
				count++
			}
		}(i)
	}
	wg.Wait()
	quitC <- struct{}{}
	err = <-resC
	if err != nil {
		return err
	}
	log.Println("Done.")
	return nil
}
