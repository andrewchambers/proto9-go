package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/andrewchambers/proto9-go"
	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

func ErrToErrno(err error) syscall.Errno {
	if err == nil {
		return 0
	}
	switch err := err.(type) {
	case *proto9.Rlerror:
		// TODO: we need an adaptor for non linux platforms.
		return syscall.Errno(err.Ecode)
	}
	return syscall.EIO
}

func StableAttrFromQid(q *proto9.Qid) fs.StableAttr {
	var mode uint32

	// TODO more types.
	if (q.Typ & proto9.QT_DIR) != 0 {
		mode = syscall.S_IFDIR
	} else {
		mode = syscall.S_IFREG
	}

	return fs.StableAttr{
		Mode: mode,
		Ino:  q.Path,
		Gen:  uint64(q.Version),
	}
}

func FillFuseAttrOutFromAttr(attr *proto9.LAttr, out *fuse.Attr) {
	out.Ino = attr.Qid.Path
	out.Size = attr.Size
	out.Blocks = attr.Blocks
	out.Blksize = uint32(attr.Blksize)
	out.Atime = attr.AtimeSec
	out.Atimensec = uint32(attr.AtimeNsec)
	out.Mtime = attr.MtimeSec
	out.Mtimensec = uint32(attr.MtimeNsec)
	out.Ctime = attr.CtimeSec
	out.Ctimensec = uint32(attr.CtimeNsec)
	out.Mode = attr.Mode
	out.Nlink = uint32(attr.Nlink)
	out.Owner.Uid = attr.Uid
	out.Owner.Gid = attr.Gid
	out.Rdev = uint32(attr.Rdev)
}

func FillFuseEntryOutFromAttr(attr *proto9.LAttr, out *fuse.EntryOut) {
	out.Ino = attr.Qid.Path
	out.Generation = uint64(attr.Qid.Version)
	FillFuseAttrOutFromAttr(attr, &out.Attr)
}

func DirEntToFuseDirent(de *proto9.DirEnt) fuse.DirEntry {
	return fuse.DirEntry{
		Mode: uint32(de.Typ),
		Ino:  de.Qid.Path,
		Name: de.Name,
	}
}

type DirHandle9 struct {
	file   *proto9.ClientDotLFile
	ents   []proto9.DirEnt
	offset uint64
	done   bool
}

func (dh *DirHandle9) HasNext() bool {
	return !dh.done
}

func (dh *DirHandle9) fill() error {
	ents, err := dh.file.Readdir(dh.offset, 0xFFFFFFFF)
	if err != nil {
		return err
	}
	dh.ents = append(dh.ents, ents...)
	if len(dh.ents) > 0 {
		dh.offset = dh.ents[len(dh.ents)-1].Offset
	} else {
		dh.done = true
	}
	return nil
}

func (dh *DirHandle9) Next() (fuse.DirEntry, syscall.Errno) {

	if dh.done {
		return fuse.DirEntry{}, syscall.EIO
	}

	// Fill initial listing, otherwise we should always have something.
	if len(dh.ents) == 0 {
		err := dh.fill()
		if err != nil {
			return fuse.DirEntry{}, ErrToErrno(err)
		}
		if dh.done {
			// Should never really happen considering . and ..
			return fuse.DirEntry{}, syscall.EIO
		}
	}

	nextEnt := dh.ents[0]
	dh.ents = dh.ents[1:]

	if len(dh.ents) == 0 {
		err := dh.fill()
		if err != nil {
			return fuse.DirEntry{}, ErrToErrno(err)
		}
	}

	return DirEntToFuseDirent(&nextEnt), 0
}

func (dh *DirHandle9) Close() {
	_ = dh.file.Clunk()
}

type FileHandle9 struct {
	file *proto9.ClientDotLFile
}

var _ = (fs.FileReader)((*FileHandle9)(nil))
var _ = (fs.FileWriter)((*FileHandle9)(nil))
var _ = (fs.FileReleaser)((*FileHandle9)(nil))
var _ = (fs.FileFsyncer)((*FileHandle9)(nil))
var _ = (fs.FileSetlker)((*FileHandle9)(nil))
var _ = (fs.FileSetlkwer)((*FileHandle9)(nil))

func (fh *FileHandle9) setlck(ctx context.Context, owner uint64, lk *fuse.FileLock, flags uint32, wait bool) syscall.Errno {

	/* XXX
	if flags != 0 {
		return syscall.ENOTSUP
	}
	*/

	typ9 := uint8(0)

	switch lk.Typ {
	case syscall.F_RDLCK:
		typ9 = proto9.L_LOCK_TYPE_RDLCK
	case syscall.F_WRLCK:
		typ9 = proto9.L_LOCK_TYPE_WRLCK
	case syscall.F_UNLCK:
		typ9 = proto9.L_LOCK_TYPE_UNLCK
	default:
		return syscall.ENOTSUP
	}

	flags9 := uint32(0)
	if wait {
		flags9 |= proto9.L_LOCK_FLAGS_BLOCK
	}

	for {
		status, err := fh.file.Lock(proto9.LSetLock{
			Typ:    typ9,
			Flags:  flags9,
			Start:  lk.Start,
			Length: lk.End - lk.Start,
			ProcId: lk.Pid,
		})
		if err != nil {
			return ErrToErrno(err)
		}

		switch status {
		case proto9.L_LOCK_SUCCESS:
			return 0
		case proto9.L_LOCK_BLOCKED:
			if wait {
				// Server doesn't seem to support blocking.
				time.Sleep(1 * time.Second)
				continue
			}
			return syscall.EAGAIN
		default:
			return syscall.EIO
		}
	}
}

func (fh *FileHandle9) Setlk(ctx context.Context, owner uint64, lk *fuse.FileLock, flags uint32) syscall.Errno {
	return fh.setlck(ctx, owner, lk, flags, false)
}

func (fh *FileHandle9) Setlkw(ctx context.Context, owner uint64, lk *fuse.FileLock, flags uint32) syscall.Errno {
	return fh.setlck(ctx, owner, lk, flags, true)
}

func (fh *FileHandle9) Fsync(ctx context.Context, flags uint32) syscall.Errno {
	err := fh.file.Fsync()
	if err != nil {
		return ErrToErrno(err)
	}
	return 0
}

func (fh *FileHandle9) Read(ctx context.Context, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	n, err := fh.file.Read(uint64(off), dest)
	if err != nil {
		return nil, ErrToErrno(err)
	}
	return fuse.ReadResultData(dest[:n]), 0
}

func (fh *FileHandle9) Write(ctx context.Context, dest []byte, off int64) (uint32, syscall.Errno) {
	n, err := fh.file.Write(uint64(off), dest)
	if err != nil {
		return n, ErrToErrno(err)
	}
	return n, 0
}

func (fh *FileHandle9) Release(ctx context.Context) syscall.Errno {
	_ = fh.file.Clunk()
	return 0
}

type Inode9 struct {
	fs.Inode
	root *proto9.ClientDotLFile
}

var _ = (fs.NodeGetattrer)((*Inode9)(nil))
var _ = (fs.NodeSetattrer)((*Inode9)(nil))
var _ = (fs.NodeLookuper)((*Inode9)(nil))
var _ = (fs.NodeCreater)((*Inode9)(nil))
var _ = (fs.NodeOpener)((*Inode9)(nil))
var _ = (fs.NodeReaddirer)((*Inode9)(nil))
var _ = (fs.NodeUnlinker)((*Inode9)(nil))
var _ = (fs.NodeMkdirer)((*Inode9)(nil))
var _ = (fs.NodeRmdirer)((*Inode9)(nil))
var _ = (fs.NodeRenamer)((*Inode9)(nil))

func (n *Inode9) pathToRoot() ([]string, bool) {
	if n.IsRoot() {
		return []string{}, true
	}
	path := make([]string, 0, 8)
	curNode := n.EmbeddedInode()
	for {
		name, nextNode := curNode.Parent()
		if nextNode == nil {
			return nil, false
		}
		path = append(path, name)
		if nextNode.IsRoot() {
			break
		}
		curNode = nextNode
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path, true
}

func (n *Inode9) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	var f *proto9.ClientDotLFile

	if fh != nil {
		f = fh.(*FileHandle9).file
	} else {
		path, found := n.pathToRoot()
		if !found {
			return syscall.EIO
		}
		newf, _, err := n.root.Walk(path)
		if err != nil {
			return ErrToErrno(err)
		}
		defer newf.Clunk()
		f = newf
	}
	attr, err := f.GetAttr(proto9.L_GETATTR_ALL)
	if err != nil {
		return ErrToErrno(err)
	}
	FillFuseAttrOutFromAttr(&attr, &out.Attr)
	// XXX should the inode come from the qid?
	//if out.Ino != n.StableAttr.Ino() {
	// XXX This is a possibility due to our network fs, do we care?
	//}
	return 0

}

func (n *Inode9) Setattr(ctx context.Context, fh fs.FileHandle, in *fuse.SetAttrIn, out *fuse.AttrOut) syscall.Errno {

	var f *proto9.ClientDotLFile

	if fh != nil {
		f = fh.(*FileHandle9).file
	} else {
		path, found := n.pathToRoot()
		if !found {
			return syscall.EIO
		}
		newf, _, err := n.root.Walk(path)
		if err != nil {
			return ErrToErrno(err)
		}
		defer newf.Clunk()
		f = newf
	}

	setAttr := proto9.LSetAttr{}

	if mtime, ok := in.GetMTime(); ok {
		setAttr.MtimeSec = uint64(mtime.Unix())
		setAttr.MtimeNsec = uint64(mtime.UnixNano() - (mtime.Unix() * 1000_000_000))
		setAttr.Valid |= proto9.L_SETATTR_MTIME
	}
	if atime, ok := in.GetATime(); ok {
		setAttr.AtimeSec = uint64(atime.Unix())
		setAttr.AtimeNsec = uint64(atime.UnixNano() - (atime.Unix() * 1000_000_000))
		setAttr.Valid |= proto9.L_SETATTR_ATIME
	}
	if size, ok := in.GetSize(); ok {
		setAttr.Size = size
		setAttr.Valid |= proto9.L_SETATTR_SIZE
	}
	if mode, ok := in.GetMode(); ok {
		setAttr.Mode = mode
		setAttr.Valid |= proto9.L_SETATTR_MODE
	}

	// TODO
	// in.GetCTime()
	// in.GetGID()
	// in.GetUID()

	err := f.SetAttr(setAttr)
	if err != nil {
		return ErrToErrno(err)
	}

	// XXX Not sure if we need to run GetAttr, we don't cache attributes anyway.
	/*
		attr, err := fh.file.GetAttr(proto9.L_GETATTR_ALL)
		if err != nil {
			return ErrToErrno(err)
		}
		FillFuseAttrOutFromAttr(&attr, &out.Attr)
	*/

	return 0
}

func (n *Inode9) Rename(ctx context.Context, name string, newParent fs.InodeEmbedder, newName string, flags uint32) syscall.Errno {

	// XXX mutex?
	// XXX handle flags
	//if flags != 0 {
	//	return syscall.ENOTSUP
	//}

	oldChildPath, found := n.pathToRoot()
	if !found {
		return syscall.EIO
	}
	oldChildPath = append(oldChildPath, name)

	f, _, err := n.root.Walk(oldChildPath)
	if err != nil {
		return ErrToErrno(err)
	}
	defer f.Clunk()

	newParent9 := newParent.(*Inode9)
	newChildPath, found := newParent9.pathToRoot()
	if !found {
		return syscall.EIO
	}
	newChildPath = append(newChildPath, newName)
	newParentPath := newChildPath[:len(newChildPath)-1]

	parentf, _, err := newParent9.root.Walk(newParentPath)
	if err != nil {
		return ErrToErrno(err)
	}
	defer parentf.Clunk()

	err = f.Rename(parentf, newName)
	if err != nil {
		return ErrToErrno(err)
	}

	// TODO update child path somehow?

	return 0
}

func (n *Inode9) Unlink(ctx context.Context, name string) syscall.Errno {
	path, found := n.pathToRoot()
	if !found {
		return syscall.EIO
	}
	path = append(path, name)
	f, _, err := n.root.Walk(path)
	if err != nil {
		return ErrToErrno(err)
	}
	return ErrToErrno(f.Remove())
}

func (n *Inode9) Rmdir(ctx context.Context, name string) syscall.Errno {
	return n.Unlink(ctx, name)
}

func (n *Inode9) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	path, found := n.pathToRoot()
	if !found {
		return nil, syscall.EIO
	}
	path = append(path, name)
	f, qids, err := n.root.Walk(path)
	if err != nil {
		return nil, ErrToErrno(err)
	}
	defer f.Clunk()

	// XXX can we get away without this getattr? it seems like fuse might not
	// strictly need it, especially since we don't do caching.
	attr, err := f.GetAttr(proto9.L_GETATTR_ALL)
	if err != nil {
		return nil, ErrToErrno(err)
	}
	FillFuseEntryOutFromAttr(&attr, out)

	newInode := n.NewInode(ctx, &Inode9{root: n.root}, StableAttrFromQid(&qids[0]))
	return newInode, 0
}

func (n *Inode9) Mkdir(ctx context.Context, name string, mode uint32, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	path, found := n.pathToRoot()
	if !found {
		return nil, syscall.EIO
	}
	path = append(path, name)
	// so we only need a single path load.
	parentPath := path[:len(path)-1]

	f, _, err := n.root.Walk(parentPath)
	if err != nil {
		return nil, ErrToErrno(err)
	}
	defer f.Clunk()

	var mode9 uint32

	mode9 = mode     // XXX convert flags?
	gid := uint32(0) // XXX correct value?

	qid, err := f.Mkdir(name, mode9, gid)
	if err != nil {
		return nil, ErrToErrno(err)
	}

	child, _, err := f.Walk([]string{name})
	if err != nil {
		return nil, ErrToErrno(err)
	}
	defer child.Clunk()

	attr, err := child.GetAttr(proto9.L_GETATTR_ALL)
	if err != nil {
		return nil, ErrToErrno(err)
	}
	FillFuseEntryOutFromAttr(&attr, out)

	newInode := n.NewInode(ctx, &Inode9{root: n.root}, StableAttrFromQid(&qid))
	return newInode, 0
}

func (n *Inode9) Create(ctx context.Context, name string, flags uint32, mode uint32, out *fuse.EntryOut) (*fs.Inode, fs.FileHandle, uint32, syscall.Errno) {
	success := false
	path, found := n.pathToRoot()
	if !found {
		return nil, nil, 0, syscall.EIO
	}
	path = append(path, name)
	// so we only need a single path load.
	parentPath := path[:len(path)-1]

	f, _, err := n.root.Walk(parentPath)
	if err != nil {
		return nil, nil, 0, ErrToErrno(err)
	}
	defer func() {
		if !success {
			_ = f.Clunk()
		}
	}()

	var flags9 uint32
	var mode9 uint32

	flags9 = flags   // XXX convert flags?
	mode9 = mode     // XXX convert flags?
	gid := uint32(0) // XXX correct value?

	qid, _, err := f.Create(name, flags9, mode9, gid)
	if err != nil {
		return nil, nil, 0, ErrToErrno(err)
	}

	attr, err := f.GetAttr(proto9.L_GETATTR_ALL)
	if err != nil {
		return nil, nil, 0, ErrToErrno(err)
	}
	FillFuseEntryOutFromAttr(&attr, out)

	newInode := n.NewInode(ctx, &Inode9{root: n.root}, StableAttrFromQid(&qid))
	success = true
	return newInode, &FileHandle9{file: f}, fuse.FOPEN_DIRECT_IO, 0
}

func (n *Inode9) Open(ctx context.Context, flags uint32) (fs.FileHandle, uint32, syscall.Errno) {

	var flags9 uint32

	if flags&syscall.O_RDWR == syscall.O_RDWR {
		flags9 |= proto9.L_O_RDWR
	} else if flags&syscall.O_WRONLY == syscall.O_WRONLY {
		flags9 |= proto9.L_O_WRONLY
	} else if flags&syscall.O_RDONLY == syscall.O_RDONLY {
		flags9 |= proto9.L_O_RDONLY
	}

	if flags&syscall.O_TRUNC == syscall.O_TRUNC {
		flags9 |= proto9.L_O_TRUNC
	}

	path, found := n.pathToRoot()
	if !found {
		return nil, 0, syscall.EIO
	}
	newf, _, err := n.root.Walk(path)
	if err != nil {
		return nil, 0, ErrToErrno(err)
	}
	success := false
	defer func() {
		if !success {
			_ = newf.Clunk()
		}
	}()
	err = newf.Open(flags9)
	if err != nil {
		return nil, 0, ErrToErrno(err)
	}
	success = true
	return &FileHandle9{file: newf}, fuse.FOPEN_DIRECT_IO, 0
}

func (n *Inode9) Readdir(ctx context.Context) (fs.DirStream, syscall.Errno) {
	path, found := n.pathToRoot()
	if !found {
		return nil, syscall.EIO
	}
	newf, _, err := n.root.Walk(path)
	if err != nil {
		return nil, ErrToErrno(err)
	}
	success := false
	defer func() {
		if !success {
			_ = newf.Clunk()
		}
	}()
	err = newf.Open(proto9.L_O_RDONLY)
	if err != nil {
		return nil, ErrToErrno(err)
	}
	dh := &DirHandle9{
		file: newf,
	}
	success = true
	return dh, 0
}

func usage() {
	fmt.Printf("9pfuse [OPTS] MOUNTPOINT\n")
	os.Exit(1)
}

func main() {

	address := flag.String("address", "localhost:1777", "address to connect to.")
	msize := flag.Uint("msize", 65536, "maximum message size.")
	aname := flag.String("aname", "", "aname to send in the attach message.")
	uname := flag.String("uname", "", "uname to send in the attach message.")

	flag.Parse()

	if len(flag.Args()) != 1 {
		usage()
	}

	mntDir := flag.Args()[0]

	conn, err := net.Dial("tcp", *address)
	if err != nil {
		log.Fatalf("unable to dial address: %s", err)
	}

	client, err := proto9.NewClient(conn, "9P2000.L", uint32(*msize))
	if err != nil {
		log.Fatalf("unable to negotiate protocol version: %s", err)
	}
	defer client.Close()

	attachPoint, err := proto9.AttachDotL(client, *aname, *uname)
	if err != nil {
		log.Fatalf("unable to attach to mount: %s", err)
	}
	defer attachPoint.Clunk()

	rootInode := &Inode9{root: attachPoint}

	server, err := fs.Mount(mntDir, rootInode,
		&fs.Options{
			MountOptions: fuse.MountOptions{
				Debug:   true,
				Name:    "9p",
				Options: []string{
					// XXX why are these not working?
					// "direct_io",
					// "hard_remove",
				},
				AllowOther:           false, // XXX option?
				DisableXAttrs:        false, // TODO implement
				EnableLocks:          true,
				IgnoreSecurityLabels: true, // option?
			},
			Logger: log.Default(),
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	// Serve the file system, until unmounted by calling fusermount -u
	server.Wait()
}
