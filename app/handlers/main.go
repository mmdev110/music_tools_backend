package handlers

import (
	"net/http"

	"gorm.io/gorm"

	"example.com/app/auth"
	mw "example.com/app/middlewares"
)

type HandlersConf struct {
	DB        *gorm.DB
	IsTesting bool      //test実行中かどうか
	SendEmail bool      //メール送信実行するか
	Auth      auth.Auth //auth関連
}

func (h *HandlersConf) Handlers() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("/_chk", h.ChkHandler)
	mux.HandleFunc("/signin", h.SignInHandler)
	mux.HandleFunc("/signup", h.SignUpHandler)
	mux.HandleFunc("/refresh", h.RefreshHandler)
	mux.HandleFunc("/signout", h.SignOutHandler)
	mux.HandleFunc("/reset_password", h.ResetPasswordHandler)
	mux.HandleFunc("/email_confirm", h.EmailConfirmationHandler)
	mux.HandleFunc("/auth_with_token", h.AuthWithTokenHandler)
	mux.HandleFunc("/user", mw.RequireAuth(h.UserHandler))
	mux.HandleFunc("/list", mw.RequireAuth(h.SearchSongsHandler))
	mux.HandleFunc("/tags", mw.RequireAuth(h.TagHandler))
	mux.HandleFunc("/genres", mw.RequireAuth(h.GenreHandler))
	mux.HandleFunc("/song/", mw.RequireAuth(h.SongHandler))
	mux.HandleFunc("/delete_song", mw.RequireAuth(h.DeleteSong))
	mux.HandleFunc("/hls/", h.HLSHandler)
	//mux.HandleFunc("/test", h.TestHandler)

	return mw.EnableCORS(mux)
}
