package thinkutils

import (
	"reflect"
	"strings"
)

type structutis struct {
}

func (this structutis) FieldNameByTag(pData interface{}, szKey, szVal string) (string, bool) {

	if nil == pData || StringUtils.IsEmpty(szKey) || StringUtils.IsEmpty(szVal) {
		return "", false
	}

	v := reflect.ValueOf(pData)

	rt := reflect.Indirect(v).Type()
	if rt.Kind() != reflect.Struct {
		return "", false
	}

	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		v := strings.Split(f.Tag.Get(szKey), ",")[0] // use split to ignore tag "options" like omitempty, etc.
		if v == szVal {
			return f.Name, true
		}
	}

	return "", false
}

func (this structutis) FieldAddrByTag(pData interface{}, szKey, szVal string) (interface{}, bool) {
	if nil == pData || StringUtils.IsEmpty(szKey) || StringUtils.IsEmpty(szVal) {
		return nil, false
	}

	structVal := reflect.ValueOf(pData)

	szFieldName, bFound := this.FieldNameByTag(pData, szKey, szVal)
	if false == bFound {
		return nil, false
	}

	fieldVal := structVal.FieldByName(szFieldName)
	if !fieldVal.IsValid() {
		return nil, false
	}

	return fieldVal.Addr().Interface(), true
}
