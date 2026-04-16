package thinkutils

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

// 全局HTTP客户端（带连接池）：懒加载+单例+线程安全
var (
	globalHTTPClient *http.Client
	clientOnce       sync.Once // 保证连接池只初始化一次
	mu               sync.Mutex
)

// 获取全局带连接池的HTTP客户端（懒加载，线程安全）
func getGlobalHTTPClient() *http.Client {
	// 双重检查锁：兼顾性能和线程安全
	if globalHTTPClient == nil {
		mu.Lock()
		defer mu.Unlock()
		if globalHTTPClient == nil {
			clientOnce.Do(func() {
				// 高并发连接池核心配置
				transport := &http.Transport{
					// 连接池容量配置（适配数万并发）
					MaxIdleConns:          10000,            // 全局最大空闲连接数
					MaxIdleConnsPerHost:   1000,             // 每个Host最大空闲连接（复用核心）
					MaxConnsPerHost:       5000,             // 每个Host最大并发连接数
					IdleConnTimeout:       90 * time.Second, // 空闲连接超时释放
					ResponseHeaderTimeout: 30 * time.Second, // 响应头读取超时
					TLSHandshakeTimeout:   10 * time.Second, // TLS握手超时
					DisableCompression:    false,            // 启用压缩（提升传输效率）
					DisableKeepAlives:     false,            // 开启Keep-Alive（必须，否则连接池失效）

					// 拨号层配置（TCP连接基础）
					DialContext: (&net.Dialer{
						Timeout:   30 * time.Second, // 拨号超时
						KeepAlive: 30 * time.Second, // TCP层Keep-Alive
						DualStack: true,             // 支持IPv4/IPv6
					}).DialContext,
				}

				// 全局客户端：复用连接池
				globalHTTPClient = &http.Client{
					Transport: transport,
					Timeout:   60 * time.Second, // 客户端整体超时（覆盖所有阶段）
				}
			})
		}
	}
	return globalHTTPClient
}

type httputils struct {
}

func (this httputils) Get(szUrl string) (string, error) {
	// 替换为全局连接池客户端
	req, err := http.NewRequest("GET", szUrl, nil)
	if err != nil {
		return "", err
	}

	resp, err := getGlobalHTTPClient().Do(req)
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

	// 替换为全局连接池客户端
	req, err := http.NewRequest("POST", szUrl, bytes.NewBufferString(url.Values(form).Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	resp, err := getGlobalHTTPClient().Do(req)
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
	// 替换原有临时client为全局连接池客户端
	resp, err := getGlobalHTTPClient().Do(request)
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

	// 替换原有临时client为全局连接池客户端
	resp, err := getGlobalHTTPClient().Do(request)
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

	// 替换原有临时client为全局连接池客户端
	resp, err := getGlobalHTTPClient().Do(request)
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

	// 替换为全局连接池客户端
	req, err := http.NewRequest("GET", szUrl, nil)
	if err != nil {
		return err
	}
	resp, err := getGlobalHTTPClient().Do(req)
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
