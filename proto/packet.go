package proto

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"errors"
	"io"
	"log"
)

// Proto is current protocol version
const (
	Version = 1

	PacketAnounce byte = 0x0A + iota
	PacketData
)

var (
	ErrPacketTypeMistmatch = errors.New("packet type mistmatch")

	ErrVersionMistmatch = errors.New("protocol version mistmatch")

	ErrFilenameLengthMistmatch = errors.New("filename length mistmatch")

	ErrFileHashLengthMistmatch = errors.New("file hash length mistmatch")

	ErrPayloadSizeMistmatch = errors.New("payload size mistmatch")
)

// Anonunce is structure for represent of sending file
type Anonunce struct {
	// Version is current protocol version
	Version byte

	// Filename is name of transfering file
	Filename string

	// FileHash is md5 sum of transfering file
	FileHash [md5.Size]byte

	// Uniq is unique random generated id of current transfer session
	UniqID uint64

	// TotalSize is file size in bytes
	TotalSize uint64

	// NumBlocks is number data pices which file is devided
	NumBlocks uint64
}

// Data is data block with payload
type Data struct {
	// UniqID is related file id
	UniqID uint64

	// BlockNum is this block number started from zero
	BlockNum uint64

	// BlockHash is md5 hash of following payload
	BlockHash [md5.Size]byte

	// Size is size of folowwing payload in bytes
	Size uint64

	Payload []byte
}

func (d *Data) MarshalBinary() ([]byte, error) {
	buf := bytes.Buffer{}

	// packet type
	err := buf.WriteByte(PacketData)
	if err != nil {
		return nil, err
	}

	// UniqID
	_, err = buf.Write(uint64buf(d.UniqID))
	if err != nil {
		return nil, err
	}

	// BlockNum
	_, err = buf.Write(uint64buf(d.BlockNum))
	if err != nil {
		return nil, err
	}

	// BlockHash
	_, err = buf.Write(d.BlockHash[:])
	if err != nil {
		return nil, err
	}

	// Size
	_, err = buf.Write(uint64buf(d.Size))
	if err != nil {
		return nil, err
	}

	// payload
	_, err = buf.Write(d.Payload)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (d *Data) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	var err error

	// check packet type
	packet, err := buf.ReadByte()
	if err != nil {
		return err
	}
	if packet != PacketData {
		return ErrPacketTypeMistmatch
	}

	//UniqID
	d.UniqID, err = readuint64(buf)
	if err != nil {
		return err
	}

	//BlockNum
	d.BlockNum, err = readuint64(buf)
	if err != nil {
		return err
	}

	//BlockHash
	blockahsh := make([]byte, md5.Size)
	n, err := buf.Read(blockahsh)
	if err != nil {
		return err
	}
	if n != md5.Size {
		return ErrFileHashLengthMistmatch
	}
	copy(d.BlockHash[:], blockahsh)

	//Size
	d.Size, err = readuint64(buf)
	if err != nil {
		return err
	}

	//Payload
	d.Payload = make([]byte, d.Size)
	n, err = buf.Read(d.Payload)
	if err != nil {
		return err
	}
	if uint64(n) != d.Size {

		log.Printf("read from data: %d d.Size: %d len(data): %d\n", n, d.Size, len(d.Payload))

		return ErrPayloadSizeMistmatch
	}

	return nil
}

func (a *Anonunce) MarshalBinary() ([]byte, error) {
	buf := bytes.Buffer{}

	// packet type
	err := buf.WriteByte(PacketAnounce)
	if err != nil {
		return nil, err
	}

	// proto version
	err = buf.WriteByte(a.Version)
	if err != nil {
		return nil, err
	}

	// filename
	_, err = buf.Write(uint16buf(uint16(len(a.Filename))))
	if err != nil {
		return nil, err
	}
	_, err = buf.WriteString(a.Filename)
	if err != nil {
		return nil, err
	}

	// file hash
	_, err = buf.Write(a.FileHash[:])
	if err != nil {
		return nil, err
	}

	// UniqID
	_, err = buf.Write(uint64buf(a.UniqID))
	if err != nil {
		return nil, err
	}

	// TotalSize
	_, err = buf.Write(uint64buf(a.TotalSize))
	if err != nil {
		return nil, err
	}

	// NumBlocks
	_, err = buf.Write(uint64buf(a.NumBlocks))
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (a *Anonunce) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	var err error

	// check packet type
	packet, err := buf.ReadByte()
	if err != nil {
		return err
	}
	if packet != PacketAnounce {
		return ErrPacketTypeMistmatch
	}

	// version
	version, err := buf.ReadByte()
	if err != nil {
		return err
	}

	if version != Version {
		return ErrVersionMistmatch
	}

	a.Version = version

	// filename
	size, err := readuint16(buf)
	if err != nil {
		return err
	}

	filename := make([]byte, size)
	n, err := buf.Read(filename)
	if err != nil {
		return err
	}
	if n != int(size) {
		return ErrFilenameLengthMistmatch
	}

	a.Filename = string(filename)

	// filehash
	filehash := make([]byte, md5.Size)
	n, err = buf.Read(filehash)
	if err != nil {
		return err
	}
	if n != md5.Size {
		return ErrFileHashLengthMistmatch
	}
	copy(a.FileHash[:], filehash)

	// UniqID
	a.UniqID, err = readuint64(buf)
	if err != nil {
		return err
	}

	// TotalSize
	a.TotalSize, err = readuint64(buf)
	if err != nil {
		return err
	}

	// NumBlocks
	a.NumBlocks, err = readuint64(buf)
	if err != nil {
		return err
	}

	return nil
}

func readuint16(r io.ByteReader) (uint16, error) {
	lo, err := r.ReadByte()
	if err != nil {
		return 0, err
	}

	hi, err := r.ReadByte()
	if err != nil {
		return 0, err
	}

	var i uint16
	i |= uint16(hi)
	i <<= 8
	i |= uint16(lo)
	return i, nil
}

func readuint64(r io.ByteReader) (uint64, error) {
	var n uint64
	var numBytes uint
	for numBytes < 8 {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}
		n |= uint64(b) << (8 * (numBytes))
		numBytes++
	}
	return n, nil
}

func uint16buf(v uint16) []byte {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, v)
	return buf
}

func uint64buf(v uint64) []byte {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, v)
	return buf
}
