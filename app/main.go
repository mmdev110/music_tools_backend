package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"example.com/app/handlers"
	"example.com/app/models"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	//DB接続
	err := models.Init(false)
	if err != nil {
		log.Fatal(err)
	}
	//playground()
	web_server()
}

// DB接続無しのwebサーバ
func web_server() {
	fmt.Println("web_simple")

	//ハンドラ登録
	mux := registerHandlers()
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
func registerHandlers() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/_chk", handlers.ChkHandler)
	mux.HandleFunc("/signin", handlers.SignInHandler)
	mux.HandleFunc("/signup", handlers.SignUpHandler)
	mux.HandleFunc("/refresh", handlers.RefreshHandler)
	mux.HandleFunc("/reset_password", handlers.ResetPasswordHandler)
	mux.HandleFunc("/email_confirm", handlers.EmailConfirmationHandler)
	mux.HandleFunc("/user", requireAuth(handlers.UserHandler))
	mux.HandleFunc("/list", requireAuth(handlers.ListHandler))
	mux.HandleFunc("/tags", requireAuth(handlers.TagHandler))
	mux.HandleFunc("/genres", requireAuth(handlers.GenreHandler))
	mux.HandleFunc("/song/", requireAuth(handlers.SongHandler))
	mux.HandleFunc("/delete_song", requireAuth(handlers.DeleteSong))
	mux.HandleFunc("/hls/", handlers.HLSHandler)
	mux.HandleFunc("/test", handlers.TestHandler)

	return enableCORS(mux)
}
