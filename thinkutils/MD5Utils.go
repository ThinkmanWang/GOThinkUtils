package thinkutils

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

type md5utils struct {
}

func (this md5utils) MD5String(szTxt string) string {
	h := md5.New()
	h.Write([]byte(szTxt))

	return strings.ToLower(hex.EncodeToString(h.Sum(nil)))
}
