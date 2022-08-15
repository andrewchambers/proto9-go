package proto9

type Qid struct {
	Typ     uint8
	Version uint32
	Path    uint64
}

type Tversion struct {
	Msize   uint32
	Version string
}

type Rversion struct {
	Msize   uint32
	Version string
}

type Tflush struct {
	OldTag uint32
}

type Rflush struct {
}

type Twalk struct {
	Fid    uint32
	NewFid uint32
	WQid   []Qid
}

type Rwalk struct {
	WQid []Qid
}

type Tread struct {
	Fid    uint32
	Offset uint64
	Count  uint32
}

type Rread struct {
	Count uint32
	Data  []byte
}

type Twrite struct {
	Fid    uint32
	Offset uint64
	Count  uint32
	Data   []byte
}

type Rwrite struct {
	Count uint32
}

type Tclunk struct {
	Fid uint32
}

type Rclunk struct {
}

type Tauth struct {
	Afid   uint32
	Uname  string
	Aname  string
	Nuname uint32
}

type Rauth struct {
	Aqid Qid
}

type Tattach struct {
	Afid    uint32
	uname   string
	aname   string
	n_uname uint32
}

type Rattach struct {
	Qid Qid
}

type Rlerror struct {
	Ecode uint32
}

type Tstatfs struct {
	Fid uint32
}

type Rstatfs struct {
	Typ     uint32
	Bsize   uint32
	Blocks  uint64
	Bfree   uint64
	Bavail  uint64
	Files   uint64
	Ffree   uint64
	Fsid    uint64
	Namelen uint32
}

type Tlopen struct {
	Fid   uint32
	Flags uint32
}

type Rlopen struct {
	Qid    Qid
	Iounit uint32
}

type Tlcreate struct {
	Fid   uint32
	Name  string
	Flags uint32
	Mode  uint32
	Gid   uint32
}

type Rlcreate struct {
	Qid    Qid
	Iounit uint32
}

type Tsymlink struct {
	Fid    uint32
	Name   string
	Target string
	Gid    uint32
}

type Rsymlink struct {
	Qid Qid
}

type Tmknod struct {
	Fid   uint32
	Name  string
	Major uint32
	Minor uint32
	Gid   uint32
}

type Rmknod struct {
	Qid Qid
}

type Trename struct {
	Fid  uint32
	Dfid uint32
	Name string
}

type Rrename struct {
}

type Treadlink struct {
	Fid uint32
}

type Rreadlink struct {
	Target string
}

const (
	GETATTR_MODE   uint64 = 0x00000001
	GETATTR_NLINK  uint64 = 0x00000002
	GETATTR_UID    uint64 = 0x00000004
	GETATTR_GID    uint64 = 0x00000008
	GETATTR_RDEV   uint64 = 0x00000010
	GETATTR_ATIME  uint64 = 0x00000020
	GETATTR_MTIME  uint64 = 0x00000040
	GETATTR_CTIME  uint64 = 0x00000080
	GETATTR_INO    uint64 = 0x00000100
	GETATTR_SIZE   uint64 = 0x00000200
	GETATTR_BLOCKS uint64 = 0x00000400

	GETATTR_BTIME        uint64 = 0x00000800
	GETATTR_GEN          uint64 = 0x00001000
	GETATTR_DATA_VERSION uint64 = 0x00002000

	GETATTR_BASIC uint64 = 0x000007ff /* Mask for fields up to BLOCKS */
	GETATTR_ALL   uint64 = 0x00003fff
)

type Tgetattr struct {
	Mask uint64
}

type Rgetattr struct {
	Valid       uint64
	Qid         Qid
	Mode        uint32
	Uid         uint32
	Gid         uint32
	Nlink       uint64
	Rdev        uint64
	Size        uint64
	Blksize     uint64
	Blocks      uint64
	AtimeSec    uint64
	AtimeNsec   uint64
	MtimeSec    uint64
	MtimeSsec   uint64
	CtimeSec    uint64
	CtimeNsec   uint64
	BtimeSec    uint64
	BtimeNsec   uint64
	Gen         uint64
	DataVersion uint64
}

const (
	SETATTR_MODE      uint32 = 0x00000001
	SETATTR_UID       uint32 = 0x00000002
	SETATTR_GID       uint32 = 0x00000004
	SETATTR_SIZE      uint32 = 0x00000008
	SETATTR_ATIME     uint32 = 0x00000010
	SETATTR_MTIME     uint32 = 0x00000020
	SETATTR_CTIME     uint32 = 0x00000040
	SETATTR_ATIME_SET uint32 = 0x00000080
	SETATTR_MTIME_SET uint32 = 0x00000100
)

type Tsetattr struct {
	Fid       uint32
	Valid     uint32
	Mode      uint32
	Uid       uint32
	Gid       uint32
	Size      uint64
	AtimeSec  uint64
	AtimeNsec uint64
	MtimeSec  uint64
	MtimeSsec uint64
}

type Rsetattr struct {
}

type Txattrwalk struct {
	Fid    uint32
	Newfid uint32
	Name   string
}

type Rxattrwalk struct {
	Size uint64
}

type Txattrcreate struct {
	Fid      uint32
	Name     string
	AttrSize uint64
	Flags    uint32
}

type Rxattrcreate struct {
}

type Treaddir struct {
	Fid    uint32
	Offset uint64
	Count  uint32
}

type Rreaddir struct {
	Data []byte
}

type DirEnt struct {
	Qid    Qid
	Offset uint64
	Typ    uint8
	Name   string
}

type Tfsync struct {
	Fid uint32
}

type Rfsync struct {
}

const (
	LOCK_TYPE_RDLCK byte = 0
	LOCK_TYPE_WRLCK byte = 1
	LOCK_TYPE_UNLCK byte = 2

	LOCK_FLAGS_BLOCK   uint32 = 1
	LOCK_FLAGS_RECLAIM uint32 = 2

	LOCK_SUCCESS byte = 0
	LOCK_BLOCKED byte = 1
	LOCK_ERROR   byte = 2
	LOCK_GRACE   byte = 3
)

type Tlock struct {
	Fid      uint32
	Typ      byte
	Flags    uint32
	Start    uint64
	Length   uint64
	ProcId   uint32
	ClientId string
}

type Rlock struct {
	Status byte
}

type Tgetlock struct {
	Fid      uint32
	Typ      byte
	Start    uint64
	Length   uint64
	ProcId   uint32
	ClientId string
}

type Rgetlock struct {
	Typ      byte
	Start    uint64
	Length   uint64
	ProcId   uint32
	ClientId string
}

type Tlink struct {
	Dfid uint32
	Fid  uint32
	Name string
}

type Rlink struct {
}

type Tmkdir struct {
	Dfid uint32
	Fid  uint32
	Name string
	Mode uint32
	Gid  uint32
}

type Rmkdir struct {
	Qid Qid
}

type Trenameat struct {
	OldDirFid uint32
	OldName   string
	NewDirFid uint32
	NewName   string
}

type Rrenameat struct {
}

type Tunlinkat struct {
	DirFid uint32
	Name   string
	Flags  uint32
}

type Runlinkat struct {
}
