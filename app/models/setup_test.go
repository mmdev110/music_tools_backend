package models

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	db, err := InitTestDB()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	TestDB = db
	//ClearTestDB(TestDB)
	os.Exit(m.Run())
}
