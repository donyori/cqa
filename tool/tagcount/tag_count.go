package tagcount

import (
	"encoding/csv"
	"log"
	"os"
	"sort"
	"strconv"
	"sync"

	"gopkg.in/mgo.v2/bson"

	"github.com/donyori/cqa/data/db"
	"github.com/donyori/cqa/data/db/mongodb"
	"github.com/donyori/cqa/data/db/wrapper"
)

type tagCount struct {
	Tag   string
	Count int
}

func TagCount() (err error) {
	goroutineNumber := GlobalSettings.GoroutineNumber
	outputFilename := GlobalSettings.OutputFilename
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
	log.Println("Start tag count. goroutine number:", goroutineNumber)
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
	qa, err := wrapper.NewQuestionAccessor(accessor)
	if err != nil {
		return err
	}
	params := mongodb.NewQueryParams()
	params.Selector = bson.M{"tags": 1}
	quitC := make(chan struct{}, 1)
	defer close(quitC)
	outC, resC, err := qa.Scan(params, uint32(goroutineNumber), quitC)
	if err != nil {
		return err
	}
	countC := make(chan *tagCount, goroutineNumber)
	log.Println("*** Start to count.")
	var wg sync.WaitGroup
	wg.Add(goroutineNumber)
	go func() {
		defer close(countC)
		defer func() {
			quitC <- struct{}{}
		}()
		wg.Wait()
	}()
	for i := 0; i < goroutineNumber; i++ {
		go func(number int) {
			defer wg.Done()
			var e error
			defer func() {
				if e != nil {
					log.Printf("*** Error occurs on %v: %v\n", number, e)
				}
			}()
			counter := make(map[string]int)
			count := 0
			for q := range outC {
				if count%logStep == 0 {
					log.Printf("*** Goroutine %v has counted %v questions.",
						number, count)
				}
				for _, tag := range q.Tags {
					counter[tag]++
				}
				count++
			}
			for tag, c := range counter {
				if c <= 0 {
					continue
				}
				countCell := &tagCount{
					Tag:   tag,
					Count: c,
				}
				countC <- countCell
			}
		}(i)
	}
	counter := make(map[string]int)
	for c := range countC {
		counter[c.Tag] += c.Count
	}
	tags := make([]string, 0, len(counter))
	for tag := range counter {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	outputFile, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer outputFile.Close() // Ignore error.
	csvWriter := csv.NewWriter(outputFile)
	defer csvWriter.Flush()
	for _, tag := range tags {
		err = csvWriter.Write([]string{tag, strconv.Itoa(counter[tag])})
		if err != nil {
			return err
		}
	}
	err = <-resC
	if err != nil {
		return err
	}
	log.Println("Done.")
	return nil
}
