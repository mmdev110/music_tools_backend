package models

import (
	"testing"

	"example.com/app/utils"
)

// 短いコードの検証用
func TestTags(t *testing.T) {
	err := Init()
	if err != nil {
		t.Fatal(err)
	}
	t.Run("create and delete", func(t *testing.T) {
		tag := UserTag{}
		tag.Name = "testtag2"
		tag.UserId = 6

		if err := tag.Create(); err != nil {
			t.Fatalf("error at tag.Create(): %v", err)
		}
		if err2 := tag.DeleteTagAndRelations([]UserTag{tag}); err2 != nil {
			t.Fatalf("error at tag.Delete(): %v", err2)
		}
	})
	t.Run("tag.GetAllByUserId()", func(t *testing.T) {
		//create
		var tags []UserTag
		userId := uint(8)
		for i, v := range []string{"tag1", "tag2", "tag3"} {
			tag := UserTag{}
			tag.UserId = userId
			tag.Name = v
			tag.SortOrder = i
			if err := tag.Create(); err != nil {
				t.Errorf("error found at tag.Create(): %v", err)
			}
			tags = append(tags, tag)
		}
		//getAll
		emptyTag := UserTag{}
		gotTags, err := emptyTag.GetAllByUserId(userId)
		if err != nil {
			t.Errorf("error found at tag.GetAll(): %v", err)
		}
		for i := range tags {
			want := tags[i]
			got := gotTags[i]
			if want.Name != got.Name {
				t.Errorf("name mismatch: got: %s want: %s", got.Name, want.Name)
			}
			if want.SortOrder != got.SortOrder {
				t.Errorf("sort_order mismatch: got: %d want: %d", got.SortOrder, want.SortOrder)
			}
		}
		//delete
		for _, v := range gotTags {
			if errDelete := v.DeleteTagAndRelations([]UserTag{v}); errDelete != nil {
				t.Errorf("error at Delete: %v", errDelete)
			}
		}
	})
}

func TestDeleteTag(t *testing.T) {
	err := Init()
	if err != nil {
		t.Fatal(err)
	}
	t.Run("get tag and its userloops", func(t *testing.T) {
		//tagid := uint(1)
		tag := UserTag{}
		DB.First(&tag, uint(1))
		result := DB.Debug().Model(&UserTag{}).Preload("UserSongSections").Find(&tag)
		if result.Error != nil {
			t.Errorf("error:= %v", result.Error)
		}
		utils.PrintStruct(tag)
	})
	t.Run("delete tag", func(t *testing.T) {
		//tagid := uint(1)
		tag := UserTag{}
		DB.First(&tag, uint(1))
		//SELECT * FROM user_loops INNER JOIN userloops_tags
		//ON user_loops.id=userloops_tags.user_loops_id
		//WHERE user_loops_tags.user_loop_tag.id=?
		//result := DB.Debug().Model(&UserSongSection{}).Joins("inner join userloops_tags on user_loops.id=userloops_tags.user_loop_id").Where("userloops_tags.user_loop_tag_id = ?", tagid).Find(&ulp)

		result := DB.Debug().Model(&tag).Association("UserSongSections").Clear()
		if result != nil {
			t.Errorf("error:= %v", result)
		}
		utils.PrintStruct(tag)
	})
}
