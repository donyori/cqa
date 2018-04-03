package mongodb

import (
	"errors"

	"gopkg.in/mgo.v2"

	"github.com/donyori/cqa/data/db/generic"
	"github.com/donyori/cqa/data/db/id"
)

type MgoMaintainer struct {
	WithMgoSession
}

var ErrNilMaintainer error = errors.New("MgoMaintainer is nil")

func NewMgoMaintainer(session generic.Session) (
	maintainer *MgoMaintainer, err error) {
	maintainer = new(MgoMaintainer)
	if session != nil {
		err = maintainer.SetSession(session)
		if err != nil {
			return nil, err
		}
	}
	return maintainer, nil
}

func (mm *MgoMaintainer) EnsureIndexes(cid id.CollectionId,
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

		tagsIndex := mgo.Index{
			Key:        []string{"tags"},
			Sparse:     true,
			Background: isBackground,
		}
		err = c.EnsureIndex(tagsIndex)
		if err != nil {
			return err
		}
	}

	return nil
}

func (mm *MgoMaintainer) EnsureDataTypes(cid id.CollectionId) error {
	if mm == nil {
		return ErrNilMaintainer
	}
	if cid != id.QuestionCollection && cid != id.QuestionVectorCollection {
		return nil
	}
	session := mm.GetSession()
	accessor, err := NewMgoAccessor(session)
	if err != nil {
		return err
	}
	quitC := make(chan struct{}, 1)
	defer close(quitC)
	outC, resC, err := accessor.Scan(cid, nil, 4, quitC, nil)
	for data := range outC {
		_, err = accessor.Save(cid, nil, data)
		if err != nil {
			quitC <- struct{}{}
			return err
		}
	}
	quitC <- struct{}{}
	return <-resC
}
