package qa

import (
	"time"

	"github.com/donyori/cqa/data/db"
	"github.com/donyori/cqa/data/db/wrapper"
	"github.com/donyori/cqa/data/model"
	"github.com/donyori/cqa/qa/qm"
)

func SearchSimilarQuestions(question string, topNumber int,
	timeLimit time.Duration) (
	res []*SimilarQuestion, isTimeout bool, err error) {
	respC, err := qm.Match(question, topNumber, timeLimit)
	if err != nil {
		return nil, false, err
	}
	resp := <-respC
	idIndexMap := make(map[model.Id]int)
	ids := make([]model.Id, len(resp.Candidates))
	for i, candidate := range resp.Candidates {
		idIndexMap[candidate.QuestionId] = i
		ids[i] = candidate.QuestionId
	}
	sess, err := db.NewSession()
	if err != nil {
		return nil, false, err
	}
	defer sess.Close()
	accessor, err := db.NewAccessor(sess)
	if err != nil {
		return nil, false, err
	}
	questionAccessor, err := wrapper.NewQuestionAccessor(accessor)
	if err != nil {
		return nil, false, err
	}
	outC, resC, err := questionAccessor.ScanByIds(ids, 5, nil)
	if err != nil {
		return nil, false, err
	}
	res = make([]*SimilarQuestion, len(resp.Candidates))
	count := 0
	for out := range outC {
		i := idIndexMap[out.QuestionId]
		res[i] = &SimilarQuestion{
			Question:   out,
			Similarity: resp.Candidates[i].Similarity,
		}
		count++
	}
	if count < len(resp.Candidates) {
		resHasNil := res
		res = make([]*SimilarQuestion, 0, len(resHasNil))
		for _, r := range resHasNil {
			if r == nil {
				continue
			}
			res = append(res, r)
		}
	}
	return res, resp.IsTimeout, <-resC
}
