package embedding

import (
	"log"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/donyori/cqa/data/db"
	"github.com/donyori/cqa/data/db/mongodb"
	"github.com/donyori/cqa/data/db/wrapper"
	"github.com/donyori/cqa/data/model"
	"github.com/donyori/cqa/nlp"
)

func QuestionEmbedding() error {
	goroutineNumber := GlobalSettings.GoroutineNumber
	minMs := GlobalSettings.MinMillisecond
	var err error
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()
	if goroutineNumber <= 0 {
		err = ErrNonPositiveGoroutineNumber
		return err
	}
	log.Println("Start question embedding. goroutine number:",
		goroutineNumber)
	session, err := db.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	log.Println("*** Succeed to connect to database.")
	accessor, err := db.NewAccessor(session)
	if err != nil {
		return err
	}
	qa, err := wrapper.NewQuestionAccessor(accessor)
	if err != nil {
		return err
	}
	params := mongodb.NewQueryParams()
	params.Selector = bson.M{"_id": 1, "title": 1}
	quitC := make(chan struct{}, 1)
	defer close(quitC)
	outC, resC, err := qa.Scan(params, uint32(goroutineNumber), quitC)
	if err != nil {
		return err
	}
	log.Println("*** Start embedding.")
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
			qva, e := wrapper.NewQuestionVectorAccessor(acc)
			if e != nil {
				return
			}
			var timer *time.Timer
			var duration time.Duration
			needDrainC := false
			if minMs > 0 {
				duration = time.Millisecond * time.Duration(minMs)
				timer = time.NewTimer(duration)
				needDrainC = true
				defer func() {
					if timer != nil && !timer.Stop() && needDrainC {
						<-timer.C
					}
				}()
			}
			count := 0
			for question := range outC {
				if count%20 == 0 {
					log.Printf("*** Goroutine %v has embedded %v questions.",
						number, count)
				}
				if timer != nil {
					if !timer.Stop() && needDrainC {
						<-timer.C
					}
					timer.Reset(duration)
					needDrainC = true
				}
				vector, e := nlp.Embedding(question.Title)
				if e != nil {
					return
				}
				qv := model.NewQuestionVector()
				qv.QuestionId = question.QuestionId
				qv.TitleVector = vector
				_, e = qva.SaveById(qv.QuestionId, qv)
				if e != nil {
					return
				}
				if timer != nil {
					<-timer.C
					needDrainC = false
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
