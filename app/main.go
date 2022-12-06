package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"example.com/app/handlers"
	"example.com/app/models"
	"example.com/app/utils"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	//AWS初期設定
	utils.ConfigureAWS()
	exec_code()
	web_simple()
}

// 単純なコード
func exec_code() {
	fmt.Println("exec_code")
	url, _ := utils.GenerateSignedUrl(1, "test.mp3", http.MethodGet, 100)
	fmt.Println(url)
}

// DB接続無しのwebサーバ
func web_simple() {
	fmt.Println("web_simple")
	//DB接続
	err := models.Init()
	if err != nil {
		log.Fatal(err)
	}

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

	log.Fatal(server.ListenAndServe())
}
func registerHandlers() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/_chk", handlers.ChkHandler)
	mux.HandleFunc("/signin", handlers.SignInHandler)
	mux.HandleFunc("/signup", handlers.SignUpHandler)
	mux.HandleFunc("/user", requireAuth(handlers.UserHandler))
	mux.HandleFunc("/list", requireAuth(handlers.ListHandler))
	mux.HandleFunc("/loop/", requireAuth(handlers.LoopHandler))

	return enableCORS(mux)
}
