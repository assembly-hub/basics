// Package set
package set

import (
	"fmt"
	"testing"
)

func TestSet_Add(t *testing.T) {
	s1 := Set[string]{}
	s2 := Set[string]{}

	s1.Add("1", "2", "3")
	s2.Add("3", "5", "6")

	fmt.Println("Union: ", s1.Union(s2).ToList())
	fmt.Println("Intersection: ", s1.Intersection(s2).ToList())
	fmt.Println("Difference: ", s1.Difference(s2).ToList())
}
