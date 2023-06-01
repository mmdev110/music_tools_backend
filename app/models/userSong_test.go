package models

import (
	"testing"
)

func TestUserSong(t *testing.T) {
	err := Init(true)
	if err != nil {
		t.Fatal(err)
	}
	defer ClearSQLiteDB()

	t.Run("save new UserSongSection", func(t *testing.T) {
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
