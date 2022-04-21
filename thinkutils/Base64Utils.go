package thinkutils

import (
	"bufio"
	"encoding/base64"
	"io/ioutil"
	"os"
)

type base64utils struct {
}

func (this base64utils) FileToBase64(szPath string) (string, error) {
	f, err := os.Open(szPath)
	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(f)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	szEncoded := base64.StdEncoding.EncodeToString(content)

	return szEncoded, nil
}

func (this base64utils) Base64ToFile(szBase64 string, szPath string) error {
	dec, err := base64.StdEncoding.DecodeString(szBase64)
	if err != nil {
		return err
	}

	f, err := os.Create(szPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(dec); err != nil {
		return err
	}

	if err := f.Sync(); err != nil {
		return err
	}

	return nil
}

func (this base64utils) Base64ToString(szBase64 string) (string, error) {
	if StringUtils.IsEmpty(szBase64) {
		return "", nil
	}

	rawDecodedText, err := base64.StdEncoding.DecodeString(szBase64)
	if err != nil {
		return "", err
	}

	return string(rawDecodedText), nil
}

func (this base64utils) StringToBase64(szTxt string) (string, error) {
	if StringUtils.IsEmpty(szTxt) {
		return "", nil
	}

	str := base64.StdEncoding.EncodeToString([]byte(szTxt))
	return str, nil
}
