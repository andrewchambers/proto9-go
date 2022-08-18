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

func FillEntryOutFromAttr(attr *proto9.LAttr, out *fuse.EntryOut) {
	out.Ino = attr.Qid.Path
	out.Generation = uint64(attr.Qid.Version)
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



type Inode9 struct {
	fs.Inode
	file *proto9.ClientDotLFile
}

type FileHandle9 struct {
	file *proto9.ClientDotLFile
}

var _ = (fs.FileReader)((*FileHandle9)(nil))

func (fh *FileHandle9) Read(ctx context.Context, dest []byte, off int64) (fuse.ReadResult, syscall.Errno) {
	n, err := fh.file.Read(uint64(off), dest)
	if err != nil {
		return nil, ErrToErrno(err)
	}
	return fuse.ReadResultData(dest[:n]), 0
}

var _ = (fs.FileWriter)((*FileHandle9)(nil))

func (fh *FileHandle9) Write(ctx context.Context, dest []byte, off int64) (uint32, syscall.Errno) {
	n, err := fh.file.Write(uint64(off), dest)
	if err != nil {
		return n, ErrToErrno(err)
	}
	return n, 0
}

var _ = (fs.FileReleaser)((*FileHandle9)(nil))

func (fh *FileHandle9) Release(ctx context.Context) syscall.Errno {
	_ = fh.file.Clunk()
	return 0
}


var _ = (fs.NodeLookuper)((*Inode9)(nil))

func (n *Inode9) Lookup(ctx context.Context, name string, out *fuse.EntryOut) (*fs.Inode, syscall.Errno) {
	newf, qids, err := n.file.Walk([]string{name})
	if err != nil {
		return nil, ErrToErrno(err)
	}
	stableAttr := StableAttrFromQid(&qids[0])
	// XXX can we get away without this getattr? it seems like fuse might not
	// strictly need it, especially since we don't do caching.
	attr, err := newf.GetAttr(proto9.L_GETATTR_ALL)
	if err != nil {
		return nil, ErrToErrno(err)
	}
	FillEntryOutFromAttr(&attr, out)
	newInode := n.NewInode(ctx, &Inode9{file: newf}, stableAttr)
	return newInode, 0
}

var _ = (fs.NodeOpener)((*Inode9)(nil))

func (n *Inode9) Open(ctx context.Context, flags uint32) (fs.FileHandle, uint32, syscall.Errno) {
	
	var flags9 uint32
	if flags&syscall.O_RDWR ==  syscall.O_RDWR {
		flags9 |= proto9.L_O_RDWR
	} else if flags&syscall.O_WRONLY ==  syscall.O_WRONLY {
		flags9 |= proto9.L_O_WRONLY
	} else if flags&syscall.O_RDONLY ==  syscall.O_RDONLY {
		flags9 |= proto9.L_O_RDONLY
	}

	newf, _, err := n.file.Walk([]string{})
	if err != nil {
		return nil, 0, ErrToErrno(err)
	}
	err = newf.Open(flags9)
	if err != nil {
		_ = newf.Clunk()
		return nil, 0, ErrToErrno(err)
	}
	fh := &FileHandle9{
		file: newf,
	}
	return fh, fuse.FOPEN_DIRECT_IO, 0
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

	rootInode := &Inode9{
		file: attachPoint,
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
				AllowOther:    false, // XXX option?
				DisableXAttrs: false, // TODO implement
				EnableLocks:   false, // TODO implement
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
