package gorial

import (
	"testing"
)

func TestNewSerial(t *testing.T) {
	s9600, err := NewSerial(9600, "COM12")
	if s9600.Mode.BaudRate != 9600 || s9600.Name != "COM12" {
		t.Fail()
	} else if err != nil {
		t.Fail()
	}

	_, err = NewSerial(38400, "ABCD")
	if err == nil {
		t.Error("new serial with and invalid name should return an error")
	}
}
