package assert

import (
	"fmt"
	"reflect"
)

// isOrdered checks that collection contains elements in order.
func (a *Assertions) isOrdered(object any, allowedComparesResults []CompareType, failMessage string, msgAndArgs ...any) {
	objKind := reflect.TypeOf(object).Kind()
	if objKind != reflect.Slice && objKind != reflect.Array {
		a.Fail(fmt.Sprintf("Can not test elements in order for type \"%s\"", objKind), msgAndArgs...)
		return
	}

	objValue := reflect.ValueOf(object)
	objLen := objValue.Len()

	if objLen <= 1 {
		return
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
			a.Fail(fmt.Sprintf("Can not compare type \"%s\" and \"%s\"", reflect.TypeOf(value), reflect.TypeOf(prevValue)), msgAndArgs...)
			return
		}

		if !containsValue(allowedComparesResults, compareResult) {
			a.Fail(fmt.Sprintf(failMessage, prevValue, value), msgAndArgs...)
			return
		}
	}
}

// IsIncreasing asserts that the collection is increasing
func (a *Assertions) IsIncreasing(object any, msgAndArgs ...any) {
	a.isOrdered(object, []CompareType{compareLess}, "\"%v\" is not less than \"%v\"", msgAndArgs...)
}

// IsNonIncreasing asserts that the collection is not increasing
func (a *Assertions) IsNonIncreasing(object any, msgAndArgs ...any) {
	a.isOrdered(object, []CompareType{compareEqual, compareGreater}, "\"%v\" is not greater than or equal to \"%v\"", msgAndArgs...)
}

// IsDecreasing asserts that the collection is decreasing
func (a *Assertions) IsDecreasing(object any, msgAndArgs ...any) {
	a.isOrdered(object, []CompareType{compareGreater}, "\"%v\" is not greater than \"%v\"", msgAndArgs...)
}

// IsNonDecreasing asserts that the collection is not decreasing
func (a *Assertions) IsNonDecreasing(object any, msgAndArgs ...any) {
	a.isOrdered(object, []CompareType{compareLess, compareEqual}, "\"%v\" is not less than or equal to \"%v\"", msgAndArgs...)
}
