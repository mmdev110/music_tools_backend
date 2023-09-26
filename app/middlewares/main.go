package middlewares

import (
	"fmt"
	"net/http"

	"example.com/app/auth"
	"example.com/app/conf"
	"example.com/app/customError"
	"example.com/app/utils"
)

func RequireAuth(next http.HandlerFunc, auth *auth.Auth) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println("auth")
		authHeader := r.Header.Get("Authorization")
		claim, err := auth.AuthCognito(authHeader)
		if err != nil {
			utils.ErrorJSON(w, customError.Others, err)
			return
		}
		uuid := claim.UUID
		email := claim.Email
		ctx := utils.SetParamsInContext(r.Context(), uuid, email)
		next(w, r.WithContext(ctx))
	}
}

func EnableCORS(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("path = %s. method = %s\n", r.URL.Path, r.Method)
		//fmt.Println(r.Header)
		w.Header().Set("Access-Control-Allow-Origin", conf.FRONTEND_URL)
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		//preflight対応
		if r.Method == http.MethodOptions {
			fmt.Println("@@@preflight response")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, X-CSRF-Token, Authorization")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			//w.Header().Set("Access-Control-Max-Age", strconv.Itoa(86400))
			utils.ResponseJSON(w, nil, http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
