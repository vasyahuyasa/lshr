package service

import (
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
	"net"
	"time"

	"github.com/vasyahuyasa/lshr/proto"
)

type state byte

const (
	statePrepare state = iota
	stateAnounce
	stateRedeliver
	stateDone
)

type Sender struct {
	port     int
	filename string
	size     uint64
	blocks   uint64
	r        io.Reader
	state    state
	conn     *net.UDPConn
}

func (s *Sender) changeState(st state) {
	s.state = st
}

func (s *Sender) broadcastState(st state) {

}

func NewSender(port int, filename string, r io.Reader, size uint64, blocks uint64) (*Sender, error) {
	conn, err := makeBroadcastConn()
	if err != nil {
		return nil, err
	}

	return &Sender{
		port:     port,
		filename: filename,
		size:     size,
		blocks:   blocks,
		r:        r,
		state:    statePrepare,
		conn:     conn,
	}, nil
}

func makeBroadcastConn() (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", net.IPv4bcast, s.port))
	if err != nil {
		return err
	}

	return net.DialUDP("udp", nil, addr)
}

func (s *Sender) Anounce() error {
	rand.Seed(time.Now().Unix())

	anounce := &proto.Anonunce{
		Version:   proto.Version,
		Filename:  s.filename,
		FileHash:  md5.Sum([]byte("test file hash")),
		UniqID:    rand.Uint64(),
		TotalSize: s.size,
		NumBlocks: s.blocks,
	}

	data, err := anounce.MarshalBinary()
	if err != nil {
		return err
	}

	s.state = stateAnounce
	for {
		_, err = conn.Write(data)
		if err != nil {
			return err
		}
		time.Sleep(time.Second)
	}
}
