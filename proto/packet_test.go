package proto

import (
	"bytes"
	"crypto/md5"
	"io"
	"testing"
)

const (
	normalNum = 43534
)

func Test_readuint16(t *testing.T) {
	type args struct {
		r io.ByteReader
	}
	tests := []struct {
		name    string
		args    args
		want    uint16
		wantErr bool
	}{
		{
			name:    "simple",
			args:    args{r: bytes.NewBuffer(uint16buf(normalNum))},
			want:    normalNum,
			wantErr: false,
		},
		{
			name:    "EOF",
			args:    args{r: &bytes.Buffer{}},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readuint16(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("readuint16() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readuint16() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readuint64(t *testing.T) {
	type args struct {
		r io.ByteReader
	}
	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{
			name:    "normal",
			args:    args{r: bytes.NewBuffer(uint64buf(normalNum))},
			want:    normalNum,
			wantErr: false,
		},
		{
			name:    "EOF",
			args:    args{r: &bytes.Buffer{}},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readuint64(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("readuint64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readuint64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnonunce_UnmarshalBinary(t *testing.T) {
	normalAnounce := Anonunce{
		Version:   Version,
		Filename:  "Test file name",
		FileHash:  [md5.Size]byte{1, 2, 3, 4, 5, 5, 6},
		UniqID:    234234,
		TotalSize: 43234234,
		NumBlocks: 34534543,
	}
	normalAnounceData, err := normalAnounce.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  Anonunce
		args    args
		wantErr bool
	}{
		{
			name:    "noraml",
			fields:  normalAnounce,
			args:    args{data: normalAnounceData},
			wantErr: false,
		},
		{
			name:    "empty data",
			fields:  normalAnounce,
			args:    args{data: []byte{}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Anonunce{
				Version:   tt.fields.Version,
				Filename:  tt.fields.Filename,
				FileHash:  tt.fields.FileHash,
				UniqID:    tt.fields.UniqID,
				TotalSize: tt.fields.TotalSize,
				NumBlocks: tt.fields.NumBlocks,
			}
			if err := a.UnmarshalBinary(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Anonunce.UnmarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestData_UnmarshalBinary(t *testing.T) {
	x := []byte{8, 4, 5, 6, 7, 4, 3, 2, 1}
	data := Data{
		UniqID:    3243534,
		BlockNum:  343534,
		BlockHash: [md5.Size]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4, 5},
		Size:      uint64(len(x)),
		Payload:   x,
	}
	b, err := data.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	anounce := Anonunce{}
	ab, err := anounce.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  Data
		args    args
		wantErr bool
	}{
		{
			name:    "normal",
			fields:  data,
			args:    args{data: b},
			wantErr: false,
		},
		{
			name:    "wrong type",
			fields:  data,
			args:    args{data: ab},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Data{
				UniqID:    tt.fields.UniqID,
				BlockNum:  tt.fields.BlockNum,
				BlockHash: tt.fields.BlockHash,
				Size:      tt.fields.Size,
				Payload:   tt.fields.Payload,
			}
			if err := d.UnmarshalBinary(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("Data.UnmarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
