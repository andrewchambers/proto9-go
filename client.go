package proto9

import (
	"errors"
	"fmt"
	"io"
	"sync"
)

var (
	ErrClientClosed    = errors.New("client closed")
	ErrTagAlreadyInUse = errors.New("tag already in use")
	ErrTagsExhausted   = errors.New("tags exhausted")
	ErrFidsExhausted   = errors.New("fids exhausted")
)

type fcallResponse struct {
	fc  Fcall
	err error
}

type Client struct {
	msize   uint32
	version string

	connWriteLock sync.Mutex
	conn          io.ReadWriteCloser

	inflightTagsLock   sync.Mutex
	inflightTags       map[uint16]chan fcallResponse
	inflightTagsClosed bool
	nextTag            uint16

	fidsLock sync.Mutex
	fids     map[uint32]struct{}
	nextFid  uint32
}

func NewClient(conn io.ReadWriteCloser, version string, msize uint32) (*Client, error) {

	c := &Client{
		conn:         conn,
		msize:        msize,
		inflightTags: make(map[uint16]chan fcallResponse),
		fids:         make(map[uint32]struct{}),
	}

	err := c.writeFcall(&Tversion{
		Tagged:  Tagged{Tag: 0xffff},
		Msize:   c.msize,
		Version: version,
	})
	resp, err := c.readFcall()
	if err != nil {
		return nil, err
	}

	rVersion, ok := resp.(*Rversion)
	if !ok || rVersion.Tag != 0xffff {
		return nil, fmt.Errorf("unexpected response from server, expected Rversion with tag 0xFFFF")
	}

	if rVersion.Version != version {
		return nil, fmt.Errorf("protocol version negotiation failed, wanted %q but got %q", version, rVersion.Version)
	}

	if rVersion.Msize > msize || rVersion.Msize < 128 {
		return nil, fmt.Errorf("protocol version negotiation failed, msize %d outside of acceptable range", rVersion.Msize)
	}

	c.msize = rVersion.Msize
	c.version = rVersion.Version

	go c.ReadWorker()
	return c, nil
}

func (c *Client) Msize() uint32 {
	return c.msize
}

func (c *Client) Version() string {
	return c.version
}

func (c *Client) writeFcall(fc Fcall) error {
	c.connWriteLock.Lock()
	defer c.connWriteLock.Unlock()
	return WriteFcall(fc, c.msize, c.conn)
}

func (c *Client) readFcall() (Fcall, error) {
	return ReadFcall(c.msize, c.conn)
}

func (c *Client) hangupInflight(err error) {
	c.inflightTagsLock.Lock()
	defer c.inflightTagsLock.Unlock()
	c.inflightTagsClosed = true
	for _, ch := range c.inflightTags {
		select {
		case ch <- fcallResponse{err: err}:
		default:
		}
	}
}

func (c *Client) ReadWorker() {
	for {
		// XXX integrate buffer pool
		fc, err := c.readFcall()
		if err != nil {
			c.hangupInflight(err)
			return
		}
		c.inflightTagsLock.Lock()
		tag := fc.GetTag()
		respChan, hasChan := c.inflightTags[tag]
		delete(c.inflightTags, tag)
		c.inflightTagsLock.Unlock()
		if hasChan {
			respChan <- fcallResponse{fc: fc}
		}
	}
}

func (c *Client) acquireTag() (uint16, chan fcallResponse, error) {
	c.inflightTagsLock.Lock()
	defer c.inflightTagsLock.Unlock()

	if c.inflightTagsClosed {
		return 0xFFFF, nil, ErrClientClosed
	}

	// No free tags, after taking into account the reserved tag.
	if len(c.inflightTags) >= (0xFFFF - 1) {
		return 0xFFFF, nil, ErrTagsExhausted
	}

	for {
		_, hasTag := c.inflightTags[c.nextTag]
		if !hasTag {
			ch := make(chan fcallResponse, 1)
			c.inflightTags[c.nextTag] = ch
			return c.nextTag, ch, nil
		}
		c.nextTag += 1
		// This tag is reserved.
		if c.nextTag == 0xFFFF {
			c.nextTag = 0
		}
	}

}

func (c *Client) releaseTag(tag uint16) {
	c.inflightTagsLock.Lock()
	defer c.inflightTagsLock.Unlock()
	delete(c.inflightTags, tag)
}

func (c *Client) AcquireFid() (uint32, error) {
	c.fidsLock.Lock()
	defer c.fidsLock.Unlock()

	if len(c.fids) >= 0xFFFFFFFF {
		return 0xFFFFFFFF, ErrFidsExhausted
	}

	for {
		_, hasFid := c.fids[c.nextFid]
		if !hasFid {
			c.fids[c.nextFid] = struct{}{}
			return c.nextFid, nil
		}
		c.nextFid += 1
		// This fid is reserved.
		if c.nextFid == 0xFFFFFFFF {
			c.nextFid = 0
		}
	}

}

func (c *Client) ReleaseFid(fid uint32) {
	c.fidsLock.Lock()
	defer c.fidsLock.Unlock()
	delete(c.fids, fid)
}

func (c *Client) Fcall(fc Fcall) (Fcall, error) {
	tag, ch, err := c.acquireTag()
	if err != nil {
		return nil, err
	}
	defer c.releaseTag(tag)

	fc.SetTag(tag)

	err = c.writeFcall(fc)
	if err != nil {
		return nil, err
	}

	resp := <-ch
	return resp.fc, resp.err
}

func (c *Client) Close() error {
	_ = c.conn.Close()
	c.hangupInflight(ErrClientClosed)
	return nil
}

func AttachDotL(c *Client, aname string, uname string) (*DotLFile, error) {

	if c.Version() != "9P2000.L" {
		return nil, fmt.Errorf("cannot attach to a 9P2000.L mount, protocol version %q", c.Version())
	}

	fid, err := c.AcquireFid()
	if err != nil {
		return nil, err
	}
	defer c.ReleaseFid(fid)

	fc, err := c.Fcall(&Tattach{
		Fid:   fid,
		Aname: aname,
		Uname: uname,
	})
	if err != nil {
		return nil, err
	}
	switch fc := fc.(type) {
	case *Rattach:
		return &DotLFile{
			c: c,
		}, nil
	case *Rlerror:
		return nil, fc
	default:
		return nil, fmt.Errorf("protocol error, expected Rattach")
	}
}

type DotLFile struct {
	c   *Client
	fid uint32
}

func (f *DotLFile) Close() error {
	defer f.c.ReleaseFid(f.fid)

	fc, err := f.c.Fcall(&Tclunk{
		Fid: f.fid,
	})
	if err != nil {
		return err
	}

	switch fc := fc.(type) {
	case *Rclunk:
		return nil
	case *Rlerror:
		return fc
	default:
		return fmt.Errorf("protocol error, expected Rclunk")
	}
}
