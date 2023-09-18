package util

import "testing"

const (
	Num = 100
)

func TestGetSonyflakeID(t *testing.T) {
	var ids [Num]uint64
	var err error
	for i := 0; i < Num; i++ {
		ids[i], err = GetSonyFlakeID()
		if err != nil {
			t.Fatal(err)
		}
		if i > 0 && ids[i] <= ids[i-1] {
			t.Error("no increasing")
		}
	}
}
