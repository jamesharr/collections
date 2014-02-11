package collections

import (
	"reflect"
	"sort"
	"fmt"
	"regexp"
	"strconv"
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

func (s *Sorter) Comparator(c Comparator) {
	s.comparator = c
}

func (s *Sorter) Natural() {
	s.comparator = NaturalComparator{}
}

func (s *Sorter) Version() {
	s.comparator = VersionComparator{}
}

func (s *Sorter) String() {
	s.comparator = StringComparator{}
}

func (s *Sorter) Sort(slice interface{}) {
	sorter := sorter{
		Sorter: s,
		slice:  reflect.ValueOf(slice),
		keys:   nil,
	}
	if s.keyFetcher == nil {
		s.ByValue()
	}
	if s.comparator == nil {
		s.comparator = StringComparator{}
	}

	// Fetch Keys
	sliceLen := sorter.slice.Len()
	sorter.keys = make([]interface{}, sliceLen)
	for i := 0; i < sliceLen; i++ {
		v := sorter.slice.Index(i).Interface()
		sorter.keys[i] = s.keyFetcher(v)
	}

	// Prepare Keys
	s.comparator.PrepareKeys(sorter.keys)

	// Sort
	sort.Sort(sorter)
}

type Comparator interface {
	Less(a, b interface{}) bool
	PrepareKeys(elements []interface{})
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
	rv := s.comparator.Less(s.keys[a], s.keys[b])
	if s.descending {
		rv = !rv
	}
	return rv
}

/* Natural sort (aka human sort)

 2.10 < 2.3 < 10.1
 */
type NaturalComparator struct{}

func (NaturalComparator) Less(a, b interface{}) bool {
	aList := a.([]interface{})
	bList := b.([]interface{})

	aLen := len(aList)
	bLen := len(bList)
	minLen := aLen
	if minLen > bLen {
		minLen = bLen
	}

	// Compare each item in list
	for i := 0; i < minLen; i++ {
		cmp := genericCompare(aList[i], bList[i])
		if cmp != equal {
			return cmp.Less()
		}
	}

	// Fall back to item length
	if aLen < bLen {
		return true
	} else {
		return false
	}
}

func (NaturalComparator) PrepareKeys(keys []interface{}) {
	for i, k := range keys {
		keys[i] = humanTokenizer(k.(string))
	}
}

/* Version sorter

 Similar to natural comparator, except there's no floating point tokens, just integers.

 2.3 < 2.10 < 10.1
 */
type VersionComparator struct{}

func (VersionComparator) Less(a, b interface{}) bool {
	// TODO
	return false
}

func (VersionComparator) PrepareKeys(keys []interface{}) {
	// TODO
}

type StringComparator struct{}

func (StringComparator) Less(a, b interface{}) bool {
	aStr := a.(string)
	bStr := b.(string)
	return aStr < bStr
}

func (StringComparator) PrepareKeys(keys []interface{}) {
	for i, k := range keys {
		if k, ok := k.(string); ok {
			// done
		} else {
			k := fmt.Sprintf("%s", k)
			keys[i] = k
		}
	}
}

type equality uint8
func (e equality) Less() bool {
	return e == less
}
const (
	less    = equality(iota)
	equal
	greater
)

func stringCompare(a, b string) equality {
	if a < b {
		return less
	} else if a > b {
		return greater
	} else {
		return equal
	}
}
func intCompare(a, b int64) equality {
	if a < b {
		return less
	} else if a > b {
		return greater
	} else {
		return equal
	}
}
func floatCompare(a, b float64) equality {
	if a < b {
		return less
	} else if a > b {
		return greater
	} else {
		return equal
	}
}

func genericCompare(a, b interface{}) equality {
	if aStr, ok := a.(string); ok {
		if bStr, ok := b.(string); ok {
			return stringCompare(aStr,bStr)
		} else if _, ok := b.(int64); ok {
			// Numbers are always less than strings
			return less
		} else if _, ok := b.(float64); ok {
			// Numbers are always less than strings
			return less
		}
	} else if aInt, ok := a.(int64); ok {
		if bInt, ok := b.(int64); ok {
			return intCompare(aInt,bInt)
		} else if bFloat, ok := b.(float64); ok {
			return floatCompare(float64(aInt), bFloat)
		} else if _, ok := b.(string); ok {
			// Numbers are always less than strings
			return greater
		}
	} else if aFloat, ok := a.(float64); ok {
		if bInt, ok := b.(int64); ok {
			return floatCompare(aFloat,float64(bInt))
		} else if bFloat, ok := b.(float64); ok {
			return floatCompare(aFloat, bFloat)
		} else if _, ok := b.(string); ok {
			// Numbers are always less than strings
			return greater
		}
	}

	aType := reflect.TypeOf(a)
	bType := reflect.TypeOf(b)
	panic(fmt.Sprintf("Don't know how to compare %v and %v", aType, bType))
	return equal
}

var humanTokenizerRegexp = regexp.MustCompile(`(\d+\.\d+|\d+|\D+)`)
func humanTokenizer(s string) []interface{} {
	tokens := humanTokenizerRegexp.FindAllString(s,-1)
	rv := make([]interface{}, 0, len(tokens))
	for _, t := range tokens {
		rv = append(rv, autoStrConv(t))
	}
	return rv
}

var versionTokenizerRegexp = regexp.MustCompile(`(\d+|\D+)`)
func versionTokenizer(s string) []interface{} {
	tokens := versionTokenizerRegexp.FindAllString(s,-1)
	rv := make([]interface{}, 0, len(tokens))
	for _, t := range tokens {
		rv = append(rv, autoStrConv(t))
	}
	return rv
}

func autoStrConv(s string) interface{} {
	if intVal, err := strconv.ParseInt(s, 10, 64); err == nil {
		return intVal
	} else if floatVal, err := strconv.ParseFloat(s, 64); err == nil {
		return floatVal
	} else {
		return s
	}
}
