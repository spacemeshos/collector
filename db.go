package collector

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/spacemeshos/go-spacemesh/events"
)

func (db *Database) StoreBlock(event *events.NewBlock) error {
	return db.inst.Insert(event)
}
func (db *Database) StoreBlockValid(event *events.ValidBlock) error {
	return db.inst.Insert(event)
}
func (db *Database) StoreTx(event *events.NewTx) error {
	return db.inst.Insert(event)
}
func (db *Database) StoreTxValid(event *events.ValidTx) error {
	return db.inst.Insert(event)
}
func (db *Database) StoreAtx(event *events.NewAtx) error {
	return db.inst.Insert(event)
}
func (db *Database) StoreAtxValid(event *events.ValidAtx) error {
	return db.inst.Insert(event)
}
func (db *Database) StoreReward(event *events.RewardReceived) error {
	return db.inst.Insert(event)
}

func (db *Database) AllBlocks() ([]*events.NewBlock, error) {
	var blocks []*events.NewBlock
	err := db.inst.Model(&blocks).Select()
	return blocks, err
}

func (db *Database) GetTransactionsFrom(origin string) ([]*events.NewTx, error) {
	var blocks []*events.NewTx
	err := db.inst.Model(&blocks).Where("origin = ?", origin).Select()
	return blocks, err

}

func (db *Database) GetTransactionsTo(origin string) ([]*events.NewTx, error) {
	var blocks []*events.NewTx
	err := db.inst.Model(&blocks).Where("origin = ?", origin).Select()
	return blocks, err

}

// Database contains a conneciton to to a postgress DB
type Database struct {
	inst *pg.DB
}

// NewDb creates a new database connection using userName and password to the database instance.
// It will connect to postgress DB on port 5432 (Postgres default port).
func NewDb(userName, pass string) *Database {
	d := Database{nil}
	d.inst = pg.Connect(&pg.Options{
		User:     userName,
		Password: pass,
	})
	return &d
}

// Start tries to initialize the tables inside the db.
func (db *Database) Start() error {
	return db.createTables(false)
}

// Close closes active connections to the db.
func (db *Database) Close() error {
	return db.inst.Close()
}

// createTables uses pg orm to create a database based on the objects provided.
// it can generate a temporary DB for testing purposes.
func (db *Database) createTables(tempDb bool) error {
	for _, model := range []interface{}{(*events.NewBlock)(nil),
		(*events.ValidBlock)(nil),
		(*events.NewTx)(nil),
		(*events.ValidTx)(nil),
		(*events.NewAtx)(nil),
		(*events.ValidAtx)(nil),
		(*events.RewardReceived)(nil)} {
		err := db.inst.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists: true,
			Temp:        tempDb,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
