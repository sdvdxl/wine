// util
package util

import (
	"crypto/md5"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"regexp"

	log "github.com/sdvdxl/log4go"
)

type JsonResult struct {
	Success bool
	Msg     string
	Data    interface{}
}

var (
	StringBlankPattern *regexp.Regexp
)

func init() {
	StringBlankPattern = regexp.MustCompile("[ \t\n]")
}

//panic ，如果err不是nil，则panic
func PanicError(err error) {
	if err != nil {
		panic(err)
	}
}

func LogError(err error, msg ...string) {
	if err != nil {
		log.Error(err, msg)
	}
}

//将传入的字符串md5后输出
func MD5(text string) string {
	result := md5.Sum([]byte(text))
	return hex.EncodeToString(result[:])
}

//MD5并加盐
func MD5AndSalt(text, salt string) string {
	return MD5(MD5(text) + MD5(salt) + salt)
}

func IsBlank(originStr string) bool {
	if len(StringBlankPattern.ReplaceAllString(originStr, "")) == 0 {
		return true
	}

	return false
}

func SHA512(text string) string {
	result := sha512.Sum512([]byte(text))
	return hex.EncodeToString(result[:])
}

func HashAndSalt(text, salt string) string {
	return SHA512(SHA512(MD5(text)) + salt)
}

func HashPid(phone, email string, timestamp int64) string {
	return SHA512(SHA512(SHA512(MD5(phone))+email) + fmt.Sprintf("%v", timestamp))
}

func GetCaptchaString(bytes []byte) string {
	runes := []rune(hex.EncodeToString(bytes))
	result := ""
	for i := 1; i < len(runes); i += 2 {
		result += string(runes[i])
	}
	return result
}
