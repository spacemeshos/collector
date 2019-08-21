package collector

import (
	"github.com/spacemeshos/go-spacemesh/events"
	"github.com/spacemeshos/go-spacemesh/log"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type MockDb struct {
	msgs map[byte]int
	total int
}

func (m *MockDb) StoreBlock(event *events.NewBlockEvent) error {
	m.msgs[1]++
	m.total++
	return nil
}

func (m *MockDb) StoreBlockValid(event *events.ValidBlockEvent) error {
	m.msgs[2]++
	m.total++
	return nil
}

func (m *MockDb) StoreTx(event *events.NewTxEvent) error {
	m.msgs[3]++
	m.total++
	return nil
}

func (m *MockDb) StoreTxValid(event *events.ValidTxEvent) error {
	m.msgs[4]++
	m.total++
	return nil
}

func (m *MockDb) StoreAtx(event *events.NewAtxEvent) error {
	m.msgs[5]++
	m.total++
	return nil
}

func (m *MockDb) StoreAtxValid(event *events.ValidAtxEvent) error {
	m.msgs[6]++
	m.total++
	return nil
}

func TestCollectEvents(t *testing.T) {
	url := "tcp://localhost:56565"
	m := &MockDb{make(map[byte]int), 0}
	c := NewCollector(m, url)


	eventPublisher, err := events.NewEventPublisher(url)
	assert.NoError(t, err)
	defer func() {
		assert.NoError(t, eventPublisher.Close())
	}()

	c.Start(false)
	time.Sleep(2 * time.Second)
	orig := events.NewBlockEvent{Layer: 1, Block: 234}
	err = eventPublisher.PublishEvent(orig)

	orig1 := events.ValidBlockEvent{Block: 234, Valid: true}
	err = eventPublisher.PublishEvent(orig1)

	orig2 := events.NewAtxEvent{AtxId: "1234"}
	err = eventPublisher.PublishEvent(orig2)

	orig3 := events.ValidAtxEvent{AtxId: "1234", Valid: true}
	err = eventPublisher.PublishEvent(orig3)

	orig4 := events.NewTxEvent{TxId: "4321", Gas: 20, Amount: 400, Destination: "1234567", Origin: "876543"}
	err = eventPublisher.PublishEvent(orig4)

	orig5 := events.ValidTxEvent{TxId:"4321", Valid:true}
	err = eventPublisher.PublishEvent(orig5)

	time.Sleep(1 * time.Second)
	c.Stop()

	log.Info("got %v", len(m.msgs))
	assert.Equal(t, 6, m.total)
	assert.Equal(t, 6, len(m.msgs))
}
