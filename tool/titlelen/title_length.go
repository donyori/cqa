package titlelen

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

type lenCount struct {
	Length int
	Count  int
}

func TitleLengthCount() (err error) {
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
	outputFilenameByLen := filepath.Join(
		outputDirectory, "title_length_count_sort_by_len.csv")
	outputFilenameByCount := filepath.Join(
		outputDirectory, "title_length_count_sort_by_count.csv")
	log.Println("Start title length count. goroutine number:", goroutineNumber)
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
	qva, err := wrapper.NewQuestionVectorAccessor(accessor)
	if err != nil {
		return err
	}
	params := mongodb.NewQueryParams()
	params.Selector = bson.M{"title_token_vectors": 1}
	quitC := make(chan struct{}, 1)
	defer close(quitC)
	outC, resC, err := qva.Scan(params, uint32(goroutineNumber), quitC)
	if err != nil {
		return err
	}
	countC := make(chan *lenCount, goroutineNumber)
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
			counter := make(map[int]int)
			count := 0
			for qv := range outC {
				if count%logStep == 0 {
					log.Printf("*** Goroutine %v has counted %d questions.",
						number, count)
				}
				length := len(qv.TitleTokenVectors)
				counter[length]++
				count++
			}
			for length, c := range counter {
				if c <= 0 {
					continue
				}
				countCell := &lenCount{
					Length: length,
					Count:  c,
				}
				countC <- countCell
			}
		}(i)
	}
	counter := make(map[int]int)
	for c := range countC {
		counter[c.Length] += c.Count
	}
	lcs := make([]*lenCount, 0, len(counter))
	for l, c := range counter {
		lcs = append(lcs, &lenCount{Length: l, Count: c})
	}
	lessByLen := func(i, j int) bool {
		return lcs[i].Length < lcs[j].Length
	}
	lessByCount := func(i, j int) bool {
		if lcs[i].Count == lcs[j].Count {
			return lcs[i].Length < lcs[j].Length
		}
		return lcs[i].Count > lcs[j].Count
	}
	sort.Slice(lcs, lessByLen)
	outputFileByLen, err := os.Create(outputFilenameByLen)
	if err != nil {
		return err
	}
	defer outputFileByLen.Close() // Ignore error.
	csvWriterByLen := csv.NewWriter(outputFileByLen)
	defer csvWriterByLen.Flush()
	for _, lc := range lcs {
		err = csvWriterByLen.Write([]string{
			strconv.Itoa(lc.Length),
			strconv.Itoa(lc.Count),
		})
		if err != nil {
			return err
		}
	}
	sort.Slice(lcs, lessByCount)
	outputFileByCount, err := os.Create(outputFilenameByCount)
	if err != nil {
		return err
	}
	defer outputFileByCount.Close() // Ignore error.
	csvWriterByCount := csv.NewWriter(outputFileByCount)
	defer csvWriterByCount.Flush()
	for _, lc := range lcs {
		err = csvWriterByCount.Write([]string{
			strconv.Itoa(lc.Length),
			strconv.Itoa(lc.Count),
		})
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
