package thinkutils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
)

type httputils struct {
}

func (this httputils) Get(szUrl string) (string, error) {
	resp, err := http.Get(szUrl)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

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
	defer func() { _ = resp.Body.Close() }()

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
	defer func() { _ = resp.Body.Close() }()

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
	defer func() { _ = resp.Body.Close() }()

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
	defer func() { _ = resp.Body.Close() }()

	byteBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return StringUtils.BytesToString(byteBody), nil
}

func (this httputils) DownloadFile(szFile, szUrl string) error {
	// Create the file
	out, err := os.Create(szFile)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(szUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
