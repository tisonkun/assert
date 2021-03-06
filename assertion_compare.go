// Copyright 2022 tison <wander4096@gmail.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package assert

import (
	"bytes"
	"fmt"
	"reflect"
	"time"
)

type CompareType int

const (
	compareLess CompareType = iota - 1
	compareEqual
	compareGreater
)

var (
	intType   = reflect.TypeOf(1)
	int8Type  = reflect.TypeOf(int8(1))
	int16Type = reflect.TypeOf(int16(1))
	int32Type = reflect.TypeOf(int32(1))
	int64Type = reflect.TypeOf(int64(1))

	uintType   = reflect.TypeOf(uint(1))
	uint8Type  = reflect.TypeOf(uint8(1))
	uint16Type = reflect.TypeOf(uint16(1))
	uint32Type = reflect.TypeOf(uint32(1))
	uint64Type = reflect.TypeOf(uint64(1))

	float32Type = reflect.TypeOf(float32(1))
	float64Type = reflect.TypeOf(float64(1))

	stringType = reflect.TypeOf("")
	timeType   = reflect.TypeOf(time.Time{})
	bytesType  = reflect.TypeOf([]byte{})
)

func compare(obj1, obj2 any, kind reflect.Kind) (CompareType, bool) {
	obj1Value := reflect.ValueOf(obj1)
	obj2Value := reflect.ValueOf(obj2)

	// throughout this switch we try and avoid calling .Convert() if possible,
	// as this has a pretty big performance impact
	switch kind {
	case reflect.Int:
		{
			intobj1, ok := obj1.(int)
			if !ok {
				intobj1 = obj1Value.Convert(intType).Interface().(int)
			}
			intobj2, ok := obj2.(int)
			if !ok {
				intobj2 = obj2Value.Convert(intType).Interface().(int)
			}
			if intobj1 > intobj2 {
				return compareGreater, true
			}
			if intobj1 == intobj2 {
				return compareEqual, true
			}
			if intobj1 < intobj2 {
				return compareLess, true
			}
		}
	case reflect.Int8:
		{
			int8obj1, ok := obj1.(int8)
			if !ok {
				int8obj1 = obj1Value.Convert(int8Type).Interface().(int8)
			}
			int8obj2, ok := obj2.(int8)
			if !ok {
				int8obj2 = obj2Value.Convert(int8Type).Interface().(int8)
			}
			if int8obj1 > int8obj2 {
				return compareGreater, true
			}
			if int8obj1 == int8obj2 {
				return compareEqual, true
			}
			if int8obj1 < int8obj2 {
				return compareLess, true
			}
		}
	case reflect.Int16:
		{
			int16obj1, ok := obj1.(int16)
			if !ok {
				int16obj1 = obj1Value.Convert(int16Type).Interface().(int16)
			}
			int16obj2, ok := obj2.(int16)
			if !ok {
				int16obj2 = obj2Value.Convert(int16Type).Interface().(int16)
			}
			if int16obj1 > int16obj2 {
				return compareGreater, true
			}
			if int16obj1 == int16obj2 {
				return compareEqual, true
			}
			if int16obj1 < int16obj2 {
				return compareLess, true
			}
		}
	case reflect.Int32:
		{
			int32obj1, ok := obj1.(int32)
			if !ok {
				int32obj1 = obj1Value.Convert(int32Type).Interface().(int32)
			}
			int32obj2, ok := obj2.(int32)
			if !ok {
				int32obj2 = obj2Value.Convert(int32Type).Interface().(int32)
			}
			if int32obj1 > int32obj2 {
				return compareGreater, true
			}
			if int32obj1 == int32obj2 {
				return compareEqual, true
			}
			if int32obj1 < int32obj2 {
				return compareLess, true
			}
		}
	case reflect.Int64:
		{
			int64obj1, ok := obj1.(int64)
			if !ok {
				int64obj1 = obj1Value.Convert(int64Type).Interface().(int64)
			}
			int64obj2, ok := obj2.(int64)
			if !ok {
				int64obj2 = obj2Value.Convert(int64Type).Interface().(int64)
			}
			if int64obj1 > int64obj2 {
				return compareGreater, true
			}
			if int64obj1 == int64obj2 {
				return compareEqual, true
			}
			if int64obj1 < int64obj2 {
				return compareLess, true
			}
		}
	case reflect.Uint:
		{
			uintobj1, ok := obj1.(uint)
			if !ok {
				uintobj1 = obj1Value.Convert(uintType).Interface().(uint)
			}
			uintobj2, ok := obj2.(uint)
			if !ok {
				uintobj2 = obj2Value.Convert(uintType).Interface().(uint)
			}
			if uintobj1 > uintobj2 {
				return compareGreater, true
			}
			if uintobj1 == uintobj2 {
				return compareEqual, true
			}
			if uintobj1 < uintobj2 {
				return compareLess, true
			}
		}
	case reflect.Uint8:
		{
			uint8obj1, ok := obj1.(uint8)
			if !ok {
				uint8obj1 = obj1Value.Convert(uint8Type).Interface().(uint8)
			}
			uint8obj2, ok := obj2.(uint8)
			if !ok {
				uint8obj2 = obj2Value.Convert(uint8Type).Interface().(uint8)
			}
			if uint8obj1 > uint8obj2 {
				return compareGreater, true
			}
			if uint8obj1 == uint8obj2 {
				return compareEqual, true
			}
			if uint8obj1 < uint8obj2 {
				return compareLess, true
			}
		}
	case reflect.Uint16:
		{
			uint16obj1, ok := obj1.(uint16)
			if !ok {
				uint16obj1 = obj1Value.Convert(uint16Type).Interface().(uint16)
			}
			uint16obj2, ok := obj2.(uint16)
			if !ok {
				uint16obj2 = obj2Value.Convert(uint16Type).Interface().(uint16)
			}
			if uint16obj1 > uint16obj2 {
				return compareGreater, true
			}
			if uint16obj1 == uint16obj2 {
				return compareEqual, true
			}
			if uint16obj1 < uint16obj2 {
				return compareLess, true
			}
		}
	case reflect.Uint32:
		{
			uint32obj1, ok := obj1.(uint32)
			if !ok {
				uint32obj1 = obj1Value.Convert(uint32Type).Interface().(uint32)
			}
			uint32obj2, ok := obj2.(uint32)
			if !ok {
				uint32obj2 = obj2Value.Convert(uint32Type).Interface().(uint32)
			}
			if uint32obj1 > uint32obj2 {
				return compareGreater, true
			}
			if uint32obj1 == uint32obj2 {
				return compareEqual, true
			}
			if uint32obj1 < uint32obj2 {
				return compareLess, true
			}
		}
	case reflect.Uint64:
		{
			uint64obj1, ok := obj1.(uint64)
			if !ok {
				uint64obj1 = obj1Value.Convert(uint64Type).Interface().(uint64)
			}
			uint64obj2, ok := obj2.(uint64)
			if !ok {
				uint64obj2 = obj2Value.Convert(uint64Type).Interface().(uint64)
			}
			if uint64obj1 > uint64obj2 {
				return compareGreater, true
			}
			if uint64obj1 == uint64obj2 {
				return compareEqual, true
			}
			if uint64obj1 < uint64obj2 {
				return compareLess, true
			}
		}
	case reflect.Float32:
		{
			float32obj1, ok := obj1.(float32)
			if !ok {
				float32obj1 = obj1Value.Convert(float32Type).Interface().(float32)
			}
			float32obj2, ok := obj2.(float32)
			if !ok {
				float32obj2 = obj2Value.Convert(float32Type).Interface().(float32)
			}
			if float32obj1 > float32obj2 {
				return compareGreater, true
			}
			if float32obj1 == float32obj2 {
				return compareEqual, true
			}
			if float32obj1 < float32obj2 {
				return compareLess, true
			}
		}
	case reflect.Float64:
		{
			float64obj1, ok := obj1.(float64)
			if !ok {
				float64obj1 = obj1Value.Convert(float64Type).Interface().(float64)
			}
			float64obj2, ok := obj2.(float64)
			if !ok {
				float64obj2 = obj2Value.Convert(float64Type).Interface().(float64)
			}
			if float64obj1 > float64obj2 {
				return compareGreater, true
			}
			if float64obj1 == float64obj2 {
				return compareEqual, true
			}
			if float64obj1 < float64obj2 {
				return compareLess, true
			}
		}
	case reflect.String:
		{
			stringobj1, ok := obj1.(string)
			if !ok {
				stringobj1 = obj1Value.Convert(stringType).Interface().(string)
			}
			stringobj2, ok := obj2.(string)
			if !ok {
				stringobj2 = obj2Value.Convert(stringType).Interface().(string)
			}
			if stringobj1 > stringobj2 {
				return compareGreater, true
			}
			if stringobj1 == stringobj2 {
				return compareEqual, true
			}
			if stringobj1 < stringobj2 {
				return compareLess, true
			}
		}
	// Check for known struct types we can check for compare results.
	case reflect.Struct:
		{
			// All structs enter here. We're not interested in most types.
			if !canConvert(obj1Value, timeType) {
				break
			}

			// time.Time can compared!
			timeObj1, ok := obj1.(time.Time)
			if !ok {
				timeObj1 = obj1Value.Convert(timeType).Interface().(time.Time)
			}

			timeObj2, ok := obj2.(time.Time)
			if !ok {
				timeObj2 = obj2Value.Convert(timeType).Interface().(time.Time)
			}

			return compare(timeObj1.UnixNano(), timeObj2.UnixNano(), reflect.Int64)
		}
	case reflect.Slice:
		{
			// We only care about the []byte type.
			if !canConvert(obj1Value, bytesType) {
				break
			}

			// []byte can be compared!
			bytesObj1, ok := obj1.([]byte)
			if !ok {
				bytesObj1 = obj1Value.Convert(bytesType).Interface().([]byte)

			}
			bytesObj2, ok := obj2.([]byte)
			if !ok {
				bytesObj2 = obj2Value.Convert(bytesType).Interface().([]byte)
			}

			return CompareType(bytes.Compare(bytesObj1, bytesObj2)), true
		}
	}

	return compareEqual, false
}

// Greater asserts that the first element is greater than the second
func (a *Assertions) Greater(e1 any, e2 any, msgAndArgs ...any) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return a.compareTwoValues(e1, e2, []CompareType{compareGreater}, "\"%v\" is not greater than \"%v\"", msgAndArgs...)
}

// GreaterOrEqual asserts that the first element is greater than or equal to the second
func (a *Assertions) GreaterOrEqual(e1 any, e2 any, msgAndArgs ...any) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return a.compareTwoValues(e1, e2, []CompareType{compareGreater, compareEqual}, "\"%v\" is not greater than or equal to \"%v\"", msgAndArgs...)
}

// Less asserts that the first element is less than the second
func (a *Assertions) Less(e1 any, e2 any, msgAndArgs ...any) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return a.compareTwoValues(e1, e2, []CompareType{compareLess}, "\"%v\" is not less than \"%v\"", msgAndArgs...)
}

// LessOrEqual asserts that the first element is less than or equal to the second
func (a *Assertions) LessOrEqual(e1 any, e2 any, msgAndArgs ...any) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	return a.compareTwoValues(e1, e2, []CompareType{compareLess, compareEqual}, "\"%v\" is not less than or equal to \"%v\"", msgAndArgs...)
}

// Positive asserts that the specified element is positive
func (a *Assertions) Positive(e any, msgAndArgs ...any) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	zero := reflect.Zero(reflect.TypeOf(e))
	return a.compareTwoValues(e, zero.Interface(), []CompareType{compareGreater}, "\"%v\" is not positive", msgAndArgs...)
}

// Negative asserts that the specified element is negative
func (a *Assertions) Negative(e any, msgAndArgs ...any) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}
	zero := reflect.Zero(reflect.TypeOf(e))
	return a.compareTwoValues(e, zero.Interface(), []CompareType{compareLess}, "\"%v\" is not negative", msgAndArgs...)
}

func (a *Assertions) compareTwoValues(e1 any, e2 any, allowedComparesResults []CompareType, failMessage string, msgAndArgs ...any) bool {
	if h, ok := a.t.(tHelper); ok {
		h.Helper()
	}

	e1Kind := reflect.ValueOf(e1).Kind()
	e2Kind := reflect.ValueOf(e2).Kind()
	if e1Kind != e2Kind {
		return a.Fail("Elements should be the same type", msgAndArgs...)
	}

	compareResult, isComparable := compare(e1, e2, e1Kind)
	if !isComparable {
		return a.Fail(fmt.Sprintf("Can not compare type \"%s\"", reflect.TypeOf(e1)), msgAndArgs...)
	}

	if !containsValue(allowedComparesResults, compareResult) {
		return a.Fail(fmt.Sprintf(failMessage, e1, e2), msgAndArgs...)
	}

	return true
}

func containsValue(values []CompareType, value CompareType) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}

	return false
}

// Wrapper around reflect.Value.CanConvert, for compatibility
// reasons.
func canConvert(value reflect.Value, to reflect.Type) bool {
	return value.CanConvert(to)
}
