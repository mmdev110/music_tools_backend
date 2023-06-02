package models

import (
	"fmt"
	"os"

	"example.com/app/conf"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

const SQLITE_FILE = "data.db"

func Init(isTesting bool) error {
	var err error
	if isTesting {
		db, err2 := connectSQLite()
		DB = db
		err = err2
	} else {
		db, err2 := connectMySQL()
		DB = db
		err = err2
	}
	if err != nil {
		return err
	}

	fmt.Println("@@@DBconnection success")
	migrateModels(DB)
	return nil
}
func connectMySQL() (*gorm.DB, error) {
	user := conf.MYSQL_USER
	password := conf.MYSQL_PASSWORD
	db_name := conf.MYSQL_DATABASE
	db_host := "db:" + conf.MYSQL_PORT
	dsn := user + ":" + password + "@tcp(" + db_host + ")/" + db_name + "?charset=utf8mb4&parseTime=True&loc=Asia%2FTokyo"
	fmt.Printf("DSN = %s\n", dsn)
	//"gorm:gorm@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True&loc=Local", // data source name
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}), &gorm.Config{})
	return db, err
}
func connectSQLite() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(SQLITE_FILE), &gorm.Config{})
	return db, err
}
func migrateModels(db *gorm.DB) {
	fmt.Println("@@@migration")
	db.AutoMigrate(&User{}, &UserSong{}, &UserSongSection{}, &UserTag{}, &UserGenre{}, &Session{})
	//db.AutoMigrate(&UserSongSection{})
	db.AutoMigrate(&UserSongAudio{})
	db.AutoMigrate(&UserSectionMidi{})
}
func ClearSQLiteDB() {
	os.Remove(SQLITE_FILE)
}
