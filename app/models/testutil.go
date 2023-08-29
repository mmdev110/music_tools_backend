package models

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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
	migrateModels(db)
	fmt.Println("@@@@connected to test db")
	return db, nil
}

type TestData struct {
	User   *User
	Tags   []UserTag
	Genres []UserGenre
	Songs  []UserSong
}

/*
return dummy data(no db inserts, data only)

テスト毎にトランザクション張ってからinsertしたいが、テストケース作成のためデータだけ先に必要というジレンマあり
なのでテストデータ取得とインサート処理を分けている
*/
func GetTestData() TestData {
	var uid = uint(9999)
	var user = User{
		ID:          uid,
		Email:       "test@test.test",
		IsConfirmed: true,
	}
	var tag1 = UserTag{
		ID:        1,
		UserId:    uid,
		Name:      "tag1",
		SortOrder: 0,
		UserSongs: []UserSong{},
	}

	var tag2 = UserTag{
		ID:        2,
		UserId:    uid,
		Name:      "tag2",
		SortOrder: 0,
		UserSongs: []UserSong{},
	}

	var tag3 = UserTag{
		ID:        3,
		UserId:    uid,
		Name:      "tag3",
		SortOrder: 0,
		UserSongs: []UserSong{},
	}

	var genre1 = UserGenre{
		ID:        1,
		UserId:    uid,
		Name:      "genre1",
		SortOrder: 0,
		UserSongs: []UserSong{},
	}

	var genre2 = UserGenre{
		ID:        2,
		UserId:    uid,
		Name:      "genre2",
		SortOrder: 0,
		UserSongs: []UserSong{},
	}

	var genre3 = UserGenre{
		ID:        3,
		UserId:    uid,
		Name:      "genre3",
		SortOrder: 0,
		UserSongs: []UserSong{},
	}
	var genre4 = UserGenre{
		ID:        4,
		UserId:    uid,
		Name:      "genre4",
		SortOrder: 0,
		UserSongs: []UserSong{},
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
	return TestData{
		User:   &user,
		Tags:   []UserTag{tag1, tag2, tag3},
		Genres: []UserGenre{genre1, genre2, genre3, genre4},
		Songs:  []UserSong{us1, us2},
	}
}

/*
APIテスト用の様々なsong
*/
type VariousSongs struct {
	Complex    UserSong
	Simple     UserSong
	NoInstUsed UserSong
}

func VariousSongData() VariousSongs {
	data := GetTestData()

	complex := UserSong{ //relationフルに使うsong
		Title:  "complex_title",
		Artist: "complex_artist",
		Tags:   []UserTag{data.Tags[0], data.Tags[2]},
		Genres: []UserGenre{data.Genres[0], data.Genres[2]},
		Instruments: []UserSongInstrument{
			{Name: "drum1"},
			{Name: "drum2"},
			{Name: "melody1"},
		},
		Audio: UserSongAudio{
			Name: "my audio1",
			Url:  Url{Get: "", Put: ""},
		},
		Sections: []UserSongSection{
			{
				Name: "intro",
				Instruments: []UserSongInstrument{
					{Name: "drum1"},
					{Name: "drum2"},
				},
				AudioRanges: []UserAudioRange{
					{Name: "full", Start: 0, End: 100},
					{Name: "transition", Start: 95, End: 100},
				},
				Midi: UserSectionMidi{
					Name: "piano_midi",
				},
			},
			{
				Name: "A melo",
				Instruments: []UserSongInstrument{
					{Name: "drum1"},
					{Name: "drum2"},
					{Name: "melody1"},
				},
				AudioRanges: []UserAudioRange{
					{Name: "full", Start: 100, End: 155},
					{Name: "transition", Start: 135, End: 155},
				},
				Midi: UserSectionMidi{
					Name: "bass_midi",
				},
			},
		},
	}
	noInstUsed := UserSong{ //instはあるが、sectionでinst指定しないパターン
		Title:  "noinstused_title",
		Artist: "noinstused_artist",
		Instruments: []UserSongInstrument{
			{Name: "drum1"},
			{Name: "drum2"},
			{Name: "melody1"},
		},
		Audio: UserSongAudio{
			Name: "my audio1",
			Url:  Url{Get: "", Put: ""},
		},
		Sections: []UserSongSection{
			{
				Name: "intro",
			},
			{
				Name: "A melo",
			},
		},
	}
	simple := UserSong{ //最小構成のsong
		Title:  "simple_title",
		Artist: "simple_artist",
	}

	return VariousSongs{
		Complex:    complex,
		Simple:     simple,
		NoInstUsed: noInstUsed,
	}
}

/*
return dummy users(no db inserts, data only)

user[0]: confirmed user with ID 10000, email "test@test.test"

user[1]: unconfirmed user with UD 10001, email "test2@test.test"
*/
func GetTestUsers() []*User {
	user := User{ID: uint(10000), Email: "tes@test.test", Password: "dummy", IsConfirmed: true}
	user2 := User{ID: uint(10001), Email: "test2@test.test", Password: "dummy", IsConfirmed: false}
	return []*User{&user, &user2}
}

/*
insert dummy users to db
*/
func InsertTestUsersOnly(db *gorm.DB) ([]*User, error) {
	users := GetTestUsers()
	for _, user := range users {
		if result := db.Create(user); result.Error != nil {
			return nil, result.Error
		}
	}

	return users, nil
}

/*
insert dummy data to DB
*/
func InsertTestData(db *gorm.DB) (*TestData, error) {
	data := GetTestData()
	if result := db.Create(data.User); result.Error != nil {
		return nil, result.Error
	}
	for _, tag := range data.Tags {
		if err := tag.Create(db); err != nil {
			return nil, fmt.Errorf("error at create %v", err)
		}
	}
	for _, genre := range data.Genres {
		if err := genre.Create(db); err != nil {
			return nil, fmt.Errorf("error at create %v", err)
		}
	}
	if err := data.Songs[0].Create(db); err != nil {
		return nil, fmt.Errorf("error at create %v", err)
	}
	if err := data.Songs[1].Create(db); err != nil {
		return nil, fmt.Errorf("error at create %v", err)
	}

	return &data, nil
}
