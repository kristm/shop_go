package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddSubscriber(t *testing.T) {
	err := AddSubscriber("test@example.com")
	assert.NoError(t, err)
}

func TestAddDuplicateSubscriber(t *testing.T) {
	err := AddSubscriber("test@example.com")
	assert.Equal(t, nil, err)

	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM subscribers").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query test db: %v", err)
	}

	expected := 1
	assert.Equal(t, expected, count)
}
