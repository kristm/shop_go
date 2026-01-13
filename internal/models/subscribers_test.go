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

func TestUnsubscribe(t *testing.T) {
	err := Unsubscribe("test@example.com")
	assert.NoError(t, err)

	var email string
	err = DB.QueryRow("SELECT email FROM subscribers where id LIKE ?", 1).Scan(&email)
	if err != nil {
		t.Fatalf("Failed to query test db: %v", err)
	}
	// email suffixed with id
	expected := "test@example.com-1"
	assert.Equal(t, expected, email)
}

func TestUnsubscribeNonExistingEmail(t *testing.T) {
	err := Unsubscribe("test2@example.com")
	assert.Error(t, err)
}
