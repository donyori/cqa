package maintainer

import (
	"log"

	"github.com/donyori/cqa/data/db"
	"github.com/donyori/cqa/data/db/id"
)

func EnsureIndexes() error {
	log.Println("Ensure indexes.")
	session, err := db.NewSession()
	if err != nil {
		log.Println(err)
		return err
	}
	defer session.Close()
	log.Println("*** Connect to database successfully.")
	maintainer, err := db.NewMaintainer(session)
	if err != nil {
		log.Println(err)
		return err
	}
	isBackground := GlobalSettings.EnsureIndexes.IsBackground
	for _, cid := range id.ValidCollectionIds {
		err = maintainer.EnsureIndexes(cid, isBackground)
		if err == nil {
			log.Printf("*** Ensure indexes of collection %q successfully.\n",
				cid.String())
		} else {
			log.Println(err)
		}
	}
	log.Println("Done.")
	return nil
}
