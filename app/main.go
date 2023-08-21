package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"example.com/app/conf"
	"example.com/app/handlers"
	"example.com/app/models"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
)

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
	//ハンドラ登録
	mux := app.registerHandlers()
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
func (app *Application) registerHandlers() http.Handler {
	h := handlers.Base{
		DB: app.DB,
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/_chk", h.ChkHandler)
	mux.HandleFunc("/signin", h.SignInHandler)
	mux.HandleFunc("/signup", h.SignUpHandler)
	mux.HandleFunc("/refresh", h.RefreshHandler)
	mux.HandleFunc("/signout", h.SignOutHandler)
	mux.HandleFunc("/reset_password", h.ResetPasswordHandler)
	mux.HandleFunc("/email_confirm", h.EmailConfirmationHandler)
	mux.HandleFunc("/user", requireAuth(h.UserHandler))
	mux.HandleFunc("/signin_with_token", requireAuth(h.SignInWithTokenHandler))
	mux.HandleFunc("/list", requireAuth(h.ListHandler))
	mux.HandleFunc("/tags", requireAuth(h.TagHandler))
	mux.HandleFunc("/genres", requireAuth(h.GenreHandler))
	mux.HandleFunc("/song/", requireAuth(h.SongHandler))
	mux.HandleFunc("/delete_song", requireAuth(h.DeleteSong))
	mux.HandleFunc("/hls/", h.HLSHandler)
	mux.HandleFunc("/test", h.TestHandler)

	return enableCORS(mux)
}
