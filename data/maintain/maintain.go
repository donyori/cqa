package maintain

import (
	"log"

	"github.com/donyori/cqa/data/db"
)

func EnsureIndexes(isBackground bool) {
	log.Println("Ensure indexes.")
	mm, err := db.NewMaintainer()
	if err != nil {
		log.Fatalln(err)
		return
	}
	err = mm.Connect()
	if err != nil {
		log.Fatalln(err)
		return
	}
	defer mm.Close()
	err = mm.EnsureIndexes(isBackground)
	if err != nil {
		log.Fatalln(err)
		return
	}
	log.Println("Finish")
}
