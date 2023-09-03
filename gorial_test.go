package gorial

import (
	"bufio"
	"testing"
)

func TestNewSerial(t *testing.T) {
	s9600, err := NewSerial(9600, "COM12")
	if s9600.Mode.BaudRate != 9600 || s9600.Name != "COM12" {
		t.Fail()
	} else if err != nil {
		t.Fail()
	}
	s9600.Port.Close()

	s38400, err := NewSerial(38400, "ABCD")
	if err == nil {
		t.Error("new serial with and invalid name should return an error")
	}
	s38400.Port.Close()
}

func TestReadPort(t *testing.T) {
	s9600, _ := NewSerial(9600, "COM12")

	c := make(chan serialResponse)
	s9600.ReadPort(c)
	message := []byte("hey I'm sending this over serial")
	expected := "hey I'm sending this over serial\r\n"
	w := bufio.NewWriter(s9600.Port)
	w.Write(message)

	got := (<-c)
	if string(got.b) != expected {
		s9600.Port.Close()
		t.Fatalf("got a message different of what expected. got %v expect %v", got, expected)
	}
	s9600.Port.Close()
}

func TestWritePort(t *testing.T) {
	s9600, _ := NewSerial(9600, "COM12")
	r := bufio.NewReader(s9600.Port)
	expected := "hello... there is someone over there?\r\n"
	got := []byte{}
	s9600.WritePort(expected)
	r.Read(got)

	if string(got) != expected {
		s9600.Port.Close()
		t.Fatalf("got a message different of what expected. got %v expect %v", got, expected)
	}
	s9600.Port.Close()
}
