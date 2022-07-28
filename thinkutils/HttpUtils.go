package thinkutils

import (
	"bytes"
	"io"
	"net/http"
)

type httputils struct {
}

func (this httputils) Get(szUrl string) (string, error) {
	resp, err := http.Get(szUrl)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close()}()

	byteBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return StringUtils.BytesToString(byteBody), nil
}

func (this httputils) PostForm(szUrl string, data map[string]string) (string, error) {
	form := make(map[string][]string)
	for k, v := range data {
		szVals := []string{v}
		form[k] = szVals
	}

	resp, err := http.PostForm(szUrl, form)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close()}()

	byteBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return StringUtils.BytesToString(byteBody), nil
}

func (this httputils) PostJSON(szUrl, szJson string) (string, error) {
	request, err := http.NewRequest("POST", szUrl, bytes.NewBuffer(StringUtils.StringToBytes(szJson)))
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close()}()

	byteBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return StringUtils.BytesToString(byteBody), nil
}

func (this httputils) GetWithHeader(szUrl string, mapHeader map[string]string) (string, error) {

	request, err := http.NewRequest("GET", szUrl, nil)
	if err != nil {
		return "", err
	}

	for k, v := range mapHeader {
		request.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close()}()

	byteBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return StringUtils.BytesToString(byteBody), nil
}

func (this httputils) PostJSONWithHeader(szUrl string, mapHeader map[string]string, szJson string) (string, error) {
	request, err := http.NewRequest("POST", szUrl, bytes.NewBuffer(StringUtils.StringToBytes(szJson)))
	if err != nil {
		return "", err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	for k, v := range mapHeader {
		request.Header.Set(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close()}()

	byteBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return StringUtils.BytesToString(byteBody), nil
}