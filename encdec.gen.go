package proto9

import (
	"bytes"
)

func (v *DirEnt) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Qid.EncodedSize()
	sz += 8 // Offset
	sz += 1 // Typ
	sz += 2 + uint64(len(v.Name))
	return sz
}

func (v *DirEnt) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Qid.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.Offset)
	if err != nil {
		return err
	}
	err = encodeUint8(b, v.Typ)
	if err != nil {
		return err
	}
	err = encodeString(b, v.Name)
	if err != nil {
		return err
	}
	return nil
}

func (v *DirEnt) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Qid.Decode(b)
	if err != nil {
		return err
	}
	v.Offset, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.Typ, err = decodeUint8(b)
	if err != nil {
		return err
	}
	v.Name, err = decodeString(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *LAttr) EncodedSize() uint64 {
	sz := uint64(0)
	sz += 8 // Valid
	sz += v.Qid.EncodedSize()
	sz += 4 // Mode
	sz += 4 // Uid
	sz += 4 // Gid
	sz += 8 // Nlink
	sz += 8 // Rdev
	sz += 8 // Size
	sz += 8 // Blksize
	sz += 8 // Blocks
	sz += 8 // AtimeSec
	sz += 8 // AtimeNsec
	sz += 8 // MtimeSec
	sz += 8 // MtimeNsec
	sz += 8 // CtimeSec
	sz += 8 // CtimeNsec
	sz += 8 // BtimeSec
	sz += 8 // BtimeNsec
	sz += 8 // Gen
	sz += 8 // DataVersion
	return sz
}

func (v *LAttr) Encode(b *bytes.Buffer) error {
	var err error
	err = encodeUint64(b, v.Valid)
	if err != nil {
		return err
	}
	err = v.Qid.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Mode)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Uid)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Gid)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.Nlink)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.Rdev)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.Size)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.Blksize)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.Blocks)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.AtimeSec)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.AtimeNsec)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.MtimeSec)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.MtimeNsec)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.CtimeSec)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.CtimeNsec)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.BtimeSec)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.BtimeNsec)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.Gen)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.DataVersion)
	if err != nil {
		return err
	}
	return nil
}

func (v *LAttr) Decode(b *bytes.Buffer) error {
	var err error
	v.Valid, err = decodeUint64(b)
	if err != nil {
		return err
	}
	err = v.Qid.Decode(b)
	if err != nil {
		return err
	}
	v.Mode, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Uid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Gid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Nlink, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.Rdev, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.Size, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.Blksize, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.Blocks, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.AtimeSec, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.AtimeNsec, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.MtimeSec, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.MtimeNsec, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.CtimeSec, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.CtimeNsec, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.BtimeSec, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.BtimeNsec, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.Gen, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.DataVersion, err = decodeUint64(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *LSetAttr) EncodedSize() uint64 {
	sz := uint64(0)
	sz += 4 // Valid
	sz += 4 // Mode
	sz += 4 // Uid
	sz += 4 // Gid
	sz += 8 // Size
	sz += 8 // AtimeSec
	sz += 8 // AtimeNsec
	sz += 8 // MtimeSec
	sz += 8 // MtimeSsec
	return sz
}

func (v *LSetAttr) Encode(b *bytes.Buffer) error {
	var err error
	err = encodeUint32(b, v.Valid)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Mode)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Uid)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Gid)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.Size)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.AtimeSec)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.AtimeNsec)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.MtimeSec)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.MtimeSsec)
	if err != nil {
		return err
	}
	return nil
}

func (v *LSetAttr) Decode(b *bytes.Buffer) error {
	var err error
	v.Valid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Mode, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Uid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Gid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Size, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.AtimeSec, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.AtimeNsec, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.MtimeSec, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.MtimeSsec, err = decodeUint64(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *LStatfs) EncodedSize() uint64 {
	sz := uint64(0)
	sz += 4 // Typ
	sz += 4 // Bsize
	sz += 8 // Blocks
	sz += 8 // Bfree
	sz += 8 // Bavail
	sz += 8 // Files
	sz += 8 // Ffree
	sz += 8 // Fsid
	sz += 4 // Namelen
	return sz
}

func (v *LStatfs) Encode(b *bytes.Buffer) error {
	var err error
	err = encodeUint32(b, v.Typ)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Bsize)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.Blocks)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.Bfree)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.Bavail)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.Files)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.Ffree)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.Fsid)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Namelen)
	if err != nil {
		return err
	}
	return nil
}

func (v *LStatfs) Decode(b *bytes.Buffer) error {
	var err error
	v.Typ, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Bsize, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Blocks, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.Bfree, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.Bavail, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.Files, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.Ffree, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.Fsid, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.Namelen, err = decodeUint32(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Qid) EncodedSize() uint64 {
	sz := uint64(0)
	sz += 1 // Typ
	sz += 4 // Version
	sz += 8 // Path
	return sz
}

func (v *Qid) Encode(b *bytes.Buffer) error {
	var err error
	err = encodeUint8(b, v.Typ)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Version)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.Path)
	if err != nil {
		return err
	}
	return nil
}

func (v *Qid) Decode(b *bytes.Buffer) error {
	var err error
	v.Typ, err = decodeUint8(b)
	if err != nil {
		return err
	}
	v.Version, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Path, err = decodeUint64(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rattach) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += v.Qid.EncodedSize()
	return sz
}

func (v *Rattach) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = v.Qid.Encode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rattach) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	err = v.Qid.Decode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rauth) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += v.Aqid.EncodedSize()
	return sz
}

func (v *Rauth) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = v.Aqid.Encode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rauth) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	err = v.Aqid.Decode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rclunk) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	return sz
}

func (v *Rclunk) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rclunk) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rflush) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	return sz
}

func (v *Rflush) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rflush) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rfsync) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	return sz
}

func (v *Rfsync) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rfsync) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rgetattr) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += v.LAttr.EncodedSize()
	return sz
}

func (v *Rgetattr) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = v.LAttr.Encode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rgetattr) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	err = v.LAttr.Decode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rgetlock) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 1 // Typ
	sz += 8 // Start
	sz += 8 // Length
	sz += 4 // ProcId
	sz += 2 + uint64(len(v.ClientId))
	return sz
}

func (v *Rgetlock) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeByte(b, v.Typ)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.Start)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.Length)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.ProcId)
	if err != nil {
		return err
	}
	err = encodeString(b, v.ClientId)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rgetlock) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Typ, err = decodeByte(b)
	if err != nil {
		return err
	}
	v.Start, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.Length, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.ProcId, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.ClientId, err = decodeString(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rlcreate) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += v.Qid.EncodedSize()
	sz += 4 // Iounit
	return sz
}

func (v *Rlcreate) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = v.Qid.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Iounit)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rlcreate) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	err = v.Qid.Decode(b)
	if err != nil {
		return err
	}
	v.Iounit, err = decodeUint32(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rlerror) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // Ecode
	return sz
}

func (v *Rlerror) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Ecode)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rlerror) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Ecode, err = decodeUint32(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rlink) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	return sz
}

func (v *Rlink) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rlink) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rlock) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 1 // Status
	return sz
}

func (v *Rlock) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeByte(b, v.Status)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rlock) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Status, err = decodeByte(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rlopen) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += v.Qid.EncodedSize()
	sz += 4 // Iounit
	return sz
}

func (v *Rlopen) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = v.Qid.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Iounit)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rlopen) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	err = v.Qid.Decode(b)
	if err != nil {
		return err
	}
	v.Iounit, err = decodeUint32(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rmkdir) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += v.Qid.EncodedSize()
	return sz
}

func (v *Rmkdir) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = v.Qid.Encode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rmkdir) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	err = v.Qid.Decode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rmknod) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += v.Qid.EncodedSize()
	return sz
}

func (v *Rmknod) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = v.Qid.Encode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rmknod) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	err = v.Qid.Decode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rread) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 + uint64(len(v.Data))
	return sz
}

func (v *Rread) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeByteSlice(b, v.Data)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rread) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Data, err = decodeByteSlice(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rreaddir) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 + uint64(len(v.Data))
	return sz
}

func (v *Rreaddir) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeByteSlice(b, v.Data)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rreaddir) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Data, err = decodeByteSlice(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rreadlink) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 2 + uint64(len(v.Target))
	return sz
}

func (v *Rreadlink) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeString(b, v.Target)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rreadlink) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Target, err = decodeString(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rremove) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	return sz
}

func (v *Rremove) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rremove) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rrename) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	return sz
}

func (v *Rrename) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rrename) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rrenameat) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	return sz
}

func (v *Rrenameat) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rrenameat) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rsetattr) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	return sz
}

func (v *Rsetattr) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rsetattr) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rstatfs) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += v.LStatfs.EncodedSize()
	return sz
}

func (v *Rstatfs) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = v.LStatfs.Encode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rstatfs) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	err = v.LStatfs.Decode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rsymlink) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += v.Qid.EncodedSize()
	return sz
}

func (v *Rsymlink) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = v.Qid.Encode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rsymlink) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	err = v.Qid.Decode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Runlinkat) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	return sz
}

func (v *Runlinkat) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Runlinkat) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rversion) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // Msize
	sz += 2 + uint64(len(v.Version))
	return sz
}

func (v *Rversion) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Msize)
	if err != nil {
		return err
	}
	err = encodeString(b, v.Version)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rversion) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Msize, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Version, err = decodeString(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rwalk) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 2 + uint64(len(v.WQids))*13
	return sz
}

func (v *Rwalk) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeQids(b, v.WQids)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rwalk) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.WQids, err = decodeQids(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rwrite) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // Count
	return sz
}

func (v *Rwrite) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Count)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rwrite) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Count, err = decodeUint32(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rxattrcreate) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	return sz
}

func (v *Rxattrcreate) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rxattrcreate) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rxattrwalk) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 8 // Size
	return sz
}

func (v *Rxattrwalk) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.Size)
	if err != nil {
		return err
	}
	return nil
}

func (v *Rxattrwalk) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Size, err = decodeUint64(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tagged) EncodedSize() uint64 {
	sz := uint64(0)
	sz += 2 // Tag
	return sz
}

func (v *Tagged) Encode(b *bytes.Buffer) error {
	var err error
	err = encodeUint16(b, v.Tag)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tagged) Decode(b *bytes.Buffer) error {
	var err error
	v.Tag, err = decodeUint16(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tattach) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // Fid
	sz += 4 // Afid
	sz += 2 + uint64(len(v.Uname))
	sz += 2 + uint64(len(v.Aname))
	sz += 4 // N_uname
	return sz
}

func (v *Tattach) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Fid)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Afid)
	if err != nil {
		return err
	}
	err = encodeString(b, v.Uname)
	if err != nil {
		return err
	}
	err = encodeString(b, v.Aname)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.N_uname)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tattach) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Fid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Afid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Uname, err = decodeString(b)
	if err != nil {
		return err
	}
	v.Aname, err = decodeString(b)
	if err != nil {
		return err
	}
	v.N_uname, err = decodeUint32(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tauth) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // Afid
	sz += 2 + uint64(len(v.Uname))
	sz += 2 + uint64(len(v.Aname))
	sz += 4 // Nuname
	return sz
}

func (v *Tauth) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Afid)
	if err != nil {
		return err
	}
	err = encodeString(b, v.Uname)
	if err != nil {
		return err
	}
	err = encodeString(b, v.Aname)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Nuname)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tauth) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Afid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Uname, err = decodeString(b)
	if err != nil {
		return err
	}
	v.Aname, err = decodeString(b)
	if err != nil {
		return err
	}
	v.Nuname, err = decodeUint32(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tclunk) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // Fid
	return sz
}

func (v *Tclunk) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Fid)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tclunk) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Fid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tflush) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 2 // OldTag
	return sz
}

func (v *Tflush) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint16(b, v.OldTag)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tflush) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.OldTag, err = decodeUint16(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tfsync) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // Fid
	return sz
}

func (v *Tfsync) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Fid)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tfsync) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Fid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tgetattr) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // Fid
	sz += 8 // Mask
	return sz
}

func (v *Tgetattr) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Fid)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.Mask)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tgetattr) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Fid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Mask, err = decodeUint64(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tgetlock) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // Fid
	sz += 1 // Typ
	sz += 8 // Start
	sz += 8 // Length
	sz += 4 // ProcId
	sz += 2 + uint64(len(v.ClientId))
	return sz
}

func (v *Tgetlock) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Fid)
	if err != nil {
		return err
	}
	err = encodeByte(b, v.Typ)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.Start)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.Length)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.ProcId)
	if err != nil {
		return err
	}
	err = encodeString(b, v.ClientId)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tgetlock) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Fid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Typ, err = decodeByte(b)
	if err != nil {
		return err
	}
	v.Start, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.Length, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.ProcId, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.ClientId, err = decodeString(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tlcreate) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // Fid
	sz += 2 + uint64(len(v.Name))
	sz += 4 // Flags
	sz += 4 // Mode
	sz += 4 // Gid
	return sz
}

func (v *Tlcreate) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Fid)
	if err != nil {
		return err
	}
	err = encodeString(b, v.Name)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Flags)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Mode)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Gid)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tlcreate) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Fid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Name, err = decodeString(b)
	if err != nil {
		return err
	}
	v.Flags, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Mode, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Gid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tlink) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // Dfid
	sz += 4 // Fid
	sz += 2 + uint64(len(v.Name))
	return sz
}

func (v *Tlink) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Dfid)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Fid)
	if err != nil {
		return err
	}
	err = encodeString(b, v.Name)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tlink) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Dfid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Fid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Name, err = decodeString(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tlock) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // Fid
	sz += 1 // Typ
	sz += 4 // Flags
	sz += 8 // Start
	sz += 8 // Length
	sz += 4 // ProcId
	sz += 2 + uint64(len(v.ClientId))
	return sz
}

func (v *Tlock) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Fid)
	if err != nil {
		return err
	}
	err = encodeByte(b, v.Typ)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Flags)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.Start)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.Length)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.ProcId)
	if err != nil {
		return err
	}
	err = encodeString(b, v.ClientId)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tlock) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Fid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Typ, err = decodeByte(b)
	if err != nil {
		return err
	}
	v.Flags, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Start, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.Length, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.ProcId, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.ClientId, err = decodeString(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tlopen) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // Fid
	sz += 4 // Flags
	return sz
}

func (v *Tlopen) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Fid)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Flags)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tlopen) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Fid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Flags, err = decodeUint32(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tmkdir) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // Dfid
	sz += 2 + uint64(len(v.Name))
	sz += 4 // Mode
	sz += 4 // Gid
	return sz
}

func (v *Tmkdir) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Dfid)
	if err != nil {
		return err
	}
	err = encodeString(b, v.Name)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Mode)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Gid)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tmkdir) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Dfid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Name, err = decodeString(b)
	if err != nil {
		return err
	}
	v.Mode, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Gid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tmknod) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // Fid
	sz += 2 + uint64(len(v.Name))
	sz += 4 // Major
	sz += 4 // Minor
	sz += 4 // Gid
	return sz
}

func (v *Tmknod) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Fid)
	if err != nil {
		return err
	}
	err = encodeString(b, v.Name)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Major)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Minor)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Gid)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tmknod) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Fid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Name, err = decodeString(b)
	if err != nil {
		return err
	}
	v.Major, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Minor, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Gid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tread) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // Fid
	sz += 8 // Offset
	sz += 4 // Count
	return sz
}

func (v *Tread) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Fid)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.Offset)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Count)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tread) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Fid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Offset, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.Count, err = decodeUint32(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Treaddir) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // Fid
	sz += 8 // Offset
	sz += 4 // Count
	return sz
}

func (v *Treaddir) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Fid)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.Offset)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Count)
	if err != nil {
		return err
	}
	return nil
}

func (v *Treaddir) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Fid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Offset, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.Count, err = decodeUint32(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Treadlink) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // Fid
	return sz
}

func (v *Treadlink) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Fid)
	if err != nil {
		return err
	}
	return nil
}

func (v *Treadlink) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Fid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tremove) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // Fid
	return sz
}

func (v *Tremove) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Fid)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tremove) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Fid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Trename) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // Fid
	sz += 4 // Dfid
	sz += 2 + uint64(len(v.Name))
	return sz
}

func (v *Trename) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Fid)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Dfid)
	if err != nil {
		return err
	}
	err = encodeString(b, v.Name)
	if err != nil {
		return err
	}
	return nil
}

func (v *Trename) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Fid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Dfid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Name, err = decodeString(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Trenameat) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // OldDfid
	sz += 2 + uint64(len(v.OldName))
	sz += 4 // NewDfid
	sz += 2 + uint64(len(v.NewName))
	return sz
}

func (v *Trenameat) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.OldDfid)
	if err != nil {
		return err
	}
	err = encodeString(b, v.OldName)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.NewDfid)
	if err != nil {
		return err
	}
	err = encodeString(b, v.NewName)
	if err != nil {
		return err
	}
	return nil
}

func (v *Trenameat) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.OldDfid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.OldName, err = decodeString(b)
	if err != nil {
		return err
	}
	v.NewDfid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.NewName, err = decodeString(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tsetattr) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // Fid
	sz += v.LSetAttr.EncodedSize()
	return sz
}

func (v *Tsetattr) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Fid)
	if err != nil {
		return err
	}
	err = v.LSetAttr.Encode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tsetattr) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Fid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	err = v.LSetAttr.Decode(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tstatfs) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // Fid
	return sz
}

func (v *Tstatfs) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Fid)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tstatfs) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Fid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tsymlink) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // Fid
	sz += 2 + uint64(len(v.Name))
	sz += 2 + uint64(len(v.Target))
	sz += 4 // Gid
	return sz
}

func (v *Tsymlink) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Fid)
	if err != nil {
		return err
	}
	err = encodeString(b, v.Name)
	if err != nil {
		return err
	}
	err = encodeString(b, v.Target)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Gid)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tsymlink) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Fid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Name, err = decodeString(b)
	if err != nil {
		return err
	}
	v.Target, err = decodeString(b)
	if err != nil {
		return err
	}
	v.Gid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tunlinkat) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // Dfid
	sz += 2 + uint64(len(v.Name))
	sz += 4 // Flags
	return sz
}

func (v *Tunlinkat) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Dfid)
	if err != nil {
		return err
	}
	err = encodeString(b, v.Name)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Flags)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tunlinkat) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Dfid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Name, err = decodeString(b)
	if err != nil {
		return err
	}
	v.Flags, err = decodeUint32(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tversion) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // Msize
	sz += 2 + uint64(len(v.Version))
	return sz
}

func (v *Tversion) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Msize)
	if err != nil {
		return err
	}
	err = encodeString(b, v.Version)
	if err != nil {
		return err
	}
	return nil
}

func (v *Tversion) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Msize, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Version, err = decodeString(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Twalk) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // Fid
	sz += 4 // NewFid
	// Wnames
	sz += 2
	for _, s := range v.Wnames {
		sz += 2 + uint64(len(s))
	}
	return sz
}

func (v *Twalk) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Fid)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.NewFid)
	if err != nil {
		return err
	}
	err = encodeStringSlice(b, v.Wnames)
	if err != nil {
		return err
	}
	return nil
}

func (v *Twalk) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Fid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.NewFid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Wnames, err = decodeStringSlice(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Twrite) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // Fid
	sz += 8 // Offset
	sz += 4 + uint64(len(v.Data))
	return sz
}

func (v *Twrite) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Fid)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.Offset)
	if err != nil {
		return err
	}
	err = encodeByteSlice(b, v.Data)
	if err != nil {
		return err
	}
	return nil
}

func (v *Twrite) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Fid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Offset, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.Data, err = decodeByteSlice(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Txattrcreate) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // Fid
	sz += 2 + uint64(len(v.Name))
	sz += 8 // AttrSize
	sz += 4 // Flags
	return sz
}

func (v *Txattrcreate) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Fid)
	if err != nil {
		return err
	}
	err = encodeString(b, v.Name)
	if err != nil {
		return err
	}
	err = encodeUint64(b, v.AttrSize)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Flags)
	if err != nil {
		return err
	}
	return nil
}

func (v *Txattrcreate) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Fid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Name, err = decodeString(b)
	if err != nil {
		return err
	}
	v.AttrSize, err = decodeUint64(b)
	if err != nil {
		return err
	}
	v.Flags, err = decodeUint32(b)
	if err != nil {
		return err
	}
	return nil
}

func (v *Txattrwalk) EncodedSize() uint64 {
	sz := uint64(0)
	sz += v.Tagged.EncodedSize()
	sz += 4 // Fid
	sz += 4 // Newfid
	sz += 2 + uint64(len(v.Name))
	return sz
}

func (v *Txattrwalk) Encode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Encode(b)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Fid)
	if err != nil {
		return err
	}
	err = encodeUint32(b, v.Newfid)
	if err != nil {
		return err
	}
	err = encodeString(b, v.Name)
	if err != nil {
		return err
	}
	return nil
}

func (v *Txattrwalk) Decode(b *bytes.Buffer) error {
	var err error
	err = v.Tagged.Decode(b)
	if err != nil {
		return err
	}
	v.Fid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Newfid, err = decodeUint32(b)
	if err != nil {
		return err
	}
	v.Name, err = decodeString(b)
	if err != nil {
		return err
	}
	return nil
}
