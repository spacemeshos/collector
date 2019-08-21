package collector

import "github.com/spacemeshos/smutil/log"

func main() {


	url := "tcp://localhost:56565"
	db := NewDb()
	err := db.createTables(false)
	if err != nil {
		log.Error("cannot create DB %v ", err)
		return
	}

	c := NewCollector(db, url)
	c.Start(true)

}
