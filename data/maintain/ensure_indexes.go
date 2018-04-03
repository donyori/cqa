package maintain

import (
	"log"

	"github.com/donyori/cqa/data/db"
	"github.com/donyori/cqa/data/db/id"
)

func EnsureIndexes() error {
	log.Println("Ensure indexes.")
	session, err := db.NewSession()
	if err != nil {
		log.Fatalln(err)
		return err
	}
	defer session.Close()
	log.Println("*** Succeed to connect to database.")
	maintainer, err := db.NewMaintainer(session)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	isBackground := GlobalSettings.EnsureIndexes.IsBackground
	for _, cid := range id.ValidCollectionIds {
		err = maintainer.EnsureIndexes(cid, isBackground)
		if err != nil {
			log.Fatalln(err)
		}
	}
	log.Println("Done.")
	return nil
}
