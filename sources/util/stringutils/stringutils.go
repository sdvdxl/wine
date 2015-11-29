package stringutils

func ToString(obj interface{}) string {
	if nil == obj {
		return ""
	}

	var result string
	switch obj.(type) {
	case string:
		result, _ = obj.(string)
	}

	return result
}
