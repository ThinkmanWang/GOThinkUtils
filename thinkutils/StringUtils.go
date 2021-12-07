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

func (this stringutils) StringToBytes(szTxt string) []byte {
	return []byte(szTxt)
}

func (this stringutils) BytesToString(byteTxt []byte) string {
	return string(byteTxt)
}
