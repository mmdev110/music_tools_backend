package models

import (
	"testing"
)

func TestUserTag(t *testing.T) {
	t.Run("can remove a tag associated with 2 songs", func(t *testing.T) {
		//data := prepareTestData(t)
		defer ClearTestDB(DB)
		var uid = uint(9999)
		var user = User{ID: uid}
		DB.Create(&user)
		tag1 := UserTag{
			UserId:    uid,
			Name:      "tag1",
			SortOrder: 0,
			UserSongs: []UserSong{},
		}
		if err := tag1.Create(DB); err != nil {
			t.Errorf("error at create %v", err)
		}
		us1 := UserSong{
			UserId: uid,
			Title:  "song1",
			Genres: []UserGenre{},
			Tags:   []UserTag{tag1},
		}
		if err := us1.Create(DB); err != nil {
			t.Errorf("error at create %v", err)
		}
		us2 := UserSong{
			UserId: uid,
			Title:  "song1",
			Genres: []UserGenre{},
			Tags:   []UserTag{tag1},
		}
		if err := us2.Create(DB); err != nil {
			t.Errorf("error at create %v", err)
		}
		if err := tag1.Delete(DB); err != nil {
			t.Errorf("error at delete tag %v", err)
		}
		//usersongからtagが消えてることを確認
		want := 0
		us := UserSong{}
		us.GetByID(DB, us1.ID, false)
		if l := len(us.Tags); l != want {
			t.Errorf("want =%d , but got =%d ", want, l)
		}
	})
}
