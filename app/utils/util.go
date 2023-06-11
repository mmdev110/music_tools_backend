package utils

import (
	"fmt"

	"golang.org/x/exp/slices"
)

func PrintStruct(struc interface{}) {
	fmt.Printf("%+v\n", struc)
}

// タグの削除
// beforeにあってafterに存在しないものを返す
// beforeがDBに保存された値、afterがフロントからの値で
// 更新時に削除する値を探す際に使用することを想定
type HasGetID interface {
	GetID() uint
}

func FindRemoved[T HasGetID](before, after []T) []T {

	removed := []T{}
	for _, bfr := range before {
		found := false
		for _, aftr := range after {
			if bfr.GetID() == aftr.GetID() {
				found = true
			}
		}
		if !found {
			removed = append(removed, bfr)
		}
	}
	return removed
}
func Uniq[T HasGetID](arr []T) (uniq []T) {

	var pushedIds []uint
	for _, v := range arr {
		if !slices.Contains(pushedIds, v.GetID()) {
			uniq = append(uniq, v)
			pushedIds = append(pushedIds, v.GetID())
		}
	}
	return
}
func Intersect[T uint | string](A, B []T) (ans []T) {
	for _, v := range A {
		if slices.Contains(B, v) {
			ans = append(ans, v)
		}
	}
	return ans
}
