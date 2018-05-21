package tagcount

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
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
	outputDirectory := GlobalSettings.OutputDirectory
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
	outputFilenameByTag := filepath.Join(outputDirectory, "sort_by_tag.csv")
	outputFilenameByCount := filepath.Join(outputDirectory, "sort_by_count.csv")
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
	log.Println("*** Start counting.")
	var wg sync.WaitGroup
	wg.Add(goroutineNumber)
	go func() {
		defer close(countC)
		defer func() {
			quitC <- struct{}{}
			// Drain outC, to avoid blocking the question scanner.
			for _ = range outC {
			}
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
					log.Printf("*** Goroutine %v has counted %d questions.\n",
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
	tcs := make([]*tagCount, 0, len(counter))
	for t, c := range counter {
		tcs = append(tcs, &tagCount{Tag: t, Count: c})
	}
	lessByTag := func(i, j int) bool {
		return tcs[i].Tag < tcs[j].Tag
	}
	lessByCount := func(i, j int) bool {
		if tcs[i].Count == tcs[j].Count {
			return tcs[i].Tag < tcs[j].Tag
		}
		return tcs[i].Count > tcs[j].Count
	}
	sort.Slice(tcs, lessByTag)
	outputFileByTag, err := os.Create(outputFilenameByTag)
	if err != nil {
		return err
	}
	defer outputFileByTag.Close() // Ignore error.
	csvWriterByTag := csv.NewWriter(outputFileByTag)
	defer csvWriterByTag.Flush()
	for _, tc := range tcs {
		err = csvWriterByTag.Write([]string{tc.Tag, strconv.Itoa(tc.Count)})
		if err != nil {
			return err
		}
	}
	sort.Slice(tcs, lessByCount)
	outputFileByCount, err := os.Create(outputFilenameByCount)
	if err != nil {
		return err
	}
	defer outputFileByCount.Close() // Ignore error.
	csvWriterByCount := csv.NewWriter(outputFileByCount)
	defer csvWriterByCount.Flush()
	for _, tc := range tcs {
		err = csvWriterByCount.Write([]string{tc.Tag, strconv.Itoa(tc.Count)})
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
