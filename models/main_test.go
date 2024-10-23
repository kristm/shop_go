package models

import (
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.Println("Test Models Main")
	ConnectDatabase()
	code := m.Run()

	os.Exit(code)
}
