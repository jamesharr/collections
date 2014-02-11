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
	s.ByField("Name")
	s.Sort(slc)
	t.Log(slc)

	t.Error("I'm pretty sure this isn't complete")
}
