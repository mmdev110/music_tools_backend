package models

import "testing"

func Test_InsertTestData(t *testing.T) {
	t.Run("testing InsertTestData", func(t *testing.T) {
		tx := TestDB.Begin()
		defer tx.Rollback()

		_, err := InsertTestData(tx)
		if err != nil {
			t.Error(err)
		}
	})

}
