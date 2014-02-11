package collections_test

import (
	"github.com/jamesharr/collections"
	"testing"
)

type Person struct {
	Name   string
	Room   string
	Height int
}

var samples = []Person{
	Person{"James", "EAB 30", 74},
	Person{"Dave", "EAB 9", 76},
	Person{"Jaime", "EAB 10", 78},
}

func TestSort_fields(t *testing.T) {
	slc := collections.Copy(samples).([]Person)

	s := collections.Sorter{}

	// s.By(func(a interface{})interface{}{
	//    return a.(Person).Room
	// })
	s.ByField("Room")

	// Value sort (default): compare normal types
	// Natural sort: "foo-2.30" < "foo-2.1"  < "foo-010.3"
	// Version sort: "foo-2.1"  < "foo-2.30" < "foo-010.3"
	s.Natural()

	// Go sort this slice
	s.Sort(slc)

	// Debug log since there's really no unit tests
	t.Log(slc)
}
