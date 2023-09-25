package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"example.com/app/auth"
	"example.com/app/conf"
	"example.com/app/handlers"
	"example.com/app/models"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
)

// TODO: Application, handlers.Base, confなどに散らばってる設定を１か所にまとめたいが、
// パッケージ跨ぐと難しい。。
type Application struct {
	DB *gorm.DB
}

func main() {

	app := Application{}
	//DB接続
	db, err := models.Init()
	if err != nil {
		log.Fatal(err)
	}
	app.DB = db

	app.web_server()
}

func (app *Application) web_server() {
	fmt.Println("web")
	//パラメータを詰め込む
	jwks, err := auth.GetJWKS(conf.AWS_REGION, conf.AWS_COGNITO_USER_POOL_ID)
	if err != nil {
		log.Fatal("error at GetJWKS")
	}
	h := handlers.HandlersConf{
		DB:        app.DB,
		SendEmail: true,
		IsTesting: false,
		JWKS:      jwks,
	}
	mux := h.Handlers()
	conf.OverRideVarsByENV()
	//サーバー起動
	server := &http.Server{
		Addr:           ":5000",
		Handler:        mux,
		ReadTimeout:    time.Duration(10 * int64(time.Second)),
		WriteTimeout:   time.Duration(600 * int64(time.Second)),
		MaxHeaderBytes: 1 << 20,
	}
	//

	log.Fatal(server.ListenAndServe())
}
