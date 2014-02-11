package collections_test

import (
	"testing"
	. "github.com/jamesharr/collections"
	"reflect"
	"strings"
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

func TestClone_slice(t *testing.T) {
	exp := []int{1,2,3}
	got := Clone(exp).([]int)
	assertDeepEq(t, got, exp, "Slice clone failed")
}

func TestClone_map(t *testing.T) {
	exp := make(map[int]bool)
	got := Clone(exp).(map[int]bool)
	assertDeepEq(t, got, exp, "Map clone failed")
}

func TestClone_badType1(t *testing.T) {
	defer expectPanic(t, "First parameter to Clone() must be")
	Clone(5)
}

func TestClone_badType2(t *testing.T) {
	defer expectPanic(t, "First parameter to Clone() must be")
	exp := []int{1,2,3}
	Clone(&exp)
}
