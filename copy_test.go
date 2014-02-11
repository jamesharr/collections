package collections_test

import (
	. "github.com/jamesharr/collections"
	"reflect"
	"strings"
	"testing"
)

func assertDeepEq(t *testing.T, got, exp interface{}, failMsg string) {
	if !reflect.DeepEqual(got, exp) {
		t.Error(failMsg)
		t.Errorf(" Expected %v", exp)
		t.Errorf(" Got %v", got)
	}
}

func expectPanic(t *testing.T, mustSubStr string) {
	r := recover()
	s, ok := r.(string)
	if r == nil {
		t.Error(" Expected: some panic")
		t.Error(" Got: none")
	} else if !ok {
		t.Error(" Expected: string panic")
		t.Errorf(" Got: %T %v", r, r)
	} else if !strings.Contains(s, mustSubStr) {
		t.Errorf(" Expected substring %q", mustSubStr)
		t.Errorf(" Got: %v", s)
	}
}

func TestCopy_slice(t *testing.T) {
	exp := []int{1, 2, 3}
	got := Copy(exp).([]int)
	assertDeepEq(t, got, exp, "Slice copy failed")
}

func TestCopy_map(t *testing.T) {
	exp := make(map[int]bool)
	got := Copy(exp).(map[int]bool)
	assertDeepEq(t, got, exp, "Map copy failed")
}

func TestCopy_badType1(t *testing.T) {
	defer expectPanic(t, "src must be an array, slice, or map.")
	Copy(5)
}

func TestCopy_badType2(t *testing.T) {
	defer expectPanic(t, "src must be an array, slice, or map.")
	exp := []int{1, 2, 3}
	Copy(&exp)
}

func TestCopyInto_slice1(t *testing.T) {
	src := []int{1, 2, 3}
	exp := []float32{1.0, 2.0, 3.0}
	var got []float32
	CopyInto(src, &got)
	assertDeepEq(t, got, exp, "Slice CopyInto failed")
}

func TestCopyInto_slice2(t *testing.T) {
	type T int
	src := []int{1, 2, 3}
	exp := []T{T(1), T(2), T(3)}
	var got []T
	CopyInto(src, &got)
	assertDeepEq(t, got, exp, "Slice CopyInto failed")
}

// Copy into same type
func TestCopyInto_map_sameType(t *testing.T) {
	src := make(map[string]int)
	src["one"] = 1
	src["two"] = 2
	src["three"] = 3
	exp := src
	var got map[string]int
	CopyInto(src, &got)
	assertDeepEq(t, got, exp, "Map CopyInto failed")
}

// Convert into generic type
func TestCopyInto_map_convert1(t *testing.T) {
	src := make(map[string]int)
	src["one"] = 1
	src["two"] = 2
	src["three"] = 3
	exp := make(map[interface{}]interface{})
	exp["one"] = 1
	exp["two"] = 2
	exp["three"] = 3
	var got map[interface{}]interface{}
	CopyInto(src, &got)
	assertDeepEq(t, got, exp, "Map CopyInto failed")
}

// Convert into specific type
func TestCopyInto_map_convert2(t *testing.T) {
	src := make(map[interface{}]interface{})
	src["one"] = int(1)
	src["two"] = int(2)
	src["three"] = int(3)
	exp := make(map[string]int)
	exp["one"] = 1
	exp["two"] = 2
	exp["three"] = 3
	var got map[string]int
	CopyInto(src, &got)
	assertDeepEq(t, got, exp, "Map CopyInto failed")
}
