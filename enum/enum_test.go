package enum

import (
	"fmt"
	"testing"
)

func TestEnumType_RawVal(t *testing.T) {
	e := New[struct {
		First  Elem[int] `code:"1" text:"123"`
		Second Elem[int] `code:"2" text:"123"`
	}]()

	fmt.Println(Code2Text[int](e))
	fmt.Println(e.First.Code, e.First.Text)
}

func TestEnumType_RawVal2(t *testing.T) {
	e := New[struct {
		First  Elem[string] `code:"1" text:"123"`
		Second Elem[string] `code:"2" text:"123"`
	}]()

	fmt.Println(Code2Text[string](e))
	fmt.Println(e.First.Code, e.First.Text)
}
