package embedding

import (
	"log"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/donyori/cqa/data/db"
	"github.com/donyori/cqa/data/db/mongodb"
	"github.com/donyori/cqa/data/dtype"
	"github.com/donyori/cqa/nlp"
)

func QuestionEmbedding() error {
	nGoroutines := GlobalSettings.GoroutineNumber
	minMs := GlobalSettings.MinMillisecond
	log.Println("Start question embedding. goroutine number:",
		nGoroutines)
	var err error
	defer func() {
		if err != nil {
			log.Fatalln(err)
		}
	}()
	qa, err := db.NewQuestionAccessor()
	if err != nil {
		return err
	}
	err = qa.Connect()
	if err != nil {
		return err
	}
	defer qa.Close()
	log.Println("*** Succeed to connect database.")
	params := mongodb.NewQueryParams()
	params.Selector = bson.M{"_id": 1, "title": 1}
	out, res, quit, err := qa.Scan(params, nGoroutines)
	if err != nil {
		return err
	}
	log.Println("*** Start embedding.")
	var wg sync.WaitGroup
	wg.Add(nGoroutines)
	for i := 0; i < nGoroutines; i++ {
		go func(number int) {
			defer wg.Done()
			var e error
			defer func() {
				if e != nil {
					log.Fatalf("*** Error occurs on %v: %v\n", number, e)
				}
			}()
			qva, e := db.NewQuestionVectorAccessor()
			if e != nil {
				return
			}
			e = qva.Connect()
			if e != nil {
				return
			}
			defer qva.Close()
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
			for question := range out {
				if count > 0 && count%100 == 0 {
					log.Printf("*** Goroutine %v has embedded %v questions.", number, count)
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
				qv := &dtype.QuestionVector{
					QuestionId:  question.QuestionId,
					TitleVector: vector,
				}
				_, e = qva.Save(qv)
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
	quit <- struct{}{}
	close(quit)
	err = <-res
	if err != nil {
		return err
	}
	log.Println("Done.")
	return nil
}
