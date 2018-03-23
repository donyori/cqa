package mongodb

import (
	"gopkg.in/mgo.v2"
)

type MgoMaintainer struct {
	MgoConnector
}

func NewMgoMaintainer(settings *MongoDbSettings) *MgoMaintainer {
	return &MgoMaintainer{MgoConnector: *NewMgoConnector(settings)}
}

func (mm *MgoMaintainer) EnsureIndexes(isBackground bool) error {
	mm.RLock()
	defer mm.RUnlock()
	if !mm.isConnectedWithoutLock() {
		return ErrNotConnected
	}
	settings := mm.getSettings()
	db := mm.Session.DB(settings.DbName)
	c := db.C(settings.CNames[MgoCNameKeyQ])

	titleIndex := mgo.Index{
		Key:        []string{"title"},
		Background: isBackground,
	}
	err := c.EnsureIndex(titleIndex)
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
	return err
}
