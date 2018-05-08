package maintainer

import (
	"log"

	"github.com/donyori/cqa/data/db"
	"github.com/donyori/cqa/data/db/id"
)

func EnsureDataTypes() error {
	log.Println("Ensure data types.")
	session, err := db.NewSession()
	if err != nil {
		log.Println(err)
		return err
	}
	defer session.Close()
	log.Println("*** Succeed to connect to database.")
	maintainer, err := db.NewMaintainer(session)
	if err != nil {
		log.Println(err)
		return err
	}
	for _, cid := range id.ValidCollectionIds {
		err = maintainer.EnsureDataTypes(cid)
		if err != nil {
			log.Println(err)
		}
	}
	log.Println("Done.")
	return nil
}
