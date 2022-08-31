package thinkutils

import (
	"github.com/ThinkmanWang/GOThinkUtils/thinkutils/reflectx"
	"reflect"
)

type structutis struct {
}

func (this structutis) FieldAddrByTag(pData interface{}, szTag, szVal string) (interface{}, bool) {
	if nil == pData || StringUtils.IsEmpty(szTag) || StringUtils.IsEmpty(szVal) {
		return nil, false
	}

	v := reflect.Indirect(reflect.ValueOf(pData))

	defaultMapper := reflectx.NewMapper(szTag)
	m := defaultMapper.FieldMap(v)

	for key, val := range m {
		//log.Info("%s", key)

		if key == szVal {
			return val.Addr().Interface(), true
		}
	}

	return "", false
}

//func (this structutis) FieldAddrByTag(pData interface{}, szKey, szVal string) (interface{}, bool) {
//	if nil == pData || StringUtils.IsEmpty(szKey) || StringUtils.IsEmpty(szVal) {
//		return nil, false
//	}
//
//	szFieldName, bFound := this.FieldNameByTag(pData, szKey, szVal)
//	if false == bFound {
//		return nil, false
//	}
//
//	return reflect.ValueOf(pData).Elem().FieldByName(szFieldName).Addr().Interface(), true
//}

//func (this structutis) FieldAddrByTag(pData interface{}, szKey, szVal string) (interface{}, bool) {
//	if nil == pData || StringUtils.IsEmpty(szKey) || StringUtils.IsEmpty(szVal) {
//		return nil, false
//	}
//
//	structVal := reflect.ValueOf(pData)
//
//	szFieldName, bFound := this.FieldNameByTag(pData, szKey, szVal)
//	if false == bFound {
//		return nil, false
//	}
//
//	fieldVal := structVal.FieldByName(szFieldName)
//	if !fieldVal.IsValid() {
//		return nil, false
//	}
//
//	return fieldVal.Addr().Interface(), true
//}
