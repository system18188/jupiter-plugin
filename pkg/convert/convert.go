package convert

import (
	"bytes"
	"errors"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"net/http"
	"strconv"
)

/*
封装类型转换相关操作
创建于2017年04月22日
*/

/*
[]byte -> string
@param b: 待转换参数
@return string
*/
func ByteListToString(b []byte) string {
	return string(b)
}

/*
float64-> string
@param n: 待转换参数
@return string
*/
func Float64ToString(n float64) string {
	return strconv.FormatFloat(n, 'f', -1, 64)
}

/*
int -> string
@param n: 待转换参数
@return string
*/
func IntToString(n int) string {
	return strconv.Itoa(n)
}

/*
int64 -> string
@param n: 待转换参数
@return string
*/
func Int64ToString(n int64) string {
	return strconv.FormatInt(n, 10)
}

/*
string -> int
@param s: 待转换参数
@return int, error
*/
func StringToInt(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.New("参数类型不正确！")
	}
	return n, nil
}

/*
string -> int
@param s: 待转换参数
@return int, default 0
*/
func StringToIntValue(s string) int {
	if n, err := StringToInt(s); err == nil {
		return n
	}
	return 0
}

/*
string -> int64
@param s: 待转换参数
@return int64, error
*/
func StringToInt64(s string) (int64, error) {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, errors.New("参数类型不正确！")
	}
	return n, nil
}

/*
string -> int64
@param s: 待转换参数
@return int64, default 0
*/
func StringToInt64Value(s string) int64 {
	if n, err := StringToInt64(s); err == nil {
		return n
	}
	return 0
}

func StringToUint32(s string) uint32 {
	if n, err := strconv.ParseUint(s, 10, 32); err == nil {
		return uint32(n)
	}
	return 0
}

func StringToUint64(s string) uint64 {
	if n, err := strconv.ParseUint(s, 10, 64); err == nil {
		return n
	}
	return 0
}

/*
string -> bool
@param s: 待转换参数
@return bool, default false
*/
func StringToBool(s string) bool {
	if n, err := strconv.ParseBool(s); err == nil {
		return n
	}
	return false
}

/*
string -> float32
@param s: 待转换参数
@return float32, error
*/
func StringToFloat32(s string) (float32, error) {
	n, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, errors.New("参数类型不正确！")
	}
	return float32(n), nil
}

/*
string -> float32
@param s: 待转换参数
@return float32, default 0
*/
func StringToFloat32Value(s string) float32 {
	if n, err := StringToFloat32(s); err == nil {
		return n
	}
	return 0
}

/*
string -> float64
@param s: 待转换参数
@return float64, error
*/
func StringToFloat64(s string) (float64, error) {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0, errors.New("参数类型不正确！")
	}
	return n, nil
}

/*
string -> float64
@param s: 待转换参数
@return float64, default 0.0
*/
func StringToFloat64Value(s string) float64 {
	if n, err := StringToFloat64(s); err == nil {
		return n
	}
	return 0.0
}

/*
float64 -> int
@param n: 待转换参数
@return int
*/
func Float64ToInt(n float64) int {
	return int(n)
}

/*
float64 -> int64
@param n: 待转换参数
@return int
*/
func Float64ToInt64(n float64) int64 {
	return int64(n)
}

/*
int64 -> int
@param n: 待转换参数
@return int
*/
func Int64ToInt(n int64) int {
	return int(n)
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func Utf8ToGbk(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

// 首字母大写其他转换成小写
func FirstCharToUpper(str string) string {
	strByte := []byte(str)
	strRes := FirstCharToUpperBytes(strByte)
	return string(strRes)
}

// 首字母大写其他转换成小写([]byte模式)
func FirstCharToUpperBytes(strByte []byte) []byte {
	if len(strByte) > 0 {
		firstChar := bytes.ToUpper(strByte[:1])
		otherChar := bytes.ToLower(strByte[1:])
		strBytes := bytes.NewBuffer(firstChar)
		strBytes.Write(otherChar)
		return strBytes.Bytes()
	} else {
		emptyByte := make([]byte, 0)
		return bytes.NewBuffer(emptyByte).Bytes()
	}
}

// 取得所以表单的键名
func FromKeyToSlice(req *http.Request) []string {
	slice := make([]string, 0)
	if err := req.ParseForm(); err != nil {
		return slice
	}
	for k := range req.PostForm {
		slice = append(slice, k)
	}
	return slice
}

func InterfaceToString(v interface{}) (string, error) {
	s, ok := v.(string)
	if !ok {
		return "", errors.New("interface to string error")
	}
	return s, nil
}
