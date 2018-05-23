package nndatasetmaker

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"sync"

	"gopkg.in/mgo.v2/bson"

	"github.com/donyori/cqa/data/db"
	dbid "github.com/donyori/cqa/data/db/id"
	"github.com/donyori/cqa/data/db/mongodb"
	"github.com/donyori/cqa/data/db/wrapper"
	"github.com/donyori/cqa/data/model"
	"github.com/donyori/cqa/qa/qc"
)

type addFlag int8

type candidate int

const (
	addFlagNowhere addFlag = iota
	addFlagToTrain
	addFlagToEval
)

func SelectQuestions() (err error) {
	epls := GlobalSettings.MaxExampleNumbersPerLabel
	defer func() {
		if err != nil {
			log.Println(err)
		}
	}()
	eplsLen := len(epls)
	for _, v := range epls {
		if v == nil {
			eplsLen--
		}
	}
	if eplsLen == 0 {
		err = errors.New("GlobalSettings.MaxExampleNumbersPerLabel is empty")
	}
	log.Println("Start selecting questions.")
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
	qca, err := wrapper.NewQuestionClassificationAccessor(accessor)
	if err != nil {
		return err
	}
	params := mongodb.NewQueryParams()
	params.Query = bson.M{"score": bson.M{"$gte": int32(0)}}
	params.Selector = bson.M{"_id": 1}
	params.SortFields = []string{"-view_count"}
	quitC := make(chan struct{})
	var qOutC <-chan *model.Question
	defer func() {
		if quitC != nil {
			close(quitC)
			quitC = nil
		}
		// Drain qOutC, to avoid blocking the question scanner.
		for _ = range qOutC {
		}
	}()
	qOutC, resC, err := qa.Scan(params, 4, quitC)
	if err != nil {
		return err
	}
	qcChans := make(map[string]chan *model.QuestionClassification)
	doneChan := make(chan string, eplsLen)
	log.Println("*** Start selecting.")
	var datasetMap sync.Map
	var wg sync.WaitGroup
	for k, v := range epls {
		if v == nil {
			continue
		}
		qcChan := make(chan *model.QuestionClassification, 8)
		qcChans[k] = qcChan
		wg.Add(1)
		go func(name string, epl *ExampleNumbersPerLabelPair,
			qcChan <-chan *model.QuestionClassification,
			quitChan <-chan struct{}) {
			defer wg.Done()
			var e error
			defer func() {
				if e != nil {
					log.Printf(
						"*** %s - Error occurs when select questions: %v.\n",
						name, e)
				}
			}()
			trainIds, evalIds, e := selectQuestionsByEpl(
				name, epl, qcChan, quitChan, doneChan)
			if e != nil {
				return
			}
			select {
			case <-quitChan:
				log.Printf("*** %s - Receive a quit signal and exit.\n", name)
				return
			default:
				if trainIds == nil {
					e = errors.New("trainIds is nil")
					return
				}
				if evalIds == nil {
					e = errors.New("evalIds is nil")
					return
				}
				idDatasets := &IdDatasets{Train: trainIds, Eval: evalIds}
				datasetMap.Store(name, idDatasets)
			}
			log.Printf("*** %s - Done and exit.\n", name)
		}(k, v, qcChan, quitC)
	}
	go func() {
		defer close(doneChan)
		wg.Wait()
	}()
	queryParam := bson.M{"_id": 0}
	params = mongodb.NewQueryParams()
	params.Query = queryParam
	params.Selector = bson.M{"_id": 1, "classification_by_tag": 1}
	params.Limit = 1
	isAllDone := false
	for !isAllDone {
		select {
		case name, ok := <-doneChan:
			if !ok {
				isAllDone = true
				break
			}
			qcChan, ok := qcChans[name]
			if !ok {
				log.Printf("*** Receive BAD done info: %q.\n", name)
				break
			}
			log.Printf("*** Receive done info: %q.\n", name)
			if qcChan != nil {
				close(qcChan)
			}
			delete(qcChans, name)
		case question, ok := <-qOutC:
			if !ok {
				qOutC = nil
				names := make([]string, 0, len(qcChans))
				for name, qcChan := range qcChans {
					names = append(names, name)
					if qcChan != nil {
						close(qcChan)
					}
				}
				for _, name := range names {
					qcChans[name] = nil
				}
				break
			}
			if question == nil || err != nil {
				break
			}
			queryParam["_id"] = question.QuestionId
			var qcm *model.QuestionClassification
			qcm, err = qca.FetchOne(params)
			if err != nil {
				log.Println(
					"*** Cannot fetch data from QuestionClassification collection.")
				if quitC != nil {
					close(quitC)
					quitC = nil
				}
				break
			}
			for _, qcChan := range qcChans {
				if qcChan == nil {
					continue
				}
				qcChan <- qcm
			}
		}
	}
	log.Println("*** All selectors are done.")
	if quitC != nil {
		close(quitC)
		quitC = nil
	}
	// Drain qOutC, to avoid blocking the question scanner.
	for _ = range qOutC {
	}
	if err != nil {
		return err
	}
	err = <-resC
	if err != nil {
		return err
	}

	log.Println("*** Create meta info and save it.")
	meta := NewMeta()
	if meta.Value == nil {
		meta.Value = NewMetaValue()
	}
	if meta.Value.IdDatasetsMap == nil {
		meta.Value.IdDatasetsMap = make(map[string]IdDatasets)
	}
	defer func() {
		if panicErr := recover(); panicErr != nil {
			e, ok := panicErr.(error)
			if ok {
				err = e
			} else {
				err = fmt.Errorf("%v", panicErr)
			}
		}
	}()
	datasetMap.Range(func(key, value interface{}) bool {
		keyStr := key.(string)
		idDatasets := value.(*IdDatasets)
		meta.Value.IdDatasetsMap[keyStr] = *idDatasets
		return true
	})
	_, err = accessor.SaveOne(dbid.MetaCollection, nil, meta)
	if err != nil {
		return err
	}

	log.Println("Done.")
	return nil
}

func selectQuestionsByEpl(name string, epl *ExampleNumbersPerLabelPair,
	qcChan <-chan *model.QuestionClassification, quitChan <-chan struct{},
	doneChan chan<- string) (
	trainIds []model.Id, evalIds []model.Id, err error) {
	doesContainNoLabelQuestions := GlobalSettings.DoesContainNoLabelQuestions
	logStep := GlobalSettings.LogStep
	defer func() {
		if panicErr := recover(); panicErr != nil {
			e, ok := panicErr.(error)
			if ok {
				err = e
			} else {
				err = fmt.Errorf("%v", panicErr)
			}
			trainIds = nil
			evalIds = nil
		}
	}()
	defer func() {
		if doneChan != nil {
			doneChan <- name
		}
		// Drain qcChan, to avoid blocking the broadcaster.
		for _ = range qcChan {
		}
	}()
	trainEpl := epl.TrainDataset
	evalEpl := epl.EvalDataset
	ter := float64(trainEpl) / float64(evalEpl)

	labelNum := int32(len(qc.KnownLabels))
	if doesContainNoLabelQuestions {
		labelNum++ // +1 for no label questions.
	}
	trainCounts := make(map[qc.Label]int32)
	evalCounts := make(map[qc.Label]int32)
	var enoughTrainCounts, enoughEvalCounts int32
	teCounts := make(map[qc.Label]float64) // +1 when add an example to train dataset and -ter when add an example to eval dataset.
	trainSet := make(map[model.Id]bool)
	evalSet := make(map[model.Id]bool)
	count := 0
	labelsBase := make([]qc.Label, len(qc.KnownLabels)+1)
	isQuit := false
	for !isQuit {
		select {
		case <-quitChan:
			return nil, nil, nil
		case qcm, ok := <-qcChan:
			if !ok {
				isQuit = true
				break
			}
			if count%logStep == 0 {
				log.Printf("*** %s - Has selected over %d questions.\n",
					name, count)
			}
			count++
			if qcm == nil {
				break
			}
			labels := labelsBase[:0]
			for _, labelStr := range qcm.ClassificationByTag {
				label := qc.ParseLabel(labelStr, true)
				if !label.IsKnown() {
					continue
				}
				labels = append(labels, label)
			}
			if len(labels) == 0 {
				if !doesContainNoLabelQuestions {
					break
				}
				labels = append(labels, qc.UnknownLabel)
			}
			var af addFlag
			var train, eval candidate
			for _, label := range labels {
				tc := trainCounts[label]
				ec := evalCounts[label]
				if tc >= trainEpl {
					train.Ban()
					if ec >= evalEpl {
						eval.Ban()
						break
					}
					eval.Vote()
				} else if ec >= evalEpl {
					eval.Ban()
					train.Vote()
				} else if teCounts[label] < ter {
					train.Vote()
				} else {
					eval.Vote()
				}
			}
			if train >= eval {
				if train > 0 {
					af = addFlagToTrain
				} else {
					af = addFlagNowhere
				}
			} else if eval > 0 {
				af = addFlagToEval
			} else {
				af = addFlagNowhere
			}
			switch af {
			case addFlagToTrain:
				trainSet[qcm.QuestionId] = true
				for _, label := range labels {
					teCounts[label]++
					trainCounts[label]++
					if trainCounts[label] > trainEpl {
						return nil, nil, fmt.Errorf(
							"overflow the max number of train examples per label(%v), question id: %v, label: %v",
							trainEpl, qcm.QuestionId, label)
					}
					if trainCounts[label] == trainEpl {
						enoughTrainCounts++
					}
				}
			case addFlagToEval:
				evalSet[qcm.QuestionId] = true
				for _, label := range labels {
					teCounts[label] -= ter
					evalCounts[label]++
					if evalCounts[label] > evalEpl {
						return nil, nil, fmt.Errorf(
							"overflow the max number of eval examples per label(%v), question id: %v, label: %v",
							evalEpl, qcm.QuestionId, label)
					}
					if evalCounts[label] == evalEpl {
						enoughEvalCounts++
					}
				}
			}
			if enoughTrainCounts >= labelNum && enoughEvalCounts >= labelNum {
				isQuit = true
			}
		}
	}

	trainList := make([]model.Id, 0, len(trainSet))
	for id, ok := range trainSet {
		if ok {
			trainList = append(trainList, id)
		}
	}
	sort.Slice(trainList, func(i, j int) bool {
		return trainList[i] < trainList[j]
	})
	evalList := make([]model.Id, 0, len(evalSet))
	for id, ok := range evalSet {
		if ok {
			evalList = append(evalList, id)
		}
	}
	sort.Slice(evalList, func(i, j int) bool {
		return evalList[i] < evalList[j]
	})

	return trainList, evalList, nil
}

func (c *candidate) Vote() {
	if *c >= 0 {
		*c++
	}
}

func (c *candidate) Ban() {
	*c = -1
}
