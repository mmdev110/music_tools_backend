package models

import (
	"fmt"
	"testing"

	"example.com/app/utils"
	"github.com/google/uuid"
)

func TestUserSong(t *testing.T) {
	t.Run("check prepareData", func(t *testing.T) {
		err := Init(true)
		if err != nil {
			t.Fatal(err)
		}
		defer ClearSQLiteDB()
		data := prepareData(t)

		fmt.Println("@@check")
		for _, song := range data.Songs {
			fmt.Printf("id = %d, uuid = %s\n", song.ID, song.UUID)
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
		if err := tag1.Create(DB); err != nil {
			t.Errorf("error at create %v", err)
		}
		tag2 := UserTag{
			UserId:    uid,
			Name:      "tag2",
			SortOrder: 0,
			UserSongs: []UserSong{},
		}
		if err := tag2.Create(DB); err != nil {
			t.Errorf("error at create %v", err)
		}
		us := UserSong{
			UserId: uid,
			Genres: []UserGenre{},
			Tags:   []UserTag{tag1, tag2},
		}
		if err := us.Create(DB); err != nil {
			t.Errorf("error at create %v", err)
		}
		song := UserSong{}
		song.GetByID(nil, us.ID, false)
		//tagのリレーション削除
		song.DeleteTagRelation(nil, &song.Tags[1])
		//tagを一つ削除
		song.Tags = append(song.Tags[:1])
		song.Update(nil)

		song2 := UserSong{}
		song2.GetByID(nil, song.ID, false)
		if l := len(song2.Tags); l != want {
			t.Errorf("want =%d , but got =%d ", want, l)
		}
	})
	t.Run("append tag to UserSong", func(t *testing.T) {
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
		if err := tag1.Create(DB); err != nil {
			t.Errorf("error at create %v", err)
		}
		tag2 := UserTag{
			UserId:    uid,
			Name:      "tag2",
			SortOrder: 0,
			UserSongs: []UserSong{},
		}
		if err := tag2.Create(DB); err != nil {
			t.Errorf("error at create %v", err)
		}
		us := UserSong{
			UserId: uid,
			Genres: []UserGenre{},
			Tags:   []UserTag{tag1},
		}
		if err := us.Create(DB); err != nil {
			t.Errorf("error at create %v", err)
		}
		song := UserSong{}
		song.GetByID(nil, us.ID, false)
		//tagを一つ追加
		song.Tags = append(song.Tags, tag2)
		song.Update(nil)

		song2 := UserSong{}
		song2.GetByID(nil, song.ID, false)
		if l := len(song2.Tags); l != want {
			t.Errorf("want =%d , but got =%d ", want, l)
		}
	})

}
func TestSearch(t *testing.T) {
	err := Init(true)
	if err != nil {
		t.Fatal(err)
	}
	defer ClearSQLiteDB()
	data := prepareData(t)
	fmt.Println("@@@TestSearch")
	type Suite struct {
		memo string
		cond SongSearchCond
		want []UserSong
	}
	uid := uint(9999)
	suites := []Suite{
		{
			memo: "return 2",
			cond: SongSearchCond{
				UserIds:     []uint{uid},
				TagIds:      []uint{1},
				GenreIds:    []uint{1},
				SectionName: "",
				OrderBy:     "",
				Ascending:   true,
			},
			want: []UserSong{data.Songs[0], data.Songs[1]},
		},
		{
			memo: "return 1",
			cond: SongSearchCond{
				UserIds:     []uint{uid},
				TagIds:      []uint{2},
				GenreIds:    []uint{2},
				SectionName: "",
				OrderBy:     "",
				Ascending:   true,
			},
			want: []UserSong{data.Songs[0]},
		},
		{
			memo: "empty condition",
			cond: SongSearchCond{
				UserIds:     []uint{uid},
				TagIds:      []uint{},
				GenreIds:    []uint{},
				SectionName: "",
				OrderBy:     "",
				Ascending:   true,
			},
			want: []UserSong{data.Songs[0], data.Songs[1]},
		},
		{
			memo: "sectionName",
			cond: SongSearchCond{
				UserIds:     []uint{uid},
				TagIds:      []uint{},
				GenreIds:    []uint{},
				SectionName: "intro1",
				OrderBy:     "",
				Ascending:   true,
			},
			want: []UserSong{data.Songs[0], data.Songs[1]},
		},
		{
			memo: "sectionName2",
			cond: SongSearchCond{
				UserIds:     []uint{uid},
				TagIds:      []uint{},
				GenreIds:    []uint{},
				SectionName: "intro2",
				OrderBy:     "",
				Ascending:   true,
			},
			want: []UserSong{data.Songs[1]},
		},
	}
	for _, s := range suites {
		t.Run(s.memo, func(t *testing.T) {
			us := UserSong{}
			songs, err := us.Search(DB, s.cond)
			if err != nil {
				t.Fatal(err)
			}
			if len(songs) != len(s.want) {
				t.Fatalf("length mismatch. want: %d, but got %d", len(s.want), len(songs))
			}
			for i, got := range songs {
				fmt.Println("====")
				//utils.PrintStruct(got)
				if got.ID != s.want[i].ID {
					t.Fatalf("want: %v, but got %v", s.want[i], got)
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
	if err := tag1.Create(DB); err != nil {
		t.Errorf("error at create %v", err)
	}
	if err := tag2.Create(DB); err != nil {
		t.Errorf("error at create %v", err)
	}
	if err := tag3.Create(DB); err != nil {
		t.Errorf("error at create %v", err)
	}
	if err := genre1.Create(DB); err != nil {
		t.Errorf("error at create %v", err)
	}
	if err := genre2.Create(DB); err != nil {
		t.Errorf("error at create %v", err)
	}
	if err := genre3.Create(DB); err != nil {
		t.Errorf("error at create %v", err)
	}
	var us1 = UserSong{
		UserId: uid,
		UUID:   uuid.NewString(),
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
			AudioRanges:     []UserAudioRange{{Name: "full", Start: 10, End: 20}},
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
			AudioRanges:     []UserAudioRange{{Name: "full", Start: 10, End: 20}},
			Instruments: []UserSongInstrument{
				{
					Name: "piano",
				},
			},
		}},
	}
	var us2 = UserSong{
		UserId: uid,
		UUID:   uuid.NewString(),
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
			AudioRanges:     []UserAudioRange{{Name: "full", Start: 10, End: 20}},
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
			AudioRanges:     []UserAudioRange{{Name: "full", Start: 10, End: 20}},
			Instruments: []UserSongInstrument{
				{
					Name: "piano2",
				},
				{
					Name: "drums2",
				},
			}},
		}}
	if err := us1.Create(DB); err != nil {
		t.Errorf("error at create %v", err)
	}
	if err := us2.Create(DB); err != nil {
		t.Errorf("error at create %v", err)
	}
	return TestData{
		uid:    uid,
		Tags:   []UserTag{tag1, tag2},
		Genres: []UserGenre{genre1, genre2},
		Songs:  []UserSong{us1, us2},
	}
}
