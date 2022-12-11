package models

import (
	"encoding/json"

	_ "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// clientとの通信に使うstruct
type UserLoopInput struct {
	ID            uint            `json:"id"`
	Progressions  []string        `json:"progressions"`
	Key           int             `json:"key"`
	Scale         string          `json:"scale"`
	Memo          string          ` json:"memo"`
	MidiRoots     []int           `json:"midi_roots"`
	UserLoopAudio []UserLoopAudio `json:"user_loop_audio"`
	UserLoopMidi  []UserLoopMidi  `json:"user_loop_midi"`
}

// DBに格納するためのstruct
// UserLoopInputの配列要素をstring化している
type UserLoop struct {
	ID     uint `gorm:"primarykey" json:"id"`
	UserId uint `gorm:"not null" json:"user_id"`
	//コード進行をcsv化したもの
	//["Am7","","","Dm7"]->"Am7,,,Dm7"
	Progressions string `json:"progressions"`
	Key          int    `json:"key"`
	Scale        string `json:"scale"`
	//オーディオファイル
	UserLoopAudio []UserLoopAudio `json:"user_loop_audio"`
	//midiファイル
	UserLoopMidi []UserLoopMidi `json:"midi_path"`
	//midiファイル内でルートとなるノートのindexをcsv化したもの
	//[1,2,3,4]->"1,2,3,4"
	MidiRoots string ` json:"midi_roots"`
	Memo      string ` json:"memo"`
	gorm.Model
}

func (ul *UserLoop) Create() error {
	result := DB.Create(&ul)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (ul *UserLoop) GetByID(id uint) error {
	result := DB.First(&ul, id)
	if result.RowsAffected == 0 {
		return nil
	}
	return nil
}
func (ul *UserLoop) GetAllByUserId(userId uint) []UserLoop {
	var loops []UserLoop
	result := DB.Where("user_id = ?", userId).Find(&loops)
	if result.RowsAffected == 0 {
		return nil
	}
	return loops
}
func (ul *UserLoop) Update() error {
	result := DB.Model(&ul).Debug().Updates(ul)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}
func (ul *UserLoop) delete() {
	DB.Delete(&ul, ul.ID)
}

func (ul *UserLoop) ApplyULInputToUL(ulInput UserLoopInput) {
	prog, _ := json.Marshal(ulInput.Progressions)
	midiroots, _ := json.Marshal(ulInput.MidiRoots)
	//ul.ID = ulInput.ID
	ul.Progressions = string(prog)
	ul.Key = ulInput.Key
	ul.Scale = ulInput.Scale
	ul.MidiRoots = string(midiroots)
	ul.Memo = ulInput.Memo
	ul.UserLoopAudio = ulInput.UserLoopAudio
	ul.UserLoopMidi = ulInput.UserLoopMidi
}
func (uli *UserLoopInput) ApplyULtoULInput(ul UserLoop) {
	var prog []string
	json.Unmarshal([]byte(ul.Progressions), &prog)
	var midiroots []int
	json.Unmarshal([]byte(ul.MidiRoots), &midiroots)
	uli.ID = ul.ID
	uli.Progressions = prog
	uli.Key = ul.Key
	uli.Scale = ul.Scale
	uli.MidiRoots = midiroots
	uli.Memo = ul.Memo
	uli.UserLoopAudio = ul.UserLoopAudio
	uli.UserLoopMidi = ul.UserLoopMidi
}
