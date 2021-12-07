package thinkutils

import (
	"math/rand"
)

type randutils struct {
}

func (this randutils) RandInt(nMin int, nMax int) int {
	return rand.Intn(nMax-nMin) + nMin
}

func (this randutils) RandPasssword(nLen int) string {
	szBase := "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789"
	szRet := ""
	for i := 0; i < nLen; i++ {
		ch := szBase[this.RandInt(0, len(szBase))]
		szRet += string(ch)
	}

	return szRet
}

func (this randutils) UUID() string {
	return UUIDUtils.New()
}
