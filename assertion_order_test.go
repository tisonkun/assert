package assert

import (
	"bytes"
	"testing"
)

func TestIsIncreasing(t *testing.T) {
	assertion := New(t, FailNowOnFailure)

	mockT := new(mockTestingT)
	mockAssertion := New(mockT, FailNowOnFailure)

	mockT.reset()
	mockAssertion.IsIncreasing([]int{1, 2})
	assertion.False(mockT.failed)

	mockT.reset()
	mockAssertion.IsIncreasing([]int{1, 2, 3, 4, 5})
	assertion.False(mockT.failed)

	mockT.reset()
	mockAssertion.IsIncreasing([]int{1, 1})
	assertion.True(mockT.failed)

	mockT.reset()
	mockAssertion.IsIncreasing([]int{2, 1})
	assertion.True(mockT.failed)

	// Check error report
	for _, currCase := range []struct {
		collection any
		msg        string
	}{
		{collection: []string{"b", "a"}, msg: `"b" is not less than "a"`},
		{collection: []int{2, 1}, msg: `"2" is not less than "1"`},
		{collection: []int{2, 1, 3, 4, 5, 6, 7}, msg: `"2" is not less than "1"`},
		{collection: []int{-1, 0, 2, 1}, msg: `"2" is not less than "1"`},
		{collection: []int8{2, 1}, msg: `"2" is not less than "1"`},
		{collection: []int16{2, 1}, msg: `"2" is not less than "1"`},
		{collection: []int32{2, 1}, msg: `"2" is not less than "1"`},
		{collection: []int64{2, 1}, msg: `"2" is not less than "1"`},
		{collection: []uint8{2, 1}, msg: `"2" is not less than "1"`},
		{collection: []uint16{2, 1}, msg: `"2" is not less than "1"`},
		{collection: []uint32{2, 1}, msg: `"2" is not less than "1"`},
		{collection: []uint64{2, 1}, msg: `"2" is not less than "1"`},
		{collection: []float32{2.34, 1.23}, msg: `"2.34" is not less than "1.23"`},
		{collection: []float64{2.34, 1.23}, msg: `"2.34" is not less than "1.23"`},
	} {
		out := &outputT{buf: bytes.NewBuffer(nil)}
		outAssertion := New(out, FailNowOnFailure)
		outAssertion.IsIncreasing(currCase.collection)
		assertion.True(out.failed)
		Contains(t, out.buf.String(), currCase.msg)
	}
}

func TestIsNonIncreasing(t *testing.T) {
	assertion := New(t, FailNowOnFailure)

	mockT := new(mockTestingT)
	mockAssertion := New(mockT, FailNowOnFailure)

	mockT.reset()
	mockAssertion.IsNonIncreasing([]int{2, 1})
	assertion.False(mockT.failed)

	mockT.reset()
	mockAssertion.IsNonIncreasing([]int{5, 4, 4, 3, 2, 1})
	assertion.False(mockT.failed)

	mockT.reset()
	mockAssertion.IsNonIncreasing([]int{1, 1})
	assertion.False(mockT.failed)

	mockT.reset()
	mockAssertion.IsNonIncreasing([]int{1, 2})
	assertion.True(mockT.failed)

	// Check error report
	for _, currCase := range []struct {
		collection any
		msg        string
	}{
		{collection: []string{"a", "b"}, msg: `"a" is not greater than or equal to "b"`},
		{collection: []int{1, 2}, msg: `"1" is not greater than or equal to "2"`},
		{collection: []int{1, 2, 7, 6, 5, 4, 3}, msg: `"1" is not greater than or equal to "2"`},
		{collection: []int{5, 4, 3, 1, 2}, msg: `"1" is not greater than or equal to "2"`},
		{collection: []int8{1, 2}, msg: `"1" is not greater than or equal to "2"`},
		{collection: []int16{1, 2}, msg: `"1" is not greater than or equal to "2"`},
		{collection: []int32{1, 2}, msg: `"1" is not greater than or equal to "2"`},
		{collection: []int64{1, 2}, msg: `"1" is not greater than or equal to "2"`},
		{collection: []uint8{1, 2}, msg: `"1" is not greater than or equal to "2"`},
		{collection: []uint16{1, 2}, msg: `"1" is not greater than or equal to "2"`},
		{collection: []uint32{1, 2}, msg: `"1" is not greater than or equal to "2"`},
		{collection: []uint64{1, 2}, msg: `"1" is not greater than or equal to "2"`},
		{collection: []float32{1.23, 2.34}, msg: `"1.23" is not greater than or equal to "2.34"`},
		{collection: []float64{1.23, 2.34}, msg: `"1.23" is not greater than or equal to "2.34"`},
	} {
		out := &outputT{buf: bytes.NewBuffer(nil)}
		outAssertion := New(out, FailNowOnFailure)
		outAssertion.IsNonIncreasing(currCase.collection)
		assertion.True(out.failed)
		Contains(t, out.buf.String(), currCase.msg)
	}
}

func TestIsDecreasing(t *testing.T) {
	assertion := New(t, FailNowOnFailure)

	mockT := new(mockTestingT)
	mockAssertion := New(mockT, FailNowOnFailure)

	mockT.reset()
	mockAssertion.IsDecreasing([]int{2, 1})
	assertion.False(mockT.failed)

	mockT.reset()
	mockAssertion.IsDecreasing([]int{5, 4, 3, 2, 1})
	assertion.False(mockT.failed)

	mockT.reset()
	mockAssertion.IsDecreasing([]int{1, 1})
	assertion.True(mockT.failed)

	mockT.reset()
	mockAssertion.IsDecreasing([]int{1, 2})
	assertion.True(mockT.failed)

	// Check error report
	for _, currCase := range []struct {
		collection any
		msg        string
	}{
		{collection: []string{"a", "b"}, msg: `"a" is not greater than "b"`},
		{collection: []int{1, 2}, msg: `"1" is not greater than "2"`},
		{collection: []int{1, 2, 7, 6, 5, 4, 3}, msg: `"1" is not greater than "2"`},
		{collection: []int{5, 4, 3, 1, 2}, msg: `"1" is not greater than "2"`},
		{collection: []int8{1, 2}, msg: `"1" is not greater than "2"`},
		{collection: []int16{1, 2}, msg: `"1" is not greater than "2"`},
		{collection: []int32{1, 2}, msg: `"1" is not greater than "2"`},
		{collection: []int64{1, 2}, msg: `"1" is not greater than "2"`},
		{collection: []uint8{1, 2}, msg: `"1" is not greater than "2"`},
		{collection: []uint16{1, 2}, msg: `"1" is not greater than "2"`},
		{collection: []uint32{1, 2}, msg: `"1" is not greater than "2"`},
		{collection: []uint64{1, 2}, msg: `"1" is not greater than "2"`},
		{collection: []float32{1.23, 2.34}, msg: `"1.23" is not greater than "2.34"`},
		{collection: []float64{1.23, 2.34}, msg: `"1.23" is not greater than "2.34"`},
	} {
		out := &outputT{buf: bytes.NewBuffer(nil)}
		outAssertion := New(out, FailNowOnFailure)
		outAssertion.IsDecreasing(currCase.collection)
		assertion.True(out.failed)
		Contains(t, out.buf.String(), currCase.msg)
	}
}

func TestIsNonDecreasing(t *testing.T) {
	assertion := New(t, FailNowOnFailure)

	mockT := new(mockTestingT)
	mockAssertion := New(mockT, FailNowOnFailure)

	mockT.reset()
	mockAssertion.IsNonDecreasing([]int{1, 2})
	assertion.False(mockT.failed)

	mockT.reset()
	mockAssertion.IsNonDecreasing([]int{1, 1, 2, 3, 4, 5})
	assertion.False(mockT.failed)

	mockT.reset()
	mockAssertion.IsNonDecreasing([]int{1, 1})
	assertion.False(mockT.failed)

	mockT.reset()
	mockAssertion.IsNonDecreasing([]int{2, 1})
	assertion.True(mockT.failed)

	// Check error report
	for _, currCase := range []struct {
		collection any
		msg        string
	}{
		{collection: []string{"b", "a"}, msg: `"b" is not less than or equal to "a"`},
		{collection: []int{2, 1}, msg: `"2" is not less than or equal to "1"`},
		{collection: []int{2, 1, 3, 4, 5, 6, 7}, msg: `"2" is not less than or equal to "1"`},
		{collection: []int{-1, 0, 2, 1}, msg: `"2" is not less than or equal to "1"`},
		{collection: []int8{2, 1}, msg: `"2" is not less than or equal to "1"`},
		{collection: []int16{2, 1}, msg: `"2" is not less than or equal to "1"`},
		{collection: []int32{2, 1}, msg: `"2" is not less than or equal to "1"`},
		{collection: []int64{2, 1}, msg: `"2" is not less than or equal to "1"`},
		{collection: []uint8{2, 1}, msg: `"2" is not less than or equal to "1"`},
		{collection: []uint16{2, 1}, msg: `"2" is not less than or equal to "1"`},
		{collection: []uint32{2, 1}, msg: `"2" is not less than or equal to "1"`},
		{collection: []uint64{2, 1}, msg: `"2" is not less than or equal to "1"`},
		{collection: []float32{2.34, 1.23}, msg: `"2.34" is not less than or equal to "1.23"`},
		{collection: []float64{2.34, 1.23}, msg: `"2.34" is not less than or equal to "1.23"`},
	} {
		out := &outputT{buf: bytes.NewBuffer(nil)}
		outAssertion := New(out, FailNowOnFailure)
		outAssertion.IsNonDecreasing(currCase.collection)
		assertion.True(out.failed)
		Contains(t, out.buf.String(), currCase.msg)
	}
}

func TestOrderingMsgAndArgsForwarding(t *testing.T) {
	assertion := New(t, FailNowOnFailure)

	msgAndArgs := []any{"format %s %x", "this", 0xc001}
	expectedOutput := "format this c001\n"
	collection := []int{1, 2, 1}
	funcs := []func(*Assertions){
		func(a *Assertions) { a.IsIncreasing(collection, msgAndArgs...) },
		func(a *Assertions) { a.IsNonIncreasing(collection, msgAndArgs...) },
		func(a *Assertions) { a.IsDecreasing(collection, msgAndArgs...) },
		func(a *Assertions) { a.IsNonDecreasing(collection, msgAndArgs...) },
	}
	for _, f := range funcs {
		out := &outputT{buf: bytes.NewBuffer(nil)}
		outAssertion := New(out, FailNowOnFailure)
		f(outAssertion)
		assertion.True(out.failed)
		Contains(t, out.buf.String(), expectedOutput)
	}
}
