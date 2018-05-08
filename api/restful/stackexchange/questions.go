package stackexchange

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/donyori/cqa/common/cmp"
	"github.com/donyori/cqa/common/conv"
)

var (
	ErrMinOrMaxNotAllowed error = errors.New("min or max param is not allowed")
)

func Questions(page int, pageSize int, fromDate *time.Time, toDate *time.Time,
	sort QuestionsSort, order Order, min interface{}, max interface{},
	tagged string) (res *QuestionsResponse, err error) {
	r := GetSeGetClient().R()
	r.SetQueryParam("page", strconv.Itoa(page))
	if pageSize > 0 {
		if pageSize > GlobalSettings.MaxPageSize {
			pageSize = GlobalSettings.MaxPageSize
		}
		r.SetQueryParam("pagesize", strconv.Itoa(pageSize))
	}
	if fromDate != nil {
		r.SetQueryParam("fromdate", strconv.FormatInt(fromDate.Unix(), 10))
	}
	if toDate != nil {
		r.SetQueryParam("todate", strconv.FormatInt(toDate.Unix(), 10))
	}
	if !sort.IsValid() {
		return nil, ErrUnknownQuestionsSort
	}
	r.SetQueryParam("sort", sort.String())
	if !order.IsValid() {
		return nil, ErrUnknownOrder
	}
	r.SetQueryParam("order", order.String())
	var minInt64, maxInt64 int64
	switch sort {
	case QuestionsSortActivity:
		fallthrough
	case QuestionsSortCreation:
		if !cmp.IsNilDeeply(min) {
			minInt64, err = conv.InterfaceToTimestamp(min)
			if err != nil {
				return nil, err
			}
		}
		if !cmp.IsNilDeeply(max) {
			maxInt64, err = conv.InterfaceToTimestamp(max)
			if err != nil {
				return nil, err
			}
		}
	case QuestionsSortVotes:
		if !cmp.IsNilDeeply(min) {
			minInt64, err = conv.InterfaceToInt64(min)
			if err != nil {
				return nil, err
			}
		}
		if !cmp.IsNilDeeply(max) {
			maxInt64, err = conv.InterfaceToInt64(max)
			if err != nil {
				return nil, err
			}
		}
	default:
		if !cmp.IsNilDeeply(min) || !cmp.IsNilDeeply(max) {
			return nil, ErrMinOrMaxNotAllowed
		}
	}
	if !cmp.IsNilDeeply(min) {
		r.SetQueryParam("min", strconv.FormatInt(minInt64, 10))
	}
	if !cmp.IsNilDeeply(max) {
		r.SetQueryParam("max", strconv.FormatInt(maxInt64, 10))
	}
	if tagged != "" {
		r.SetQueryParam("tagged", tagged)
	}
	resp, err := r.Get(
		"/" + GlobalSettings.ApiVersion + GlobalSettings.QuestionsPath)
	if err != nil {
		return nil, err
	}
	res = new(QuestionsResponse)
	err = json.Unmarshal(resp.Body(), res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
