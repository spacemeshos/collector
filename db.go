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

func (db *Database) AllBlocks() ([]*events.NewBlock, error){
	var blocks []*events.NewBlock
	err := db.inst.Model(&blocks).Select()
	return blocks, err
}



type Database struct {
	inst *pg.DB
}

func NewDb() *Database {
	d := Database{nil}
	d.inst = pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "mysecretpassword",
	})
	return &d
}


func (db *Database) Close() error{
	return db.inst.Close()
}

func (db *Database) createTables(tempDb bool) error{
	for _, model := range []interface{}{(*events.NewBlock)(nil),
		(*events.ValidBlock)(nil),
		(*events.NewTx)(nil),
		(*events.ValidTx)(nil),
		(*events.NewAtx)(nil),
		(*events.ValidAtx)(nil)} {
		err := db.inst.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists:true,
			Temp:tempDb,
		})
		if err != nil {
			return err
		}
	}
	return nil
}


