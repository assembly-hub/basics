package util

import (
	"fmt"
	"testing"
	"time"
)

func TestMap2Json(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			t.Error()
		}
	}()

	m := map[string]interface{}{
		"key1": "1",
		"key2": "2",
		"key3": 123,
	}
	js, err := Map2JSON(m)
	if err != nil {
		panic(err)
	}

	fmt.Println(js)
}

func TestMap2JsonBytes(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			t.Error()
		}
	}()

	m := map[string]interface{}{
		"key1": "1",
		"key2": "2",
		"key3": 123,
	}
	js, err := Map2Bytes(m)
	if err != nil {
		panic(err)
	}

	fmt.Println(js)
}

func TestStartWith(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			t.Error()
		}
	}()

	s := "123456789Abc"
	sub1 := "123"
	sub2 := "123456787654321234565432123"
	sub3 := "123456789Abc"
	sub4 := "123456789ABC"

	if !StartWith(s, sub1, false) {
		panic("")
	}
	if StartWith(s, sub2, false) {
		panic("")
	}
	if !StartWith(s, sub3, false) {
		panic("")
	}
	if StartWith(s, sub4, false) {
		panic("")
	}

	if !StartWith(s, sub1, true) {
		panic("")
	}
	if StartWith(s, sub2, true) {
		panic("")
	}
	if !StartWith(s, sub3, true) {
		panic("")
	}
	if !StartWith(s, sub4, true) {
		panic("")
	}
}

func TestEndWith(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			t.Error()
		}
	}()

	s := "123456789Abc"
	sub1 := "abc"
	sub2 := "123456787654321234565432123"
	sub3 := "123456789Abc"
	sub4 := "123456789ABC"

	if EndWith(s, sub1, false) {
		panic("")
	}
	if EndWith(s, sub2, false) {
		panic("")
	}
	if !EndWith(s, sub3, false) {
		panic("")
	}
	if EndWith(s, sub4, false) {
		panic("")
	}

	if !EndWith(s, sub1, true) {
		panic("")
	}
	if EndWith(s, sub2, true) {
		panic("")
	}
	if !EndWith(s, sub3, true) {
		panic("")
	}
	if !EndWith(s, sub4, true) {
		panic("")
	}
}

func TestElemInArr(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			t.Error()
		}
	}()

	arr := []int{1, 2, 3}
	if !ElemIn(1, arr) {
		panic("")
	}
	if !ElemIn(2, arr) {
		panic("")
	}
	if !ElemIn(3, arr) {
		panic("")
	}
}

func TestIntersectionStr(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			t.Error()
		}
	}()

	arr1 := []string{"1", "2", "3"}
	arr2 := []string{"2", "3", "5"}
	arr3 := []string{"4", "5", "6"}
	if len(Intersection(arr1, arr2)) != 2 {
		panic("")
	}
	if len(Intersection(arr1, arr3)) != 0 {
		panic("")
	}
}

func TestUnionStr(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			t.Error()
		}
	}()

	arr1 := []string{"1"}
	arr2 := []string{"2"}
	arr3 := Union(arr1, arr2)
	if len(arr3) != 2 {
		panic("")
	}
	if !ElemIn("1", arr3) || !ElemIn("2", arr3) {
		panic("")
	}
}

func TestDifferenceStr(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			t.Error()
		}
	}()

	arr1 := []string{"1", "2", "3"}
	arr2 := []string{"2", "3", "5"}
	if Difference(arr1, arr2)[0] != "1" {
		panic("")
	}
}

func TestJoinIntArr(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			t.Error()
		}
	}()

	arr := []int{1, 2, 3}
	s := JoinArr(arr, ",")
	if s != "1,2,3" {
		panic("")
	}
	fmt.Println(s)
}

func TestJoinStrArr(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			t.Error()
		}
	}()

	arr := []string{"1", "2", "3"}
	s := JoinArr(arr, ",")
	if s != "1,2,3" {
		panic("")
	}
	fmt.Println(s)
}

func TestGetMapValue(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			t.Error()
		}
	}()

	m := map[string]interface{}{
		"test": "123",
	}

	if GetMapValue(m, "test", "") != "123" {
		panic("")
	}
}

func TestParamsSorted(t *testing.T) {
	data := map[string]string{
		"1": "1",
		"2": "2",
		"3": "3",
	}
	s := ParamsSorted(data)
	fmt.Println(s)
}

func TestStrArrSplit(t *testing.T) {
	arr := []string{"1", "2", "3"}
	fmt.Println(ArrSplit(arr, 1))
}

func TestRandomInt64(t *testing.T) {
	r := RandomInt64(50, 100)

	fmt.Println(r)
}

func TestHumpFormatToUnderLine(t *testing.T) {
	fmt.Println(HumpFormatToUnderLine("QweHabSEd1"))
}

func TestMapHumpFormatToUnderLine(t *testing.T) {
	m := map[string]interface{}{
		"tpOrderId": "123",
		"money":     1.01,
	}
	fmt.Println(MapHumpFormatToUnderLine(m))
}

func TestMergeMapArr(t *testing.T) {
	m := map[string]interface{}{
		"test": 1,
	}

	m = MergeMap[string, interface{}](nil, m, nil, map[string]interface{}{
		"test2": 2,
	})

	fmt.Println(m)
}

func TestStr2Split(t *testing.T) {
	s := "sn in [1,2,3] && ver >= 1.2.3.4 || ver == 1.0.0.0 || channel == 123 && sn == 1 && !(sn in [1])"
	seq := []string{"in", "&&", "||", "nin", "==", ">=", "<=", "!=", ">", "<", "!"}

	arr, seqArr := StrSplit(s, seq, true)

	fmt.Println(arr, seqArr)
}

func TestGo(t *testing.T) {
	SafeGo(func() {
		panic(123)
	})

	fmt.Println("123")
	time.Sleep(time.Second * 2)
	fmt.Println(111)
}

func TestMerge(t *testing.T) {
	s := "//q3123123123/12312/312/321/3/21/3/21/3/123/12/3///32/13/123123"
	s2 := string(CharMerge(s, byte('/')))
	fmt.Println(s2)
}
