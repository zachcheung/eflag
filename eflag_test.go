package eflag

import (
	"os"
	"testing"
)

func TestParseWithPrefix(t *testing.T) {
	var myBool bool
	var myInt int
	var myString string

	os.Setenv("PREFIX_MYBOOL", "true")
	os.Setenv("PREFIX_MY_INT_ENV", "1")
	os.Setenv("PREFIX_MYSTRING", "custom_value")

	Var(&myBool, "mybool", false, "Description for mybool flag", "")
	Var(&myInt, "myint", 0, "Description for myint flag", "PREFIX_MY_INT_ENV")
	Var(&myString, "mystring", "default", "Description for mystring flag", "-")

	SetPrefix("PREFIX_")
	Parse()

	if !myBool {
		t.Error("Expected myBool to be true, but it's false.")
	}

	if myInt != 1 {
		t.Errorf("Expected myInt to be 1, but got %d", myInt)
	}

	if myString != "default" {
		t.Errorf("Expected myString to be 'default', but got '%s'", myString)
	}
}

func TestParseWithoutPrefix(t *testing.T) {
	var testBool bool
	var testInt int
	var testString string

	os.Setenv("TESTBOOL", "true")
	os.Setenv("TEST_INT_ENV", "1")
	os.Setenv("TESTSTRING", "custom_value")

	Var(&testBool, "testbool", false, "Description for testbool flag", "")
	Var(&testInt, "testint", 0, "Description for testint flag", "TEST_INT_ENV")
	Var(&testString, "teststring", "default", "Description for teststring flag", "-")

	// reset prefix here
	SetPrefix("")
	Parse()

	if !testBool {
		t.Error("Expected testBool to be true, but it's false.")
	}

	if testInt != 1 {
		t.Errorf("Expected testInt to be 1, but got %d", testInt)
	}

	if testString != "default" {
		t.Errorf("Expected testString to be 'default', but got '%s'", testString)
	}
}
