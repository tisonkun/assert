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
	"runtime"
	"testing"
	"time"
)

func TestCompare(t *testing.T) {
	type customInt int
	type customInt8 int8
	type customInt16 int16
	type customInt32 int32
	type customInt64 int64
	type customUInt uint
	type customUInt8 uint8
	type customUInt16 uint16
	type customUInt32 uint32
	type customUInt64 uint64
	type customFloat32 float32
	type customFloat64 float64
	type customString string
	type customTime time.Time
	type customBytes []byte
	for _, currCase := range []struct {
		less    any
		greater any
		cType   string
	}{
		{less: customString("a"), greater: customString("b"), cType: "string"},
		{less: "a", greater: "b", cType: "string"},
		{less: customInt(1), greater: customInt(2), cType: "int"},
		{less: 1, greater: 2, cType: "int"},
		{less: customInt8(1), greater: customInt8(2), cType: "int8"},
		{less: int8(1), greater: int8(2), cType: "int8"},
		{less: customInt16(1), greater: customInt16(2), cType: "int16"},
		{less: int16(1), greater: int16(2), cType: "int16"},
		{less: customInt32(1), greater: customInt32(2), cType: "int32"},
		{less: int32(1), greater: int32(2), cType: "int32"},
		{less: customInt64(1), greater: customInt64(2), cType: "int64"},
		{less: int64(1), greater: int64(2), cType: "int64"},
		{less: customUInt(1), greater: customUInt(2), cType: "uint"},
		{less: uint8(1), greater: uint8(2), cType: "uint8"},
		{less: customUInt8(1), greater: customUInt8(2), cType: "uint8"},
		{less: uint16(1), greater: uint16(2), cType: "uint16"},
		{less: customUInt16(1), greater: customUInt16(2), cType: "uint16"},
		{less: uint32(1), greater: uint32(2), cType: "uint32"},
		{less: customUInt32(1), greater: customUInt32(2), cType: "uint32"},
		{less: uint64(1), greater: uint64(2), cType: "uint64"},
		{less: customUInt64(1), greater: customUInt64(2), cType: "uint64"},
		{less: float32(1.23), greater: float32(2.34), cType: "float32"},
		{less: customFloat32(1.23), greater: customFloat32(2.23), cType: "float32"},
		{less: 1.23, greater: 2.34, cType: "float64"},
		{less: customFloat64(1.23), greater: customFloat64(2.34), cType: "float64"},
		{less: time.Now(), greater: time.Now().Add(time.Hour), cType: "time.Time"},
		{less: customTime(time.Now()), greater: customTime(time.Now().Add(time.Hour)), cType: "time.Time"},
		{less: []byte{1, 1}, greater: []byte{1, 2}, cType: "[]byte"},
		{less: customBytes([]byte{1, 1}), greater: customBytes([]byte{1, 2}), cType: "[]byte"},
	} {
		resLess, isComparable := compare(currCase.less, currCase.greater, reflect.ValueOf(currCase.less).Kind())
		if !isComparable {
			t.Error("object should be comparable for type " + currCase.cType)
		}

		if resLess != compareLess {
			t.Errorf("object less (%v) should be less than greater (%v) for type "+currCase.cType,
				currCase.less, currCase.greater)
		}

		resGreater, isComparable := compare(currCase.greater, currCase.less, reflect.ValueOf(currCase.less).Kind())
		if !isComparable {
			t.Error("object are comparable for type " + currCase.cType)
		}

		if resGreater != compareGreater {
			t.Errorf("object greater should be greater than less for type " + currCase.cType)
		}

		resEqual, isComparable := compare(currCase.less, currCase.less, reflect.ValueOf(currCase.less).Kind())
		if !isComparable {
			t.Error("object are comparable for type " + currCase.cType)
		}

		if resEqual != 0 {
			t.Errorf("objects should be equal for type " + currCase.cType)
		}
	}
}

type outputT struct {
	buf     *bytes.Buffer
	helpers map[string]struct{}
}

// Implements TestingT
func (t *outputT) Errorf(format string, args ...any) {
	s := fmt.Sprintf(format, args...)
	t.buf.WriteString(s)
}

func (t *outputT) FailNow() {}

func (t *outputT) Helper() {
	if t.helpers == nil {
		t.helpers = make(map[string]struct{})
	}
	t.helpers[callerName(1)] = struct{}{}
}

// callerName gives the function name (qualified with a package path)
// for the caller after skip frames (where 0 means the current function).
func callerName(skip int) string {
	// Make room for the skip PC.
	var pc [1]uintptr
	n := runtime.Callers(skip+2, pc[:]) // skip + runtime.Callers + callerName
	if n == 0 {
		panic("testing: zero callers found")
	}
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame.Function
}

func TestGreater(t *testing.T) {
	mockAssertion := NewWithOnFailureNoop(new(testing.T))
	New(t).True(mockAssertion.Greater(2, 1))
	New(t).False(mockAssertion.Greater(1, 1))
	New(t).False(mockAssertion.Greater(1, 2))

	// Check error report
	for _, currCase := range []struct {
		less    any
		greater any
		msg     string
	}{
		{less: "a", greater: "b", msg: `"a" is not greater than "b"`},
		{less: 1, greater: 2, msg: `"1" is not greater than "2"`},
		{less: int8(1), greater: int8(2), msg: `"1" is not greater than "2"`},
		{less: int16(1), greater: int16(2), msg: `"1" is not greater than "2"`},
		{less: int32(1), greater: int32(2), msg: `"1" is not greater than "2"`},
		{less: int64(1), greater: int64(2), msg: `"1" is not greater than "2"`},
		{less: uint8(1), greater: uint8(2), msg: `"1" is not greater than "2"`},
		{less: uint16(1), greater: uint16(2), msg: `"1" is not greater than "2"`},
		{less: uint32(1), greater: uint32(2), msg: `"1" is not greater than "2"`},
		{less: uint64(1), greater: uint64(2), msg: `"1" is not greater than "2"`},
		{less: float32(1.23), greater: float32(2.34), msg: `"1.23" is not greater than "2.34"`},
		{less: 1.23, greater: 2.34, msg: `"1.23" is not greater than "2.34"`},
		{less: []byte{1, 1}, greater: []byte{1, 2}, msg: `"[1 1]" is not greater than "[1 2]"`},
		{less: time.Time{}, greater: time.Time{}.Add(time.Hour), msg: `"0001-01-01 00:00:00 +0000 UTC" is not greater than "0001-01-01 01:00:00 +0000 UTC"`},
	} {
		out := &outputT{buf: bytes.NewBuffer(nil)}
		outAssertion := New(out)
		New(t).False(outAssertion.Greater(currCase.less, currCase.greater))
		New(t).Contains(out.buf.String(), currCase.msg)
		New(t).Contains(out.helpers, "github.com/tisonkun/assert.(*Assertions).Greater")
	}
}

func TestGreaterOrEqual(t *testing.T) {
	mockAssertion := NewWithOnFailureNoop(new(testing.T))
	New(t).True(mockAssertion.GreaterOrEqual(2, 1))
	New(t).True(mockAssertion.GreaterOrEqual(1, 1))
	New(t).False(mockAssertion.GreaterOrEqual(1, 2))

	// Check error report
	for _, currCase := range []struct {
		less    any
		greater any
		msg     string
	}{
		{less: "a", greater: "b", msg: `"a" is not greater than or equal to "b"`},
		{less: 1, greater: 2, msg: `"1" is not greater than or equal to "2"`},
		{less: int8(1), greater: int8(2), msg: `"1" is not greater than or equal to "2"`},
		{less: int16(1), greater: int16(2), msg: `"1" is not greater than or equal to "2"`},
		{less: int32(1), greater: int32(2), msg: `"1" is not greater than or equal to "2"`},
		{less: int64(1), greater: int64(2), msg: `"1" is not greater than or equal to "2"`},
		{less: uint8(1), greater: uint8(2), msg: `"1" is not greater than or equal to "2"`},
		{less: uint16(1), greater: uint16(2), msg: `"1" is not greater than or equal to "2"`},
		{less: uint32(1), greater: uint32(2), msg: `"1" is not greater than or equal to "2"`},
		{less: uint64(1), greater: uint64(2), msg: `"1" is not greater than or equal to "2"`},
		{less: float32(1.23), greater: float32(2.34), msg: `"1.23" is not greater than or equal to "2.34"`},
		{less: 1.23, greater: 2.34, msg: `"1.23" is not greater than or equal to "2.34"`},
		{less: []byte{1, 1}, greater: []byte{1, 2}, msg: `"[1 1]" is not greater than or equal to "[1 2]"`},
		{less: time.Time{}, greater: time.Time{}.Add(time.Hour), msg: `"0001-01-01 00:00:00 +0000 UTC" is not greater than or equal to "0001-01-01 01:00:00 +0000 UTC"`},
	} {
		out := &outputT{buf: bytes.NewBuffer(nil)}
		outAssertion := New(out)
		New(t).False(outAssertion.GreaterOrEqual(currCase.less, currCase.greater))
		New(t).Contains(out.buf.String(), currCase.msg)
		New(t).Contains(out.helpers, "github.com/tisonkun/assert.(*Assertions).GreaterOrEqual")
	}
}

func TestLess(t *testing.T) {
	mockAssertion := NewWithOnFailureNoop(new(testing.T))
	New(t).False(mockAssertion.Less(2, 1))
	New(t).False(mockAssertion.Less(1, 1))
	New(t).True(mockAssertion.Less(1, 2))

	// Check error report
	for _, currCase := range []struct {
		less    any
		greater any
		msg     string
	}{
		{less: "a", greater: "b", msg: `"b" is not less than "a"`},
		{less: 1, greater: 2, msg: `"2" is not less than "1"`},
		{less: int8(1), greater: int8(2), msg: `"2" is not less than "1"`},
		{less: int16(1), greater: int16(2), msg: `"2" is not less than "1"`},
		{less: int32(1), greater: int32(2), msg: `"2" is not less than "1"`},
		{less: int64(1), greater: int64(2), msg: `"2" is not less than "1"`},
		{less: uint8(1), greater: uint8(2), msg: `"2" is not less than "1"`},
		{less: uint16(1), greater: uint16(2), msg: `"2" is not less than "1"`},
		{less: uint32(1), greater: uint32(2), msg: `"2" is not less than "1"`},
		{less: uint64(1), greater: uint64(2), msg: `"2" is not less than "1"`},
		{less: float32(1.23), greater: float32(2.34), msg: `"2.34" is not less than "1.23"`},
		{less: 1.23, greater: 2.34, msg: `"2.34" is not less than "1.23"`},
		{less: []byte{1, 1}, greater: []byte{1, 2}, msg: `"[1 2]" is not less than "[1 1]"`},
		{less: time.Time{}, greater: time.Time{}.Add(time.Hour), msg: `"0001-01-01 01:00:00 +0000 UTC" is not less than "0001-01-01 00:00:00 +0000 UTC"`},
	} {
		out := &outputT{buf: bytes.NewBuffer(nil)}
		outAssertion := New(out)
		New(t).False(outAssertion.Less(currCase.greater, currCase.less))
		New(t).Contains(out.buf.String(), currCase.msg)
		New(t).Contains(out.helpers, "github.com/tisonkun/assert.(*Assertions).Less")
	}
}

func TestLessOrEqual(t *testing.T) {
	mockAssertion := NewWithOnFailureNoop(new(testing.T))
	New(t).False(mockAssertion.LessOrEqual(2, 1))
	New(t).True(mockAssertion.LessOrEqual(1, 1))
	New(t).True(mockAssertion.LessOrEqual(1, 2))

	// Check error report
	for _, currCase := range []struct {
		less    any
		greater any
		msg     string
	}{
		{less: "a", greater: "b", msg: `"b" is not less than or equal to "a"`},
		{less: 1, greater: 2, msg: `"2" is not less than or equal to "1"`},
		{less: int8(1), greater: int8(2), msg: `"2" is not less than or equal to "1"`},
		{less: int16(1), greater: int16(2), msg: `"2" is not less than or equal to "1"`},
		{less: int32(1), greater: int32(2), msg: `"2" is not less than or equal to "1"`},
		{less: int64(1), greater: int64(2), msg: `"2" is not less than or equal to "1"`},
		{less: uint8(1), greater: uint8(2), msg: `"2" is not less than or equal to "1"`},
		{less: uint16(1), greater: uint16(2), msg: `"2" is not less than or equal to "1"`},
		{less: uint32(1), greater: uint32(2), msg: `"2" is not less than or equal to "1"`},
		{less: uint64(1), greater: uint64(2), msg: `"2" is not less than or equal to "1"`},
		{less: float32(1.23), greater: float32(2.34), msg: `"2.34" is not less than or equal to "1.23"`},
		{less: 1.23, greater: 2.34, msg: `"2.34" is not less than or equal to "1.23"`},
		{less: []byte{1, 1}, greater: []byte{1, 2}, msg: `"[1 2]" is not less than or equal to "[1 1]"`},
		{less: time.Time{}, greater: time.Time{}.Add(time.Hour), msg: `"0001-01-01 01:00:00 +0000 UTC" is not less than or equal to "0001-01-01 00:00:00 +0000 UTC"`},
	} {
		out := &outputT{buf: bytes.NewBuffer(nil)}
		outAssertion := New(out)
		New(t).False(outAssertion.LessOrEqual(currCase.greater, currCase.less))
		New(t).Contains(out.buf.String(), currCase.msg)
		New(t).Contains(out.helpers, "github.com/tisonkun/assert.(*Assertions).LessOrEqual")
	}
}

func TestPositive(t *testing.T) {
	mockAssertion := NewWithOnFailureNoop(new(testing.T))
	New(t).True(mockAssertion.Positive(1))
	New(t).True(mockAssertion.Positive(1.23))
	New(t).False(mockAssertion.Positive(0))
	New(t).False(mockAssertion.Positive(-1))
	New(t).False(mockAssertion.Positive(-1.23))

	// Check error report
	for _, currCase := range []struct {
		e   any
		msg string
	}{
		{e: -1, msg: `"-1" is not positive`},
		{e: int8(-1), msg: `"-1" is not positive`},
		{e: int16(-1), msg: `"-1" is not positive`},
		{e: int32(-1), msg: `"-1" is not positive`},
		{e: int64(-1), msg: `"-1" is not positive`},
		{e: float32(-1.23), msg: `"-1.23" is not positive`},
		{e: -1.23, msg: `"-1.23" is not positive`},
	} {
		out := &outputT{buf: bytes.NewBuffer(nil)}
		outAssertion := New(out)
		New(t).False(outAssertion.Positive(currCase.e))
		New(t).Contains(out.buf.String(), currCase.msg)
		New(t).Contains(out.helpers, "github.com/tisonkun/assert.(*Assertions).Positive")
	}
}

func TestNegative(t *testing.T) {
	mockAssertion := NewWithOnFailureNoop(new(testing.T))
	New(t).False(mockAssertion.Negative(1))
	New(t).False(mockAssertion.Negative(1.23))
	New(t).False(mockAssertion.Negative(0))
	New(t).True(mockAssertion.Negative(-1))
	New(t).True(mockAssertion.Negative(-1.23))

	// Check error report
	for _, currCase := range []struct {
		e   any
		msg string
	}{
		{e: 1, msg: `"1" is not negative`},
		{e: int8(1), msg: `"1" is not negative`},
		{e: int16(1), msg: `"1" is not negative`},
		{e: int32(1), msg: `"1" is not negative`},
		{e: int64(1), msg: `"1" is not negative`},
		{e: float32(1.23), msg: `"1.23" is not negative`},
		{e: 1.23, msg: `"1.23" is not negative`},
	} {
		out := &outputT{buf: bytes.NewBuffer(nil)}
		outAssertion := New(out)
		New(t).False(outAssertion.Negative(currCase.e))
		New(t).Contains(out.buf.String(), currCase.msg)
		New(t).Contains(out.helpers, "github.com/tisonkun/assert.(*Assertions).Negative")
	}
}

func Test_compareTwoValuesDifferentValuesTypes(t *testing.T) {
	mockAssertion := NewWithOnFailureNoop(new(testing.T))
	for _, currCase := range []struct {
		v1            any
		v2            any
		compareResult bool
	}{
		{v1: 123, v2: "abc"},
		{v1: "abc", v2: 123456},
		{v1: float64(12), v2: "123"},
		{v1: "float(12)", v2: float64(1)},
	} {
		New(t).False(mockAssertion.compareTwoValues(currCase.v1, currCase.v2, []CompareType{compareLess, compareEqual, compareGreater}, "testFailMessage"))
	}
}

func Test_compareTwoValuesNotComparableValues(t *testing.T) {
	mockAssertion := NewWithOnFailureNoop(new(testing.T))
	type CompareStruct struct{}
	for _, currCase := range []struct {
		v1 any
		v2 any
	}{
		{v1: CompareStruct{}, v2: CompareStruct{}},
		{v1: map[string]int{}, v2: map[string]int{}},
		{v1: make([]int, 5), v2: make([]int, 5)},
	} {
		New(t).False(mockAssertion.compareTwoValues(currCase.v1, currCase.v2, []CompareType{compareLess, compareEqual, compareGreater}, "testFailMessage"))
	}
}

func Test_compareTwoValuesCorrectCompareResult(t *testing.T) {
	mockAssertion := NewWithOnFailureNoop(new(testing.T))
	for _, currCase := range []struct {
		v1           any
		v2           any
		compareTypes []CompareType
	}{
		{v1: 1, v2: 2, compareTypes: []CompareType{compareLess}},
		{v1: 1, v2: 2, compareTypes: []CompareType{compareLess, compareEqual}},
		{v1: 2, v2: 2, compareTypes: []CompareType{compareGreater, compareEqual}},
		{v1: 2, v2: 2, compareTypes: []CompareType{compareEqual}},
		{v1: 2, v2: 1, compareTypes: []CompareType{compareEqual, compareGreater}},
		{v1: 2, v2: 1, compareTypes: []CompareType{compareGreater}},
	} {
		New(t).True(mockAssertion.compareTwoValues(currCase.v1, currCase.v2, currCase.compareTypes, "testFailMessage"))
	}
}

func Test_containsValue(t *testing.T) {
	for _, currCase := range []struct {
		values []CompareType
		value  CompareType
		result bool
	}{
		{values: []CompareType{compareGreater}, value: compareGreater, result: true},
		{values: []CompareType{compareGreater, compareLess}, value: compareGreater, result: true},
		{values: []CompareType{compareGreater, compareLess}, value: compareLess, result: true},
		{values: []CompareType{compareGreater, compareLess}, value: compareEqual, result: false},
	} {
		compareResult := containsValue(currCase.values, currCase.value)
		New(t).Equal(currCase.result, compareResult)
	}
}

func TestComparingMsgAndArgsForwarding(t *testing.T) {
	msgAndArgs := []any{"format %s %x", "this", 0xc001}
	expectedOutput := "format this c001\n"
	funcs := []func(*Assertions){
		func(a *Assertions) { a.Greater(1, 2, msgAndArgs...) },
		func(a *Assertions) { a.GreaterOrEqual(1, 2, msgAndArgs...) },
		func(a *Assertions) { a.Less(2, 1, msgAndArgs...) },
		func(a *Assertions) { a.LessOrEqual(2, 1, msgAndArgs...) },
		func(a *Assertions) { a.Positive(0, msgAndArgs...) },
		func(a *Assertions) { a.Negative(0, msgAndArgs...) },
	}
	for _, f := range funcs {
		out := &outputT{buf: bytes.NewBuffer(nil)}
		outAssertion := New(out)
		f(outAssertion)
		New(t).Contains(out.buf.String(), expectedOutput)
	}
}
