package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dio-av/gorial"
)

var (
	baud     int
	portName string
)

func init() {
	flag.IntVar(&baud, "baud", 115_200, "select com port baudrate")
	flag.StringVar(&portName, "port", "", "select wich port to open")
	flag.Parse()
}

func main() {
	_, err := gorial.GetPorts()
	if err != nil {
		log.Fatalln(err)
	}

	if portName == "" {
		log.Fatalln("port name cannot be empty")
	}
	s, err := gorial.NewSerial(baud, portName)
	if err != nil {
		log.Fatalln(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	s.StartAsyncSerial()

	interrupt := <-c
	fmt.Println("Got signal:", interrupt)
	close(c)
	err = (*s.Port).Close()
	fmt.Println("closing COM port", s.Name)
	if err != nil {
		log.Fatalln("error closing serial:", s.Name, err)
	}
	os.Exit(1)
}
