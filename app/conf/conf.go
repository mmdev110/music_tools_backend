package conf

import (
	"os"
	"time"
)

// 定数、環境変数など
// アクセストークン有効期間
var TOKEN_DURATION = 1 * time.Minute

// リフレッシュトークン有効期間
var REFRESH_DURATION = 10 * time.Minute
var COOKIE_DOMAIN = "localhost"
var SESSION_ID_KEY = "_session_id"
var BACKEND_URL = os.Getenv("BACKEND_URL")
var FRONTEND_URL = os.Getenv("FRONTEND_URL")

// openssl genrsa -out private-key.pem 2048
var HMAC_SECRET_KEY = os.Getenv("HMAC_SECRET_KEY")

// MYSQL
var MYSQL_ROOT_PASSWORD = os.Getenv("MYSQL_ROOT_PASSWORD")
var MYSQL_DATABASE = os.Getenv("MYSQL_DATABASE")
var MYSQL_USER = os.Getenv("MYSQL_USER")
var MYSQL_PASSWORD = os.Getenv("MYSQL_PASSWORD")
var MYSQL_PORT = os.Getenv("MYSQL_PORT")
var MYSQL_HOST = os.Getenv("MYSQL_HOST")

// AWS
var AWS_REGION = os.Getenv("AWS_REGION")
var AWS_ACCESS_KEY_ID = os.Getenv("AWS_ACCESS_KEY_ID")
var AWS_SECRET_ACCESS_KEY = os.Getenv("AWS_SECRET_ACCESS_KEY")
var AWS_BUCKET_NAME = os.Getenv("AWS_BUCKET_NAME")
var AWS_MEDIACONVERT_ENDPOINT = os.Getenv("AWS_MEDIACONVERT_ENDPOINT")
var AWS_CLOUDFRONT_DOMAIN = os.Getenv("AWS_CLOUDFRONT_DOMAIN")
var SUPPORT_EMAIL = "support@" + os.Getenv("SUPPORT_EMAIL_DOMAIN")
var PRESIGNED_DURATION = time.Duration(15 * time.Minute)

func OverRideVarsByENV() {
	if os.Getenv("ENV") == "local" {
		TOKEN_DURATION = 1 * time.Minute

		// リフレッシュトークン有効期間
		REFRESH_DURATION = 10 * time.Minute
	}
}
