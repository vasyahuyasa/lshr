package main

import (
	"fmt"
	"log"
	"os"

	"github.com/vasyahuyasa/lshr/reciver"
	"github.com/vasyahuyasa/lshr/sender"
)

const port = 23543

func main() {
	var err error

	// os.Args[0] is always programm name
	if len(os.Args) == 1 {
		err = recive()
	} else {
		filename := os.Args[1]
		err = send(filename)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func send(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("can not send file %s: %v", filename, err)
	}

	info, err := f.Stat()
	if err != nil {
		return err
	}

	s := sender.New(port, filename, f, uint64(info.Size()), 1)
	log.Printf("Start anounce %s to %d port\n", filename, port)
	return s.Anounce()
}

func recive() error {
	log.Printf("Wait incoming files on port %d\n", port)
	r := reciver.New(port)
	return r.Recive()
}
