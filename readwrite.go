package proto9

import (
	"bytes"
	"errors"
	"io"
)

func WriteFcall(fc Fcall, msize uint32, w io.Writer) error {
	// TODO No copy fast path for things like Twrite and Rread.
	// TODO Take buffer from a pool.
	var b bytes.Buffer

	if uint32(b.Cap()) < msize {
		b.Grow(int(msize - uint32(b.Cap())))
	}

	hdr := [5]byte{0, 0, 0, 0, fc.Kind()}
	_, err := b.Write(hdr[:])
	if err != nil {
		return err
	}

	err = fc.Encode(&b)
	if err != nil {
		return err
	}

	buf := b.Bytes()
	l := b.Len()
	buf[0] = byte(l & 0xff)
	buf[1] = byte((l & 0xff00) >> 8)
	buf[2] = byte((l & 0xff0000) >> 16)
	buf[3] = byte((l & 0xff000000) >> 24)

	_, err = w.Write(buf)
	return err
}

func ReadFcallInto(msize uint32, r io.Reader, b *bytes.Buffer) (Fcall, error) {

	lr := io.LimitedReader{
		R: r,
		N: 5,
	}

	b.Reset()
	if uint32(b.Cap()) < msize {
		b.Grow(int(msize - uint32(b.Cap())))
	}

	nRead, err := b.ReadFrom(&lr)
	if nRead != 5 {
		return nil, io.EOF
	}

	hdr := b.Next(5)

	sz := uint32(hdr[0]) | (uint32(hdr[1]) << 8) | (uint32(hdr[2]) << 16) | (uint32(hdr[3]) << 24)
	if sz <= 5 || sz > msize {
		return nil, errors.New("message size is outside valid range")
	}

	fc, err := FcallFromKind(hdr[4])
	if err != nil {
		return nil, err
	}

	toRead := int64(sz - 5)
	lr.N = toRead

	nRead, err = b.ReadFrom(&lr)
	if nRead != toRead {
		if err == nil {
			err = io.EOF
		}
		return nil, err
	}

	err = fc.Decode(b)
	if err != nil {
		return nil, err
	}

	return fc, err
}

func ReadFcall(msize uint32, r io.Reader) (Fcall, error) {
	b := bytes.NewBuffer(make([]byte, 0, msize))
	return ReadFcallInto(msize, r, b)
}
