package collector

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/spacemeshos/go-spacemesh/events"
)

func (db *db) StoreBlock(event events.NewBlockEvent) error {
	return db.inst.Insert(event)
}
func (db *db) StoreBlockValid(event events.ValidBlockEvent) error {
	return db.inst.Insert(event)
}
func (db *db) StoreTx(event events.NewTxEvent) error {
	return db.inst.Insert(event)
}
func (db *db) StoreTxValid(event events.ValidTxEvent) error {
	return db.inst.Insert(event)
}
func (db *db) StoreAtx(event events.NewAtxEvent) error {
	return db.inst.Insert(event)
}
func (db *db) StoreAtxValid(event events.ValidAtxEvent) error {
	return db.inst.Insert(event)
}

func (db *db) AllBlocks() ([]events.NewBlockEvent, error){
	var blocks []events.NewBlockEvent
	err := db.inst.Select(&blocks)
	return blocks, err
}



type db struct {
	inst *pg.DB
}

func NewDb() *db {
	d := db{nil}
	d.inst = pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "mysecretpassword",
	})
	return &d
}


func (db *db) Close() error{
	return db.inst.Close()
}

func (db *db) createTables() error{
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


