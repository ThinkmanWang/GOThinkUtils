package thinkutils

import "strings"

type stringutils struct {
}

func (this stringutils) IsEmpty(szTxt string) bool {
	if len(strings.TrimSpace(szTxt)) <= 0 {
		return true
	}

	return false
}

func (this stringutils) IsEmptyPtr(szTxt *string) bool {
	if nil == szTxt || len(strings.TrimSpace(*szTxt)) <= 0 {
		return true
	}

	return false
}
