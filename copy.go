package collections

import (
	r "reflect"
)

// Copy() makes a shallow copy of an array, slice, or map
func Copy(src interface{}) interface{} {
	// Make new container
	dstPtrValue := r.New(r.TypeOf(src))

	// Copy
	CopyInto(src, dstPtrValue.Interface())

	// Return dest
	return r.Indirect(dstPtrValue).Interface()
}

// CopyInto() copies+converts a src array, slice, or map into a destination.
//
// The element type of src must be convertible to the element type of dst, but don't have
// to be the same type.
func CopyInto(src interface{}, dstPtr interface{}) {
	srcValue := r.ValueOf(src)

	// Check kind of src
	kind := srcValue.Kind()
	if kind != r.Array && kind != r.Slice && kind != r.Map {
		panic("CopyInto(src, dstPtr): src must be an array, slice, or map")
	}

	// Check kind of dst
	dstPtrValue := r.ValueOf(dstPtr)
	if dstKind := dstPtrValue.Kind(); dstKind != r.Ptr {
		panic("CopyInto(src, dstPtr): dstPtr must be a pointer to a slice or map")
	}
	dstValue := r.Indirect(dstPtrValue)
	if k := dstValue.Kind(); k != r.Slice && k != r.Map {
		panic("CopyInto(src, dstPtr): dstPtr must be a pointer to a slice or map")
	}

	// TODO: Check for element / key+value type convertibility.
	// Isn't a high priority. The result either way is a panic. By us handling
	// the type checking, we'll just end up with better messages.

	if kind == r.Slice || kind == r.Array {
		dstElemType := dstValue.Type().Elem()

		// Prepare dst
		srcLen := srcValue.Len()
		if srcLen < dstValue.Cap() {
			dstValue.SetLen(srcLen)
		} else {
			dstValue.Set(r.MakeSlice(dstValue.Type(), srcLen, srcLen))
		}

		// Copy
		for i := 0; i < srcLen; i++ {
			e := srcValue.Index(i)
			e = e.Convert(dstElemType)
			dstValue.Index(i).Set(e)
		}
	} else if kind == r.Map {
		dstType := dstValue.Type()
		dstKeyType := dstType.Key()
		dstElemType := dstType.Elem()

		// Clear destination
		if dstValue.Len() != 0 || dstValue.IsNil() {
			dstValue.Set(r.MakeMap(dstType))
		}

		// Copy
		for _, k := range srcValue.MapKeys() {
			v := srcValue.MapIndex(k)

			// Copy value and convert
			k = r.ValueOf(k.Interface())
			k = k.Convert(dstKeyType)

			// Copy value and convert
			v = r.ValueOf(v.Interface())
			v = v.Convert(dstElemType)

			// Set in destination map
			dstValue.SetMapIndex(k, v)
		}
	}
}
