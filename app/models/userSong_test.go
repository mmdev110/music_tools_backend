package models

import (
	"fmt"
	"testing"

	"example.com/app/utils"
)

func TestUserSong(t *testing.T) {

	t.Run("save new UserSong with full associations", func(t *testing.T) {
		t.Skip()
		err := Init(true)
		if err != nil {
			t.Fatal(err)
		}
		defer ClearSQLiteDB()

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
				Key:             1,
				BPM:             120,
				Scale:           "メジャー",
				Memo:            "sectionMemo1",
				LoopRange:       LoopRange{Start: 10, End: 20},
			}, {
				Name:            "intro2",
				ProgressionsCSV: "Am7,F,G,C",
				Key:             1,
				BPM:             140,
				Scale:           "マイナー",
				Memo:            "sectionMemo2",
				LoopRange:       LoopRange{Start: 30, End: 40},
			}},
		}
		if err := us.Create(); err != nil {
			t.Errorf("error at create %v", err)
		}
		utils.PrintStruct(us.Genres)

		fmt.Println("====update genre")
		song := UserSong{}
		song.GetByID(us.ID)
		song.Genres[0].Name = "NewGenre"
		song.Update()
		fmt.Println("====check updated genre")
		song2 := UserSong{}
		song2.GetByID(song.ID)
		utils.PrintStruct(song2.Genres)
	})
	t.Run("delete tag from UserSong", func(t *testing.T) {
		t.Skip()
		want := 1
		err := Init(true)
		if err != nil {
			t.Fatal(err)
		}
		defer ClearSQLiteDB()

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
		us := UserSong{
			UserId: uid,
			Genres: []UserGenre{},
			Tags:   []UserTag{tag1, tag2},
		}
		if err := us.Create(); err != nil {
			t.Errorf("error at create %v", err)
		}
		song := UserSong{}
		song.GetByID(us.ID)
		//tagのリレーション削除
		song.DeleteTagRelation(&song.Tags[1])
		//tagを一つ削除
		song.Tags = append(song.Tags[:1])
		song.Update()

		song2 := UserSong{}
		song2.GetByID(song.ID)
		if l := len(song2.Tags); l != want {
			t.Errorf("want =%d , but got =%d ", want, l)
		}
	})
	t.Run("append tag to UserSong", func(t *testing.T) {
		t.Skip()
		want := 2
		err := Init(true)
		if err != nil {
			t.Fatal(err)
		}
		defer ClearSQLiteDB()

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
		us := UserSong{
			UserId: uid,
			Genres: []UserGenre{},
			Tags:   []UserTag{tag1},
		}
		if err := us.Create(); err != nil {
			t.Errorf("error at create %v", err)
		}
		song := UserSong{}
		song.GetByID(us.ID)
		//tagを一つ追加
		song.Tags = append(song.Tags, tag2)
		song.Update()

		song2 := UserSong{}
		song2.GetByID(song.ID)
		if l := len(song2.Tags); l != want {
			t.Errorf("want =%d , but got =%d ", want, l)
		}
	})

}

func TestGetByUserId(t *testing.T) {
	t.Run("get userSong with search conditions", func(t *testing.T) {
		err := Init(true)
		if err != nil {
			t.Fatal(err)
		}
		defer ClearSQLiteDB()
		prepareData(t)
		fmt.Println("====get usersong with condition")
		song := UserSong{}
		cond := &SongSearchCond{
			TagIds:      []uint{1, 2, 3},
			GenreIds:    []uint{2, 3},
			SectionName: "",
		}
		//cond := &SongSearchCond{}
		songs, _ := song.GetByUserId(9999, cond)
		for _, v := range songs {
			utils.PrintStruct(v)
		}
	})
	t.Run("getSongByTagIds", func(t *testing.T) {
		err := Init(true)
		if err != nil {
			t.Fatal(err)
		}
		defer ClearSQLiteDB()
		prepareData(t)
		fmt.Println("====get usersong with condition")
		us := UserSong{}
		tagIds := []uint{1, 2, 3}
		//cond := &SongSearchCond{}
		songs, _ := us.getSongByTagIds(9999, tagIds)
		for _, v := range songs {
			utils.PrintStruct(v.ID)
			utils.PrintStruct(v.Tags)
		}
	})
}

func prepareData(t *testing.T) {
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
	tag3 := UserTag{
		UserId:    uid,
		Name:      "tag3",
		SortOrder: 0,
		UserSongs: []UserSong{},
	}
	if err := tag3.Create(); err != nil {
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
	genre3 := UserGenre{
		UserId:    uid,
		Name:      "genre3",
		SortOrder: 0,
		UserSongs: []UserSong{},
	}
	if err := genre3.Create(); err != nil {
		t.Errorf("error at create %v", err)
	}
	us1 := UserSong{
		UserId: uid,
		Title:  "title1",
		Artist: "artist1",
		Memo:   "memo1",
		Genres: []UserGenre{genre1, genre2, genre3},
		Tags:   []UserTag{tag1, tag2, tag3},
		Audio: UserSongAudio{
			Name: "song1",
			Url:  Url{Get: "get1", Put: "put1"},
		},
		Sections: []UserSongSection{{
			Name:            "intro1",
			ProgressionsCSV: "Am7,F,G,C",
			Key:             1,
			BPM:             120,
			Scale:           "メジャー",
			Memo:            "sectionMemo1",
			LoopRange:       LoopRange{Start: 10, End: 20},
		}, {
			Name:            "intro3",
			ProgressionsCSV: "Am7,F,G,C",
			Key:             1,
			BPM:             140,
			Scale:           "マイナー",
			Memo:            "sectionMemo2",
			LoopRange:       LoopRange{Start: 30, End: 40},
		}},
	}
	if err := us1.Create(); err != nil {
		t.Errorf("error at create %v", err)
	}
	us2 := UserSong{
		UserId: uid,
		Title:  "title1",
		Artist: "artist1",
		Memo:   "memo1",
		Genres: []UserGenre{genre1},
		Tags:   []UserTag{tag1},
		Audio: UserSongAudio{
			Name: "song1",
			Url:  Url{Get: "get1", Put: "put1"},
		},
		Sections: []UserSongSection{{
			Name:            "intro1",
			ProgressionsCSV: "Am7,F,G,C",
			Key:             1,
			BPM:             120,
			Scale:           "メジャー",
			Memo:            "sectionMemo1",
			LoopRange:       LoopRange{Start: 10, End: 20},
		}, {
			Name:            "intro2",
			ProgressionsCSV: "Am7,F,G,C",
			Key:             1,
			BPM:             140,
			Scale:           "マイナー",
			Memo:            "sectionMemo2",
			LoopRange:       LoopRange{Start: 30, End: 40},
		}},
	}
	if err := us2.Create(); err != nil {
		t.Errorf("error at create %v", err)
	}

}
