package collector

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/spacemeshos/go-spacemesh/events"
)

func (db *Database) StoreBlock(event *events.NewBlockEvent) error {
	return db.inst.Insert(event)
}
func (db *Database) StoreBlockValid(event *events.ValidBlockEvent) error {
	return db.inst.Insert(event)
}
func (db *Database) StoreTx(event *events.NewTxEvent) error {
	return db.inst.Insert(event)
}
func (db *Database) StoreTxValid(event *events.ValidTxEvent) error {
	return db.inst.Insert(event)
}
func (db *Database) StoreAtx(event *events.NewAtxEvent) error {
	return db.inst.Insert(event)
}
func (db *Database) StoreAtxValid(event *events.ValidAtxEvent) error {
	return db.inst.Insert(event)
}

func (db *Database) AllBlocks() ([]*events.NewBlockEvent, error){
	var blocks []*events.NewBlockEvent
	err := db.inst.Model(&blocks).Select()
	return blocks, err
}



type Database struct {
	inst *pg.DB
}

func NewDbNow() *Database {
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

func (db *Database) createTables() error{
	for _, model := range []interface{}{(*events.NewBlockEvent)(nil),
		(*events.ValidBlockEvent)(nil),
		(*events.NewTxEvent)(nil),
		(*events.ValidTxEvent)(nil),
		(*events.NewAtxEvent)(nil),
		(*events.ValidAtxEvent)(nil)} {
		err := db.inst.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists:true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}


