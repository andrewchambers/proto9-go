package proto9

import (
	"errors"
	"fmt"
	"io"
	"sync"
)

type ClientDotLFile struct {
	Client    *Client
	Fid       uint32
	clunkOnce sync.Once
}

func AttachDotL(c *Client, aname string, uname string) (*ClientDotLFile, Qid, error) {
	if c.Version() != "9P2000.L" {
		return nil, Qid{}, fmt.Errorf("cannot attach to mount, protocol version %q", c.Version())
	}
	fid, err := c.AcquireFid()
	if err != nil {
		return nil, Qid{}, err
	}
	success := false
	defer func() {
		if !success {
			c.ReleaseFid(fid)
		}
	}()

	fc, err := c.Fcall(&Tattach{
		Fid:     fid,
		Afid:    NOFID,
		Aname:   aname,
		Uname:   uname,
		N_uname: 0xFFFFFFFF,
	})

	if err != nil {
		return nil, Qid{}, err
	}
	switch fc := fc.(type) {
	case *Rattach:
		success = true
		return &ClientDotLFile{
			Client: c,
			Fid:    fid,
		}, fc.Qid, nil
	case *Rlerror:
		return nil, Qid{}, fc
	default:
		return nil, Qid{}, fmt.Errorf("protocol error, expected Rattach")
	}
}

func (f *ClientDotLFile) Remove() error {
	var removeErr error
	f.clunkOnce.Do(func() {
		defer f.Client.ReleaseFid(f.Fid)
		fc, err := f.Client.Fcall(&Tremove{
			Fid: f.Fid,
		})
		if err != nil {
			removeErr = err
			return
		}
		switch fc := fc.(type) {
		case *Rremove:
		case *Rlerror:
			removeErr = fc
		default:
			removeErr = fmt.Errorf("protocol error, expected Rremove")
		}
	})
	return removeErr
}

func (f *ClientDotLFile) Clunk() error {
	var clunkErr error
	f.clunkOnce.Do(func() {
		defer f.Client.ReleaseFid(f.Fid)
		fc, err := f.Client.Fcall(&Tclunk{
			Fid: f.Fid,
		})
		if err != nil {
			clunkErr = err
			return
		}
		switch fc := fc.(type) {
		case *Rclunk:
		case *Rlerror:
			clunkErr = fc
		default:
			clunkErr = fmt.Errorf("protocol error, expected Rclunk")
		}
	})
	return clunkErr
}

func (f *ClientDotLFile) walk(wnames []string) (*ClientDotLFile, []Qid, error) {
	fid, err := f.Client.AcquireFid()
	if err != nil {
		return nil, nil, err
	}
	success := false
	defer func() {
		if !success {
			f.Client.ReleaseFid(fid)
		}
	}()
	fc, err := f.Client.Fcall(&Twalk{
		Fid:    f.Fid,
		NewFid: fid,
		Wnames: wnames,
	})
	if err != nil {
		return nil, nil, err
	}
	switch fc := fc.(type) {
	case *Rwalk:
		if len(fc.WQids) != len(wnames) {
			return nil, fc.WQids, ErrShortWalk
		}
		success = true
		return &ClientDotLFile{
			Client: f.Client,
			Fid:    fid,
		}, fc.WQids, nil
	case *Rlerror:
		return nil, nil, fc
	default:
		return nil, nil, fmt.Errorf("protocol error, expected Rattach")
	}
}

func (f *ClientDotLFile) Walk(wnames []string) (*ClientDotLFile, []Qid, error) {

	if len(wnames) == 0 {
		return f.walk(wnames)
	}

	wFile := f
	qids := []Qid{}

	for len(wnames) != 0 {
		batchSize := 13 // From spec.
		if len(wnames) < batchSize {
			batchSize = len(wnames)
		}
		batch := wnames[:batchSize]
		wnames = wnames[batchSize:]
		newWFile, newQids, err := wFile.walk(batch)
		if len(newQids) != 0 {
			qids = append(qids, newQids...)
		}
		if wFile != f {
			_ = wFile.Clunk()
		}
		if err != nil {
			return nil, qids, err
		}
		wFile = newWFile
	}

	return wFile, qids, nil
}

func (f *ClientDotLFile) Open(flags uint32) error {
	fc, err := f.Client.Fcall(&Tlopen{
		Fid:   f.Fid,
		Flags: flags,
	})
	if err != nil {
		return err
	}
	switch fc := fc.(type) {
	case *Rlopen:
		return nil
	case *Rlerror:
		return fc
	default:
		return errors.New("protocol error, expected Rlopen")
	}
}

func (f *ClientDotLFile) Fsync() error {
	fc, err := f.Client.Fcall(&Tfsync{
		Fid: f.Fid,
	})
	if err != nil {
		return err
	}
	switch fc := fc.(type) {
	case *Rfsync:
		return nil
	case *Rlerror:
		return fc
	default:
		return errors.New("protocol error, expected Rfsync")
	}
}

func (f *ClientDotLFile) Read(offset uint64, buf []byte) (uint32, error) {
	if uint32(len(buf)) > (f.Client.Msize() - IOHDRSZ) {
		buf = buf[:int(f.Client.Msize()-IOHDRSZ)]
	}
	fc, err := f.Client.Fcall(&Tread{
		Fid:    f.Fid,
		Offset: offset,
		Count:  uint32(len(buf)),
	})
	if err != nil {
		return 0, err
	}
	switch fc := fc.(type) {
	case *Rread:
		if len(fc.Data) > len(buf) {
			return 0, errors.New("returned data exceeds buffer")
		}
		buf = buf[:len(fc.Data)]
		copy(buf, fc.Data)
		return uint32(len(fc.Data)), nil
	case *Rlerror:
		return 0, fc
	default:
		return 0, errors.New("protocol error, expected Rread")
	}
}

func (f *ClientDotLFile) Write(offset uint64, buf []byte) (uint32, error) {
	if uint32(len(buf)) > (f.Client.Msize() - IOHDRSZ) {
		buf = buf[:int(f.Client.Msize()-IOHDRSZ)]
	}
	fc, err := f.Client.Fcall(&Twrite{
		Fid:    f.Fid,
		Offset: offset,
		Data:   buf,
	})
	if err != nil {
		return 0, err
	}
	switch fc := fc.(type) {
	case *Rwrite:
		return fc.Count, nil
	case *Rlerror:
		return 0, fc
	default:
		return 0, errors.New("protocol error, expected Rwrite")
	}
}

func (f *ClientDotLFile) Create(name string, flags uint32, mode uint32, gid uint32) (Qid, uint32, error) {
	fc, err := f.Client.Fcall(&Tlcreate{
		Fid:   f.Fid,
		Name:  name,
		Flags: flags,
		Mode:  mode,
		Gid:   gid,
	})
	if err != nil {
		return Qid{}, 0, err
	}
	switch fc := fc.(type) {
	case *Rlcreate:
		return fc.Qid, fc.Iounit, nil
	case *Rlerror:
		return Qid{}, 0, fc
	default:
		return Qid{}, 0, errors.New("protocol error, expected Rlcreate")
	}
}

func (f *ClientDotLFile) GetAttr(mask uint64) (LAttr, error) {
	fc, err := f.Client.Fcall(&Tgetattr{
		Fid:  f.Fid,
		Mask: mask,
	})
	if err != nil {
		return LAttr{}, err
	}
	switch fc := fc.(type) {
	case *Rgetattr:
		return fc.LAttr, nil
	case *Rlerror:
		return LAttr{}, fc
	default:
		return LAttr{}, errors.New("protocol error, expected Rgetattr")
	}
}

func (f *ClientDotLFile) SetAttr(attr LSetAttr) error {
	fc, err := f.Client.Fcall(&Tsetattr{
		Fid:      f.Fid,
		LSetAttr: attr,
	})
	if err != nil {
		return err
	}
	switch fc := fc.(type) {
	case *Rsetattr:
		return nil
	case *Rlerror:
		return fc
	default:
		return errors.New("protocol error, expected Rsetattr")
	}
}

func (f *ClientDotLFile) Rename(dir *ClientDotLFile, name string) error {
	fc, err := f.Client.Fcall(&Trename{
		Fid:  f.Fid,
		Dfid: dir.Fid,
		Name: name,
	})
	if err != nil {
		return err
	}
	switch fc := fc.(type) {
	case *Rrename:
		return nil
	case *Rlerror:
		return fc
	default:
		return errors.New("protocol error, expected Rrename")
	}
}

func (f *ClientDotLFile) Mkdir(name string, mode uint32, gid uint32) (Qid, error) {
	fc, err := f.Client.Fcall(&Tmkdir{
		Dfid: f.Fid,
		Name: name,
		Mode: mode,
		Gid:  gid,
	})
	if err != nil {
		return Qid{}, err
	}
	switch fc := fc.(type) {
	case *Rmkdir:
		return fc.Qid, nil
	case *Rlerror:
		return Qid{}, fc
	default:
		return Qid{}, errors.New("protocol error, expected Rmkdir")
	}
}

func (f *ClientDotLFile) Statfs() (LStatfs, error) {
	fc, err := f.Client.Fcall(&Tstatfs{
		Fid: f.Fid,
	})
	if err != nil {
		return LStatfs{}, err
	}
	switch fc := fc.(type) {
	case *Rstatfs:
		return fc.LStatfs, nil
	case *Rlerror:
		return LStatfs{}, fc
	default:
		return LStatfs{}, errors.New("protocol error, expected Rstatfs")
	}
}

type DotLDirIter struct {
	file   *ClientDotLFile
	ents   []DirEnt
	offset uint64
	done   bool
}

func (di *DotLDirIter) HasNext() bool {
	return !di.done
}

func (di *DotLDirIter) fill() error {
	ents, err := di.file.Readdir(di.offset, 0xFFFFFFFF)
	if err != nil {
		return err
	}

	// Reverse entries so we can pop them off in the right order.
	for i, j := 0, len(ents)-1; i < j; i, j = i+1, j-1 {
		ents[i], ents[j] = ents[j], ents[i]
	}

	di.ents = ents
	if len(di.ents) > 0 {
		di.offset = di.ents[0].Offset
	} else {
		di.done = true
	}
	return nil
}

func (di *DotLDirIter) Next() (DirEnt, error) {

	if di.done {
		return DirEnt{}, io.EOF
	}

	// Fill initial listing, otherwise we should always have something.
	if len(di.ents) == 0 {
		err := di.fill()
		if err != nil {
			return DirEnt{}, err
		}
	}

	if di.done {
		// Should never really happen considering . and ..
		return DirEnt{}, io.EOF
	}

	nextEnt := di.ents[len(di.ents)-1]
	di.ents = di.ents[:len(di.ents)-1]

	if len(di.ents) == 0 {
		err := di.fill()
		if err != nil {
			return nextEnt, err
		}
	}

	return nextEnt, nil
}

func (di *DotLDirIter) Unget(ent DirEnt) {
	di.ents = append(di.ents, ent)
	di.done = false
}

func (f *ClientDotLFile) DirIter() *DotLDirIter {
	return &DotLDirIter{
		file: f,
	}
}

func (f *ClientDotLFile) Readdir(offset uint64, count uint32) ([]DirEnt, error) {
	maxCount := f.Client.Msize() - READDIRHDRSZ
	if count > maxCount {
		count = maxCount
	}
	fc, err := f.Client.Fcall(&Treaddir{
		Fid:    f.Fid,
		Offset: offset,
		Count:  count,
	})
	if err != nil {
		return nil, err
	}
	switch fc := fc.(type) {
	case *Rreaddir:
		return fc.Data, nil
	case *Rlerror:
		return nil, fc
	default:
		return nil, errors.New("protocol error, expected Rreaddir")
	}
}

func (f *ClientDotLFile) ReaddirAll() ([]DirEnt, error) {
	allEnts := make([]DirEnt, 0, 8)
	offset := uint64(0)
	for {
		ents, err := f.Readdir(offset, 0xFFFFFFFF)
		if err != nil {
			return nil, err
		}
		if len(ents) == 0 {
			break
		}
		allEnts = append(allEnts, ents...)
		offset = allEnts[len(allEnts)-1].Offset
	}
	return allEnts, nil
}

func (f *ClientDotLFile) Lock(l LSetLock) (byte, error) {
	fc, err := f.Client.Fcall(&Tlock{
		Fid:      f.Fid,
		LSetLock: l,
	})
	if err != nil {
		return L_LOCK_ERROR, err
	}
	switch fc := fc.(type) {
	case *Rlock:
		return fc.Status, nil
	case *Rlerror:
		return L_LOCK_ERROR, fc
	default:
		return L_LOCK_ERROR, errors.New("protocol error, expected Rlock")
	}
}
