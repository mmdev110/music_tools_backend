package handlers

import "gorm.io/gorm"

type Base struct {
	DB        *gorm.DB
	IsTesting bool //test実行中かどうか
	SendEmail bool //メール送信実行するか
}
