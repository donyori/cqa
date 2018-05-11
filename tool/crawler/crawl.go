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

func CrawlQuestions() (err error) {
	defer log.Println("Crawling questions finish.")
	tags := GlobalSettings.CrawlTags
	logStep := GlobalSettings.LogStep
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
	log.Println("*** Connect to database successfully.")
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
	pageSize := se.GlobalSettings.MaxPageSize
	backoffUnit := se.GlobalSettings.BackoffUnit
	var count int = 0
	defer func() {
		// It's different from "defer log.Printf(...)"
		log.Printf("*** Finally, crawled %d questions.\n", count)
	}()
	for _, tag := range tags {
		log.Printf(
			"*** Start to crawl questions tagged %q."+
				" Already crawled %d questions.\n",
			tag, count)
		page := se.GlobalSettings.StartPage
		if meta.Value.LastActivityDates == nil {
			meta.Value.LastActivityDates = make(map[string]time.Time)
		}
		var min *time.Time = nil
		lastDate, ok := meta.Value.LastActivityDates[tag]
		if ok {
			minTime := lastDate
			min = &minTime
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
					"*** Crawled %d questions, last activity date: %v, quota: remaining = %d, max = %d\n",
					count, lastDate, res.QuotaRemaining, res.QuotaMax)
			}
			if res.QuotaRemaining == 0 {
				err = ErrNoQuotaRemaining
				return err
			}
			if res.Backoff > 0 {
				sleepDuration := backoffUnit * time.Duration(res.Backoff)
				log.Printf("*** Backoff = %d, sleep %v.\n",
					res.Backoff, sleepDuration)
				time.Sleep(sleepDuration)
				log.Println("*** Go on.")
			}
			page++
		}
	}
	return nil
}
