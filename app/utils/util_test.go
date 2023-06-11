package utils

import "testing"

type Temp struct {
	ID uint
}

func (t Temp) GetID() uint {
	return t.ID
}
func TestFindRemoved(t *testing.T) {
	type Suite struct {
		memo   string
		before []Temp
		after  []Temp
		want   []Temp
	}
	suites := []Suite{
		{
			memo:   "before 3, after 2",
			before: []Temp{{ID: 1}, {ID: 2}, {ID: 3}},
			after:  []Temp{{ID: 1}, {ID: 2}},
			want:   []Temp{{ID: 3}},
		},
		{
			memo:   "before 3, after 3",
			before: []Temp{{ID: 1}, {ID: 2}, {ID: 3}},
			after:  []Temp{{ID: 1}, {ID: 2}, {ID: 3}},
			want:   []Temp{},
		},
		{
			memo:   "When after is empty, return all values",
			before: []Temp{{ID: 1}, {ID: 2}, {ID: 3}},
			after:  []Temp{},
			want:   []Temp{{ID: 1}, {ID: 2}, {ID: 3}},
		},
	}
	for _, suite := range suites {
		t.Run(suite.memo, func(t *testing.T) {
			got := FindRemoved(suite.before, suite.after)
			if len(got) != len(suite.want) {
				t.Errorf("want: %v, but got: %v", suite.want, got)
			}
			for i, v := range got {
				if v.ID != suite.want[i].ID {
					t.Errorf("want: %v, but got: %v", suite.want, got)
				}
			}
		})
	}
}
func TestUniq(t *testing.T) {
	type Suite struct {
		memo string
		arr  []Temp
		want []Temp
	}
	suites := []Suite{
		{
			memo: "arr 3, after 2",
			arr:  []Temp{{ID: 1}, {ID: 1}, {ID: 1}, {ID: 2}, {ID: 3}},
			want: []Temp{{ID: 1}, {ID: 2}, {ID: 3}},
		},
		{
			memo: "arr 3, after 3",
			arr:  []Temp{{ID: 1}, {ID: 1}, {ID: 1}},
			want: []Temp{{ID: 1}},
		},
		{
			memo: "When after is empty, return all values",
			arr:  []Temp{{ID: 1}, {ID: 2}, {ID: 2}, {ID: 2}, {ID: 3}, {ID: 1}},
			want: []Temp{{ID: 1}, {ID: 2}, {ID: 3}},
		},
	}
	for _, suite := range suites {
		t.Run(suite.memo, func(t *testing.T) {
			got := Uniq(suite.arr)
			if len(got) != len(suite.want) {
				t.Errorf("want: %v, but got: %v", suite.want, got)
			}
			for i, v := range got {
				if v.ID != suite.want[i].ID {
					t.Errorf("want: %v, but got: %v", suite.want, got)
				}
			}
		})
	}
}
func TestIntersect(t *testing.T) {
	type Suite struct {
		memo string
		A    []uint
		B    []uint
		want []uint
	}
	suites := []Suite{
		{
			memo: "1",
			A:    []uint{1, 2, 3},
			B:    []uint{1, 2},
			want: []uint{1, 2},
		},
		{
			memo: "2",
			A:    []uint{1, 2, 3},
			B:    []uint{1, 2, 3},
			want: []uint{1, 2, 3},
		},
		{
			memo: "3",
			A:    []uint{1, 2, 3},
			B:    []uint{},
			want: []uint{},
		},
	}
	for _, suite := range suites {
		t.Run(suite.memo, func(t *testing.T) {
			got := Intersect(suite.A, suite.B)
			if len(got) != len(suite.want) {
				t.Errorf("want: %v, but got: %v", suite.want, got)
			}
			for i, v := range got {
				if v != suite.want[i] {
					t.Errorf("want: %v, but got: %v", suite.want, got)
				}
			}
		})
	}
}
