package proto

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"errors"
	"io"
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
	Size uint32

	Payload []byte
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
