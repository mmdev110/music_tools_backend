package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/app/models"
	"example.com/app/utils"
	"golang.org/x/exp/slices"
)

func Test_SaveTags(t *testing.T) {
	data := models.GetTestData()

	newTags := append([]models.UserTag{}, data.Tags...)
	newTags = append(newTags, []models.UserTag{
		{UserId: data.User.ID, Name: "new tag1"},
		{UserId: data.User.ID, Name: "new Tag2"},
	}...)

	updatedTags := append([]models.UserTag{}, data.Tags...)
	updatedTags[0].Name = "updated tag1"
	updatedTags[2].Name = "updated tag2"

	tests := []struct {
		name         string
		tags         []models.UserTag
		authRequired bool
		statusCode   int
	}{
		{"can save new tags", newTags, true, http.StatusOK},
		{"can update existing tags", updatedTags, true, http.StatusOK},
		{"error with no tags", nil, true, http.StatusBadRequest},
		{"error with no authorization", nil, false, http.StatusBadRequest},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h.DB = TestDB.Begin()
			defer h.DB.Rollback()

			models.InsertTestData(h.DB)

			str, err := utils.ToJSON(test.tags)
			if err != nil {
				t.Error(err)
			}

			req := httptest.NewRequest(http.MethodPost, ts.URL+"/tags", strings.NewReader(str))
			req.RequestURI = ""
			if test.authRequired {
				addAuthorizationHeader(req, data.User)
			}
			res, err2 := ts.Client().Do(req)
			if err2 != nil {
				t.Error(err2)
			}
			defer res.Body.Close()

			got_status := res.StatusCode
			want_status := test.statusCode
			if got_status != want_status {
				t.Errorf("status_code: got %d, want %d", got_status, want_status)
			}
			if test.name == "can save new tags" {
				tmp := models.UserTag{}
				tags, _ := tmp.GetAllByUserId(h.DB, data.User.ID)
				//utils.PrintStruct(tags)
				got_num := len(tags)
				want_num := len(newTags)
				if got_num != want_num {
					t.Errorf("tags_num: got %d, want %d", got_num, want_num)
				}
			}
			if test.name == "can update existing tags" {
				tmp := models.UserTag{}
				tags, _ := tmp.GetAllByUserId(h.DB, data.User.ID)
				//utils.PrintStruct(tags)
				got_num := len(tags)
				want_num := len(updatedTags)
				if got_num != want_num {
					t.Errorf("tags_num: got %d, want %d", got_num, want_num)
				}
				names := []string{}
				for _, tag := range tags {
					names = append(names, tag.Name)
				}
				fmt.Println(names)
				if !slices.Contains(names, updatedTags[0].Name) {
					t.Errorf("%s failed", updatedTags[0].Name)
				}
				if !slices.Contains(names, updatedTags[2].Name) {
					t.Errorf("%s failed", updatedTags[2].Name)
				}
			}

		})
	}

}
func Test_GetTags(t *testing.T) {

	t.Run("can get tags", func(t *testing.T) {
		h.DB = TestDB.Begin()
		defer h.DB.Rollback()

		data, err := models.InsertTestData(h.DB)
		if err != nil {
			t.Error(err)
		}
		req := httptest.NewRequest(http.MethodGet, ts.URL+"/tags", nil)
		req.RequestURI = ""
		addAuthorizationHeader(req, data.User)

		res, err := ts.Client().Do(req)
		if err != nil {
			t.Error(err)
		}
		defer res.Body.Close()

		got_code := res.StatusCode
		want_code := http.StatusOK
		if got_code != want_code {
			t.Errorf("status_code: got %d, want %d", got_code, want_code)
		}
		var res_tags []models.UserTag
		if err := utils.BodyToStruct(res.Body, &res_tags); err != nil {
			t.Error(err)
		}
		//utils.PrintStruct(res_tags)
		got_num := len(res_tags)
		want_num := 3
		//utils.PrintStruct(res_genres)
		if got_num != want_num {
			t.Errorf("status_code: got %d, want %d", got_num, want_num)
		}
	})

}
