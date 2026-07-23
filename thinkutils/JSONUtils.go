package thinkutils

import (
	jsoniter "github.com/json-iterator/go"
)

// json 采用 jsoniter 的"标准库兼容"配置：行为与 encoding/json 完全一致
// (含 HTML 转义、json.Marshaler/Unmarshaler 接口)，但反序列化性能通常快 2~3 倍。
// 该库已在依赖树中(由 gin 间接引入)，切换不增加任何新依赖。
var json = jsoniter.ConfigCompatibleWithStandardLibrary

type jsonutils struct {
}

func (this jsonutils) ToJson(v interface{}) string {
	byteJson, err := json.Marshal(v)
	if err != nil {
		return ""
	}

	return string(byteJson)
}

// ToJsonBytes 直接返回 []byte，避免调用方再做 string->[]byte 的多余拷贝。
// 适用于写入 BigCache 等只需要字节的场景。
func (this jsonutils) ToJsonBytes(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (this jsonutils) FromJson(szJson string, v interface{}) error {
	return json.Unmarshal(StringUtils.StringToBytes(szJson), v)
}

// FromJsonBytes 直接从 []byte 反序列化，避免 []byte->string->[]byte 的来回拷贝。
// 适用于从 BigCache 读取后直接解码的热路径。
func (this jsonutils) FromJsonBytes(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (this jsonutils) IsJSONString(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func (this jsonutils) TrimJSON(szJson string) string {
	var js map[string]interface{}
	err := json.Unmarshal([]byte(szJson), &js)
	if err != nil {
		return szJson
	}

	return this.ToJson(js)
}