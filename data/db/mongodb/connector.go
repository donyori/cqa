package mongodb

import (
	"errors"
	"sync"

	"gopkg.in/mgo.v2"
)

type MgoConnector struct {
	sync.RWMutex
	Session *mgo.Session

	settings *MongoDbSettings
}

var (
	ErrNotConnected error = errors.New("MongoDB is not connected")
)

func NewMgoConnector(settings *MongoDbSettings) *MgoConnector {
	return &MgoConnector{settings: settings}
}

func (mc *MgoConnector) IsConnected() bool {
	if mc == nil {
		return false
	}
	mc.RLock()
	defer mc.RUnlock()
	return mc.isConnectedWithoutLock()
}

func (mc *MgoConnector) Connect() error {
	if mc == nil {
		return errors.New("MongoDB connector is nil")
	}
	mc.Lock()
	defer mc.Unlock()
	sess, err := mgo.Dial(mc.getSettings().Url)
	if err != nil {
		return err
	}
	mc.Session = sess
	return nil
}

func (mc *MgoConnector) Close() {
	if mc == nil {
		return
	}
	mc.Lock()
	defer mc.Unlock()
	if mc.Session != nil {
		mc.Session.Close()
		mc.Session = nil
	}
}

func (mc *MgoConnector) isConnectedWithoutLock() bool {
	if mc == nil {
		return false
	}
	if mc.Session == nil {
		return false
	}
	return mc.Session.Ping() == nil
}

func (mc *MgoConnector) getSettings() *MongoDbSettings {
	if mc != nil && mc.settings != nil {
		return mc.settings
	} else {
		return &GlobalSettings
	}
}
