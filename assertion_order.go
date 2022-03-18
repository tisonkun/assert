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
	"fmt"
	"reflect"
)

// isOrdered checks that collection contains elements in order.
func (a *Assertions) isOrdered(object any, allowedComparesResults []CompareType, failMessage string, msgAndArgs ...any) bool {
	objKind := reflect.TypeOf(object).Kind()
	if objKind != reflect.Slice && objKind != reflect.Array {
		return a.Fail(fmt.Sprintf("Can not test elements in order for type \"%s\"", objKind), msgAndArgs...)
	}

	objValue := reflect.ValueOf(object)
	objLen := objValue.Len()

	if objLen <= 1 {
		return true
	}

	value := objValue.Index(0)
	valueInterface := value.Interface()
	firstValueKind := value.Kind()

	for i := 1; i < objLen; i++ {
		prevValue := value
		prevValueInterface := valueInterface

		value = objValue.Index(i)
		valueInterface = value.Interface()

		compareResult, isComparable := compare(prevValueInterface, valueInterface, firstValueKind)

		if !isComparable {
			return a.Fail(fmt.Sprintf("Can not compare type \"%s\" and \"%s\"", reflect.TypeOf(value), reflect.TypeOf(prevValue)), msgAndArgs...)
		}

		if !containsValue(allowedComparesResults, compareResult) {
			return a.Fail(fmt.Sprintf(failMessage, prevValue, value), msgAndArgs...)
		}
	}

	return true
}

// IsIncreasing asserts that the collection is increasing
func (a *Assertions) IsIncreasing(object any, msgAndArgs ...any) bool {
	return a.isOrdered(object, []CompareType{compareLess}, "\"%v\" is not less than \"%v\"", msgAndArgs...)
}

// IsNonIncreasing asserts that the collection is not increasing
func (a *Assertions) IsNonIncreasing(object any, msgAndArgs ...any) bool {
	return a.isOrdered(object, []CompareType{compareEqual, compareGreater}, "\"%v\" is not greater than or equal to \"%v\"", msgAndArgs...)
}

// IsDecreasing asserts that the collection is decreasing
func (a *Assertions) IsDecreasing(object any, msgAndArgs ...any) bool {
	return a.isOrdered(object, []CompareType{compareGreater}, "\"%v\" is not greater than \"%v\"", msgAndArgs...)
}

// IsNonDecreasing asserts that the collection is not decreasing
func (a *Assertions) IsNonDecreasing(object any, msgAndArgs ...any) bool {
	return a.isOrdered(object, []CompareType{compareLess, compareEqual}, "\"%v\" is not less than or equal to \"%v\"", msgAndArgs...)
}
