package number
import (
	"fmt"
	"strconv"
)

func DefaultInt(obj interface{}, defaultValue int) int {
	if returnValue, ok:=  obj.(int) ;ok{
		return returnValue
	} else {
		returnValue, err:= strconv.Atoi(fmt.Sprint(obj))
		if err!=nil {
			return defaultValue
		}

		return returnValue
	}
}

func DefaultInt64(obj interface{}, defaultValue int64) int64 {
	if returnValue, ok:=  obj.(int64) ;ok{
		return returnValue
	} else {
		returnValue, err:= strconv.Atoi(fmt.Sprint(obj))
		if err!=nil {
			return int64(defaultValue)
		}

		return int64(returnValue)
	}
}
