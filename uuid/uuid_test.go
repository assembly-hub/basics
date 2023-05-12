package uuid

import (
	"fmt"
	"testing"
)

func TestNewV4(t *testing.T) {
	uu, err := NewV4()
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(uu.String())
}
