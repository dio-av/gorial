package gorial

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"go.bug.st/serial"
)

const (
	BUFFER_READ = 64
)

type Serial struct {
	Mode serial.Mode
	Port serial.Port
	Name string
}

func NewSerial(baud int, com string) (*Serial, error) {
	var s Serial
	s.Mode.BaudRate = baud
	s.Name = com
	p, err := serial.Open(com, &s.Mode)
	if err != nil {
		return &Serial{}, fmt.Errorf("error creating new serial: %s  %q", com, err)
	}
	s.Port = p
	return &s, nil
}

func (s *Serial) ChangeMode(m *serial.Mode) error {
	err := s.Port.SetMode(m)
	return err
}

func (s *Serial) WritePort(message string) error {
	n, err := s.Port.Write([]byte(message + "\n\r"))
	if err != nil {
		return fmt.Errorf("error %q writing to port %v", err, s.Name)
	}
	fmt.Printf("%d written to %v port\n", n, s.Name)
	return nil
}

func (s *Serial) ReadPort() {
	buff := make([]byte, BUFFER_READ)
	reader := bufio.NewReader(s.Port)
	for {
		n, err := s.Port.Read(buff)
		if err != nil {
			log.Fatal(err)
			break
		}
		if n == 0 {
			fmt.Println("\nEOF")
			break
		} else if n > BUFFER_READ {
			fmt.Println("input size too big")
			break
		}
		reply, err := reader.ReadBytes('\r')
		if err != nil {
			log.Fatal(err)
		}
		r := strings.Trim(string(reply), "\r")
		fmt.Printf("Read %d bytes from %v: %q\n", len(reply), s.Name, string(r))
	}
}

func GetPorts() ([]string, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return []string{}, fmt.Errorf("error get ports list: %q", err)
	}
	if len(ports) == 0 {
		return []string{}, fmt.Errorf("no serial ports found")
	}
	for _, port := range ports {
		fmt.Printf("Found port: %v\n", port)
	}
	return ports, nil
}

func (s *Serial) StartAsyncSerial() {
	go s.ReadPort()
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("write the message you want to send to port: ", s.Name)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Printf("your input message is %s\n", m)
		go s.WritePort(m)
	}
}

func PortsCleanUp(ss []Serial) error {
	for _, s := range ss {
		err := s.Port.Close()
		fmt.Println("closing COM port", s.Port)
		if err != nil {
			log.Fatalln("error closing serial:", s.Name, err)
		}
	}
	return nil
}
