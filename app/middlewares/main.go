package middlewares

import (
	"fmt"
	"net/http"

	"example.com/app/conf"
	"example.com/app/customError"
	"example.com/app/utils"
)

func RequireAuth(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println("auth")
		authHeader := r.Header.Get("Authorization")
		claim, err := utils.Authenticate(authHeader, "access")
		//for key, value := range r.Header {
		//	fmt.Printf("%v: %v\n", key, value)
		//}
		if err != nil {
			//w.WriteHeader(http.StatusUnauthorized)
			utils.ErrorJSON(w, customError.Others, err)
			return

		}
		userId := claim.UserId
		//fmt.Println(userId)
		ctx := utils.SetUIDInContext(r.Context(), userId)
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