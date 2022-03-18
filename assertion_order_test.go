package assert

import (
	"bytes"
	"testing"
)

func TestIsIncreasing(t *testing.T) {
	mockAssertion := New(new(testing.T)).WithOnFailure(func(t TestingT) any { return -1 })
	New(t).Nil(mockAssertion.IsIncreasing([]int{1, 2}))
	New(t).Nil(mockAssertion.IsIncreasing([]int{1, 2, 3, 4, 5}))
	New(t).NotNil(mockAssertion.IsIncreasing([]int{1, 1}))
	New(t).NotNil(mockAssertion.IsIncreasing([]int{2, 1}))

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
		outAssertion := New(out).WithOnFailure(func(t TestingT) any { return -1 })
		New(t).NotNil(outAssertion.IsIncreasing(currCase.collection))
		New(t).Contains(out.buf.String(), currCase.msg)
	}
}

func TestIsNonIncreasing(t *testing.T) {
	mockAssertion := New(new(testing.T)).WithOnFailure(func(t TestingT) any { return -1 })
	New(t).Nil(mockAssertion.IsNonIncreasing([]int{2, 1}))
	New(t).Nil(mockAssertion.IsNonIncreasing([]int{5, 4, 4, 3, 2, 1}))
	New(t).Nil(mockAssertion.IsNonIncreasing([]int{1, 1}))
	New(t).NotNil(mockAssertion.IsNonIncreasing([]int{1, 2}))

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
		outAssertion := New(out).WithOnFailure(func(t TestingT) any { return -1 })
		New(t).NotNil(outAssertion.IsNonIncreasing(currCase.collection))
		New(t).Contains(out.buf.String(), currCase.msg)
	}
}

func TestIsDecreasing(t *testing.T) {
	mockAssertion := New(new(testing.T)).WithOnFailure(func(t TestingT) any { return -1 })
	New(t).Nil(mockAssertion.IsDecreasing([]int{2, 1}))
	New(t).Nil(mockAssertion.IsDecreasing([]int{5, 4, 3, 2, 1}))
	New(t).NotNil(mockAssertion.IsDecreasing([]int{1, 1}))
	New(t).NotNil(mockAssertion.IsDecreasing([]int{1, 2}))

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
		outAssertion := New(out).WithOnFailure(func(t TestingT) any { return -1 })
		New(t).NotNil(outAssertion.IsDecreasing(currCase.collection))
		New(t).Contains(out.buf.String(), currCase.msg)
	}
}

func TestIsNonDecreasing(t *testing.T) {
	mockAssertion := New(new(testing.T)).WithOnFailure(func(t TestingT) any { return -1 })
	New(t).Nil(mockAssertion.IsNonDecreasing([]int{1, 2}))
	New(t).Nil(mockAssertion.IsNonDecreasing([]int{1, 1, 2, 3, 4, 5}))
	New(t).Nil(mockAssertion.IsNonDecreasing([]int{1, 1}))
	New(t).NotNil(mockAssertion.IsNonDecreasing([]int{2, 1}))

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
		outAssertion := New(out).WithOnFailure(func(t TestingT) any { return -1 })
		New(t).NotNil(outAssertion.IsNonDecreasing(currCase.collection))
		New(t).Contains(out.buf.String(), currCase.msg)
	}
}

func TestOrderingMsgAndArgsForwarding(t *testing.T) {
	msgAndArgs := []any{"format %s %x", "this", 0xc001}
	expectedOutput := "format this c001\n"
	collection := []int{1, 2, 1}
	funcs := []func(*Assertions) any{
		func(a *Assertions) any { return a.IsIncreasing(collection, msgAndArgs...) },
		func(a *Assertions) any { return a.IsNonIncreasing(collection, msgAndArgs...) },
		func(a *Assertions) any { return a.IsDecreasing(collection, msgAndArgs...) },
		func(a *Assertions) any { return a.IsNonDecreasing(collection, msgAndArgs...) },
	}
	for _, f := range funcs {
		out := &outputT{buf: bytes.NewBuffer(nil)}
		outAssertion := New(out).WithOnFailure(func(t TestingT) any { return -1 })
		New(t).NotNil(f(outAssertion))
		New(t).Contains(out.buf.String(), expectedOutput)
	}
}
