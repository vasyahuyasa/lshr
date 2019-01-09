package service

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/dustin/go-humanize"

	"github.com/vasyahuyasa/lshr/proto"
)

const bufSize = 1024 * 1024 * 10 // Buffer size 10MBi for incoming packets

type Reciver struct {
	port   int
	waitMu sync.Mutex
	wait   map[uint64]*proto.Anonunce
}

func NewReciver(port int) *Reciver {
	return &Reciver{
		port: port,
		wait: map[uint64]*proto.Anonunce{},
	}
}

func (recv *Reciver) Recive() error {
	conn, err := net.ListenPacket("udp", fmt.Sprintf(":%d", recv.port))
	if err != nil {
		return err
	}
	defer conn.Close()

	//simple read
	data := make([]byte, bufSize)
	for {
		size, addr, err := conn.ReadFrom(data)
		if err != nil {
			return err
		}
		if size == 0 {
			return errors.New("recived zero length packet")
		}
		log.Printf("recived %d bytes size packet from %s\n", size, addr)

		switch data[0] {
		case proto.PacketAnounce:
			err = recv.processAnonce(data[:size])
			if err != nil {
				log.Println("can not process anonunce packet:", err)
			}
			break
		case proto.PacketData:
			break
		default:
			log.Println("recived unknow packet type", data[0])
		}
	}

	return nil
}

func (recv *Reciver) processAnonce(data []byte) error {
	an := &proto.Anonunce{}

	err := an.UnmarshalBinary(data)
	if err != nil {
		return err
	}

	// if we already waiting for file do nothing
	recv.waitMu.Lock()
	defer recv.waitMu.Unlock()
	if _, ok := recv.wait[an.UniqID]; ok {
		return nil
	}

	log.Printf("id: %d filename: %s hash: %s size: %s, blocks: %d\nAccept y/n (n)? ", an.UniqID, an.Filename, hex.EncodeToString(an.FileHash[:]), humanize.Bytes(an.TotalSize), an.NumBlocks)

	// ask user to accept this file
	reader := bufio.NewReader(os.Stdin)
	yn, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	yn = strings.ToLower(yn)

	if len(yn) > 0 && yn[0] == 'y' {
		log.Printf("file %d [%s] (%s) will be accepted\n", an.UniqID, an.Filename, hex.EncodeToString(an.FileHash[:]))
		recv.wait[an.UniqID] = an
	} else {
		log.Printf("file %d [%s] (%s) will be ignored\n", an.UniqID, an.Filename, hex.EncodeToString(an.FileHash[:]))
	}

	return nil
}
