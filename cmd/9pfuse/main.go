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
	path []string
}

var _ = (fs.NodeGetattrer)((*Inode9)(nil))
var _ = (fs.NodeLookuper)((*Inode9)(nil))
var _ = (fs.NodeCreater)((*Inode9)(nil))
var _ = (fs.NodeOpener)((*Inode9)(nil))
var _ = (fs.NodeReaddirer)((*Inode9)(nil))

func (n *Inode9) Getattr(ctx context.Context, fh fs.FileHandle, out *fuse.AttrOut) syscall.Errno {
	if fh != nil {
		fh9 := fh.(*FileHandle9)
		attr, err := fh9.file.GetAttr(proto9.L_GETATTR_ALL)
		if err != nil {
			return ErrToErrno(err)
		}
		FillFuseAttrOutFromAttr(&attr, &out.Attr)
		return 0
	} else {
		f, _, err := n.root.Walk(n.path)
		if err != nil {
			return ErrToErrno(err)
		}
		defer f.Clunk()

		attr, err := f.GetAttr(proto9.L_GETATTR_ALL)
		if err != nil {
			return ErrToErrno(err)
		}
		FillFuseAttrOutFromAttr(&attr, &out.Attr)
		return 0
	}

}

func (n *Inode9) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	path := make([]string, 0, len(n.path)+1)
	path = append([]string{}, n.path...)
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

	newInode := n.NewInode(ctx, &Inode9{root: n.root, path: path}, StableAttrFromQid(&qids[0]))

	return newInode, 0
}

func (n *Inode9) Create(ctx context.Context, name string, flags uint32, mode uint32, out *fuse.EntryOut) (*fs.Inode, fs.FileHandle, uint32, syscall.Errno) {

	path := make([]string, 0, len(n.path)+1)
	path = append([]string{}, n.path...)
	path = append(path, name)

	f, _, err := n.root.Walk(path)
	if err != nil {
		return nil, nil, 0, ErrToErrno(err)
	}

	success := false
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
	newInode := n.NewInode(ctx, &Inode9{root: n.root, path: path}, StableAttrFromQid(&qid))
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

	newf, _, err := n.root.Walk(n.path)
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
	newf, _, err := n.root.Walk(n.path)
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

	rootInode := &Inode9{
		root: attachPoint,
		path: []string{},
	}

	zeroSeconds := time.Duration(0)

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
				EnableLocks:          false, // TODO implement
				IgnoreSecurityLabels: true,  // option?
			},
			EntryTimeout:    &zeroSeconds,
			AttrTimeout:     &zeroSeconds,
			NegativeTimeout: &zeroSeconds,
			Logger:          log.Default(),
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	// Serve the file system, until unmounted by calling fusermount -u
	server.Wait()
}
