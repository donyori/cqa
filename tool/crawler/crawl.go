package crawler

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"time"

	se "github.com/donyori/cqa/api/restful/stackexchange"
	"github.com/donyori/cqa/data/db"
	dbid "github.com/donyori/cqa/data/db/id"
	"github.com/donyori/cqa/data/db/wrapper"
)

var (
	ErrCannotGetMeta    error = errors.New("cannot get crawler meta")
	ErrNoQuotaRemaining error = errors.New("quota remaining is zero")
)

func CrawlQuestions() error {
	defer log.Println("Crawling questions finish.")
	tags := GlobalSettings.CrawlTags
	logStep := GlobalSettings.LogStep
	var err error
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()
	log.Println("Start to crawl questions.")
	sess, err := db.NewSession()
	if err != nil {
		return err
	}
	defer sess.Close()
	accessor, err := db.NewAccessor(sess)
	if err != nil {
		return err
	}
	questionAccessor, err := wrapper.NewQuestionAccessor(accessor)
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
	if meta == nil {
		meta = NewMeta()
	}
	if meta.Value == nil {
		meta.Value = NewMetaValue()
	}
	log.Println("*** Load meta successfully.")
	defer func() {
		if err != nil {
			log.Println(err)
			err = nil
		}
		meta.Value.LastCrawlTime = time.Now()
		_, err = accessor.SaveOne(dbid.MetaCollection, nil, meta)
	}()
	page := se.GlobalSettings.StartPage
	pageSize := se.GlobalSettings.MaxPageSize
	backoffUnit := se.GlobalSettings.BackoffUnit
	var count int = 0
	defer func() {
		// It's different from "defer log.Printf(...)"
		log.Printf("*** Finally, crawled %d questions\n", count)
	}()
	for _, tag := range tags {
		if meta.Value.LastActivityDates == nil {
			meta.Value.LastActivityDates = make(map[string]time.Time)
		}
		var min *time.Time = nil
		lastDate, ok := meta.Value.LastActivityDates[tag]
		if ok {
			min = &lastDate
		}
		hasMore := true
		for hasMore {
			var res *se.QuestionsResponse
			res, err = se.Questions(page, pageSize, nil, nil,
				se.QuestionsSortActivity, se.OrderAsc, min, nil, tag)
			if err != nil {
				return err
			}
			hasMore = res.HasMore
			if res.ErrorName != "" || res.ErrorMessage != "" {
				err = fmt.Errorf("error %d - %v: %v",
					res.ErrorId, res.ErrorName, res.ErrorMessage)
				return err
			}
			questions := res.ExtractItems()
			if len(questions) > 0 {
				err = nil
				for _, q := range questions {
					_, err = questionAccessor.SaveOne(nil, q)
					if err != nil {
						break
					}
					count++
					if q.LastActivityDate != nil &&
						q.LastActivityDate.After(lastDate) {
						lastDate = *q.LastActivityDate
					}
				}
				meta.Value.LastCrawlTime = time.Now()
				meta.Value.LastActivityDates[tag] = lastDate
				accessor.SaveOne(dbid.MetaCollection, nil, meta) // Ignore error.
				if err != nil {
					return err
				}
			}
			if page%logStep == 0 {
				log.Printf(
					"*** Crawled %d questions, quota: remaining = %d, max = %d\n",
					count, res.QuotaRemaining, res.QuotaMax)
			}
			if res.QuotaRemaining == 0 {
				err = ErrNoQuotaRemaining
				return err
			}
			if res.Backoff > 0 {
				time.Sleep(backoffUnit * time.Duration(res.Backoff))
			}
			page++
		}
	}
	return nil
}
