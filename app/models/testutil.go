package models

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var connectedToTestDB = false
var TestDB *gorm.DB

func InitTestDB() (*gorm.DB, error) {
	//docker-composeのtest_db参照
	user := "testuser"
	password := "testpassword"
	db_name := "test_db"
	db_host := "test_db:3307"
	dsn := user + ":" + password + "@tcp(" + db_host + ")/" + db_name + "?charset=utf8mb4&parseTime=True"
	db, err := connectMySQL(dsn)
	if err != nil {
		return nil, err
	}
	connectedToTestDB = true
	migrateModels(db)
	fmt.Println("@@@@connected to test db")
	return db, nil
}

func ClearTestDB(db *gorm.DB) {
	//テスト用DBのクリア
	//テーブル消してマイグレーションし直す
	//怖すぎ
	if !connectedToTestDB {
		return
	}
	if os.Getenv("ENV") == "prod" {
		return
	}
	fmt.Println("clearing DB")
	db.Exec("DROP TABLE usersongs_genres")
	db.Exec("DROP TABLE usersongs_tags")
	db.Exec("DROP TABLE user_tags")
	db.Exec("DROP TABLE user_genres")
	db.Exec("DROP TABLE user_audio_ranges")
	db.Exec("DROP TABLE user_section_midis")
	db.Exec("DROP TABLE sections_instruments")
	db.Exec("DROP TABLE user_song_instruments")
	db.Exec("DROP TABLE user_song_audios")
	db.Exec("DROP TABLE sessions")
	db.Exec("DROP TABLE user_song_sections")
	db.Exec("DROP TABLE user_songs")
	db.Exec("DROP TABLE users")
	migrateModels(db)
}

type TestData struct {
	uid    uint
	Tags   []UserTag
	Genres []UserGenre
	Songs  []UserSong
}

/*
return dummy users

user[0]: confirmed user with email "test@test.test"

user[1]: unconfirmed user with email "test2@test.test"
*/
func PrepareTestUsersOnly(db *gorm.DB) ([]*User, error) {

	user := User{ID: uint(9999), Email: "test@test.test", Password: "dummy", IsConfirmed: true}
	user2 := User{ID: uint(9998), Email: "test2@test.test", Password: "dummy", IsConfirmed: false}

	if result := db.Create(&user); result.Error != nil {
		return nil, result.Error
	}
	if result := db.Create(&user2); result.Error != nil {
		return nil, result.Error
	}
	return []*User{&user, &user2}, nil
}

func PrepareTestData(t *testing.T, db *gorm.DB) TestData {
	var uid = uint(9999)
	var user = User{
		ID:          uid,
		Email:       "test@test.test",
		IsConfirmed: true,
	}
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
	if res := db.Create(&user); res.Error != nil {
		t.Errorf("error at create %v ", res.Error)
	}
	if err := tag1.Create(db); err != nil {
		t.Errorf("error at create %v", err)
	}
	if err := tag2.Create(db); err != nil {
		t.Errorf("error at create %v", err)
	}
	if err := tag3.Create(db); err != nil {
		t.Errorf("error at create %v", err)
	}
	if err := genre1.Create(db); err != nil {
		t.Errorf("error at create %v", err)
	}
	if err := genre2.Create(db); err != nil {
		t.Errorf("error at create %v", err)
	}
	if err := genre3.Create(db); err != nil {
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
	if err := us1.Create(db); err != nil {
		t.Errorf("error at create %v", err)
	}
	if err := us2.Create(db); err != nil {
		t.Errorf("error at create %v", err)
	}
	return TestData{
		uid:    uid,
		Tags:   []UserTag{tag1, tag2},
		Genres: []UserGenre{genre1, genre2},
		Songs:  []UserSong{us1, us2},
	}
}
