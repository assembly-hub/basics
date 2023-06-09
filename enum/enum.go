package enum

import (
	"reflect"

	"github.com/assembly-hub/basics/set"
	"github.com/assembly-hub/basics/util"
)

type elemType interface {
	int | string
}

type Elem[T elemType] struct {
	Code T
	Text string
}

func New[T any]() *T {
	var e T
	loadEnumData(&e)
	return &e
}

func loadEnumData(enumPtr interface{}) {
	dataVal := reflect.ValueOf(enumPtr)
	if dataVal.Kind() != reflect.Ptr {
		panic("Enum's type need struct ptr")
	}
	dataVal = dataVal.Elem()
	if dataVal.Kind() != reflect.Struct {
		panic("Enum's type need struct ptr")
	}

	for i := 0; i < dataVal.NumField(); i++ {
		if !dataVal.Field(i).CanSet() {
			continue
		}
		code := dataVal.Type().Field(i).Tag.Get("code")
		text := dataVal.Type().Field(i).Tag.Get("text")

		if dataVal.Type().Field(i).Type.Kind() != reflect.Struct {
			panic("Enum Field type need Elem[int | string]")
		}

		eType, ok := dataVal.Type().Field(i).Type.FieldByName("Code")
		if !ok {
			panic("Enum Field type need Elem[int | string]")
		}

		var v interface{}
		switch eType.Type.Kind() {
		case reflect.Int:
			i, err := util.Str2Int[int](code)
			if err != nil {
				panic(err)
			}

			v = Elem[int]{
				Code: i,
				Text: text,
			}
		case reflect.String:
			v = Elem[string]{
				Code: code,
				Text: text,
			}
		default:
			panic("EnumData Field type need Elem[int | string]")
		}

		dataVal.Field(i).Set(reflect.ValueOf(v))
	}
}

func Code2Text[T elemType](enumPtr interface{}) map[T]string {
	dataVal := reflect.ValueOf(enumPtr)
	if dataVal.Kind() != reflect.Ptr {
		panic("Code2TextForInt param type need struct ptr")
	}
	dataVal = dataVal.Elem()
	if dataVal.Kind() != reflect.Struct {
		panic("Code2TextForInt param type need struct ptr")
	}

	m := map[T]string{}
	for i := 0; i < dataVal.NumField(); i++ {
		if dataVal.Type().Field(i).Type.Kind() != reflect.Struct {
			panic("Code2Text Field type need Elem[int | string]")
		}

		enumType, ok := dataVal.Type().Field(i).Type.FieldByName("Code")
		if !ok {
			panic("Code2Text Field type need Elem[int | string]")
		}

		switch enumType.Type.Kind() {
		case reflect.Int, reflect.String:
			v := dataVal.Field(i).Interface().(Elem[T])
			m[v.Code] = v.Text
		default:
			panic("Code2Text Field type need Elem[int | string]")
		}
	}
	return m
}

// ID2NameList 返回数据结构为：
//
//	[{
//	     "id": enum code,
//	     "name": enum text
//	 },{
//	     ...
//	 }]
func ID2NameList(enumPtr interface{}) []map[string]interface{} {
	dataVal := reflect.ValueOf(enumPtr)
	if dataVal.Kind() != reflect.Ptr {
		panic("Id2NameList param type need struct ptr")
	}
	dataVal = dataVal.Elem()
	if dataVal.Kind() != reflect.Struct {
		panic("Id2NameList param type need struct ptr")
	}

	var arr []map[string]interface{}
	for i := 0; i < dataVal.NumField(); i++ {
		if dataVal.Type().Field(i).Type.Kind() != reflect.Struct {
			panic("ID2NameList Field type need Elem[int | string]")
		}

		enumType, ok := dataVal.Type().Field(i).Type.FieldByName("Code")
		if !ok {
			panic("ID2NameList Field type need Elem[int | string]")
		}

		switch enumType.Type.Kind() {
		case reflect.Int:
			v := dataVal.Field(i).Interface().(Elem[int])
			arr = append(arr, map[string]interface{}{
				"id":   v.Code,
				"name": v.Text,
			})
		case reflect.String:
			v := dataVal.Field(i).Interface().(Elem[string])
			arr = append(arr, map[string]interface{}{
				"id":   v.Code,
				"name": v.Text,
			})
		default:
			panic("ID2NameList Field type need Elem[int | string]")
		}
	}
	return arr
}

func CodeSet[T elemType](enumPtr interface{}) set.Set[T] {
	dataVal := reflect.ValueOf(enumPtr)
	if dataVal.Kind() != reflect.Ptr {
		panic("Code2TextForInt param type need struct ptr")
	}
	dataVal = dataVal.Elem()
	if dataVal.Kind() != reflect.Struct {
		panic("Code2TextForInt param type need struct ptr")
	}

	s := set.New[T]()
	for i := 0; i < dataVal.NumField(); i++ {
		if dataVal.Type().Field(i).Type.Kind() != reflect.Struct {
			panic("Code2TextForInt Field type need Elem[int | string]")
		}

		enumType, ok := dataVal.Type().Field(i).Type.FieldByName("Code")
		if !ok {
			panic("Code2TextForInt Field type need Elem[int | string]")
		}

		switch enumType.Type.Kind() {
		case reflect.Int, reflect.String:
			v := dataVal.Field(i).Interface().(Elem[T])
			s.Add(v.Code)
		default:
			panic("Code2TextForInt Field type need Elem[int | string]")
		}
	}
	return s
}
