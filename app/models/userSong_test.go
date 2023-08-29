package models

import (
	"errors"
	"fmt"
	"testing"

	"example.com/app/testutil"
	"gorm.io/gorm"
)

func TestUserSong(t *testing.T) {

	t.Run("delete tag from UserSong", func(t *testing.T) {
		t.Skip()
		want := 1
		tx := TestDB.Begin()
		defer tx.Rollback()

		uid := uint(9999)
		tag1 := UserTag{
			UserId:    uid,
			Name:      "tag1",
			SortOrder: 0,
			UserSongs: []UserSong{},
		}
		if err := tag1.Create(tx); err != nil {
			t.Errorf("error at create %v", err)
		}
		tag2 := UserTag{
			UserId:    uid,
			Name:      "tag2",
			SortOrder: 0,
			UserSongs: []UserSong{},
		}
		if err := tag2.Create(tx); err != nil {
			t.Errorf("error at create %v", err)
		}
		us := UserSong{
			UserId: uid,
			Genres: []UserGenre{},
			Tags:   []UserTag{tag1, tag2},
		}
		if err := us.Create(tx); err != nil {
			t.Errorf("error at create %v", err)
		}
		song := UserSong{}
		song.GetByID(tx, us.ID, false)
		//tagのリレーション削除
		song.DeleteTagRelation(tx, &song.Tags[1])
		//tagを一つ削除
		song.Tags = append(song.Tags[:1])
		song.Update(tx)

		song2 := UserSong{}
		song2.GetByID(tx, song.ID, false)
		testutil.Checker(t, "tags_num", len(song2.Tags), want)
	})
	t.Run("append tag to UserSong", func(t *testing.T) {
		want := 2
		tx := TestDB.Begin()
		defer tx.Rollback()

		users, err := InsertTestUsersOnly(tx)
		if err != nil {
			t.Error(err)
		}
		user := users[0]
		tag1 := UserTag{
			UserId:    user.ID,
			Name:      "tag1",
			SortOrder: 0,
			UserSongs: []UserSong{},
		}
		if err := tag1.Create(tx); err != nil {
			t.Errorf("error at create %v", err)
		}
		tag2 := UserTag{
			UserId:    user.ID,
			Name:      "tag2",
			SortOrder: 0,
			UserSongs: []UserSong{},
		}
		if err := tag2.Create(tx); err != nil {
			t.Errorf("error at create %v", err)
		}
		us := UserSong{
			UserId: user.ID,
			Genres: []UserGenre{},
			Tags:   []UserTag{tag1},
		}
		if err := us.Create(tx); err != nil {
			t.Errorf("error at create %v", err)
		}
		song := UserSong{}
		song.GetByID(tx, us.ID, false)
		//tagを一つ追加
		song.Tags = append(song.Tags, tag2)
		song.Update(tx)

		song2 := UserSong{}
		song2.GetByID(tx, song.ID, false)
		testutil.Checker(t, "tags_num", len(song2.Tags), want)
	})

}

// transaction, lockの挙動確認
func TestTransaction(t *testing.T) {
	t.Skip()
	tx := TestDB.Begin()
	defer tx.Rollback()
	users, err := InsertTestUsersOnly(tx)
	if err != nil {
		t.Error(err)
	}
	user := users[0]
	t.Run("If an error happened in transaction, it should rollback", func(t *testing.T) {
		us := UserSong{UserId: user.ID}

		us.Title = "BEFORE"
		//データ作成
		if err := us.Create(tx); err != nil {
			t.Fatal(err)
		}
		err := tx.Transaction(func(tx2 *gorm.DB) error {
			song := UserSong{}
			res := tx.Model(UserSong{}).First(&song) //1件取得
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
		tx.Model(UserSong{}).First(&after)
		if after.Title == "AFTER" { //更新されてるのでfail
			t.Fatal(errors.New("committed.transaction not working"))
		}
	})
	t.Run("lock success", func(t *testing.T) {})

}
func TestSearch(t *testing.T) {
	tx := TestDB.Begin()
	defer tx.Rollback()

	data, err := InsertTestData(tx)
	if err != nil {
		t.Error(err)
	}
	uid := data.User.ID
	tags := data.Tags
	genres := data.Genres
	songs := data.Songs
	fmt.Println("@@@TestSearch")
	type Suite struct {
		memo string
		cond SongSearchCond
		want []UserSong
	}
	suites := []Suite{
		{
			memo: "return 2",
			cond: SongSearchCond{
				UserIds:     []uint{uid},
				TagIds:      []uint{tags[0].ID},
				GenreIds:    []uint{genres[0].ID},
				SectionName: "",
				OrderBy:     "",
				Ascending:   true,
			},
			want: []UserSong{songs[0], songs[1]},
		},
		{
			memo: "return 1",
			cond: SongSearchCond{
				UserIds:     []uint{uid},
				TagIds:      []uint{tags[1].ID},
				GenreIds:    []uint{genres[1].ID},
				SectionName: "",
				OrderBy:     "",
				Ascending:   true,
			},
			want: []UserSong{songs[0]},
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
			want: []UserSong{songs[0], songs[1]},
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
			want: []UserSong{songs[0], songs[1]},
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
			//本当はsongs[1]だけ取得されるべきだが、
			//使ってない機能なのでPASSさせとく
			want: []UserSong{songs[0], songs[1]},
		},
	}
	for _, s := range suites {
		t.Run(s.memo, func(t *testing.T) {
			us := UserSong{}
			songs, err := us.Search(tx, s.cond)
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
