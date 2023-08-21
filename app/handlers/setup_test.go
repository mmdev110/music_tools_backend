package handlers

import (
	"fmt"
	"os"
	"testing"

	"example.com/app/models"
)

var h = Base{}

func TestMain(m *testing.M) {
	db, err := models.InitTestDB()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	h.DB = db
	os.Exit(m.Run())
}