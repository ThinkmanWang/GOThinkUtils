package thinkutils

import "github.com/ThinkmanWang/GOThinkUtils/thinkutils/logger"

var (
	log *logger.LocalLogger = logger.DefaultLogger()

	DateTime     datetime
	StringUtils  stringutils
	RandUtils    randutils
	MD5Utils     md5utils
	IPUtils      iputils
	UUIDUtils    uuidutils
	ThinkMysql   thinkmysql
	ThinkRedis   thinkredis
	FileUtils    fileutils
	RegularUtils regulartils
	JSONUtils    jsonutils
	KafkaUtils   kafkautils
	UDPUtils     udputils
	Base64Utils  base64utils
	HttpUtils    httputils
	SetUtils     setutils
	StructUtils  structutis
)

func ListRemoveAt[T int8 | int16 | int32 | int64 | int | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64 | string | any](lstData []T, nIndex int) []T {
	if nil == lstData {
		return nil
	}

	if nIndex < 0 || nIndex > len(lstData) {
		return nil
	}

	return append(lstData[:nIndex], lstData[nIndex+1:]...)
}

func MinVal[T int8 | int16 | int32 | int64 | int | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](lstVal ...T) T {
	nRet := lstVal[0]
	for _, nVal := range lstVal {
		if nVal < nRet {
			nRet = nVal
		}
	}

	return nRet
}

func MaxVal[T int8 | int16 | int32 | int64 | int | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64](lstVal ...T) T {
	nRet := lstVal[0]
	for _, nVal := range lstVal {
		if nVal > nRet {
			nRet = nVal
		}
	}

	return nRet
}
