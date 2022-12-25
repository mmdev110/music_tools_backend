package models

import (
	"fmt"
	"testing"

	"example.com/app/utils"
)

func TestUserLoop(t *testing.T) {
	err := Init()
	if err != nil {
		t.Fatal(err)
	}
	t.Run("UserLoop.GetByID", func(t *testing.T) {
		var ul = UserLoop{}
		if err := ul.GetByID(1); err != nil {
			t.Errorf("error at GetByID: %v", err)
		}
		fmt.Println("=======")
		utils.PrintStruct(ul.UserLoopTags)
	})
	t.Run("UserLoop.GetByUserId", func(t *testing.T) {
		var ul = UserLoop{}
		loops, err := ul.GetByUserId(4, ULSearchCond{})
		if err != nil {
			t.Errorf("error at GetByID: %v", err)
		}
		fmt.Println("=======")
		for _, v := range loops {
			utils.PrintStruct(v.UserLoopTags)
		}
	})
	t.Run("save new UserLoop", func(t *testing.T) {
		t.Skip()
		userId := uint(2)
		ul := UserLoop{}
		ul.UserId = userId
		ul.Name = "test save"
		ul.UserLoopTags = []UserLoopTag{{
			ID:     uint(79),
			UserId: userId,
			Name:   "user2",
		}}
		ul.Create()
	})

}
func TestGetUserLoop(t *testing.T) {
	err := Init()
	if err != nil {
		t.Fatal(err)
	}
	t.Run("UserLoop.GetByID", func(t *testing.T) {
		var ul = UserLoop{}
		if err := ul.GetByID(1); err != nil {
			t.Errorf("error at GetByID: %v", err)
		}
		fmt.Println("=======")
		utils.PrintStruct(ul.UserLoopTags)
	})

}

func TestUpdateUserLoop(t *testing.T) {
	err := Init()
	if err != nil {
		t.Fatal(err)
	}
	t.Run("update existing UserLoop", func(t *testing.T) {
		ul := UserLoop{}
		if err := ul.GetByID(1); err != nil {
			t.Errorf("%v", err)
		}
		tag_to_delete := ul.UserLoopTags[0]
		fmt.Println("tag_to_delete")
		utils.PrintStruct(tag_to_delete)

		ul.DeleteTagRelations([]UserLoopTag{tag_to_delete})

	})
}
