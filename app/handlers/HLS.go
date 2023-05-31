package handlers

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"example.com/app/conf"
	"example.com/app/models"
	"example.com/app/utils"
)

var r_aac = regexp.MustCompile(`\.aac`)

// S3からm3u8ファイルを取得し、その中身のaacファイルのアドレスをpresigned urlに置き換えた上で返す
func HLSHandler(w http.ResponseWriter, r *http.Request) {
	//TODO
	//認証つけた方が良い

	//user := getUserFromContext(r.Context())
	//fmt.Printf("userid in handler = %d\n", user.ID)
	//パラメータ取り出し
	str := strings.TrimPrefix(r.URL.Path, "/hls/")
	int, _ := strconv.Atoi(str)
	userSongId := uint(int)

	var us = models.UserSong{}
	us.GetByID(userSongId)
	//if us.UserId != user.ID {
	//	utils.ErrorJSON(w, errors.New("user mismatch"))
	//}
	fmt.Println(us.GetFolderName() + us.GetHLSName())
	presignedUrl, err := utils.GenerateSignedUrl(us.GetFolderName()+us.GetHLSName(), http.MethodGet, conf.PRESIGNED_DURATION)
	if err != nil {
		utils.ErrorJSON(w, err)
	}
	resp, err := http.Get(presignedUrl)
	if err != nil {
		utils.ErrorJSON(w, err)
	}
	defer resp.Body.Close()

	//m3u8ファイルの書き換え
	scanner := bufio.NewScanner(resp.Body)
	//b, _ := io.ReadAll(resp.Body)
	//oldHLS := string(b)
	newHLS := ""
	for scanner.Scan() {
		line := scanner.Text()
		if r_aac.MatchString(line) {
			aac, _ := url.QueryUnescape(line)
			presigned, _ := utils.GenerateSignedUrl(us.GetFolderName()+aac, http.MethodGet, conf.PRESIGNED_DURATION)
			newHLS = newHLS + presigned
		} else {
			newHLS = newHLS + line
		}
		newHLS = newHLS + "\n"
	}
	fmt.Println(newHLS)
	w.Header().Set("Content-Type", "application/x-mpegURL")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(newHLS))
	return
}
