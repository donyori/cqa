package embedding

import (
	"errors"
	"log"
	"reflect"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/donyori/cqa/data/db"
	dbid "github.com/donyori/cqa/data/db/id"
	"github.com/donyori/cqa/data/db/mongodb"
	"github.com/donyori/cqa/data/db/wrapper"
	"github.com/donyori/cqa/data/model"
	"github.com/donyori/cqa/nlp"
)

var ErrCannotGetMeta error = errors.New("cannot get crawler meta")

func QuestionEmbedding() (err error) {
	goroutineNumber := GlobalSettings.GoroutineNumber
	minMs := GlobalSettings.MinMillisecond
	enableQuestionBuffer := GlobalSettings.EnableQuestionBuffer
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
	log.Println("Start question embedding. goroutine number:", goroutineNumber)
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
	metaRes, err := accessor.FetchOneById(
		dbid.MetaCollection, MetaKey, reflect.TypeOf(Meta{}))
	if err != nil {
		return err
	}
	var meta *Meta = nil
	if metaRes != nil {
		var ok bool
		meta, ok = metaRes.(*Meta)
		if !ok {
			err = ErrCannotGetMeta
			return err
		}
	}
	hasMetaBefore := true
	if meta == nil {
		meta = NewMeta()
		hasMetaBefore = false
	}
	if meta.Value == nil {
		meta.Value = NewMetaValue()
		hasMetaBefore = false
	}
	log.Println("*** Load meta successfully.")
	defer func() {
		if err != nil {
			log.Println(err)
			err = nil
		}
		meta.Value.QuestionLastEmbeddingTime = time.Now()
		_, err = accessor.SaveOne(dbid.MetaCollection, nil, meta)
	}()
	qa, err := wrapper.NewQuestionAccessor(accessor)
	if err != nil {
		return err
	}
	params := mongodb.NewQueryParams()
	if hasMetaBefore {
		params.Query = bson.M{"last_create_or_edit_date": bson.M{
			"$gte": meta.Value.QuestionLastCreateOrEditDate,
		}}
	}
	params.Selector = bson.M{
		"_id": 1, "title": 1, "last_create_or_edit_date": 1,
	}
	params.SortFields = []string{"last_create_or_edit_date"}
	var outC <-chan *model.Question
	var resC <-chan error
	quitC := make(chan struct{})
	defer close(quitC)
	if enableQuestionBuffer {
		log.Println("*** Start to buffer questions.")
		var qs []*model.Question
		qs, err = qa.FetchAll(params)
		if err != nil {
			return err
		}
		log.Println("*** Buffer questions successfully.")
		outChannel := make(chan *model.Question, uint32(goroutineNumber))
		resChannel := make(chan error, 1)
		outC = outChannel
		resC = resChannel
		go func() {
			defer close(resChannel)
			defer close(outChannel)
			for _, q := range qs {
				select {
				case <-quitC:
					return
				default:
					outChannel <- q
				}
			}
			resChannel <- nil
		}()
	} else {
		log.Println("*** Disable to buffer questions.")
		outC, resC, err = qa.Scan(params, uint32(goroutineNumber), quitC)
		if err != nil {
			return err
		}
	}
	writeMetaC := make(chan time.Time, goroutineNumber)
	writeMetaDoneC := make(chan struct{})
	log.Println("*** Start embedding.")
	go func() {
		defer close(writeMetaDoneC)
		var lastT time.Time
		isFirst := true
		for t := range writeMetaC {
			if !isFirst && !t.After(lastT) {
				continue
			}
			isFirst = false
			lastT = t
			meta.Value.QuestionLastEmbeddingTime = time.Now()
			meta.Value.QuestionLastCreateOrEditDate = lastT
			accessor.SaveOne(dbid.MetaCollection, nil, meta) // Ignore error.
		}
	}()
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
				if count%logStep == 0 {
					log.Printf("*** Goroutine %v has embedded %v questions.",
						number, count)
				}
				var isExisted bool
				isExisted, e = qva.IsExistedById(question.QuestionId)
				if e != nil {
					return
				}
				if isExisted {
					continue
				}
				if timer != nil {
					if !timer.Stop() && needDrainC {
						<-timer.C
					}
					timer.Reset(duration)
					needDrainC = true
				}
				var vector *model.Vector32
				var tokenVectors []*model.TokenVector
				vector, tokenVectors, e =
					nlp.EmbeddingWithTokens(question.Title)
				if e != nil {
					return
				}
				qv := model.NewQuestionVector()
				qv.QuestionId = question.QuestionId
				qv.TitleVector = vector
				qv.TitleTokenVectors = tokenVectors
				_, e = qva.SaveOne(nil, qv)
				if e != nil {
					return
				}
				if question.LastCreateOrEditDate != nil {
					writeMetaC <- *question.LastCreateOrEditDate
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
	close(writeMetaC)
	err = <-resC
	<-writeMetaDoneC
	if err != nil {
		return err
	}
	log.Println("Done.")
	return nil
}
