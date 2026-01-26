package utils

import "reflect"

func NewStringPtr(s string) *string {
	return &s
}

func NewIntPtr(i int) *int {
	return &i
}

func NewBoolPtr(b bool) *bool {
	return &b
}

func NewFloat64Ptr(f float64) *float64 {
	return &f
}

func NewInt64Ptr(i int64) *int64 {
	return &i
}

func NewInt32Ptr(i int32) *int32 {
	return &i
}

func NewInt8Ptr(i int8) *int8 {
	return &i
}

func NewInt16Ptr(i int16) *int16 {
	return &i
}

func NewUintPtr(i uint) *uint {
	return &i
}

func NewUint64Ptr(i uint64) *uint64 {
	return &i
}

func NewUint32Ptr(i uint32) *uint32 {
	return &i
}

func NewUint16Ptr(i uint16) *uint16 {
	return &i
}

func NewUint8Ptr(i uint8) *uint8 {
	return &i
}

func NewBytePtr(b byte) *byte {
	return &b
}

func NewRunePtr(r rune) *rune {
	return &r
}

func NewComplex64Ptr(c complex64) *complex64 {
	return &c
}
func NewPtr[T any](v T) *T {
	return &v
}

func IsValueNil(i interface{}) bool {
	if i == nil {
		return true
	}

	// Lấy reflect.Value của i
	v := reflect.ValueOf(i)

	// Nếu không phải kiểu có thể nil (pointer, map, slice, func, chan, interface) thì luôn luôn không nil
	switch v.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Interface, reflect.Func, reflect.Chan:
		return v.IsNil()
	}

	return false
}
