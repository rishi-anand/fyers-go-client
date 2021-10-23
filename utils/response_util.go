package utils

import "github.com/tidwall/gjson"

const (
	statusMsgKey    = "s"
	successMsgValue = "ok"
	failedMsgValue  = "error"
)

func IsSuccessResponse(in []byte) bool {
	if len(in) == 0 {
		return false
	}
	return GetJsonValueAtPath(in, statusMsgKey) == successMsgValue
}

func IsFailedResponse(in []byte) bool {
	return GetJsonValueAtPath(in, statusMsgKey) == failedMsgValue
}

func GetJsonValueAtPath(in []byte, path string) string {
	return gjson.Get(string(in), path).String()
}
