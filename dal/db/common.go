package db

import (
	"strconv"
	"strings"
)

type void struct{}

var empty void

// SetToString 将 Set 集合转化为 SQL 中的查询条件，如果为空则返回 NULL
func SetToString(intSet map[uint64]struct{}) string {
	if len(intSet) == 0 {
		return "NULL"
	}
	var s strings.Builder
	for k := range intSet {
		s.WriteString(strconv.Itoa(int(k)) + ",")
	}
	return strings.Trim(s.String(), ",")
}
