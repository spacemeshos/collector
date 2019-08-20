package collector

import (
	"github.com/spacemeshos/go-spacemesh/events"
	"github.com/spacemeshos/go-spacemesh/log"
	"github.com/spacemeshos/go-spacemesh/types"
)

type collector struct{
	stop chan struct{}
	db Db
}

func NewCollector(db Db) *collector{
	return &collector{make(chan struct{}),db}
}

type Db interface {
	StoreBlock(event events.NewBlockEvent) error
	StoreBlockValid(event events.ValidBlockEvent) error
	StoreTx(event events.NewTxEvent) error
	StoreTxValid(event events.ValidTxEvent) error
	StoreAtx(event events.NewAtxEvent) error
	StoreAtxValid(event events.ValidAtxEvent) error
}

func (c *collector) collectEvents(url string) {
	sub, err := events.NewSubscriber(url)
	if err != nil {
		log.Info("cannot start subscriber")
		return
	}
	blocks, err := sub.Subscribe(events.NewBlock)
	blocksValid, err := sub.Subscribe(events.BlockValid)
	txs, err := sub.Subscribe(events.NewTx)
	txValid, err := sub.Subscribe(events.TxValid)
	atxs, err := sub.Subscribe(events.NewAtx)
	atxsValid, err := sub.Subscribe(events.AtxValid)
	if err != nil {
		log.Info("cannot start subscriber")
		return
	}

	loop:
	for {
		select {
		case data  := <- blocks:
			var e events.NewBlockEvent
			err := types.BytesToInterface(data, e)
			if err != nil {
				log.Error("cannot parse received message %v", err)
			}
			err = c.db.StoreBlock(e)
			if err != nil {
				log.Error("cannot write message %v", err)
			}
		case data  := <- blocksValid:
			var e events.ValidBlockEvent
			err := types.BytesToInterface(data, e)
			if err != nil {
				log.Error("cannot parse received message %v", err)
			}
			err = c.db.StoreBlockValid(e)
			if err != nil {
				log.Error("cannot write message %v", err)
			}
		case data  := <- txs:
			var e events.NewTxEvent
			err := types.BytesToInterface(data, e)
			if err != nil {
				log.Error("cannot parse received message %v", err)
			}
			err = c.db.StoreTx(e)
			if err != nil {
				log.Error("cannot write message %v", err)
			}
		case data  := <- txValid:
			var e events.ValidTxEvent
			err := types.BytesToInterface(data, e)
			if err != nil {
				log.Error("cannot parse received message %v", err)
			}
			err = c.db.StoreTxValid(e)
			if err != nil {
				log.Error("cannot write message %v", err)
			}
		case data  := <- atxs:
			var e events.NewAtxEvent
			err := types.BytesToInterface(data, e)
			if err != nil {
				log.Error("cannot parse received message %v", err)
			}
			err = c.db.StoreAtx(e)
			if err != nil {
				log.Error("cannot write message %v", err)
			}
		case data  := <- atxsValid:
			var e events.ValidAtxEvent
			err := types.BytesToInterface(data, e)
			if err != nil {
				log.Error("cannot parse received message %v", err)
			}
			err = c.db.StoreAtxValid(e)
			if err != nil {
				log.Error("cannot write message %v", err)
			}
			case <-c.stop:
				break loop
		}
	}
}