package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/andrewchambers/proto9-go"
	"github.com/hanwen/go-fuse/v2/fuse"
)

func ErrToStatus(err error) fuse.Status {
	if err == nil {
		return fuse.OK
	}
	switch err := err.(type) {
	case *proto9.Rlerror:
		// TODO: we need an adaptor for non linux platforms.
		return fuse.Status(err.Ecode)
	}
	return fuse.Status(syscall.EIO)
}

func FillFuseAttrFromAttr(attr *proto9.LAttr, out *fuse.Attr) {
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
	out.Generation = uint64(attr.Qid.Version)
	FillFuseAttrFromAttr(attr, &out.Attr)
}

type Inode9 struct {
	nodeId uint64
	refs   uint64
	qid    proto9.Qid
	f      *proto9.ClientDotLFile
}

func (i *Inode9) IncRef(n uint64) uint64 {
	return atomic.AddUint64(&i.refs, n)
}

func (i *Inode9) RefCount() uint64 {
	return atomic.LoadUint64(&i.refs)
}

func (i *Inode9) DecRef(n uint64) uint64 {
	return atomic.AddUint64(&i.refs, ^(n - 1))
}

type Proto9FS struct {
	fuse.RawFileSystem

	nodeIdCounter uint64

	lock sync.RWMutex

	n2i map[uint64]*Inode9
	p2i map[uint64]*Inode9
}

func NewProto9FS(rootFile *proto9.ClientDotLFile, rootQid proto9.Qid) *Proto9FS {
	fs := &Proto9FS{
		RawFileSystem: fuse.NewDefaultRawFileSystem(),
		nodeIdCounter: 1,
		n2i:           make(map[uint64]*Inode9),
		p2i:           make(map[uint64]*Inode9),
	}

	rootInode := &Inode9{
		nodeId: 1,
		refs:   1,
		qid:    rootQid,
		f:      rootFile,
	}

	fs.n2i[rootInode.nodeId] = rootInode
	fs.p2i[rootInode.qid.Path] = rootInode

	return fs
}

func (fs *Proto9FS) nextNodeId() uint64 {
	return atomic.AddUint64(&fs.nodeIdCounter, 1)
}

// A lookup request is sent by the kernel when the VFS wants to know
// about a child inode. Many lookups calls can occur in parallel,
// but only one call happens for each (inode, name) pair.
func (fs *Proto9FS) Lookup(cancel <-chan struct{}, header *fuse.InHeader, name string, out *fuse.EntryOut) fuse.Status {

	success := false

	fs.lock.Lock()
	parent := fs.n2i[header.NodeId]
	fs.lock.Unlock()

	f, qids, err := parent.f.Walk([]string{name})
	if err != nil {
		return ErrToStatus(err)
	}
	qid := qids[0]
	defer func() {
		if !success {
			f.Clunk()
		}
	}()

	attr, err := f.GetAttr(proto9.L_GETATTR_ALL)
	if err != nil {
		return ErrToStatus(err)
	}
	FillFuseEntryOutFromAttr(&attr, out)

	fs.lock.Lock()
	inode, ok := fs.p2i[qid.Path]
	if !ok || true {
		inode = &Inode9{
			nodeId: fs.nextNodeId(),
			qid:    qid,
			f:      f,
			refs:   1,
		}
		fs.n2i[inode.nodeId] = inode
		fs.p2i[inode.qid.Path] = inode
	} else {
		inode.IncRef(1)
	}
	fs.lock.Unlock()

	out.NodeId = inode.nodeId

	log.Printf("XXX lookup %d rc=%d", out.NodeId, inode.RefCount())

	success = true
	return fuse.OK
}

// A forget request is sent by the kernel when it is no
// longer interested in an inode.
func (fs *Proto9FS) Forget(nodeId, nlookup uint64) {

	fs.lock.Lock()
	inode := fs.n2i[nodeId]

	if inode == nil {
		// XXX happens due to go-fuse epoll hack.
		fs.lock.Unlock()
		return
	}

	rc := inode.DecRef(nlookup)
	if rc == 0 {
		delete(fs.n2i, nodeId)
		pi := fs.p2i[inode.qid.Path]
		if inode == pi {
			delete(fs.p2i, inode.qid.Path)
		}
	}
	fs.lock.Unlock()

	log.Printf("XXX: forget nodeId=%d, rc=%d", nodeId, rc)

	if rc == 0 {
		inode.f.Clunk()
	}
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

	rootFile, rootQid, err := proto9.AttachDotL(client, *aname, *uname)
	if err != nil {
		log.Fatalf("unable to attach to mount: %s", err)
	}
	defer rootFile.Clunk()

	fs := NewProto9FS(rootFile, rootQid)

	server, err := fuse.NewServer(fs, mntDir, &fuse.MountOptions{
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
		Debug:                true,
	})
	if err != nil {
		log.Fatalf("unable to create fuse server: %s", err)
	}

	go server.Serve()

	err = server.WaitMount()
	if err != nil {
		log.Fatalf("unable wait for mount: %s", err)
	}

	// Serve the file system, until unmounted by calling fusermount -u
	server.Wait()
}
