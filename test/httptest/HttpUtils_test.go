package testing

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/ThinkmanWang/GOThinkUtils/thinkutils"
)

// 初始化被测试的 httputils 实例
var httpUtil = thinkutils.HttpUtils

// -------------------------- 测试 Get 方法 --------------------------
func TestHttputils_Get(t *testing.T) {
	// 测试用例：带 GET 参数的 postman-echo 接口
	testUrl := "https://postman-echo.com/get?name=test&age=20"

	// 调用被测试的 Get 方法
	resp, err := httpUtil.Get(testUrl)
	// 1. 验证是否有错误
	if err != nil {
		// Fatalf 会终止当前测试，因为连请求都失败了，后续验证无意义
		t.Fatalf("Get 方法调用失败: %v", err)
	}

	// 2. 验证响应非空
	if resp == "" {
		t.Error("Get 响应为空，不符合预期")
	}

	// 3. 进阶：解析响应 JSON，验证参数是否正确传递（精准断言）
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(resp), &result); err != nil {
		t.Fatalf("解析 Get 响应失败: %v", err)
	}
	// 验证 GET 参数是否正确
	args := result["args"].(map[string]interface{})
	if args["name"] != "test" {
		t.Errorf("Get 参数 name 预期为 test，实际为 %v", args["name"])
	}
	if args["age"] != "20" {
		t.Errorf("Get 参数 age 预期为 20，实际为 %v", args["age"])
	}

	t.Logf("Get 测试通过，响应：%s", resp)
}

// -------------------------- 测试 PostForm 方法 --------------------------
func TestHttputils_PostForm(t *testing.T) {
	testUrl := "https://postman-echo.com/post"
	// 构造表单参数
	formData := map[string]string{
		"username": "test_user",
		"password": "123456",
		"gender":   "male",
	}

	// 调用 PostForm 方法
	resp, err := httpUtil.PostForm(testUrl, formData)
	if err != nil {
		t.Fatalf("PostForm 调用失败: %v", err)
	}

	// 解析响应验证参数
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(resp), &result); err != nil {
		t.Fatalf("解析 PostForm 响应失败: %v", err)
	}
	form := result["form"].(map[string]interface{})
	// 逐个验证表单参数
	if form["username"] != "test_user" {
		t.Errorf("PostForm 参数 username 预期为 test_user，实际为 %v", form["username"])
	}
	if form["password"] != "123456" {
		t.Errorf("PostForm 参数 password 预期为 123456，实际为 %v", form["password"])
	}

	t.Logf("PostForm 测试通过，响应：%s", resp)
}

// -------------------------- 测试 PostJSON 方法 --------------------------
func TestHttputils_PostJSON(t *testing.T) {
	testUrl := "https://postman-echo.com/post"
	// 构造 JSON 字符串
	jsonStr := `{
		"name": "张三",
		"age": 25,
		"hobbies": ["编程", "阅读"],
		"isVip": true
	}`

	// 调用 PostJSON 方法
	resp, err := httpUtil.PostJSON(testUrl, jsonStr)
	if err != nil {
		t.Fatalf("PostJSON 调用失败: %v", err)
	}

	// 解析响应验证
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(resp), &result); err != nil {
		t.Fatalf("解析 PostJSON 响应失败: %v", err)
	}
	// postman-echo 会把 JSON 数据放在 data 字段
	data := result["data"].(map[string]interface{})
	if data["name"] != "张三" {
		t.Errorf("PostJSON 参数 name 预期为 张三，实际为 %v", data["name"])
	}
	if int(data["age"].(float64)) != 25 { // JSON 数字默认解析为 float64
		t.Errorf("PostJSON 参数 age 预期为 25，实际为 %v", data["age"])
	}

	t.Logf("PostJSON 测试通过，响应：%s", resp)
}

// -------------------------- 测试 GetWithHeader 方法 --------------------------
func TestHttputils_GetWithHeader(t *testing.T) {
	testUrl := "https://postman-echo.com/get"
	// 构造自定义请求头
	headers := map[string]string{
		"Token":      "test_token_123",
		"User-Agent": "Go-Testing/1.0",
	}

	// 调用 GetWithHeader 方法
	resp, err := httpUtil.GetWithHeader(testUrl, headers)
	if err != nil {
		t.Fatalf("GetWithHeader 调用失败: %v", err)
	}

	// 验证请求头是否传递成功
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(resp), &result); err != nil {
		t.Fatalf("解析 GetWithHeader 响应失败: %v", err)
	}
	headersResp := result["headers"].(map[string]interface{})
	// 注意：http 头会被转为小写，所以验证时用小写 key
	if headersResp["token"] != "test_token_123" {
		t.Errorf("GetWithHeader 头 Token 预期为 test_token_123，实际为 %v", headersResp["token"])
	}

	t.Logf("GetWithHeader 测试通过，响应：%s", resp)
}

// -------------------------- 测试 PostJSONWithHeader 方法 --------------------------
func TestHttputils_PostJSONWithHeader(t *testing.T) {
	testUrl := "https://postman-echo.com/post"
	headers := map[string]string{
		"Token":       "post_json_token_456",
		"App-Version": "1.0.0",
	}
	jsonStr := `{"orderId": "123456", "amount": 99.9}`

	// 调用方法
	resp, err := httpUtil.PostJSONWithHeader(testUrl, headers, jsonStr)
	if err != nil {
		t.Fatalf("PostJSONWithHeader 调用失败: %v", err)
	}

	// 验证头和参数
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(resp), &result); err != nil {
		t.Fatalf("解析 PostJSONWithHeader 响应失败: %v", err)
	}
	// 验证请求头
	headersResp := result["headers"].(map[string]interface{})
	if headersResp["token"] != "post_json_token_456" {
		t.Errorf("PostJSONWithHeader 头 Token 预期为 post_json_token_456，实际为 %v", headersResp["token"])
	}
	// 验证 JSON 参数
	data := result["data"].(map[string]interface{})
	if data["orderId"] != "123456" {
		t.Errorf("PostJSONWithHeader 参数 orderId 预期为 123456，实际为 %v", data["orderId"])
	}

	t.Logf("PostJSONWithHeader 测试通过，响应：%s", resp)
}

// -------------------------- 测试 DownloadFile 方法 --------------------------
func TestHttputils_DownloadFile(t *testing.T) {
	// 测试地址：postman-echo 的测试文本文件
	testUrl := "https://cdn.zllinks.com/iaa/sdk/SDK%E6%9B%B4%E6%96%B0_ys_2026/XIAOMI/2025-05-27_222.159u2/z_ys_SDKUtils_xiaomi-release_1.0.1.aar"
	// 临时文件路径（测试完成后可删除）
	testFile := "./z_ys_SDKUtils_xiaomi-release_1.0.1.aar"

	// 调用下载方法
	err := httpUtil.DownloadFile(testFile, testUrl)
	if err != nil {
		t.Fatalf("DownloadFile 调用失败: %v", err)
	}

	// 验证文件是否存在
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Error("DownloadFile 下载的文件不存在，不符合预期")
	}

	// 清理测试文件
	defer func() {
		if err := os.Remove(testFile); err != nil {
			t.Logf("清理测试文件失败: %v", err)
		}
	}()

	t.Log("DownloadFile 测试通过")
}
