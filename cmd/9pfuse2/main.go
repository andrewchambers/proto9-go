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
	"time"

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

func qidToMode(q *proto9.Qid) uint32 {
	if q.Typ&proto9.QT_DIR != 0 {
		return syscall.S_IFDIR
	} else if q.Typ&proto9.QT_SYMLINK != 0 {
		return syscall.S_IFLNK
	} else {
		return syscall.S_IFREG
	}
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
	_f     atomic.Value
}

func (i *Inode9) GetFile() (*proto9.ClientDotLFile, bool) {
	v := i._f.Load()
	if v == nil {
		return nil, false
	}
	f := v.(*proto9.ClientDotLFile)
	if f == nil {
		return nil, false
	}
	return f, true
}

func (i *Inode9) SwapFile(f *proto9.ClientDotLFile) *proto9.ClientDotLFile {
	v := i._f.Swap(f)
	if v == nil {
		return nil
	}
	return v.(*proto9.ClientDotLFile)
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

type OpenFile struct {
	inode  *Inode9
	f      *proto9.ClientDotLFile
	diLock sync.Mutex
	di     *proto9.DotLDirIter
}

type Proto9FS struct {
	fuse.RawFileSystem

	server *fuse.Server

	nodeIdCounter uint64

	fileHandleCounter uint64

	lock sync.Mutex

	n2Inode     map[uint64]*Inode9
	p2Inode     map[uint64]*Inode9
	fh2OpenFile map[uint64]*OpenFile

	// TODO
	// dentsLock sync.Mutex
	// dents map[uint64]map[string]struct{}
}

func NewProto9FS(rootFile *proto9.ClientDotLFile, rootQid proto9.Qid) *Proto9FS {
	fs := &Proto9FS{
		RawFileSystem:     fuse.NewDefaultRawFileSystem(),
		nodeIdCounter:     1,
		fileHandleCounter: 0,
		n2Inode:           make(map[uint64]*Inode9),
		p2Inode:           make(map[uint64]*Inode9),
		fh2OpenFile:       make(map[uint64]*OpenFile),
		dirEnts:           make(map[uint64]map[string]struct{}),
	}

	rootInode := &Inode9{
		nodeId: 1,
		refs:   1,
		qid:    rootQid,
	}
	rootInode.SwapFile(rootFile)
	fs.n2Inode[rootInode.nodeId] = rootInode
	fs.p2Inode[rootInode.qid.Path] = rootInode
	return fs
}

func (fs *Proto9FS) nextNodeId() uint64 {
	return atomic.AddUint64(&fs.nodeIdCounter, 1)
}

func (fs *Proto9FS) nextFileHandle() uint64 {
	return atomic.AddUint64(&fs.fileHandleCounter, 1)
}

func (fs *Proto9FS) Init(server *fuse.Server) {
	fs.server = server
}

// A lookup request is sent by the kernel when the VFS wants to know
// about a child inode. Many lookups calls can occur in parallel,
// but only one call happens for each (inode, name) pair.
func (fs *Proto9FS) Lookup(cancel <-chan struct{}, header *fuse.InHeader, name string, out *fuse.EntryOut) fuse.Status {
	fs.lock.Lock()
	parent := fs.n2Inode[header.NodeId]
	fs.lock.Unlock()

	parentf, ok := parent.GetFile()
	if !ok {
		return fuse.EIO
	}

	f, qids, err := parentf.Walk([]string{name})
	if err != nil {
		return ErrToStatus(err)
	}
	qid := qids[0]
	defer func() {
		if f != nil {
			_ = f.Clunk()
		}
	}()

	attr, err := f.GetAttr(proto9.L_GETATTR_ALL)
	if err != nil {
		return ErrToStatus(err)
	}
	FillFuseEntryOutFromAttr(&attr, out)

	fs.lock.Lock()

	inode, ok := fs.p2Inode[qid.Path]
	if !ok {
		inode = &Inode9{
			nodeId: fs.nextNodeId(),
			qid:    qid,
			refs:   1,
		}
		f = inode.SwapFile(f)
		fs.n2Inode[inode.nodeId] = inode
		fs.p2Inode[inode.qid.Path] = inode
	} else {
		inode.IncRef(1)
		// Using the new file works around the fact
		// that diod doesn't update fids properly on rename.
		f = inode.SwapFile(f)
	}
	fs.lock.Unlock()

	out.NodeId = inode.nodeId

	/*
		// TODO
		// The idea is to maintain a list of all the dents we have
		// created so we can manually issue forget messages to free
		// kernel memory.
		fs.dentsLock.Lock()
		names, ok := fs.dents[parent.nodeId]
		if !ok {
			names = make(map[string]time.Time)
			fs.dents[parent.nodeId] = names
		}
		names[name] = struct{}{}
		fs.dentsLock.Unlock()
	*/

	// log.Printf("XXX lookup %d rc=%d", out.NodeId, inode.RefCount())
	return fuse.OK
}

// A forget request is sent by the kernel when it is no
// longer interested in an inode.
func (fs *Proto9FS) Forget(nodeId, nlookup uint64) {

	// log.Printf("XXX forget %d nlookup=%d", nodeId, nlookup)

	if nodeId == ^uint64(0) {
		// go-fuse uses this inode for its own purposes (epoll bug fix).
		return
	}

	fs.lock.Lock()
	inode := fs.n2Inode[nodeId]
	rc := inode.DecRef(nlookup)
	if rc == 0 {
		delete(fs.n2Inode, nodeId)
		pi := fs.p2Inode[inode.qid.Path]
		if inode == pi {
			delete(fs.p2Inode, inode.qid.Path)
		}
	}
	fs.lock.Unlock()

	if rc == 0 {
		f, ok := inode.GetFile()
		if ok {
			_ = f.Clunk()
		}
	}
}

func (fs *Proto9FS) GetAttr(cancel <-chan struct{}, in *fuse.GetAttrIn, out *fuse.AttrOut) fuse.Status {
	fs.lock.Lock()
	inode := fs.n2Inode[in.NodeId]
	fs.lock.Unlock()
	f, ok := inode.GetFile()
	if !ok {
		return fuse.EIO
	}
	attr, err := f.GetAttr(proto9.L_GETATTR_ALL)
	if err != nil {
		return ErrToStatus(err)
	}
	FillFuseAttrFromAttr(&attr, &out.Attr)
	return fuse.OK
}

func (fs *Proto9FS) SetAttr(cancel <-chan struct{}, in *fuse.SetAttrIn, out *fuse.AttrOut) fuse.Status {

	fs.lock.Lock()
	inode := fs.n2Inode[in.NodeId]
	fs.lock.Unlock()

	setAttr := proto9.LSetAttr{}

	if mtime, ok := in.GetMTime(); ok {
		setAttr.MtimeSec = uint64(mtime.Unix())
		setAttr.MtimeNsec = uint64(mtime.UnixNano() % 1000_000_000)
		setAttr.Valid |= proto9.L_SETATTR_MTIME
	}
	if atime, ok := in.GetATime(); ok {
		setAttr.AtimeSec = uint64(atime.Unix())
		setAttr.AtimeNsec = uint64(atime.UnixNano() % 1000_000_000)
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

	f, ok := inode.GetFile()
	if !ok {
		return fuse.EIO
	}

	err := f.SetAttr(setAttr)
	if err != nil {
		return ErrToStatus(err)
	}

	// XXX a full getattr might not be necessary.
	attr, err := f.GetAttr(proto9.L_GETATTR_ALL)
	if err != nil {
		return ErrToStatus(err)
	}
	FillFuseAttrFromAttr(&attr, &out.Attr)

	return fuse.OK
}

func openFlagsTo9(flags uint32) (uint32, bool) {
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

	if flags&syscall.O_EXCL == syscall.O_TRUNC {
		flags9 |= proto9.L_O_EXCL
	}

	// XXX more flags or errors for unsupported

	return flags, true
}

func (fs *Proto9FS) Open(cancel <-chan struct{}, in *fuse.OpenIn, out *fuse.OpenOut) fuse.Status {

	fs.lock.Lock()
	inode := fs.n2Inode[in.NodeId]
	fs.lock.Unlock()

	f, ok := inode.GetFile()
	if !ok {
		return fuse.EIO
	}

	f, _, err := f.Walk([]string{})
	if err != nil {
		return ErrToStatus(err)
	}

	defer func() {
		if f != nil {
			_ = f.Clunk()
		}
	}()

	flags9, ok := openFlagsTo9(in.Flags)
	if !ok {
		return fuse.ENOTSUP
	}

	err = f.Open(flags9)
	if err != nil {
		return ErrToStatus(err)
	}

	out.Fh = fs.nextFileHandle()
	out.OpenFlags |= fuse.FOPEN_DIRECT_IO

	fs.lock.Lock()
	fs.fh2OpenFile[out.Fh] = &OpenFile{
		inode: inode,
		f:     f,
	}
	f = nil
	fs.lock.Unlock()

	return fuse.OK
}

func (fs *Proto9FS) Create(cancel <-chan struct{}, in *fuse.CreateIn, name string, out *fuse.CreateOut) fuse.Status {

	fs.lock.Lock()
	inode := fs.n2Inode[in.NodeId]
	fs.lock.Unlock()

	f, ok := inode.GetFile()
	if !ok {
		return fuse.EIO
	}

	f, _, err := f.Walk([]string{})
	if err != nil {
		return ErrToStatus(err)
	}
	defer func() {
		if f != nil {
			_ = f.Clunk()
		}
	}()

	flags9, ok := openFlagsTo9(in.Flags)
	if !ok {
		return fuse.ENOTSUP
	}

	qid, _, err := f.Create(name, flags9, in.Mode, in.Caller.Gid)
	if err != nil {
		return ErrToStatus(err)
	}

	inodef, _, err := f.Walk([]string{})
	if err != nil {
		return ErrToStatus(err)
	}

	newInode := &Inode9{
		nodeId: fs.nextNodeId(),
		qid:    qid,
		refs:   1,
	}
	newInode.SwapFile(inodef)

	out.NodeId = newInode.nodeId
	out.Ino = newInode.qid.Path
	out.Generation = uint64(newInode.qid.Version)
	out.Mode = qidToMode(&newInode.qid)
	out.OpenFlags |= fuse.FOPEN_DIRECT_IO
	out.Fh = fs.nextFileHandle()

	fs.lock.Lock()
	fs.n2Inode[newInode.nodeId] = newInode
	fs.p2Inode[newInode.qid.Path] = newInode
	fs.fh2OpenFile[out.Fh] = &OpenFile{
		inode: inode,
		f:     f,
	}
	f = nil
	fs.lock.Unlock()

	return fuse.OK
}

func (fs *Proto9FS) Rename(cancel <-chan struct{}, in *fuse.RenameIn, oldName string, newName string) fuse.Status {
	fs.lock.Lock()
	parentInode := fs.n2Inode[in.NodeId]
	fs.lock.Unlock()

	parentf, ok := parentInode.GetFile()
	if !ok {
		return fuse.EIO
	}

	f, _, err := parentf.Walk([]string{oldName})
	if err != nil {
		return ErrToStatus(err)
	}
	defer f.Clunk()

	err = f.Rename(parentf, newName)
	if err != nil {
		return ErrToStatus(err)
	}

	return fuse.OK
}

func (fs *Proto9FS) Read(cancel <-chan struct{}, in *fuse.ReadIn, buf []byte) (fuse.ReadResult, fuse.Status) {
	fs.lock.Lock()
	f := fs.fh2OpenFile[in.Fh]
	fs.lock.Unlock()

	n, err := f.f.Read(uint64(in.Offset), buf)
	if err != nil {
		return nil, ErrToStatus(err)
	}
	return fuse.ReadResultData(buf[:n]), fuse.OK
}

func (fs *Proto9FS) Write(cancel <-chan struct{}, in *fuse.WriteIn, buf []byte) (uint32, fuse.Status) {
	fs.lock.Lock()
	f := fs.fh2OpenFile[in.Fh]
	fs.lock.Unlock()

	n, err := f.f.Write(uint64(in.Offset), buf)
	if err != nil {
		return 0, ErrToStatus(err)
	}
	return n, fuse.OK
}

func (fs *Proto9FS) setLk(cancel <-chan struct{}, in *fuse.LkIn, wait bool) fuse.Status {

	fs.lock.Lock()
	f := fs.fh2OpenFile[in.Fh]
	fs.lock.Unlock()

	typ9 := uint8(0)

	switch in.Lk.Typ {
	case syscall.F_RDLCK:
		typ9 = proto9.L_LOCK_TYPE_RDLCK
	case syscall.F_WRLCK:
		typ9 = proto9.L_LOCK_TYPE_WRLCK
	case syscall.F_UNLCK:
		typ9 = proto9.L_LOCK_TYPE_UNLCK
	default:
		return fuse.ENOTSUP
	}

	flags9 := uint32(0)
	if wait {
		flags9 |= proto9.L_LOCK_FLAGS_BLOCK
	}

	for {
		status, err := f.f.Lock(proto9.LSetLock{
			Typ:    typ9,
			Flags:  flags9,
			Start:  in.Lk.Start,
			Length: in.Lk.End - in.Lk.Start,
			ProcId: in.Lk.Pid,
		})
		if err != nil {
			return ErrToStatus(err)
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
			return fuse.EAGAIN
		default:
			return fuse.EIO
		}
	}
}

func (fs *Proto9FS) SetLk(cancel <-chan struct{}, in *fuse.LkIn) fuse.Status {
	return fs.setLk(cancel, in, true)
}

func (fs *Proto9FS) SetLkw(cancel <-chan struct{}, in *fuse.LkIn) fuse.Status {
	return fs.setLk(cancel, in, true)
}

func (fs *Proto9FS) Release(cancel <-chan struct{}, in *fuse.ReleaseIn) {
	fs.lock.Lock()
	f := fs.fh2OpenFile[in.Fh]
	delete(fs.fh2OpenFile, in.Fh)
	fs.lock.Unlock()
	_ = f.f.Clunk()
	return
}

func (fs *Proto9FS) remove(cancel <-chan struct{}, header *fuse.InHeader, name string) fuse.Status {
	fs.lock.Lock()
	inode := fs.n2Inode[header.NodeId]
	fs.lock.Unlock()

	f, ok := inode.GetFile()
	if !ok {
		return fuse.EIO
	}

	f, _, err := f.Walk([]string{name})
	if err != nil {
		return ErrToStatus(err)
	}

	err = f.Remove()
	if err != nil {
		return ErrToStatus(err)
	}
	return fuse.OK
}
func (fs *Proto9FS) Unlink(cancel <-chan struct{}, header *fuse.InHeader, name string) fuse.Status {
	return fs.remove(cancel, header, name)
}

func (fs *Proto9FS) Rmdir(cancel <-chan struct{}, header *fuse.InHeader, name string) fuse.Status {
	return fs.remove(cancel, header, name)
}

func (fs *Proto9FS) Mkdir(cancel <-chan struct{}, header *fuse.MkdirIn, name string, out *fuse.EntryOut) fuse.Status {

	fs.lock.Lock()
	parentInode := fs.n2Inode[header.NodeId]
	fs.lock.Unlock()

	parentf, ok := parentInode.GetFile()
	if !ok {
		return fuse.EIO
	}

	qid, err := parentf.Mkdir(name, header.Mode, header.InHeader.Caller.Gid)
	if err != nil {
		return ErrToStatus(err)
	}

	newInode := &Inode9{
		nodeId: fs.nextNodeId(),
		qid:    qid,
		refs:   1,
	}
	out.NodeId = newInode.nodeId
	out.Ino = newInode.qid.Path
	out.Generation = uint64(newInode.qid.Version)
	out.Mode = qidToMode(&newInode.qid)

	fs.lock.Lock()
	fs.n2Inode[newInode.nodeId] = newInode
	fs.p2Inode[newInode.qid.Path] = newInode
	fs.lock.Unlock()

	return fuse.OK
}

func (fs *Proto9FS) OpenDir(cancel <-chan struct{}, in *fuse.OpenIn, out *fuse.OpenOut) fuse.Status {

	fs.lock.Lock()
	inode := fs.n2Inode[in.NodeId]
	fs.lock.Unlock()

	dirf, ok := inode.GetFile()
	if !ok {
		return fuse.EIO
	}

	dirf, _, err := dirf.Walk([]string{})
	if err != nil {
		return ErrToStatus(err)
	}
	defer func() {
		if dirf != nil {
			_ = dirf.Clunk()
		}
	}()

	flags9, ok := openFlagsTo9(in.Flags)
	if !ok {
		return fuse.ENOTSUP
	}

	err = dirf.Open(flags9)
	if err != nil {
		return ErrToStatus(err)
	}

	out.Fh = fs.nextFileHandle()
	out.OpenFlags |= fuse.FOPEN_DIRECT_IO

	fs.lock.Lock()
	fs.fh2OpenFile[out.Fh] = &OpenFile{
		inode: inode,
		f:     dirf,
		di:    dirf.DirIter(),
	}
	dirf = nil
	fs.lock.Unlock()

	return fuse.OK
}

func (fs *Proto9FS) readDir(cancel <-chan struct{}, in *fuse.ReadIn, out *fuse.DirEntryList, plus bool) fuse.Status {
	fs.lock.Lock()
	d := fs.fh2OpenFile[in.Fh]
	fs.lock.Unlock()

	if d.di == nil {
		return fuse.EBADF
	}

	d.diLock.Lock()
	defer d.diLock.Unlock()

	// TODO verify offset is correct.
	for d.di.HasNext() {
		ent, err := d.di.Next()
		if err != nil {
			return ErrToStatus(err)
		}
		fuseDirEnt := fuse.DirEntry{
			Name: ent.Name,
			Mode: qidToMode(&ent.Qid),
			Ino:  ent.Qid.Path,
		}
		if plus {
			entryOut := out.AddDirLookupEntry(fuseDirEnt)
			if entryOut != nil {
				wnames := []string{}

				if ent.Name != "." {
					wnames = append(wnames, ent.Name)
				}

				dirf, ok := d.inode.GetFile()
				if !ok {
					return fuse.EIO
				}

				f, _, err := dirf.Walk(wnames)
				if err != nil {
					return ErrToStatus(err)
				}

				attr, err := f.GetAttr(proto9.L_GETATTR_ALL)
				_ = f.Clunk()
				if err != nil {
					return ErrToStatus(err)
				}
				FillFuseEntryOutFromAttr(&attr, entryOut)
			} else {
				d.di.Unget(ent)
				break
			}
		} else {
			if !out.AddDirEntry(fuseDirEnt) {
				d.di.Unget(ent)
				break
			}
		}
	}
	return fuse.OK
}

func (fs *Proto9FS) ReadDir(cancel <-chan struct{}, in *fuse.ReadIn, out *fuse.DirEntryList) fuse.Status {
	return fs.readDir(cancel, in, out, false)
}

func (fs *Proto9FS) ReadDirPlus(cancel <-chan struct{}, in *fuse.ReadIn, out *fuse.DirEntryList) fuse.Status {
	return fs.readDir(cancel, in, out, true)
}

func (fs *Proto9FS) fsync(cancel <-chan struct{}, in *fuse.FsyncIn) fuse.Status {
	fs.lock.Lock()
	f := fs.fh2OpenFile[in.Fh]
	fs.lock.Unlock()

	err := f.f.Fsync()
	if err != nil {
		return ErrToStatus(err)
	}
	return fuse.OK
}

func (fs *Proto9FS) Fsync(cancel <-chan struct{}, in *fuse.FsyncIn) fuse.Status {
	return fs.fsync(cancel, in)
}

func (fs *Proto9FS) FsyncDir(cancel <-chan struct{}, in *fuse.FsyncIn) fuse.Status {
	return fs.fsync(cancel, in)
}

func (fs *Proto9FS) ReleaseDir(in *fuse.ReleaseIn) {
	fs.lock.Lock()
	f := fs.fh2OpenFile[in.Fh]
	delete(fs.fh2OpenFile, in.Fh)
	fs.lock.Unlock()
	_ = f.f.Clunk()
	return
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
