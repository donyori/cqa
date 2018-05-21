package mongodb

import (
	"errors"

	"gopkg.in/mgo.v2"

	"github.com/donyori/cqa/data/db/generic"
	"github.com/donyori/cqa/data/db/id"
)

type Maintainer struct {
	WithSession
}

var ErrNilMaintainer error = errors.New("MongoDB maintainer is nil")

func NewMaintainer(session generic.Session) (
	maintainer *Maintainer, err error) {
	maintainer = new(Maintainer)
	if session != nil {
		err = maintainer.SetSession(session)
		if err != nil {
			return nil, err
		}
	}
	return maintainer, nil
}

func (mm *Maintainer) EnsureIndexes(cid id.CollectionId,
	isBackground bool) error {
	if mm == nil {
		return ErrNilMaintainer
	}
	session, c, err := mm.aquireSessionAndCollection(cid)
	if err != nil {
		return err
	}
	defer session.Release()

	switch cid {
	case id.QuestionCollection:
		titleIndex := mgo.Index{
			Key:        []string{"title"},
			Background: isBackground,
		}
		err = c.EnsureIndex(titleIndex)
		if err != nil {
			return err
		}

		scoreIndex := mgo.Index{
			Key:        []string{"score"},
			Background: isBackground,
		}
		err = c.EnsureIndex(scoreIndex)
		if err != nil {
			return err
		}

		viewCountIndex := mgo.Index{
			Key:        []string{"view_count"},
			Background: isBackground,
		}
		err = c.EnsureIndex(viewCountIndex)
		if err != nil {
			return err
		}

		creationDateIndex := mgo.Index{
			Key:        []string{"creation_date"},
			Background: isBackground,
		}
		err = c.EnsureIndex(creationDateIndex)
		if err != nil {
			return err
		}

		lastActivityDateIndex := mgo.Index{
			Key:        []string{"last_activity_date"},
			Background: isBackground,
		}
		err = c.EnsureIndex(lastActivityDateIndex)
		if err != nil {
			return err
		}

		lastEditDateIndex := mgo.Index{
			Key:        []string{"last_edit_date"},
			Background: isBackground,
		}
		err = c.EnsureIndex(lastEditDateIndex)
		if err != nil {
			return err
		}

		lastCreateOrEditDateIndex := mgo.Index{
			Key:        []string{"last_create_or_edit_date"},
			Background: isBackground,
		}
		err = c.EnsureIndex(lastCreateOrEditDateIndex)
		if err != nil {
			return err
		}

		tagsIndex := mgo.Index{
			Key:        []string{"tags"},
			Sparse:     true,
			Background: isBackground,
		}
		err = c.EnsureIndex(tagsIndex)
		if err != nil {
			return err
		}

	case id.QuestionClassificationCollection:
		cByTagIndex := mgo.Index{
			Key:        []string{"classification_by_tag"},
			Sparse:     true,
			Background: isBackground,
		}
		err = c.EnsureIndex(cByTagIndex)
		if err != nil {
			return err
		}

		cByNnIndex := mgo.Index{
			Key:        []string{"classification_by_nn"},
			Sparse:     true,
			Background: isBackground,
		}
		err = c.EnsureIndex(cByNnIndex)
		if err != nil {
			return err
		}
	}

	return nil
}

func (mm *Maintainer) EnsureDataTypes(cid id.CollectionId) error {
	if mm == nil {
		return ErrNilMaintainer
	}
	if cid != id.QuestionCollection &&
		cid != id.QuestionVectorCollection &&
		cid != id.QuestionClassificationCollection {
		return nil
	}
	session := mm.GetSession()
	accessor, err := NewAccessor(session)
	if err != nil {
		return err
	}
	quitC := make(chan struct{}, 1)
	defer close(quitC)
	outC, resC, err := accessor.Scan(cid, nil, 4, quitC, nil)
	for data := range outC {
		_, err = accessor.SaveOne(cid, nil, data)
		if err != nil {
			quitC <- struct{}{}
			return err
		}
	}
	quitC <- struct{}{}
	return <-resC
}
