package collections

import (
	"reflect"
	"sort"
)

type Sorter struct {
	descending bool
	keyFetcher KeyFetcher
	comparator Comparator
}

func (s *Sorter) Ascend() {
	s.descending = false
}

func (s *Sorter) Descend() {
	s.descending = true
}

func (s *Sorter) By(kf KeyFetcher) {
	s.keyFetcher = kf
}

func (s *Sorter) ByValue() {
	s.keyFetcher = func(elem interface{}) interface{} {
		return elem
	}
}

func (s *Sorter) ByField(fieldName string) {
	s.keyFetcher = func(elem interface{}) interface{} {
		objValue := reflect.ValueOf(elem)
		if objValue.Kind() == reflect.Ptr {
			objValue = reflect.Indirect(objValue)
		}
		return objValue.FieldByName(fieldName).Interface()
	}
}

func (s *Sorter) ByMapKey(key interface{}) {
	s.keyFetcher = func(elem interface{}) interface{} {
		keyNameValue := reflect.ValueOf(key)
		mapValue := reflect.ValueOf(elem)
		return mapValue.MapIndex(keyNameValue).Interface()
	}
}

func (s *Sorter) BySliceIndex(index int) {
	s.keyFetcher = func(elem interface{}) interface{} {
		listValue := reflect.ValueOf(elem)
		return listValue.Index(index).Interface()
	}
}

func (s *Sorter) Comparator(kw Comparator) {
	s.keyWrapper = kw
}

func (s *Sorter) Natural() {
	// TODO
}

func (s *Sorter) Numerical() {
	// TODO
}

func (s *Sorter) String() {
	// TODO
}

func (s *Sorter) Sort(slice interface{}) {
	if s.keyFetcher == nil {
		// TODO default key fetcher
	}
	if s.comparator == nil {
		// TODO default key wrapper
	}
	sorter := sorter{
		Sorter: s,
		slice:  reflect.ValueOf(slice),
		keys:   nil,
	}
	sort.Sort(sorter)
}

type Comparator interface {
	Less(a, b interface{}) bool
	PrepareKeys(elements []interface{}) []interface{}
}

type KeyFetcher func(interface{}) interface{}

type sorter struct {
	*Sorter
	slice reflect.Value
	keys  []interface{}
}

func (s sorter) Len() int {
	return len(s.keys)
}

func (s sorter) Swap(a, b int) {
	// Swap in keys
	s.keys[a], s.keys[b] = s.keys[b], s.keys[a]

	// Swap in slice
	aValue := s.slice.Index(a)
	bValue := s.slice.Index(b)
	tmpValue := reflect.Indirect(reflect.New(aValue.Type()))
	tmpValue.Set(aValue)
	aValue.Set(bValue)
	bValue.Set(tmpValue)
}

func (s sorter) Less(a, b int) bool {
	rv := s.keys[a].Less(s.keys[b])
	if s.descending {
		rv = !rv
	}
	return rv
}

// Natural comparator (AKA human / version ordering)
type NaturalComparator struct{}

func (NaturalComparator) Less(a, b interface{}) bool {
	return false
}

func (NaturalComparator) PrepareKeys(keys []interface{}) []interface{} {
	return keys
}

// NumberComparator attempts to convert
type NumberComparator struct{}

func (NumberComparator) Less(a, b interface{}) bool {
	return false
}

func (NumberComparator) PrepareKeys(keys []interface{}) []interface{} {
	return keys
}

type StringComparator struct{}

func (StringComparator) Less(a, b interface{}) bool {
	return false
}

func (StringComparator) PrepareKeys(keys []interface{}) []interface{} {
	return keys
}
