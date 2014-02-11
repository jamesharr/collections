package collections

import . "reflect"

// Clone() makes a shallow copy of an array, slice
func Clone(src interface{}) interface{} {
	srcValue := ValueOf(src)
	srcType := srcValue.Type()

	// Check type
	kind := srcValue.Kind()
	if kind != Array && kind != Slice && kind != Map {
		panic("First parameter to Clone() must be an array, slice, or map.")
	}

	// Make new container
	dstValue := Indirect(New(srcType))

	// Copy items
	if kind == Map {
		dstValue.Set(MakeMap(srcType))
		for _, key := range srcValue.MapKeys() {
			value := srcValue.MapIndex(key)
			dstValue.SetMapIndex(key, value)
		}
	} else {
		srcLen := srcValue.Len()
		dstValue.Set(MakeSlice(srcType, srcLen, srcLen))
		for i := 0; i < srcLen; i++ {
			dstValue.Index(i).Set(srcValue.Index(i))
		}
	}

	// Return dest
	return dstValue.Interface()
}
