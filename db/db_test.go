package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	db, err := connect(":memory:")
	assert.Nil(t, err)
	assert.True(t, ping(*db))
}

func TestSetup(t *testing.T) {
	db, err := connect(":memory:")
	assert.Nil(t, err)
	assert.Nil(t, setup(db))
}
