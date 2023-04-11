// Package util
package util

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"
	"time"
)

// util tool 基础函数工具集合

func Map2JSON[K comparable, V any](m map[K]V) (string, error) {
	mapJSON, err := Map2Bytes(m)
	if err != nil {
		return "", err
	}
	return string(mapJSON), nil
}

func Interface2JSON(i interface{}) (string, error) {
	interfaceJSON, err := json.Marshal(i)
	if err != nil {
		return "", err
	}
	return string(interfaceJSON), nil
}

func Map2Bytes[K comparable, V any](m map[K]V) ([]byte, error) {
	mapJSON, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return mapJSON, nil
}

type intType interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

func Str2Int[T intType](s string) (T, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return T(i), nil
}

func IntToStr[T intType](i T) string {
	return strconv.FormatInt(int64(i), 10)
}

func Str2IntArr[T intType](s string, sep string) ([]T, error) {
	arr := strings.Split(s, sep)
	arrLen := len(arr)
	if arrLen <= 0 {
		return nil, errors.New("have no data")
	}
	intArr := make([]T, arrLen)
	var err error
	for i, v := range arr {
		intArr[i], err = Str2Int[T](v)
		if err != nil {
			return nil, err
		}
	}
	return intArr, nil
}

type uintType interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func Str2Uint[T uintType](s string) (T, error) {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return T(i), nil
}

func UintToStr[T uintType](i T) string {
	return strconv.FormatUint(uint64(i), 10)
}

func Str2FloatArr[T floatType](s string, sep string) ([]T, error) {
	arr := strings.Split(s, sep)
	arrLen := len(arr)
	if arrLen <= 0 {
		return nil, errors.New("have no data")
	}
	float64Arr := make([]T, arrLen)
	var err error
	for i, v := range arr {
		float64Arr[i], err = Str2Float[T](v)
		if err != nil {
			return nil, err
		}
	}
	return float64Arr, nil
}

func Str2Bool(s string) (bool, error) {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return false, err
	}

	return b, nil
}

type floatType interface {
	~float32 | ~float64
}

func Str2Float[T floatType](s string) (T, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}

	return T(f), nil
}

func FloatToStr[T floatType](f T) string {
	return strconv.FormatFloat(float64(f), 'f', -1, 64)
}

func BoolToStr(b bool) string {
	return strconv.FormatBool(b)
}

func StartWith(s string, sub string, ignoreCase bool) (found bool) {
	if len(s) < len(sub) {
		return false
	}
	if ignoreCase {
		return strings.ToUpper(s[:len(sub)]) == strings.ToUpper(sub)
	}
	return s[:len(sub)] == sub
}

func EndWith(s string, sub string, ignoreCase bool) (found bool) {
	if len(s) < len(sub) {
		return false
	}
	startIndex := len(s) - len(sub)
	if ignoreCase {
		return strings.ToUpper(s[startIndex:]) == strings.ToUpper(sub)
	}
	return s[startIndex:] == sub
}

// ElemIn 判断数据是否在数组中
func ElemIn[T comparable](target T, arr []T) bool {
	for _, e := range arr {
		if target == e {
			return true
		}
	}
	return false
}

// Intersection 交集
func Intersection[T comparable](arr1 []T, arr2 []T) []T {
	len1 := len(arr1)
	len2 := len(arr2)
	if len1 <= 0 || len2 <= 0 {
		return []T{}
	}

	minLen := len1
	if minLen > len2 {
		minLen = len2
	}

	map1 := map[T]struct{}{}
	for _, v := range arr1 {
		map1[v] = struct{}{}
	}

	if minLen > len(map1) {
		minLen = len(map1)
	}

	resultArr := make([]T, minLen)
	resultLen := 0
	for _, v := range arr2 {
		_, ok := map1[v]
		if ok {
			resultArr[resultLen] = v
			resultLen++
		}
	}

	if resultLen <= 0 {
		return []T{}
	}
	return resultArr[:resultLen]
}

// Union 并集
func Union[T comparable](arr1 []T, arr2 []T) []T {
	len1 := len(arr1)
	len2 := len(arr2)
	if len1 <= 0 && len2 <= 0 {
		return []T{}
	}

	allMap := map[T]struct{}{}
	for _, v := range arr1 {
		allMap[v] = struct{}{}
	}

	for _, v := range arr2 {
		allMap[v] = struct{}{}
	}

	maxLen := len(allMap)
	if maxLen <= 0 {
		return []T{}
	}

	resultArr := make([]T, maxLen)
	resultCount := 0
	for k := range allMap {
		resultArr[resultCount] = k
		resultCount++
	}

	return resultArr
}

// Difference 差集
func Difference[T comparable](arr1 []T, arr2 []T) []T {
	len1 := len(arr1)
	if len1 <= 0 {
		return []T{}
	}

	len2 := len(arr2)
	if len2 <= 0 {
		return arr1
	}

	map2 := map[T]struct{}{}
	for _, v := range arr2 {
		map2[v] = struct{}{}
	}

	resultArr := make([]T, len1)
	resultCount := 0
	map1 := map[T]struct{}{}
	for _, v := range arr1 {
		_, ok := map2[v]
		if !ok {
			_, ok = map1[v]
			if !ok {
				map1[v] = struct{}{}
				resultArr[resultCount] = v
				resultCount++
			}

		}
	}

	return resultArr[:resultCount]
}

// JoinArr 链接数组
func JoinArr[T any](arr []T, linkStr string) string {
	var ret strings.Builder
	if len(arr) >= 1 {
		ret.WriteString(InterfaceToString(arr[0]))
	}
	for i, lenArr := 1, len(arr); i < lenArr; i++ {
		ret.WriteString(linkStr)
		ret.WriteString(InterfaceToString(arr[i]))
	}
	return ret.String()
}

// GetMapValue 获取map数据
func GetMapValue[K comparable, V any](raw map[K]V, key K, def V) V {
	if raw == nil {
		return def
	}
	val, ok := raw[key]
	if ok {
		return val
	}
	return def
}

// InterfaceToString 任意数据转string
func InterfaceToString(v interface{}) (ret string) {
	defer func() {
		if p := recover(); p != nil {
			log.Println(debug.Stack())
			ret = ""
		}
	}()

	ret = ""
	if v == nil {
		return ret
	}
	switch v := v.(type) {
	case float64:
		ret = strconv.FormatFloat(v, 'f', -1, 64)
	case float32:
		ret = strconv.FormatFloat(float64(v), 'f', -1, 64)
	case int:
		ret = strconv.Itoa(v)
	case uint:
		ret = strconv.Itoa(int(v))
	case int8:
		ret = strconv.Itoa(int(v))
	case uint8:
		ret = strconv.Itoa(int(v))
	case int16:
		ret = strconv.Itoa(int(v))
	case uint16:
		ret = strconv.Itoa(int(v))
	case int32:
		ret = strconv.Itoa(int(v))
	case uint32:
		ret = strconv.Itoa(int(v))
	case int64:
		ret = strconv.FormatInt(v, 10)
	case uint64:
		ret = strconv.FormatUint(v, 10)
	case string:
		ret = v
	case []byte:
		ret = string(v)
	default:
		newValue, _ := json.Marshal(v)
		ret = string(newValue)
	}
	return ret
}

// MergeMap 合并map
func MergeMap[K comparable, V any](data ...map[K]V) map[K]V {
	if len(data) <= 0 {
		return nil
	}

	var raw map[K]V
	for _, e := range data {
		if e == nil {
			continue
		}

		if raw == nil {
			raw = map[K]V{}
		}

		for k, v := range e {
			raw[k] = v
		}
	}

	return raw
}

// Interface2Interface 数据结构转换
func Interface2Interface(data interface{}, obj interface{}) error {
	marshal, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = json.Unmarshal(marshal, obj)
	if err != nil {
		return err
	}
	return nil
}

// ParamsSorted 参数计算排序字符串
func ParamsSorted[V any](params map[string]V) string {
	if len(params) <= 0 {
		return ""
	}

	keys := make([]string, len(params))
	i := 0
	for k := range params {
		keys[i] = k
		i++
	}

	sort.Strings(keys)

	var ret strings.Builder
	for ii, k := range keys {
		if ii == 0 {
			ret.WriteString(k)
			ret.WriteByte('=')
			ret.WriteString(InterfaceToString(params[k]))
		} else {
			ret.WriteByte('&')
			ret.WriteString(k)
			ret.WriteByte('=')
			ret.WriteString(InterfaceToString(params[k]))
		}
	}

	return ret.String()
}

// ArrSplit 按照给定的子数组长度进行数组切分，返回二维数组
func ArrSplit[T any](arr []T, batchSize int) [][]T {
	s := len(arr)
	if s <= 0 {
		return nil
	}

	if batchSize <= 0 {
		return [][]T{arr}
	}

	splitSize := s / batchSize
	if s%batchSize != 0 {
		splitSize++
	}
	if splitSize <= 0 {
		return nil
	}

	var ret [][]T
	for i := 0; i < splitSize; i++ {
		if i == splitSize-1 {
			ret = append(ret, arr[i*batchSize:])
		} else {
			ret = append(ret, arr[i*batchSize:(i+1)*batchSize])
		}
	}
	return ret
}

type numType interface {
	~int | ~int8 | ~int16 | ~int32 | int64 | uint | ~uint8 | ~uint16 | ~uint32 | uint64 | ~float32 | ~float64
}

// NumArrTransfer 数字型的数组转换
func NumArrTransfer[T1 numType, T2 numType](arr []T1) []T2 {
	if len(arr) <= 0 {
		return nil
	}

	ret := make([]T2, 0, len(arr))
	for _, v := range arr {
		ret = append(ret, T2(v))
	}
	return ret
}

// RandomInt64 returns, as an int64, a non-negative pseudo-random number in the half-open interval [min,max)
func RandomInt64(min, max int64) int64 {
	if min < 0 {
		panic("min must gte 0")
	}
	if max <= 0 {
		panic("max must gt 0")
	}
	if min >= max {
		panic("max must gt min")
	}
	rand.Seed(time.Now().UnixNano())
	r := rand.Int63n(max - min)
	return r + min
}

// HumpFormatToUnderLine 驼峰转下滑线命名
func HumpFormatToUnderLine(s string) string {
	lowerStr := strings.ToLower(s)
	var index []int
	for i := 0; i < len(s); i++ {
		if lowerStr[i] != s[i] {
			index = append(index, i)
		}
	}

	var newStr strings.Builder
	prefixIndex := 0
	for _, i := range index {
		if newStr.Len() > 0 {
			newStr.WriteByte('_')
		}

		newStr.WriteString(lowerStr[prefixIndex:i])
		prefixIndex = i
	}

	newStr.WriteByte('_')
	newStr.WriteString(lowerStr[prefixIndex:])
	return strings.TrimLeft(newStr.String(), "_")
}

// MapHumpFormatToUnderLine map key 驼峰转下滑线命名
func MapHumpFormatToUnderLine[V any](raw map[string]V) map[string]V {
	newMap := map[string]V{}
	for k, v := range raw {
		newMap[HumpFormatToUnderLine(k)] = v
	}
	return newMap
}

// StrSplit 多分隔符切分字符串
// 返回拆分的字符串数组以及实际的分割标识数组
func StrSplit(s string, seq []string, trimSpace bool) ([]string, []string) {
	if len(s) <= 0 {
		return nil, seq
	}

	if len(seq) <= 0 {
		return []string{s}, nil
	}

	var seqArr []string
	var arr []string
	sLen := len(s)
	lastIndex := 0
	for i := 0; i < sLen; {
		tagSeq := ""
		for _, ss := range seq {
			l := len(ss)
			if i+l >= sLen {
				continue
			}

			if s[i:i+l] == ss {
				tagSeq = ss
				break
			}
		}

		if len(tagSeq) <= 0 {
			i++
		} else {
			seqArr = append(seqArr, tagSeq)
			if trimSpace {
				arr = append(arr, strings.TrimSpace(s[lastIndex:i]))
			} else {
				arr = append(arr, s[lastIndex:i])
			}

			i += len(tagSeq)
			lastIndex = i
		}
	}

	if trimSpace {
		arr = append(arr, strings.TrimSpace(s[lastIndex:]))
	} else {
		arr = append(arr, s[lastIndex:])
	}
	return arr, seqArr
}

func SafeGo(f func()) {
	go func() {
		defer func() {
			if p := recover(); p != nil {
				log.Printf("%v", p)
			}
		}()

		f()
	}()
}

type char interface {
	byte | rune
}

// CharMerge 合并字符串中相邻的多个指定的字符
func CharMerge[T char](s string, ch T) []T {
	if s == "" {
		return nil
	}

	str := make([]T, 1, len(s))
	str[0] = T(s[0])
	for i := 1; i < len(s); i++ {
		c := T(s[i])
		if c != ch || c != str[len(str)-1] {
			str = append(str, c)
		}
	}
	return str
}
