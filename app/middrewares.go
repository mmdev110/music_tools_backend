package main

import (
	"fmt"
	"net/http"

	"example.com/app/utils"
)

func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("auth middreware")
		authHeader := r.Header.Get("Authorization")
		claim, err := utils.Authenticate(authHeader)
		//for key, value := range r.Header {
		//	fmt.Printf("%v: %v\n", key, value)
		//}
		if err != nil {
			//w.WriteHeader(http.StatusUnauthorized)
			utils.ErrorJSON(w, err)
			return
		}
		userId := claim.UserId
		fmt.Println(userId)
		ctx := utils.SetUIDInContext(r.Context(), userId)
		next(w, r.WithContext(ctx))
	}
}

func enableCORS(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("path = %s. method = %s\n", r.URL.Path, r.Method)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		//preflight対応
		if r.Method == http.MethodOptions {
			fmt.Println("@@@preflight response")
			//w.Header().Set("Access-Control-Max-Age", strconv.Itoa(86400))
			utils.ResponseJSON(w, nil, http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}
