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
)

type fcallResponse struct {
	fc  Fcall
	err error
}

type Client struct {
	Msize   uint32
	Version string

	inflightTagsLock   sync.Mutex
	inflightTags       map[uint16]chan fcallResponse
	inflightTagsClosed bool
	nextTag            uint16

	connWriteLock sync.Mutex
	conn          io.ReadWriteCloser
}

func (c *Client) acquireFreeTag() (uint16, chan fcallResponse, error) {
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

func (c *Client) acquireTag(t uint16) (chan fcallResponse, error) {
	c.inflightTagsLock.Lock()
	defer c.inflightTagsLock.Unlock()

	if c.inflightTagsClosed {
		return nil, ErrClientClosed
	}

	_, hasTag := c.inflightTags[t]
	if hasTag {
		return nil, ErrTagAlreadyInUse
	}

	ch := make(chan fcallResponse, 1)
	c.inflightTags[t] = ch
	return ch, nil
}

func (c *Client) releaseTag(tag uint16) {
	c.inflightTagsLock.Lock()
	defer c.inflightTagsLock.Unlock()
	delete(c.inflightTags, tag)
}

func (c *Client) writeFcall(fc Fcall) error {
	c.connWriteLock.Lock()
	defer c.connWriteLock.Unlock()
	return WriteFcall(fc, c.Msize, c.conn)
}

func (c *Client) FcallWithTag(fc Fcall, tag uint16) (Fcall, error) {
	ch, err := c.acquireTag(tag)
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

func (c *Client) Fcall(fc Fcall) (Fcall, error) {
	tag, ch, err := c.acquireFreeTag()
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

func (c *Client) Negotiate(version string, msize uint32) error {
	resp, err := c.FcallWithTag(&Tversion{
		Msize:   msize,
		Version: version,
	}, 0xffff)
	if err != nil {
		return err
	}
	switch resp := resp.(type) {
	case *Rversion:
		c.Msize = resp.Msize
		c.Version = resp.Version
	default:
		return fmt.Errorf("unexpected response from server, expected Rversion")
	}

	if c.Version != version {
		return fmt.Errorf("protocol version negotiation failed, wanted %q but got %q", version, c.Version)
	}

	if c.Msize > msize || c.Msize < 128 {
		return fmt.Errorf("protocol version negotiation failed, msize %d outside of acceptable range", c.Msize)
	}

	return nil
}

func (c *Client) Close() error {

	_ = c.conn.Close()

	c.inflightTagsLock.Lock()
	defer c.inflightTagsLock.Unlock()
	c.inflightTagsClosed = true
	for _, ch := range c.inflightTags {
		select {
		case ch <- fcallResponse{err: ErrClientClosed}:
		default:
		}
	}

	return nil
}
