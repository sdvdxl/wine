package json

import (
	j "encoding/json"
	. "github.com/sdvdxl/wine/sources/util/log"
)

//解析对象为json字符串，如果obj为nil，那么返回 空字符串
func ToJSONString(obj interface{}) string {
	if nil == obj {
		return ""
	}

	bytes, err := j.Marshal(obj)
	if err != nil {
		Logger.Error("error occured when encode %v", obj)
		Logger.Error(err)
		return ""

	}

	return string(bytes)
}

//将json字符串转换为对象，如果字符串为空，则obj为nil，err也是nil
//如果不是空，则调用json.Unmarshal并返回
func ToJSONObject(jsonString string, obj interface{}) error {
	if jsonString == "" {
		obj = nil
		return nil
	}

	bytes := []byte(jsonString)
	return j.Unmarshal(bytes, obj)
}
