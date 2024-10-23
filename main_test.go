package main

import (
	"log"
	"os"
	"testing"

	"github.com/go-playground/assert/v2"
)

func setupSuite(tb testing.TB) func(tb testing.TB) {
	log.Println("setup suite")

	return func(tb testing.TB) {
		log.Println("teardown suite")
	}
}

func setupTest(tb testing.TB) func(tb testing.TB) {
	log.Println("setup test")

	return func(tb testing.TB) {
		log.Println("teardown test")
	}
}

func TestMain(m *testing.M) {
	log.Println("Test Main")
	code := m.Run()
	os.Exit(code)
}

func TestNothing(t *testing.T) {
	assert.Equal(t, true, true)
}
