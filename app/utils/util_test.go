package utils

import "testing"

type Temp struct {
	ID uint
}

func (t Temp) PrintID() uint {
	return t.ID
}
func TestFindRemoved(t *testing.T) {

	t.Run("before 3, after 2", func(t *testing.T) {
		before := []Temp{{ID: 1}, {ID: 2}, {ID: 3}}
		after := []Temp{{ID: 1}, {ID: 2}}
		removed := FindRemoved(before, after)
		if removed[0].ID != 3 {
			t.Error("failed!!")
		}
	})
	t.Run("before 3, after 3", func(t *testing.T) {
		before := []Temp{{ID: 1}, {ID: 2}, {ID: 3}}
		after := []Temp{{ID: 1}, {ID: 2}, {ID: 3}}
		removed := FindRemoved(before, after)
		if len(removed) != 0 {
			t.Errorf("failed!! %v", removed)
		}
	})
	t.Run("When after is empty, return all values", func(t *testing.T) {
		before := []Temp{{ID: 1}, {ID: 2}, {ID: 3}}
		after := []Temp{}
		removed := FindRemoved(before, after)
		if len(removed) != 3 {
			t.Error("failed!!")
		}
	})

}
