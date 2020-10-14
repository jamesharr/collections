# PROJECT STATUS

This project should not be used in production. Feel free to examine this code for reference and/or use it if you dare. It's probably not a good idea though because it's not very Go-like.

# Collection routines for Go

[![build status](https://secure.travis-ci.org/jamesharr/collections.png)](http://travis-ci.org/jamesharr/collections)

A small library that does nasty things with reflection as an API experiment.

## Highlights

* `Copy(src) -> copy)` that can do allocation for you.
* `CopyInto(src, &dstPtr)` that can perform type conversion on the fly.
* Flexible sort API. (up and coming)

## Examples

```go

// Copy a slice
src := []int{1, 2, 3}
dst := Copy(src).([]int)

// Copy a slice, convert the type
src := []int{1, 2, 3}
var dst []float32
CopyInto(src, &dst)

// or not
src := []int{1, 2, 3}
var dst []int
CopyInto(src, &dst)

// Copy a slice, convert the type
type Foo int
func (f Foo) SayHi() {
	fmt.Println("Hello ", f)
}
src := []int{1, 2, 3}
var dst []Foo
CopyInto(src, &dst)
for _, v := range dst {
	v.SayHi()
}

// Sort a slice by a key
// Sort a list by human ordering

```

## To Do List

* Sorting
* Benchmarks

## Commentary

If performance or type saftey are your goals, you may want to look somewhere else.

Yes, This library does disgusting things with reflection. It's going to be slower than
writing things by hand and/or using type assertions. It's not my intention to build a
dynamic language with Go. My intention is to reduce the number of alloc+copy+cast
loops I need to write while throwing together small demo apps.

I would much rather do all this with type safety, but in its current revision, Go does
not facilitate this, so I must resort to reflection.
