package thinkutils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

type md5utils struct {
}

func (this md5utils) MD5String(szTxt string) string {
	h := md5.New()
	h.Write([]byte(szTxt))

	return strings.ToLower(hex.EncodeToString(h.Sum(nil)))
}

func (this md5utils) MD5File(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	body, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	md5 := fmt.Sprintf("%x", md5.Sum(body))
	runtime.GC()

	return strings.ToLower(strings.TrimSpace(md5)), nil
}
