package models

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"example.com/app/conf"
	"example.com/app/utils"
	_ "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// song情報
type UserSong struct {
	ID       uint              `gorm:"primarykey" json:"id"`
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
	Tags        []UserTag            `gorm:"many2many:usersongs_tags" json:"tags"`
	Instruments []UserSongInstrument `json:"instruments"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func (us *UserSong) Create() error {
	//Sections.Instrumentsがsongs.Instrumentsに依存しており、
	//同時に生成できないため、CREATE文を分ける
	sections := us.Sections
	us.Sections = []UserSongSection{}
	result := DB.Debug().Omit("Tags.*", "Genres.*").Create(&us)
	if result.Error != nil {
		return result.Error
	}
	//Instrumentsに付与されたIDをsectionsに紐付けてCREATE
	us.Sections = sections
	instruments := us.Instruments
	utils.PrintStruct(instruments)
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
		sec.Create()
		us.Sections[i] = sec
	}

	us.SetMediaUrls()
	return nil
}

// songを返す
func (us *UserSong) GetByID(id uint) *gorm.DB {
	result := DB.Model(&UserSong{}).Preload("Audio").Preload("Instruments").Preload("Sections").Preload("Sections.Midi").Preload("Tags").Preload("Genres").Debug().First(&us, id)
	if result.RowsAffected == 0 {
		return result
	}
	err := us.SetMediaUrls()
	if err != nil {
		fmt.Println(err)
	}
	return result
}

type SongSearchCond struct {
	TagIds      []uint `json:"tag_ids"`
	GenreIds    []uint `json:"genre_ids"`
	SectionName string `json:"section_name"`
}

// userIdに紐づくsong(検索条件があればそれも考慮する)
func (us *UserSong) GetByUserId(userId uint, cond SongSearchCond) ([]UserSong, error) {
	var songs []UserSong
	var result *gorm.DB
	//うまいやり方を考える
	//if len(cond.TagIds) == 0 && len(cond.GenreIds) == 0 && cond.SectionName == "" { //検索条件なし=uidのみで検索
	//	fmt.Println("no search cond")
	//	result = DB.Preload("Audio").Preload("Sections").Preload("Sections.Midi").Preload("Tags").Preload("Genres").Debug().Where("user_id=?", userId).Find(&songs)
	//} else {
	//	result = DB.Preload("Audio").Preload("Sections", "name=?", cond.SectionName).Joins("INNER JOIN user_song_sections sec ON sec.user_song_id=user_songs.id ").Preload("Tags", "id IN (?)", cond.TagIds).Preload("Genres", "id IN (?)", cond.GenreIds).Debug().Where("user_id=? AND sec.name=?", userId, cond.SectionName).Find(&songs)
	//	//result = DB.Debug().Joins("INNER JOIN userloops_tags ult ON user_loops.id=ult.user_loop_id").Joins("INNER JOIN user_loop_tags tags ON tags.id=ult.user_loop_tag_id").Where("user_loops.user_id=? AND tags.id IN ?", userId, condition.TagIds).Find(&songs)
	//}
	isTagConditonActive := len(cond.TagIds) > 0
	isGenreConditionActive := len(cond.GenreIds) > 0
	//tagIdsを持ってるsongを検索
	var songIdsWithTags []uint
	if isTagConditonActive {
		songWithTags, _ := us.getSongByTagIds(userId, cond.TagIds)
		for _, v := range songWithTags {
			songIdsWithTags = append(songIdsWithTags, v.ID)
		}
	}
	fmt.Println("songIdsWithTags: ", songIdsWithTags)

	//genreIdsを持ってるsongを検索
	var songIdsWithGenres []uint
	if isGenreConditionActive {
		songWithGenres, _ := us.getSongByGenreIds(userId, cond.GenreIds)
		//DB.Debug().Joins("INNER JOIN usersongs_genres usg ON user_songs.id=usg.user_song_id").Where("user_id=? AND usg.user_genre_id IN ?", userId, cond.GenreIds).Find(&songWithGenres)
		for _, v := range songWithGenres {
			songIdsWithGenres = append(songIdsWithGenres, v.ID)
		}
	}
	fmt.Println("songIdsWithGenres: ", songIdsWithGenres)
	//両方に含まれるidを抽出
	var commonIds []uint
	if isTagConditonActive && isGenreConditionActive {
		commonIds = utils.Intersect(songIdsWithTags, songIdsWithGenres)
	} else if !isTagConditonActive {
		commonIds = songIdsWithGenres
	} else if !isGenreConditionActive {
		commonIds = songIdsWithTags
	}
	fmt.Println(isTagConditonActive)
	fmt.Println(isGenreConditionActive)
	fmt.Println("commonIds: ", commonIds)

	//そのidの中から、sectionNameで絞り込み
	db := DB.Debug().Preload("Audio").Preload("Tags").Preload("Genres").Preload("Instruments")
	query := "user_id=?"
	args := []interface{}{userId}

	if isTagConditonActive || isGenreConditionActive { //タグ、ジャンル検索条件がある場合、userSongId条件追加
		query += " AND user_songs.id in ?"
		args = append(args, commonIds)
	}
	//sectionName
	if cond.SectionName != "" { //sectionName指定がある場合
		db.Preload("Sections.Instruments", "name=?", cond.SectionName).Joins("INNER JOIN user_song_sections sec ON sec.user_song_id=user_songs.id ")
		query += " AND sec.name=?"
		args = append(args, cond.SectionName)
	} else {
		db.Preload("Sections.Instruments")
	}
	result = db.Where(query, args...).Find(&songs)

	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	for i := range songs {
		err := songs[i].SetMediaUrls()
		if err != nil {
			fmt.Println(err)
		}
	}
	return songs, nil
}
func (us *UserSong) Update() error {
	fmt.Println("@@@@update")
	result := DB.Debug().Session(&gorm.Session{FullSaveAssociations: true}).Omit("Tags.*", "Genres.*", "created_at").Save(&us)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}
func (us *UserSong) Delete() error {
	//audio,midiもまとめて削除
	result := DB.Debug().Omit("Tags.*", "Genres.*").Delete(&us)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

var PlaylistSuffix = "_hls"

// s3ファイルの格納場所を返す
func (us *UserSong) SetMediaUrls() error {

	//audio
	if err := us.setAudioUrl(); err != nil {
		return err
	}
	//midi
	if err := us.setMidiUrl(); err != nil {
		return err
	}
	return nil
}
func (us *UserSong) setAudioUrl() error {
	Backend := conf.BACKEND_URL
	fmt.Println("@@@@setmediaurl")
	//fmt.Println(us.Audio)
	//audio
	if us.Audio.Name != "" {
		audio := &us.Audio

		//put urlはpresigned URL
		//get urlはm3u8ファイルを書き換える必要があるためバックエンドを指定する
		get := Backend + "/hls/" + strconv.Itoa(int(us.ID))
		audio.Url.Get = get
		put, err := utils.GenerateSignedUrl(us.GetFolderName()+audio.Name, http.MethodPut, conf.PRESIGNED_DURATION)
		if err != nil {
			return err
		}
		fmt.Println("get: ", get)
		fmt.Println("put: ", put)
		audio.Url.Put = put
	}
	return nil
}

// midiファイルの格納場所
// presigned urlを返す
func (us *UserSong) setMidiUrl() error {
	for _, section := range us.Sections {
		if section.Midi.Name != "" {
			midi := &section.Midi
			path := strconv.Itoa(int(us.UserId)) + "/" + strconv.Itoa(int(section.ID)) + "/" + midi.Name
			get, err := utils.GenerateSignedUrl(path, http.MethodGet, conf.PRESIGNED_DURATION)
			if err != nil {
				return err
			}
			put, err2 := utils.GenerateSignedUrl(path, http.MethodPut, conf.PRESIGNED_DURATION)
			if err != nil {
				return err2
			}
			midi.Url.Get = get
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
func (us *UserSong) DeleteTagRelation(tag *UserTag) error {
	if err := DB.Model(us).Association("Tags").Delete(tag); err != nil {
		return err
	}
	return nil
}
func (us *UserSong) DeleteGenreRelation(genre *UserGenre) error {
	if err := DB.Model(us).Association("Genres").Delete(genre); err != nil {
		return err
	}
	return nil
}
func (us UserSong) GetID() uint {
	return us.ID
}

func (us *UserSong) getSongByTagIds(userId uint, tagIds []uint) ([]UserSong, error) {
	fmt.Println("getSongByTagIds")
	var songWithTags []UserSong
	result := DB.Debug().Preload("Tags").Joins("INNER JOIN usersongs_tags ust ON user_songs.id=ust.user_song_id AND ust.user_tag_id IN ?", tagIds).Where("user_id=?", userId).Find(&songWithTags)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	//同じsongIdの要素が複数取得されてしまうので、ユニークにする
	uniq := utils.Uniq(songWithTags)
	return uniq, nil
}
func (us *UserSong) getSongByGenreIds(userId uint, genreIds []uint) ([]UserSong, error) {
	fmt.Println("getSongByGenreIds")
	var songWithGenres []UserSong
	result := DB.Debug().Preload("Genres").Joins("INNER JOIN usersongs_genres usg ON user_songs.id=usg.user_song_id AND usg.user_genre_id IN ?", genreIds).Where("user_id=?", userId).Find(&songWithGenres)
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	//同じsongIdの要素が複数取得されてしまうので、ユニークにする
	uniq := utils.Uniq(songWithGenres)
	return uniq, nil
}
