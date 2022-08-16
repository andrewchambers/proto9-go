package proto9

import (
	"io"
	"net"
	"sync"
)

type Filesystem interface {
	Fcall(Fcall) Fcall
	Clunk()
}

func Serve(l net.Listener, makeFilesystem func() Filesystem) error {
	wg := &sync.WaitGroup{}
	defer wg.Done()
	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			return err
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			ServeConn(c, makeFilesystem())
		}()
	}
}

func ServeConn(rwc io.ReadWriteCloser, fs Filesystem) {
	wg := &sync.WaitGroup{}

	defer func() {
		_ = rwc.Close()
		wg.Done()
		fs.Clunk()
	}()

	msize := uint32(4096)

	fc, err := ReadFcall(msize, rwc)
	switch fc := fc.(type) {
	case *Tversion:
		switch rVersion := fs.Fcall(fc).(type) {
		case *Rversion:
			msize = rVersion.Msize
			err = WriteFcall(rVersion, msize, rwc)
			if err != nil || rVersion.Version == "unknown" {
				return
			}
		default:
			return
		}
	default:
		return
	}

	for {
		// XXX integrate buffer pool.
		fc, err := ReadFcall(msize, rwc)
		if err != nil {
			return
		}
		wg.Add(1)
		// XXX investigate performance of reusing goroutines.
		go func() {
			defer wg.Done()
			resp := fs.Fcall(fc)
			_ = WriteFcall(resp, msize, rwc)
		}()
	}
}

type DotLFile interface {
	Remove() error
	Clunk() error
}

type DotLFilesystem struct {
	Msize uint32

	filesLock sync.RWMutex
	files     map[uint32]DotLFile
}

func (fs *DotLFilesystem) Fcall(fc Fcall) Fcall {
	switch fc := fc.(type) {
	case *Tversion:
		rVersion := &Rversion{
			Tagged: fc.Tagged,
		}
		if fs.Msize < fc.Msize {
			rVersion.Msize = fs.Msize
		} else {
			rVersion.Msize = fc.Msize
		}
		if fc.Version == "9P2000.L" {
			rVersion.Version = fc.Version
		} else {
			rVersion.Version = "unknown"
		}
		return rVersion
	case *Tattach:
		panic("TODO")
	default:
		return &Rlerror{
			Ecode: 123, // XXX TODO error codes.
		}
	}
}

func (fs *DotLFilesystem) Clunk() {
	fs.filesLock.RLock()
	defer fs.filesLock.RUnlock()

	for _, f := range fs.files {
		f.Clunk()
	}
}
