package models

import (
	"fmt"
	"testing"

	"example.com/app/utils"
)

func TestUserSongSection(t *testing.T) {
	err := Init()
	if err != nil {
		t.Fatal(err)
	}
	t.Run("UserSongSection.GetByID", func(t *testing.T) {
		var us = UserSong{}
		if err := us.GetByID(1); err != nil {
			t.Errorf("error at GetByID: %v", err)
		}
		fmt.Println("=======")
		utils.PrintStruct(us.Tags)
	})
	t.Run("UserSongSection.GetByUserId", func(t *testing.T) {
		var us = UserSong{}
		loops, err := us.GetByUserId(4, ULSearchCond{})
		if err != nil {
			t.Errorf("error at GetByID: %v", err)
		}
		fmt.Println("=======")
		for _, v := range loops {
			utils.PrintStruct(v.Tags)
		}
	})
	t.Run("save new UserSongSection", func(t *testing.T) {
		t.Skip()
		userId := uint(2)
		us := UserSong{}
		us.UserId = userId
		us.Title = "test save"
		us.Tags = []UserTag{{
			ID:     uint(79),
			UserId: userId,
			Name:   "user2",
		}}
		us.Create()
	})

}
func TestGetUserSongSection(t *testing.T) {
	err := Init()
	if err != nil {
		t.Fatal(err)
	}
	t.Run("UserSongSection.GetByID", func(t *testing.T) {
		var us = UserSong{}
		if err := us.GetByID(1); err != nil {
			t.Errorf("error at GetByID: %v", err)
		}
		fmt.Println("=======")
		utils.PrintStruct(us.Tags)
	})

}

func TestUpdateUserSongSection(t *testing.T) {
	err := Init()
	if err != nil {
		t.Fatal(err)
	}
	t.Run("update existing UserSongSection", func(t *testing.T) {
		us := UserSong{}
		if err := us.GetByID(1); err != nil {
			t.Errorf("%v", err)
		}
		tag_to_delete := us.Tags[0]
		fmt.Println("tag_to_delete")
		utils.PrintStruct(tag_to_delete)

		us.DeleteTagRelations([]UserTag{tag_to_delete})

	})
}
