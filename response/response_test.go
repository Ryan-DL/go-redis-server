package response

import (
	"testing"
)

func TestSimpleStringSerialize(t *testing.T) {
	expected := "+OK\r\n"
	actual := SimpleString("OK").Serialize()
	if actual != expected {
		t.Errorf("SimpleString Serialize() failed. Expected: %q, got: %q", expected, actual)
	}
}

func TestErrorTypeSerialize(t *testing.T) {
	expected := "-Error message\r\n"
	actual := ErrorType("Error message").Serialize()
	if actual != expected {
		t.Errorf("ErrorType Serialize() failed. Expected: %q, got: %q", expected, actual)
	}
}

func TestIntegerTypeSerialize(t *testing.T) {
	expected := ":123\r\n"
	actual := IntegerType(123).Serialize()
	if actual != expected {
		t.Errorf("IntegerType Serialize() failed. Expected: %q, got: %q", expected, actual)
	}
}

func TestBulkStringTypeSerialize(t *testing.T) {
	expected := "$11\r\nHello World\r\n"
	actual := BulkStringType("Hello World").Serialize()
	if actual != expected {
		t.Errorf("BulkStringType Serialize() failed. Expected: %q, got: %q", expected, actual)
	}
}

func TestNullBulkStringSerialize(t *testing.T) {
	expected := "$-1\r\n"
	actual := NullBulkString{}.Serialize()
	if actual != expected {
		t.Errorf("NullBulkString Serialize() failed. Expected: %q, got: %q", expected, actual)
	}
}

func TestArrayTypeSerialize(t *testing.T) {
	// Test an array with multiple elements
	expected := "*3\r\n+OK\r\n:123\r\n$11\r\nHello World\r\n"
	array := ArrayType{
		SimpleString("OK"),
		IntegerType(123),
		BulkStringType("Hello World"),
	}
	actual := array.Serialize()
	if actual != expected {
		t.Errorf("ArrayType Serialize() failed. Expected: %q, got: %q", expected, actual)
	}

	// Test an empty array
	expected = "*0\r\n"
	array = ArrayType{}
	actual = array.Serialize()
	if actual != expected {
		t.Errorf("ArrayType Serialize() failed for empty array. Expected: %q, got: %q", expected, actual)
	}

	// Test a null array
	expected = "*-1\r\n"
	array = nil
	actual = array.Serialize()
	if actual != expected {
		t.Errorf("ArrayType Serialize() failed for null array. Expected: %q, got: %q", expected, actual)
	}
}
