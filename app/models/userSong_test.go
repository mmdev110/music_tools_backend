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
		uid := uint(9999)
		tag1 := UserTag{
			UserId:    uid,
			Name:      "tag1",
			SortOrder: 0,
			UserSongs: []UserSong{},
		}
		if err := tag1.Create(); err != nil {
			t.Errorf("error at create %v", err)
		}
		tag2 := UserTag{
			UserId:    uid,
			Name:      "tag2",
			SortOrder: 0,
			UserSongs: []UserSong{},
		}
		if err := tag2.Create(); err != nil {
			t.Errorf("error at create %v", err)
		}
		genre1 := UserGenre{
			UserId:    uid,
			Name:      "genre1",
			SortOrder: 0,
			UserSongs: []UserSong{},
		}
		if err := genre1.Create(); err != nil {
			t.Errorf("error at create %v", err)
		}
		genre2 := UserGenre{
			UserId:    uid,
			Name:      "genre2",
			SortOrder: 0,
			UserSongs: []UserSong{},
		}
		if err := genre2.Create(); err != nil {
			t.Errorf("error at create %v", err)
		}
		us := UserSong{
			UserId: uid,
			Title:  "title1",
			Artist: "artist1",
			Memo:   "memo1",
			Genres: []UserGenre{genre1, genre2},
			Tags:   []UserTag{tag1, tag2},
			Audio: UserSongAudio{
				Name: "song1",
				Url:  Url{Get: "get1", Put: "put1"},
			},
			Sections: []UserSongSection{{
				Name:            "intro1",
				ProgressionsCSV: "Am7,F,G,C",
				Key:             0,
				BPM:             120,
				Scale:           "メジャー",
				Memo:            "sectionMemo1",
				LoopRange:       LoopRange{Start: 10, End: 20},
			}, {
				Name:            "intro2",
				ProgressionsCSV: "Am7,F,G,C",
				Key:             0,
				BPM:             140,
				Scale:           "マイナー",
				Memo:            "sectionMemo2",
				LoopRange:       LoopRange{Start: 30, End: 40},
			}},
		}
		if err := us.Create(); err != nil {
			t.Errorf("error at create %v", err)
		}
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
