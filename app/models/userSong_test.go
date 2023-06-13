package models

import (
	"fmt"
	"testing"

	"example.com/app/utils"
)

func TestUserSong(t *testing.T) {
	t.Run("check prepareData", func(t *testing.T) {
		t.Skip()
		err := Init(true)
		if err != nil {
			t.Fatal(err)
		}
		defer ClearSQLiteDB()
		data := prepareData(t)

		for _, song := range data.Songs {
			utils.PrintStruct(song.Instruments)
			for _, section := range song.Sections {
				utils.PrintStruct(section.Instruments)
			}
		}
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
	err := Init(true)
	if err != nil {
		t.Fatal(err)
	}
	defer ClearSQLiteDB()
	data := prepareData(t)
	type Suite struct {
		memo string
		cond SongSearchCond
		want []UserSong
	}
	suites := []Suite{
		{
			memo: "return 2",
			cond: SongSearchCond{
				TagIds:      []uint{1},
				GenreIds:    []uint{1},
				SectionName: "",
			},
			want: []UserSong{data.Songs[0], data.Songs[1]},
		},
		{
			memo: "return 1",
			cond: SongSearchCond{
				TagIds:      []uint{2},
				GenreIds:    []uint{2},
				SectionName: "",
			},
			want: []UserSong{data.Songs[0]},
		},
		{
			memo: "empty condition",
			cond: SongSearchCond{
				TagIds:      []uint{},
				GenreIds:    []uint{},
				SectionName: "",
			},
			want: []UserSong{data.Songs[0], data.Songs[1]},
		},
		{
			memo: "sectionName",
			cond: SongSearchCond{
				TagIds:      []uint{},
				GenreIds:    []uint{},
				SectionName: "intro1",
			},
			want: []UserSong{data.Songs[0], data.Songs[1]},
		},
		{
			memo: "sectionName",
			cond: SongSearchCond{
				TagIds:      []uint{},
				GenreIds:    []uint{},
				SectionName: "intro2",
			},
			want: []UserSong{data.Songs[1]},
		},
	}
	for _, s := range suites {
		t.Run(s.memo, func(t *testing.T) {
			us := UserSong{}
			songs, err := us.GetByUserId(data.uid, s.cond)
			if err != nil {
				t.Error(err)
			}
			if len(songs) != len(s.want) {
				t.Errorf("length mismatch. want: %d, but got %d", len(s.want), len(songs))
			}
			for i, got := range songs {
				fmt.Println("====")
				utils.PrintStruct(got)
				if got.ID != s.want[i].ID {
					t.Errorf("want: %v, but got %v", s.want[i], got)
				}
			}
		})
	}

}
func TestGetSongByTagIds(t *testing.T) {
	err := Init(true)
	if err != nil {
		t.Fatal(err)
	}
	defer ClearSQLiteDB()
	data := prepareData(t)
	us := UserSong{}
	type Suite struct {
		memo   string
		tagIds []uint
		want   []UserSong
	}
	suites := []Suite{
		{
			memo:   "123",
			tagIds: []uint{1, 2, 3},
			want:   []UserSong{data.Songs[0], data.Songs[1]},
		},
		{memo: "2",
			tagIds: []uint{2},
			want:   []UserSong{data.Songs[0]},
		},
	}
	for _, s := range suites {
		t.Run(s.memo, func(t *testing.T) {
			songs, err := us.getSongByTagIds(data.uid, s.tagIds)
			if err != nil {
				t.Error(err)
			}
			if len(songs) != len(s.want) {
				t.Errorf("length mismatch. want: %d, but got %d", len(s.want), len(songs))
			}
			for i, got := range songs {
				if got.ID != s.want[i].ID {
					t.Errorf("want: %v, but got %v", s.want[i], got)
				}
			}
		})
	}
}

type TestData struct {
	uid    uint
	Tags   []UserTag
	Genres []UserGenre
	Songs  []UserSong
}

func prepareData(t *testing.T) TestData {
	var uid = uint(9999)
	var tag1 = UserTag{
		UserId:    uid,
		Name:      "tag1",
		SortOrder: 0,
		UserSongs: []UserSong{},
	}

	var tag2 = UserTag{
		UserId:    uid,
		Name:      "tag2",
		SortOrder: 0,
		UserSongs: []UserSong{},
	}

	var tag3 = UserTag{
		UserId:    uid,
		Name:      "tag3",
		SortOrder: 0,
		UserSongs: []UserSong{},
	}

	var genre1 = UserGenre{
		UserId:    uid,
		Name:      "genre1",
		SortOrder: 0,
		UserSongs: []UserSong{},
	}

	var genre2 = UserGenre{
		UserId:    uid,
		Name:      "genre2",
		SortOrder: 0,
		UserSongs: []UserSong{},
	}

	var genre3 = UserGenre{
		UserId:    uid,
		Name:      "genre3",
		SortOrder: 0,
		UserSongs: []UserSong{},
	}

	fmt.Println("preparing data")
	if err := tag1.Create(); err != nil {
		t.Errorf("error at create %v", err)
	}
	if err := tag2.Create(); err != nil {
		t.Errorf("error at create %v", err)
	}
	if err := tag3.Create(); err != nil {
		t.Errorf("error at create %v", err)
	}
	if err := genre1.Create(); err != nil {
		t.Errorf("error at create %v", err)
	}
	if err := genre2.Create(); err != nil {
		t.Errorf("error at create %v", err)
	}
	if err := genre3.Create(); err != nil {
		t.Errorf("error at create %v", err)
	}
	var us1 = UserSong{
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
		Instruments: []UserSongInstrument{
			{
				Name:      "guitar",
				SortOrder: 0,
			},
			{
				Name:      "piano",
				SortOrder: 1,
			},
			{
				Name:      "drums",
				SortOrder: 2,
			},
		},
		Sections: []UserSongSection{{
			Name:            "intro1",
			ProgressionsCSV: "Am7,F,G,C",
			Key:             1,
			BPM:             120,
			Scale:           "メジャー",
			Memo:            "sectionMemo1",
			LoopRange:       LoopRange{Start: 10, End: 20},
			Instruments: []UserSongInstrument{
				{
					Name: "guitar",
				},
				{
					Name: "drums",
				},
			},
		}, {
			Name:            "intro3",
			ProgressionsCSV: "Am7,F,G,C",
			Key:             1,
			BPM:             140,
			Scale:           "マイナー",
			Memo:            "sectionMemo2",
			LoopRange:       LoopRange{Start: 30, End: 40},
			Instruments: []UserSongInstrument{
				{
					Name: "piano",
				},
			},
		}},
	}
	var us2 = UserSong{
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
		Instruments: []UserSongInstrument{
			{
				Name:      "guitar2",
				SortOrder: 0,
			},
			{
				Name:      "piano2",
				SortOrder: 1,
			},
			{
				Name:      "drums2",
				SortOrder: 2,
			},
		},
		Sections: []UserSongSection{{
			Name:            "intro1",
			ProgressionsCSV: "Am7,F,G,C",
			Key:             1,
			BPM:             120,
			Scale:           "メジャー",
			Memo:            "sectionMemo1",
			LoopRange:       LoopRange{Start: 10, End: 20},
			Instruments: []UserSongInstrument{
				{
					Name: "piano2",
				},
			}}, {
			Name:            "intro2",
			ProgressionsCSV: "Am7,F,G,C",
			Key:             1,
			BPM:             140,
			Scale:           "マイナー",
			Memo:            "sectionMemo2",
			LoopRange:       LoopRange{Start: 30, End: 40},
			Instruments: []UserSongInstrument{
				{
					Name: "piano2",
				},
				{
					Name: "drums2",
				},
			}},
		}}
	if err := us1.Create(); err != nil {
		t.Errorf("error at create %v", err)
	}
	if err := us2.Create(); err != nil {
		t.Errorf("error at create %v", err)
	}
	return TestData{
		uid:    uid,
		Tags:   []UserTag{tag1, tag2},
		Genres: []UserGenre{genre1, genre2},
		Songs:  []UserSong{us1, us2},
	}
}
