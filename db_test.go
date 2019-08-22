package collector

import (
	"github.com/spacemeshos/go-spacemesh/events"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_WriteObjectToDB(t *testing.T){
	db := NewDb("postgres", "mysecretpassword")
	defer func() { assert.NoError(t, db.Close()) }()
	err := db.createTables(true)
	assert.NoError(t, err)

	err = db.StoreBlock(&events.NewBlock{Id:1234,Layer:1,Atx:"1234"})
	assert.NoError(t, err)

	blocks, err := db.AllBlocks()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(blocks))

	err = db.StoreBlockValid(&events.ValidBlock{Id:1234,Valid:true})
	assert.NoError(t, err)

	err = db.StoreAtx(&events.NewAtx{Id:"1234"})
	assert.NoError(t, err)

	err = db.StoreReward(&events.RewardReceived{Coinbase:"1234", Amount:6})
	assert.NoError(t, err)
}

func Test_ReadQuery(t *testing.T){
	db := NewDb("postgres", "mysecretpassword")
	defer func() { assert.NoError(t, db.Close()) }()
	err := db.createTables(true)
	assert.NoError(t, err)

	err = db.StoreTx(&events.NewTx{Id:"1234", Origin:"1111", Destination:"2222", Gas:10, Amount:20})
	assert.NoError(t, err)
	err = db.StoreTx(&events.NewTx{Id:"1235", Origin:"1111", Destination:"3333", Gas:10, Amount:20})
	assert.NoError(t, err)
	err = db.StoreTx(&events.NewTx{Id:"1236", Origin:"2222", Destination:"2222", Gas:10, Amount:20})
	assert.NoError(t, err)
	err = db.StoreTx(&events.NewTx{Id:"1237", Origin:"7777", Destination:"2222", Gas:10, Amount:20})
	assert.NoError(t, err)

	txs, err := db.GetTransactionsFrom("1111")
	assert.Equal(t, 2, len(txs))

	txs, err = db.GetTransactionsTo("2222")
	assert.Equal(t, 2, len(txs))

}
