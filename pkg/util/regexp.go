package util

import (
	"regexp"
)

// CheckPasswordLever 用于校验密码强度
// 密码必须包括字⺟⼤⼩写、数字、特殊符号，长度6-32位
func CheckPasswordLever(ps string) bool {
	if len(ps) < 5 || len(ps) > 32 {
		return false
	}
	num := `[0-9]{1}`
	lowercase := `[a-z]{1}`
	uppercase := `[A-Z]{1}`
	symbol := `[!@#~$%^&*()+|_]{1}`
	if b, err := regexp.MatchString(num, ps); !b || err != nil {
		return false
	}
	if b, err := regexp.MatchString(lowercase, ps); !b || err != nil {
		return false
	}
	if b, err := regexp.MatchString(uppercase, ps); !b || err != nil {
		return false
	}
	if b, err := regexp.MatchString(symbol, ps); !b || err != nil {
		return false
	}
	return true
}
