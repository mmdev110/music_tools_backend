package models

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"example.com/app/conf"
	"example.com/app/utils"
	"github.com/google/uuid"
	_ "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// song情報
type UserSong struct {
	ID       uint              `gorm:"primarykey" json:"id"`
	UUID     string            `gorm:"index:idx_uuid,unique;not null" json:"uuid"` //ユーザー表示用のid
	UserId   uint              `gorm:"not null" json:"user_id"`
	Title    string            `json:"title"`
	Artist   string            `json:"artist"`
	Sections []UserSongSection `json:"sections"`
	Memo     string            `json:"memo"`
	//オーディオファイル
	Audio UserSongAudio `json:"audio"`
	//ジャンル
	Genres []UserGenre `gorm:"many2many:usersongs_genres" json:"genres"`
	//タグ
	Tags           []UserTag            `gorm:"many2many:usersongs_tags" json:"tags"`
	Instruments    []UserSongInstrument `json:"instruments"`
	ViewTimes      uint                 `gorm:"not null" json:"view_times"`
	LastModifiedAt time.Time            `json:"-"`
	//current_timestamp(3)について
	//https://github.com/go-gorm/mysql/issues/58
	LastViewedAt time.Time      `json:"-"`
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

func (us *UserSong) Create(db *gorm.DB) error {
	//Sections.Instrumentsがsongs.Instrumentsに依存しており、
	//同時に生成できないため、CREATE文を分ける
	//CREATE UserSong
	sections := us.Sections
	us.Sections = []UserSongSection{} //一旦Sectionsを空にする
	//uuid付与
	us.UUID = uuid.NewString()
	us.LastModifiedAt = time.Now()
	us.LastViewedAt = time.Now()
	for {
		result := db.Omit("Tags.*", "Genres.*").Create(&us)
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			//uuidが衝突してるので更新して再実行
			us.UUID = uuid.NewString()
			continue
		}
		if result.Error != nil { //その他のエラー
			return result.Error
		}
		break
	}

	//CREATE Sections
	//Instrumentsに付与されたIDをsectionsに紐付けてCREATE
	us.Sections = sections
	instruments := us.Instruments
	//for地獄
	for i, sec := range us.Sections {
		sec.UserSongId = us.ID
		for j, instSec := range sec.Instruments {
			for _, inst := range instruments {
				if inst.Name == instSec.Name {
					instSec.ID = inst.ID
					instSec.UserSongId = inst.UserSongId
				}
			}
			sec.Instruments[j] = instSec
		}
		sec.Create(db)
		us.Sections[i] = sec
	}
	return nil
}

// songを返す
func (us *UserSong) GetByID(db *gorm.DB, id uint, lock bool) *gorm.DB {
	if lock {
		fmt.Println("lock!")
		db = db.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	result := db.Model(&UserSong{}).
		Preload("Audio").
		Preload("Instruments", func(db *gorm.DB) *gorm.DB {
			return db.Order("user_song_instruments.sort_order ASC")
		}).
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("user_song_sections.sort_order ASC")
		}).
		Preload("Sections.AudioRanges", func(db *gorm.DB) *gorm.DB {
			return db.Order("user_audio_ranges.sort_order ASC")
		}).
		Preload("Sections.Midi").
		Preload("Sections.Instruments", func(db *gorm.DB) *gorm.DB {
			return db.Order("user_song_instruments.sort_order ASC")
		}).
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Order("user_tags.sort_order ASC")
		}).
		Preload("Genres", func(db *gorm.DB) *gorm.DB {
			return db.Order("user_genres.sort_order ASC")
		}).
		First(&us, id)
	if result.RowsAffected == 0 {
		return result
	}
	return result
}
func (us *UserSong) GetByUUID(db *gorm.DB, uuid string, lock bool) (err error, isFound bool) {
	if lock {
		fmt.Println("lock!")
		db = db.Clauses(clause.Locking{Strength: "UPDATE"})
	}
	result := db.Model(&UserSong{}).
		Preload("Audio").
		Preload("Instruments", func(db *gorm.DB) *gorm.DB {
			return db.Order("user_song_instruments.sort_order ASC")
		}).
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("user_song_sections.sort_order ASC")
		}).
		Preload("Sections.AudioRanges", func(db *gorm.DB) *gorm.DB {
			return db.Order("user_audio_ranges.sort_order ASC")
		}).
		Preload("Sections.Midi").
		Preload("Sections.Instruments", func(db *gorm.DB) *gorm.DB {
			return db.Order("user_song_instruments.sort_order ASC")
		}).
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Order("user_tags.sort_order ASC")
		}).
		Preload("Genres", func(db *gorm.DB) *gorm.DB {
			return db.Order("user_genres.sort_order ASC")
		}).
		Where("uuid = ?", uuid).
		First(&us)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, false
		}
		return result.Error, false
	}
	return nil, true
}

// songを返す
func (us *UserSong) GetByUserId(db *gorm.DB, userId uint) ([]UserSong, error) {
	var songs []UserSong
	result := db.Model(&UserSong{}).
		Preload("Audio").
		Preload("Instruments", func(db *gorm.DB) *gorm.DB {
			return db.Order("user_song_instruments.sort_order ASC")
		}).
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("user_song_sections.sort_order ASC")
		}).
		Preload("Sections.AudioRanges", func(db *gorm.DB) *gorm.DB {
			return db.Order("user_audio_ranges.sort_order ASC")
		}).
		Preload("Sections.Midi").
		Preload("Sections.Instruments", func(db *gorm.DB) *gorm.DB {
			return db.Order("user_song_instruments.sort_order ASC")
		}).
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Order("user_tags.sort_order ASC")
		}).
		Preload("Genres", func(db *gorm.DB) *gorm.DB {
			return db.Order("user_genres.sort_order ASC")
		}).
		Where("user_id = ?", userId).
		Find(&songs)
	if result.Error != nil {
		return nil, result.Error
	}
	return songs, nil
}

type SongSearchCond struct {
	UserIds     []uint `json:"user_ids"`
	TagIds      []uint `json:"tag_ids"`
	GenreIds    []uint `json:"genre_ids"`
	SectionName string `json:"section_name"`
	OrderBy     string `json:"order_by"`
	Ascending   bool   `json:"ascending"`
}

// ORDER句の引数を生成(create_at ASCなど)
func (cond SongSearchCond) buildOrderArg() string {
	order := "ASC"
	orderColumn := cond.OrderBy
	if cond.OrderBy == "" {
		orderColumn = "created_at"
	}

	if cond.Ascending {
		order = "ASC"
	} else {
		order = "DESC"
	}
	orderArg := orderColumn + " " + order //"created_at DESC"
	return orderArg
}

// 検索
func (us *UserSong) Search(db *gorm.DB, cond SongSearchCond) ([]UserSong, error) {
	var songs []UserSong
	var result *gorm.DB

	//tag, genreからsong_idを絞る
	var songIds []uint
	tmpSongs, _ := us.preSearchByUIdTagIdsAndGenreIds(db, cond.UserIds, cond.TagIds, cond.GenreIds)
	for _, v := range tmpSongs {
		songIds = append(songIds, v.ID)
	}
	fmt.Println("pre search songIds = ", songIds)

	//songIdsとsectionNameで再検索
	query := "id IN(?) AND user_id IN (?)"
	args := []interface{}{songIds, cond.UserIds}
	orderArg := cond.buildOrderArg() //"created_at DESC"

	db = db.Preload("Audio").
		Preload("Instruments", func(db *gorm.DB) *gorm.DB {
			return db.Order("user_song_instruments.sort_order ASC")
		}).
		Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("user_song_sections.sort_order ASC")
		}).
		Preload("Sections.Midi").
		Preload("Tags", func(db *gorm.DB) *gorm.DB {
			return db.Order("user_tags.sort_order ASC")
		}).
		Preload("Genres", func(db *gorm.DB) *gorm.DB {
			return db.Order("user_genres.sort_order ASC")
		})
	if cond.SectionName != "" { //sectionName指定がある場合
		db.Preload("Sections", "name=?", cond.SectionName, func(db *gorm.DB) *gorm.DB {
			return db.Order("user_song_sections.sort_order ASC")
		})
	} else {
		db.Preload("Sections", func(db *gorm.DB) *gorm.DB {
			return db.Order("user_song_sections.sort_order ASC")
		})
	}
	db.Preload("Sections.AudioRanges", func(db *gorm.DB) *gorm.DB {
		return db.Order("user_audio_ranges.sort_order ASC")
	}).Preload("Sections.Instruments", func(db *gorm.DB) *gorm.DB {
		return db.Order("user_song_instruments.sort_order ASC")
	})
	result = db.Where(query, args...).Order(orderArg).Find(&songs)

	if result.RowsAffected == 0 {
		return []UserSong{}, nil
	}
	if result.Error != nil {
		return []UserSong{}, result.Error
	}
	return songs, nil
}

// Searchのための処理
func (us *UserSong) preSearchByUIdTagIdsAndGenreIds(db *gorm.DB, userIds []uint, tagIds []uint, genreIds []uint) ([]UserSong, error) {
	var songs []UserSong
	//var song *UserSong
	var result *gorm.DB

	isTagConditonActive := len(tagIds) > 0
	isGenreConditionActive := len(genreIds) > 0

	//tag, genreからsong_idを絞る
	query := "user_songs.user_id IN (?)"
	args := []interface{}{userIds}
	db = db.Model(&UserSong{}).Distinct("user_songs.id")
	if isTagConditonActive {
		db.Joins(
			"INNER JOIN usersongs_tags ust ON user_songs.id=ust.user_song_id " +
				"INNER JOIN user_tags tags ON ust.user_tag_id=tags.id")
		query += " AND tags.id IN(?)"
		args = append(args, tagIds)
	}
	if isGenreConditionActive {
		db.Joins(
			"INNER JOIN usersongs_genres usg ON user_songs.id=usg.user_song_id " +
				"INNER JOIN user_genres genres ON usg.user_genre_id=genres.id")
		query += " AND genres.id IN(?)"
		args = append(args, genreIds)
	}

	result = db.Where(query, args...).Find(&songs)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return songs, nil
}
func (us *UserSong) Update(db *gorm.DB) error {
	fmt.Println("@@@@update")
	//Sections.InstrumentsとSong.Instrumentsを同時に作成できないため、Sectionsを後で保存する
	result := db.Session(&gorm.Session{FullSaveAssociations: true}).Omit("Tags.*", "Genres.*", "created_at", "Sections").Save(&us)
	if err := result.Error; err != nil {
		return err
	}
	if len(us.Sections) > 0 {
		//song.Instrumentsに付与されたIDをsections.Instrumentsに紐付けてSAVE
		instruments := us.Instruments
		//for地獄
		for i, sec := range us.Sections {
			sec.UserSongId = us.ID
			for j, instSec := range sec.Instruments {
				for _, inst := range instruments {
					if inst.Name == instSec.Name {
						instSec.ID = inst.ID
						instSec.UserSongId = inst.UserSongId
					}
				}
				sec.Instruments[j] = instSec
			}
			us.Sections[i] = sec
		}
		result2 := db.Model(&UserSongSection{}).Session(&gorm.Session{FullSaveAssociations: true}).Omit("Instruments.*").Save(&us.Sections)
		if err := result2.Error; err != nil {
			return err
		}
	}
	return nil
}
func (us *UserSong) Delete(db *gorm.DB) error {
	//audio,midiもまとめて削除
	result := db.Omit("Tags.*", "Genres.*").Delete(&us)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

var PlaylistSuffix = "_hls"

// s3ファイルの格納場所を返す
func (us *UserSong) SetMediaUrls() error {

	//audio
	if err := us.SetAudioUrlGet(); err != nil {
		return err
	}
	if err := us.SetAudioUrlPut(); err != nil {
		return err
	}
	//midi
	if err := us.SetMidiUrlGet(); err != nil {
		return err
	}
	if err := us.SetMidiUrlPut(); err != nil {
		return err
	}
	return nil
}
func (us *UserSong) SetAudioUrlGet() error {
	Backend := conf.BACKEND_URL
	fmt.Println("@@@@SetAudioUrlGet")
	if us.Audio.Name != "" {
		audio := &us.Audio
		//get urlはm3u8ファイルを書き換える必要があるためバックエンドを指定する
		get := Backend + "/hls/" + strconv.Itoa(int(us.ID))
		audio.Url.Get = get
	}
	return nil
}
func (us *UserSong) SetAudioUrlPut() error {
	fmt.Println("@@@@SetAudioUrlPut")
	if us.Audio.Name != "" {
		audio := &us.Audio
		put, err := utils.GenerateSignedUrl(us.GetFolderName()+audio.Name, http.MethodPut, conf.PRESIGNED_DURATION)
		if err != nil {
			return err
		}
		audio.Url.Put = put
	}
	return nil
}

// midiファイルの格納場所
// presigned urlを返す
func (us *UserSong) SetMidiUrlGet() error {
	for _, section := range us.Sections {
		if section.Midi.Name != "" {
			midi := &section.Midi
			path := us.GetFolderName() + midi.Name
			get, err := utils.GenerateSignedUrl(path, http.MethodGet, conf.PRESIGNED_DURATION)
			if err != nil {
				return err
			}
			midi.Url.Get = get
		}
	}
	return nil
}
func (us *UserSong) SetMidiUrlPut() error {
	for _, section := range us.Sections {
		if section.Midi.Name != "" {
			midi := &section.Midi
			path := us.GetFolderName() + midi.Name
			put, err2 := utils.GenerateSignedUrl(path, http.MethodPut, conf.PRESIGNED_DURATION)
			if err2 != nil {
				return err2
			}
			midi.Url.Put = put
		}
	}
	return nil
}

// m3u8ファイルの名前を返す
// (オーディオファイル)_hls.m3u8というルールになっている
func (us *UserSong) GetHLSName() string {
	audio := &us.Audio
	n := strings.ReplaceAll(audio.Name, filepath.Ext(audio.Name), "")
	return n + PlaylistSuffix + ".m3u8"
}
func (us *UserSong) GetFolderName() string {
	folder := strconv.Itoa(int(us.UserId)) + "/"
	return folder
}

// 中間テーブルのrelationを削除
func (us *UserSong) DeleteTagRelation(db *gorm.DB, tag *UserTag) error {
	if err := db.Model(us).Association("Tags").Delete(tag); err != nil {
		return err
	}
	return nil
}
func (us *UserSong) DeleteGenreRelation(db *gorm.DB, genre *UserGenre) error {
	if err := db.Model(us).Association("Genres").Delete(genre); err != nil {
		return err
	}
	return nil
}
func (us UserSong) GetID() uint {
	return us.ID
}
