package utils

import (
	"fmt"
)

func PrintStruct(struc interface{}) {
	fmt.Printf("%+v\n", struc)
}

// タグの削除
// beforeにあってafterに存在しないものを返す
// beforeがDBに保存された値、afterがフロントからの値で
// 更新時に削除する値を探す際に使用することを想定
type HasPrintID interface {
	PrintID() uint
}

func FindRemoved[T HasPrintID](before, after []T) []T {

	removed := []T{}
	for _, bfr := range before {
		found := false
		for _, aftr := range after {
			if bfr.PrintID() == aftr.PrintID() {
				found = true
			}
		}
		if !found {
			removed = append(removed, bfr)
		}
	}
	return removed
}
