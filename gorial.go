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

type serialResponse struct {
	b   []byte
	err error
}

func NewSerial(baud int, com string) (*Serial, error) {
	p, err := serial.Open(com, &serial.Mode{BaudRate: baud})
	if err != nil {
		return &Serial{}, fmt.Errorf("error creating new serial: %s %q", com, err)
	}
	return &Serial{
		Mode: serial.Mode{
			BaudRate: baud,
		},
		Name: com,
		Port: p,
	}, nil
}

func (s *Serial) ChangeMode(m *serial.Mode) error {
	err := s.Port.SetMode(m)
	return fmt.Errorf("error changing serial mode %q", err)
}

func (s *Serial) WritePort(message string) error {
	n, err := s.Port.Write([]byte(message + "\r\n"))
	if err != nil {
		return fmt.Errorf("error %q writing to port %v", err, s.Name)
	}
	fmt.Printf("%d bytes written to %v port\n", n, s.Name)
	return nil
}

func (s *Serial) ReadPort(sr chan serialResponse) {
	buff := make([]byte, BUFFER_READ)
	r := &serialResponse{}
	for {
		n, err := s.Port.Read(buff)
		if err != nil {
			r.b = []byte{}
			r.err = fmt.Errorf("error reading port %s %q", s.Name, err)
			sr <- *r
		}
		if n == 0 {
			fmt.Println("\nEOF")
			break
		}
		fmt.Printf("Received %d bytes from %s: %s\n", n, s.Name, buff)
		if strings.ContainsRune(string(buff), '\r') {
			r.b = buff
			sr <- *r
		}
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
	c := make(chan serialResponse)
	go s.ReadPort(c)
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
