package models

import (
	"errors"
	"fmt"
	"testing"

	"example.com/app/utils"
	"gorm.io/gorm"
)

func Test_PrepareData(t *testing.T) {
	t.Run("check prepareTestData", func(t *testing.T) {
		defer ClearTestDB(DB)
		data := prepareTestData(t)

		fmt.Println("@@check")
		for _, song := range data.Songs {
			fmt.Printf("id = %d, uuid = %s\n", song.ID, song.UUID)
			utils.PrintStruct(song.Instruments)
			for _, section := range song.Sections {
				utils.PrintStruct(section.Instruments)
			}
		}
	})
}

func TestUserSong(t *testing.T) {

	t.Run("delete tag from UserSong", func(t *testing.T) {
		t.Skip()
		want := 1
		err := InitTestDB()
		if err != nil {
			t.Fatal(err)
		}
		defer ClearTestDB(DB)

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
		song.GetByID(DB, us.ID, false)
		//tagのリレーション削除
		song.DeleteTagRelation(DB, &song.Tags[1])
		//tagを一つ削除
		song.Tags = append(song.Tags[:1])
		song.Update(DB)

		song2 := UserSong{}
		song2.GetByID(DB, song.ID, false)
		if l := len(song2.Tags); l != want {
			t.Errorf("want =%d , but got =%d ", want, l)
		}
	})
	t.Run("append tag to UserSong", func(t *testing.T) {
		want := 2
		err := InitTestDB()
		if err != nil {
			t.Fatal(err)
		}
		defer ClearTestDB(DB)

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
		song.GetByID(DB, us.ID, false)
		//tagを一つ追加
		song.Tags = append(song.Tags, tag2)
		song.Update(DB)

		song2 := UserSong{}
		song2.GetByID(DB, song.ID, false)
		if l := len(song2.Tags); l != want {
			t.Errorf("want =%d , but got =%d ", want, l)
		}
	})

}

// transaction, lockの挙動確認
func TestTransaction(t *testing.T) {
	err := InitTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer ClearTestDB(DB)
	us := UserSong{}
	us.Title = "BEFORE"
	//データ作成
	if err := us.Create(DB); err != nil {
		t.Fatal(err)
	}
	t.Run("If an error happened in transaction, it should rollback", func(t *testing.T) {
		us := UserSong{}
		us.Title = "BEFORE"
		//データ作成
		if err := us.Create(DB); err != nil {
			t.Fatal(err)
		}
		err := DB.Transaction(func(tx *gorm.DB) error {
			song := UserSong{}
			res := tx.Model(UserSong{}).First(&song)
			if res.RowsAffected == 0 {
				return errors.New("test data not found")
			}
			song.Title = "AFTER"
			if err := song.Update(tx); err != nil {
				return err
			}
			//最後にエラーを返してtransactionを失敗させる
			return errors.New("intentional")
		})
		if err.Error() != "intentional" { //意図してないエラーなのでfail
			t.Fatal(err)
		}
		//再取得してtitleを確認
		after := UserSong{}
		DB.Model(UserSong{}).First(&after)
		fmt.Println(after.Title)
		if after.Title == "AFTER" { //更新されてるのでfail
			t.Fatal(errors.New("transaction not working"))
		}
	})
	t.Run("lock success", func(t *testing.T) {})

}
func TestSearch(t *testing.T) {
	err := InitTestDB()
	if err != nil {
		t.Fatal(err)
	}
	defer ClearTestDB(DB)
	data := prepareTestData(t)
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
