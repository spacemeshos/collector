package collector

import (
	"github.com/spacemeshos/go-spacemesh/events"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_WriteObjectToDB(t *testing.T){
	db := NewDb()
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
}
