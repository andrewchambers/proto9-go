package proto9

import (
	"bytes"
	"fmt"
)

type Fcall interface {
	Kind() uint8
	GetTag() uint16
	SetTag(uint16)
	EncodedSize() uint64
	Encode(*bytes.Buffer) error
	Decode(*bytes.Buffer) error
}

type Tagged struct {
	Tag uint16
}

func (t *Tagged) GetTag() uint16 {
	return t.Tag
}

func (t *Tagged) SetTag(tag uint16) {
	t.Tag = tag
}

type Qid struct {
	Typ     uint8
	Version uint32
	Path    uint64
}

type Tversion struct {
	Tagged
	Msize   uint32
	Version string
}

type Rversion struct {
	Tagged
	Msize   uint32
	Version string
}

type Tflush struct {
	Tagged
	OldTag uint16
}

type Rflush struct {
	Tagged
}

type Twalk struct {
	Tagged
	Fid    uint32
	NewFid uint32
	Wnames []string
}

type Rwalk struct {
	Tagged
	WQids []Qid
}

type Tread struct {
	Tagged
	Fid    uint32
	Offset uint64
	Count  uint32
}

type Rread struct {
	Tagged
	Data []byte
}

type Twrite struct {
	Tagged
	Fid    uint32
	Offset uint64
	Data   []byte
}

type Rwrite struct {
	Tagged
	Count uint32
}

type Tclunk struct {
	Tagged
	Fid uint32
}

type Rclunk struct {
	Tagged
}

type Tremove struct {
	Tagged
	Fid uint32
}

type Rremove struct {
	Tagged
}

type Tauth struct {
	Tagged
	Afid   uint32
	Uname  string
	Aname  string
	Nuname uint32
}

type Rauth struct {
	Tagged
	Aqid Qid
}

type Tattach struct {
	Tagged
	Fid     uint32
	Afid    uint32
	Uname   string
	Aname   string
	N_uname uint32
}

type Rattach struct {
	Tagged
	Qid Qid
}

type Rlerror struct {
	Tagged
	Ecode uint32
}

type Tstatfs struct {
	Tagged
	Fid uint32
}

type LStatfs struct {
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

type Rstatfs struct {
	Tagged
	LStatfs
}

type Tlopen struct {
	Tagged
	Fid   uint32
	Flags uint32
}

type Rlopen struct {
	Tagged
	Qid    Qid
	Iounit uint32
}

type Tlcreate struct {
	Tagged
	Fid   uint32
	Name  string
	Flags uint32
	Mode  uint32
	Gid   uint32
}

type Rlcreate struct {
	Tagged
	Qid    Qid
	Iounit uint32
}

type Tsymlink struct {
	Tagged
	Fid    uint32
	Name   string
	Target string
	Gid    uint32
}

type Rsymlink struct {
	Tagged
	Qid Qid
}

type Tmknod struct {
	Tagged
	Fid   uint32
	Name  string
	Major uint32
	Minor uint32
	Gid   uint32
}

type Rmknod struct {
	Tagged
	Qid Qid
}

type Trename struct {
	Tagged
	Fid  uint32
	Dfid uint32
	Name string
}

type Rrename struct {
	Tagged
}

type Treadlink struct {
	Tagged
	Fid uint32
}

type Rreadlink struct {
	Tagged
	Target string
}

type Tgetattr struct {
	Tagged
	Fid  uint32
	Mask uint64
}

type LAttr struct {
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
	MtimeNsec   uint64
	CtimeSec    uint64
	CtimeNsec   uint64
	BtimeSec    uint64
	BtimeNsec   uint64
	Gen         uint64
	DataVersion uint64
}

type Rgetattr struct {
	Tagged
	LAttr
}

type LSetAttr struct {
	Valid     uint32
	Mode      uint32
	Uid       uint32
	Gid       uint32
	Size      uint64
	AtimeSec  uint64
	AtimeNsec uint64
	MtimeSec  uint64
	MtimeNsec uint64
}

type Tsetattr struct {
	Tagged
	Fid uint32
	LSetAttr
}

type Rsetattr struct {
	Tagged
}

type Txattrwalk struct {
	Tagged
	Fid    uint32
	Newfid uint32
	Name   string
}

type Rxattrwalk struct {
	Tagged
	Size uint64
}

type Txattrcreate struct {
	Tagged
	Fid      uint32
	Name     string
	AttrSize uint64
	Flags    uint32
}

type Rxattrcreate struct {
	Tagged
}

type Treaddir struct {
	Tagged
	Fid    uint32
	Offset uint64
	Count  uint32
}

type Rreaddir struct {
	Tagged
	Data []DirEnt
}

type DirEnt struct {
	Qid    Qid
	Offset uint64
	Typ    uint8
	Name   string
}

type Tfsync struct {
	Tagged
	Fid uint32
}

type Rfsync struct {
	Tagged
}

type LSetLock struct {
	Typ      byte
	Flags    uint32
	Start    uint64
	Length   uint64
	ProcId   uint32
	ClientId string
}

type Tlock struct {
	Tagged
	Fid uint32
	LSetLock
}

type Rlock struct {
	Tagged
	Status byte
}

type LGetLock struct {
	Typ      byte
	Start    uint64
	Length   uint64
	ProcId   uint32
	ClientId string
}

type Tgetlock struct {
	Tagged
	Fid uint32
	LGetLock
}

type Rgetlock struct {
	Tagged
	LGetLock
}

type Tlink struct {
	Tagged
	Dfid uint32
	Fid  uint32
	Name string
}

type Rlink struct {
	Tagged
}

type Tmkdir struct {
	Tagged
	Dfid uint32
	Name string
	Mode uint32
	Gid  uint32
}

type Rmkdir struct {
	Tagged
	Qid Qid
}

type Trenameat struct {
	Tagged
	OldDfid uint32
	OldName string
	NewDfid uint32
	NewName string
}

type Rrenameat struct {
	Tagged
}

type Tunlinkat struct {
	Tagged
	Dfid  uint32
	Name  string
	Flags uint32
}

type Runlinkat struct {
	Tagged
}

func FcallFromKind(kind uint8) (Fcall, error) {
	switch kind {
	// 9P2000.L
	case 7:
		return &Rlerror{}, nil
	case 8:
		return &Tstatfs{}, nil
	case 9:
		return &Rstatfs{}, nil
	case 12:
		return &Tlopen{}, nil
	case 13:
		return &Rlopen{}, nil
	case 14:
		return &Tlcreate{}, nil
	case 15:
		return &Rlcreate{}, nil
	case 16:
		return &Tsymlink{}, nil
	case 17:
		return &Rsymlink{}, nil
	case 18:
		return &Tmknod{}, nil
	case 19:
		return &Rmknod{}, nil
	case 20:
		return &Trename{}, nil
	case 21:
		return &Rrename{}, nil
	case 22:
		return &Treadlink{}, nil
	case 23:
		return &Rreadlink{}, nil
	case 24:
		return &Tgetattr{}, nil
	case 25:
		return &Rgetattr{}, nil
	case 26:
		return &Tsetattr{}, nil
	case 27:
		return &Rsetattr{}, nil
	case 30:
		return &Txattrwalk{}, nil
	case 31:
		return &Rxattrwalk{}, nil
	case 32:
		return &Txattrcreate{}, nil
	case 33:
		return &Rxattrcreate{}, nil
	case 40:
		return &Treaddir{}, nil
	case 41:
		return &Rreaddir{}, nil
	case 50:
		return &Tfsync{}, nil
	case 51:
		return &Rfsync{}, nil
	case 52:
		return &Tlock{}, nil
	case 53:
		return &Rlock{}, nil
	case 54:
		return &Tgetlock{}, nil
	case 55:
		return &Rgetlock{}, nil
	case 70:
		return &Tlink{}, nil
	case 71:
		return &Rlink{}, nil
	case 72:
		return &Tmkdir{}, nil
	case 73:
		return &Rmkdir{}, nil
	case 74:
		return &Trenameat{}, nil
	case 75:
		return &Rrenameat{}, nil
	case 76:
		return &Tunlinkat{}, nil
	case 77:
		return &Runlinkat{}, nil
	// 9P2000
	case 100:
		return &Tversion{}, nil
	case 101:
		return &Rversion{}, nil
	case 102:
		return &Tauth{}, nil
	case 103:
		return &Rauth{}, nil
	case 104:
		return &Tattach{}, nil
	case 105:
		return &Rattach{}, nil
	// case 106:
	//	return &Terror{}, nil
	// case 107:
	// 	return &Rerror{}, nil
	case 108:
		return &Tflush{}, nil
	case 109:
		return &Rflush{}, nil
	case 110:
		return &Twalk{}, nil
	case 111:
		return &Rwalk{}, nil
	// case 112:
	//	return &Topen{}, nil
	// case 113:
	//	return &Ropen{}, nil
	// case 114:
	//	return &Tcreate{}, nil
	// case 115:
	//	return &Rcreate {}, nil
	case 116:
		return &Tread{}, nil
	case 117:
		return &Rread{}, nil
	case 118:
		return &Twrite{}, nil
	case 119:
		return &Rwrite{}, nil
	case 120:
		return &Tclunk{}, nil
	case 121:
		return &Rclunk{}, nil
	case 122:
		return &Tremove{}, nil
	case 123:
		return &Rremove{}, nil
	// case 124:
	//	return &Tstat{}, nil
	// case 125:
	//	return &Rstat{}, nil
	// case 126:
	//	return &Twstat{}, nil
	// case 126:
	//	return &Rwstat{}, nil
	default:
		return nil, fmt.Errorf("unknown message kind: %d", kind)
	}
}

func (m *Rlerror) Kind() uint8      { return 7 }
func (m *Tstatfs) Kind() uint8      { return 8 }
func (m *Rstatfs) Kind() uint8      { return 9 }
func (m *Tlopen) Kind() uint8       { return 12 }
func (m *Rlopen) Kind() uint8       { return 13 }
func (m *Tlcreate) Kind() uint8     { return 14 }
func (m *Rlcreate) Kind() uint8     { return 15 }
func (m *Tsymlink) Kind() uint8     { return 16 }
func (m *Rsymlink) Kind() uint8     { return 17 }
func (m *Tmknod) Kind() uint8       { return 18 }
func (m *Rmknod) Kind() uint8       { return 19 }
func (m *Trename) Kind() uint8      { return 20 }
func (m *Rrename) Kind() uint8      { return 21 }
func (m *Treadlink) Kind() uint8    { return 22 }
func (m *Rreadlink) Kind() uint8    { return 23 }
func (m *Tgetattr) Kind() uint8     { return 24 }
func (m *Rgetattr) Kind() uint8     { return 25 }
func (m *Tsetattr) Kind() uint8     { return 26 }
func (m *Rsetattr) Kind() uint8     { return 27 }
func (m *Txattrwalk) Kind() uint8   { return 30 }
func (m *Rxattrwalk) Kind() uint8   { return 31 }
func (m *Txattrcreate) Kind() uint8 { return 32 }
func (m *Rxattrcreate) Kind() uint8 { return 33 }
func (m *Treaddir) Kind() uint8     { return 40 }
func (m *Rreaddir) Kind() uint8     { return 41 }
func (m *Tfsync) Kind() uint8       { return 50 }
func (m *Rfsync) Kind() uint8       { return 51 }
func (m *Tlock) Kind() uint8        { return 52 }
func (m *Rlock) Kind() uint8        { return 53 }
func (m *Tgetlock) Kind() uint8     { return 54 }
func (m *Rgetlock) Kind() uint8     { return 55 }
func (m *Tlink) Kind() uint8        { return 70 }
func (m *Rlink) Kind() uint8        { return 71 }
func (m *Tmkdir) Kind() uint8       { return 72 }
func (m *Rmkdir) Kind() uint8       { return 73 }
func (m *Trenameat) Kind() uint8    { return 74 }
func (m *Rrenameat) Kind() uint8    { return 75 }
func (m *Tunlinkat) Kind() uint8    { return 76 }
func (m *Runlinkat) Kind() uint8    { return 77 }
func (m *Tversion) Kind() uint8     { return 100 }
func (m *Rversion) Kind() uint8     { return 101 }
func (m *Tauth) Kind() uint8        { return 102 }
func (m *Rauth) Kind() uint8        { return 103 }
func (m *Tattach) Kind() uint8      { return 104 }
func (m *Rattach) Kind() uint8      { return 105 }

// func (m *Terror) Kind() uint8       { return 106 }
// func (m *Rerror) Kind() uint8       { return 107 }
func (m *Tflush) Kind() uint8 { return 108 }
func (m *Rflush) Kind() uint8 { return 109 }
func (m *Twalk) Kind() uint8  { return 110 }
func (m *Rwalk) Kind() uint8  { return 111 }

// func (m *Topen) Kind() uint8        { return 112 }
// func (m *Ropen) Kind() uint8        { return 113 }
// func (m *Tcreate) Kind() uint8      { return 114 }
// func (m *Rcreate) Kind() uint8      { return 115 }
func (m *Tread) Kind() uint8   { return 116 }
func (m *Rread) Kind() uint8   { return 117 }
func (m *Twrite) Kind() uint8  { return 118 }
func (m *Rwrite) Kind() uint8  { return 119 }
func (m *Tclunk) Kind() uint8  { return 120 }
func (m *Rclunk) Kind() uint8  { return 121 }
func (m *Tremove) Kind() uint8 { return 122 }
func (m *Rremove) Kind() uint8 { return 123 }

// func (m *Tstat) Kind() uint8        { return 124 }
// func (m *Rstat) Kind() uint8        { return 125 }
// func (m *Twstat) Kind() uint8       { return 126 }
// func (m *Rwstat) Kind() uint8       { return 126 }
