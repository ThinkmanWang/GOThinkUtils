package thinkutils

import (
	"regexp"
)

type regulartils struct {
}

func (this regulartils) IsEmail(email string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z]\\w{5,10}@([0-9a-zA-Z]{3,5}\\.){1,3}[a-z]{3}$")
	return emailRegex.MatchString(email)
}

func (this regulartils) IsPhone(szPhone string) bool {
	emailRegex := regexp.MustCompile("^1[0-9][0-9]{9}$")
	return emailRegex.MatchString(szPhone)
}
