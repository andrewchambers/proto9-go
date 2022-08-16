package proto9

//go:generate go run ./cmd/encdec-codegen/main.go github.com/andrewchambers/proto9-go encdec.gen.go
//go:generate gofmt -w encdec.gen.go

import (
	"bytes"
	"errors"
)

var (
	ErrValueTooLong = errors.New("value too large for 9p message")
)

func encodeByte(b *bytes.Buffer, v byte) error {
	return b.WriteByte(v)
}

func encodeUint8(b *bytes.Buffer, v uint8) error {
	return b.WriteByte(v)
}

func encodeUint16(b *bytes.Buffer, v uint16) error {
	buf := [2]byte{byte(v & 0x00ff), byte((v & 0xff00) >> 8)}
	_, err := b.Write(buf[:])
	return err
}

func encodeUint32(b *bytes.Buffer, v uint32) error {
	buf := [4]byte{
		byte(v & 0xff),
		byte((v & 0xff00) >> 8),
		byte((v & 0xff0000) >> 16),
		byte((v & 0xff000000) >> 24),
	}
	_, err := b.Write(buf[:])
	return err
}

func encodeUint64(b *bytes.Buffer, v uint64) error {
	buf := [8]byte{
		byte(v & 0xff),
		byte((v & 0xff00) >> 8),
		byte((v & 0xff0000) >> 16),
		byte((v & 0xff000000) >> 24),
		byte((v & 0xff00000000) >> 32),
		byte((v & 0xff0000000000) >> 40),
		byte((v & 0xff000000000000) >> 48),
		byte((v & 0xff00000000000000) >> 56),
	}
	_, err := b.Write(buf[:])
	return err
}

func encodeString(b *bytes.Buffer, v string) error {
	if len(v) > 0xffff {
		return ErrValueTooLong
	}
	err := encodeUint32(b, uint32(len(v)))
	if err != nil {
		return err
	}
	_, err = b.WriteString(v)
	return err
}

func encodeByteSlice(b *bytes.Buffer, v []byte) error {
	if len(v) > 0xffffffff {
		return ErrValueTooLong
	}
	return nil
}

func encodeQids(b *bytes.Buffer, v []Qid) error {
	if len(v) > 13 {
		return ErrValueTooLong
	}
	err := encodeUint32(b, uint32(len(v)))
	if err != nil {
		return err
	}
	for i := 0; i < len(v); i++ {
		err = v[i].Encode(b)
		if err != nil {
			return err
		}
	}
	return nil
}

func decodeByte(b *bytes.Buffer) (byte, error) {
	return b.ReadByte()
}

func decodeUint8(b *bytes.Buffer) (uint8, error) {
	return b.ReadByte()
}

func decodeUint16(b *bytes.Buffer) (uint16, error) {
	buf := [2]byte{}
	_, err := b.Read(buf[:])
	if err != nil {
		return 0, err
	}
	return uint16(buf[0]) | (uint16(buf[1]) << 8), nil
}

func decodeUint32(b *bytes.Buffer) (uint32, error) {
	buf := [4]byte{}
	_, err := b.Read(buf[:])
	if err != nil {
		return 0, err
	}
	return uint32(buf[0]) | (uint32(buf[1]) << 8) | (uint32(buf[2]) << 16) | (uint32(buf[3]) << 24), nil
}

func decodeUint64(b *bytes.Buffer) (uint64, error) {
	buf := [8]byte{}
	_, err := b.Read(buf[:])
	if err != nil {
		return 0, err
	}
	return uint64(buf[0]) | (uint64(buf[1]) << 8) | (uint64(buf[2]) << 16) | (uint64(buf[3]) << 24) | (uint64(buf[4]) << 32) | (uint64(buf[5]) << 40) | (uint64(buf[6]) << 48) | (uint64(buf[7]) << 56), nil
}

func decodeString(b *bytes.Buffer) (string, error) {
	l, err := decodeUint16(b)
	if err != nil {
		return "", err
	}
	return string(b.Next(int(l))), nil
}

func decodeByteSlice(b *bytes.Buffer) ([]byte, error) {
	l, err := decodeUint32(b)
	if err != nil {
		return nil, err
	}
	// XXX Technically this is not allowed because the buffer documentation
	// says the buffer is no longer valid on the next call to read... however
	// the implementation doesn't actually invalidate it and we avoid the copy.
	return b.Next(int(l)), nil
}

func decodeQids(b *bytes.Buffer) ([]Qid, error) {
	l, err := decodeUint16(b)
	if err != nil {
		return nil, err
	}
	qids := []Qid{}
	qid := Qid{}
	for i := 0; i < int(l); i++ {
		err = qid.Decode(b)
		if err != nil {
			return nil, err
		}
		qids = append(qids, qid)
	}
	return qids, nil
}
