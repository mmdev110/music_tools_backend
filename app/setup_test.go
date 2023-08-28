package main

import (
	"os"
	"testing"

	"gorm.io/gorm"
)

var TestDB *gorm.DB

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}
