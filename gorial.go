package gorial

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"go.bug.st/serial"
)

const (
	BUFFER_READ = 256
)

type Serial struct {
	Mode *serial.Mode
	Port *serial.Port
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
		Mode: &serial.Mode{
			BaudRate: baud,
		},
		Name: com,
		Port: &p,
	}, nil
}

func (s *Serial) ChangeMode(m *serial.Mode) error {
	err := (*s.Port).SetMode(m)
	return fmt.Errorf("error changing serial mode %q", err)
}

func (s *Serial) WritePort(message string) error {
	n, err := (*s.Port).Write([]byte(message + "\r\n"))
	if err != nil {
		return fmt.Errorf("error %q writing to port %v", err, s.Name)
	}
	fmt.Printf("%d bytes written to %v port\n", n, s.Name)
	return nil
}

func (s *Serial) ReadPort(data chan []byte) error {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		bufferTotal := make([]byte, 4096)
		defer wg.Done()
		for {
			buff := make([]byte, 1024)
			if s.Port != nil {
				n, err := (*s.Port).Read(buff)
				if err != nil {
					log.Println("error reading serial", s.Name, err)
				}
				if n == 0 {
					break
				}
				fmt.Printf("received %s\n", string(buff[:n]))
				buff = []byte(strings.Trim(string(buff), "\x00/"))
				bufferTotal = append(bufferTotal, buff...)

				if strings.Contains(string(bufferTotal), "\n") {
					fmt.Printf("buffer total: %s\n", string(bufferTotal))
					data <- bufferTotal
					bufferTotal = []byte{}
				}
			}
		}
	}()
	wg.Wait()
	return nil
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
	c := make(chan []byte)
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
		err := (*s.Port).Close()
		fmt.Println("closing COM port", s.Port)
		if err != nil {
			log.Fatalln("error closing serial:", s.Name, err)
		}
	}
	return nil
}
