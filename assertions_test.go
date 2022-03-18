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
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"testing"
	"time"
)

func NewWithOnFailureNoop(t TestingT) *Assertions {
	return New(t).WithOnFailure(func(TestingT) {})
}

var (
	i     any
	zeros = []any{
		false,
		byte(0),
		complex64(0),
		complex128(0),
		float32(0),
		float64(0),
		0,
		int8(0),
		int16(0),
		int32(0),
		int64(0),
		rune(0),
		uint(0),
		uint8(0),
		uint16(0),
		uint32(0),
		uint64(0),
		uintptr(0),
		"",
		[0]any{},
		[]any(nil),
		struct{ x int }{},
		(*any)(nil),
		(func())(nil),
		nil,
		any(nil),
		map[any]any(nil),
		(chan any)(nil),
		(<-chan any)(nil),
		(chan<- any)(nil),
	}
	nonZeros = []any{
		true,
		byte(1),
		complex64(1),
		complex128(1),
		float32(1),
		float64(1),
		1,
		int8(1),
		int16(1),
		int32(1),
		int64(1),
		rune(1),
		uint(1),
		uint8(1),
		uint16(1),
		uint32(1),
		uint64(1),
		uintptr(1),
		"s",
		[1]any{1},
		[]any{},
		struct{ x int }{1},
		&i,
		func() {},
		any(1),
		map[any]any{},
		make(chan any),
		(<-chan any)(make(chan any)),
		(chan<- any)(make(chan any)),
	}
)

// AssertionTesterInterface defines an interface to be used for testing assertion methods
type AssertionTesterInterface interface {
	TestMethod()
}

// AssertionTesterConformingObject is an object that conforms to the AssertionTesterInterface interface
type AssertionTesterConformingObject struct {
}

func (a *AssertionTesterConformingObject) TestMethod() {
}

// AssertionTesterNonConformingObject is an object that does not conform to the AssertionTesterInterface interface
type AssertionTesterNonConformingObject struct {
}

func TestObjectsAreEqual(t *testing.T) {
	cases := []struct {
		expected any
		actual   any
		result   bool
	}{
		// cases that are expected to be equal
		{"Hello World", "Hello World", true},
		{123, 123, true},
		{123.5, 123.5, true},
		{[]byte("Hello World"), []byte("Hello World"), true},
		{nil, nil, true},

		// cases that are expected not to be equal
		{map[int]int{5: 10}, map[int]int{10: 20}, false},
		{'x', "x", false},
		{"x", 'x', false},
		{0, 0.1, false},
		{0.1, 0, false},
		{time.Now, time.Now, false},
		{func() {}, func() {}, false},
		{uint32(10), int32(10), false},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("ObjectsAreEqual(%#v, %#v)", c.expected, c.actual), func(t *testing.T) {
			res := ObjectsAreEqual(c.expected, c.actual)

			if res != c.result {
				t.Errorf("ObjectsAreEqual(%#v, %#v) should return %#v", c.expected, c.actual, c.result)
			}

		})
	}

	// Cases where type differ but values are equal
	if !ObjectsAreEqualValues(uint32(10), int32(10)) {
		t.Error("ObjectsAreEqualValues should return true")
	}
	if ObjectsAreEqualValues(0, nil) {
		t.Fail()
	}
	if ObjectsAreEqualValues(nil, 0) {
		t.Fail()
	}

}

func TestImplements(t *testing.T) {

	mockT := new(testing.T)

	if !Implements(mockT, (*AssertionTesterInterface)(nil), new(AssertionTesterConformingObject)) {
		t.Error("Implements method should return true: AssertionTesterConformingObject implements AssertionTesterInterface")
	}
	if Implements(mockT, (*AssertionTesterInterface)(nil), new(AssertionTesterNonConformingObject)) {
		t.Error("Implements method should return false: AssertionTesterNonConformingObject does not implements AssertionTesterInterface")
	}
	if Implements(mockT, (*AssertionTesterInterface)(nil), nil) {
		t.Error("Implements method should return false: nil does not implement AssertionTesterInterface")
	}

}

func TestIsType(t *testing.T) {

	mockT := new(testing.T)

	if !IsType(mockT, new(AssertionTesterConformingObject), new(AssertionTesterConformingObject)) {
		t.Error("IsType should return true: AssertionTesterConformingObject is the same type as AssertionTesterConformingObject")
	}
	if IsType(mockT, new(AssertionTesterConformingObject), new(AssertionTesterNonConformingObject)) {
		t.Error("IsType should return false: AssertionTesterConformingObject is not the same type as AssertionTesterNonConformingObject")
	}

}

func TestEqual(t *testing.T) {
	type myType string

	mockAssertion := NewWithOnFailureNoop(new(testing.T))
	var m map[string]any

	cases := []struct {
		expected any
		actual   any
		result   bool
		remark   string
	}{
		{"Hello World", "Hello World", true, ""},
		{123, 123, true, ""},
		{123.5, 123.5, true, ""},
		{[]byte("Hello World"), []byte("Hello World"), true, ""},
		{nil, nil, true, ""},
		{int32(123), int32(123), true, ""},
		{uint64(123), uint64(123), true, ""},
		{myType("1"), myType("1"), true, ""},
		{&struct{}{}, &struct{}{}, true, "pointer equality is based on equality of underlying value"},

		// Not expected to be equal
		{m["bar"], "something", false, ""},
		{myType("1"), myType("2"), false, ""},

		// A case that might be confusing, especially with numeric literals
		{10, uint(10), false, ""},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("Equal(%#v, %#v)", c.expected, c.actual), func(t *testing.T) {
			res := mockAssertion.Equal(c.expected, c.actual)

			if res != c.result {
				t.Errorf("Equal(%#v, %#v) should return %#v: %s", c.expected, c.actual, c.result, c.remark)
			}
		})
	}
}

func ptr(i int) *int {
	return &i
}

func TestSame(t *testing.T) {

	mockT := new(testing.T)

	if Same(mockT, ptr(1), ptr(1)) {
		t.Error("Same should return false")
	}
	if Same(mockT, 1, 1) {
		t.Error("Same should return false")
	}
	p := ptr(2)
	if Same(mockT, p, *p) {
		t.Error("Same should return false")
	}
	if !Same(mockT, p, p) {
		t.Error("Same should return true")
	}
}

func TestNotSame(t *testing.T) {
	mockT := new(testing.T)

	if !NotSame(mockT, ptr(1), ptr(1)) {
		t.Error("NotSame should return true; different pointers")
	}
	if !NotSame(mockT, 1, 1) {
		t.Error("NotSame should return true; constant inputs")
	}
	p := ptr(2)
	if !NotSame(mockT, p, *p) {
		t.Error("NotSame should return true; mixed-type inputs")
	}
	if NotSame(mockT, p, p) {
		t.Error("NotSame should return false")
	}
}

func Test_samePointers(t *testing.T) {
	p := ptr(2)

	type args struct {
		first  any
		second any
	}
	tests := []struct {
		name      string
		args      args
		assertion func(*Assertions, bool, ...any) bool
	}{
		{
			name:      "1 != 2",
			args:      args{first: 1, second: 2},
			assertion: (*Assertions).False,
		},
		{
			name:      "1 != 1 (not same ptr)",
			args:      args{first: 1, second: 1},
			assertion: (*Assertions).False,
		},
		{
			name:      "ptr(1) == ptr(1)",
			args:      args{first: p, second: p},
			assertion: (*Assertions).True,
		},
		{
			name:      "int(1) != float32(1)",
			args:      args{first: 1, second: float32(1)},
			assertion: (*Assertions).False,
		},
		{
			name:      "array != slice",
			args:      args{first: [2]int{1, 2}, second: []int{1, 2}},
			assertion: (*Assertions).False,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertion := New(t)
			tt.assertion(assertion, samePointers(tt.args.first, tt.args.second))
		})
	}
}

// bufferT implements TestingT. Its implementation of Errorf writes the output that would be produced by
// testing.T.Errorf to an internal bytes.Buffer.
type bufferT struct {
	buf bytes.Buffer
}

func (t *bufferT) Errorf(format string, args ...any) {
	// implementation of decorate is copied from testing.T
	decorate := func(s string) string {
		_, file, line, ok := runtime.Caller(3) // decorate + log + public function.
		if ok {
			// Truncate file name at last file name separator.
			if index := strings.LastIndex(file, "/"); index >= 0 {
				file = file[index+1:]
			} else if index = strings.LastIndex(file, "\\"); index >= 0 {
				file = file[index+1:]
			}
		} else {
			file = "???"
			line = 1
		}
		buf := new(bytes.Buffer)
		// Every line is indented at least one tab.
		buf.WriteByte('\t')
		_, _ = fmt.Fprintf(buf, "%s:%d: ", file, line)
		lines := strings.Split(s, "\n")
		if l := len(lines); l > 1 && lines[l-1] == "" {
			lines = lines[:l-1]
		}
		for i, line := range lines {
			if i > 0 {
				// Second and subsequent lines are indented an extra tab.
				buf.WriteString("\n\t\t")
			}
			buf.WriteString(line)
		}
		buf.WriteByte('\n')
		return buf.String()
	}
	t.buf.WriteString(decorate(fmt.Sprintf(format, args...)))
}

func (t *bufferT) FailNow() {}

func TestStringEqual(t *testing.T) {
	for i, currCase := range []struct {
		equalWant  string
		equalGot   string
		msgAndArgs []any
		want       string
	}{
		{equalWant: "hi, \nmy name is", equalGot: "what,\nmy name is", want: "\tassertions.go:\\d+: \n\t+Error Trace:\t\n\t+Error:\\s+Not equal:\\s+\n\\s+expected: \"hi, \\\\nmy name is\"\n\\s+actual\\s+: \"what,\\\\nmy name is\"\n\\s+Diff:\n\\s+-+ Expected\n\\s+\\++ Actual\n\\s+@@ -1,2 \\+1,2 @@\n\\s+-hi, \n\\s+\\+what,\n\\s+my name is"},
	} {
		mockT := &bufferT{}
		New(mockT).Equal(currCase.equalWant, currCase.equalGot, currCase.msgAndArgs...)
		Regexp(t, regexp.MustCompile(currCase.want), mockT.buf.String(), "Case %d", i)
	}
}

func TestEqualFormatting(t *testing.T) {
	for i, currCase := range []struct {
		equalWant  string
		equalGot   string
		msgAndArgs []any
		want       string
	}{
		{equalWant: "want", equalGot: "got", want: "\tassertions.go:\\d+: \n\t+Error Trace:\t\n\t+Error:\\s+Not equal:\\s+\n\\s+expected: \"want\"\n\\s+actual\\s+: \"got\"\n\\s+Diff:\n\\s+-+ Expected\n\\s+\\++ Actual\n\\s+@@ -1 \\+1 @@\n\\s+-want\n\\s+\\+got\n"},
		{equalWant: "want", equalGot: "got", msgAndArgs: []any{"hello, %v!", "world"}, want: "\tassertions.go:[0-9]+: \n\t+Error Trace:\t\n\t+Error:\\s+Not equal:\\s+\n\\s+expected: \"want\"\n\\s+actual\\s+: \"got\"\n\\s+Diff:\n\\s+-+ Expected\n\\s+\\++ Actual\n\\s+@@ -1 \\+1 @@\n\\s+-want\n\\s+\\+got\n\\s+Messages:\\s+hello, world!\n"},
		{equalWant: "want", equalGot: "got", msgAndArgs: []any{123}, want: "\tassertions.go:[0-9]+: \n\t+Error Trace:\t\n\t+Error:\\s+Not equal:\\s+\n\\s+expected: \"want\"\n\\s+actual\\s+: \"got\"\n\\s+Diff:\n\\s+-+ Expected\n\\s+\\++ Actual\n\\s+@@ -1 \\+1 @@\n\\s+-want\n\\s+\\+got\n\\s+Messages:\\s+123\n"},
		{equalWant: "want", equalGot: "got", msgAndArgs: []any{struct{ a string }{"hello"}}, want: "\tassertions.go:[0-9]+: \n\t+Error Trace:\t\n\t+Error:\\s+Not equal:\\s+\n\\s+expected: \"want\"\n\\s+actual\\s+: \"got\"\n\\s+Diff:\n\\s+-+ Expected\n\\s+\\++ Actual\n\\s+@@ -1 \\+1 @@\n\\s+-want\n\\s+\\+got\n\\s+Messages:\\s+{a:hello}\n"},
	} {
		mockT := &bufferT{}
		New(mockT).Equal(currCase.equalWant, currCase.equalGot, currCase.msgAndArgs...)
		Regexp(t, regexp.MustCompile(currCase.want), mockT.buf.String(), "Case %d", i)
	}
}

func TestFormatUnequalValues(t *testing.T) {
	expected, actual := formatUnequalValues("foo", "bar")
	New(t).Equal(`"foo"`, expected, "value should not include type")
	New(t).Equal(`"bar"`, actual, "value should not include type")

	expected, actual = formatUnequalValues(123, 123)
	New(t).Equal(`123`, expected, "value should not include type")
	New(t).Equal(`123`, actual, "value should not include type")

	expected, actual = formatUnequalValues(int64(123), int32(123))
	New(t).Equal(`int64(123)`, expected, "value should include type")
	New(t).Equal(`int32(123)`, actual, "value should include type")

	expected, actual = formatUnequalValues(int64(123), nil)
	New(t).Equal(`int64(123)`, expected, "value should include type")
	New(t).Equal(`<nil>(<nil>)`, actual, "value should include type")

	type testStructType struct {
		Val string
	}

	expected, actual = formatUnequalValues(&testStructType{Val: "test"}, &testStructType{Val: "test"})
	New(t).Equal(`&assert.testStructType{Val:"test"}`, expected, "value should not include type annotation")
	New(t).Equal(`&assert.testStructType{Val:"test"}`, actual, "value should not include type annotation")
}

func TestNotNil(t *testing.T) {
	mockAssertion := NewWithOnFailureNoop(new(testing.T))
	if !mockAssertion.NotNil(new(AssertionTesterConformingObject)) {
		t.Error("NotNil should return nil: object is not nil")
	}
	if mockAssertion.NotNil(nil) {
		t.Error("NotNil should return not-nil value: object is nil")
	}
	if mockAssertion.NotNil((*struct{})(nil)) {
		t.Error("NotNil should return not-nil value: object is (*struct{})(nil)")
	}
}

func TestNil(t *testing.T) {
	mockAssertion := NewWithOnFailureNoop(new(testing.T))
	if !mockAssertion.Nil(nil) {
		t.Error("Nil should return nil: object is nil")
	}
	if !mockAssertion.Nil((*struct{})(nil)) {
		t.Error("Nil should return nil: object is (*struct{})(nil)")
	}
	if mockAssertion.Nil(new(AssertionTesterConformingObject)) {
		t.Error("Nil should return not-nil value: object is not nil")
	}
}

func TestTrue(t *testing.T) {
	mockAssertion := NewWithOnFailureNoop(new(testing.T))
	if !mockAssertion.True(true) {
		t.Error("True should return true")
	}
	if mockAssertion.True(false) {
		t.Error("True should return false")
	}
}

func TestFalse(t *testing.T) {
	mockAssertion := NewWithOnFailureNoop(new(testing.T))
	if !mockAssertion.False(false) {
		t.Error("False should return true")
	}
	if mockAssertion.False(true) {
		t.Error("False should return false")
	}
}

func TestExactly(t *testing.T) {
	mockAssertion := NewWithOnFailureNoop(new(testing.T))

	a := float32(1)
	b := float64(1)
	c := float32(1)
	d := float32(2)
	cases := []struct {
		expected any
		actual   any
		result   bool
	}{
		{a, b, false},
		{a, d, false},
		{a, c, true},
		{nil, a, false},
		{a, nil, false},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("Exactly(%#v, %#v)", c.expected, c.actual), func(t *testing.T) {
			res := mockAssertion.Exactly(c.expected, c.actual)

			if res != c.result {
				t.Errorf("Exactly(%#v, %#v) should return %#v", c.expected, c.actual, c.result)
			}
		})
	}
}

func TestNotEqual(t *testing.T) {
	mockAssertion := NewWithOnFailureNoop(new(testing.T))

	cases := []struct {
		expected any
		actual   any
		result   bool
	}{
		// cases that are expected not to match
		{"Hello World", "Hello World!", true},
		{123, 1234, true},
		{123.5, 123.55, true},
		{[]byte("Hello World"), []byte("Hello World!"), true},
		{nil, new(AssertionTesterConformingObject), true},

		// cases that are expected to match
		{nil, nil, false},
		{"Hello World", "Hello World", false},
		{123, 123, false},
		{123.5, 123.5, false},
		{[]byte("Hello World"), []byte("Hello World"), false},
		{new(AssertionTesterConformingObject), new(AssertionTesterConformingObject), false},
		{&struct{}{}, &struct{}{}, false},
		{func() int { return 23 }, func() int { return 24 }, false},
		// A case that might be confusing, especially with numeric literals
		{10, uint(10), true},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("NotEqual(%#v, %#v)", c.expected, c.actual), func(t *testing.T) {
			res := mockAssertion.NotEqual(c.expected, c.actual)

			if res != c.result {
				t.Errorf("NotEqual(%#v, %#v) should return %#v", c.expected, c.actual, c.result)
			}
		})
	}
}

func TestNotEqualValues(t *testing.T) {
	mockAssertion := NewWithOnFailureNoop(new(testing.T))

	cases := []struct {
		expected any
		actual   any
		result   bool
	}{
		// cases that are expected not to match
		{"Hello World", "Hello World!", true},
		{123, 1234, true},
		{123.5, 123.55, true},
		{[]byte("Hello World"), []byte("Hello World!"), true},
		{nil, new(AssertionTesterConformingObject), true},

		// cases that are expected to match
		{nil, nil, false},
		{"Hello World", "Hello World", false},
		{123, 123, false},
		{123.5, 123.5, false},
		{[]byte("Hello World"), []byte("Hello World"), false},
		{new(AssertionTesterConformingObject), new(AssertionTesterConformingObject), false},
		{&struct{}{}, &struct{}{}, false},

		// Different behaviour from NotEqual()
		{func() int { return 23 }, func() int { return 24 }, true},
		{10, 11, true},
		{10, uint(10), false},

		{struct{}{}, struct{}{}, false},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("NotEqualValues(%#v, %#v)", c.expected, c.actual), func(t *testing.T) {
			if mockAssertion.NotEqualValues(c.expected, c.actual) != c.result {
				t.Errorf("NotEqualValues(%#v, %#v) should return %#v", c.expected, c.actual, c.result)
			}
			if mockAssertion.EqualValues(c.expected, c.actual) == c.result {
				t.Errorf("EqualValues(%#v, %#v) should return %#v", c.expected, c.actual, !c.result)
			}
		})
	}
}

func TestContainsNotContains(t *testing.T) {
	type A struct {
		Name, Value string
	}
	list := []string{"Foo", "Bar"}

	complexList := []*A{
		{"b", "c"},
		{"d", "e"},
		{"g", "h"},
		{"j", "k"},
	}
	simpleMap := map[any]any{"Foo": "Bar"}
	var zeroMap map[any]any

	cases := []struct {
		expected any
		actual   any
		result   bool
	}{
		{"Hello World", "Hello", true},
		{"Hello World", "Salut", false},
		{list, "Bar", true},
		{list, "Salut", false},
		{complexList, &A{"g", "h"}, true},
		{complexList, &A{"g", "e"}, false},
		{simpleMap, "Foo", true},
		{simpleMap, "Bar", false},
		{zeroMap, "Bar", false},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("Contains(%#v, %#v)", c.expected, c.actual), func(t *testing.T) {
			mockAssertion := NewWithOnFailureNoop(new(testing.T))
			res := mockAssertion.Contains(c.expected, c.actual)

			if res != c.result {
				if res {
					t.Errorf("Contains(%#v, %#v) should return nil:\n\t%#v contains %#v", c.expected, c.actual, c.expected, c.actual)
				} else {
					t.Errorf("Contains(%#v, %#v) should return not-nil value:\n\t%#v does not contain %#v", c.expected, c.actual, c.expected, c.actual)
				}
			}
		})
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("NotContains(%#v, %#v)", c.expected, c.actual), func(t *testing.T) {
			mockAssertion := NewWithOnFailureNoop(new(testing.T))
			res := mockAssertion.NotContains(c.expected, c.actual)

			// NotContains should be inverse of Contains. If it's not, something is wrong
			if res == mockAssertion.Contains(c.expected, c.actual) {
				if res {
					t.Errorf("NotContains(%#v, %#v) should return nil:\n\t%#v does not contains %#v", c.expected, c.actual, c.expected, c.actual)
				} else {
					t.Errorf("NotContains(%#v, %#v) should return not-nil value:\n\t%#v contains %#v", c.expected, c.actual, c.expected, c.actual)
				}
			}
		})
	}
}

func TestContainsFailMessage(t *testing.T) {
	out := &outputT{buf: bytes.NewBuffer(nil)}
	outAssertion := New(out)
	outAssertion.Contains("hello world", errors.New("hello"))
	expectedFail := "\"hello world\" does not contain &errors.errorString{s:\"hello\"}"
	actualFail := out.buf.String()
	if !strings.Contains(actualFail, expectedFail) {
		t.Errorf("Contains failure should include %q but was %q", expectedFail, actualFail)
	}
}

func TestContainsNotContainsOnNilValue(t *testing.T) {
	out := &outputT{buf: bytes.NewBuffer(nil)}
	outAssertion := New(out)
	outAssertion.Contains(nil, "key")
	expectedFail := "<nil> could not be applied builtin len()"
	actualFail := out.buf.String()
	if !strings.Contains(actualFail, expectedFail) {
		t.Errorf("Contains failure should include %q but was %q", expectedFail, actualFail)
	}

	out = &outputT{buf: bytes.NewBuffer(nil)}
	outAssertion = New(out)
	outAssertion.NotContains(nil, "key")
	expectedFail = "\"%!s(<nil>)\" could not be applied builtin len()"
	actualFail = out.buf.String()
	if !strings.Contains(actualFail, expectedFail) {
		t.Errorf("Contains failure should include %q but was %q", expectedFail, actualFail)
	}
}

func TestSubsetNotSubset(t *testing.T) {

	// MTestCase adds a custom message to the case
	cases := []struct {
		expected any
		actual   any
		result   bool
		message  string
	}{
		// cases that are expected to contain
		{[]int{1, 2, 3}, nil, true, "given subset is nil"},
		{[]int{1, 2, 3}, []int{}, true, "any set contains the nil set"},
		{[]int{1, 2, 3}, []int{1, 2}, true, "[1, 2, 3] contains [1, 2]"},
		{[]int{1, 2, 3}, []int{1, 2, 3}, true, "[1, 2, 3] contains [1, 2, 3"},
		{[]string{"hello", "world"}, []string{"hello"}, true, "[\"hello\", \"world\"] contains [\"hello\"]"},

		// cases that are expected not to contain
		{[]string{"hello", "world"}, []string{"hello", "testify"}, false, "[\"hello\", \"world\"] does not contain [\"hello\", \"testify\"]"},
		{[]int{1, 2, 3}, []int{4, 5}, false, "[1, 2, 3] does not contain [4, 5"},
		{[]int{1, 2, 3}, []int{1, 5}, false, "[1, 2, 3] does not contain [1, 5]"},
	}

	for _, c := range cases {
		t.Run("SubSet: "+c.message, func(t *testing.T) {

			mockT := new(testing.T)
			res := Subset(mockT, c.expected, c.actual)

			if res != c.result {
				if res {
					t.Errorf("Subset should return true: %s", c.message)
				} else {
					t.Errorf("Subset should return false: %s", c.message)
				}
			}
		})
	}
	for _, c := range cases {
		t.Run("NotSubSet: "+c.message, func(t *testing.T) {
			mockT := new(testing.T)
			res := NotSubset(mockT, c.expected, c.actual)

			// NotSubset should match the inverse of Subset. If it doesn't, something is wrong
			if res == Subset(mockT, c.expected, c.actual) {
				if res {
					t.Errorf("NotSubset should return true: %s", c.message)
				} else {
					t.Errorf("NotSubset should return false: %s", c.message)
				}
			}
		})
	}
}

func TestNotSubsetNil(t *testing.T) {
	mockT := new(testing.T)
	NotSubset(mockT, []string{"foo"}, nil)
	if !mockT.Failed() {
		t.Error("NotSubset on nil set should have failed the test")
	}
}

func Test_containsElement(t *testing.T) {
	assertion := New(t)
	list1 := []string{"Foo", "Bar"}
	list2 := []int{1, 2}
	simpleMap := map[any]any{"Foo": "Bar"}

	ok, found := containsElement("Hello World", "World")
	assertion.True(ok)
	assertion.True(found)

	ok, found = containsElement(list1, "Foo")
	assertion.True(ok)
	assertion.True(found)

	ok, found = containsElement(list1, "Bar")
	assertion.True(ok)
	assertion.True(found)

	ok, found = containsElement(list2, 1)
	assertion.True(ok)
	assertion.True(found)

	ok, found = containsElement(list2, 2)
	assertion.True(ok)
	assertion.True(found)

	ok, found = containsElement(list1, "Foo!")
	assertion.True(ok)
	assertion.False(found)

	ok, found = containsElement(list2, 3)
	assertion.True(ok)
	assertion.False(found)

	ok, found = containsElement(list2, "1")
	assertion.True(ok)
	assertion.False(found)

	ok, found = containsElement(simpleMap, "Foo")
	assertion.True(ok)
	assertion.True(found)

	ok, found = containsElement(simpleMap, "Bar")
	assertion.True(ok)
	assertion.False(found)

	ok, found = containsElement(1433, "1")
	assertion.False(ok)
	assertion.False(found)
}

func TestElementsMatch(t *testing.T) {
	mockT := new(testing.T)

	cases := []struct {
		expected any
		actual   any
		result   bool
	}{
		// matching
		{nil, nil, true},

		{nil, nil, true},
		{[]int{}, []int{}, true},
		{[]int{1}, []int{1}, true},
		{[]int{1, 1}, []int{1, 1}, true},
		{[]int{1, 2}, []int{1, 2}, true},
		{[]int{1, 2}, []int{2, 1}, true},
		{[2]int{1, 2}, [2]int{2, 1}, true},
		{[]string{"hello", "world"}, []string{"world", "hello"}, true},
		{[]string{"hello", "hello"}, []string{"hello", "hello"}, true},
		{[]string{"hello", "hello", "world"}, []string{"hello", "world", "hello"}, true},
		{[3]string{"hello", "hello", "world"}, [3]string{"hello", "world", "hello"}, true},
		{[]int{}, nil, true},

		// not matching
		{[]int{1}, []int{1, 1}, false},
		{[]int{1, 2}, []int{2, 2}, false},
		{[]string{"hello", "hello"}, []string{"hello"}, false},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("ElementsMatch(%#v, %#v)", c.expected, c.actual), func(t *testing.T) {
			res := ElementsMatch(mockT, c.actual, c.expected)

			if res != c.result {
				t.Errorf("ElementsMatch(%#v, %#v) should return %v", c.actual, c.expected, c.result)
			}
		})
	}
}

func TestDiffLists(t *testing.T) {
	tests := []struct {
		name   string
		listA  any
		listB  any
		extraA []any
		extraB []any
	}{
		{
			name:   "equal empty",
			listA:  []string{},
			listB:  []string{},
			extraA: nil,
			extraB: nil,
		},
		{
			name:   "equal same order",
			listA:  []string{"hello", "world"},
			listB:  []string{"hello", "world"},
			extraA: nil,
			extraB: nil,
		},
		{
			name:   "equal different order",
			listA:  []string{"hello", "world"},
			listB:  []string{"world", "hello"},
			extraA: nil,
			extraB: nil,
		},
		{
			name:   "extra A",
			listA:  []string{"hello", "hello", "world"},
			listB:  []string{"hello", "world"},
			extraA: []any{"hello"},
			extraB: nil,
		},
		{
			name:   "extra A twice",
			listA:  []string{"hello", "hello", "hello", "world"},
			listB:  []string{"hello", "world"},
			extraA: []any{"hello", "hello"},
			extraB: nil,
		},
		{
			name:   "extra B",
			listA:  []string{"hello", "world"},
			listB:  []string{"hello", "hello", "world"},
			extraA: nil,
			extraB: []any{"hello"},
		},
		{
			name:   "extra B twice",
			listA:  []string{"hello", "world"},
			listB:  []string{"hello", "hello", "world", "hello"},
			extraA: nil,
			extraB: []any{"hello", "hello"},
		},
		{
			name:   "integers 1",
			listA:  []int{1, 2, 3, 4, 5},
			listB:  []int{5, 4, 3, 2, 1},
			extraA: nil,
			extraB: nil,
		},
		{
			name:   "integers 2",
			listA:  []int{1, 2, 1, 2, 1},
			listB:  []int{2, 1, 2, 1, 2},
			extraA: []any{1},
			extraB: []any{2},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			actualExtraA, actualExtraB := diffLists(test.listA, test.listB)
			New(t).Equal(test.extraA, actualExtraA, "extra A does not match for listA=%v listB=%v",
				test.listA, test.listB)
			New(t).Equal(test.extraB, actualExtraB, "extra B does not match for listA=%v listB=%v",
				test.listA, test.listB)
		})
	}
}

func TestCondition(t *testing.T) {
	mockT := new(testing.T)

	if !Condition(mockT, func() bool { return true }, "Truth") {
		t.Error("Condition should return true")
	}

	if Condition(mockT, func() bool { return false }, "Lie") {
		t.Error("Condition should return false")
	}

}

func TestDidPanic(t *testing.T) {

	if funcDidPanic, _, _ := didPanic(func() {
		panic("Panic!")
	}); !funcDidPanic {
		t.Error("didPanic should return true")
	}

	if funcDidPanic, _, _ := didPanic(func() {
	}); funcDidPanic {
		t.Error("didPanic should return false")
	}

}

func TestPanics(t *testing.T) {

	mockT := new(testing.T)

	if !Panics(mockT, func() {
		panic("Panic!")
	}) {
		t.Error("Panics should return true")
	}

	if Panics(mockT, func() {
	}) {
		t.Error("Panics should return false")
	}

}

func TestPanicsWithValue(t *testing.T) {

	mockT := new(testing.T)

	if !PanicsWithValue(mockT, "Panic!", func() {
		panic("Panic!")
	}) {
		t.Error("PanicsWithValue should return true")
	}

	if PanicsWithValue(mockT, "Panic!", func() {
	}) {
		t.Error("PanicsWithValue should return false")
	}

	if PanicsWithValue(mockT, "at the disco", func() {
		panic("Panic!")
	}) {
		t.Error("PanicsWithValue should return false")
	}
}

func TestPanicsWithError(t *testing.T) {

	mockT := new(testing.T)

	if !PanicsWithError(mockT, "panic", func() {
		panic(errors.New("panic"))
	}) {
		t.Error("PanicsWithError should return true")
	}

	if PanicsWithError(mockT, "Panic!", func() {
	}) {
		t.Error("PanicsWithError should return false")
	}

	if PanicsWithError(mockT, "at the disco", func() {
		panic(errors.New("panic"))
	}) {
		t.Error("PanicsWithError should return false")
	}

	if PanicsWithError(mockT, "Panic!", func() {
		panic("panic")
	}) {
		t.Error("PanicsWithError should return false")
	}
}

func TestNotPanics(t *testing.T) {

	mockT := new(testing.T)

	if !NotPanics(mockT, func() {
	}) {
		t.Error("NotPanics should return true")
	}

	if NotPanics(mockT, func() {
		panic("Panic!")
	}) {
		t.Error("NotPanics should return false")
	}

}

func TestNoError(t *testing.T) {
	mockAssertion := NewWithOnFailureNoop(new(testing.T))

	// start with a nil error
	var err error

	New(t).True(mockAssertion.NoError(err), "NoError should return True for nil arg")

	// now set an error
	err = errors.New("some error")

	New(t).False(mockAssertion.NoError(err), "NoError with error should return False")

	// returning an empty error interface
	err = func() error {
		var err *customError
		return err
	}()

	if err == nil { // err is not nil here!
		t.Errorf("Error should be nil due to empty interface: %s", err)
	}

	New(t).False(mockAssertion.NoError(err), "NoError should fail with empty error interface")
}

type customError struct{}

func (*customError) Error() string { return "fail" }

func TestError(t *testing.T) {
	mockAssertion := NewWithOnFailureNoop(new(testing.T))

	// start with a nil error
	var err error

	New(t).False(mockAssertion.Error(err), "Error should return False for nil arg")

	// now set an error
	err = errors.New("some error")

	New(t).True(mockAssertion.Error(err), "Error with error should return True")

	// returning an empty error interface
	err = func() error {
		var err *customError
		return err
	}()

	if err == nil { // err is not nil here!
		t.Errorf("Error should be nil due to empty interface: %s", err)
	}

	New(t).True(mockAssertion.Error(err), "Error should pass with empty error interface")
}

func TestEqualError(t *testing.T) {
	mockAssertion := NewWithOnFailureNoop(new(testing.T))

	// start with a nil error
	var err error
	New(t).False(mockAssertion.EqualError(err, ""), "EqualError should return false for nil arg")

	// now set an error
	err = errors.New("some error")
	New(t).False(mockAssertion.EqualError(err, "Not some error"), "EqualError should return false for different error string")
	New(t).True(mockAssertion.EqualError(err, "some error"), "EqualError should return true")
}

func TestErrorContains(t *testing.T) {
	mockAssertion := NewWithOnFailureNoop(new(testing.T))

	// start with a nil error
	var err error
	New(t).False(mockAssertion.ErrorContains(err, ""), "ErrorContains should return false for nil arg")

	// now set an error
	err = errors.New("some error: another error")
	New(t).False(mockAssertion.ErrorContains(err, "bad error"), "ErrorContains should return false for different error string")
	New(t).True(mockAssertion.ErrorContains(err, "some error"), "ErrorContains should return true")
	New(t).True(mockAssertion.ErrorContains(err, "another error"), "ErrorContains should return true")
}

func TestErrorRegexp(t *testing.T) {
	mockAssertion := NewWithOnFailureNoop(new(testing.T))

	// start with a nil error
	var err error
	New(t).False(mockAssertion.ErrorRegexp(err, ""), "ErrorContains should return false for nil arg")

	// now set an error
	err = errors.New("some error: another error")
	New(t).False(mockAssertion.ErrorRegexp(err, "bad error"), "ErrorContains should return false for different error string")
	New(t).True(mockAssertion.ErrorRegexp(err, "some error"), "ErrorContains should return true")
	New(t).True(mockAssertion.ErrorRegexp(err, "another error"), "ErrorContains should return true")
}

func TestIsEmpty(t *testing.T) {
	assertion := New(t)
	chWithValue := make(chan struct{}, 1)
	chWithValue <- struct{}{}

	assertion.True(isEmpty(""))
	assertion.True(isEmpty(nil))
	assertion.True(isEmpty([]string{}))
	assertion.True(isEmpty(0))
	assertion.True(isEmpty(int32(0)))
	assertion.True(isEmpty(int64(0)))
	assertion.True(isEmpty(false))
	assertion.True(isEmpty(map[string]string{}))
	assertion.True(isEmpty(new(time.Time)))
	assertion.True(isEmpty(time.Time{}))
	assertion.True(isEmpty(make(chan struct{})))
	assertion.False(isEmpty("something"))
	assertion.False(isEmpty(errors.New("something")))
	assertion.False(isEmpty([]string{"something"}))
	assertion.False(isEmpty(1))
	assertion.False(isEmpty(true))
	assertion.False(isEmpty(map[string]string{"Hello": "World"}))
	assertion.False(isEmpty(chWithValue))

}

func TestEmpty(t *testing.T) {
	mockT := new(testing.T)
	assertion := New(t)

	chWithValue := make(chan struct{}, 1)
	chWithValue <- struct{}{}
	var tiP *time.Time
	var tiNP time.Time
	var s *string
	var f *os.File
	sP := &s
	x := 1
	xP := &x

	type TString string
	type TStruct struct {
		x int
	}

	assertion.True(Empty(mockT, ""), "Empty string is empty")
	assertion.True(Empty(mockT, nil), "Nil is empty")
	assertion.True(Empty(mockT, []string{}), "Empty string array is empty")
	assertion.True(Empty(mockT, 0), "Zero int value is empty")
	assertion.True(Empty(mockT, false), "False value is empty")
	assertion.True(Empty(mockT, make(chan struct{})), "Channel without values is empty")
	assertion.True(Empty(mockT, s), "Nil string pointer is empty")
	assertion.True(Empty(mockT, f), "Nil os.File pointer is empty")
	assertion.True(Empty(mockT, tiP), "Nil time.Time pointer is empty")
	assertion.True(Empty(mockT, tiNP), "time.Time is empty")
	assertion.True(Empty(mockT, TStruct{}), "struct with zero values is empty")
	assertion.True(Empty(mockT, TString("")), "empty aliased string is empty")
	assertion.True(Empty(mockT, sP), "ptr to nil value is empty")

	assertion.False(Empty(mockT, "something"), "Non Empty string is not empty")
	assertion.False(Empty(mockT, errors.New("something")), "Non nil object is not empty")
	assertion.False(Empty(mockT, []string{"something"}), "Non empty string array is not empty")
	assertion.False(Empty(mockT, 1), "Non-zero int value is not empty")
	assertion.False(Empty(mockT, true), "True value is not empty")
	assertion.False(Empty(mockT, chWithValue), "Channel with values is not empty")
	assertion.False(Empty(mockT, TStruct{x: 1}), "struct with initialized values is empty")
	assertion.False(Empty(mockT, TString("abc")), "non-empty aliased string is empty")
	assertion.False(Empty(mockT, xP), "ptr to non-nil value is not empty")
}

func TestNotEmpty(t *testing.T) {
	mockT := new(testing.T)
	assertion := New(t)

	chWithValue := make(chan struct{}, 1)
	chWithValue <- struct{}{}

	assertion.False(NotEmpty(mockT, ""), "Empty string is empty")
	assertion.False(NotEmpty(mockT, nil), "Nil is empty")
	assertion.False(NotEmpty(mockT, []string{}), "Empty string array is empty")
	assertion.False(NotEmpty(mockT, 0), "Zero int value is empty")
	assertion.False(NotEmpty(mockT, false), "False value is empty")
	assertion.False(NotEmpty(mockT, make(chan struct{})), "Channel without values is empty")

	assertion.True(NotEmpty(mockT, "something"), "Non Empty string is not empty")
	assertion.True(NotEmpty(mockT, errors.New("something")), "Non nil object is not empty")
	assertion.True(NotEmpty(mockT, []string{"something"}), "Non empty string array is not empty")
	assertion.True(NotEmpty(mockT, 1), "Non-zero int value is not empty")
	assertion.True(NotEmpty(mockT, true), "True value is not empty")
	assertion.True(NotEmpty(mockT, chWithValue), "Channel with values is not empty")
}

func Test_getLen(t *testing.T) {
	assertion := New(t)
	falseCases := []any{
		nil,
		0,
		true,
		false,
		'A',
		struct{}{},
	}
	for _, v := range falseCases {
		ok, l := getLen(v)
		assertion.False(ok, "Expected getLen fail to get length of %#v", v)
		New(t).Equal(0, l, "getLen should return 0 for %#v", v)
	}

	ch := make(chan int, 5)
	ch <- 1
	ch <- 2
	ch <- 3
	trueCases := []struct {
		v any
		l int
	}{
		{[]int{1, 2, 3}, 3},
		{[...]int{1, 2, 3}, 3},
		{"ABC", 3},
		{map[int]int{1: 2, 2: 4, 3: 6}, 3},
		{ch, 3},

		{[]int{}, 0},
		{map[int]int{}, 0},
		{make(chan int), 0},

		{[]int(nil), 0},
		{map[int]int(nil), 0},
		{(chan int)(nil), 0},
	}

	for _, c := range trueCases {
		ok, l := getLen(c.v)
		assertion.True(ok, "Expected getLen success to get length of %#v", c.v)
		New(t).Equal(c.l, l)
	}
}

func TestLen(t *testing.T) {
	mockT := new(testing.T)
	assertion := New(t)

	assertion.False(Len(mockT, nil, 0), "nil does not have length")
	assertion.False(Len(mockT, 0, 0), "int does not have length")
	assertion.False(Len(mockT, true, 0), "true does not have length")
	assertion.False(Len(mockT, false, 0), "false does not have length")
	assertion.False(Len(mockT, 'A', 0), "Rune does not have length")
	assertion.False(Len(mockT, struct{}{}, 0), "Struct does not have length")

	ch := make(chan int, 5)
	ch <- 1
	ch <- 2
	ch <- 3

	cases := []struct {
		v any
		l int
	}{
		{[]int{1, 2, 3}, 3},
		{[...]int{1, 2, 3}, 3},
		{"ABC", 3},
		{map[int]int{1: 2, 2: 4, 3: 6}, 3},
		{ch, 3},

		{[]int{}, 0},
		{map[int]int{}, 0},
		{make(chan int), 0},

		{[]int(nil), 0},
		{map[int]int(nil), 0},
		{(chan int)(nil), 0},
	}

	for _, c := range cases {
		assertion.True(Len(mockT, c.v, c.l), "%#v have %d items", c.v, c.l)
	}

	cases = []struct {
		v any
		l int
	}{
		{[]int{1, 2, 3}, 4},
		{[...]int{1, 2, 3}, 2},
		{"ABC", 2},
		{map[int]int{1: 2, 2: 4, 3: 6}, 4},
		{ch, 2},

		{[]int{}, 1},
		{map[int]int{}, 1},
		{make(chan int), 1},

		{[]int(nil), 1},
		{map[int]int(nil), 1},
		{(chan int)(nil), 1},
	}

	for _, c := range cases {
		assertion.False(Len(mockT, c.v, c.l), "%#v have %d items", c.v, c.l)
	}
}

func TestWithinDuration(t *testing.T) {
	mockT := new(testing.T)
	assertion := New(t)
	a := time.Now()
	b := a.Add(10 * time.Second)

	assertion.True(WithinDuration(mockT, a, b, 10*time.Second), "A 10s difference is within a 10s time difference")
	assertion.True(WithinDuration(mockT, b, a, 10*time.Second), "A 10s difference is within a 10s time difference")

	assertion.False(WithinDuration(mockT, a, b, 9*time.Second), "A 10s difference is not within a 9s time difference")
	assertion.False(WithinDuration(mockT, b, a, 9*time.Second), "A 10s difference is not within a 9s time difference")

	assertion.False(WithinDuration(mockT, a, b, -9*time.Second), "A 10s difference is not within a 9s time difference")
	assertion.False(WithinDuration(mockT, b, a, -9*time.Second), "A 10s difference is not within a 9s time difference")

	assertion.False(WithinDuration(mockT, a, b, -11*time.Second), "A 10s difference is not within a 9s time difference")
	assertion.False(WithinDuration(mockT, b, a, -11*time.Second), "A 10s difference is not within a 9s time difference")
}

func TestInDelta(t *testing.T) {
	mockT := new(testing.T)
	assertion := New(t)

	assertion.True(InDelta(mockT, 1.001, 1, 0.01), "|1.001 - 1| <= 0.01")
	assertion.True(InDelta(mockT, 1, 1.001, 0.01), "|1 - 1.001| <= 0.01")
	assertion.True(InDelta(mockT, 1, 2, 1), "|1 - 2| <= 1")
	assertion.False(InDelta(mockT, 1, 2, 0.5), "Expected |1 - 2| <= 0.5 to fail")
	assertion.False(InDelta(mockT, 2, 1, 0.5), "Expected |2 - 1| <= 0.5 to fail")
	assertion.False(InDelta(mockT, "", nil, 1), "Expected non numerals to fail")
	assertion.False(InDelta(mockT, 42, math.NaN(), 0.01), "Expected NaN for actual to fail")
	assertion.False(InDelta(mockT, math.NaN(), 42, 0.01), "Expected NaN for expected to fail")
	assertion.True(InDelta(mockT, math.NaN(), math.NaN(), 0.01), "Expected NaN for both to pass")

	cases := []struct {
		a, b  any
		delta float64
	}{
		{uint(2), uint(1), 1},
		{uint8(2), uint8(1), 1},
		{uint16(2), uint16(1), 1},
		{uint32(2), uint32(1), 1},
		{uint64(2), uint64(1), 1},

		{2, 1, 1},
		{int8(2), int8(1), 1},
		{int16(2), int16(1), 1},
		{int32(2), int32(1), 1},
		{int64(2), int64(1), 1},

		{float32(2), float32(1), 1},
		{float64(2), float64(1), 1},
	}

	for _, tc := range cases {
		assertion.True(InDelta(mockT, tc.a, tc.b, tc.delta), "Expected |%V - %V| <= %v", tc.a, tc.b, tc.delta)
	}
}

func TestInDeltaSlice(t *testing.T) {
	mockT := new(testing.T)
	assertion := New(t)

	assertion.True(InDeltaSlice(mockT,
		[]float64{1.001, math.NaN(), 0.999},
		[]float64{1, math.NaN(), 1},
		0.1), "{1.001, NaN, 0.009} is element-wise close to {1, NaN, 1} in delta=0.1")

	assertion.True(InDeltaSlice(mockT,
		[]float64{1, math.NaN(), 2},
		[]float64{0, math.NaN(), 3},
		1), "{1, NaN, 2} is element-wise close to {0, NaN, 3} in delta=1")

	assertion.False(InDeltaSlice(mockT,
		[]float64{1, math.NaN(), 2},
		[]float64{0, math.NaN(), 3},
		0.1), "{1, NaN, 2} is not element-wise close to {0, NaN, 3} in delta=0.1")

	assertion.False(InDeltaSlice(mockT, "", nil, 1), "Expected non numeral slices to fail")
}

func TestInDeltaMapValues(t *testing.T) {
	mockT := new(testing.T)
	assertion := New(t)

	for _, tc := range []struct {
		title  string
		expect any
		actual any
		f      func(bool, ...any) bool
		delta  float64
	}{
		{
			title: "Within delta",
			expect: map[string]float64{
				"foo": 1.0,
				"bar": 2.0,
				"baz": math.NaN(),
			},
			actual: map[string]float64{
				"foo": 1.01,
				"bar": 1.99,
				"baz": math.NaN(),
			},
			delta: 0.1,
			f:     assertion.True,
		},
		{
			title: "Within delta",
			expect: map[int]float64{
				1: 1.0,
				2: 2.0,
			},
			actual: map[int]float64{
				1: 1.0,
				2: 1.99,
			},
			delta: 0.1,
			f:     assertion.True,
		},
		{
			title: "Different number of keys",
			expect: map[int]float64{
				1: 1.0,
				2: 2.0,
			},
			actual: map[int]float64{
				1: 1.0,
			},
			delta: 0.1,
			f:     assertion.False,
		},
		{
			title: "Within delta with zero value",
			expect: map[string]float64{
				"zero": 0,
			},
			actual: map[string]float64{
				"zero": 0,
			},
			delta: 0.1,
			f:     assertion.True,
		},
		{
			title: "With missing key with zero value",
			expect: map[string]float64{
				"zero": 0,
				"foo":  0,
			},
			actual: map[string]float64{
				"zero": 0,
				"bar":  0,
			},
			f: assertion.False,
		},
	} {
		tc.f(InDeltaMapValues(mockT, tc.expect, tc.actual, tc.delta), tc.title+"\n"+diff(tc.expect, tc.actual))
	}
}

func TestInEpsilon(t *testing.T) {
	mockT := new(testing.T)
	assertion := New(t)

	cases := []struct {
		a, b    any
		epsilon float64
	}{
		{uint8(2), uint16(2), .001},
		{2.1, 2.2, 0.1},
		{2.2, 2.1, 0.1},
		{-2.1, -2.2, 0.1},
		{-2.2, -2.1, 0.1},
		{uint64(100), uint8(101), 0.01},
		{0.1, -0.1, 2},
		{0.1, 0, 2},
		{math.NaN(), math.NaN(), 1},
		{time.Second, time.Second + time.Millisecond, 0.002},
	}

	for _, tc := range cases {
		assertion.True(InEpsilon(t, tc.a, tc.b, tc.epsilon, "Expected %V and %V to have a relative difference of %v", tc.a, tc.b, tc.epsilon), "test: %q", tc)
	}

	cases = []struct {
		a, b    any
		epsilon float64
	}{
		{uint8(2), int16(-2), .001},
		{uint64(100), uint8(102), 0.01},
		{2.1, 2.2, 0.001},
		{2.2, 2.1, 0.001},
		{2.1, -2.2, 1},
		{2.1, "bla-bla", 0},
		{0.1, -0.1, 1.99},
		{0, 0.1, 2}, // expected must be different to zero
		{time.Second, time.Second + 10*time.Millisecond, 0.002},
		{math.NaN(), 0, 1},
		{0, math.NaN(), 1},
		{0, 0, math.NaN()},
	}

	for _, tc := range cases {
		assertion.False(InEpsilon(mockT, tc.a, tc.b, tc.epsilon, "Expected %V and %V to have a relative difference of %v", tc.a, tc.b, tc.epsilon))
	}

}

func TestInEpsilonSlice(t *testing.T) {
	mockT := new(testing.T)
	assertion := New(t)

	assertion.True(InEpsilonSlice(mockT,
		[]float64{2.2, math.NaN(), 2.0},
		[]float64{2.1, math.NaN(), 2.1},
		0.06), "{2.2, NaN, 2.0} is element-wise close to {2.1, NaN, 2.1} in espilon=0.06")

	assertion.False(InEpsilonSlice(mockT,
		[]float64{2.2, 2.0},
		[]float64{2.1, 2.1},
		0.04), "{2.2, 2.0} is not element-wise close to {2.1, 2.1} in espilon=0.04")

	assertion.False(InEpsilonSlice(mockT, "", nil, 1), "Expected non numeral slices to fail")
}

func TestRegexp(t *testing.T) {
	mockT := new(testing.T)
	assertion := New(t)

	cases := []struct {
		rx, str string
	}{
		{"^start", "start of the line"},
		{"end$", "in the end"},
		{"[0-9]{3}[.-]?[0-9]{2}[.-]?[0-9]{2}", "My phone number is 650.12.34"},
	}

	for _, tc := range cases {
		assertion.True(Regexp(mockT, tc.rx, tc.str))
		assertion.True(Regexp(mockT, regexp.MustCompile(tc.rx), tc.str))
		assertion.False(NotRegexp(mockT, tc.rx, tc.str))
		assertion.False(NotRegexp(mockT, regexp.MustCompile(tc.rx), tc.str))
	}

	cases = []struct {
		rx, str string
	}{
		{"^asdfastart", "Not the start of the line"},
		{"end$", "in the end."},
		{"[0-9]{3}[.-]?[0-9]{2}[.-]?[0-9]{2}", "My phone number is 650.12a.34"},
	}

	for _, tc := range cases {
		assertion.False(Regexp(mockT, tc.rx, tc.str), "Expected \"%s\" to not match \"%s\"", tc.rx, tc.str)
		assertion.False(Regexp(mockT, regexp.MustCompile(tc.rx), tc.str))
		assertion.True(NotRegexp(mockT, tc.rx, tc.str))
		assertion.True(NotRegexp(mockT, regexp.MustCompile(tc.rx), tc.str))
	}
}

func testAutogeneratedFunction() {
	defer func() {
		if err := recover(); err == nil {
			panic("did not panic")
		}
		CallerInfo()
	}()
	t := struct {
		io.Closer
	}{}
	c := t
	_ = c.Close()
}

func TestCallerInfoWithAutogeneratedFunctions(t *testing.T) {
	NotPanics(t, func() {
		testAutogeneratedFunction()
	})
}

func TestZero(t *testing.T) {
	mockT := new(testing.T)
	assertion := New(t)

	for _, test := range zeros {
		assertion.True(Zero(mockT, test, "%#v is not the %v zero value", test, reflect.TypeOf(test)))
	}

	for _, test := range nonZeros {
		assertion.False(Zero(mockT, test, "%#v is not the %v zero value", test, reflect.TypeOf(test)))
	}
}

func TestNotZero(t *testing.T) {
	mockT := new(testing.T)
	assertion := New(t)

	for _, test := range zeros {
		assertion.False(NotZero(mockT, test, "%#v is not the %v zero value", test, reflect.TypeOf(test)))
	}

	for _, test := range nonZeros {
		assertion.True(NotZero(mockT, test, "%#v is not the %v zero value", test, reflect.TypeOf(test)))
	}
}

func TestFileExists(t *testing.T) {
	assertion := New(t)

	mockT := new(testing.T)
	assertion.True(FileExists(mockT, "assertions.go"))

	mockT = new(testing.T)
	assertion.False(FileExists(mockT, "random_file"))

	mockT = new(testing.T)
	assertion.False(FileExists(mockT, "../_codegen"))

	var tempFiles []string

	link, err := getTempSymlinkPath("assertions.go")
	if err != nil {
		t.Fatal("could not create temp symlink, err:", err)
	}
	tempFiles = append(tempFiles, link)
	mockT = new(testing.T)
	assertion.True(FileExists(mockT, link))

	link, err = getTempSymlinkPath("non_existent_file")
	if err != nil {
		t.Fatal("could not create temp symlink, err:", err)
	}
	tempFiles = append(tempFiles, link)
	mockT = new(testing.T)
	assertion.True(FileExists(mockT, link))

	errs := cleanUpTempFiles(tempFiles)
	if len(errs) > 0 {
		t.Fatal("could not clean up temporary files")
	}
}

func TestNoFileExists(t *testing.T) {
	assertion := New(t)

	mockT := new(testing.T)
	assertion.False(NoFileExists(mockT, "assertions.go"))

	mockT = new(testing.T)
	assertion.True(NoFileExists(mockT, "non_existent_file"))

	mockT = new(testing.T)
	assertion.True(NoFileExists(mockT, "../_codegen"))

	var tempFiles []string

	link, err := getTempSymlinkPath("assertions.go")
	if err != nil {
		t.Fatal("could not create temp symlink, err:", err)
	}
	tempFiles = append(tempFiles, link)
	mockT = new(testing.T)
	assertion.False(NoFileExists(mockT, link))

	link, err = getTempSymlinkPath("non_existent_file")
	if err != nil {
		t.Fatal("could not create temp symlink, err:", err)
	}
	tempFiles = append(tempFiles, link)
	mockT = new(testing.T)
	assertion.False(NoFileExists(mockT, link))

	errs := cleanUpTempFiles(tempFiles)
	if len(errs) > 0 {
		t.Fatal("could not clean up temporary files")
	}
}

func getTempSymlinkPath(file string) (string, error) {
	link := file + "_symlink"
	err := os.Symlink(file, link)
	return link, err
}

func cleanUpTempFiles(paths []string) []error {
	var res []error
	for _, path := range paths {
		err := os.Remove(path)
		if err != nil {
			res = append(res, err)
		}
	}
	return res
}

func TestDirExists(t *testing.T) {
	assertion := New(t)

	mockT := new(testing.T)
	assertion.False(DirExists(mockT, "assertions.go"))

	mockT = new(testing.T)
	assertion.False(DirExists(mockT, "non_existent_dir"))

	mockT = new(testing.T)
	assertion.True(DirExists(mockT, "."))

	var tempFiles []string

	link, err := getTempSymlinkPath("assertions.go")
	if err != nil {
		t.Fatal("could not create temp symlink, err:", err)
	}
	tempFiles = append(tempFiles, link)
	mockT = new(testing.T)
	assertion.False(DirExists(mockT, link))

	link, err = getTempSymlinkPath("non_existent_dir")
	if err != nil {
		t.Fatal("could not create temp symlink, err:", err)
	}
	tempFiles = append(tempFiles, link)
	mockT = new(testing.T)
	assertion.False(DirExists(mockT, link))

	errs := cleanUpTempFiles(tempFiles)
	if len(errs) > 0 {
		t.Fatal("could not clean up temporary files")
	}
}

func TestNoDirExists(t *testing.T) {
	assertion := New(t)

	mockT := new(testing.T)
	assertion.True(NoDirExists(mockT, "assertions.go"))

	mockT = new(testing.T)
	assertion.True(NoDirExists(mockT, "non_existent_dir"))

	mockT = new(testing.T)
	assertion.False(NoDirExists(mockT, "."))

	var tempFiles []string

	link, err := getTempSymlinkPath("assertions.go")
	if err != nil {
		t.Fatal("could not create temp symlink, err:", err)
	}
	tempFiles = append(tempFiles, link)
	mockT = new(testing.T)
	assertion.True(NoDirExists(mockT, link))

	link, err = getTempSymlinkPath("non_existent_dir")
	if err != nil {
		t.Fatal("could not create temp symlink, err:", err)
	}
	tempFiles = append(tempFiles, link)
	mockT = new(testing.T)
	assertion.True(NoDirExists(mockT, link))

	errs := cleanUpTempFiles(tempFiles)
	if len(errs) > 0 {
		t.Fatal("could not clean up temporary files")
	}
}

func TestJSONEq(t *testing.T) {
	mockAssertion := NewWithOnFailureNoop(new(testing.T))
	for _, test := range []struct {
		name     string
		expected string
		actual   string
		result   bool
	}{
		{"EqualSONString", `{"hello": "world", "foo": "bar"}`, `{"hello": "world", "foo": "bar"}`, true},
		{"EquivalentButNotEqual", `{"hello": "world", "foo": "bar"}`, `{"foo": "bar", "hello": "world"}`, true},
		{"HashOfArraysAndHashes", "{\r\n\t\"numeric\": 1.5,\r\n\t\"array\": [{\"foo\": \"bar\"}, 1, \"string\", [\"nested\", \"array\", 5.5]],\r\n\t\"hash\": {\"nested\": \"hash\", \"nested_slice\": [\"this\", \"is\", \"nested\"]},\r\n\t\"string\": \"foo\"\r\n}", "{\r\n\t\"numeric\": 1.5,\r\n\t\"hash\": {\"nested\": \"hash\", \"nested_slice\": [\"this\", \"is\", \"nested\"]},\r\n\t\"string\": \"foo\",\r\n\t\"array\": [{\"foo\": \"bar\"}, 1, \"string\", [\"nested\", \"array\", 5.5]]\r\n}", true},
		{"Array", `["foo", {"hello": "world", "nested": "hash"}]`, `["foo", {"nested": "hash", "hello": "world"}]`, true},
		{"HashAndArrayNotEquivalent", `["foo", {"hello": "world", "nested": "hash"}]`, `{"foo": "bar", {"nested": "hash", "hello": "world"}}`, false},
		{"HashesNotEquivalent", `{"foo": "bar"}`, `{"foo": "bar", "hello": "world"}`, false},
		{"ActualIsNotJSON", `{"foo": "bar"}`, "Not JSON", false},
		{"ExpectedIsNotJSON", "Not JSON", `{"foo": "bar", "hello": "world"}`, false},
		{"ExpectedAndActualNotJSON", "Not JSON", "Not JSON", false},
		{"ArraysOfDifferentOrder", `["foo", {"hello": "world", "nested": "hash"}]`, `[{ "hello": "world", "nested": "hash"}, "foo"]`, false},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			New(t).Equal(test.result, mockAssertion.JSONEq(test.expected, test.actual))
		})
	}
}

func TestYAMLEq(t *testing.T) {
	mockAssertion := NewWithOnFailureNoop(new(testing.T))

	hashOfArraysAndHashesExpected := `
numeric: 1.5
array:
  - foo: bar
  - 1
  - "string"
  - ["nested", "array", 5.5]
hash:
  nested: hash
  nested_slice: [this, is, nested]
string: "foo"
`
	hashOfArraysAndHashesActual := `
numeric: 1.5
hash:
  nested: hash
  nested_slice: [this, is, nested]
string: "foo"
array:
  - foo: bar
  - 1
  - "string"
  - ["nested", "array", 5.5]
`

	for _, test := range []struct {
		name     string
		expected string
		actual   string
		result   bool
	}{
		{"EqualYAMLString", `{"hello": "world", "foo": "bar"}`, `{"hello": "world", "foo": "bar"}`, true},
		{"EquivalentButNotEqual", `{"hello": "world", "foo": "bar"}`, `{"foo": "bar", "hello": "world"}`, true},
		{"HashOfArraysAndHashes", hashOfArraysAndHashesExpected, hashOfArraysAndHashesActual, true},
		{"Array", `["foo", {"hello": "world", "nested": "hash"}]`, `["foo", {"nested": "hash", "hello": "world"}]`, true},
		{"HashAndArrayNotEquivalent", `["foo", {"hello": "world", "nested": "hash"}]`, `{"foo": "bar", {"nested": "hash", "hello": "world"}}`, false},
		{"HashesNotEquivalent", `{"foo": "bar"}`, `{"foo": "bar", "hello": "world"}`, false},
		{"ActualIsSimpleString", `{"foo": "bar"}`, "Simple String", false},
		{"ExpectedIsSimpleString", "Simple String", `{"foo": "bar", "hello": "world"}`, false},
		{"ExpectedAndActualSimpleString", "Simple String", "Simple String", true},
		{"ArraysOfDifferentOrder", `["foo", {"hello": "world", "nested": "hash"}]`, `[{ "hello": "world", "nested": "hash"}, "foo"]`, false},
	} {
		test := test
		t.Run(test.name, func(t *testing.T) {
			New(t).Equal(test.result, mockAssertion.YAMLEq(test.expected, test.actual))
		})
	}
}

type diffTestingStruct struct {
	A string
	B int
}

func (d *diffTestingStruct) String() string {
	return d.A
}

func TestDiff(t *testing.T) {
	expected := `

Diff:
--- Expected
+++ Actual
@@ -1,3 +1,3 @@
 (struct { foo string }) {
- foo: (string) (len=5) "hello"
+ foo: (string) (len=3) "bar"
 }
`
	actual := diff(
		struct{ foo string }{"hello"},
		struct{ foo string }{"bar"},
	)
	New(t).Equal(expected, actual)

	expected = `

Diff:
--- Expected
+++ Actual
@@ -2,5 +2,5 @@
  (int) 1,
- (int) 2,
  (int) 3,
- (int) 4
+ (int) 5,
+ (int) 7
 }
`
	actual = diff(
		[]int{1, 2, 3, 4},
		[]int{1, 3, 5, 7},
	)
	New(t).Equal(expected, actual)

	expected = `

Diff:
--- Expected
+++ Actual
@@ -2,4 +2,4 @@
  (int) 1,
- (int) 2,
- (int) 3
+ (int) 3,
+ (int) 5
 }
`
	actual = diff(
		[]int{1, 2, 3, 4}[0:3],
		[]int{1, 3, 5, 7}[0:3],
	)
	New(t).Equal(expected, actual)

	expected = `

Diff:
--- Expected
+++ Actual
@@ -1,6 +1,6 @@
 (map[string]int) (len=4) {
- (string) (len=4) "four": (int) 4,
+ (string) (len=4) "five": (int) 5,
  (string) (len=3) "one": (int) 1,
- (string) (len=5) "three": (int) 3,
- (string) (len=3) "two": (int) 2
+ (string) (len=5) "seven": (int) 7,
+ (string) (len=5) "three": (int) 3
 }
`

	actual = diff(
		map[string]int{"one": 1, "two": 2, "three": 3, "four": 4},
		map[string]int{"one": 1, "three": 3, "five": 5, "seven": 7},
	)
	New(t).Equal(expected, actual)

	expected = `

Diff:
--- Expected
+++ Actual
@@ -1,3 +1,3 @@
 (*errors.errorString)({
- s: (string) (len=19) "some expected error"
+ s: (string) (len=12) "actual error"
 })
`

	actual = diff(
		errors.New("some expected error"),
		errors.New("actual error"),
	)
	New(t).Equal(expected, actual)

	expected = `

Diff:
--- Expected
+++ Actual
@@ -2,3 +2,3 @@
  A: (string) (len=11) "some string",
- B: (int) 10
+ B: (int) 15
 }
`

	actual = diff(
		diffTestingStruct{A: "some string", B: 10},
		diffTestingStruct{A: "some string", B: 15},
	)
	New(t).Equal(expected, actual)

	expected = `

Diff:
--- Expected
+++ Actual
@@ -1,2 +1,2 @@
-(time.Time) 2020-09-24 00:00:00 +0000 UTC
+(time.Time) 2020-09-25 00:00:00 +0000 UTC
 
`

	actual = diff(
		time.Date(2020, 9, 24, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 9, 25, 0, 0, 0, 0, time.UTC),
	)
	New(t).Equal(expected, actual)
}

func TestTimeEqualityErrorFormatting(t *testing.T) {
	out := &outputT{buf: bytes.NewBuffer(nil)}
	New(out).Equal(time.Second*2, time.Millisecond)
	expectedErr := "\\s+Error Trace:\\s+Error:\\s+Not equal:\\s+\n\\s+expected: 2s\n\\s+actual\\s+: 1ms\n"
	Regexp(t, regexp.MustCompile(expectedErr), out.buf.String())
}

func TestDiffEmptyCases(t *testing.T) {
	New(t).Equal("", diff(nil, nil))
	New(t).Equal("", diff(struct{ foo string }{}, nil))
	New(t).Equal("", diff(nil, struct{ foo string }{}))
	New(t).Equal("", diff(1, 2))
	New(t).Equal("", diff(1, 2))
	New(t).Equal("", diff([]int{1}, []bool{true}))
}

// Ensure there are no data races
func TestDiffRace(t *testing.T) {
	t.Parallel()

	expected := map[string]string{
		"a": "A",
		"b": "B",
		"c": "C",
	}

	actual := map[string]string{
		"d": "D",
		"e": "E",
		"f": "F",
	}

	// run diffs in parallel simulating tests with t.Parallel()
	numRoutines := 10
	rChans := make([]chan string, numRoutines)
	for idx := range rChans {
		rChans[idx] = make(chan string)
		go func(ch chan string) {
			defer close(ch)
			ch <- diff(expected, actual)
		}(rChans[idx])
	}

	for _, ch := range rChans {
		for msg := range ch {
			NotZero(t, msg) // dummy assert
		}
	}
}

func TestFailNow(t *testing.T) {
	out := &outputT{buf: bytes.NewBuffer(nil)}
	mockAssertion := New(out)
	New(t).False(mockAssertion.FailNow("failed"))
}

func TestBytesEqual(t *testing.T) {
	var cases = []struct {
		a, b []byte
	}{
		{make([]byte, 2), make([]byte, 2)},
		{make([]byte, 2), make([]byte, 2, 3)},
		{nil, make([]byte, 0)},
	}
	for i, c := range cases {
		New(t).Equal(reflect.DeepEqual(c.a, c.b), ObjectsAreEqual(c.a, c.b), "case %d failed", i+1)
	}
}

func BenchmarkBytesEqual(b *testing.B) {
	const size = 1024 * 8
	s := make([]byte, size)
	for i := range s {
		s[i] = byte(i % 255)
	}
	s2 := make([]byte, size)
	copy(s2, s)

	mockT := new(testing.T)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		NewWithOnFailureNoop(mockT).Equal(s, s2)
	}
}

func BenchmarkNotNil(b *testing.B) {
	for i := 0; i < b.N; i++ {
		New(b).NotNil(b)
	}
}

func TestEventuallyFalse(t *testing.T) {
	assertion := New(t)
	mockT := new(testing.T)

	condition := func() bool {
		return false
	}

	assertion.False(Eventually(mockT, condition, 100*time.Millisecond, 20*time.Millisecond))
}

func TestEventuallyTrue(t *testing.T) {
	assertion := New(t)
	state := 0
	condition := func() bool {
		defer func() {
			state += 1
		}()
		return state == 2
	}

	assertion.True(Eventually(t, condition, 100*time.Millisecond, 20*time.Millisecond))
}

func TestNeverFalse(t *testing.T) {
	assertion := New(t)
	condition := func() bool {
		return false
	}

	assertion.True(Never(t, condition, 100*time.Millisecond, 20*time.Millisecond))
}

func TestNeverTrue(t *testing.T) {
	assertion := New(t)
	mockT := new(testing.T)
	state := 0
	condition := func() bool {
		defer func() {
			state = state + 1
		}()
		return state == 2
	}

	assertion.False(Never(mockT, condition, 100*time.Millisecond, 20*time.Millisecond))
}

func TestEventuallyIssue805(t *testing.T) {
	assertion := New(t)
	mockT := new(testing.T)

	NotPanics(t, func() {
		condition := func() bool { <-time.After(time.Millisecond); return true }
		assertion.False(Eventually(mockT, condition, time.Millisecond, time.Microsecond))
	})
}

func Test_validateEqualArgs(t *testing.T) {
	if validateEqualArgs(func() {}, func() {}) == nil {
		t.Error("non-nil functions should error")
	}

	if validateEqualArgs(func() {}, func() {}) == nil {
		t.Error("non-nil functions should error")
	}

	if validateEqualArgs(nil, nil) != nil {
		t.Error("nil functions are equal")
	}
}

func TestTruncatingFormat(t *testing.T) {
	original := strings.Repeat("a", bufio.MaxScanTokenSize-102)
	result := truncatingFormat(original)
	New(t).Equal(fmt.Sprintf("%#v", original), result, "string should not be truncated")

	original = original + "x"
	result = truncatingFormat(original)
	New(t).NotEqual(fmt.Sprintf("%#v", original), result, "string should have been truncated.")

	if !strings.HasSuffix(result, "<... truncated>") {
		t.Error("truncated string should have <... truncated> suffix")
	}
}

func TestErrorIs(t *testing.T) {
	mockT := new(testing.T)
	tests := []struct {
		err    error
		target error
		result bool
	}{
		{io.EOF, io.EOF, true},
		{fmt.Errorf("wrap: %w", io.EOF), io.EOF, true},
		{io.EOF, io.ErrClosedPipe, false},
		{nil, io.EOF, false},
		{io.EOF, nil, false},
		{nil, nil, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("ErrorIs(%#v,%#v)", tt.err, tt.target), func(t *testing.T) {
			res := ErrorIs(mockT, tt.err, tt.target)
			if res != tt.result {
				t.Errorf("ErrorIs(%#v,%#v) should return %t", tt.err, tt.target, tt.result)
			}
		})
	}
}

func TestNotErrorIs(t *testing.T) {
	mockT := new(testing.T)
	tests := []struct {
		err    error
		target error
		result bool
	}{
		{io.EOF, io.EOF, false},
		{fmt.Errorf("wrap: %w", io.EOF), io.EOF, false},
		{io.EOF, io.ErrClosedPipe, true},
		{nil, io.EOF, true},
		{io.EOF, nil, true},
		{nil, nil, false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(fmt.Sprintf("NotErrorIs(%#v,%#v)", tt.err, tt.target), func(t *testing.T) {
			res := NotErrorIs(mockT, tt.err, tt.target)
			if res != tt.result {
				t.Errorf("NotErrorIs(%#v,%#v) should return %t", tt.err, tt.target, tt.result)
			}
		})
	}
}

func TestErrorAs(t *testing.T) {
	mockT := new(testing.T)
	tests := []struct {
		err    error
		result bool
	}{
		{fmt.Errorf("wrap: %w", &customError{}), true},
		{io.EOF, false},
		{nil, false},
	}
	for _, tt := range tests {
		tt := tt
		var target *customError
		t.Run(fmt.Sprintf("ErrorAs(%#v,%#v)", tt.err, target), func(t *testing.T) {
			res := ErrorAs(mockT, tt.err, &target)
			if res != tt.result {
				t.Errorf("ErrorAs(%#v,%#v) should return %t)", tt.err, target, tt.result)
			}
		})
	}
}
