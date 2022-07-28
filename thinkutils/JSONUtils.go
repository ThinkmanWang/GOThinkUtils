package thinkutils

import "encoding/json"

type jsonutils struct {
}

func (this jsonutils) ToJson(v interface{}) string {
	byteJson, err := json.Marshal(v)
	if err != nil {
		return ""
	}

	return string(byteJson)
}

func (this jsonutils) FromJson(szJson string, v interface{}) error {
	return json.Unmarshal(StringUtils.StringToBytes(szJson), v)
}

func (this jsonutils) IsJSONString(s string) bool {
	var js string
	return json.Unmarshal([]byte(s), &js) == nil
}